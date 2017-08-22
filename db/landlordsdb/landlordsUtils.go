package landlordsdb

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	BASE_CARD_NUM     = 17
	LANDLORD_CARD_NUM = 3
	PLAYER_NUM        = 3
	COLOR_NUM         = 4
	CARD_SUM_NUM      = 54
	ORDINARY_CARD_SUM = 13
	BLACK_JOKER_INDEX = 52
	TWO_INDEX         = 48
	//权重
	WEIGHT_NONE        = 0
	WEIGHT_THREE       = 1
	WEIGHT_FOUR        = 2
	WEIGHT_FIVE        = 3
	WEIGHT_SIX         = 4
	WEIGHT_SEVEN       = 5
	WEIGHT_EIGHT       = 6
	WEIGHT_NINE        = 7
	WEIGHT_TEN         = 8
	WEIGHT_JACK        = 9
	WEIGHT_QUEEN       = 10
	WEIGHT_KING        = 11
	WEIGHT_ACE         = 12
	WEIGHT_TWO         = 13
	WEIGHT_BLACK_JOKER = 15
	WEIGHT_RED_JOKER   = 16

	//出牌类型
	TYPE_NONE                   = 0
	TYPE_ONE                    = 1
	TYPE_PAIR                   = 2
	TYPE_THREE                  = 3
	TYPE_THREE_WITH_ONE         = 4
	TYPE_THREE_WITH_PAIR        = 5
	TYPE_FOUR_WITH_TWO          = 6
	TYPE_FOUR_WITH_TWO_PAIR     = 7
	TYPE_SINGLE_STRAIGHT        = 8
	TYPE_DOUBLE_STRAIGHT        = 9
	TYPE_AIRPLANE               = 10
	TYPE_AIRPLANE_WITH_SINGLE   = 11
	TYPE_AIRPLANE_WITH_PAIR     = 12
	TYPE_THREE_SHUN             = 13
	TYPE_THREE_SHUN_WITH_SINGLE = 14
	TYPE_THREE_SHUN_WITH_PAIR   = 15
	TYPE_BOMB_ORDINARY          = 16 //炸弹
	TYPE_ROCKET                 = 17 //双王火箭

	DOWN_PLAYER  = 1 //当前用户
	RIGHT_PLAYER = 2 //当前用户的下一个
	LEFT_PLAYER  = 3 //right用户的下一个

)

var (
	AllTypes []string = []string{"TYPE_NONE", "TYPE_ONE", "TYPE_PAIR", "TYPE_THREE", "TYPE_THREE_WITH_ONE", "TYPE_THREE_WITH_PAIR",
		"TYPE_FOUR_WITH_TWO", "TYPE_FOUR_WITH_TWO_PAIR", "TYPE_SINGLE_STRAIGHT", "TYPE_DOUBLE_STRAIGHT", "TYPE_AIRPLANE",
		"TYPE_AIRPLANE_WITH_SINGLE", "TYPE_AIRPLANE_WITH_PAIR", "TYPE_THREE_SHUN", "TYPE_THREE_SHUN_WITH_SINGLE", "TYPE_THREE_SHUN_WITH_PAIR",
		"TYPE_BOMB_ORDINARY", "TYPE_ROCKET"}

	AllCards       []int             = make([]int, CARD_SUM_NUM)
	PlayersCards   map[string]string = make(map[string]string)
	ACards         string            = ""
	BCards         string            = ""
	CCards         string            = ""
	LandlordCards  string            = ""
	LandlordPlayer string            = ""
	WeightPop      int               = 0
	WeightLast     int               = 0
	TypeLast       int               = 0
	CardsPopLast   string            = ""
	TurnPlayer     string            = "" //轮流
)

/**
初始化参数
*/
func InitData() {
	AllCards = make([]int, CARD_SUM_NUM)
	PlayersCards = make(map[string]string)
	ACards = ""
	CCards = ""
	LandlordCards = ""
	LandlordPlayer = ""
	WeightPop = 0
	WeightLast = 0
	TypeLast = 0
	CardsPopLast = ""
	TurnPlayer = ""
}

/**
初始化牌
*/
func InitCards() {
	//3-K,A,2 1-4,5-8...,53,54
	for i := 0; i < CARD_SUM_NUM; i++ {
		AllCards[i] = i
	}
	rand.Seed(time.Now().UnixNano())
	temp := 0
	for i := 0; i < CARD_SUM_NUM; i++ {
		temp = rand.Intn(i + 1)
		if temp != i {
			swap := AllCards[i]
			AllCards[i] = AllCards[temp]
			AllCards[temp] = swap
		}

	}

}

//当最后一个用户进入房间后发起随机分牌
func RandomDistributeCards(players []util.Map) ([]util.Map, string) {
	InitData()
	InitCards()
	var a_cards []int = make([]int, BASE_CARD_NUM)
	var b_cards []int = make([]int, BASE_CARD_NUM)
	var c_cards []int = make([]int, BASE_CARD_NUM)
	//发牌
	for i := 0; i < BASE_CARD_NUM; i++ {
		a_cards[i] = AllCards[i*PLAYER_NUM]
		b_cards[i] = AllCards[i*PLAYER_NUM+1]
		c_cards[i] = AllCards[i*PLAYER_NUM+2]

	}
	//排序
	sort.Ints(a_cards)
	sort.Ints(b_cards)
	sort.Ints(c_cards)

	for i := 0; i < BASE_CARD_NUM; i++ {
		if i == BASE_CARD_NUM-1 {
			ACards += fmt.Sprintf("%v", a_cards[i])
			BCards += fmt.Sprintf("%v", b_cards[i])
			CCards += fmt.Sprintf("%v", c_cards[i])
		} else {
			ACards += fmt.Sprintf("%v,", a_cards[i])
			BCards += fmt.Sprintf("%v,", b_cards[i])
			CCards += fmt.Sprintf("%v,", c_cards[i])
		}
	}

	players_u:=[]util.Map{util.Map{"uid":players[0].StrVal("uid"),"cards":ACards},
		util.Map{"uid":players[1].StrVal("uid"),"cards":BCards},
		util.Map{"uid":players[2].StrVal("uid"),"cards":CCards}}

	for i := BASE_CARD_NUM * PLAYER_NUM; i < CARD_SUM_NUM; i++ {
		if i == CARD_SUM_NUM-1 {
			LandlordCards += fmt.Sprintf("%v", AllCards[i])
		} else {
			LandlordCards += fmt.Sprintf("%v,", AllCards[i])
		}
	}

	return players_u, LandlordCards
}

/**
确定先手玩家
*/
func ConfirmFirstPlayer() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(3)
}

/**
抢地主规则
*/
func GrabLandlordRule(uid string, landLordInfo, info *LandlordInfo) *LandlordInfo {
	players:=[]util.Map{}
	for i := 0; i < len(info.Players); i++ {
		players = append(players,util.Map{"uid":info.Players[i].StrVal("uid")})
		if  uid == info.Players[i].StrVal("uid"){
			players[i][LD_GRAB_SCORE] = landLordInfo.Multiple.IntVal(LD_GRAB_SCORE)
		}
	}
	//任意一个玩家抢地主分数为最高分或者最后一个抢地主且分数高于之前的则确定抢地主玩家
	if landLordInfo.Multiple.IntVal(LD_GRAB_SCORE) == 3 ||
		(landLordInfo.Multiple.IntVal(LD_GRAB_SCORE) > info.Multiple.IntVal(LD_GRAB_SCORE) && info.OperateNum == PLAYER_NUM-1) {
		for i := 0; i < len(info.Players); i++ {
			if info.Players[i].StrVal("uid") == uid {
				players[i]["cards"] = AddLordsCardsForLand(info.Players[i].StrVal("cards"),info.LandlordCards)
				//landLordInfo.Players[i]["cards"] = AddLordsCardsForLand(info.Players[i].StrVal("cards"),info.LandlordCards)
			}
		}
		landLordInfo.LandlordUser = uid
		landLordInfo.TurnUser = uid
		landLordInfo.Status = LS_DOUBLING
		landLordInfo.DoubleTime = util.Now()
		landLordInfo.OperateNum = 0
	} else if landLordInfo.Multiple.IntVal(LD_GRAB_SCORE) < 3 && info.OperateNum < PLAYER_NUM-1 {
		//非最后一个玩家抢地主且抢地主分数小于3且比之前的分数要高
		if landLordInfo.Multiple.IntVal(LD_GRAB_SCORE) > info.Multiple.IntVal(LD_GRAB_SCORE)  {
			for i := 0; i < len(info.Players); i++ {
				if info.Players[i].StrVal("uid") == uid {
					landLordInfo.TurnUser = info.Players[(i+1)%3].StrVal("uid")
				}
			}
		}else {
			landLordInfo.Multiple[LD_GRAB_SCORE] = 0
		}

		landLordInfo.OperateNum = info.OperateNum + 1
		landLordInfo.GrabTime = util.Now() //重新计算时间
	}

	//最后一个不抢
	if landLordInfo.Multiple.IntVal(LD_GRAB_SCORE) == 0&& info.OperateNum == PLAYER_NUM-1 {
		if info.Multiple.IntVal(LD_GRAB_SCORE) == 0 {//重新发牌
			log.D("都不抢地主,重新发牌")
			return nil
		}
		//从之前的玩家选择最高分的
		for i := 0; i < 2; i++ {
			if info.Players[i].IntVal(LD_GRAB_SCORE) == info.Multiple.IntVal(LD_GRAB_SCORE) {
				players[i]["cards"] = AddLordsCardsForLand(info.Players[i].StrVal("cards"),info.LandlordCards)
				landLordInfo.LandlordUser = players[i].StrVal("uid")
				landLordInfo.TurnUser = players[i].StrVal("uid")
			}
		}
		landLordInfo.Status = LS_DOUBLING
		landLordInfo.DoubleTime = util.Now()
		landLordInfo.OperateNum = 0

	}
	landLordInfo.Players = players
	return landLordInfo
}

/**
加倍规则
*/
func DoubleLandlordRule(uid string, landLordInfo, info *LandlordInfo) *LandlordInfo {
	if info.OperateNum == PLAYER_NUM-1 {
		landLordInfo.OperateNum = 0
		landLordInfo.Status = LS_FIGHTING
		landLordInfo.PopCardTime = util.Now()
		landLordInfo.TurnUser = info.LandlordUser
	} else if info.OperateNum < PLAYER_NUM-1 {
		for i := 0; i < len(info.Players); i++ {
			if info.Players[i].StrVal("uid") == uid {
				landLordInfo.TurnUser = info.Players[(i+1)%3].StrVal("uid")
				break
			}
		}
		landLordInfo.OperateNum = info.OperateNum + 1
		landLordInfo.DoubleTime = util.Now()
	}
	if landLordInfo.Multiple.StrVal(LD_DOUBLE_USERS)==uid {//加倍
		landLordInfo.Multiple[LD_DOUBLE_USERS] = uid
	}
	return landLordInfo
}

/**
出牌规则
*/
func FightLandlordRule(uid, popCards string, landLordInfo, info *LandlordInfo) *LandlordInfo {
	pop_type, _ := JudgeCardsTypeAndWeight(popCards)
	multiple := landLordInfo.Multiple
	if len(multiple)<1 {
		multiple = util.Map{}
	}
	if pop_type == TYPE_BOMB_ORDINARY || pop_type == TYPE_ROCKET { //炸弹
		multiple[LD_BOMBS] = info.Multiple.IntVal(LD_BOMBS) + 1
	}
	flag_over := false
	players:=[]util.Map{util.Map{"uid":info.Players[0].StrVal("uid")},util.Map{"uid":info.Players[1].StrVal("uid")},util.Map{"uid":info.Players[2].StrVal("uid")}}
	for i := 0; i < len(info.Players); i++ {
		if info.Players[i].StrVal("uid") == uid {
			players[i]["pop_cards"] = popCards
			cards_remove := RemovePopCards(info.Players[i].StrVal("cards"),popCards)
			players[i]["cards"] = cards_remove
			players[i]["pop_times"] =info.Players[i].IntVal("pop_times")+1
			if len(cards_remove) == 0 {
				flag_over = true//出完牌没
			}
		}
	}

	//是否结束
	if flag_over {
		log.D("OverGame")
		//结束
		//record := &LandlordRecord{
		//	Id:               bson.NewObjectId().Hex(),
		//	Users:            []string{info.Players[0].StrVal("uid"),info.Players[1].StrVal("uid"),info.Players[2].StrVal("uid")},
		//	LandlordUser:     info.LandlordUser,
		//	Category:         info.Category,
		//	Time:             util.Now(),
		//	BombMultiple:     landLordInfo.Multiple.IntVal(LD_BOMBS),
		//	DoubleMultiple:   int64(len(info.Multiple.StrVal(LD_DOUBLE_USERS))),
		//	LandlordMultiple: info.Multiple.IntVal(LD_GRAB_SCORE),
		//}
		////统计战果
		//landLordInfo.OperateNum = 0
		//landLordInfo.Status = LS_SHOW_RESULT
		//landLordInfo.OverGameTime = util.Now()
		////判断春天反春天
		//if uid == info.LandlordUser && IsMeetSpring(uid, info.Players) {
		//	record.WinRole = LW_LANDLORD
		//	record.SpringMultiple = 1
		//	//multiple[LD_SPRING] = 1
		//}
		//if uid != info.LandlordUser && IsMeetAntiSpring(info.LandlordUser, info.Players) {
		//	record.WinRole = LW_FARMER
		//	record.AntiSpringMultiple = 1
		//	//multiple[LD_ANTI_SPRING] = 1
		//}
		//record.SumMultiple = record.LandlordMultiple
		//sum := int(record.DoubleMultiple + record.BombMultiple + record.SpringMultiple + record.AntiSpringMultiple)
		//for i := 0; i < sum; i++ {
		//	record.SumMultiple *= 2
		//}
		//if info.Multiple.IntVal(LD_MING) > 0 {
		//	record.SumMultiple *= info.Multiple.IntVal(LD_MING)
		//}
		//err := db.C(CN_LANDLORD_RECORD).Insert(&record)
		//if err != nil {
		//	log.E("insert landlord record err(%v)", err)
		//	return nil
		//}
		//players_u:=info.Players
		//for i := 0; i < len(players_u); i++ {
		//	players_u[i]["pop_times"] = 0
		//	players_u[i]["cards"] = ""
		//	players_u[i]["pop_cards"] = LC_PASS
		//	players_u[i]["ming"] = ""
		//	if players_u[i].IntVal("status") == LUS_COLLOCATION {
		//		players_u[i]["status"] = LUS_ONLINE
		//	}
		//}
		////更新房间
		//err=db.C(CN_LANDLORD_INFO).Update(bson.M{"_id":info.Id},bson.M{
		//	"$set":bson.M{
		//		"players":players_u,
		//		"turn_user":"",
		//		"operate_num" : 0,
		//		"landlord_cards" : "",
		//		"last_pop_cards" : LC_PASS,
		//		"note_pop_cards" : LC_PASS,
		//		"status" : LS_SHOW_RESULT,
		//		"landlord_user" : "",
		//		"multiple" :util.Map{},
		//		"last" :util.Now(),
		//		"queue_time" : 0,
		//		"grab_time" : 0,
		//		"double_time" : 0,
		//		"pop_card_time" : 0,
		//		"over_game_time" :util.Now(),
		//		"ming_time" : 0,
		//		"size" : 0,
		//	}})
		//if err != nil {
		//	log.E("update landlordInfo err(%v)", err)
		//	return nil
		//}
		////离线移除
		//for i := 0; i < len(info.Players); i++ {
		//	if info.Players[i].IntVal("status") == LUS_OFFLINE {
		//		u:=info.Players[i].StrVal("uid")
		//		err=db.C(CN_LANDLORD_INFO).Update(bson.M{"_id":info.Id,"players.uid":u,"players.status":LUS_OFFLINE},bson.M{
		//			"$pull":bson.M{
		//				"players.uid":u,
		//			}})
		//		if err != nil {
		//			log.E("update landlordInfo err(%v)", err)
		//			return nil
		//		}
		//	}
		//}

		return nil
	} else {
		landLordInfo.OperateNum = (info.OperateNum + 1) % 3
		landLordInfo.PopCardTime = util.Now()
		landLordInfo.LastPopCards = popCards
		for i := 0; i < len(info.Players); i++ {
			if info.Players[i].StrVal("uid") == uid {
				landLordInfo.TurnUser = info.Players[(i+1)%3].StrVal("uid")
				break
			}
		}
	}
	landLordInfo.NotePopCards = LC_PASS
	landLordInfo.Players = players
	landLordInfo.Multiple = multiple
	return landLordInfo
}

/**
不出牌规则
*/
func PassCardRule(uid string, landlordInfo, info *LandlordInfo) *LandlordInfo {
	players:=[]util.Map{util.Map{"uid":info.Players[0].StrVal("uid")},util.Map{"uid":info.Players[1].StrVal("uid")},util.Map{"uid":info.Players[2].StrVal("uid")}}
	for i := 0; i < len(info.Players); i++ {
		if uid == info.Players[i].StrVal("uid") {
			players[i]["pop_cards"] = LC_PASS
			if info.Players[(i+2)%3].StrVal("pop_cards") == LC_PASS||info.Players[(i+2)%3].StrVal("pop_cards") == "" {//如果上家跟自己都不出，那么lastPopCards为空
				landlordInfo.LastPopCards = LC_PASS
			}
			landlordInfo.TurnUser = info.Players[(i+1)%3].StrVal("uid")
			break
		}
	}
	landlordInfo.Players = players
	landlordInfo.NotePopCards = LC_PASS
	landlordInfo.OperateNum = (info.OperateNum + 1) % 3
	landlordInfo.PopCardTime = util.Now()
	return landlordInfo
}

/**
提示出牌规则
*/
//func NotePopCardRule(landlordInfo, info *LandlordInfo) *LandlordInfo {
//	//1.
//
//}

/**
获取提示出牌
*/
func GetNoteCardsRule(uid string, info *LandlordInfo) string {
	//根据上家出牌以及中途出现的不出牌开启新回合
	lastPopCards := info.LastPopCards
	if lastPopCards == LC_PASS{
		lastPopCards = ""
	}
	cards_player := ""
	for i := 0; i < PLAYER_NUM; i++ {
		if info.Players[i].StrVal("uid") == uid {
			if info.Players[(i+1)%3].StrVal("pop_cards") == LC_PASS && info.Players[(i+2)%3].StrVal("pop_cards") == LC_PASS{
				lastPopCards = ""
			}
			cards_player = info.Players[i].StrVal("cards")
		}
	}
	note_cards := ""
	note_cards_last := info.NotePopCards
	if note_cards_last == LC_PASS {
		note_cards_last = ""
	}
	if len(note_cards_last) < 1 { //第一次提示
		if len(lastPopCards) < 1 { //首发出牌
			return FirstPopAndFirstNote(cards_player)//首发第一次出牌
		} else { //非首发出牌
			return NotFirstPopAndFirstNote(cards_player,lastPopCards)//非首发第一次出牌
		}
	}else {
		if len(lastPopCards)<1 {
			//首发非第一次提示
			note_cards = FirstPopAndNotFirstNote(cards_player,note_cards_last)
			if len(note_cards)<1 {
				return FirstPopAndFirstNote(cards_player)
			}
		}else {
			//非首发非首次提示
			return NotFirstPopAndNotFirstNote(cards_player,lastPopCards,note_cards_last)
		}
	}
	return LC_PASS
}

/**
首发第一次提示
 */
func FirstPopAndFirstNote(cards_player string) string {
	cards_player_int := StringToIntCards(cards_player)
	cards_type, _ := JudgeCardsTypeAndWeight(cards_player)
	note_cards:=""
	//判断其是否符合出牌规则
	if cards_type > TYPE_NONE {
		return cards_player
	}
	note_cards = getNotePopCardsForTypeONEPure(cards_player_int, 0, 0)
	if len(note_cards) > 0 {
		return note_cards
	}
	note_cards = getNotePopCardsForTypePAIRPure(cards_player_int, 0, 0)
	if len(note_cards) > 0 {
		return note_cards
	}
	note_cards = getNotePopCardsForTypeTHREE(cards_player_int, 0, 0)
	if len(note_cards) > 0 {
		return note_cards
	}
	note_cards = getNotePopCardsForTypeBOMB(cards_player_int, 0, 0)
	if len(note_cards) > 0 {
		return note_cards
	}
	note_cards = getNotePopCardsForTypeRocket(cards_player_int)
	if len(note_cards) > 0 {
		return note_cards
	}
	return ""
}
/**
首发非第一次提示
 */
func FirstPopAndNotFirstNote(cards_player,lastNoteCards string) string {
	//已经提示过
	note_cards := ""
	last_note_type, last_note_weight := JudgeCardsTypeAndWeight(lastNoteCards)
	cards_player_int := StringToIntCards(cards_player)
	last_note_size := len(StringToIntCards(lastNoteCards))
	switch last_note_type {
	case TYPE_ONE:
		note_cards = getNotePopCardsForTypeONE(cards_player_int, last_note_weight, 0)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_PAIR:
		note_cards = getNotePopCardsForTypePAIR(cards_player_int, last_note_weight, 0)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_THREE:
		note_cards = getNotePopCardsForTypeTHREE(cards_player_int, last_note_weight, 0)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_THREE_WITH_ONE:
		note_cards = getNotePopCardsForTypeThreeWithOne(cards_player_int, last_note_weight, 0)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_BOMB_ORDINARY:
		note_cards = getNotePopCardsForTypeBOMB(cards_player_int, last_note_weight, 0)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_THREE_WITH_PAIR:
		note_cards = getNotePopCardsForTypeThreeWithPair(cards_player_int, last_note_weight, 0)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_FOUR_WITH_TWO:
		note_cards = getNotePopCardsForTypeFourWithSingle(cards_player_int, last_note_weight, 0)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_FOUR_WITH_TWO_PAIR:
		note_cards = getNotePopCardsForTypeFourWithTwoPair(cards_player_int, last_note_weight, 0)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_SINGLE_STRAIGHT:
		note_cards = getNotePopCardsForTypeSingleStraight(cards_player_int, last_note_weight, 0, 0, last_note_size)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_DOUBLE_STRAIGHT:
		last_note_size = last_note_size / 2
		note_cards = getNotePopCardsForTypeDoubleStraight(cards_player_int, last_note_weight, 0, 0, last_note_size)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_THREE_SHUN:
		last_note_size = last_note_size / 3
		note_cards = getNotePopCardsForTypeThreeShun(cards_player_int, last_note_weight, 0, 0, last_note_size)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_THREE_SHUN_WITH_SINGLE:
		last_note_size = last_note_size / 4
		note_cards = getNotePopCardsForTypeThreeShunWithSingle(cards_player_int, last_note_weight, 0, 0, last_note_size)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_THREE_SHUN_WITH_PAIR:
		last_note_size = last_note_size / 5
		note_cards = getNotePopCardsForTypeThreeShunWithPair(cards_player_int, last_note_weight, 0, 0, last_note_size)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_AIRPLANE:
		note_cards = getNotePopCardsForTypeAirplane(cards_player_int, last_note_weight, 0)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_AIRPLANE_WITH_SINGLE:
		note_cards = getNotePopCardsForTypeAirplaneWithSingle(cards_player_int, last_note_weight, 0)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_AIRPLANE_WITH_PAIR:
		note_cards = getNotePopCardsForTypeAirplaneWithPair(cards_player_int, last_note_weight, 0)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	}
	note_cards = getNotePopCardsForTypeBOMB(cards_player_int, last_note_weight, 0)
	if len(note_cards) > 0 {
		return note_cards
	}
	note_cards = getNotePopCardsForTypeRocket(cards_player_int)
	if len(note_cards) > 0 {
		return note_cards
	}
	return ""
}
/**
非首发第一次提示
 */
func NotFirstPopAndFirstNote(cards_player,lastPopCards string) string {
	note_cards:=""
	last_pop_type, last_pop_weight := JudgeCardsTypeAndWeight(lastPopCards)
	cards_player_int := StringToIntCards(cards_player)
	last_pop_size := len(StringToIntCards(lastPopCards))
	switch last_pop_type {
	case TYPE_ONE:
		note_cards = getNotePopCardsForTypeONE(cards_player_int, 0, last_pop_weight)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_PAIR:
		note_cards = getNotePopCardsForTypePAIR(cards_player_int, 0, last_pop_weight)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_THREE:
		note_cards = getNotePopCardsForTypeTHREE(cards_player_int, 0, last_pop_weight)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_THREE_WITH_ONE:
		note_cards = getNotePopCardsForTypeThreeWithOne(cards_player_int, 0, last_pop_weight)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_BOMB_ORDINARY:
		note_cards = getNotePopCardsForTypeBOMB(cards_player_int, 0, last_pop_weight)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_THREE_WITH_PAIR:
		note_cards = getNotePopCardsForTypeThreeWithPair(cards_player_int, 0, last_pop_weight)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_FOUR_WITH_TWO:
		note_cards = getNotePopCardsForTypeFourWithSingle(cards_player_int, 0, last_pop_weight)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_FOUR_WITH_TWO_PAIR:
		note_cards = getNotePopCardsForTypeFourWithTwoPair(cards_player_int, 0, last_pop_weight)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_SINGLE_STRAIGHT:
		note_cards = getNotePopCardsForTypeSingleStraight(cards_player_int, 0, 0, last_pop_weight, last_pop_size)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_DOUBLE_STRAIGHT:
		last_pop_size = last_pop_size / 2
		note_cards = getNotePopCardsForTypeDoubleStraight(cards_player_int, 0, 0, last_pop_weight, last_pop_size)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_THREE_SHUN:
		last_pop_size = last_pop_size / 3
		note_cards = getNotePopCardsForTypeThreeShun(cards_player_int, 0, 0, last_pop_weight, last_pop_size)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_THREE_SHUN_WITH_SINGLE:
		last_pop_size = last_pop_size / 4
		note_cards = getNotePopCardsForTypeThreeShunWithSingle(cards_player_int, 0, 0, last_pop_weight, last_pop_size)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_THREE_SHUN_WITH_PAIR:
		last_pop_size = last_pop_size / 5
		note_cards = getNotePopCardsForTypeThreeShunWithPair(cards_player_int, 0, 0, last_pop_weight, last_pop_size)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_AIRPLANE:
		note_cards = getNotePopCardsForTypeAirplane(cards_player_int, 0, last_pop_weight)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_AIRPLANE_WITH_SINGLE:
		note_cards = getNotePopCardsForTypeAirplaneWithSingle(cards_player_int, 0, last_pop_weight)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	case TYPE_AIRPLANE_WITH_PAIR:
		note_cards = getNotePopCardsForTypeAirplaneWithPair(cards_player_int, 0, last_pop_weight)
		if len(note_cards) > 0 {
			return note_cards
		}
		break
	}
	//打不过就炸
	note_cards = getNotePopCardsForTypeBOMB(cards_player_int, 0, last_pop_weight)
	if len(note_cards) > 0 {
		return note_cards
	}
	note_cards = getNotePopCardsForTypeRocket(cards_player_int)
	if len(note_cards) > 0 {
		return note_cards
	}
	return ""
}
/**
非首发非第一次提示
 */
func NotFirstPopAndNotFirstNote(cards_player,lastPopCards,lastNoteCards string) string  {
	note:=FirstPopAndNotFirstNote(cards_player,lastNoteCards)
	if len(note)>0 {
		return note
	}
	note=NotFirstPopAndFirstNote(cards_player,lastPopCards)
	return note
}

/**
春天
*/
func IsMeetAntiSpring(landlord string, players []util.Map) bool {
	for i := 0; i < len(players); i++ {
		if landlord == players[i].StrVal("uid") && players[i].IntVal("pop_times") == 1 {
			return true
		}
	}
	return false
}

/**
反春天
*/
func IsMeetSpring(uid string, players []util.Map) bool {
	for i := 0; i < len(players); i++ {
		if uid == players[i].StrVal("uid") && players[(i+1)%3].IntVal("pop_times") == 0 &&
			players[(i+2)%3].IntVal("pop_times")== 0 {
			return true
		}
	}
	return false
}

/**
添加地主牌
*/
func AddLordsCardsForLand(player_cards, landlord_cards string) string {
	cards := fmt.Sprintf("%v,%v", player_cards, landlord_cards)
	cards_int := StringToIntCards(cards)
	if cards_int == nil {
		return ""
	}
	sort.Ints(cards_int)
	return IntToStringCards(cards_int)
}

/**
出牌
*/
func RemovePopCards(player_cards, pop_cards string) string {
	if pop_cards==LC_PASS || pop_cards==""{
		return player_cards
	}
	player_cards_int := StringToIntCards(player_cards)
	pop_cards_int := StringToIntCards(pop_cards)
	new_cards_int := make([]int,len(player_cards_int)-len(pop_cards_int))
	flag := 0
	log.D("pop:%v,cards:%v",pop_cards,player_cards)
	isPop:=false
	for i := 0; i < len(player_cards_int); i++ {
		isPop = false
		for j := 0; j < len(pop_cards_int); j++ {
			if player_cards_int[i] == pop_cards_int[j] {
				isPop = true
				break
			}
		}

		if  !isPop{
			new_cards_int[flag] = player_cards_int[i]
			flag++
		}
	}
	if new_cards_int == nil {
		return ""
	}
	sort.Ints(new_cards_int)
	return IntToStringCards(new_cards_int)

}

/**
判断游戏是否结束
*/
func IsGameOver() {

}

/**
是否满足出牌条件
同类型比较权重
不同类型
1.王炸 唯一性谁有王炸则谁大
2.普通炸 谁出普通炸则谁大
*/
func IsMeetPopLogic(popCards, lastPopCards string) bool {
	if len(popCards) < 1 {
		return false
	}

	type_pop, weight_pop := JudgeCardsTypeAndWeight(popCards)
	type_last, weight_last := JudgeCardsTypeAndWeight(lastPopCards)
	//王炸
	if type_last == TYPE_ROCKET {
		return false
	}
	if type_pop == TYPE_ROCKET {
		return true
	}
	//同类型比较
	if type_pop == type_last {
		if weight_pop < weight_last {
			return false
		}
	} else {
		//炸弹
		if type_last == TYPE_BOMB_ORDINARY {
			return false
		}
		if type_pop == TYPE_BOMB_ORDINARY {
			return true
		}
	}
	return true
}

/**
判断出牌类型和权重
 */
func JudgeCardsTypeAndWeight(cards_str string) (int, int) {
	if len(cards_str) < 1 {
		return TYPE_NONE, 0
	}
	cards := StringToIntCards(cards_str)
	sort.Ints(cards)
	type_cards := TYPE_NONE
	WeightPop = 0
	weight_cards := WEIGHT_NONE
	size := len(cards)
	switch size {
	case 1:
		type_cards = TYPE_ONE
		WeightPop = cards[0]/COLOR_NUM + 1
		break
	case 2:
		if hasSameTwo(cards) {
			WeightPop = cards[0]/COLOR_NUM + 1
			type_cards = TYPE_PAIR
		} else if hasRocket(cards) {
			WeightPop = WEIGHT_RED_JOKER
			type_cards = TYPE_ROCKET
		}
		break
	case 3:
		if hasSameThree(cards) {
			WeightPop = cards[0]/COLOR_NUM + 1
			type_cards = TYPE_THREE
		}
		break
	case 4:
		if hasThreeWithOne(cards) {
			WeightPop = cards[1]/COLOR_NUM + 1
			type_cards = TYPE_THREE_WITH_ONE
		} else if hasBombOrdinary(cards) {
			WeightPop = cards[0]/COLOR_NUM + 1
			type_cards = TYPE_BOMB_ORDINARY
		}
		break
	case 5:
		if hasThreeWithPair(cards) {
			WeightPop = cards[2]/COLOR_NUM + 1
			type_cards = TYPE_THREE_WITH_PAIR
		}
		break
	case 6:
		if hasAirplane(cards) {
			WeightPop = cards[5]/COLOR_NUM + 1
			type_cards = TYPE_AIRPLANE
		} else if hasFourWithTwoSingle(cards) {
			WeightPop = cards[2]/COLOR_NUM + 1
			type_cards = TYPE_FOUR_WITH_TWO
		}
		break
	case 8:
		if hasAirplaneWithSingle(cards) {
			type_cards = TYPE_AIRPLANE_WITH_SINGLE
		} else if hasFourWithTwoPair(cards) {
			for i := 0; i < 5; i += 2 {
				if cards[i]/COLOR_NUM == cards[i+3]/COLOR_NUM {
					WeightPop = cards[i]/COLOR_NUM + 1
				}
			}
			type_cards = TYPE_FOUR_WITH_TWO_PAIR
		}
		break
	case 10:
		if hasAirplaneWithPair(cards) {
			type_cards = TYPE_AIRPLANE_WITH_PAIR
		}
		break
	default:
		break
	}
	if type_cards == TYPE_NONE {
		if hasSingleShun(cards) {
			WeightPop = cards[len(cards)-1]/COLOR_NUM + 1
			type_cards = TYPE_SINGLE_STRAIGHT
		} else if hasPairShun(cards) {
			WeightPop = cards[len(cards)-1]/COLOR_NUM + 1
			type_cards = TYPE_DOUBLE_STRAIGHT
		} else if hasThreeShun(cards) {
			WeightPop = cards[len(cards)-1]/COLOR_NUM + 1
			type_cards = TYPE_THREE_SHUN
		} else if hasThreeShunWithSingle(cards) {
			type_cards = TYPE_THREE_SHUN_WITH_SINGLE

		} else if hasThreeShunWithPair(cards) {
			type_cards = TYPE_THREE_SHUN_WITH_PAIR

		}
	}
	if type_cards != TYPE_NONE {
		weight_cards = WeightPop
	}
	return type_cards, weight_cards
}

/**
对
*/
func hasSameTwo(cards []int) bool {
	return cards[0]/COLOR_NUM == cards[1]/COLOR_NUM && cards[0] < BLACK_JOKER_INDEX && cards[1] < BLACK_JOKER_INDEX
}

/**
王炸
*/
func hasRocket(cards []int) bool {
	return cards[0] != cards[1] && cards[0] >= BLACK_JOKER_INDEX && cards[1] >= BLACK_JOKER_INDEX
}

/**
三条
*/
func hasSameThree(cards []int) bool {
	return cards[0]/COLOR_NUM == cards[1]/COLOR_NUM && cards[0]/COLOR_NUM == cards[2]/COLOR_NUM &&
		cards[0] < BLACK_JOKER_INDEX && cards[1] < BLACK_JOKER_INDEX && cards[2] < BLACK_JOKER_INDEX
}

/**
三带一
*/
func hasThreeWithOne(cards []int) bool {
	return cards[1]/COLOR_NUM == cards[2]/COLOR_NUM && ((cards[0]/COLOR_NUM == cards[1]/COLOR_NUM && cards[1]/COLOR_NUM != cards[3]/COLOR_NUM) ||
		(cards[1]/COLOR_NUM == cards[3]/COLOR_NUM && cards[0]/COLOR_NUM != cards[1]/COLOR_NUM))
}

/**
四条
*/
func hasBombOrdinary(cards []int) bool {
	return cards[0]/COLOR_NUM == cards[1]/COLOR_NUM && cards[2]/COLOR_NUM == cards[3]/COLOR_NUM &&
		cards[0]/COLOR_NUM == cards[2]/COLOR_NUM && cards[3] < BLACK_JOKER_INDEX
}

/**
三带一对
*/
func hasThreeWithPair(cards []int) bool {
	return cards[0]/COLOR_NUM == cards[1]/COLOR_NUM && cards[3]/COLOR_NUM == cards[4]/COLOR_NUM &&
		((cards[0]/COLOR_NUM == cards[2]/COLOR_NUM && cards[3] < BLACK_JOKER_INDEX) ||
			(cards[3]/COLOR_NUM == cards[2]/COLOR_NUM && cards[1] < BLACK_JOKER_INDEX))

}

/**
四带二
*/
func hasFourWithTwoSingle(cards []int) bool {
	return (cards[3]/COLOR_NUM == cards[0]/COLOR_NUM && cards[1]/COLOR_NUM == cards[2]/COLOR_NUM && cards[1]/COLOR_NUM == cards[3]/COLOR_NUM && cards[4] < BLACK_JOKER_INDEX) ||
		(cards[3]/COLOR_NUM == cards[1]/COLOR_NUM && cards[2]/COLOR_NUM == cards[4]/COLOR_NUM && cards[2]/COLOR_NUM == cards[3]/COLOR_NUM) ||
		(cards[3]/COLOR_NUM == cards[2]/COLOR_NUM && cards[3]/COLOR_NUM == cards[4]/COLOR_NUM && cards[2]/COLOR_NUM == cards[5]/COLOR_NUM && cards[1] < BLACK_JOKER_INDEX)

}

/**
四带两对
*/
func hasFourWithTwoPair(cards []int) bool {
	return cards[0]/COLOR_NUM == cards[1]/COLOR_NUM && cards[2]/COLOR_NUM == cards[3]/COLOR_NUM && cards[4]/COLOR_NUM == cards[5]/COLOR_NUM && cards[6]/COLOR_NUM == cards[7]/COLOR_NUM &&
		((cards[0]/COLOR_NUM == cards[2]/COLOR_NUM && cards[6] < BLACK_JOKER_INDEX) ||
			(cards[2]/COLOR_NUM == cards[4]/COLOR_NUM && cards[6] < BLACK_JOKER_INDEX && cards[1] < BLACK_JOKER_INDEX) ||
			(cards[4]/COLOR_NUM == cards[6]/COLOR_NUM && cards[1] < BLACK_JOKER_INDEX))
}

/**
单顺
*/
func hasSingleShun(cards []int) bool {
	size := len(cards)
	if size < 5 || size > 12 {
		return false
	}
	for i := 0; i < size-1; i++ {
		if CalcAbs(cards[i]/COLOR_NUM-cards[i+1]/COLOR_NUM) != 1 || cards[i] >= TWO_INDEX {
			return false
		}
	}
	return true
}

/**
对顺
*/
func hasPairShun(cards []int) bool {
	size := len(cards)
	if size < 6 || size%2 != 0 || size > 20 {
		return false
	}
	for i := 0; i < size-1; i += 2 {
		if cards[i] >= TWO_INDEX || cards[i+1] >= TWO_INDEX || cards[i]/COLOR_NUM != cards[i+1]/COLOR_NUM {
			return false
		}
		if i < size-2 && CalcAbs(cards[i]/COLOR_NUM-cards[i+2]/COLOR_NUM) != 1 {
			return false
		}
	}
	return true
}

/**
三顺
*/
func hasThreeShun(cards []int) bool {
	size := len(cards)
	if size < 9 || size%3 != 0 && size > 18 {
		return false
	}
	for i := 0; i < size-2; i += 3 {
		if cards[i] > TWO_INDEX || cards[i+1] > TWO_INDEX || cards[i+2] > TWO_INDEX ||
			cards[i]/COLOR_NUM != cards[i+1]/COLOR_NUM || cards[i]/COLOR_NUM != cards[i+2]/COLOR_NUM ||
			cards[i+1]/COLOR_NUM != cards[i+2]/COLOR_NUM {
			return false
		}
		if i < size-3 && CalcAbs(cards[i]/COLOR_NUM-cards[i+3]/COLOR_NUM) != 1 {
			return false
		}
	}
	return true
}

/**
三顺带单
*/
func hasThreeShunWithSingle(cards []int) bool {
	size := len(cards)
	if size < 8 && size > 18 {
		return false
	}
	shun := 0
	single_front := 0
	for i := 0; i < size-single_front*3; i++ {
		if cards[i]/COLOR_NUM != cards[i+1]/COLOR_NUM {
			single_front++
		} else {
			break
		}
	}
	single_rear := 0
	for i := single_front; i < size-2-single_rear; i += 3 {
		if cards[i] > TWO_INDEX || cards[i+1] > TWO_INDEX || cards[i+2] > TWO_INDEX ||
			cards[i]/COLOR_NUM != cards[i+1]/COLOR_NUM || cards[i]/COLOR_NUM != cards[i+2]/COLOR_NUM ||
			cards[i+1]/COLOR_NUM != cards[i+2]/COLOR_NUM {
			return false
		}
		if i < size-3-single_rear && CalcAbs(cards[i]/COLOR_NUM-cards[i+3]/COLOR_NUM) != 1 {
			return false
		}
		if cards[size-single_rear-1]/COLOR_NUM != cards[size-single_rear-2]/COLOR_NUM {
			single_rear++
		}
		shun++
		//fmt.Println(single_front,",",single_rear,",",shun)

	}
	if shun != single_front+single_rear {
		return false
	}
	WeightPop = cards[size-1-single_rear]/COLOR_NUM + 1
	return true
}

/**
三顺带对
*/
func hasThreeShunWithPair(cards []int) bool {
	size := len(cards)
	if size < 10 && size > 20 || cards[size-2] > BLACK_JOKER_INDEX {
		return false
	}
	pair_front := 0
	pair_rear := 0
	shun := 0
	for i := 0; i < size-5*pair_front; i += 2 {
		if cards[i]/COLOR_NUM == cards[i+1]/COLOR_NUM && cards[i]/COLOR_NUM != cards[i+2]/COLOR_NUM {
			pair_front++
		} else {
			break
		}
	}
	for i := pair_front * 2; i < size-2-pair_rear*2; i += 3 {
		if cards[i] > TWO_INDEX || cards[i+1] > TWO_INDEX || cards[i+2] > TWO_INDEX ||
			cards[i]/COLOR_NUM != cards[i+1]/COLOR_NUM ||
			cards[i]/COLOR_NUM != cards[i+2]/COLOR_NUM || cards[i+1]/COLOR_NUM != cards[i+2]/COLOR_NUM {
			return false
		}
		if i < size-3-pair_rear*2 && CalcAbs(cards[i]/COLOR_NUM-cards[i+3]/COLOR_NUM) != 1 {
			return false
		}
		if cards[size-pair_rear*2-1]/COLOR_NUM == cards[size-pair_rear*2-2]/COLOR_NUM && cards[size-pair_rear*2-1]/COLOR_NUM != cards[size-pair_rear*2-3]/COLOR_NUM {
			pair_rear++
		}
		shun++

		//fmt.Println(pair_front,",",pair_rear,",",shun)
	}
	if pair_front+pair_rear != shun {
		return false
	}
	//权重
	WeightPop = cards[size-1-2*pair_rear]/COLOR_NUM + 1
	return true
}

/**
飞机
*/
func hasAirplane(cards []int) bool {
	for i := 0; i < len(cards); i++ {
		if cards[i] >= TWO_INDEX {
			return false
		}
	}
	return cards[0]/COLOR_NUM == cards[2]/COLOR_NUM && cards[1]/COLOR_NUM == cards[2]/COLOR_NUM &&
		CalcAbs(cards[2]/COLOR_NUM-cards[3]/COLOR_NUM) == 1 &&
		cards[3]/COLOR_NUM == cards[4]/COLOR_NUM && cards[4]/COLOR_NUM == cards[5]/COLOR_NUM
}

/**
飞机带单
*/
func hasAirplaneWithSingle(cards []int) bool {
	return hasThreeShunWithSingle(cards)
}

/**
飞机带双
*/
func hasAirplaneWithPair(cards []int) bool {
	return hasThreeShunWithPair(cards)
}

/**
获得提示出牌,拆或者不拆
*/
/*
NONE 提示出牌
*/
func getNotePopCardsForTypeNONE(player_cards []int, last_note_weight int) string {
	//1.单，2.对，3.。。。
	for i := 0; i < len(player_cards); i++ {
		if player_cards[i]/COLOR_NUM+1 > last_note_weight {
			return fmt.Sprintf("%v", player_cards[i])
		}
	}
	return ""
}

/**
单 ONE
*/
func getNotePopCardsForTypeONEPure(player_cards []int, last_note_weight, last_pop_weight int) string {
	weight := 0
	size := len(player_cards)
	isSingle := true
	for i := 0; i < size-1; i++ {
		weight = player_cards[i]/COLOR_NUM + 1
		if weight > last_note_weight && weight > last_pop_weight && player_cards[i] < BLACK_JOKER_INDEX-1 {
			if weight != player_cards[i+1]/COLOR_NUM+1 {
				if isSingle {
					return fmt.Sprintf("%v", player_cards[i])
				} else {
					if i == size-2 {
						return fmt.Sprintf("%v", player_cards[i+1])
					}
					isSingle = true
				}
			} else {
				isSingle = false
			}
		}
	}
	return ""
}

/**
带单 ONE
*/
func getNotePopCardsForTypeONE(player_cards []int, last_note_weight, last_pop_weight int) string {
	weight := 0
	size := len(player_cards)
	isSingle := true
	for i := 0; i < size-1; i++ {
		weight = player_cards[i]/COLOR_NUM + 1
		if weight > last_note_weight && weight > last_pop_weight && player_cards[i] != BLACK_JOKER_INDEX+1 {
			if weight != player_cards[i+1]/COLOR_NUM+1 {
				if isSingle {
					return fmt.Sprintf("%v", player_cards[i])
				} else {
					if i == size-2 {
						return fmt.Sprintf("%v", player_cards[i+1])
					}
					isSingle = true
				}
			} else {
				isSingle = false
			}
		}
	}
	one := ""
	two := strings.Split(getNotePopCardsForTypePAIR(player_cards, last_note_weight, last_pop_weight), ",")
	one = two[0]
	if len(one) < 1 {
		three := strings.Split(getNotePopCardsForTypeTHREE(player_cards, last_note_weight, last_pop_weight), ",")
		one = three[0]
	}
	if len(one) < 1 {
		return ""
	}
	return one
}

/**
对 PAIR
*/
func getNotePopCardsForTypePAIR(player_cards []int, last_note_weight, last_pop_weight int) string {
	weight := 0
	size := len(player_cards)
	for i := 0; i < size-1; i++ {
		weight = player_cards[i]/COLOR_NUM + 1
		if weight > last_note_weight && weight > last_pop_weight {
			if weight == player_cards[i+1]/COLOR_NUM+1 && player_cards[i] < BLACK_JOKER_INDEX {
				if i < size-2 && weight == player_cards[i+2]/COLOR_NUM+1 {
					if i < size-3 && weight == player_cards[i+2]/COLOR_NUM+1 {
						i++
					}
					i++
				} else if i == size-2 || weight != player_cards[i+2]/COLOR_NUM+1 {
					return fmt.Sprintf("%v,%v", player_cards[i], player_cards[i+1])
				}
			}
		}
	}
	two := ""
	three := strings.Split(getNotePopCardsForTypeTHREE(player_cards, last_note_weight, last_pop_weight), ",")
	if len(three) == 3 {
		two = fmt.Sprintf("%v,%v", three[0], three[1])
	}
	return two
}

/**
拆三条
*/
func getNotePopCardsForTypePAIRPure(player_cards []int, last_note_weight, last_pop_weight int) string {
	weight := 0
	size := len(player_cards)
	for i := 0; i < size-1; i++ {
		weight = player_cards[i]/COLOR_NUM + 1
		if weight > last_note_weight && weight > last_pop_weight {
			if weight == player_cards[i+1]/COLOR_NUM+1 && player_cards[i] < BLACK_JOKER_INDEX {
				if i < size-2 && weight == player_cards[i+2]/COLOR_NUM+1 {
					if i < size-3 && weight == player_cards[i+2]/COLOR_NUM+1 {
						i++
					}
					i++
				} else if i == size-2 || weight != player_cards[i+2]/COLOR_NUM+1 {
					return fmt.Sprintf("%v,%v", player_cards[i], player_cards[i+1])
				}
			}
		}
	}
	//weight:=0
	//for i := 0; i < len(player_cards)-1; i++ {
	//	weight = player_cards[i]/COLOR_NUM+1
	//	if weight==player_cards[i+1]/COLOR_NUM+1&&player_cards[i]<BLACK_JOKER_INDEX&&weight>last_pop_weight&&weight>last_note_weight {
	//		return fmt.Sprintf("%v,%v",player_cards[i],player_cards[i+1])
	//	}
	//}
	return ""
}

/**
三 THREE
*/
func getNotePopCardsForTypeTHREE(player_cards []int, last_note_weight, last_pop_weight int) string {
	weight := 0
	size := len(player_cards)
	for i := 0; i < size-2; i++ {
		weight = player_cards[i]/COLOR_NUM + 1
		if weight > last_note_weight && weight > last_pop_weight {
			if weight == player_cards[i+1]/COLOR_NUM+1 && weight == player_cards[i+2]/COLOR_NUM+1 {
				if i < size-3 && weight == player_cards[i+3]/COLOR_NUM+1 {
					if i < size-4 && weight == player_cards[i+3] {
						i++
					}
					i++
				} else if i == size-3 || weight != player_cards[i+3]/COLOR_NUM+1 {
					return fmt.Sprintf("%v,%v,%v", player_cards[i], player_cards[i+1], player_cards[i+2])
				}
				i++
			}
		}
	}

	//weight:=0
	//for i := 0; i < len(player_cards)-2; i++ {
	//	weight = player_cards[i]/COLOR_NUM+1
	//	if weight==player_cards[i+1]/COLOR_NUM+1&&weight==player_cards[i+2]/COLOR_NUM+1&&
	//		player_cards[i]<BLACK_JOKER_INDEX&&weight>last_pop_weight&&weight>last_note_weight {
	//		return fmt.Sprintf("%v,%v,%v",player_cards[i],player_cards[i+1],player_cards[i+1])
	//	}
	//}
	return ""
}

/**
四	BOMB
*/
func getNotePopCardsForTypeBOMB(player_cards []int, last_note_weight, last_pop_weight int) string {
	weight := 0
	for i := 0; i < len(player_cards)-3; i++ {
		weight = player_cards[i]/COLOR_NUM + 1
		if weight == player_cards[i+1]/COLOR_NUM+1 && weight == player_cards[i+2]/COLOR_NUM+1 && weight == player_cards[i+3]/COLOR_NUM+1 &&
			player_cards[i] < BLACK_JOKER_INDEX && weight > last_pop_weight && weight > last_note_weight {
			return fmt.Sprintf("%v,%v,%v,%v", player_cards[i], player_cards[i+1], player_cards[i+2], player_cards[i+3])
		}
	}
	return ""
}

/**
三带一	ThreeWithOne
*/
func getNotePopCardsForTypeThreeWithOne(player_cards []int, last_note_weight, last_pop_weight int) string {
	three := getNotePopCardsForTypeTHREE(player_cards, last_note_weight, last_pop_weight)
	if len(three) < 1 {
		return ""
	}
	one := ""
	cards_remove := RemoveStringCards(player_cards, three)
	one = getNotePopCardsForTypeONE(cards_remove, 0, 0)
	if len(one) < 1 {
		return ""
	}
	return fmt.Sprintf("%v,%v", three, one)
}

/**
三带一对	ThreeWithPair
*/
func getNotePopCardsForTypeThreeWithPair(player_cards []int, last_note_weight, last_pop_weight int) string {
	three := getNotePopCardsForTypeTHREE(player_cards, last_note_weight, last_pop_weight)
	if len(three) < 1 {
		return ""
	}
	pair := getNotePopCardsForTypePAIR(RemoveStringCards(player_cards, three), 0, 0)
	if len(pair) < 1 {
		return ""
	}
	return fmt.Sprintf("%v,%v", three, pair)
}

/**
四带二	FourWithSingle
*/
func getNotePopCardsForTypeFourWithSingle(player_cards []int, last_note_weight, last_pop_weight int) string {
	four := getNotePopCardsForTypeBOMB(player_cards, last_note_weight, last_pop_weight)
	if len(four) < 1 {
		return ""
	}
	//四带单
	cards_remove := RemoveStringCards(player_cards, four)
	one := getNotePopCardsForTypeONE(cards_remove, 0, 0)
	if len(one) < 1 {
		return ""
	}
	one_int := StringToIntCards(one)
	one_str := fmt.Sprintf("%v", one)
	for i := 0; i < len(cards_remove); i++ {
		if cards_remove[i]/COLOR_NUM == one_int[0]/COLOR_NUM && cards_remove[i] != one_int[0] {
			one_str = fmt.Sprintf("%v,%v", one_str, cards_remove[i])
		} else if cards_remove[i]/COLOR_NUM > one_int[0]/COLOR_NUM {
			break
		}
	}
	cards_remove = RemoveStringCards(player_cards, one_str)
	other := getNotePopCardsForTypeONE(cards_remove, 0, 0)
	if len(one) > 0 && len(other) > 0 {
		return fmt.Sprintf("%v,%v,%v", four, one, other)
	}
	return ""
}

/**
四带二对	FourWithTwoPair
*/
func getNotePopCardsForTypeFourWithTwoPair(player_cards []int, last_note_weight, last_pop_weight int) string {
	four := getNotePopCardsForTypeBOMB(player_cards, last_note_weight, last_pop_weight)
	if len(four) < 1 {
		return ""
	}
	//四带两对
	cards_remove := RemoveStringCards(player_cards, four)
	one := getNotePopCardsForTypePAIR(cards_remove, 0, 0)
	if len(one) < 1 {
		return ""
	}
	one_int := StringToIntCards(one)
	one_str := fmt.Sprintf("%v", one_int[0])
	for i := 0; i < len(cards_remove); i++ {
		if cards_remove[i]/COLOR_NUM == one_int[0]/COLOR_NUM && cards_remove[i] != one_int[0] {
			one_str = fmt.Sprintf("%v,%v", one_str, cards_remove[i])
		} else if cards_remove[i]/COLOR_NUM > one_int[0]/COLOR_NUM {
			break
		}
	}
	cards_remove = RemoveStringCards(cards_remove, one_str)
	other := getNotePopCardsForTypePAIR(cards_remove, 0, 0)
	if len(other) < 1 {
		return ""
	}
	return fmt.Sprintf("%v,%v,%v", four, one, other)
}

/**
王炸
*/
func getNotePopCardsForTypeRocket(player_cards []int) string {
	for i := 0; i < len(player_cards)-1; i++ {
		if player_cards[i] == BLACK_JOKER_INDEX && player_cards[i+1] == BLACK_JOKER_INDEX+1 {
			return fmt.Sprintf("%v,%v", BLACK_JOKER_INDEX, BLACK_JOKER_INDEX+1)
		}
	}
	return ""
}

/**
单顺
TYPE_SINGLE_STRAIGHT      = 8
	TYPE_DOUBLE_STRAIGHT      = 9
	TYPE_AIRPLANE	= 10
	TYPE_AIRPLANE_WITH_SINGLE = 11
	TYPE_AIRPLANE_WITH_PAIR   = 12
	TYPE_THREE_SHUN		= 13
	TYPE_THREE_SHUN_WITH_SINGLE = 14
	TYPE_THREE_SHUN_WITH_PAIR = 15
*/
func getNotePopCardsForTypeSingleStraight(player_cards []int, last_note_weight, last_pop_weight, last_note_size, last_pop_size int) string {
	size_player_cards := len(player_cards)
	if size_player_cards < last_pop_size || size_player_cards < last_note_size {
		return ""
	}
	shun := ""
	flag := 0
	weight := 0
	for i := 0; i < size_player_cards-1; i++ {
		weight = player_cards[i]/COLOR_NUM + 1
		if weight > last_pop_weight && weight > last_note_weight && player_cards[i+1] < TWO_INDEX {
			if CalcAbs(player_cards[i+1]/COLOR_NUM-player_cards[i]/COLOR_NUM) == 1 {
				if flag == 0 {
					shun = fmt.Sprintf("%v", player_cards[i])
				} else {
					shun = fmt.Sprintf("%v,%v", shun, player_cards[i])
				}
				flag++
				if flag >= 5 && (flag == last_note_size || flag == last_pop_size) {
					return shun
				}
				if i == size_player_cards-2 {
					shun = fmt.Sprintf("%v,%v", shun, player_cards[i+1])
					flag++
				}
				if flag >= 5 && (flag == last_note_size || flag == last_pop_size) {
					return shun
				}
			} else if CalcAbs(player_cards[i]/COLOR_NUM-player_cards[i+1]/COLOR_NUM) > 1 {
				shun = ""
				flag = 0
			}
		}
	}

	return ""
}

/**
连对
*/
func getNotePopCardsForTypeDoubleStraight(player_cards []int, last_note_weight, last_pop_weight, last_note_size, last_pop_size int) string {
	size_player_cards := len(player_cards)
	if size_player_cards < last_pop_size*2 || size_player_cards < last_note_size*2 {
		return ""
	}
	shun := ""
	flag := 0
	weight := 0
	for i := 0; i < size_player_cards-2; i++ {
		weight = player_cards[i]/COLOR_NUM + 1
		if weight > last_pop_weight && weight > last_note_weight && player_cards[i+2] < TWO_INDEX {
			if player_cards[i]/COLOR_NUM == player_cards[i+1]/COLOR_NUM &&
				CalcAbs(player_cards[i]/COLOR_NUM-player_cards[i+2]/COLOR_NUM) == 1 {
				if flag == 0 {
					shun = fmt.Sprintf("%v,%v", player_cards[i], player_cards[i+1])
				} else {
					shun = fmt.Sprintf("%v,%v,%v", shun, player_cards[i], player_cards[i+1])
				}
				flag++
				if flag >= 2 && (flag == last_note_size || flag == last_pop_size) {
					return shun
				}
				if i == size_player_cards-4 && player_cards[i+2]/COLOR_NUM == player_cards[i+3]/COLOR_NUM {
					shun = fmt.Sprintf("%v,%v,%v", shun, player_cards[i+2], player_cards[i+3])
					flag++
				}
				if flag >= 2 && (flag == last_note_size || flag == last_pop_size) {
					return shun
				}
				i++
			} else if CalcAbs(player_cards[i]/COLOR_NUM-player_cards[i+1]/COLOR_NUM) > 0 {
				shun = ""
				flag = 0
			} else if player_cards[i]/COLOR_NUM == player_cards[i+1]/COLOR_NUM &&
				CalcAbs(player_cards[i]/COLOR_NUM-player_cards[i+2]/COLOR_NUM) > 0 {
				shun = ""
				flag = 0
				i++
			} else if player_cards[i]/COLOR_NUM == player_cards[i+1]/COLOR_NUM &&
				player_cards[i]/COLOR_NUM == player_cards[i+2]/COLOR_NUM &&
				CalcAbs(player_cards[i]/COLOR_NUM-player_cards[i+3]/COLOR_NUM) > 1 {
				shun = ""
				flag = 0
				i += 2
			} else if player_cards[i]/COLOR_NUM == player_cards[i+1]/COLOR_NUM &&
				player_cards[i]/COLOR_NUM == player_cards[i+2]/COLOR_NUM &&
				player_cards[i]/COLOR_NUM == player_cards[i+3]/COLOR_NUM {
				shun = ""
				flag = 0
				i += 3
			}

		}
	}
	return ""
}

/**
三顺
*/
func getNotePopCardsForTypeThreeShun(player_cards []int, last_note_weight, last_pop_weight, last_note_size, last_pop_size int) string {
	size_player_cards := len(player_cards)
	if size_player_cards < last_pop_size*3 || size_player_cards < last_note_size*3 {
		return ""
	}
	shun := ""
	flag := 0
	weight := 0
	for i := 0; i < size_player_cards-3; i++ {
		weight = player_cards[i]/COLOR_NUM + 1
		if weight > last_pop_weight && weight > last_note_weight && player_cards[i+2] < TWO_INDEX {
			if player_cards[i]/COLOR_NUM == player_cards[i+1]/COLOR_NUM &&
				player_cards[i]/COLOR_NUM == player_cards[i+2]/COLOR_NUM &&
				CalcAbs(player_cards[i]/COLOR_NUM-player_cards[i+3]/COLOR_NUM) == 1 {
				if flag == 0 {
					shun = fmt.Sprintf("%v,%v,%v", player_cards[i], player_cards[i+1], player_cards[i+2])
				} else {
					shun = fmt.Sprintf("%v,%v,%v,%v", shun, player_cards[i], player_cards[i+1], player_cards[i+2])
				}
				flag++
				if flag >= 2 && (flag == last_note_size || flag == last_pop_size) {
					return shun
				}
				if i == size_player_cards-6 && player_cards[i+3]/COLOR_NUM == player_cards[i+4]/COLOR_NUM &&
					player_cards[i+3]/COLOR_NUM == player_cards[i+5]/COLOR_NUM {
					shun = fmt.Sprintf("%v,%v,%v,%v", shun, player_cards[i+3], player_cards[i+4], player_cards[i+5])
					flag++
				}
				if flag >= 2 && (flag == last_note_size || flag == last_pop_size) {
					return shun
				}
				i += 2
			} else if CalcAbs(player_cards[i]/COLOR_NUM-player_cards[i+1]/COLOR_NUM) > 0 {
				shun = ""
				flag = 0
			} else if player_cards[i]/COLOR_NUM == player_cards[i+1]/COLOR_NUM &&
				CalcAbs(player_cards[i]/COLOR_NUM-player_cards[i+2]/COLOR_NUM) > 1 {
				shun = ""
				flag = 0
				i++
			} else if player_cards[i]/COLOR_NUM == player_cards[i+1]/COLOR_NUM &&
				player_cards[i]/COLOR_NUM == player_cards[i+2]/COLOR_NUM &&
				CalcAbs(player_cards[i]/COLOR_NUM-player_cards[i+3]/COLOR_NUM) > 1 {
				shun = ""
				flag = 0
				i += 2
			} else if player_cards[i]/COLOR_NUM == player_cards[i+3]/COLOR_NUM {
				shun = ""
				flag = 0
				i += 3
			}
		}
	}
	return ""

}

/**
三顺带单
*/
func getNotePopCardsForTypeThreeShunWithSingle(player_cards []int, last_note_weight, last_pop_weight, last_note_size, last_pop_size int) string {
	size_player_cards := len(player_cards)
	if size_player_cards < last_pop_size*4 || size_player_cards < last_note_size*4 {
		return ""
	}
	shun := getNotePopCardsForTypeThreeShun(player_cards, last_note_weight, last_pop_weight, last_note_size, last_pop_size)
	if len(shun) < 1 {
		return ""
	}
	cards_remove := RemoveStringCards(player_cards, shun)
	flag := 0
	if last_pop_size > 0 {
		flag = last_pop_size
	}
	if last_note_size > 0 {
		flag = last_note_size
	}
	one := getNotePopCardsForTypeONE(cards_remove, 0, 0)
	if len(one) < 1 {
		return ""
	}
	shun = fmt.Sprintf("%v,%v", shun, one)
	one_int := StringToIntCards(one)
	one_str := fmt.Sprintf("%v", one)
	for i := 0; i < flag-1; i++ {
		for i := 0; i < len(cards_remove); i++ {
			if cards_remove[i]/COLOR_NUM == one_int[0]/COLOR_NUM && cards_remove[i] != one_int[0] {
				one_str = fmt.Sprintf("%v,%v", one_str, cards_remove[i])
			} else if cards_remove[i]/COLOR_NUM > one_int[0]/COLOR_NUM {
				break
			}
		}
		cards_remove = RemoveStringCards(cards_remove, one_str)
		one = getNotePopCardsForTypeONE(cards_remove, 0, 0)
		if len(one) < 1 {
			return ""
		}
		shun = fmt.Sprintf("%v,%v", shun, one)
		one_int = StringToIntCards(one)
		one_str = fmt.Sprintf("%v", one)
	}
	return shun

}

/**
三顺带对
*/
func getNotePopCardsForTypeThreeShunWithPair(player_cards []int, last_note_weight, last_pop_weight, last_note_size, last_pop_size int) string {
	size_player_cards := len(player_cards)
	if size_player_cards < last_pop_size*5 || size_player_cards < last_note_size*5 {
		return ""
	}
	shun := getNotePopCardsForTypeThreeShun(player_cards, last_note_weight, last_pop_weight, last_note_size, last_pop_size)
	if len(shun) < 1 {
		return ""
	}
	cards_remove := RemoveStringCards(player_cards, shun)
	flag := 0
	if last_pop_size > 0 {
		flag = last_pop_size
	}
	if last_note_size > 0 {
		flag = last_note_size
	}
	pair := getNotePopCardsForTypePAIR(cards_remove, 0, 0)
	if len(pair) < 1 {
		return ""
	}
	shun = fmt.Sprintf("%v,%v", shun, pair)
	pair_int := StringToIntCards(pair)
	pair_str := fmt.Sprintf("%v", pair_int[0])
	for i := 0; i < flag-1; i++ {
		log.D("shun:%v", shun)
		for i := 0; i < len(cards_remove); i++ {
			if cards_remove[i]/COLOR_NUM == pair_int[0]/COLOR_NUM && cards_remove[i] != pair_int[0] {
				pair_str = fmt.Sprintf("%v,%v", pair_str, cards_remove[i])
			} else if cards_remove[i]/COLOR_NUM > pair_int[0]/COLOR_NUM {
				break
			}
		}
		cards_remove = RemoveStringCards(cards_remove, pair_str)
		pair = getNotePopCardsForTypePAIR(cards_remove, 0, 0)
		if len(pair) < 1 {
			return ""
		}
		shun = fmt.Sprintf("%v,%v", shun, pair)
		pair_int = StringToIntCards(pair)
		pair_str = fmt.Sprintf("%v", pair_int[0])
	}
	return ""

}

/**
飞机
*/
func getNotePopCardsForTypeAirplane(player_cards []int, last_note_weight, last_pop_weight int) string {
	last_note_size := 0
	last_pop_size := 0
	if last_pop_weight > 0 {
		last_pop_size = 2
	}
	if last_note_weight >0 {
		last_note_size = 2
	}
	return getNotePopCardsForTypeThreeShun(player_cards, last_note_weight, last_pop_weight, last_note_size, last_pop_size)

}

/**
飞机带单
*/
func getNotePopCardsForTypeAirplaneWithSingle(player_cards []int, last_note_weight, last_pop_weight int) string {
	last_note_size := 0
	last_pop_size := 0
	if last_pop_weight > 0 {
		last_pop_size = 2
	}
	if last_note_weight >0{
		last_note_size = 2
	}
	return getNotePopCardsForTypeThreeShunWithSingle(player_cards, last_note_weight, last_pop_weight, last_note_size, last_pop_size)

}

/**
飞机带对
*/
func getNotePopCardsForTypeAirplaneWithPair(player_cards []int, last_note_weight, last_pop_weight int) string {
	last_note_size := 0
	last_pop_size := 0
	if last_pop_weight > 0 {
		last_pop_size = 2
	}
	if last_note_weight >0{
		last_note_size = 2
	}
	return getNotePopCardsForTypeThreeShunWithPair(player_cards, last_note_weight, last_pop_weight, last_note_size, last_pop_size)

}

//转换
func IntToStringCards(cards []int) string {
	sort.Ints(cards)
	var cards_str string = ""
	for i := 0; i < len(cards); i++ {
		if i == len(cards)-1 {
			cards_str += fmt.Sprintf("%v", cards[i])
		} else {
			cards_str += fmt.Sprintf("%v,", cards[i])
		}
	}
	return cards_str
}
func StringToIntCards(cards string) []int {
	if len(cards)<1 ||cards == LC_PASS{
		return nil
	}
	cards_array := strings.Split(cards, ",")
	cards_int := make([]int, len(cards_array))
	var err error
	for i := 0; i < len(cards_array); i++ {
		cards_int[i], err = strconv.Atoi(cards_array[i])
		if err != nil {
			log.E("StringToIntCards failed err(%v)", err)
			return nil
		}
	}
	return cards_int
}

/**
移除牌
*/
func RemoveStringCards(cards_int []int, cards_remove string) []int {
	cards_remove_int := StringToIntCards(cards_remove)
	if len(cards_remove_int) < 1 {
		return cards_int
	}
	if len(cards_int) < 1 {
		return cards_int
	}
	if len(cards_int) == len(cards_remove_int) {
		return nil
	}
	new_cards := make([]int, len(cards_int)-len(cards_remove_int))
	flag := 0
	for i := 0; i < len(cards_int); i++ {
		for j := 0; j < len(cards_remove_int); j++ {
			if cards_int[i] == cards_remove_int[j] {
				break
			}
			if j+1 == len(cards_remove_int) && cards_int[i] != cards_remove_int[j] {
				new_cards[flag] = cards_int[i]
				flag++
			}
		}
	}
	return new_cards
}

/**
求绝对值
*/
func CalcAbs(a int) (ret int) {
	ret = (a ^ a>>31) - a>>31
	return
}
