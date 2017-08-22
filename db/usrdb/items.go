package usrdb

import "github.com/Centny/gwf/util"

type Usr struct {
	Id       string   `bson:"_id" json:"id"`                          //the user id
	Account  string   `bson:"account" json:"account,omitempty"`       //the user account default
	Accounts []string `bson:"accounts" json:"accounts,omitempty"`     //the user name list
	Pwd      string   `bson:"pwd" json:"pwd,omitempty"`               //the user password
	Image    string   `bson:"image" json:"image,omitempty"`           //the user image
	Phone    string   `bson:"phone,omitempty" json:"phone,omitempty"` //the user phone bind
	Email    string   `bson:"email,omitempty" json:"email,omitempty"` //the user email bind
	Name     string   `bson:"name" json:"name,omitempty"`             //the user name
	Ip       string   `bson:"ip" json:"ip,omitempty"`                 //the user ip
	Attrs    util.Map `bson:"attrs" json:"attrs,omitempty"`           //the user attribute list split by group
	Status   int      `bson:"status" json:"status,omitempty"`         //the user status
	Type     int      `bson:"type" json:"type,omitempty"`             //the account type 1默认，2手机，3邮箱
	Role     int      `bson:"role" json:"role,omitempty"`             //the  user role
	Last     int64    `bson:"last" json:"last,omitempty"`             //the last updated time
	Time     int64    `bson:"time" json:"time,omitempty"`             //the create time
	Size     int      `bson:"size" json:"size,omitempty"`             //the accounts size
}

//账号，昵称，头像，性别，年龄，生日，所在地，故乡所在地，姓名，职业
//attrs: nickname	sex	age	birthday	location	hometown	profession

type PhoneCode struct {
	Phone  string `bson:"_id"`
	Code   int    `bson:"code"`
	Time   int64  `bson:"time"`
	Type   int    `bson:"type"`
	Status string `bson:"status"`
}

const (
	//phone type
	BIND       = 1
	MODIFYBIND = 2
	REGISTER   = 3
	VERIFY     = 4
	VERIFYBIND = 5
	RESETPWD   = 6
	LOGIN      = 7
	UNBIND     = 8
)
const (
	//phone code status
	PCS_NORMAL = "N"
	PCS_DEL    = "D"
)

type VerificationCode struct {
	Send     string `bson:"_id"`
	Category string `bson:"category"`
	Code     int    `bson:"code"`
	Time     int64  `bson:"time"`
	Type     int    `bson:"type"`
	Status   string `bson:"status"`
}

var (
	//VerificationCode.category
	EMAIL = "email"
	PHONE = "phone"
)

type Sequence struct {
	Id  string `bson:"_id" json:"id"`  //the sequenc id.
	Val uint64 `bson:"val" json:"val"` //the current sequene value
}

const (
	USR_ADMIN    = 1
	USR_ORDINARY = 2
)

const (
	//account type
	ACCOUNT_DEFAULT = 1
	ACCOUNT_PHONE   = 2
	ACCOUNT_EMAIL   = 3
)

const (
	//user status,(not equal 0)
	USR_S_N = 10  //NORMAL
	USR_S_D = -1  //DELETE
	USR_S_Z = -10 //FREEZE
)

const (
	DEFAULT_SEARCH = 1
	NICK_SEARCH    = 2
	PHONE_SEARCH   = 3
	TIME_SEARCH    = 4
)

const (
	DEFAULT_SORT      = 1
	REGISTER_SORT     = 2
	UPDATE_SORT       = 3
	ACCOUNT_SIZE_SORT = 4
	MATCH_SORT        = 5
)

type Session struct {
	Id string `bson:"_id" json:"id,omitempty"` //the session id
	//Kvs   util.Map `bson:"kvs" json:"kvs,omitempty"`   //the session values.
	Uid    string   `bson:"uid" json:"uid,omitempty"`     //the user id
	Token  string   `bson:"token" json:"token,omitempty"` //the  user token
	Last   int64    `bson:"last" json:"last,omitempty"`   //last update time
	Time   int64    `bson:"time" json:"time,omitempty"`   //create time
	update util.Map `bson:"-" json:"-"`                   //updated list, using on flush
}
