package landlordsdb

import (
	"github.com/Centny/gwf/log"
	"testing"
	"github.com/Centny/gwf/util"
)

//测试提示出牌规则
func TestGetNoteCards(t *testing.T) {
	player_cards:=[]int{0,1,2,3,4,8,9,16,17,20}
	note_cards:=getNotePopCardsForTypeNONE(player_cards,0)
	log.D("Note Cards %v",note_cards)

	player_cards=[]int{4,5,6,7,8,9,10,11,12,13,14}
	note_cards=getNotePopCardsForTypeONE(player_cards,0,0)
	log.D("Note Cards %v",note_cards)

	//player_cards=[]int{0,1,2,4,5}
	//player_cards=[]int{0,4,5}
	//player_cards=[]int{0,1,2,4,5,8,9,10}
	//player_cards=[]int{0,4,8,9}
	player_cards=[]int{0,1,2,4,5,6}
	player_cards=[]int{0,1,2,3,4,5,6,7,9,10}
	note_cards=getNotePopCardsForTypePAIR(player_cards,0,0)
	log.D("Note Cards %v",note_cards)

	player_cards=[]int{0,1,2,4,5,6}
	//player_cards=[]int{0,1,4,5,6}
	//player_cards=[]int{0,1,4,8,9,10}
	note_cards=getNotePopCardsForTypeTHREE(player_cards,0,0)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,2}
	note_cards=getNotePopCardsForTypeTHREE(player_cards,0,0)
	log.D("Note Cards %v",note_cards)

	player_cards=[]int{0,1,2,4,5,6,7}
	note_cards=getNotePopCardsForTypeBOMB(player_cards,0,0)
	log.D("Note Cards %v",note_cards)

	player_cards=[]int{0,1,2,4,5,6,7,8}
	note_cards=getNotePopCardsForTypeThreeWithOne(player_cards,0,0)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,2,4,5,6,7,8,9}
	note_cards=getNotePopCardsForTypeThreeWithOne(player_cards,0,0)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,5,6,8}
	note_cards=getNotePopCardsForTypeThreeWithOne(player_cards,0,0)
	log.D("Note Cards %v",note_cards)

	player_cards=[]int{0,1,2,4,5,6,7,8,9}
	note_cards=getNotePopCardsForTypeThreeWithPair(player_cards,0,0)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,2,4,5,6,8,9}
	note_cards=getNotePopCardsForTypeThreeWithPair(player_cards,0,0)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,2,4,5,6}
	note_cards=getNotePopCardsForTypeThreeWithPair(player_cards,0,0)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,2,4,5,6,7}
	note_cards=getNotePopCardsForTypeThreeWithPair(player_cards,0,0)
	log.D("Note Cards %v",note_cards)

	player_cards=[]int{0,1,2,4,5,6,7,8,12}
	note_cards=getNotePopCardsForTypeFourWithSingle(player_cards,0,0)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,2,3,4,5,6,7,12,13}
	note_cards=getNotePopCardsForTypeFourWithSingle(player_cards,0,0)
	log.D("Note Cards %v",note_cards)


	player_cards=[]int{0,1,2,4,5,6,7,8,9,12,13}
	note_cards=getNotePopCardsForTypeFourWithTwoPair(player_cards,0,0)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,2,4,5,6,7,8,9}
	note_cards=getNotePopCardsForTypeFourWithTwoPair(player_cards,0,0)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,2,4,5,6,7,8,9,12,16}
	note_cards=getNotePopCardsForTypeFourWithTwoPair(player_cards,0,0)
	log.D("Note Cards %v",note_cards)

	player_cards=[]int{0,4,8,12,16}
	note_cards=getNotePopCardsForTypeSingleStraight(player_cards,0,0,5,5)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,2,4,5,12,16,20}
	note_cards=getNotePopCardsForTypeSingleStraight(player_cards,0,0,5,5)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,2,3,4,8,12,16,20}
	note_cards=getNotePopCardsForTypeSingleStraight(player_cards,0,0,5,5)
	log.D("Note Cards %v",note_cards)

	player_cards=[]int{0,1,4,5,8,9}
	note_cards=getNotePopCardsForTypeDoubleStraight(player_cards,0,0,3,3)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,4,5,8}
	note_cards=getNotePopCardsForTypeDoubleStraight(player_cards,0,0,3,3)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,2,8,9,12,13}
	note_cards=getNotePopCardsForTypeDoubleStraight(player_cards,0,0,3,3)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,2,8,9,12,13,16,17}
	note_cards=getNotePopCardsForTypeDoubleStraight(player_cards,0,0,3,3)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,4,5,8,12,13,16,17}
	note_cards=getNotePopCardsForTypeDoubleStraight(player_cards,0,0,3,3)
	log.D("Note Cards %v",note_cards)

	player_cards=[]int{0,1,2,4,5,6,8,9,10}
	size:=3
	note_cards=getNotePopCardsForTypeThreeShun(player_cards,0,0,size,size)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,2,3,4,5,6,8,9,10,12,13,14}
	note_cards=getNotePopCardsForTypeThreeShun(player_cards,0,0,size,size)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,4,5,6,8,9,10}
	note_cards=getNotePopCardsForTypeThreeShun(player_cards,0,0,size,size)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,2,3,4,5,6,8,9,10}
	note_cards=getNotePopCardsForTypeThreeShun(player_cards,0,0,size,size)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,2,8,9,10,12,13}
	note_cards=getNotePopCardsForTypeThreeShun(player_cards,0,0,size,size)
	log.D("Note Cards %v",note_cards)

	player_cards=[]int{0,1,2,4,5,6,8,9,10,12,13,14,15,16,20}
	size = 3
	note_cards=getNotePopCardsForTypeThreeShunWithSingle(player_cards,0,0,size,size)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,2,4,5,6,8,9,10,12,13,14,15,16,20,24}
	note_cards=getNotePopCardsForTypeThreeShunWithSingle(player_cards,0,0,size,size)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,2,4,5,8,9,10,12,13,14,16,17,18,20,24}
	note_cards=getNotePopCardsForTypeThreeShunWithSingle(player_cards,0,0,size,size)
	log.D("Note Cards %v",note_cards)

	player_cards=[]int{0,1,2,4,5,6,8,9,10,12,13,14,16,17,24,26}
	size = 3
	note_cards=getNotePopCardsForTypeThreeShunWithPair(player_cards,0,0,size,size)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,2,4,5,6,8,9,10,12,13,14,15,16,20,21,24,25}
	note_cards=getNotePopCardsForTypeThreeShunWithPair(player_cards,0,0,size,size)
	log.D("Note Cards %v",note_cards)
	player_cards=[]int{0,1,4,8,9,10,12,13,14,16,17,18,20,24}
	note_cards=getNotePopCardsForTypeThreeShunWithPair(player_cards,0,0,size,size)
	log.D("Note Cards %v",note_cards)
}

//测试提示出牌
func TestGetNoteCardsRule(t *testing.T) {
	info,err := FindRoomV("u2","","",0,0,nil)
	if err!=nil {
		t.Error(err)
		return
	}
	log.D("LandlordInfo:%v",util.S2Json(info))
	note:= GetNoteCardsRule(info.TurnUser,info)
	log.D("Note:%v",note)
}

func TestString(t *testing.T) {

	//InitData()
	//players:=[]string{"01","02","03"}
	////cur_player := "02"
	////发牌
	//RandomDistributeCards(players)
	//log.D("landlordsCards: %v",util.S2Json(LandlordCards))
	//log.D("playersCards: %v",util.S2Json(PlayersCards))
	////添加地主牌
	////landlord:="01"
	//LandlordPlayer = "01"
	//AddLordsCardsForLand(players)
	//log.D("playersCards: %v",util.S2Json(PlayersCards))

	//出牌
	//1.单		0
	card_pop := "51"
	card_last:=card_pop
	type_pop := TYPE_NONE
	weight_pop:=WEIGHT_NONE
	type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	card_pop = "52,53"
	type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	log.D("Pop Cards card_pop(%v) larger than card_last(%v) result: %v",card_pop,card_last,IsMeetPopLogic(card_pop,card_last))


	////2.对		0,1
	//card_pop = "0,1"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//
	////3.王炸		52,53
	//card_pop = "52,53"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//
	////4.三		0,1,2
	//card_pop = "0,1,2"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//
	////5.三带一	0,1,2,4
	//card_pop = "0,1,2,4"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//card_pop = "0,8,9,10"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//
	////6.四炸		0,1,2,3 xx
	//card_pop = "0,1,2,3"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	////
	////7.三带对	0,1,2,4,5
	//card_pop = "0,1,2,4,5"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type 33344: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//card_pop = "0,1,8,9,10"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type 33444: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//card_pop = "0,1,2,52,53"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type 333ww: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//
	////8.四带二单	0,1,2,3,4,8
	//card_pop="0,1,2,3,4,8"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type 333345: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//card_pop="0,8,12,13,14,15"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type 356666: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//card_pop="0,12,13,14,15,24"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type 366668: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//
	////9.四带两对	0,1,2,3,4,5,8,9 xx
	//card_pop="0,1,2,3,4,5,8,9"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type 33334455: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//card_pop="0,1,8,9,10,11,16,17"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type: 33555577 : %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//card_pop="0,1,8,9,16,17,18,19"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type: 33557777: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//
	////10.飞机	0,1,2,4,5,6
	//card_pop="0,1,2,4,5,6"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//card_pop="44,45,46,48,49,50"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//
	////11.飞机带单	0,1,2,4,5,6,8,12
	//card_pop="0,1,2,4,5,6,8,12"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//card_pop="0,44,45,46,48,49,50,52"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//
	////12飞机带队	0,1,2,4,5,6,8,9,12,13
	//card_pop="0,1,2,4,5,6,8,9,12,13"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//card_pop="0,1,4,5,44,45,46,48,49,50"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//
	////13.单顺	0,4,8,12,16,20
	//card_pop="0,4,8,12,16,20"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//card_pop="0,4,8,12,16,20,24,28,32,36,40,44,48"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//
	////14.对顺	0,1,4,5,8,9
	//card_pop="0,1,4,5,8,9"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//card_pop="40,41,44,45,48,49"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type: kk1122: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//log.D("Pop Cards Type:%v: %v,Weight: %v",util.S2Json(card_pop),AllTypes[type_pop],weight_pop)
	//
	////15.三顺	0,1,2,4,5,6,8,9,10
	//card_pop="0,1,2,4,5,6,8,9,10"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//card_pop="0,1,7,4,5,6,8,9,10"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//
	////16.三顺带单	0,1,2,4,5,6,8,9,10,12,16,20 xx
	//card_pop="0,1,2,4,5,6,8,9,10,12,16,20"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type 333444555678:%v: %v,Weight: %v",util.S2Json(card_pop),AllTypes[type_pop],weight_pop)
	//card_pop="0,4,8,12,13,14,16,17,18,20,21,22"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type 345666777888:%v: %v,Weight: %v",util.S2Json(card_pop),AllTypes[type_pop],weight_pop)
	//card_pop="0,4,8,9,10,12,13,14,16,17,18,44"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type 34555666777K:%v: %v,Weight: %v",util.S2Json(card_pop),AllTypes[type_pop],weight_pop)
	////wrong
	//card_pop="0,4,8,9,10,12,13,14,20,21,22,44"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type 34555666777K:%v: %v,Weight: %v",util.S2Json(card_pop),AllTypes[type_pop],weight_pop)
	////log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
	//
	////17.三顺带对	0,1,2,4,5,6,8,9,10,12,13,16,17,20,21 xx
	//card_pop="0,1,2,4,5,6,8,9,10,12,13,16,17,20,21"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type:%v: %v,Weight: %v",util.S2Json(card_pop),AllTypes[type_pop],weight_pop)
	//card_pop="0,1,4,5,8,9,12,13,14,16,17,18,20,21,22"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type 334455666777888:%v: %v,Weight: %v",util.S2Json(card_pop),AllTypes[type_pop],weight_pop)
	//card_pop="0,1,4,5,8,9,10,12,13,14,16,17,18,20,21"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type 334455566677788:%v: %v,Weight: %v",util.S2Json(card_pop),AllTypes[type_pop],weight_pop)
	////wrong
	//card_pop="0,1,2,4,5,6,8,9,10,12,13,16,17,20,24"
	//type_pop,weight_pop = JudgeCardsTypeAndWeight(card_pop)
	//log.D("Pop Cards Type:%v: %v,Weight: %v",util.S2Json(card_pop),AllTypes[type_pop],weight_pop)
	//
	//log.D("Pop Cards Type: %v,Weight: %v", AllTypes[type_pop],weight_pop)
}

