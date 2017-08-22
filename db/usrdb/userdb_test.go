package usrdb

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"strings"
	"testing"
	//"time"
	"fmt"
	"github.com/anacrolix/sync"
	"yule/db"
)

func TestUser(t *testing.T) {
	Remove()
	//register
	pwd := "123"
	account := "a1"
	user := &Usr{
		Account: account,
		Pwd:     pwd,
		Type:    1,
	}
	err := AddUser(user)
	if err != nil {
		t.Error(err)
		return
	}
	//register err account ""
	user_err := &Usr{}
	user_err.Account = ""
	user.Pwd = "123"
	err = AddUser(user_err)
	if err == nil {
		t.Error(err)
		return
	}
	//register err pwd ""
	user_err = &Usr{}
	user_err.Account = "123456"
	user.Pwd = ""
	err = AddUser(user_err)
	if err == nil {
		t.Error(err)
		return
	}
	//register user exist
	user_err = &Usr{}
	user_err.Account = "100"
	user.Pwd = "123"
	err = AddUser(user_err)
	if err == nil {
		t.Error(err)
		return
	}

	//login
	user_new, err := FindUserByAccountPwd(account, pwd)
	if err != nil {
		t.Error(err)
		return
	}
	if user_new.Account != user.Account {
		t.Error("login err")
		return
	}

	//err account not exist
	_, err = FindUserByAccountPwd("b1", pwd)
	if err == nil {
		t.Error(err)
		return
	}

	_, err = FindUserByAccountPwd("100", "abc")
	if err == nil {
		t.Error(err)
		return
	}

	//update
	user_update := &Usr{
		Account: "account",
		Pwd:     "123456789",
		Status:  USR_S_Z,
		Attrs: util.Map{
			"location": "Guangzhou",
			"hometown": "Xinyi",
		},
	}
	user_update.Id = "u1"
	err = UpdateUserInfo(user_update)
	if err != nil {
		t.Error(err)
		return
	}
	if user_update.Attrs.StrVal("location") != "Guangzhou" || user_update.Attrs.StrVal("hometown") != "Xinyi" {
		t.Errorf("update error: %v", util.S2Json(user_new))
		return
	}

	//update err uid=""
	user_err.Id = ""
	user_err.Account = ""
	user_err.Pwd = ""
	err = UpdateUserInfo(user_err)
	if err == nil {
		t.Error(err)
		return
	}
	//err account has been existed
	user_err.Id = "u1"
	user_err.Account = "100"
	user_err.Pwd = "123"
	err = UpdateUserInfo(user_err)
	if err == nil {
		t.Error(err)
		return
	}

	//change pwd
	user_err.Id = "u1"
	user_err.Account = ""
	user_err.Pwd = "1234"
	err = UpdateUserInfo(user_err)
	if err != nil {
		t.Error(err)
		return
	}

	//getUserInfo
	uid := "u1"
	_, err = GetUserInfo(uid)
	if err != nil {
		t.Error(err)
		return
	}
	//err user not esist
	uid_err := "q1"
	_, err = GetUserInfo(uid_err)
	if err == nil {
		t.Error(err)
		return
	}

	//findUsers
	query := bson.M{"_id": "u1"}
	selector := bson.M{"phone": 1}
	_, err = FindUsers(query, selector)
	if err != nil {
		t.Error(err)
		return
	}

	//db err
	//register
	account = "adb"
	pwd = "123"
	user_db := &Usr{
		Account: account,
		Pwd:     pwd,
		Type:    1,
	}
	mgo.Mock = true
	mgo.SetMckC("Query-Apply", 0)
	err = AddUserV(user_db)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()

	account = "adb"
	pwd = "123"
	user_db1 := &Usr{
		Account: account,
		Pwd:     pwd,
		Type:    1,
	}
	mgo.SetMckC("Query-Apply", 0)
	err = AddUser(user_db1)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()

	account = "100"
	pwd = "123"
	user_db2 := &Usr{
		Account: account,
		Pwd:     pwd,
		Type:    1,
	}
	mgo.SetMckC("Query-Apply", 1)
	err = AddUserV(user_db2)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()

	//login
	mgo.SetMckC("Query-All", 0)
	account = "1"
	_, err = FindUserByAccountPwd(account, pwd)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()
	//update
	mgo.SetMckC("Query-All", 0)
	mgo.SetMckC("Collection-Update", 0)
	err = UpdateUserInfo(user_new)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()
	//find
	mgo.SetMckC("Query-All", 0)
	_, err = GetUserInfo(user_new.Id)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()
	//findUsers
	mgo.SetMckC("Query-All", 0)
	_, err = FindUsers(query, selector)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()

}

func TestAdmin(t *testing.T) {
	Remove()

	var u *Usr
	var err error
	u = &Usr{
		Account: "2",
		Pwd:     "123",
		Type:    1,
	}
	for i := 0; i < 10; i++ {
		u.Account += "1"
		u.Pwd = "123"
		err = AddUser(u)
		if err != nil {
			t.Error(err)
			return
		}
	}
	//usersList
	total, users, err := ListUsersOrdinary("u100", 1, 3)
	if err != nil {
		t.Error(err)
	}
	if len(users) != 3 || total != 10 {
		t.Errorf("listUsers err %v", util.S2Json(users))
	}

	//err uid0=""
	uid_err := ""
	total, users, err = ListUsersOrdinary(uid_err, 1, 3)
	if total != 0 && users != nil && err == nil {
		t.Error(err)
	}
	//err uid0 not exist
	uid_err = "ads"
	total, users, err = ListUsersOrdinary(uid_err, 1, 3)
	if total != 0 && users != nil && err == nil {
		t.Error(err)
	}
	//err uid0 exist but not admin
	uid_err = "u1"
	total, users, err = ListUsersOrdinary(uid_err, 1, 3)
	if total != 0 && users != nil && err == nil {
		t.Error(err)
	}

	//err user nil
	u_err := &Usr{}
	total, users, err = ListUserV(u_err, false, false, bson.M{}, 0, 0, 0, 0, 0, 0)
	if total != 0 && users != nil && err == nil {
		t.Error(err)
	}

	//db err list
	mgo.Mock = true
	mgo.SetMckC("Query-All", 0)
	//mgo.SetMckC("Query-All", 1)
	_, _, err = ListUsersOrdinary("u100", 1, 3)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()
	mgo.SetMckC("Query-All", 1)
	//mgo.SetMckC("Query-All", 1)
	_, _, err = ListUsersOrdinary("u100", 1, 3)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()
	mgo.SetMckC("Query-Count", 0)
	//mgo.SetMckC("Query-All", 1)
	_, _, err = ListUsersOrdinary("u100", 1, 3)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()

}

func TestSearchUsersOrdinary(t *testing.T) {
	var err error
	PrepareData()
	db.C(CN_SESSION).Remove(nil)
	var uid = "u100"
	var searchMethod, skip, limit, sort int
	var startTime, overTime int64
	var nickKey, phoneKey string
	skip = 0
	limit = 5
	searchMethod, sort = DEFAULT_SEARCH, DEFAULT_SORT
	startTime, overTime = 0, 0
	//1.NICK_SEARCH DEFAULT_SORT
	searchMethod = NICK_SEARCH
	nickKey = "q"
	sort = DEFAULT_SORT
	total, users, err := SearchUsersOrdinary(uid, nickKey, phoneKey,
		searchMethod, sort, skip, limit, startTime, overTime)
	if err != nil {
		t.Error(err)
		return
	}
	log.D("total:%v", total)
	if total != 10 {
		t.Errorf("err users(1.NICK_SEARCH DEFAULT_SORT): %v", util.S2Json(users))
		return
	}
	if len(users) != 5 {
		t.Errorf("err users(1.NICK_SEARCH DEFAULT_SORT): %v", util.S2Json(users))
		return
	}
	for i := 0; i < 5; i++ {
		if !strings.Contains(users[i].Attrs.StrVal("nickname"), nickKey) {
			t.Errorf("err users(1.NICK_SEARCH DEFAULT_SORT): %v", util.S2Json(users))
			return
		}
	}
	t.Logf("users(1.NICK_SEARCH DEFAULT_SORT): %v", util.S2Json(users))

	//err uid ""
	uid0 := ""
	_, _, err = SearchUsersOrdinary(uid0, nickKey, phoneKey,
		searchMethod, sort, skip, limit, startTime, overTime)
	if err == nil {
		t.Error(err)
		return
	}
	//err uid not exist
	uid0 = "abs"
	_, _, err = SearchUsersOrdinary(uid0, nickKey, phoneKey,
		searchMethod, sort, skip, limit, startTime, overTime)
	if err == nil {
		t.Error(err)
		return
	}
	//err uid  exist but not admin
	uid0 = "u1"
	_, _, err = SearchUsersOrdinary(uid0, nickKey, phoneKey,
		searchMethod, sort, skip, limit, startTime, overTime)
	if err == nil {
		t.Error(err)
		return
	}

	//db
	mgo.Mock = true
	mgo.SetMckV("Query-All", 0, 1)
	_, _, err = SearchUsersOrdinary(uid, nickKey, phoneKey,
		searchMethod, sort, skip, limit, startTime, overTime)
	if err == nil {
		t.Error(err)
		return
	}
	mgo.ClearMock()

	//2.NICK_SEARCH REGISTER_SORT 降序
	searchMethod = NICK_SEARCH
	nickKey = "q"
	sort = REGISTER_SORT
	total, users, err = SearchUsersOrdinary(uid, nickKey, phoneKey,
		searchMethod, sort, skip, limit, startTime, overTime)
	if err != nil {
		t.Error(err)
		return
	}
	if total != 10 {
		t.Errorf("err users(1.NICK_SEARCH DEFAULT_SORT): %v", util.S2Json(users))
		return
	}
	if len(users) != 5 {
		t.Errorf("err users(2.NICK_SEARCH DEFAULT_SORT): %v", util.S2Json(users))
		return
	}
	for i := 0; i < 5; i++ {
		if !strings.Contains(users[i].Attrs.StrVal("nickname"), nickKey) {
			t.Errorf("err users(2.NICK_SEARCH DEFAULT_SORT): %v", util.S2Json(users))
			return
		}
	}
	for i := 0; i < 4; i++ {
		if users[i].Time < users[i+1].Time {
			t.Errorf("err users(2.NICK_SEARCH REGISTER_SORT 降序): %v", util.S2Json(users))
			return
		}
	}

	t.Logf("users(2.NICK_SEARCH REGISTER_SORT): %v", util.S2Json(users))

	//3.NICK_SEARCH UPDATE_SORT
	searchMethod = NICK_SEARCH
	nickKey = "q"
	sort = UPDATE_SORT
	total, users, err = SearchUsersOrdinary(uid, nickKey, phoneKey,
		searchMethod, sort, skip, limit, startTime, overTime)
	if err != nil {
		t.Error(err)
		return
	}
	if total != 10 {
		t.Errorf("err users(1.NICK_SEARCH DEFAULT_SORT): %v", util.S2Json(users))
		return
	}
	if len(users) != 5 {
		t.Errorf("err users(3.NICK_SEARCH UPDATE_SORT): %v", util.S2Json(users))
		return
	}
	for i := 0; i < 5; i++ {
		if !strings.Contains(users[i].Attrs.StrVal("nickname"), nickKey) {
			t.Errorf("err users(3.NICK_SEARCH UPDATE_SORT): %v", util.S2Json(users))
			return
		}
	}
	for i := 0; i < 4; i++ {
		if users[i].Last < users[i+1].Last {
			t.Errorf("err users(3.NICK_SEARCH UPDATE_SORT 降序): %v", util.S2Json(users))
			return
		}
	}

	t.Logf("users(3.NICK_SEARCH UPDATE_SORT): %v", util.S2Json(users))

	//4.NICK_SEARCH ACCOUNT_SIZE_SORT
	searchMethod = NICK_SEARCH
	nickKey = "q"
	sort = ACCOUNT_SIZE_SORT
	_, users, err = SearchUsersOrdinary(uid, nickKey, phoneKey,
		searchMethod, sort, skip, limit, startTime, overTime)
	if err != nil {
		t.Error(err)
		return
	}
	if len(users) != 5 {
		t.Errorf("err users(4.NICK_SEARCH ACCOUNT_SIZE_SORT): %v", util.S2Json(users))
		return
	}
	for i := 0; i < 5; i++ {
		if !strings.Contains(users[i].Attrs.StrVal("nickname"), nickKey) {
			t.Errorf("err users(4.NICK_SEARCH ACCOUNT_SIZE_SORT): %v", util.S2Json(users))
			return
		}
	}
	for i := 0; i < 4; i++ {
		if users[i].Size < users[i+1].Size {
			t.Errorf("err users(4.NICK_SEARCH ACCOUNT_SIZE_SORT): %v", util.S2Json(users))
			return
		}
	}

	t.Logf("users(4.NICK_SEARCH ACCOUNT_SIZE_SORT): %v", util.S2Json(users))

	//5. REGISTER_SORT DEFAULT_SORT
	searchMethod = PHONE_SEARCH
	phoneKey = "1"
	sort = REGISTER_SORT
	_, users, err = SearchUsersOrdinary(uid, nickKey, phoneKey,
		searchMethod, sort, skip, limit, startTime, overTime)
	if err != nil {
		t.Error(err)
		return
	}
	if len(users) != 5 {
		t.Errorf("err users(5.PHONE_SEARCH REGISTER_SORT): %v", util.S2Json(users))
		return
	}
	for i := 0; i < 4; i++ {
		if users[i].Time < users[i+1].Time {
			t.Errorf("err users(5.PHONE_SEARCH REGISTER_SORT): %v", util.S2Json(users))
			return
		}
	}
	for i := 0; i < 5; i++ {
		if !strings.Contains(users[i].Phone, phoneKey) {
			t.Errorf("err users(5.PHONE_SEARCH REGISTER_SORT): %v", util.S2Json(users))
			return
		}
	}
	t.Logf("users(5.PHONE_SEARCH REGISTER_SORT): %v", util.S2Json(users))

	//6. TIME_SEARCH REGISTER_SORT
	searchMethod = TIME_SEARCH
	startTime = users[0].Time - 1000000
	overTime = users[0].Time + 1000000
	sort = REGISTER_SORT
	total, users, err = SearchUsersOrdinary(uid, nickKey, phoneKey,
		searchMethod, sort, skip, limit, startTime, overTime)
	if err != nil {
		t.Error(err)
		return
	}
	if total != limit*10 {
		t.Errorf("err users(1.NICK_SEARCH DEFAULT_SORT): %v", util.S2Json(users))
		return
	}
	if len(users) != 5 {
		t.Errorf("err users(6. TIME_SEARCH REGISTER_SORT): %v", util.S2Json(users))
		return
	}
	for i := 0; i < 4; i++ {
		if users[i].Time < users[i+1].Time {
			t.Errorf("err users(6. TIME_SEARCH REGISTER_SORT): %v", util.S2Json(users))
			return
		}
	}
	for i := 0; i < 5; i++ {
		if startTime > users[i].Time || users[i].Time > overTime {
			t.Errorf("err users(6. TIME_SEARCH REGISTER_SORT): %v", util.S2Json(users))
			return
		}
	}

	t.Logf("users(6. TIME_SEARCH REGISTER_SORT): %v", util.S2Json(users))

	//7.NICK_SEARCH MATCH_SORT
	searchMethod = NICK_SEARCH
	nickKey = "q"
	sort = MATCH_SORT
	_, users, err = SearchUsersOrdinary(uid, nickKey, phoneKey,
		searchMethod, sort, skip, limit, startTime, overTime)
	if err != nil {
		t.Error(err)
		return
	}
	if len(users) != 5 {
		t.Errorf("err users(7.NICK_SEARCH MATCH_SORT): %v", util.S2Json(users))
		return
	}
	for i := 0; i < 5; i++ {
		if !strings.Contains(users[i].Attrs.StrVal("nickname"), nickKey) {
			t.Errorf("err users(7.NICK_SEARCH MATCH_SORT): %v", util.S2Json(users))
			return
		}
	}
	t.Logf("users(7.NICK_SEARCH MATCH_SORT): %v", util.S2Json(users))

	for i := 0; i < 4; i++ {
		if len(users[i].Attrs.StrVal("nickname")) > len(users[i+1].Attrs.StrVal("nickname")) {
			t.Errorf("err users(7.NICK_SEARCH MATCH_SORT): %v,%v", len(users[i].Attrs.StrVal("nickname")), len(users[i+1].Attrs.StrVal("nickname")))
			return
		}
	}
}

func PrepareData() {
	Remove()
	u := &Usr{}
	var err error
	u = &Usr{
		Account: "",
		Pwd:     "123",
		Type:    1,
	}
	for i := 0; i < 20; i++ {
		u.Account += "1"
		u.Pwd = "123"
		u.Type = 1
		err = AddUser(u)
		if err != nil {
			log.E("err: %v", err)
			return
		}
	}
	for i := 0; i < 79; i++ {
		u.Account = NewStringLen(20)
		u.Pwd = NewStringLen(6)
		u.Type = 1
		err = AddUser(u)
		if err != nil {
			log.E("err: %v", err)
			return
		}
	}

	uid := []string{"u1", "u2", "u3", "u4", "u5", "u6", "u7", "u8", "u9", "u10", "u11", "u12", "u13", "u14", "u15", "u16", "u17", "u18", "u19"}
	u = &Usr{}
	index := 0
	nickname := ""
	for i := 0; i < 10; i++ {
		u.Email = ""
		u.Id = uid[index]
		index++
		nickname = "q" + NewStringLen(5) + NewDigitLen(10-i)
		if i < 2 {
			u.Email = NewStringLen(5) + "@" + NewStringLen(3) + "." + NewStringLen(3)
		}
		u.Phone = NewDigitLen(11)
		if i == 3 || i == 4 {
			u.Phone = ""
		}
		u.Attrs = util.Map{
			"nickname": nickname,
		}
		UpdateUserInfo(u)
	}
	for i := 0; i < 5; i++ {
		u.Phone = "1" + NewDigitLen(10)
		u.Id = uid[index]
		index++
		nickname = "x" + NewDigitLen(8)
		u.Attrs = util.Map{
			"nickname": nickname,
		}
		UpdateUserInfo(u)
	}
	u.Phone = ""
	for i := 0; i < 4; i++ {
		//u.Phone = NewDigitLen(11)
		u.Id = uid[index]
		index++
		nickname = "y" + NewDigitLen(8)
		u.Attrs = util.Map{
			"nickname": nickname,
		}
		UpdateUserInfo(u)
	}

}

var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
var StdDigit = []byte("0123456789")

func NewDigitLen(length int) string {
	return NewLenChars(length, StdDigit)
}

func NewStringLen(length int) string {
	return NewLenChars(length, StdChars)
}
func NewLenChars(length int, chars []byte) string {
	if length == 0 {
		return ""
	}
	clen := len(chars)
	if clen < 2 || clen > 256 {
		panic("Wrong charset length for NewLenChars()")
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			panic("Error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				continue // Skip this number to avoid modulo bias.
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}

//test data
//func TestData(t *testing.T) {
//	//Remove()
//	//account:= ""
//	//pwd:="123"
//	////管理员
//	//for i := 0; i < 10; i++ {
//	//	_,uid,_ := NewUid()
//	//	account = NewStringLen(10)
//	//	C(CN_USER).Insert(Usr{Id: uid, Account: account, Pwd: Sha(pwd), Role: 1, Time: util.Now()})
//	//}
//	var w sync.WaitGroup
//	for i := 0; i < 100; i++ {
//		w.Add(1)
//		//go DataTest1(w)
//		go DataTest2(w)
//	}
//	w.Wait()
//
//	//DataTest()
//}

//func DataTest( sync.WaitGroup)  {
//	u:=&Usr{}
//	var err error
//	u = &Usr{
//		Account: "",
//		Pwd:     "123",
//		Type:    1,
//	}
//	for i := 0; i < 400; i++ {
//		u.Account = NewStringLen(10)
//		u.Pwd = "123"
//		u.Type = 1
//		err = AddUser(u)
//		if err != nil {
//			log.E("err: %v", err)
//			return
//		}
//	}
//
//}

func DataTest1(w sync.WaitGroup) {
	u := &Usr{}
	var err error
	u = &Usr{
		Account: "",
		Pwd:     "123",
		Type:    1,
	}
	for i := 0; i < 1000; i++ {
		u.Account = NewStringLen(13)
		u.Pwd = "123"
		u.Type = 1
		err = AddUser(u)
		if err != nil {
			log.E("err: %v", err)
			return
		}
	}
	w.Done()
}

func DataTest2(w sync.WaitGroup) {
	var err error
	u := &Usr{}
	for i := 1; i < 400; i++ {

		u.Id = "u" + NewDigitLen(8)
		u.Account = NewStringLen(15)
		u.Type = 1
		u.Pwd = "123"
		err = AddUserV(u)
		if err != nil {
			log.E("err: %v", err)
			return
		}
	}
	w.Done()

}

func TestC(t *testing.T) {
	var uid, account string
	for i := 1000; i < 300000; i++ {
		uid = fmt.Sprintf("u%v", i)
		account = fmt.Sprintf("a%v", i)
		db.C(CN_USER).Insert(Usr{Id: uid, Account: account, Pwd: Sha("123"), Role: 2, Attrs: util.Map{"nickname": uid}, Time: util.Now()})
	}
}
