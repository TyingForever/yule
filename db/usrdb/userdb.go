package usrdb

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"yule/db"
)

//defined error for user alread exist
var UserExistErr = util.Err("User Existed")

//defined error for user not found
var UserNotFound = util.Err("User Not Found")

//defined error for admin not found
var AdminNotFound = util.Err("Admin Not Found")

//adding user by User struct, see User comment for detail.
func AddUser(u *Usr) error {
	var err error
	_, uid, err := NewUid()
	if err != nil {
		return err
	}
	u.Id = uid
	return AddUserV(u)
}

func AddUserV(u *Usr) error {
	if len(u.Pwd) < 1 {
		return util.Err("AddUser error->pwd field is required, but empty")
	}

	InitUser(u)
	updated, err := db.C(CN_USER).Find(bson.M{"account": u.Account}).Select(bson.M{"pwd": 0}).Apply(mgo.Change{
		Update: bson.M{
			"$setOnInsert": u,
		},
		ReturnNew: true,
		Upsert:    true,
	}, u)
	if err != nil {
		err = util.Err("AddUser add user error(%v) by %v", err, util.S2Json(u))
		log.E("%v", err)
		return err
	}
	if updated.UpsertedId == nil {
		return UserExistErr
	}
	log.D("AddUser add user success by %v", util.S2Json(u))
	u.Pwd = ""
	return err
}

//通过账号密码获取用户
func FindUserByAccountPwd(account, pwd string) (*Usr, error) {
	selector := bson.M{"account": 1, "type": 1, "role": 1}
	return FindUserByUsrPwdV(account, pwd, selector)
}

func FindUserByUsrPwdV(account, pwd string, selector bson.M) (*Usr, error) {
	_, usrs, err := ListUserV(&Usr{
		Account: account,
		Pwd:     pwd,
	}, true, false, selector, 0, 1, 1, 1, 0, 0)

	if err != nil {
		return nil, err
	}
	if len(usrs) < 1 {
		return nil, UserNotFound
	} else {
		return usrs[0], nil
	}
}

//获取用户信息
func GetUserInfo(uid string) (*Usr, error) {
	//selector := bson.M{"pwd": 0,"score":bson.M{"$meta":"textScore"}}
	selector := bson.M{"pwd": 0}
	return GetUserInfoByUid(uid, selector)
}

func GetUserInfoByUid(uid string, selector bson.M) (*Usr, error) {
	_, usrs, err := ListUserV(&Usr{
		Id: uid,
	}, true, false, selector, 0, 1, 1, 1, 0, 0)

	if err != nil {
		return nil, err
	}
	if len(usrs) < 1 {
		return nil, UserNotFound
	} else {
		return usrs[0], nil
	}
}

func FindUsers(query, selector bson.M) (users []*Usr, err error) {
	err = db.C(CN_USER).Find(query).Select(selector).All(&users)
	return
}

//list user by query and selector from mongodb. return total and user list.
//	u: the query fields, all filed in User struct is supported.
//	usr_in: query User.usr field as $in or $all mode. true is $in, false is $all
//	do_total: query total for matched count
//	selector: the result field selector
//	skip: the skip count for result paged.
//	limit: the limit count for result paged.
func ListUserV(u *Usr, usr_in, do_total bool, selector bson.M, skip, limit, searchMethod, sortMethod int, startTime, overTime int64) (int, []*Usr, error) {
	and := []bson.M{}
	if len(u.Id) > 0 {
		and = append(and, bson.M{"_id": u.Id})
	}
	if len(u.Account) > 0 {
		and = append(and, bson.M{
			"$or": []bson.M{
				bson.M{"account": u.Account},
				bson.M{"phone": u.Account},
				bson.M{"email": u.Account},
			},
		})
	}
	if len(u.Pwd) > 0 {
		and = append(and, bson.M{
			"$or": []bson.M{
				bson.M{"pwd": Sha(u.Pwd)},
				bson.M{"pwd": Sha_err(u.Pwd)},
				bson.M{"pwd": Md5(u.Pwd)},
			},
		})
	}

	//if len(u.Email) > 0 {
	//	and = append(and, bson.M{"email":u.Email,"attrs.privated.email":u.Email,"$inc":bson.M{"size":1}})
	//}
	//if len(u.Phone) > 0 {
	//	and = append(and, bson.M{"phone":u.Phone,"attrs.privated.phone":u.Phone,"$inc":bson.M{"size":1}})
	//}
	switch searchMethod {
	case PHONE_SEARCH: //电话
		if len(u.Phone) > 0 {
			and = append(and, bson.M{"phone": bson.M{"$regex": u.Phone}})
		}
	case NICK_SEARCH: //昵称
		if len(u.Attrs.StrVal("nickname")) > 0 {
			log.D("ListUserV nickname(%v)", u.Attrs.StrVal("nickname"))
			and = append(and, bson.M{"attrs.nickname": bson.M{"$regex": u.Attrs.StrVal("nickname"), "$options": "$i"}})
		}
	case TIME_SEARCH: //注册时间段
		if startTime > 0 && overTime > 0 {
			and = append(and, bson.M{"time": bson.M{"$gte": startTime, "$lte": overTime}})
		}
	}

	if u.Status != 0 {
		and = append(and, bson.M{"status": u.Status})
	}
	//if u.Last > 0 {
	//	and = append(and, bson.M{"last": bson.M{"$gte": u.Last}})
	//} else if u.Last < 0 {
	//	and = append(and, bson.M{"last": bson.M{"$lte": -u.Last}})
	//}
	//if u.Time > 0 {
	//	and = append(and, bson.M{"time": bson.M{"$gte": u.Time}})
	//} else if u.Time < 0 {
	//	and = append(and, bson.M{"time": bson.M{"$lte": -u.Time}})
	//}

	if u.Role == USR_ORDINARY {
		and = append(and, bson.M{"role": USR_ORDINARY})
	}
	if len(and) < 1 {
		return 0, nil, util.Err("at last one arguments must be setted on User struct, but User(%v)", util.S2Json(u))
	}

	var usr []*Usr
	var fargs = bson.M{"$and": and}
	//if ShowLog {
	//	log.D("ListUserV args(%v)", util.S2Json(fargs))
	//}
	var Q = db.C(CN_USER).Find(fargs).Select(selector)
	if skip > 0 {
		Q = Q.Skip(skip)
	}
	if limit > 0 {
		Q = Q.Limit(limit)
	}
	switch sortMethod {
	case REGISTER_SORT:
		Q.Sort("-time") //注册时间排序--降序
	case UPDATE_SORT:
		Q.Sort("-last") //最后更改时间排序--降序
	case ACCOUNT_SIZE_SORT:
		Q.Sort("-size") //可登陆账号数目--降序
	case MATCH_SORT:
		//Q.Sort("")
	}
	var err = Q.All(&usr)
	if err != nil {
		err = util.Err("ListUserV list user error(%v) by selector(%v),args(%v)", err, util.S2Json(selector), util.S2Json(fargs))
		log.E("%v", err)
		return 0, nil, err
	}
	if ShowLog {
		log.D("ListUserV list user(%v found) by selector(%v),args(%v)", len(usr), util.S2Json(selector), util.S2Json(fargs))
	}
	var total int = 0
	if do_total {
		total, err = db.C(CN_USER).Find(fargs).Select(selector).Skip(skip).Limit(limit * 10).Count()
		if err != nil {
			err = util.Err("ListUserV count user error(%v) by selector(%v),args(%v)", err, util.S2Json(selector), util.S2Json(fargs))
			log.E("%v", err)
			return 0, nil, err
		} else if ShowLog {
			log.D("ListUserV count user total(%v) by selector(%v),args(%v)", total, util.S2Json(selector), util.S2Json(fargs))
		}
	}
	total += skip
	return total, usr, err
}

//修改用户信息
func UpdateUserInfo(u *Usr) error {
	if u.Id == "" {
		return util.Err("the user id is emtpy")
	}
	var args = bson.M{}
	//var usrs = bson.M{}

	if len(u.Account) > 0 {
		_, us, err := ListUserV(&Usr{Account: u.Account}, true, false, bson.M{"status": 1}, 0, 1, 1, 1, 0, 1)
		if err != nil {
			return err
		}
		if len(us) > 0 {
			return UserExistErr
		}
	}

	//if len(u.Pwd) > 0 {
	//	and = append(and, bson.M{
	//		"$or": []bson.M{
	//			bson.M{"pwd": Sha(u.Pwd)},
	//			bson.M{"pwd": Sha_err(u.Pwd)},
	//			bson.M{"pwd": Md5(u.Pwd)},
	//		},
	//	})
	//}
	//
	//if len(u.Email) > 0 {
	//	and = append(and, bson.M{"email":u.Email,"attrs.privated.email":u.Email,"$inc":bson.M{"size":1}})
	//}
	//if len(u.Phone) > 0 {
	//	and = append(and, bson.M{"phone":u.Phone,"attrs.privated.phone":u.Phone,"$inc":bson.M{"size":1}})
	//}

	if len(u.Phone) > 0 {
		args["phone"] = u.Phone
	}
	if len(u.Email) > 0 {
		args["email"] = u.Email
	}

	if len(u.Pwd) > 0 {
		args["pwd"] = Sha(u.Pwd)
	}

	if len(u.Attrs) > 0 {
		for key, val := range u.Attrs {
			args["attrs."+key] = val
		}
	}

	if u.Status != 0 {
		args["status"] = u.Status
	}
	args["last"] = util.Now()
	var fargs = bson.M{"$set": args}
	//if len(usrs) > 0 {
	//	fargs["$addToSet"] = usrs
	//}

	if ShowLog {
		log.D("updateUser: %v", util.S2Json(u))
	}
	err := db.C(CN_USER).Update(bson.M{"_id": u.Id}, fargs)
	if err != nil {
		log.E("UpdateUser update user error(%v) by id(%v) to %v", err, u.Id, util.S2Json(fargs))
		return err
	}
	if ShowLog {
		log.D("UpdateUser update user success by id(%v) to %v ", u.Id, util.S2Json(fargs))
	}
	return err
}

//获取所有的用户信息
func ListUsersOrdinary(uid string, skip, limit int) (int, []*Usr, error) {
	if uid == "" {
		return 0, nil, util.Err("the user id is emtpy")
	}
	_, us, err := ListUserV(&Usr{Id: uid}, true, false, bson.M{"account": 1, "role": 1}, 0, 1, 1, 1, 0, 0)
	if err != nil {
		return 0, nil, err
	}
	if len(us) < 1 {
		return 0, nil, UserNotFound
	}
	if us[0].Role != USR_ADMIN {
		return 0, nil, AdminNotFound
	}
	return ListUserV(&Usr{Role: USR_ORDINARY}, true, true, bson.M{"pwd": 0, "attr": 0}, skip, limit, 1, 1, 0, 0)
}

//搜索的用户信息
func SearchUsersOrdinary(uid, nickKey, phoneKey string, searchMethod, sort, skip, limit int, startTime, overTime int64) (int, []*Usr, error) {
	if uid == "" {
		return 0, nil, util.Err("the user id is emtpy")
	}
	_, us, err := ListUserV(&Usr{Id: uid}, true, false, bson.M{"account": 1, "role": 1}, 0, 1, 1, 1, 0, 0)
	if err != nil {
		return 0, nil, err
	}
	if len(us) < 1 {
		return 0, nil, UserNotFound
	}
	if us[0].Role != USR_ADMIN {
		return 0, nil, AdminNotFound
	}
	attrs := util.Map{
		"nickname": nickKey,
	}
	if searchMethod == NICK_SEARCH && sort == MATCH_SORT {
		return PipeUserV(&Usr{Role: USR_ORDINARY, Phone: phoneKey, Attrs: attrs}, true, true, bson.M{"pwd": 0, "attr": 0}, skip, limit, searchMethod, sort, startTime, overTime)
	}

	return ListUserV(&Usr{Role: USR_ORDINARY, Phone: phoneKey, Attrs: attrs}, true, true, bson.M{"pwd": 0, "attr": 0}, skip, limit, searchMethod, sort, startTime, overTime)
}

//管道查询
func PipeUserV(u *Usr, usr_in, do_total bool, selector bson.M, skip, limit, searchMethod, sortMethod int, startTime, overTime int64) (int, []*Usr, error) {
	and := []bson.M{}

	if len(u.Attrs.StrVal("nickname")) > 0 {
		log.D("ListUserV nickname(%v)", u.Attrs.StrVal("nickname"))
		and = append(and, bson.M{"attrs.nickname": bson.M{"$regex": u.Attrs.StrVal("nickname"), "$options": "$i"}})
	}

	if u.Role == USR_ORDINARY {
		and = append(and, bson.M{"role": USR_ORDINARY})
	}
	if len(and) < 1 {
		return 0, nil, util.Err("at last one arguments must be setted on User struct, but User(%v)", util.S2Json(u))
	}

	var usr []*Usr
	var fargs = bson.M{"$and": and}
	var o1 = bson.M{"$match": bson.M{"role": USR_ORDINARY, "attrs.nickname": bson.M{"$regex": u.Attrs.StrVal("nickname"), "$options": "$i"}}}
	var o2 = bson.M{"$project": bson.M{
		"uid":     1,
		"account": 1,
		"phone":   1,
		"email":   1,
		"attrs":   1,
		"last":    1,
		"time":    1,
		"role":    1,
		"nickLen": bson.M{"$strLenBytes": "$attrs.nickname"},
	}}
	var o3 = bson.M{
		"$sort": bson.M{
			"nickLen": 1,
		},
	}
	var o4 = bson.M{
		"$skip": skip,
	}
	var o5 = bson.M{
		"$limit": limit,
	}

	operations := []bson.M{o1, o2, o3, o4, o5}
	var err = db.C(CN_USER).Pipe(operations).All(&usr)
	if err != nil {
		err = util.Err("ListUserV list user error(%v) by project(%v),args(%v)", err, util.S2Json(o2), util.S2Json(o1))
		log.E("%v", err)
		return 0, nil, err
	}
	if ShowLog {
		log.D("ListUserV list user(%v found) by selector(%v),args(%v)", len(usr), util.S2Json(o2), util.S2Json(o1))
	}
	var total int = 0
	if do_total {
		total, err = db.C(CN_USER).Find(fargs).Select(selector).Skip(skip).Limit(limit * 10).Count()
		if err != nil {
			err = util.Err("ListUserV count user error(%v) by selector(%v),args(%v)", err, util.S2Json(selector), util.S2Json(fargs))
			log.E("%v", err)
			return 0, nil, err
		} else if ShowLog {
			log.D("ListUserV count user total(%v) by selector(%v),args(%v)", total, util.S2Json(selector), util.S2Json(fargs))
		}
	}
	total += skip
	return total, usr, err
}

func InitUser(u *Usr) {
	u.Pwd = Sha(u.Pwd)
	u.Status = USR_S_N
	u.Last = util.Now()
	u.Time = u.Last
	u.Role = USR_ORDINARY
}
