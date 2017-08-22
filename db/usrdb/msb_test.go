package usrdb

import (
	"gopkg.in/mgo.v2"
	"testing"
	"yule/db"
)

func TestSession(t *testing.T) {
	db.C(CN_SESSION).RemoveAll(nil)
	//add
	uid := "1"
	token := "123"
	err := AddUserSession(uid, token)
	if err != nil {
		t.Error(err)
		return
	}
	//find
	session, err := FindUserSession(token)
	if err != nil {
		t.Error(err)
		return
	}
	if session.Uid != uid || session.Token != token {
		t.Error("find err")
		return
	}

	//update
	token = "234"
	err = UpdateUserSession(uid, token)
	if err != nil {
		t.Error(err)
	}
	session, err = FindUserSession(token)
	if err != nil {
		t.Error(err)
		return
	}
	if session.Uid != uid || session.Token != token {
		t.Error("update err")
		return
	}
	//remove
	err = Logout(uid)
	if err != nil {
		t.Error(err)
	}
	_, err = FindUserSession(token)
	if err == nil {
		t.Error(err)
		return
	}
	//db err
	mgo.Mock = true
	mgo.SetMckC("Collection-Insert", 0)
	uid = "1"
	token = "123"
	err = AddUserSession(uid, token)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()

	mgo.SetMckC("Collection-Update", 0)
	uid = "1"
	token = "234"
	err = UpdateUserSession(uid, token)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()

	mgo.SetMckC("Query-One", 0)
	_, err = FindUserSession(uid)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()

	mgo.SetMckC("Collection-RemoveAll", 0)
	err = Logout(uid)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()

}
