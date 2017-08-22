package main

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"os"
	"runtime"
	"strings"
	"yule"
)

func main() {
	wd, _ := os.Getwd()
	ulimit, _ := util.Exec("ulimit", "-n")
	umask, _ := util.Exec("umask")
	fmt.Sprintf("Env:\n\twd:%v\n\tulimit:%v\n\tumask:%v\n\n",
		wd, strings.TrimSpace(ulimit), strings.TrimSpace(umask))
	runtime.GOMAXPROCS(util.CPU())

	var conf = "yule.properties"
	if len(os.Args) > 1 {
		conf = os.Args[1]
	}
	var cfg = util.NewFcfg3()

	cfg.InitWithFilePath2(conf, true)
	cfg.Print()
	log.RedirectV(cfg.Val2("YULE_OUT_L", ""), cfg.Val2("YULE_ERR_L", ""), false)

	INT, _ := routing.NewJsonINT("conf/")
	INT.Default = "zh"
	routing.Shared.INT = INT

	//uap.Cfg = cfg
	//uap.StartLoopCacheTimeout()
	//err := uap.StartRunner()
	//if err != nil {
	//	log.E("uap StartRunner error:%v", err)
	//	fmt.Println(err.Error())
	//	return
	//}

	yule.RunSrv(cfg)
}
