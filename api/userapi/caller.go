package userapi

import (
	"bytes"
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	"yule/db/usrdb"
)

//注册
func DoRegister(u *usrdb.Usr, login int) (util.Map, error) {
	_, res, err := util.HPostN2(
		fmt.Sprintf("%v/pub/api/register?login=%v", SrvAddr(), login),
		"application/json", bytes.NewBufferString(util.S2Json(u)),
	)
	log.D("res: %v", util.S2Json(res))
	if err != nil {
		return nil, err
	}
	if res.IntVal("code") == 0 {
		return res.MapVal("data"), err
	} else {

		return nil, util.Err("DoLogin by error->%v, %v", err, util.S2Json(res))
	}
}

//用户登录
func DoLogin(account, pwd string) (util.Map, error) {
	res, err := util.HGet2("%v/pub/api/login?account=%v&pwd=%v", SrvAddr(), account, pwd)
	if err != nil {
		return nil, err
	}
	if usrdb.ShowLog {
		log.D("DoLogin res: %v", util.S2Json(res))
	}
	if res.IntVal("code") == 0 {
		return res.MapVal("data"), err
	} else {
		return nil, util.Err("DoLogin by error->%v, %v", err, util.S2Json(res))
	}
}

//获取用户信息
func DoGetUser(token string) (util.Map, error) {
	res, err := util.HGet2("%v/usr/api/getUser?token=%v", SrvAddr(), token)
	if err != nil {
		return nil, err
	}
	if usrdb.ShowLog {
		log.D("DoGetUser res: %v", util.S2Json(res))
	}
	if res.IntVal("code") == 0 {
		return res.MapVal("data"), err
	} else {
		return nil, util.Err("DoGetUser by error->%v, %v", err, util.S2Json(res))
	}
}

//更改用户信息
func DoUpdateUser(token string, u *usrdb.Usr) error {
	_, res, err := util.HPostN2(
		fmt.Sprintf("%v/usr/api/updateUser?token=%v", SrvAddr(), token),
		"application/json", bytes.NewBufferString(util.S2Json(u)),
	)
	if err != nil {
		return err
	}
	if usrdb.ShowLog {
		log.D("DoUpdateUser res: %v", util.S2Json(res))
	}
	if res.IntVal("code") == 0 {
		return nil
	} else {
		return util.Err("update by data(%v) error->%v", util.S2Json(u), util.S2Json(res))
	}
}

//获取所有用户信息
func DoGetUserList(token string, skip, limit int) (util.Map, error) {
	res, err := util.HGet2("%v/usr/api/findUsers?token=%v&skip=%v&limit=%v", SrvAddr(), token, skip, limit)

	if err != nil {
		return nil, err
	}
	if usrdb.ShowLog {
		log.D("DoGetUserList res: %v", util.S2Json(res))
	}
	if res.IntVal("code") == 0 {
		return res.MapVal("data"), err
	} else {

		return nil, util.Err("DoGetUserList by error->%v, %v", err, util.S2Json(res))
	}
}

//关键字搜索用户
func DoSearchUserList(token string, nickKey, phoneKey string, searchMethod, skip, limit, sort int, startTime, overTime int64) (util.Map, error) {
	res, err := util.HGet2("%v/usr/api/searchUsers?token=%v&nickKey=%v&phoneKey=%v"+
		"&searchMethod=%v&skip=%v&limit=%v&sort=%v&startTime=%v&overTime=%v", SrvAddr(), token, nickKey, phoneKey,
		searchMethod, skip, limit, sort, startTime, overTime)
	if err != nil {
		return nil, err
	}
	if usrdb.ShowLog {
		log.D("DoSearchUserList res: %v", util.S2Json(res))
	}
	if res.IntVal("code") == 0 {
		return res.MapVal("data"), err
	} else {

		return nil, util.Err("DoSearchUserList by error->%v, %v", err, util.S2Json(res))
	}
}

//注销
func DoLogout(token string) error {
	res, err := util.HGet2("%v/usr/api/logout?token=%v", SrvAddr(), token)
	if err != nil {
		return err
	}
	if res.IntVal("code") == 0 {
		return nil
	} else {
		return util.Err("logout  error->%v", util.S2Json(res))
	}
}

//绑定手机号码
func DoBindPhone(token, phone, phoneOld string, types, vcode int) (util.Map, error) {
	var res util.Map
	var err error
	if token != "" {
		res, err = util.HGet2("%v/usr/api/bindPhone?token=%v&phone=%v&phoneOld=%v&types=%v&pcode=%v", SrvAddr(), token, phone, phoneOld, types, vcode)
	} else {
		res, err = util.HGet2("%v/pub/api/bindPhone?token=%v&phone=%v&phoneOld=%v&types=%v&pcode=%v", SrvAddr(), token, phone, phoneOld, types, vcode)
	}
	if err != nil {
		return nil, err
	}
	if usrdb.ShowLog {
		log.D("DoBindPhone res: %v", util.S2Json(res))
	}
	if res.IntVal("code") == 0 {
		return res.MapVal("data"), err
	} else {
		return nil, util.Err("DoBindPhone by error->%v, %v", err, util.S2Json(res))
	}
}

//发送短信
func DoSendMsg(token, phone string, types, vcode, mark, sign int) (util.Map, error) {
	res, err := util.HGet2("%v/usr/api/sendMessage?token=%v&phone=%v&types=%v&vcode=%v&mark=%v&sign=%v", SrvAddr(), token, phone, types, vcode, mark, sign)
	if err != nil {
		return nil, err
	}
	if usrdb.ShowLog {
		log.D("DoSendMsg res: %v", util.S2Json(res))
	}
	if res.IntVal("code") == 0 {
		return res.MapVal("data"), err
	} else {

		return nil, util.Err("DoSendMsg by error->%v, %v", err, util.S2Json(res))
	}
}

//邮件
//绑定邮箱号码
func DoBindEmail(token, email, emailOld string, types, vcode int) (util.Map, error) {
	var res util.Map
	var err error
	if token != "" {
		res, err = util.HGet2("%v/usr/api/bindEmail?token=%v&email=%v&emailOld=%v&types=%v&ecode=%v", SrvAddr(), token, email, emailOld, types, vcode)
	} else {
		res, err = util.HGet2("%v/pub/api/bindEmail?token=%v&email=%v&emailOld=%v&types=%v&ecode=%v", SrvAddr(), token, email, emailOld, types, vcode)
	}
	if err != nil {
		return nil, err
	}
	if usrdb.ShowLog {
		log.D("DoBindEmail res: %v", util.S2Json(res))
	}
	if res.IntVal("code") == 0 {
		return res.MapVal("data"), err
	} else {
		return nil, util.Err("DoBindEmail by error->%v, %v", err, util.S2Json(res))
	}
}

//发送短信
func DoSendEmail(token, email string, types, vcode, mark, sign int) (util.Map, error) {
	res, err := util.HGet2("%v/usr/api/sendEmail?token=%v&email=%v&types=%v&vcode=%v&mark=%v&sign=%v", SrvAddr(), token, email, types, vcode, mark, sign)
	if err != nil {
		return nil, err
	}
	if usrdb.ShowLog {
		log.D("DoSendEmail res: %v", util.S2Json(res))
	}
	if res.IntVal("code") == 0 {
		return res.MapVal("data"), err
	} else {

		return nil, util.Err("DoSendEmail by error->%v, %v", err, util.S2Json(res))
	}
}
