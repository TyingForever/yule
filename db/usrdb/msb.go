package usrdb

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2/bson"
	"yule/db"
)

func AddUserSession(uid, token string) error {
	session := &Session{
		Id:    bson.NewObjectId().Hex(),
		Uid:   uid,
		Token: token,
		Time:  util.Now(),
	}
	return db.C(CN_SESSION).Insert(session)
}

func FindUserSession(token string) (Session, error) {
	var session Session
	err := db.C(CN_SESSION).Find(bson.M{"token": token}).One(&session)
	if ShowLog {
		log.D("session: %v", util.S2Json(session))
	}
	return session, err
}

func UpdateUserSession(uid, token string) error {
	return db.C(CN_SESSION).Update(bson.M{"uid": uid}, bson.M{"$set": bson.M{"token": token, "time": util.Now()}})
}

//清除缓存
func Logout(uid string) error {
	_, err := db.C(CN_SESSION).RemoveAll(bson.M{"uid": uid})
	return err
}
