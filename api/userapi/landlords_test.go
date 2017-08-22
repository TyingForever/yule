package userapi

import (
	"testing"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	"yule/db/usrdb"
	"fmt"
	"strings"
)

var token []string = make([]string,10)
var account0 []string = make([]string,10)
var pwd string = "123456"
func prepareDataR() {
	Remove()
	for i := 0; i < 10; i++ {
		account0[i] = fmt.Sprintf("aaa%v",i+1)
		u:=&usrdb.Usr{
			Account:account0[i],
			Pwd:pwd,
			Type:1,
		}
		_,err:=DoRegister(u,0)
		if err!=nil {
			log.E(TAG_LANDLORDS_TEST+"err(%v)",err)
			return
		}
		res,err:=DoLogin(account0[i],pwd)
		if err!=nil {
			log.E(TAG_LANDLORDS_TEST+"err(%)",err)
			return
		}
		token[i] = res.StrVal("token")
	}
}

func TestDoEntryLandlords(t *testing.T) {
	prepareDataR()

	gameType:=1
	res,err:=DoEntryLandlords(token[0],gameType)
	if err!=nil {
		log.E(TAG_LANDLORDS_TEST+"err(%v)",err)
		t.Error(err)
		return
	}
	if LOG_API_TEST {
		log.D(TAG_LANDLORDS_TEST+"DoEntryLandlords: %v",util.S2Json(res))
	}
	operate:=1
	categoryId:=1
	res,err=DoStartNewLandlords(token[0],operate,categoryId)
	if err != nil {
		log.E(TAG_LANDLORDS_TEST+"err(%v)",err)
		t.Error(err)
		return
	}
	if LOG_API_TEST {
		log.D(TAG_LANDLORDS_TEST+"DoStartNewLandlords: %v",util.S2Json(res))
	}

	lid:="123456"
	res,err=DoGetLandlords(token[0],lid)
	if err!=nil {
		log.E(TAG_LANDLORDS_TEST+"err(%v)",err)
		t.Error(err)
		return
	}
	if LOG_API_TEST {
		log.D(TAG_LANDLORDS_TEST+"DoGetLandlords: %v",util.S2Json(res))
	}

	pop_cards:="1,2,3,4,5"
	res,err=DoOperateLandlords(token[0],lid,operate,pop_cards)
	if err != nil {
		log.E(TAG_LANDLORDS_TEST+"err(%v)",err)
		t.Error(err)
		return
	}
	if LOG_API_TEST {
		log.D(TAG_LANDLORDS_TEST+"DoOperateLandlords: %v",util.S2Json(res))
	}
}
