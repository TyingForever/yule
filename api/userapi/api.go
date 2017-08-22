package userapi

import (
	"github.com/Centny/gwf/routing"
)

var SrvAddr = func() string {
	panic("yule server is not init")
}

const (
	LOG_API = true
	LOG_API_TEST = true
	TAG_LANDLORDS = "landlords--"
	TAG_LANDLORDS_TEST = "landlords_test--"
)

func Hand(pre string, mux *routing.SessionMux) {
	mux.HFilterFunc("^"+pre+"/usr/.*$", LoginFilter)

	//user
	mux.HFunc("^"+pre+"/pub/api/register(\\?.*)?$", Register)

	//先发送短信，获得验证码->发送发证码进行验证完成操作
	mux.HFunc("^"+pre+"/usr/api/sendMessage(\\?.*)?$", SendMessage)
	mux.HFunc("^"+pre+"/usr/api/bindPhone(\\?.*)?$", BindPhone)
	mux.HFunc("^"+pre+"/pub/api/bindPhone(\\?.*)?$", BindPhone)

	mux.HFunc("^"+pre+"/usr/api/bindEmail(\\?.*)?$", BindEmail)
	mux.HFunc("^"+pre+"/usr/api/sendEmail(\\?.*)?$", SendEmail)

	//
	//mux.HFunc("^"+pre+"/pub/api/loginPhone(\\?.*)?$", LoginPhone)
	//mux.HFunc("^"+pre+"/pub/api/loginEmail(\\?.*)?$", LoginEmail)

	mux.HFunc("^"+pre+"/pub/api/login(\\?.*)?$", Login)
	mux.HFunc("^"+pre+"/usr/api/getUser(\\?.*)?$", GetUser)
	mux.HFunc("^"+pre+"/usr/api/updateUser(\\?.*)?$", UpdateUser)
	mux.HFunc("^"+pre+"/usr/api/findUsers(\\?.*)?$", FindUsers)
	mux.HFunc("^"+pre+"/usr/api/logout(\\?.*)?$", Logout)
	mux.HFunc("^"+pre+"/usr/api/searchUsers(\\?.*)?$", SearchUsers)

	//landlords
	mux.HFunc("^"+pre+"/usr/api/entryLandlords(\\?.*)?$", EntryLandlords)
	mux.HFunc("^"+pre+"/usr/api/startNewLandlords(\\?.*)?$", StartNewLandlords)
	mux.HFunc("^"+pre+"/usr/api/getLandlords(\\?.*)?$", GetLandlords)
	mux.HFunc("^"+pre+"/usr/api/operateLandlords(\\?.*)?$", OperateLandlords)
}
