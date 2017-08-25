package userapi

import (
	"testing"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	"yule/db/usrdb"
	"fmt"
	"yule/db/landlordsdb"
	"time"
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
	//创建游戏
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
	//加入房间
	res,err=DoStartNewLandlords(token[1],operate,categoryId)
	if err != nil {
		log.E(TAG_LANDLORDS_TEST+"err(%v)",err)
		t.Error(err)
		return
	}
	if LOG_API_TEST {
		log.D(TAG_LANDLORDS_TEST+"DoStartNewLandlords: %v",util.S2Json(res))
	}
	res,err=DoStartNewLandlords(token[2],operate,categoryId)
	if err != nil {
		log.E(TAG_LANDLORDS_TEST+"err(%v)",err)
		t.Error(err)
		return
	}
	if LOG_API_TEST {
		log.D(TAG_LANDLORDS_TEST+"DoStartNewLandlords: %v",util.S2Json(res))
	}

	//游戏开始
	//u1明牌
	res,err=DoOperateLandlords(token[0],res.MapVal("landlordInfo").StrVal("lid"),
		landlordsdb.OP_MING,landlordsdb.LM_TD_BEFORE_DEAL,0,0,"")
	if err != nil {
		log.E(TAG_LANDLORDS_TEST+"err(%v)",err)
		t.Error(err)
		return
	}
	if LOG_API_TEST {
		log.D(TAG_LANDLORDS_TEST+"DoOperateLandlords: %v",util.S2Json(res))
	}
	time.Sleep(1*time.Second)

	//发牌
	res,err=DoOperateLandlords(token[0],res.MapVal("landlordInfo").StrVal("lid"),
		landlordsdb.OP_DEAL,0,0,0,"")
	if err != nil {
		log.E(TAG_LANDLORDS_TEST+"err(%v)",err)
		t.Error(err)
		return
	}
	if LOG_API_TEST {
		log.D(TAG_LANDLORDS_TEST+"DoOperateLandlords: %v",util.S2Json(res))
	}
	//抢地主1 u1
	turn_user_token:= ""
	for i := 0; i < 3; i++ {
		if fmt.Sprintf("u%v",i+1)==res.MapVal("landlordInfo").StrVal("turn_user") {
			turn_user_token = token[i]
		}
	}
	res,err=DoOperateLandlords(turn_user_token,res.MapVal("landlordInfo").StrVal("lid"),
		landlordsdb.OP_GRAB,0,0,1,"")
	if err != nil {
		log.E(TAG_LANDLORDS_TEST+"err(%v)",err)
		t.Error(err)
		return
	}
	if LOG_API_TEST {
		log.D(TAG_LANDLORDS_TEST+"DoOperateLandlords: %v",util.S2Json(res))
	}
	//抢地主2
	for i := 0; i < 3; i++ {
		if fmt.Sprintf("u%v",i+1)==res.MapVal("landlordInfo").StrVal("turn_user") {
			turn_user_token = token[i]
		}
	}
	res,err=DoOperateLandlords(turn_user_token,res.MapVal("landlordInfo").StrVal("lid"),
		landlordsdb.OP_GRAB,0,0,2,"")
	if err != nil {
		log.E(TAG_LANDLORDS_TEST+"err(%v)",err)
		t.Error(err)
		return
	}
	if LOG_API_TEST {
		log.D(TAG_LANDLORDS_TEST+"DoOperateLandlords: %v",util.S2Json(res))
	}
	//抢地主0分
	for i := 0; i < 3; i++ {
		if fmt.Sprintf("u%v",i+1)==res.MapVal("landlordInfo").StrVal("turn_user") {
			turn_user_token = token[i]
		}
	}
	res,err=DoOperateLandlords(turn_user_token,res.MapVal("landlordInfo").StrVal("lid"),
		landlordsdb.OP_GRAB,0,0,0,"")
	if err != nil {
		log.E(TAG_LANDLORDS_TEST+"err(%v)",err)
		t.Error(err)
		return
	}
	if LOG_API_TEST {
		log.D(TAG_LANDLORDS_TEST+"DoOperateLandlords: %v",util.S2Json(res))
	}
	//加倍
	for i := 0; i < 3; i++ {
		if fmt.Sprintf("u%v",i+1)==res.MapVal("landlordInfo").StrVal("turn_user") {
			turn_user_token = token[i]
		}
	}
	res,err=DoOperateLandlords(turn_user_token,res.MapVal("landlordInfo").StrVal("lid"),
		landlordsdb.OP_DOUBLE,0,1,0,"")
	if err != nil {
		log.E(TAG_LANDLORDS_TEST+"err(%v)",err)
		t.Error(err)
		return
	}
	if LOG_API_TEST {
		log.D(TAG_LANDLORDS_TEST+"DoOperateLandlords: %v",util.S2Json(res))
	}
	//加倍
	for i := 0; i < 3; i++ {
		if fmt.Sprintf("u%v",i+1)==res.MapVal("landlordInfo").StrVal("turn_user") {
			turn_user_token = token[i]
		}
	}
	res,err=DoOperateLandlords(turn_user_token,res.MapVal("landlordInfo").StrVal("lid"),
		landlordsdb.OP_DOUBLE,0,1,0,"")
	if err != nil {
		log.E(TAG_LANDLORDS_TEST+"err(%v)",err)
		t.Error(err)
		return
	}
	if LOG_API_TEST {
		log.D(TAG_LANDLORDS_TEST+"DoOperateLandlords: %v",util.S2Json(res))
	}
	//不加倍
	for i := 0; i < 3; i++ {
		if fmt.Sprintf("u%v",i+1)==res.MapVal("landlordInfo").StrVal("turn_user") {
			turn_user_token = token[i]
		}
	}
	res,err=DoOperateLandlords(turn_user_token,res.MapVal("landlordInfo").StrVal("lid"),
		landlordsdb.OP_DOUBLE,0,0,0,"")
	if err != nil {
		log.E(TAG_LANDLORDS_TEST+"err(%v)",err)
		t.Error(err)
		return
	}
	if LOG_API_TEST {
		log.D(TAG_LANDLORDS_TEST+"DoOperateLandlords: %v",util.S2Json(res))
	}

	log.D("----------出牌----------")
	//出牌
	for {
		for j := 0; j < 3; j++ {
			if fmt.Sprintf("u%v",j+1)==res.MapVal("landlordInfo").StrVal("turn_user") {
				turn_user_token = token[j]
			}
		}
		res,err=DoOperateLandlords(turn_user_token,res.MapVal("landlordInfo").StrVal("lid"),
			landlordsdb.OP_GET_NOTE,0,0,0,"")
		if err != nil {
			t.Error(err)
			return
		}
		notePopCards:= res.MapVal("landlordInfo").StrVal("note_pop_cards")
		log.D("player:%v,noteCards:%v",turn_user_token,notePopCards)
		if len(notePopCards) >0&&notePopCards!=landlordsdb.LC_PASS {
			res,err=DoOperateLandlords(turn_user_token,res.MapVal("landlordInfo").StrVal("lid"),
				landlordsdb.OP_POP_CARD,0,0,0,notePopCards)
			if err != nil {
				t.Error(err)
				return
			}
			log.D("player:%v,popCards:%v", turn_user_token,  notePopCards)
			if res.MapVal("landlordInfo").IntVal("status") == landlordsdb.LS_SHOW_RESULT {
				log.D("Game over: %v", util.S2Json(res))
				break
			}
		} else {
			res,err=DoOperateLandlords(turn_user_token,res.MapVal("landlordInfo").StrVal("lid"),
				landlordsdb.OP_PASS_CARD,0,0,0,notePopCards)
			if err != nil {
				t.Error(err)
				return
			}
			log.D("player:%v,popCards:%v", turn_user_token,  notePopCards)
		}
	}

}
