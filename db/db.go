package db

import (
	"gopkg.in/mgo.v2"
)

var C = func(name string) *mgo.Collection {
	panic("the yule database collection handler function is not initial")
}
