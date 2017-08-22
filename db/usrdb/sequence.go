package usrdb

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"yule/db"
)

const (
	SEQ_USR = "user"
)

func QuerySequence(id string) (uint64, error) {
	var seq Sequence
	_, err := db.C(CN_SEQUENCE).Find(bson.M{"_id": id}).Apply(
		mgo.Change{
			Update:    bson.M{"$inc": bson.M{"val": 1}},
			Upsert:    true,
			ReturnNew: true,
		}, &seq,
	)
	return seq.Val, err
}

func NewUid() (uint64, string, error) {
	uid, err := QuerySequence(SEQ_USR)
	return uid, fmt.Sprintf("u%v", uid), err
}
