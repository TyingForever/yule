package usrdb

import (
	"fmt"
	"github.com/Centny/dbm/mgo"
	"yule/db"
)

func init() {
	func() {
		defer func() {
			fmt.Println("init recover: ", recover())
		}()
		db.C("fake")
	}()
	mgo.AddDefault2("yule:123@loc.m:27017/yule")
	db.C = mgo.C
}
