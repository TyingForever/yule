package usrdb

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"regexp"
	"strings"
	"yule/db"
)

func UpsertCode(send, category string, code int, validTime int64, Type int) (int, error) {
	_, err := db.C(CN_VERIFICATION_CODE).Find(bson.M{"_id": send,
		"$or": []bson.M{
			{"status": PCS_DEL},
			{"type": bson.M{"$ne": Type}},
			{"time": bson.M{"$lte": util.Now() - validTime}}}, //数据库中的时间<=当前时间-有限时间，即失效
	}).Apply(mgo.Change{
		Update: bson.M{
			"$set": bson.M{
				"_id":      send,
				"category": category,
				"time":     util.Now(),
				"code":     code,
				"type":     Type,
				"status":   PCS_NORMAL,
			},
		},
		Upsert: true,
	}, nil)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			log.W("UpsertCode send(%v) time(%v) err(%v)", send, validTime, err)
			return 7, util.Err("验证码已发送，请稍后再试")
		}
		return 1, err
	}
	return 0, nil
}

func VerifyCode(send, category string, code int, validTime int64) (int, error) {
	var verificationCode VerificationCode
	_, err := db.C(CN_VERIFICATION_CODE).Find(
		bson.M{"_id": send, "category": category, "code": code, "time": bson.M{"$gte": util.Now() - validTime}},
	).Apply(mgo.Change{
		Update:    bson.M{"$set": bson.M{"status": PCS_DEL}},
		Upsert:    false,
		ReturnNew: false,
	}, &verificationCode)
	if err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			//return 9, util.Err("验证码未发送或已过期，请重新申请发送")
			return 9, util.Err("验证码错误")
		}
		return 1, err
	}
	//if phoneCode.Code != code {
	//	return 5, util.Err("验证码错误")
	//}
	if verificationCode.Status != PCS_NORMAL {
		return 5, util.Err("验证码错误")
	}
	return 0, nil
}

//改绑定手机
func ChangeBindPhone(uid, phone, phoneOld string) (int, string, error) {
	if !CheckPhoneFormat(phone, phoneOld) {
		return 1, "手机格式有误", util.Err("phone(%v) or phoneOld(%v) invalid,please input correct phone", phone, phoneOld)
	}
	if phone == phoneOld {
		return 1, "新旧手机不能一样", util.Err("phone(%v) is equal to  phoneOld(%v)", phone, phoneOld)
	}
	err := db.C(CN_USER).Update(
		bson.M{"_id": uid, "attrs.privated.phone": phoneOld, "phone": phoneOld},
		bson.M{
			"$unset": bson.M{"attrs.privated.phonePending": ""},
			"$set":   bson.M{"attrs.privated.phone": phone, "phone": phone, "last": util.Now()},
		})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return 4, "手机已被绑定", util.Err("phone(%v) had been bind", phone)
		}
		if strings.Contains(err.Error(), "not found") {
			return 3, "用户已绑定手机或旧绑定手机填写错误", util.Err("uid(%v) had bind phone", uid)
		}
		log.E("Change Bind phone Update user's phone(%v) by uid(%v)  error(%v)", phone, uid, err)
		return 2, "服务器错误", err
	}
	return 0, "", nil
}

//绑定手机号码
func BindPhone(uid, phone string) (int, string, error) {
	if !CheckPhoneFormat(phone) {
		return 1, "手机格式有误", util.Err("phone(%v)  invalid,please input correct phone", phone)
	}
	//$exists 只查询存在该字段的
	err := db.C(CN_USER).Update(
		bson.M{"_id": uid, "attrs.privated.phone": bson.M{"$exists": false}},
		bson.M{
			"$set": bson.M{"attrs.privated.phone": phone, "last": util.Now(), "phone": phone},
			"$inc": bson.M{"size": 1},
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return 4, "手机已被绑定", util.Err("phone(%v) had been bind", phone)
		}
		if strings.Contains(err.Error(), "not found") {
			return 3, "用户已绑定手机", util.Err("uid(%v) had bind phone", uid)
		}
		log.E("BindPhone Update user's phone(%v) by uid(%v) error(%v)", phone, uid, err)
		return 2, "服务器错误", err
	}
	return 0, "", nil
}

//解绑手机
func UnBindPhone(uid, phone string) (int, string, error) {
	if !CheckPhoneFormat(phone) {
		return 1, "手机格式有误", util.Err("phone(%v)  invalid,please input correct phone", phone)
	}
	err := db.C(CN_USER).Update(
		bson.M{"_id": uid, "phone": phone},
		bson.M{
			"$unset": bson.M{"attrs.privated.phonePending": "", "attrs.privated.phone": "", "phone": ""},
			"$set":   bson.M{"last": util.Now()},
			"$inc":   bson.M{"size": -1},
		},
	)
	if err != nil {
		log.E("BindPhone Update user's phone(%v) by uid(%v) error(%v)", phone, uid, err)
		return 2, "服务器错误", err
	}
	return 0, "", nil
}

//email
//改绑定邮箱
func ChangeBindEmail(uid, email, emailOld string) (int, string, error) {
	if !CheckEmailFormat(email, emailOld) {
		return 1, "邮箱格式有误", util.Err("email(%v) or emailOld(%v) invalid,please input correct email", email, emailOld)
	}
	if email == emailOld {
		return 1, "新旧邮箱不能一样", util.Err("email(%v) is equal to  emailOld(%v)", email, emailOld)
	}
	err := db.C(CN_USER).Update(
		bson.M{"_id": uid, "attrs.privated.email": emailOld, "email": emailOld},
		bson.M{
			"$unset": bson.M{"attrs.privated.emailPending": ""},
			"$set":   bson.M{"attrs.privated.email": email, "email": email, "last": util.Now()},
		})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return 4, "邮箱已被绑定", util.Err("email(%v) had been bind", email)
		}
		if strings.Contains(err.Error(), "not found") {
			return 3, "用户已绑定邮箱或旧绑定邮箱填写错误", util.Err("uid(%v) had bind email", uid)
		}
		log.E("Change Bind email Update user's email(%v) by uid(%v)  error(%v)", email, uid, err)
		return 2, "服务器错误", err
	}
	return 0, "", nil
}

//绑定邮箱号码
func BindEmail(uid, email string) (int, string, error) {

	if !CheckEmailFormat(email) {
		return 1, "邮箱格式有误", util.Err("email(%v)  invalid,please input correct email", email)
	}
	//$exists 只查询存在该字段的
	err := db.C(CN_USER).Update(
		bson.M{"_id": uid, "attrs.privated.email": bson.M{"$exists": false}},
		bson.M{
			"$set": bson.M{"attrs.privated.email": email, "last": util.Now(), "email": email},
			"$inc": bson.M{"size": 1},
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return 4, "邮箱已被绑定", util.Err("email(%v) had been bind", email)
		}
		if strings.Contains(err.Error(), "not found") {
			return 3, "用户已绑定邮箱", util.Err("uid(%v) had bind Email", uid)
		}
		log.E("BindEmail Update user's email(%v) by uid(%v) error(%v)", email, uid, err)
		return 2, "服务器错误", err
	}
	return 0, "", nil
}

//解绑邮箱
func UnBindEmail(uid, email string) (int, string, error) {
	if !CheckEmailFormat(email) {
		return 1, "邮箱格式有误", util.Err("email(%v)  invalid,please input correct email", email)
	}
	err := db.C(CN_USER).Update(
		bson.M{"_id": uid, "email": email},
		bson.M{
			"$unset": bson.M{"attrs.privated.emailPending": "", "attrs.privated.email": "", "email": ""},
			"$set":   bson.M{"last": util.Now()},
			"$inc":   bson.M{"size": -1},
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return 3, "用户未绑定邮箱", util.Err("uid(%v) had bind Email", uid)
		}
		log.E("Bindemail Update user's email(%v) by uid(%v) error(%v)", email, uid, err)
		return 2, "服务器错误", err
	}
	return 0, "", nil
}

//检查手机号码的格式
func CheckPhoneFormat(phone ...string) bool {
	reg := regexp.MustCompile("^[\\d]{11}$")
	for _, item := range phone {
		if !reg.MatchString(item) {
			return false
		}
	}
	return true
}

//检查邮箱号码的格式
func CheckEmailFormat(email ...string) bool {
	reg := regexp.MustCompile("\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*")
	for _, item := range email {
		if !reg.MatchString(item) {
			return false
		}
	}
	return true
}

func CreateIndex(key string, do_unique bool) error {
	index := mgo.Index{
		Key:        []string{key},
		Unique:     do_unique,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := db.C(CN_USER).EnsureIndex(index)
	//err := db.C(CN_USER).EnsureIndex(mgo.Index{Key:[]string{"email","phone"},Unique:true,Sparse:true,DropDups:true,Background: true})
	if err != nil {
		log.E("create encureIndex err(%v)", err)
		return err
	}
	return err
}
