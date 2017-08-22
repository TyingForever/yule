package landlordsdb

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2/bson"
	"testing"
	"yule/db"
	"time"
)

/**
开始游戏
1.等待：非轮流，若有空闲房间则加入，否则创建新的游戏房间等待其他用户加入房间。当用户是最后一个加入则启动广播通知其他两个用户开始斗地主。
当玩家在等待时间内仍未开始游戏，则系统自动给改房间的，玩家匹配新的房间，并通知该房间的玩家进入新的房间等待游戏的开始。
2.明牌：非轮流，当开始游戏前5秒有明牌时间，5秒后开始发牌，随着牌数的增加而降低明牌的倍数，需要注意明牌条件，需要根据金额来确定倍率，默认不明牌
3.抢地主：轮流，发完牌则进入抢地主，随机确认先手抢地主，抢地主规则：每人只能叫一次，叫牌1，2，3，不叫。后叫牌者要比前面的分数高或者不叫。
叫牌分数最高者为地主，在叫牌过程中有叫3分者为地主并结束叫牌，若所有玩家都不叫地主则重新开始
4.明牌：非轮流，在首次出牌前明牌，若加倍完成仍未明牌则默认不明牌
5.加倍：非轮流，每人在限定时间进行一次加倍，可以加倍，不加倍，默认不加倍
6.出牌：轮流，地主首先出牌，首轮无上次出牌，则上次类型默认为NONE，若自己的牌在本回合最大，则上次出牌类型为清零NONE
7.结束：当玩家第一个出完牌则结束游戏，玩家角色为赢家，开始清算低分，倍数统计输赢金币，并显示结果一定时间，用户在规定时间可以选择继续游戏或者回到大厅，时间已过默认回到大厅
8.继续游戏：->1
*/
var uid0 []string = []string{"u1", "u2", "u3", "u4", "u5", "u6", "u7", "u8", "u9", "u10"}

func TestOperateLandlordInfo(t *testing.T) {
	Remove()
	cid:=1
	//首个用户，创建新房间
	landlordInfo, err := GetRoomRandom(uid0[0],cid)
	if err != nil {
		t.Error(err)
		return
	}
	log.D("GetRoomRandom create new room %v", util.S2Json(landlordInfo))
	//加入已创建的房间
	landlordInfo, err = GetRoomRandom(uid0[1],cid)
	if err != nil {
		t.Error(err)
		return
	}
	log.D("GetRoomRandom join in room %v", util.S2Json(landlordInfo))
	//加入已创建的房间
	landlordInfo, err = GetRoomRandom(uid0[2], cid)
	if err != nil {
		t.Error(err)
		return
	}
	log.D("GetRoomRandom join in room %v", util.S2Json(landlordInfo))
	//获取房间信息
	land, err := FindRoomV(uid0[0], "", LC_PRIMER_LEVEL,cid, LS_MING, nil)
	if err != nil {
		t.Error(err)
		return
	}
	log.D("FindRoom: %v", util.S2Json(land))
	//u1明牌
	landlord := &LandlordInfo{Id: land.Id}
	info, err := OperateLandlordInfo(uid0[0], "", landlord, OP_MING, LM_TD_BEFORE_DEAL)
	if err != nil {
		t.Error(err)
		return
	}
	log.D("MING:%v", util.S2Json(info))
	//发牌
	time.Sleep(1 * time.Second)
	info, err = OperateLandlordInfo(uid0[0], "", landlord, OP_DEAL, 0)
	if err != nil {
		t.Error(err)
		return
	}
	log.D("Deal :%v", util.S2Json(info))
	//抢地主
	landlord.Multiple = util.Map{"grab_score": 0}
	info, err = OperateLandlordInfo(info.TurnUser, "", landlord, OP_GRAB, 0)
	if err != nil {
		t.Error(err)
		return
	}
	log.D("Grab :%v", util.S2Json(info))
	landlord.Multiple = util.Map{"grab_score": 0}
	info, err = OperateLandlordInfo(info.TurnUser, "", landlord, OP_GRAB, 0)
	if err != nil {
		t.Error(err)
		return
	}
	log.D("Deal :%v", util.S2Json(info))
	landlord.Multiple = util.Map{"grab_score": 2}
	info, err = OperateLandlordInfo(info.TurnUser, "", landlord, OP_GRAB, 0)
	if err != nil {
		t.Error(err)
		return
	}
	log.D("Deal :%v", util.S2Json(info))

	//加倍
	landlord = &LandlordInfo{
		Id: info.Id,
	}
	landlord.Multiple = util.Map{LD_DOUBLE_USERS: info.TurnUser}
	info, err = OperateLandlordInfo(info.TurnUser, "", landlord, OP_DOUBLE, 0)
	if err != nil {
		t.Error(err)
		return
	}
	log.D("Double :%v", util.S2Json(info))
	landlord.Multiple = util.Map{}
	info, err = OperateLandlordInfo(info.TurnUser, "", landlord, OP_DOUBLE, 0)
	if err != nil {
		t.Error(err)
		return
	}
	log.D("Double :%v", util.S2Json(info))
	landlord.Multiple = util.Map{LD_DOUBLE_USERS: info.TurnUser}
	info, err = OperateLandlordInfo(info.TurnUser, "", landlord, OP_DOUBLE, 0)
	if err != nil {
		t.Error(err)
		return
	}
	log.D("Double :%v", util.S2Json(info))
	log.D("----------------------")
	//出牌
	//提示
	landlord = &LandlordInfo{
		Id: info.Id,
	}
	for i:=0;i<10;i++{
		turn := info.TurnUser

		info, err = OperateLandlordInfo(turn, "", landlord, OP_GET_NOTE, 0)
		if err != nil {
			t.Error(err)
			return
		}
		log.D("player:%v,noteCards:%v",turn, info.NotePopCards)
		if len(info.NotePopCards) >0&&info.NotePopCards!=LC_PASS {
			info, err = OperateLandlordInfo(info.TurnUser, info.NotePopCards, landlord, OP_POP_CARD, 0)
			if err != nil {
				t.Error(err)
				return
			}
			log.D("player:%v,popCards:%v", turn,  info.LastPopCards)
			if info.Status == LS_SHOW_RESULT {
				log.D("Game over: %v", util.S2Json(info))
				break
			}
		} else {
			info, err = OperateLandlordInfo(info.TurnUser, "", landlord, OP_PASS_CARD, 0)
			if err != nil {
				t.Error(err)
				return
			}
			log.D("player:%v,noteCards:%v", turn, info.NotePopCards)

		}
	}
}

func TestPopCards(t *testing.T)  {
	info:=&LandlordInfo{}
	err:=db.C(CN_LANDLORD_INFO).Find(bson.M{"players.uid":"u1"}).One(&info)
	if err!=nil {
		t.Error(err)
		return
	}
	landlord := &LandlordInfo{
		Id: info.Id,
	}
	flag:=0
	for {
		turn := info.TurnUser
		info, err = OperateLandlordInfo(turn, "", landlord, OP_GET_NOTE, 0)
		if err != nil {
			t.Error(err)
			return
		}
		log.D("player:%v,noteCards:%v", turn, info.NotePopCards)
		if len(info.NotePopCards) > 0 &&info.NotePopCards!=LC_PASS{
			info, err = OperateLandlordInfo(info.TurnUser, info.NotePopCards, landlord, OP_POP_CARD, 0)
			if err != nil {
				t.Error(err)
				return
			}
			log.D("player:%v,popCards:%v", turn, info.LastPopCards)
			if info.Status == LS_SHOW_RESULT {
				log.D("Game over: %v", util.S2Json(info))
				break
			}
		} else {
			info, err = OperateLandlordInfo(info.TurnUser, "", landlord, OP_PASS_CARD, 0)
			if err != nil {
				t.Error(err)
				return
			}
			log.D("player:%v,noteCards:%v", turn, info.NotePopCards)

		}
		if flag==102{
			break
		}
		flag++
	}
}


func TestQueue(t *testing.T) {
	//Remove()
	room, err := GetRoomRandom(uid0[2], 1)
	if err != nil {
		t.Error(err)
		return
	}
	if SHOW_LOG_DB_TEST {
		log.D("GetRoomRandom: room %v", util.S2Json(room))
	}
}

//快速加入房间
func TestGetRoomRandom(t *testing.T) {
	Remove()
	//首个用户，创建新房间
	landlordInfo, err := GetRoomRandom(uid0[0], 1)
	if err != nil {
		t.Error(err)
		return
	}
	log.D("GetRoomRandom create new room %v", util.S2Json(landlordInfo))
	//加入已创建的房间
	landlordInfo, err = GetRoomRandom(uid0[1], 1)
	if err != nil {
		t.Error(err)
		return
	}
	log.D("GetRoomRandom  room %v", util.S2Json(landlordInfo))
	//加入已创建的房间
	landlordInfo, err = GetRoomRandom(uid0[2], 1)
	if err != nil {
		t.Error(err)
		return
	}
	log.D("GetRoomRandom join in room %v", util.S2Json(landlordInfo))
}

/**
获取斗地主大厅的设置测试
*/
func TestFindLandlordsHall(t *testing.T) {
	db.C(CN_LANDLORD_SET).RemoveAll(nil)
	set := &LandlordSet{
		Id:          bson.NewObjectId().Hex(),
		SoundEffect: 1,
		Music:       1,
		Categories: []util.Map{
			util.Map{
				"cid":1,
				"title":     LC_PRIMER_TITLE,
				"level":     LC_PRIMER_LEVEL,
				"end_score": LC_PRIMER_SCORE_END,
				"min_gold":  LC_PRIMER_GOLD_MIN,
			},
			util.Map{
				"cid":2,
				"title":     LC_SUPERIOR_TITLE,
				"level":     LC_SUPERIOR_LEVEL,
				"end_score": LC_SUPERIOR_SCORE_END,
				"min_gold":  LC_SUPERIOR_GOLD_MIN,
			},
			util.Map{
				"cid":3,
				"title":     LC_MASTER_TITLE,
				"level":     LC_MASTER_LEVEL,
				"end_score": LC_MASTER_SCORE_END,
				"min_gold":  LC_MASTER_GOLD_MIN,
			},
			util.Map{
				"cid":4,
				"title":     LC_TOP_TITLE,
				"level":     LC_TOP_LEVEL,
				"end_score": LC_TOP_SCORE_END,
				"min_gold":  LC_TOP_GOLD_MIN,
			},
		},
		QueueTime:    LT_QUEUE,
		MingTime:     LT_MING,
		GrabTime:     LT_GRAB,
		DoubleTime:   LT_DOUBLE,
		PopCardTime:  LT_POP_CARD,
		OverGameTime: LT_GAME_OVER_SHOW,
	}
	err := db.C(CN_LANDLORD_SET).Insert(set)
	landlordSet, err := FindLandlordsHall(nil)
	if err != nil {
		t.Error(err)
		return
	}
	if SHOW_LOG_DB_TEST {
		log.D(TAG_LANDLORD_DB+",FindLandlordsHall,Find success landlord_set: %v", util.S2Json(landlordSet))
	}
}
