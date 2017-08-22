package usrdb

import (
	"gopkg.in/mgo.v2"
	"testing"
	"time"
)

func TestA(t *testing.T) {
	PrepareData()
}

func TestBind(t *testing.T) {
	PrepareData()
	uid := "u100"
	uid0 := "u3"
	uid1 := "u4"
	email1 := "843232800@qq.com"
	email2 := "843232801@qq.com"
	email3 := "843232802@qq.com"
	email_err := "138001384"

	//bind
	re, s, err := BindEmail(uid, email1)
	if re != 0 || s != "" || err != nil {
		t.Error(err)
		return
	}
	re, s, err = BindEmail(uid0, email3)
	if re != 0 || s != "" || err != nil {
		t.Error(err)
		return
	}
	//err email err
	re, s, err = BindEmail(uid, email_err)
	if re == 0 && s == "" && err == nil {
		t.Error(err)
		return
	}
	//err this user has been bound email
	re, s, err = BindEmail(uid, email1)
	if re == 0 && s == "" && err == nil {
		t.Error(err)
		return
	}
	//err this email has been bound
	re, s, err = BindEmail(uid1, email1)
	if re == 0 && s == "" && err == nil {
		t.Error(err)
		return
	}

	//changeBind
	//err email err
	re, s, err = ChangeBindEmail(uid, email_err, email1)
	if re == 0 && s == "" && err == nil {
		t.Error(err)
		return
	}
	//err email same
	re, s, err = ChangeBindEmail(uid, email1, email1)
	if re == 0 && s == "" && err == nil {
		t.Error(err)
		return
	}
	//err old email wrong
	re, s, err = ChangeBindEmail(uid1, email2, email3)
	if re == 0 && s == "" && err == nil {
		t.Error(err)
		return
	}
	//err this email has been bound
	re, s, err = ChangeBindEmail(uid, email3, email1)
	if re == 0 && s == "" && err == nil {
		t.Error(err)
		return
	}

	//server err
	re, s, err = ChangeBindEmail("xxx", email2, email1)
	if re == 0 && s == "" && err == nil {
		t.Error(err)
		return
	}
	//ok
	re, s, err = ChangeBindEmail(uid, email2, email1)
	if re != 0 || s != "" || err != nil {
		t.Error(err)
		return
	}

	//unbind
	//err email err
	re, s, err = UnBindEmail(uid, email_err)
	if re == 0 && s == "" && err == nil {
		t.Error(err)
		return
	}
	//err no email
	re, s, err = UnBindEmail(uid1, "a@b.c")
	if re == 0 && s == "" && err == nil {
		t.Error(err)
		return
	}

	//unbind
	re, s, err = UnBindEmail(uid, email2)
	if re != 0 || s != "" || err != nil {
		t.Error(err)
		return
	}

	phone1 := "13800138000"
	phone2 := "13800138001"
	phone3 := "13800138002"
	phone_err := "138001384"

	//bind
	re, s, err = BindPhone(uid, phone1)
	if re != 0 || s != "" || err != nil {
		t.Error(err)
		return
	}
	re, s, err = BindPhone(uid0, phone3)
	if re != 0 || s != "" || err != nil {
		t.Error(err)
		return
	}
	//err phone err
	re, s, err = BindPhone(uid, phone_err)
	if re == 0 && s == "" && err == nil {
		t.Error(err)
		return
	}
	//err this user has been bound phone
	re, s, err = BindPhone(uid, phone1)
	if re == 0 && s == "" && err == nil {
		t.Error(err)
		return
	}
	//err this phone has been bound
	re, s, err = BindPhone(uid1, phone1)
	if re == 0 && s == "" && err == nil {
		t.Error(err)
		return
	}

	//changeBind
	//err phone err
	re, s, err = ChangeBindPhone(uid, phone_err, phone1)
	if re == 0 && s == "" && err == nil {
		t.Error(err)
		return
	}
	//err phone same
	re, s, err = ChangeBindPhone(uid, phone1, phone1)
	if re == 0 && s == "" && err == nil {
		t.Error(err)
		return
	}
	//err old phone wrong
	re, s, err = ChangeBindPhone(uid1, phone2, phone3)
	if re == 0 && s == "" && err == nil {
		t.Error(err)
		return
	}
	//err this phone has been bound
	re, s, err = ChangeBindPhone(uid, phone3, phone1)
	if re == 0 && s == "" && err == nil {
		t.Error(err)
		return
	}

	//server err
	re, s, err = ChangeBindPhone("xxx", phone2, phone1)
	if re == 0 && s == "" && err == nil {
		t.Error(err)
		return
	}
	//ok
	re, s, err = ChangeBindPhone(uid, phone2, phone1)
	if re != 0 || s != "" || err != nil {
		t.Error(err)
		return
	}

	//unbind
	//err phone err
	re, s, err = UnBindPhone(uid, phone_err)
	if re == 0 && s == "" && err == nil {
		t.Error(err)
		return
	}
	//err no phone
	re, s, err = UnBindPhone(uid1, phone1)
	if re == 0 && s == "" && err == nil {
		t.Error(err)
		return
	}

	//unbind
	re, s, err = UnBindPhone(uid, phone2)
	if re != 0 || s != "" || err != nil {
		t.Error(err)
		return
	}

	//verifyCode
	send := "13800138010"
	send1 := "13800138011"
	//send2:="13800138012"

	category := "phone"
	code := 123456
	validTime := 60 * time.Second.Nanoseconds()
	Type := 1
	re, err = UpsertCode(send, category, code, validTime, Type)
	if re != 0 || err != nil {
		t.Error(err)
		return
	}
	//err time
	re, err = UpsertCode(send, category, code, validTime, Type)
	if re == 0 && err == nil {
		t.Error(err)
		return
	}

	//verify
	//err not found
	re, err = VerifyCode(send1, category, code, validTime)
	if re == 0 && err == nil {
		t.Error(err)
		return
	}
	//ok

	re, err = VerifyCode(send, category, code, validTime)
	if re != 0 || err != nil {
		t.Error(err)
		return
	}
	//code has been verified
	re, err = VerifyCode(send, category, code, validTime)
	if re == 0 && err == nil {
		t.Error(err)
		return
	}

	//db err
	//email
	mgo.Mock = true
	mgo.SetMckC("Collection-Update", 0)
	_, _, err = BindEmail(uid, email2)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()
	mgo.SetMckC("Collection-Update", 0)
	_, _, err = ChangeBindEmail(uid, email1, email2)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()
	mgo.SetMckC("Collection-Update", 0)
	_, _, err = UnBindEmail(uid, email1)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()

	//db phone err
	mgo.SetMckC("Collection-Update", 0)
	_, _, err = BindPhone(uid, phone2)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()

	mgo.SetMckC("Collection-Update", 0)
	_, _, err = ChangeBindPhone(uid, phone1, phone2)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()

	mgo.SetMckC("Collection-Update", 0)
	_, _, err = UnBindPhone(uid, phone1)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()

	//db upsertCode
	mgo.SetMckC("Query-Apply", 0)
	_, err = UpsertCode(send, category, code, validTime, Type)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()
	mgo.SetMckC("Query-Apply", 0)
	_, err = VerifyCode(send, category, code, validTime)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()

	//db index
	err = CreateIndex("", false)
	if err == nil {
		t.Error(err)
		return
	}

}
