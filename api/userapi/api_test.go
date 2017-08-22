package userapi

import (
	"github.com/Centny/dbm/mgo"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2/bson"
	"os"
	"runtime"
	"w.gdy.io/dyf/uas"
	"w.gdy.io/dyf/uas/uap"
	"w.gdy.io/dyf/uas/ucs/sso"
	"yule/db"
	"yule/db/usrdb"
	"yule/db/landlordsdb"
)

var ts *httptest.Server

func init() {
	func() {
		defer func() {
			recover()
		}()
		SrvAddr()
	}()

	os.Setenv("N_RCP_ALL_CFG", "")
	runtime.GOMAXPROCS(util.CPU())

	ts = httptest.NewMuxServer()
	rc := uap.NewRCAuthFilterV("http://127.0.0.1", "")
	rc.Optioned = true
	ts.Mux.HFilter("^.*$", rc)
	Hand("", ts.Mux)

	cfg := util.NewFcfg3()
	cfg.SetVal("SSO_LOGIN_PRE", "http://127.0.0.1/sso/index.html?url=%v")
	cfg.SetVal("SSO_LOGIN_URL", "http://127.0.0.1:7808/sso/api/auth?token=%v")

	sso.Cfg = cfg
	uap.Cfg = cfg

	err := uas.StartTestSrv("yule:123@loc.m:27017/yule", "cny", ":1879", ":1878")
	if err != nil {
		panic(err)
	}

	//err = order.StartTestSrv("../usr_test.properties")
	//if err != nil {
	//	panic(err)
	//	return
	//}

	SrvAddr = func() string {
		return ts.URL
	}

	db.C = mgo.C
}

func rs() string {
	return bson.NewObjectId().Hex()
}

//import (
//	"github.com/Centny/dbm/mgo"
//	"github.com/Centny/gwf/log"
//	"github.com/Centny/gwf/routing/httptest"
//	"github.com/Centny/gwf/util"
//	"os"
//	"runtime"
//	"w.gdy.io/dyf/uas/uap"
//	"yule/db/usrdb"
//	"testing"
//	"gopkg.in/mgo.v2/bson"
//)
//
//var ts *httptest.Server
//
//func init() {
//	func() {
//		defer func() {
//			recover()
//		}()
//		SrvAddr()
//	}()
//
//	os.Setenv("N_RCP_ALL_CFG", "")
//	runtime.GOMAXPROCS(util.CPU())
//
//	ts = httptest.NewMuxServer()
//	rc := uap.NewRCAuthFilterV("http://127.0.0.1", "")
//	//rc.Optioned = true
//	ts.Mux.HFilter("^.*$", rc)
//	Hand("", ts.Mux)
//
//	SrvAddr = func() string {
//		return ts.URL
//	}
//
//	usrdb.C = mgo.C
//	log.D("mgo is %v", util.S2Json(mgo.C))
//}
//
func Remove() {
	db.C(usrdb.CN_USER).RemoveAll(nil)
	db.C(usrdb.CN_SEQUENCE).RemoveAll(nil)
	db.C(usrdb.CN_SESSION).RemoveAll(nil)
	db.C(usrdb.CN_USER).Insert(usrdb.Usr{Id: "u100", Account: "100", Pwd: usrdb.Sha("123"), Role: 1, Time: util.Now()})
	db.C(landlordsdb.CN_LANDLORD_INFO).RemoveAll(nil)
}
