package userapi

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2"
	"math/rand"
	"strings"
	"testing"
	"yule/db"
	"yule/db/usrdb"
)

func TestRe(t *testing.T) {
	Remove()
	var u *usrdb.Usr
	u = &usrdb.Usr{
		Account: "1234",
		Pwd:     "123456",
		Type:    1,
	}
	log.D("---------------",util.S2Json(u))
	//
	login := 0
	//
	_, err := DoRegister(u, login)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestRegister(t *testing.T) {
	Remove()
	//注册 不登录
	var u *usrdb.Usr
	u = &usrdb.Usr{
		Account: "1",
		Pwd:     "123",
		Type:    1,
	}
	//
	login := 0
	//
	u.Account = "1"
	res, err := DoRegister(u, login)
	if err != nil {
		t.Error(err)
		return
	}
	res_user := res.MapVal("usr")
	//log.D("user: %v ", util.S2Json(res.MapVal("usr")))
	if u.Account != res_user.StrVal("account") || u.Type != (int)(res_user.IntVal("type")) {
		t.Error("data err")
		return
	}

	//err exist
	u.Account = "1"
	_, err = DoRegister(u, login)
	if err == nil {
		t.Error(err)
		return
	}

	//注册并登录
	u.Account = "2"
	login = 1
	res, err = DoRegister(u, login)
	if err != nil {
		t.Error(err)
		return
	}
	res_user = res.MapVal("usr")
	token := res.StrVal("token")
	if u.Account != res_user.StrVal("account") || res_user.IntVal("role") != 2 || u.Type != (int)(res_user.IntVal("type")) || token == "" {
		t.Error("data err")
		return
	}

	//更改用户信息
	token = res.StrVal("token")
	update_user := &usrdb.Usr{}
	update_user.Attrs = util.Map{
		"age":        18,
		"birthday":   util.Now(),
		"location":   "GuangDong",
		"hometown":   "MaoMing",
		"profession": "Programmer",
		"nickname":   "King",
	}
	err = DoUpdateUser(token, update_user)
	if err != nil {
		t.Error(err)
		return
	}

	//获取用户信息
	token = res.StrVal("token")
	res, err = DoGetUser(token)
	if err != nil {
		t.Error(err)
		return
	}
	res_user = res.MapVal("user")
	if u.Account != res_user.StrVal("account") || res_user.IntVal("role") != 2 ||
		u.Type != (int)(res_user.IntVal("type")) || token == "" ||
		update_user.Attrs.IntVal("age") != res_user.MapVal("attrs").IntVal("age") ||
		update_user.Attrs.IntVal("birthday") != res_user.MapVal("attrs").IntVal("birthday") ||
		update_user.Attrs.StrVal("location") != res_user.MapVal("attrs").StrVal("location") ||
		update_user.Attrs.StrVal("hometown") != res_user.MapVal("attrs").StrVal("hometown") ||
		update_user.Attrs.StrVal("nickname") != res_user.MapVal("attrs").StrVal("nickname") {
		t.Errorf("data err %v", util.S2Json(res))
		return
	}

	err = DoLogout(token)
	if err != nil {
		t.Error(err)
		return
	}

	//db err
	mgo.Mock = true
	mgo.SetMckC("Query-Apply", 0)
	u.Account = "a1"
	u.Pwd = "123"
	_, err = DoRegister(u, 0)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()

	mgo.SetMckC("Query-Apply", 0)
	u.Account = "a2"
	u.Pwd = "123"
	_, err = DoRegister(u, 1)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()

	//db update
	mgo.SetMckC("Collection-Update", 0)
	err = DoUpdateUser(token, u)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()

	mgo.SetMckC("Collection-RemoveAll", 0)
	err = DoLogout(token)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()

}

func TestLogin(t *testing.T) {
	Remove()
	account := "100"
	pwd := "123"
	res, err := DoLogin(account, pwd)
	if err != nil {
		t.Error(err)
		return
	}
	log.D("res: %v", util.S2Json(res))
	res_user := res.MapVal("usr")
	log.D("user: %v ", util.S2Json(res.MapVal("usr")))
	token := res.StrVal("token")
	if account != res_user.StrVal("account") || res_user.IntVal("role") != 1 || token == "" {
		t.Error("data err")
		return
	}

	//login err
	account = "1"
	pwd = "1234"
	_, err = DoLogin(account, pwd)
	if err == nil {
		t.Error(err)
		return
	}

}

func TestDoGetUserList(t *testing.T) {
	Remove()
	prepareData()
	var err error

	db.C(usrdb.CN_SESSION).Remove(nil)
	account := "100"
	pwd := "123"
	res, err := DoLogin(account, pwd)
	if err != nil {
		t.Error(err)
		return
	}
	log.D("res: %v", util.S2Json(res))

	token := res.StrVal("token")
	skip := 0
	limit := 10
	res, err = DoGetUserList(token, skip, limit)
	if err != nil {
		t.Error(err)
		return
	}
	if len(res.AryMapVal("users")) != 10 {
		t.Errorf("DoGetUserList err %v", util.S2Json(res))
	}

	t.Logf("res: %v", util.S2Json(res))
}

//手机绑定
func TestBindPhone(t *testing.T) {
	prepareData()
	db.C(usrdb.CN_SESSION).Remove(nil)
	account := "11111"
	pwd := "123"
	res, err := DoLogin(account, pwd)
	if err != nil {
		t.Error(err)
		return
	}
	token := res.StrVal("token")

	//bind
	if usrdb.ShowLog {
		log.D("res: %v", util.S2Json(res))
	}
	res, err = DoSendMsg(token, "13800138004", usrdb.BIND, 0, 0, 0)
	if err != nil {
		t.Error(err)
		return
	}
	//sendMessage err
	//手机格式有误
	_, err = DoSendMsg(token, "138001384", usrdb.BIND, 0, 0, 0)
	if err == nil {
		t.Error(err)
		return
	}
	//参数有误
	_, err = DoSendMsg(token, "13800138004", 0, 0, 0, 0)
	if err == nil {
		t.Error(err)
		return
	}

	if usrdb.ShowLog {
		log.D("res msg: %v", util.S2Json(res))
	}
	res, err = DoBindPhone(token, res.MapVal("phoneCode").StrVal("Send"), "", int(res.MapVal("phoneCode").IntVal("Type")), int(res.MapVal("phoneCode").IntVal("Code")))
	if err != nil {
		t.Error(err)
		return
	}
	if usrdb.ShowLog {
		log.D("res msg: %v", util.S2Json(res))
	}
	//err 验证码出错
	_, err = DoBindPhone(token, res.MapVal("phoneCode").StrVal("Send"), "", int(res.MapVal("phoneCode").IntVal("Type")), int(res.MapVal("phoneCode").IntVal("Type")))
	if err == nil {
		t.Error(err)
		return
	}
	//err 已被绑定
	_, err = DoSendMsg(token, "13800138004", usrdb.BIND, 0, 0, 0)
	if err == nil {
		t.Error(err)
		return
	}

	//改绑定
	res, err = DoSendMsg(token, "13800138002", usrdb.MODIFYBIND, 0, 0, 0)
	if err != nil {
		t.Error(err)
		return
	}
	if usrdb.ShowLog {
		log.D("res msg: %v", util.S2Json(res))
	}
	//sendMessage err
	//已存在该手机号
	_, err = DoSendMsg(token, "13800138004", usrdb.MODIFYBIND, 0, 0, 0)
	if err == nil {
		t.Error(err)
		return
	}
	res, err = DoBindPhone(token, res.MapVal("phoneCode").StrVal("Send"), "13800138004", int(res.MapVal("phoneCode").IntVal("Type")), int(res.MapVal("phoneCode").IntVal("Code")))
	if err != nil {
		t.Error(err)
		return
	}
	if usrdb.ShowLog {
		log.D("res msg: %v", util.S2Json(res))
	}
	//err 手机绑定时原手机号有误
	res, err = DoBindPhone(token, res.MapVal("phoneCode").StrVal("Send"), "13800138000", int(res.MapVal("phoneCode").IntVal("Type")), int(res.MapVal("phoneCode").IntVal("Code")))
	if err == nil {
		t.Error(err)
		return
	}
	//新旧手机一样
	res, err = DoBindPhone(token, res.MapVal("phoneCode").StrVal("Send"), res.MapVal("phoneCode").StrVal("Send"), int(res.MapVal("phoneCode").IntVal("Type")), int(res.MapVal("phoneCode").IntVal("Code")))
	if err == nil {
		t.Error(err)
		return
	}

	//unbind
	res, err = DoSendMsg(token, "13800138002", usrdb.UNBIND, 0, 0, 0)
	if err != nil {
		t.Error(err)
		return
	}
	if usrdb.ShowLog {
		log.D("res msg: %v", util.S2Json(res))
	}
	//sendMessage err
	//该用户不存在该手机号
	_, err = DoSendMsg(token, "13800138003", usrdb.UNBIND, 0, 0, 0)
	if err == nil {
		t.Error(err)
		return
	}
	res, err = DoBindPhone(token, res.MapVal("phoneCode").StrVal("Send"), "", int(res.MapVal("phoneCode").IntVal("Type")), int(res.MapVal("phoneCode").IntVal("Code")))
	if err != nil {
		t.Error(err)
		return
	}
	if usrdb.ShowLog {
		log.D("res msg: %v", util.S2Json(res))
	}

	//err

}

func TestBindEmail(t *testing.T) {
	prepareData()
	//usrdb.C(usrdb.CN_SESSION).Remove(nil)
	account := "1111"
	pwd := "123"
	res, err := DoLogin(account, pwd)
	if err != nil {
		t.Error(err)
		return
	}
	token := res.StrVal("token")

	if usrdb.ShowLog {
		log.D("res: %v", util.S2Json(res))
	}
	email1 := "843232800@qq.com"
	email2 := "843232801@qq.com"
	email3 := "843232802@qq.com"
	email_err := "138001384"
	//bind
	if usrdb.ShowLog {
		log.D("res: %v", util.S2Json(res))
	}
	res, err = DoSendEmail(token, email1, usrdb.BIND, 0, 0, 0)
	if err != nil {
		t.Error(err)
		return
	}
	//sendMessage err
	//邮箱格式有误
	_, err = DoSendEmail(token, email_err, usrdb.BIND, 0, 0, 0)
	if err == nil {
		t.Error(err)
		return
	}
	//参数有误
	_, err = DoSendEmail(token, email1, 0, 0, 0, 0)
	if err == nil {
		t.Error(err)
		return
	}

	if usrdb.ShowLog {
		log.D("res msg: %v", util.S2Json(res))
	}
	res, err = DoBindEmail(token, res.MapVal("emailCode").StrVal("Send"), "", int(res.MapVal("emailCode").IntVal("Type")), int(res.MapVal("emailCode").IntVal("Code")))
	if err != nil {
		t.Error(err)
		return
	}
	if usrdb.ShowLog {
		log.D("res msg: %v", util.S2Json(res))
	}
	//err 验证码出错
	_, err = DoBindEmail(token, res.MapVal("emailCode").StrVal("Send"), "", int(res.MapVal("emailCode").IntVal("Type")), int(res.MapVal("emailCode").IntVal("Type")))
	if err == nil {
		t.Error(err)
		return
	}
	//err 已被绑定
	_, err = DoSendEmail(token, email1, usrdb.BIND, 0, 0, 0)
	if err == nil {
		t.Error(err)
		return
	}

	//改绑定
	res, err = DoSendEmail(token, email2, usrdb.MODIFYBIND, 0, 0, 0)
	if err != nil {
		t.Error(err)
		return
	}
	if usrdb.ShowLog {
		log.D("res msg: %v", util.S2Json(res))
	}
	//sendMessage err
	//已存在该邮箱号
	_, err = DoSendEmail(token, email1, usrdb.MODIFYBIND, 0, 0, 0)
	if err == nil {
		t.Error(err)
		return
	}
	res, err = DoBindEmail(token, res.MapVal("emailCode").StrVal("Send"), email1, int(res.MapVal("emailCode").IntVal("Type")), int(res.MapVal("emailCode").IntVal("Code")))
	if err != nil {
		t.Error(err)
		return
	}
	if usrdb.ShowLog {
		log.D("res msg: %v", util.S2Json(res))
	}
	//err 邮箱绑定时原邮箱号有误
	res, err = DoBindEmail(token, res.MapVal("emailCode").StrVal("Send"), email_err, int(res.MapVal("emailCode").IntVal("Type")), int(res.MapVal("emailCode").IntVal("Code")))
	if err == nil {
		t.Error(err)
		return
	}
	//新旧邮箱一样
	res, err = DoBindEmail(token, res.MapVal("emailCode").StrVal("Send"), res.MapVal("emailCode").StrVal("Send"), int(res.MapVal("emailCode").IntVal("Type")), int(res.MapVal("emailCode").IntVal("Code")))
	if err == nil {
		t.Error(err)
		return
	}

	//unbind
	res, err = DoSendEmail(token, email2, usrdb.UNBIND, 0, 0, 0)
	if err != nil {
		t.Error(err)
		return
	}
	if usrdb.ShowLog {
		log.D("res msg: %v", util.S2Json(res))
	}
	//sendMessage err
	//该用户不存在该邮箱号
	_, err = DoSendEmail(token, email3, usrdb.UNBIND, 0, 0, 0)
	if err == nil {
		t.Error(err)
		return
	}
	res, err = DoBindEmail(token, res.MapVal("emailCode").StrVal("Send"), "", int(res.MapVal("emailCode").IntVal("Type")), int(res.MapVal("emailCode").IntVal("Code")))
	if err != nil {
		t.Error(err)
		return
	}
	if usrdb.ShowLog {
		log.D("res msg: %v", util.S2Json(res))
	}

	//email := "843232829@qq.com"
	//res, err = DoSendEmail(token, email, 1, 0, 0, 0)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//if usrdb.ShowLog {
	//	log.D("res msg: %v", util.S2Json(res))
	//}
	//
	//log.D("res %v", res.MapVal("emailCode").IntVal("Code"))
	//res, err = DoBindEmail(token, res.MapVal("emailCode").StrVal("Send"), "", int(res.MapVal("emailCode").IntVal("Type")), int(res.MapVal("emailCode").IntVal("Code")))
	////res,err = DoBindEmail(token,res.MapVal("emailCode").StrVal("Email"),"",int(res.MapVal("emailCode").IntVal("Type")),123)
	//
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//if usrdb.ShowLog {
	//	log.D("res msg: %v", util.S2Json(res))
	//}

}

//关键字搜索
func TestSearchUsers(t *testing.T) {
	var err error
	prepareData()
	db.C(usrdb.CN_SESSION).Remove(nil)
	account := "100"
	pwd := "123"
	res, err := DoLogin(account, pwd)
	if err != nil {
		t.Error(err)
		return
	}
	log.D("res: %v", util.S2Json(res))

	token := res.StrVal("token")
	var searchMethod, skip, limit, sort int
	var startTime, overTime int64
	var nickKey, phoneKey string
	skip = 0
	limit = 5
	searchMethod, sort = usrdb.DEFAULT_SEARCH, usrdb.DEFAULT_SORT
	startTime, overTime = 0, 0
	//1.NICK_SEARCH DEFAULT_SORT
	searchMethod = usrdb.NICK_SEARCH
	nickKey = "q"
	sort = usrdb.DEFAULT_SORT
	res, err = DoSearchUserList(token, nickKey, phoneKey,
		searchMethod, skip, limit, sort, startTime, overTime)
	if err != nil {
		t.Error(err)
		return
	}

	if len(res.AryMapVal("users")) != 5 {
		t.Errorf("err res(1.NICK_SEARCH DEFAULT_SORT): %v", util.S2Json(res))
		return
	}
	for i := 0; i < 5; i++ {
		if !strings.Contains(res.AryMapVal("users")[i].MapVal("attrs").StrVal("nickname"), nickKey) {
			t.Errorf("err res(1.NICK_SEARCH DEFAULT_SORT): %v", util.S2Json(res))
			return
		}
	}
	t.Logf("res(1.NICK_SEARCH DEFAULT_SORT): %v", util.S2Json(res))

	//2.NICK_SEARCH REGISTER_SORT 降序
	searchMethod = usrdb.NICK_SEARCH
	nickKey = "q"
	sort = usrdb.REGISTER_SORT
	res, err = DoSearchUserList(token, nickKey, phoneKey,
		searchMethod, skip, limit, sort, startTime, overTime)
	if err != nil {
		t.Error(err)
		return
	}
	if len(res.AryMapVal("users")) != 5 {
		t.Errorf("err res(2.NICK_SEARCH DEFAULT_SORT): %v", util.S2Json(res))
		return
	}
	for i := 0; i < 5; i++ {
		if !strings.Contains(res.AryMapVal("users")[i].MapVal("attrs").StrVal("nickname"), nickKey) {
			t.Errorf("err res(2.NICK_SEARCH DEFAULT_SORT): %v", util.S2Json(res))
			return
		}
	}
	for i := 0; i < 4; i++ {
		if res.AryMapVal("users")[i].IntVal("time") < res.AryMapVal("users")[i+1].IntVal("time") {
			t.Errorf("err res(2.NICK_SEARCH REGISTER_SORT 降序): %v", util.S2Json(res))
			return
		}
	}

	t.Logf("res(2.NICK_SEARCH REGISTER_SORT): %v", util.S2Json(res))

	//3.NICK_SEARCH UPDATE_SORT
	searchMethod = usrdb.NICK_SEARCH
	nickKey = "q"
	sort = usrdb.UPDATE_SORT
	res, err = DoSearchUserList(token, nickKey, phoneKey,
		searchMethod, skip, limit, sort, startTime, overTime)
	if err != nil {
		t.Error(err)
		return
	}
	if len(res.AryMapVal("users")) != 5 {
		t.Errorf("err res(3.NICK_SEARCH UPDATE_SORT): %v", util.S2Json(res))
		return
	}
	for i := 0; i < 5; i++ {
		if !strings.Contains(res.AryMapVal("users")[i].MapVal("attrs").StrVal("nickname"), nickKey) {
			t.Errorf("err res(3.NICK_SEARCH UPDATE_SORT): %v", util.S2Json(res))
			return
		}
	}
	for i := 0; i < 4; i++ {
		if res.AryMapVal("users")[i].IntVal("last") < res.AryMapVal("users")[i+1].IntVal("last") {
			t.Errorf("err res(3.NICK_SEARCH UPDATE_SORT 降序): %v", util.S2Json(res))
			return
		}
	}

	t.Logf("res(3.NICK_SEARCH UPDATE_SORT): %v", util.S2Json(res))

	//4.NICK_SEARCH ACCOUNT_SIZE_SORT
	searchMethod = usrdb.NICK_SEARCH
	nickKey = "q"
	sort = usrdb.ACCOUNT_SIZE_SORT
	res, err = DoSearchUserList(token, nickKey, phoneKey,
		searchMethod, skip, limit, sort, startTime, overTime)
	if err != nil {
		t.Error(err)
		return
	}
	if len(res.AryMapVal("users")) != 5 {
		t.Errorf("err res(4.NICK_SEARCH ACCOUNT_SIZE_SORT): %v", util.S2Json(res))
		return
	}
	for i := 0; i < 5; i++ {
		if !strings.Contains(res.AryMapVal("users")[i].MapVal("attrs").StrVal("nickname"), nickKey) {
			t.Errorf("err res(4.NICK_SEARCH ACCOUNT_SIZE_SORT): %v", util.S2Json(res))
			return
		}
	}
	for i := 0; i < 4; i++ {
		if res.AryMapVal("users")[i].IntVal("size") < res.AryMapVal("users")[i+1].IntVal("size") {
			t.Errorf("err res(4.NICK_SEARCH ACCOUNT_SIZE_SORT): %v", util.S2Json(res))
			return
		}
	}

	t.Logf("res(4.NICK_SEARCH ACCOUNT_SIZE_SORT): %v", util.S2Json(res))

	//5. REGISTER_SORT DEFAULT_SORT
	searchMethod = usrdb.PHONE_SEARCH
	phoneKey = "1"
	sort = usrdb.REGISTER_SORT
	res, err = DoSearchUserList(token, nickKey, phoneKey,
		searchMethod, skip, limit, sort, startTime, overTime)
	if err != nil {
		t.Error(err)
		return
	}
	if len(res.AryMapVal("users")) != 5 {
		t.Errorf("err res(5.PHONE_SEARCH REGISTER_SORT): %v", util.S2Json(res))
		return
	}
	for i := 0; i < 4; i++ {
		if res.AryMapVal("users")[i].IntVal("time") < res.AryMapVal("users")[i+1].IntVal("time") {
			t.Errorf("err res(5.PHONE_SEARCH REGISTER_SORT): %v", util.S2Json(res))
			return
		}
	}
	for i := 0; i < 5; i++ {
		if !strings.Contains(res.AryMapVal("users")[i].StrVal("phone"), phoneKey) {
			t.Errorf("err res(5.PHONE_SEARCH REGISTER_SORT): %v", util.S2Json(res))
			return
		}
	}
	t.Logf("res(5.PHONE_SEARCH REGISTER_SORT): %v", util.S2Json(res))

	//6. TIME_SEARCH REGISTER_SORT
	searchMethod = usrdb.TIME_SEARCH
	startTime = res.AryMapVal("users")[0].IntVal("time") - 1000000
	overTime = res.AryMapVal("users")[0].IntVal("time") + 1000000
	sort = usrdb.REGISTER_SORT
	res, err = DoSearchUserList(token, nickKey, phoneKey,
		searchMethod, skip, limit, sort, startTime, overTime)
	if err != nil {
		t.Error(err)
		return
	}
	if len(res.AryMapVal("users")) != 5 {
		t.Errorf("err res(6. TIME_SEARCH REGISTER_SORT): %v", util.S2Json(res))
		return
	}
	for i := 0; i < 4; i++ {
		if res.AryMapVal("users")[i].IntVal("time") < res.AryMapVal("users")[i+1].IntVal("time") {
			t.Errorf("err res(6. TIME_SEARCH REGISTER_SORT): %v", util.S2Json(res))
			return
		}
	}
	for i := 0; i < 5; i++ {
		if startTime > res.AryMapVal("users")[i].IntVal("time") || res.AryMapVal("users")[i].IntVal("time") > overTime {
			t.Errorf("err res(6. TIME_SEARCH REGISTER_SORT): %v", util.S2Json(res))
			return
		}
	}

	t.Logf("res(6. TIME_SEARCH REGISTER_SORT): %v", util.S2Json(res))

}

func prepareData() {
	Remove()
	u := &usrdb.Usr{}
	var err error
	u = &usrdb.Usr{
		Account: "",
		Pwd:     "123",
		Type:    1,
	}
	login := 0
	for i := 0; i < 20; i++ {
		u.Account += "1"
		_, err = DoRegister(u, login)
		if err != nil {
			log.E("err: %v", err)
			return
		}
	}

	uid := []string{"u1", "u2", "u3", "u4", "u5", "u6", "u7", "u8", "u9", "u10", "u11", "u12", "u13", "u14", "u15", "u16", "u17", "u18", "u19"}
	u = &usrdb.Usr{}
	index := 0
	nickname := "q" + NewStringLen(5)
	for i := 0; i < 10; i++ {
		u.Email = ""
		u.Id = uid[index]
		index++
		nickname += "1"
		u.Attrs = util.Map{
			"nickname": nickname,
		}
		if i < 2 {
			u.Phone = NewDigitLen(11)
			u.Email = NewStringLen(5) + "@" + NewStringLen(3) + "." + NewStringLen(3)
		} else {
			u.Phone = ""
			u.Email = ""
		}
		usrdb.UpdateUserInfo(u)
	}
	nickname = "x" + NewStringLen(5)
	for i := 0; i < 5; i++ {
		u.Phone = "1" + NewDigitLen(10)
		u.Id = uid[index]
		index++
		u.Phone = NewDigitLen(11)
		nickname += "2"
		u.Attrs = util.Map{
			"nickname": nickname,
		}
		usrdb.UpdateUserInfo(u)
	}
	nickname = "y" + NewStringLen(5)
	for i := 0; i < 4; i++ {
		u.Phone = NewDigitLen(11)
		u.Id = uid[index]
		index++
		nickname += "3"
		u.Attrs = util.Map{
			"nickname": nickname,
		}
		usrdb.UpdateUserInfo(u)
	}
}

var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
var StdDigit = []byte("0123456789")

func NewDigitLen(length int) string {
	return NewLenChars(length, StdDigit)
}
func NewStringLen(length int) string {
	return NewLenChars(length, StdChars)
}
func NewLenChars(length int, chars []byte) string {
	if length == 0 {
		return ""
	}
	clen := len(chars)
	if clen < 2 || clen > 256 {
		panic("Wrong charset length for NewLenChars()")
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			panic("Error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				continue // Skip this number to avoid modulo bias.
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}
