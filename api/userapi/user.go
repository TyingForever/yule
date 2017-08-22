package userapi

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"yule/db/usrdb"
)

//注册用户
//通过用户对象的相关字段注册用户
//@url,需求登录，POST请求
//	~/pub/api/createUser		POST	application/json
//@arg,json对象中的参数以及query中参数
//	account		R	用户名列表，如果普通用户名、手机、邮箱等
//	pwd		R	用户密码
/*
	样例	~/pub/api/createUser?login=1	1为注册完成自动登录，0为不登录
	{
		"account": "2",
		"pwd": "123",
		"type":1
	}
*/
//@ret,返回通用code/data
//	code	I	0：登录成功，1：参数错误，2：json body错误，3：注册用户失败
//	token	S	成功登录的token
//	usr	O	注册成功的用户对象
//	account	S	账号
//	id	S	用户id
//	last	I	上次更新的时间
//	ip	S	标志ip
//	role	I	角色1:管理员；2:用户
//	status	I	状态
//	time	I	用户注册时间
//	type	I	用户登录类型
/*	样例
		{
	    "code":0,
	    "data":{
		"token":"59130F1EA8C72E29C0632100",
		"usr":{
		    "account":"2",
		    "id":"u2",
		    "ip":"127.0.0.1",
		    "last":1494421278511,
		    "role":2,
		    "status":10,
		    "time":1494421278511,
		    "type":1
		}
	    }
	}
*/
//@tag,用户,注册
//@author,zhnagyq,2017-05-10
//@case,demo1
func Register(hs *routing.HTTPSession) routing.HResult {
	var login int = 0
	var err error
	err = hs.ValidF(`
	login,O|I,R:-1;
	`, &login)

	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}

	user := &usrdb.Usr{}
	err = hs.UnmarshalJ(user)
	if err != nil {
		log.E("CreateUser decoding request json body fail by error(%v)", err)
		return hs.MsgResErr2(1, "arg-err", util.Err("CreateUser->decoding request json body fail by error(%v)", err))
	}
	var addr = strings.Split(hs.R.Header.Get("X-Real-IP"), ":")[0]
	if len(addr) < 1 {
		addr = strings.Split(hs.R.RemoteAddr, ":")[0]
	}
	user.Ip = addr
	//log.D("user: %v", util.S2Json(user))

	FilterUserAttrs(user)
	err = usrdb.AddUser(user)
	if err != nil {
		if err.Error() == usrdb.UserExistErr.Error() {
			return hs.MsgResErr(4, "该用户已存在", err)
		}
		log.E("CreateUser->adding user by data(%v) error(%v)", util.S2Json(user), err)
		return hs.MsgResErr2(3, "srv-err", util.Err("CreateUser->add user by (%v) error(%v)", util.S2Json(user), err))
	}
	if login < 1 {
		log.D("CreateAutoUser-> adding user success by data(%v),login(%v)", util.S2Json(user), login)
		return hs.MsgRes(util.Map{
			"usr": user,
		})
	}

	//如果选择注册完成立即登录，则获取token
	token, err := do_login(hs, user)
	if err != nil {
		log.E("CreateAutoUser-> adding user success, but login fail by data(%v) error(%v)", util.S2Json(user), err)
		return hs.MsgRes2(4, util.Map{
			"usr": user,
		})
	}
	if usrdb.ShowLog {
		log.D("CreateAutoUser-> adding user and login success by data(%v)", util.S2Json(user))
	}
	return hs.MsgRes(util.Map{
		"token": token,
		"usr":   user,
	})

}

//登录
//用户登录
//@url,不需求登录，普通Get请求
//	~/pub/api/login		GET
//
//@arg,普通query参数
//	account		R	用户名列表，如普通用户名、手机、邮箱等
//	pwd		R	用户密码
/*
	样例1 登陆普通用户
	~/pub/api/login?account=xx&pwd=xx

*/
//@ret,返回通用code/data
//	code	I	0：登录成功，1：参数错误，2：json body错误，3：登录用户失败
//	token	S	成功登录的token
//	usr	O	注册成功的用户对象
//	account	S	账号
//	id	S	用户id
//	last	I	上次更新的时间
//	ip	S	标志ip
//	role	I	角色1:管理员；2:用户
//	status	I	状态
//	time	I	用户注册时间
//	type	I	用户登录类型
/*	样例
		{
		    "code":0,
		    "data":{
			"token":"59131D1DA8C72E37EC0F123A",
			"usr":{
			    "account":"1",
			    "id":"u1",
			    "ip":"127.0.0.1",
			    "last":1494422126940,
			    "role":2,
			    "status":10,
			    "time":1494422126940,
			    "type":1
			}
		    }
}
*/
//@tag,登录,用户
//@author,zhangyq,2017-05-05
//@case,demo1
func Login(hs *routing.HTTPSession) routing.HResult {
	var account, pwd string
	err := hs.ValidCheckVal(`
		account,R|S,L:0;
		pwd,R|S,L:0;
		`, &account, &pwd)
	if err != nil {
		log.E("arg-err,%v", err)
		return hs.MsgResErr(1, "arg-err", err)
	}
	//log.D("account,pwd: %v ,%v", account, pwd)

	user, err := usrdb.FindUserByAccountPwd(account, pwd)
	if err != nil {
		return hs.MsgResErr(10, "密码有误", util.Err("account(%v) pwd(%v)", account, pwd))
	}
	//log.D("user :%v", util.S2Json(user))
	if user.Status == usrdb.USR_S_D {
		return hs.MsgResErr(11, "该用户不存在", util.Err("uid(%v) status(%v)", user.Id, user.Status))
	}

	token, err := do_login(hs, user)
	if err != nil {
		return hs.MsgResErr2(3, "登录失败", util.Err("flush uid to session error(%v)", err))
	}
	//重登功能
	redirect := hs.CheckVal("url")
	if len(redirect) < 1 {
		usrdb.AddUserSession(user.Id, token)
		return hs.MsgRes(util.Map{
			"token": token,
			"usr":   user,
		})
	} else {
		hs.Redirect(token_url(redirect, token))
		return routing.HRES_RETURN
	}
}

//获取用户信息
//通过token获取用户信息
//@url,需求登录，GET请求
//	~/usr/api/getUser	GET
//@arg,普通Query参数
//	oid		O	管理员查询的普通用户id
//	method		O	获取用户信息的方式
/*	样例
//	用户获取自己的用户信息
	~/usr/api/getUser?token=xxxx
//	管理员获取普通用户的信息
	~/usr/api/getUser?token=xxxx&oid=xx&method=x
*/
//@ret,返回能用code/data
//	code								I	0:获取成功,1:未登录,2:服务器错误
//	usr								O	已经登录的用户信息
//	last								I	最后更新时间
//	id								I	账号
//	account								I	密码
//	image								S	头像
//	ip								S	ip身份
//	name								S	姓名
//	role								I	角色 1.管理员 2.普通用户
//	status								I	状态
//	create								I	创建时间
//	attrs								O	用户属性
//	attrs.nickname							S	用户属性昵称
//	attrs.sex							S	用户属性性别
//	attrs.age							I	用户属性年龄
//	attrs.birthday							I	用户属性生日
//	attrs.hometown							S	用户属性家乡
//	attrs.location							S	用户属性现住址
//	attrs.profession						S	用户属性专业
/*	样例
	{
	    "code":0,
	    "data":{
		"usr":{
		    "id":"u3",
		    "account":"3",
		    "pwd":"40bd001563085fc35165329ea1ff5c5ecbdbbeef",
		    "ip":"127.0.0.1",
		    "attrs":{
			"age":12,
			"birthday":1494297238378,
			"hometown":"Xinyi",
			"location":"Guangzhou",
			"profession":"Programmer"
		    },
		    "status":10,
		    "type":1,
		    "role":1,
		    "last":1494296723965,
		    "time":1494296723965
		}
	    }
}
*/
//@tag,用户,信息
//@author,zhangyq,2017-05-22
//@case,demo1
func GetUser(hs *routing.HTTPSession) routing.HResult {
	var oid string
	var method int
	var err error
	err = hs.ValidF(`
	oid,O|S,L:0;
	method,O|I,R:-1;
	`, &oid, &method)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}
	uid := hs.StrVal("uid")
	user, err := usrdb.GetUserInfo(uid)
	if err != nil {
		return hs.MsgResErr2(3, "查找失败", util.Err("flush uid to session error(%v)", err))
	}
	if user.Role == usrdb.USR_ADMIN {
		if oid != "" && method == 1 {
			user, err = usrdb.GetUserInfo(oid)
			if err != nil {
				return hs.MsgResErr2(3, "查找失败", util.Err("flush uid to session error(%v)", err))
			}
			return hs.MsgRes(util.Map{"user": user})
		}
	}
	return hs.MsgRes(util.Map{"user": user})
}

//更改用户信息
//通过用户对象的相关字段更新用户
//@url,需求登录，POST请求
//	~/usr/api/updateUser		POST	application/json
//@arg,json对象中的参数以及query中参数
//	oid		O	管理员查询的普通用户id
//	method		O	获取用户信息的方式
//	role		O	请求更改用户信息的角色
//	account		O	用户名列表，如果普通用户名、手机、邮箱等，仅用于添加登录名
//	pwd		O	用户密码
//	email		O	用户邮箱
//	phone		O	用户手机
//	attrs		O	用户自定义属性
//	attrs.nickname							O	用户属性昵称
//	attrs.sex							O	用户属性性别
//	attrs.age							O	用户属性年龄
//	attrs.birthday							O	用户属性生日
//	attrs.hometown							O	用户属性家乡
//	attrs.location							O	用户属性现住址
//	attrs.profession						O	用户属性专业
/*	样例
	//更改用户自己的基本信息
	~/usr/api/updateUser?token=xxxx
	{
		"attrs": {
				"nickname":"abc",
				"age":12,
				"birthday":1494297238378,
				"hometown":"Xinyi",
				"location":"Guangzhou",
				"profession":"Programmer"
		}
	}
	//管理员更改普通用户的信息
	~/usr/api/updateUser?token=xxxx&oid=xx&method=x&role=x
	{
		"phone":"12345678901",
		"email":"123@abc.com",
		"name":"abcde",
		"pwd":"123",
		"attrs": {
				"nickname":"abc",
				"age":12,
				"birthday":1494297238378,
				"hometown":"Xinyi",
				"location":"Guangzhou",
				"profession":"Programmer"
		}
	}
*/
//@ret,返回通用code/data
//	code	I	0：登录成功，1：json body错误，2：更新用户失败，401：无权限
//	token	S	成功登录的token
/*	样例
	{
		"code": 0,
		"data": "OK"
	}
*/
//@tag,用户,更新
//@author,zhangyq,2017-05-22
//@case,demo1
func UpdateUser(hs *routing.HTTPSession) routing.HResult {
	var oid string
	var method int = 0
	var role int = 0
	var err error
	err = hs.ValidF(`
	oid,O|S,L:0;
	method,O|I,R:-1;
	role,O|I,R:-1;
	`, &oid, &method, &role)
	log.D("params: %v,%v,%v", oid, method, role)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}
	user := &usrdb.Usr{}
	err = hs.UnmarshalJ(user)
	if err != nil {
		log.E("CreateUser decoding request json body fail by error(%v)", err)
		return hs.MsgResErr2(1, "arg-err", util.Err("UpdateUser->decoding request json body fail by error(%v)", err))
	}
	if role == usrdb.USR_ADMIN && method == 1 {
		if oid != "" && method == 1 {
			user.Id = oid
			err = usrdb.UpdateUserInfo(user)
			if err != nil {
				return hs.MsgResErr2(3, "查找失败", util.Err("flush uid to session error(%v)", err))
			}
			return hs.MsgRes(util.Map{"user": user})
		}
	} else {
		uid := hs.StrVal("uid")
		user.Id = uid
		err = usrdb.UpdateUserInfo(user)
		if err != nil {
			log.E("UpdateUser failed %v", err)
			return hs.MsgResErr2(7, "mongodb-err", util.Err("UpdateUser->update fail by error(%v)", err))

		}
	}
	return hs.MsgRes(util.Map{"data": "OK"})
}

//管理员获取用户信息
//通过token获取用户信息
//@url,需求登录，GET请求
//	~/usr/api/findUsers	GET
//@arg,普通Query参数
//	token		R	用户名列表，如果普通用户名、手机、邮箱等，仅用于添加登录名
//	skip		O	跳过搜索数
//	limit		O	搜索最大数目
/*	样例
	~/usr/api/findUsers?token=xxx
*/
//@ret,返回能用code/data
//	code								I	0:获取成功,1:未登录,2:服务器错误
//	usrs								A	已经登录的用户信息
//	last								I	最后更新时间
//	id								I	账号
//	account								I	密码
//	image								S	头像
//	ip								S	ip身份
//	name								S	姓名
//	role								I	角色 1.管理员 2.普通用户
//	status								I	状态
//	time								I	创建时间
/*	样例
{
    "code":0,
    "data":{
        "users":[
            {
                "account":"11",
                "id":"u2",
                "ip":"127.0.0.1",
                "last":1494579940261,
                "role":2,
                "status":10,
                "time":1494579940261,
                "type":1
            },
            {
                "account":"111",
                "id":"u3",
                "ip":"127.0.0.1",
                "last":1494579940267,
                "role":2,
                "status":10,
                "time":1494579940267,
                "type":1
            },
            {
                "account":"1111",
                "id":"u4",
                "ip":"127.0.0.1",
                "last":1494579940276,
                "role":2,
                "status":10,
                "time":1494579940276,
                "type":1
            },
            {
                "account":"11111",
                "id":"u5",
                "ip":"127.0.0.1",
                "last":1494579940284,
                "role":2,
                "status":10,
                "time":1494579940284,
                "type":1
            },
            {
                "account":"111111",
                "id":"u6",
                "ip":"127.0.0.1",
                "last":1494579940290,
                "role":2,
                "status":10,
                "time":1494579940290,
                "type":1
            }
        ]
    }
}
*/
//@tag,列出用户
//@author,zhangyq,2017-05-10
//@case,demo1
func FindUsers(hs *routing.HTTPSession) routing.HResult {
	var skip, limit int
	err := hs.ValidF(`
	skip,O|I,R:-1;
	limit,O|I,R:-1;
	`, &skip, &limit)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}
	uid := hs.StrVal("uid")
	total, users, err := usrdb.ListUsersOrdinary(uid, skip, limit)
	if err != nil {
		log.E("FindUsers get users err(%v)", err)
		return hs.MsgResErr2(3, "查找失败", util.Err("flush uid to session error(%v)", err))
	}
	if usrdb.ShowLog {
		log.D("Total: %v", total)
	}
	return hs.MsgRes(util.Map{"users": users, "total": total})

}

//管理员通过搜索获取用户信息并进行排序
//通过token获取用户信息
//@url,需求登录，GET请求
//	~/usr/api/searchUsers	GET
//@arg,普通Query参数
//	token		R	用户名列表，如果普通用户名、手机、邮箱等，仅用于添加登录名
//	searchMethod	O	关键字搜索方式，默认为1,昵称搜索为2，手机号搜索为3，注册时间段搜索为4
//	nickKey		O	关键字昵称搜索
//	phoneKey	O	关键字手机号搜索
//	startTime		O	关键字注册时间段开始时间
//	overTime		O	关键字注册时间段结束时间
//	skip	O	跳过搜索数
//	limit		O	搜索最大数目
//	sort		O	排序方式，1,,默认；2,创建时间；3,修改时间;4,登录账号数目
/*	样例
	~/usr/api/searchUsers?token=591BA02EA8C72E1BCC3BADE8&nickKey=x&phoneKey=x&searchMethod=2&skip=0&limit=100&sort=2&startTime=0&overTime=0
*/
//@ret,返回能用code/data
//	code								I	0:获取成功,1:未登录,2:服务器错误
//	usrs								A	已经登录的用户信息
//	last								I	最后更新时间
//	id								I	账号
//	account								I	密码
//	image								S	头像
//	ip								S	ip身份
//	name								S	姓名
//	role								I	角色 1.管理员 2.普通用户
//	status								I	状态
//	time								I	创建时间
/*	样例
{
    "code":0,
    "data":{
        "users":[
            {
                "id":"u15",
                "account":"111111111111111",
                "phone":"19854775499",
                "ip":"127.0.0.1",
                "attrs":{
                    "nickname":"x2R24222222"
                },
                "status":10,
                "type":1,
                "role":2,
                "last":1494925656919,
                "time":1494925656839
            },
            {
                "id":"u14",
                "account":"11111111111111",
                "phone":"15648673524",
                "ip":"127.0.0.1",
                "attrs":{
                    "nickname":"x2R2422222"
                },
                "status":10,
                "type":1,
                "role":2,
                "last":1494925656917,
                "time":1494925656833
            },
            {
                "id":"u13",
                "account":"1111111111111",
                "phone":"15517536347",
                "ip":"127.0.0.1",
                "attrs":{
                    "nickname":"x2R242222"
                },
                "status":10,
                "type":1,
                "role":2,
                "last":1494925656915,
                "time":1494925656826
            },
            {
                "id":"u12",
                "account":"111111111111",
                "phone":"19231019961",
                "ip":"127.0.0.1",
                "attrs":{
                    "nickname":"x2R24222"
                },
                "status":10,
                "type":1,
                "role":2,
                "last":1494925656912,
                "time":1494925656816
            },
            {
                "id":"u11",
                "account":"11111111111",
                "phone":"19746736678",
                "ip":"127.0.0.1",
                "attrs":{
                    "nickname":"x2R2422"
                },
                "status":10,
                "type":1,
                "role":2,
                "last":1494925656909,
                "time":1494925656807
            }
        ]
    }
}
*/
//@tag,列出用户,搜索,排序
//@author,zhangyq,2017-05-10
//@case,demo1
func SearchUsers(hs *routing.HTTPSession) routing.HResult {
	var searchMethod, skip, limit, sort int
	var startTime, overTime int64
	var nickKey, phoneKey string
	err := hs.ValidF(`
	searchMethod,O|I,R:0;
	nickKey,O|S,L:0;
	phoneKey,O|S,L:0;
	startTime,O|I,R:-1;
	overTime,O|I,R:-1;
	skip,O|I,R:-1;
	limit,O|I,R:0;
	sort,O|I,R:0;
	`, &searchMethod, &nickKey, &phoneKey, &startTime, &overTime, &skip, &limit, &sort)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}

	if usrdb.ShowLog {
		log.D("getParams: nickname:%v,searchMethod:%v,sortMetod:%v,skip:%v,limit:%v", nickKey, searchMethod, sort, skip, limit)
	}

	uid := hs.StrVal("uid")
	total, users, err := usrdb.SearchUsersOrdinary(uid, nickKey, phoneKey, searchMethod, sort, skip, limit, startTime, overTime)
	if err != nil {
		log.E("FindUsers get users err(%v)", err)
		return hs.MsgResErr2(3, "查找失败", util.Err("flush uid to session error(%v)", err))
	}
	if usrdb.ShowLog {
		log.D("Total: %v", total)
	}
	return hs.MsgRes(util.Map{"users": users, "total": total})

}

//注销用户
//通过用户对象的相关字段注销用户
//@url,需求登录，GET请求
//	~/usr/api/logout		Get	application/json
//@arg,普通Query参数
/*	样例
	~/usr/api/logout?token=xxxx
*/
//@ret,返回通用code/data
//	code	I	0：登录成功，1：json body错误，2：更新用户失败，401：无权限
//	data	S	返回结果
/*	样例
	{
		"code": 0,
		"data": "OK"
	}
*/
//@tag,用户,注销
//@author,zhangyq,2017-05-09
//@case,demo1
func Logout(hs *routing.HTTPSession) routing.HResult {
	uid := hs.StrVal("uid")
	err := usrdb.Logout(uid)
	if err != nil {
		log.E("FindUsers get users err(%v)", err)
		return hs.MsgResErr2(3, "删除失败", util.Err("flush uid to session error(%v)", err))
	}
	return hs.MsgRes(util.Map{"data": "OK"})
}

//绑定手机号码
//通过用户对象的相关字段绑定手机号码
//@url,需求登录，Get请求
//	~/usr/api/bindPhone		Get
//@arg,json对象中的参数以及query中参数
//	types	R	BIND=1,MODIFYBIND=2,REGISTER=3,VERIFY=4,VERIFYBIND=5,RESETPWD=6,LOGIN=7,UNBIND=8
//	phone	R	手机号码,eg:13513513535
//	phoneOld	R	手机号码,eg:13513513535
//	pcode	O	短信验证码的值
/*
	样例	~/usr/api/bindPhone?token=xx&phone=xx&phoneOld=xx&types=xx&pcode=xx"	1为登录用户绑定手机号码,2为登录用户改绑手机号码
*/
//@ret,返回通用code/data
//	code	I	0：绑定成功，1：参数错误，2：json body错误，
/*	样例
{
    "code":0,
    "data":{
        "data":"OK"
    }
}
*/
//@tag,用户,手机
//@author,zhnagyq,2017-05-11
//@case,demo1
func BindPhone(hs *routing.HTTPSession) routing.HResult {
	var phone, phoneOld string
	var pcode, login, t int
	err := hs.ValidCheckVal(`
		pcode,R|I,R:0;
		phone,R|S,L:0;
		phoneOld,O|S,L:0;
		types,R|I,R:0;
		login,O|I,R:-1;
		`, &pcode, &phone, &phoneOld, &t, &login)
	if err != nil {
		return hs.MsgResErr2(1, "参数错误", err)
	}
	uid := hs.StrVal("uid")
	log.D("BindUserPhone receive phone(%v), phoneOld(%v), pcode(%v), t(%v), uid(%v), token(%v), do_login(%v)",
		phone, phoneOld, pcode, t, uid, hs.StrVal("token"), login)

	if uid == "" && (t == usrdb.BIND || t == usrdb.MODIFYBIND || t == usrdb.UNBIND) {
		return hs.MsgResE(6, "please login first")
	}

	code, err := usrdb.VerifyCode(phone, usrdb.PHONE, pcode, 300000)
	if err != nil {
		log.E("BindUserPhone verify code phone(%v) pcode(%v) err(%v)", phone, pcode, err)
		return hs.MsgResE(code, err.Error())
	}

	var msg string
	var u *usrdb.Usr
	if t == usrdb.MODIFYBIND {
		code, msg, err = usrdb.ChangeBindPhone(uid, phone, phoneOld)
	} else if t == usrdb.BIND {
		code, msg, err = usrdb.BindPhone(uid, phone)
	} else if t == usrdb.UNBIND {
		code, msg, err = usrdb.UnBindPhone(uid, phone)
	} else {
		return hs.MsgResErr(1, "参数错误", util.Err("type(%v) err", t))
	}

	if err != nil {
		log.E("BindPhone uid(%v) err(%v)", uid, err)
		return hs.MsgResErr2(code, msg, err)
	}

	if login > 0 {
		token, err := do_login(hs, u)
		if err != nil {
			log.E("BindPhone Do_login uid(%v) err(%v)", u.Id, err)
			return hs.MsgResErr2(11, "验证成功但登录失败", err)
		}
		return hs.MsgRes(util.Map{"token": token, "usr": u})
	}

	return hs.MsgRes(util.Map{"data": "OK"})
}

//发送信息
//发短息到手机
//@url,不需求登录，GET请求
//	~/usr/api/sendMessage		GET
//@arg,query中参数
//	types	R	BIND=1,MODIFYBIND=2,REGISTER=3,VERIFY=4,VERIFYBIND=5,RESETPWD=6,LOGIN=7,UNBIND=8
//	phone	R	手机号码,eg:13513513535
//	mark	O	验证码的标记id
//	vcode	O	图片验证码的值
//	sign	O	签名，有值时使用该值验证合法性，否则用图片验证码验证。生成方法，将字符串 key=对应的密钥&mobile=对应的手机号&types=对应的类型 进行md5加密
/*
	样例	~/usr/api/sendMessage?types=1&token=xxx&phone=xxxx	1为登录用户绑定手机号码时发送短信,2为登录用户改绑手机号码发送短信
*/
//@ret,返回通用code/data
//	code	I	0：发送短信成功，1：参数错误，2：请输入正确手机号码，3：手机绑定出错，4：类型错误
//	Type	I	BIND=1,MODIFYBIND=2,REGISTER=3,VERIFY=4,VERIFYBIND=5,RESETPWD=6,LOGIN=7,UNBIND=8
//	Category	S	"phone","email"
//	Send	S	手机号码,eg:13513513535
//	Status	I	验证码的标记id
//	Code	I	图片验证码的值
//	Time	I	发送时间
/*	样例
		{
    "code":0,
    "data":{
        "phoneCode":{
            "Category":"phone",
            "Code":185988,
            "Send":"13800138004",
            "Status":"N",
            "Time":1494497830792,
            "Type":1
        }
    }
}
*/
//@tag,用户,手机
//@author,zhnagyq,2017-05-11
//@case,demo1
func SendMessage(hs *routing.HTTPSession) routing.HResult {
	var phone string
	var types, vcode, mark, sign int
	err := hs.ValidCheckVal(`
		phone,R|S,L:0;
		types,R|I,R:0;
		vcode,O|I,R:-1;
		mark,O|I,R:-1;
		sign,O|I,R:-1;
		`, &phone, &types, &vcode, &mark, &sign)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}

	//log.D("SendPhoneMessage receive mobile(%s),types(%v),vcode(%v),mark(%v),sign(%v)", phone, types, vcode, mark, sign)

	reg := regexp.MustCompile("^[\\d]{11}$")
	if !reg.MatchString(phone) {
		log.W("SendPhoneMessage receive param not phone number, phone(%v), types(%v)", phone, types)
		return hs.MsgResErr(2, "请输入正确手机号码", util.Err("phone(%v) err", phone))
	}

	users, err := usrdb.FindUsers(
		bson.M{"$or": []bson.M{
			bson.M{"phone": phone},
			bson.M{"attrs.privated.phone": phone},
		}},
		bson.M{"_id": 1})

	if usrdb.ShowLog {
		log.D("FindUser %v", util.S2Json(users))
	}
	if err != nil {
		log.E("%v", fmt.Sprintf("SendPhoneMessage query user by phone(%v) error", err))
		return hs.MsgResE(1, fmt.Sprintf("SendPhoneMessage query user by phone(%v) error", err))
	}

	if types == usrdb.BIND || types == usrdb.REGISTER || types == usrdb.MODIFYBIND {
		if len(users) > 0 {
			log.W("%v", fmt.Sprintf("phone(%v) had been register or binded", phone))
			return hs.MsgResErr2(3, "手机已注册或被绑定", util.Err("phone(%v) had been register or binded", phone))
		}
	} else if types == usrdb.RESETPWD || types == usrdb.LOGIN || types == usrdb.VERIFY || types == usrdb.UNBIND {
		if len(users) == 0 {
			log.W("%v", fmt.Sprintf("phone(%v) had not been register or binded", phone))
			return hs.MsgResErr2(3, "手机未被注册或绑定", util.Err("phone(%v) had not been register or binded", phone))
		}
	} else {
		return hs.MsgResErr2(4, "类型错误", util.Err("type(%v) is invalid", types))
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(899999) + 100000
	code, err := usrdb.UpsertCode(phone, usrdb.PHONE, num, 30000, types)
	if err != nil {
		log.E("SendPhoneMessage UpsertCode phone(%v) num(%v) types(%v) err(%v)", phone, num, types, err)
		return hs.MsgResErr2(code, "服务器错误", err)
	}

	phoneCode := &usrdb.VerificationCode{
		Send:     phone,
		Category: usrdb.PHONE,
		Code:     num,
		Status:   usrdb.PCS_NORMAL,
		Time:     util.Now(),
		Type:     types,
	}
	return hs.MsgRes(util.Map{"phoneCode": phoneCode})
}

//发送信息
//发信息到邮箱
//@url,不需求登录，GET请求
//	~/usr/api/sendEmail		GET
//@arg,query中参数
//	types	R	BIND=1,MODIFYBIND=2,REGISTER=3,VERIFY=4,VERIFYBIND=5,RESETPWD=6,LOGIN=7,UNBIND=8
//	email	R	邮箱号码,eg:84323229@qq.com
//	mark	O	验证码的标记id
//	vcode	O	图片验证码的值
//	sign	O	签名，有值时使用该值验证合法性，否则用图片验证码验证。生成方法，将字符串 key=对应的密钥&mobile=对应的邮箱号&types=对应的类型 进行md5加密
//
/*
	样例	~/usr/api/sendEmail?types=1&token=xxx&email=xxxx	1为登录用户绑定邮箱号码时发送短信,2为登录用户改绑邮箱号码发送短信
*/
//@ret,返回通用code/data
//	code	I	0：发送短信成功，1：参数错误，2：请输入正确手机号码，3：手机绑定出错，4：类型错误
//	Type	I	BIND=1,MODIFYBIND=2,REGISTER=3,VERIFY=4,VERIFYBIND=5,RESETPWD=6,LOGIN=7,UNBIND=8
//	Category	S	"phone","email"
//	Send	S	邮箱号码,eg:84323229@qq.com
//	Status	I	验证码的标记id
//	Code	I	图片验证码的值
//	Time	I	发送时间
/*	样例
{
    "code":0,
    "data":{
        "phoneCode":{
           "emailCode":{
        "Category":"email",
        "Code":883644,
        "Send":"84323229@qq.com",
        "Status":"N",
        "Time":1494497752129,
        "Type":1
        }
    }
}
}
*/
//@tag,用户,邮箱
//@author,zhnagyq,2017-05-11
//@case,demo1
func SendEmail(hs *routing.HTTPSession) routing.HResult {
	var email string
	var types, vcode, mark, sign int
	err := hs.ValidCheckVal(`
		email,R|S,L:0;
		types,R|I,R:0;
		vcode,O|I,R:-1;
		mark,O|I,R:-1;
		sign,O|I,R:-1;
		`, &email, &types, &vcode, &mark, &sign)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}

	//log.D("SendemailMessage receive mobile(%s),types(%v),vcode(%v),mark(%v),sign(%v)", email, types, vcode, mark, sign)

	reg := regexp.MustCompile("\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*")
	if !reg.MatchString(email) {
		log.W("SendemailMessage receive param not email number, email(%v), types(%v)", email, types)
		return hs.MsgResErr(2, "请输入正确邮箱号码", util.Err("email(%v) err", email))
	}

	users, err := usrdb.FindUsers(
		bson.M{"$or": []bson.M{
			bson.M{"email": email},
			bson.M{"attrs.privated.email": email},
		}},
		bson.M{"_id": 1})

	if usrdb.ShowLog {
		log.D("FindUser %v", util.S2Json(users))
	}
	if err != nil {
		log.E("%v", fmt.Sprintf("SendemailMessage query user by email(%v) error", err))
		return hs.MsgResE(1, fmt.Sprintf("SendemailMessage query user by email(%v) error", err))
	}

	if types == usrdb.BIND || types == usrdb.REGISTER || types == usrdb.MODIFYBIND {
		if len(users) > 0 {
			log.W("%v", fmt.Sprintf("email(%v) had been register or binded", email))
			return hs.MsgResErr2(3, "邮箱已注册或被绑定", util.Err("email(%v) had been register or binded", email))
		}
	} else if types == usrdb.RESETPWD || types == usrdb.LOGIN || types == usrdb.VERIFY || types == usrdb.UNBIND {
		if len(users) == 0 {
			log.W("%v", fmt.Sprintf("email(%v) had not been register or binded", email))
			return hs.MsgResErr2(3, "邮箱未被注册或绑定", util.Err("email(%v) had not been register or binded", email))
		}
	} else {
		return hs.MsgResErr2(4, "类型错误", util.Err("type(%v) is invalid", types))
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(899999) + 100000
	code, err := usrdb.UpsertCode(email, usrdb.EMAIL, num, 30000, types)
	if err != nil {
		log.E("SendemailMessage UpsertCode email(%v) num(%v) types(%v) err(%v)", email, num, types, err)
		return hs.MsgResErr2(code, "服务器错误", err)
	}

	emailCode := &usrdb.VerificationCode{
		Send:     email,
		Category: usrdb.EMAIL,
		Code:     num,
		Status:   usrdb.PCS_NORMAL,
		Time:     util.Now(),
		Type:     types,
	}
	return hs.MsgRes(util.Map{"emailCode": emailCode})
}

//绑定邮箱号码
//通过用户对象的相关字段绑定邮箱号码
//@url,需求登录，Get请求
//	~/usr/api/bindEmail		Get
//@arg,json对象中的参数以及query中参数
//	types	R	BIND=1,MODIFYBIND=2,REGISTER=3,VERIFY=4,VERIFYBIND=5,RESETPWD=6,LOGIN=7,UNBIND=8
//	email	R	邮箱号,eg:aa11@aa.com
//	emailOld	R	邮箱号码,eg:eg:aa11@aa.com
//	ecode	O	邮箱验证码的值
/*
	样例	~/usr/api/bindemail?token=xx&email=xx&emailOld=xx&types=xx&ecode=xx"	1为登录用户绑定邮箱号码,2为登录用户改绑邮箱号码
*/
//@ret,返回通用code/data
//	code	I	0：绑定成功，1：参数错误，2：json body错误，
/*	样例
{
    "code":0,
    "data":{
        "data":"OK"
    }
}
*/
//@tag,用户,邮箱
//@author,zhnagyq,2017-05-11
//@case,demo1
func BindEmail(hs *routing.HTTPSession) routing.HResult {
	var email, emailOld string
	var ecode, login, t int
	err := hs.ValidCheckVal(`
		ecode,R|I,R:0;
		email,R|S,L:0;
		emailOld,O|S,L:0;
		types,R|I,R:0;
		login,O|I,R:-1;
		`, &ecode, &email, &emailOld, &t, &login)
	if err != nil {
		return hs.MsgResErr2(1, "参数错误", err)
	}
	uid := hs.StrVal("uid")
	log.D("BindUseremail receive email(%v), emailOld(%v), ecode(%v), t(%v), uid(%v), token(%v), do_login(%v)",
		email, emailOld, ecode, t, uid, hs.StrVal("token"), login)

	if uid == "" && (t == usrdb.BIND || t == usrdb.MODIFYBIND || t == usrdb.UNBIND) {
		return hs.MsgResE(6, "please login first")
	}

	code, err := usrdb.VerifyCode(email, usrdb.EMAIL, ecode, 300000)
	if err != nil {
		log.E("BindUseremail verify code email(%v) ecode(%v) err(%v)", email, ecode, err)
		return hs.MsgResE(code, err.Error())
	}

	var msg string
	var u *usrdb.Usr
	if t == usrdb.MODIFYBIND {
		code, msg, err = usrdb.ChangeBindEmail(uid, email, emailOld)
	} else if t == usrdb.BIND {
		code, msg, err = usrdb.BindEmail(uid, email)
	} else if t == usrdb.UNBIND {
		code, msg, err = usrdb.UnBindEmail(uid, email)
	} else {
		return hs.MsgResErr(1, "参数错误", util.Err("type(%v) err", t))
	}

	if err != nil {
		log.E("Bindemail uid(%v) err(%v)", uid, err)
		return hs.MsgResErr2(code, msg, err)
	}

	if login > 0 {
		token, err := do_login(hs, u)
		if err != nil {
			log.E("BindEmail Do_login uid(%v) err(%v)", u.Id, err)
			return hs.MsgResErr2(11, "验证成功但登录失败", err)
		}
		return hs.MsgRes(util.Map{"token": token, "usr": u})
	}

	return hs.MsgRes(util.Map{"data": "OK"})
}

//create token
func NewToken() string {
	return strings.ToUpper(bson.NewObjectId().Hex())
}

//do login
/*
@arg
	hs:http封装的httpSession对象
	user:用户信息对象
@desc
	用户登录,保存在session数据库中,在/msb/msb.go的Session实现
@ret
	token: 返回用户的有效token
	err: 登录的错误信息
@author
	zhangyq modify on 2017-05-10
*/
func do_login(hs *routing.HTTPSession, user *usrdb.Usr) (string, error) {
	token := NewToken()
	err := usrdb.AddUserSession(user.Id, token)
	if err != nil {
		return "", err
	}
	return token, hs.Flush()
}

//the user logined filter
func LoginFilter(hs *routing.HTTPSession) routing.HResult {
	var token_res string
	err := hs.ValidF(`
	token,R|S,L:0;
	`, &token_res)
	if err != nil {
		log.E("arg-err,%v", err)
		return hs.MsgResErr(1, "arg-err", err)
	}
	session, err := usrdb.FindUserSession(token_res)
	if err != nil {
		return hs.MsgResErr(7, "query-err", err)
	}
	//log.D("session: %v",util.S2Json(session))
	token := session.Token
	hs.SetVal("uid", session.Uid)
	//log.D("token,token_req %v, %v", token, token_res)
	if len(token) < 1 || token != token_res {
		return hs.MsgResE3(301, "arg-err", "not login")
	} else {
		return routing.HRES_CONTINUE
	}
}

//USR_PUB_GRP 允许用户更新的属性组
var USR_PUB_GRP = map[string]int{
	"nickname":   1,
	"sex":        1,
	"age":        1,
	"birthday":   1,
	"location":   1,
	"hometown":   1,
	"profession": 1,
}

func FilterUserAttrs(usr *usrdb.Usr) {
	for key := range usr.Attrs {
		if _, ok := USR_PUB_GRP[key]; ok {
			continue
		}
		log.W("receive unsupported attribute group(%v),do delete", key)
		delete(usr.Attrs, key)
	}
}

func token_url(redirect string, token string) string {
	var url string = ""
	if ok, _ := regexp.MatchString("^.*\\?.*$", redirect); ok {
		url = fmt.Sprintf("%s&token=%s", redirect, token)
	} else {
		url = fmt.Sprintf("%s?token=%s", redirect, token)
	}
	return url
}
