package usrdb

import (
	"crypto/sha1"
	"fmt"
	"github.com/Centny/gwf/util"
	"yule/db"
)

var ShowLog bool = true

const (
	//CN_USER = "usr_user"
	CN_USER              = "usr_user"
	CN_SEQUENCE          = "usr_sequence"
	CN_SESSION           = "usr_session"
	CN_PHONE             = "usr_phone"
	CN_VERIFICATION_CODE = "usr_verification_code"
)

//the password encrypt func
var Sha = func(val string) string {
	return util.Sha1_b([]byte(val))
}

var Sha_err = func(val string) string {
	return fmt.Sprintf("%x", sha1.New().Sum([]byte(val)))
}

var Md5 = func(val string) string {
	return util.Md5_b([]byte(val))
}

func Remove() {
	db.C(CN_USER).RemoveAll(nil)
	db.C(CN_SEQUENCE).RemoveAll(nil)
	db.C(CN_SESSION).RemoveAll(nil)
	db.C(CN_VERIFICATION_CODE).RemoveAll(nil)
	db.C(CN_USER).Insert(Usr{Id: "u100", Account: "100", Pwd: Sha("123"), Role: 1, Time: util.Now()})
}
