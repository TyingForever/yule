package usrdb

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"testing"
)

func TestQuerySequence(t *testing.T) {
	_, uid, err := NewUid()
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(uid) < 1 {
		t.Error("error")
		return
	}
	fmt.Println("uid->", uid)

	//db err
	mgo.Mock = true
	mgo.SetMckC("Query-Apply", 0)
	_, _, err = NewUid()
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()
}
