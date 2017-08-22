package yule

import (
	"fmt"
	"github.com/Centny/dbm/mgo"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/filter"
	"github.com/Centny/gwf/util"
	tmgo "gopkg.in/mgo.v2"
	"yule/api/userapi"
	"yule/db"
)

func RunSrv(fcfg *util.Fcfg) {
	defer func() {
		log.Flush()
	}()

	fcfg.Print()

	var err error
	//初始化db
	dbCon := "yule:123@loc.m:27017/yule"
	mgo.AddDefault2(dbCon)
	db.C = mgo.C

	//创建索引
	//err = mgo.ChkIdx(mgo.C, db.AnswerIndexes)

	monitor := filter.NewMonitorH()
	cors := filter.NewCORS_All()
	routing.HFilter("^/.*$", cors)
	routing.HFilterFunc("^/.*$", filter.NoCacheFilter)
	routing.Shared.StartMonitor()

	userapi.Hand("", routing.Shared)

	monitor.AddMonitor("mgo", tmgo.M)
	monitor.AddMonitor("http", routing.Shared)
	routing.H("^/adm/status(\\?.*)?", monitor)
	routing.Shared.Print()
	routing.Shared.ShowLog = true
	//host := "192.168.191.5:"
	host:= "192.168.10.147:"
	listenPort := host + "8080"
	log.I("running web server on %v", listenPort)

	err = routing.ListenAndServe(listenPort)
	if err != nil {
		fmt.Println("RunAnswer listen serve err")
		log.E("RunAnswer listen serve err(%v)", err)
		return
	}
}
