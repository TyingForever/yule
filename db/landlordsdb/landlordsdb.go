package landlordsdb

import (
	"yule/db"
	"github.com/Centny/gwf/log"
	"gopkg.in/mgo.v2/bson"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2"
	"strings"
	"fmt"
)

/**进入游戏
获取斗地主游戏大厅信息
 */
func FindLandlordsHall(selector bson.M)  (landlordSet *LandlordSet,err error) {
	err = db.C(CN_LANDLORD_SET).Find(nil).Select(selector).One(&landlordSet)
	if err!=nil {
		log.E(TAG_LANDLORD_DB+",FindLandlordsHall,Find err(%v)",err)
		return nil,err
	}
	return landlordSet,err
}


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

/**
创建新房间
 */
func CreateNewRoom(uid string,categoryId int)(*LandlordInfo,error)  {
	room,err:=InitRoom(uid,categoryId)
	if err!=nil {
		return nil,err
	}
	_,err = db.C(CN_LANDLORD_INFO).Find(bson.M{"players.uid":uid}).Apply(mgo.Change{
		Update:bson.M{"$setOnInsert":room},
		ReturnNew:true,
		Upsert:true,

	},room)
	if err!=nil {
		log.E(TAG_LANDLORD_DB,",CreateNewRoom create err(%v)",err)
		return nil,err
	}
	return room,err
}

/**
随机安排房间
 */
func GetRoomRandom(uid string,cid int) (*LandlordInfo,error) {
	//判断用户是否在该房间
	var user_room []*LandlordInfo
	err:=db.C(CN_LANDLORD_INFO).Find(bson.M{"players.uid":uid,"status":bson.M{"$ne":LS_INVALID}}).All(&user_room)
	if err != nil {
		log.E(TAG_LANDLORD_DB+" Find landlordInfo err(%v)",err)
		return nil,err
	}
	if len(user_room)>0 {
		log.E(TAG_LANDLORD_DB+" user exists in landlord room")
		return nil,util.Err(TAG_LANDLORD_DB+" user exists in landlord room")
	}
	//用户不在斗地主房间则开始找房间
	landlordInfo,err:=FindRoomV("","","",cid,LS_QUEUE,nil)
	if err!=nil {
		if strings.Contains(err.Error(),"not found") {
			return CreateNewRoom(uid,cid)
		}
		return nil,err
	}
	landlord_set,err := FindLandlordsHall(nil)
	if err!=nil {
		log.E(TAG_LANDLORD_DB+" FindLandlordsHall err(%v)",err)
		return nil,err
	}
	//如果找到房间则添加用户至该房间,先判断是否等待超时
	if  util.Now() -landlordInfo.QueueTime >landlord_set.QueueTime {
		//解散该房间
		err=db.C(CN_LANDLORD_INFO).Remove(bson.M{"_id":landlordInfo.Id})
		if err!=nil {
			log.E(TAG_LANDLORD_DB+" Remove the room err %v",err)
			return nil,err
		}
		//解散了并通知，消息推送
		return nil,util.Err(TAG_LANDLORD_DB+"queue over time")
	}
	//添加玩家
	err=db.C(CN_LANDLORD_INFO).Update(bson.M{"_id":landlordInfo.Id},bson.M{"$push":bson.M{"players":bson.M{
		"uid":uid,"cards":"","pop_cards":LC_PASS,"pop_times":0,"status":LUS_ONLINE,"host":"127.0.0.1"}}})
	if err!=nil {
		log.E(TAG_LANDLORD_DB+" Add player error(%v)",err)
		return nil,err
	}
	info:=&LandlordInfo{
		Id:landlordInfo.Id,
		//Players:[]util.Map{util.Map{"uid":uid,"cards":"","pop_cards":"","pop_times":0,"status":1,"host":"127.0.0.1"}},
		Size:1,
	}

	//若添加成功则判断是否满员，满员则开始，明牌选择
	if  landlordInfo.Size+1 == PLAYER_NUM {
		log.D(TAG_LANDLORD_DB+ " GetRoomRandom full players,start choose ming cards")
		info.MingTime = util.Now()
		info.Status = LS_MING
		//info.PlayersPopTimes = []int{0,0,0}
	}
	err = ChangeLandlordInfo(uid,info)
	if err!=nil {
		log.E(TAG_LANDLORD_DB+ " GetRoomRandom change err(%v)",err)
		return nil,err
	}
	landlord_new,err:=FindRoomV(uid,landlordInfo.Id,"",0,0,nil)
	if err!=nil {
		log.E(TAG_LANDLORD_DB+ " GetRoomRandom get landlord after change failed err(%v)",err)
		return nil,err
	}
	return landlord_new,err
}

/**
找房间
 */
func FindRoomV(uid,lid ,level string,categoryId,status int ,selector bson.M)  (landlordInfo *LandlordInfo,err error){
	and:=[]bson.M{}
	if len(lid)>0 {
		and = append(and,bson.M{"_id":lid})
	}
	if len(uid)>0 {
		and = append(and,bson.M{"players.uid":uid})
	}
	if status>0 {
		and = append(and,bson.M{"status":status})
	}
	if len(level)>0 {
		and = append(and,bson.M{"category.level":level})
	}
	if categoryId>0 {
		and = append(and,bson.M{"category.cid":categoryId})
	}
	if len(and) < 1 {
		return  nil, util.Err("at last one arguments must be setted ")
	}
	var fargs = bson.M{"$and":and}
	var Q = db.C(CN_LANDLORD_INFO).Find(fargs).Select(selector)
	Q.Sort("-size")
	err = Q.One(&landlordInfo)
	if err!=nil {
		err = util.Err("FindRoomV list room error(%v) by selector(%v),args(%v)", err, util.S2Json(selector), util.S2Json(fargs))
		log.E(TAG_LANDLORD_DB+" err(%v)",err)
		return nil,err
	}
	return landlordInfo,nil
}


/**
玩家操作
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
func OperateLandlordInfo(uid,popCards string,landLordInfo *LandlordInfo,operate,ming_type int64) (*LandlordInfo,error) {
	if len(landLordInfo.Id)<1 {
		return nil,util.Err("Landlord Id is null")
	}
	if operate<1 {
		return nil,util.Err("Landlord operate is null")
	}
	multiple:=landLordInfo.Multiple
	if len(multiple)<1 {
		multiple = util.Map{}
	}
	if popCards==LC_PASS {
		popCards =""
	}
	var info *LandlordInfo
	err:=db.C(CN_LANDLORD_INFO).Find(bson.M{"_id":landLordInfo.Id}).One(&info)
	if len(info.Id)<1 {
		log.E("LandlordInfo find err %v",err)
		return nil,err
	}
	players:=[]util.Map{util.Map{"uid":info.Players[0].StrVal("uid")},util.Map{"uid":info.Players[1].StrVal("uid")},util.Map{"uid":info.Players[2].StrVal("uid")}}
	switch operate {
	case OP_MING://明牌
		if (ming_type == LM_TD_BEFORE_DEAL&&util.Now()-info.MingTime>LT_MING)||info.Status<LS_MING||info.Status>LS_FIGHTING {
			log.E("ming err")
			return nil,util.Err("ming err")
		}
		for i := 0; i < len(players); i++ {
			if players[i].StrVal("uid") == uid {
				players[i]["ming"] = LM_TD_BEFORE_DEAL
			}
		}
		if ming_type == LM_TD_BEFORE_DEAL&&util.Now()-info.MingTime<=LT_MING {
			multiple[LD_MING] = LM_TD_BEFORE_DEAL
		}else {
			multiple[LD_MING] = ming_type
		}
		landLordInfo.Multiple = multiple
		landLordInfo.Players = players
		err:=ChangeLandlordInfo(uid,landLordInfo)
		if err!=nil {
			log.E(TAG_LANDLORD_DB+" ChangeLandlordInfo err(%v)",err)
			return nil,err
		}
		break
	case OP_DEAL://发牌
		if info.Status != LS_MING||util.Now()-info.MingTime<=LT_MING {
			log.E("deal err")
			return nil,util.Err("deal err")
		}
		players,landlord_cards:=RandomDistributeCards(info.Players)//发牌
		landLordInfo.Players = players
		turn_user:=info.Players[ConfirmFirstPlayer()].StrVal("uid")
		landLordInfo.LandlordCards = landlord_cards
		landLordInfo.TurnUser = turn_user
		landLordInfo.Status = LS_GRABBING
		landLordInfo.GrabTime = util.Now()
		err:=ChangeLandlordInfo(uid,landLordInfo)
		if err!=nil {
			log.E(TAG_LANDLORD_DB+" ChangeLandlordInfo err(%v)",err)
			return nil,err
		}
		break
	case OP_GRAB:
		//进入抢地主 必须是在抢地主时间内，必须是轮到当前用户，状态为抢地主，每个用户最多只能轮流一次（即总共最多三次），抢地主分数必须比之前的分数要高
		if util.Now() - info.GrabTime > LT_GRAB || uid != info.TurnUser || info.Status != LS_GRABBING||info.OperateNum>=PLAYER_NUM ||
			landLordInfo.Multiple.IntVal(LD_GRAB_SCORE) <= info.Multiple.IntVal(LD_GRAB_SCORE){
			log.E("grab err")
			return nil,util.Err("grab err")
		}
		//抢地主规则
		landLordInfo = GrabLandlordRule(uid,landLordInfo,info)
		err = ChangeLandlordInfo(uid,landLordInfo)
		if err!=nil {
			log.E(TAG_LANDLORD_DB+" ChangeLandlordInfo err(%v)",err)
			return nil,err
		}
		break
	case OP_DOUBLE:
		//加倍 必须是在加倍时间内，必须是轮到当前用户，状态为抢地主，每个用户最多只能轮流一次（即总共最多三次），抢地主分数必须比之前的分数要高
		if util.Now() - info.DoubleTime > LT_DOUBLE || uid != info.TurnUser || info.Status != LS_DOUBLING||info.OperateNum>=PLAYER_NUM {
			log.E("double err")
			return nil,util.Err("double err")
		}
		//加倍规则
		landLordInfo = DoubleLandlordRule(uid,landLordInfo,info)
		err = ChangeLandlordInfo(uid,landLordInfo)
		if err!=nil {
			log.E(TAG_LANDLORD_DB+" ChangeLandlordInfo err(%v)",err)
			return nil,err
		}
		break
	case OP_POP_CARD:
		//出牌
		if util.Now()-info.PopCardTime>LT_POP_CARD||uid!=info.TurnUser||info.Status!=LS_FIGHTING{
			log.E("pop card err")
			return nil,util.Err("pop card err")
		}
		//根据上家出牌以及中途出现的不出牌开启新回合
		lastPopCards:=info.LastPopCards
		for i := 0; i < PLAYER_NUM; i++ {
			if info.Players[i].StrVal("uid")==uid {
				if info.Players[(i+1)%3].StrVal("pop_cards") == LC_PASS&&
					info.Players[(i+2)%3].StrVal("pop_cards") == LC_PASS{//假如另外两家都不出牌，则轮到自己为首发出牌
					lastPopCards = ""
				}
			}
		}
		if info.LastPopCards==LC_PASS {
			lastPopCards = ""
		}
		//不符合出牌规则
		if !IsMeetPopLogic(popCards,lastPopCards){
			log.E("pop card err")
			return nil,util.Err("pop card err")
		}
		landLordInfo = FightLandlordRule(uid,popCards,landLordInfo,info)
		if landLordInfo!=nil {
			landLordInfo.NotePopCards = ""
			err = ChangeLandlordInfo(uid,landLordInfo)
			if err!=nil {
				log.E(TAG_LANDLORD_DB+" ChangeLandlordInfo err(%v)",err)
				return nil,err
			}
		}

	case OP_PASS_CARD:
		if util.Now()-info.PopCardTime>LT_POP_CARD||uid!=info.TurnUser||info.Status!=LS_FIGHTING{
			log.E("pop card err")
			return nil,util.Err("pass card err")
		}
		landLordInfo = PassCardRule(uid,landLordInfo,info)
		landLordInfo.NotePopCards = ""
		err = ChangeLandlordInfo(uid,landLordInfo)
		if err!=nil {
			log.E(TAG_LANDLORD_DB+" ChangeLandlordInfo err(%v)",err)
			return nil,err
		}
		break
	case OP_GET_NOTE:
		if util.Now()-info.PopCardTime>LT_POP_CARD||uid!=info.TurnUser||info.Status!=LS_FIGHTING{
			log.E("get note card err")
			return nil,util.Err("get note card err")
		}
		noteCards:=GetNoteCardsRule(uid,info)
		if len(noteCards)<1 {
			noteCards = LC_PASS
		}
		landLordInfo.NotePopCards =noteCards
		err = ChangeLandlordInfo(uid,landLordInfo)
		if err!=nil {
			log.E(TAG_LANDLORD_DB+" ChangeLandlordInfo err(%v)",err)
			return nil,err
		}
		break
	case OP_CONTINUE:
		if util.Now()-info.OverGameTime>LT_GAME_OVER_SHOW||info.Status!=LS_SHOW_RESULT {
			log.E("continue err")
			return nil,util.Err("continue card err")
		}
		//初始化
		landLordInfo = InitRoomContinue(uid , info)
		err = ChangeLandlordInfo(uid,landLordInfo)
		if err!=nil {
			log.E(TAG_LANDLORD_DB+" ChangeLandlordInfo err(%v)",err)
			return nil,err
		}
	case OP_RETURN_HALL:
		if info.Status!=LS_SHOW_RESULT {
			log.E("return hall err")
			return nil,util.Err("return hall err")
		}
		err=db.C(CN_LANDLORD_INFO).Update(bson.M{"_id":info.Id,"players.uid":uid,"players.status":LUS_ONLINE},bson.M{
			"$pull":bson.M{
				"players.uid":uid,
			}})
		if err != nil {
			log.E(TAG_LANDLORD_DB+" ChangeLandlordInfo err(%v)",err)
			return nil,err
		}

	}
	info_new,err:=FindRoomV(uid,info.Id,"",0,0,nil)
	if err!=nil {
		log.E(TAG_LANDLORD_DB+" Find the landlordInfo err(%v)",err)
		return nil,err
	}
	return info_new,nil
}



/**
修改斗地主信息
 */
func ChangeLandlordInfo(uid string,landlordInfo *LandlordInfo) error {
	and:=[]bson.M{}
	if len(landlordInfo.Id)<1{
		return util.Err("landlord id is empty")
	}
	and = append(and,bson.M{"_id":landlordInfo.Id})
	arg:=bson.M{}
	fargs:=bson.M{}
	if len(landlordInfo.Players)>0 {
		for i := 0; i < len(landlordInfo.Players); i++ {
			for key,val:= range landlordInfo.Players[i]{
				name:=fmt.Sprintf("players.%v.%v",i,key)
				arg[name] = val
			}
		}
	}

	if landlordInfo.Size>0 {
		fargs["$inc"]=bson.M{"size":1}
	}

	//if len(landlordInfo.Users)>0 {
	//	fargs["$addToSet"] = bson.M{"users":landlordInfo.Users[0]}
	//	fargs["$inc"] = bson.M{"size":1}
	//	and = append(and,bson.M{"status":LS_QUEUE})
	//	and = append(and,bson.M{"queue_time":bson.M{"$gte":util.Now()-LT_QUEUE}})
	//	and = append(and,bson.M{"size":bson.M{"$lt":3}})
	//}
	//if len(landlordInfo.ACards)>0 {
	//	arg["a_cards"] = landlordInfo.ACards
	//}
	//if len(landlordInfo.BCards)>0 {
	//	arg["b_cards"] = landlordInfo.BCards
	//}
	//if len(landlordInfo.CCards)>0 {
	//	arg["c_cards"] = landlordInfo.CCards
	//}
	//if len(landlordInfo.PlayersPopTimes)>0 {
	//	arg["players_pop_times"] = landlordInfo.PlayersPopTimes
	//}
	//if len(landlordInfo.MingCardUsers)>0 {
	//	fargs["$addToSet"] = bson.M{"users":landlordInfo.MingCardUsers[0]}
	//}
	//if len(landlordInfo.ACardsPop)>0 {
	//	arg["a_cards_pop"] = landlordInfo.ACardsPop
	//}
	//if len(landlordInfo.BCardsPop)>0 {
	//	arg["b_cards_pop"] = landlordInfo.BCardsPop
	//}
	//if len(landlordInfo.BCardsPop)>0 {
	//	arg["c_cards_pop"] = landlordInfo.CCardsPop
	//}

	if len(landlordInfo.LastPopCards)>0 {
		arg["last_pop_cards"] = landlordInfo.LastPopCards
	}
	if len(landlordInfo.LandlordCards)>0 {
		arg["landlord_cards"] = landlordInfo.LandlordCards
	}
	if len(landlordInfo.LandlordUser)>0 {
		arg["landlord_user"] = landlordInfo.LandlordUser
		and = append(and,bson.M{"$or":[]bson.M{bson.M{"status":LS_MING},bson.M{"status":LS_GRABBING}}})
	}
	if landlordInfo.Status>0 {
		arg["status"] = landlordInfo.Status
	}
	if len(landlordInfo.Multiple)>0 {
		for key, val := range landlordInfo.Multiple {
			if key==LD_DOUBLE_USERS {
				fargs["$addToSet"] = bson.M{"multiple."+LD_DOUBLE_USERS:landlordInfo.Multiple.StrVal(LD_DOUBLE_USERS)}
			}else {
				arg["multiple."+key] = val
			}
			//if key == LD_MING{//若为明牌则需要判断有没有明牌
			//	and = append(and,bson.M{"multiple."+key:{"$exists":false}})
			//}
		}
	}

	if landlordInfo.QueueTime>0 {
		arg["queue_time"] = landlordInfo.QueueTime
	}
	if landlordInfo.OverGameTime>0 {
		arg["over_game_time"] = landlordInfo.OverGameTime
	}
	if landlordInfo.DoubleTime > 0 {
		arg["double_time"] = landlordInfo.DoubleTime
	}
	if landlordInfo.PopCardTime >0  {
		arg["pop_card_time"] = landlordInfo.PopCardTime
	}
	if landlordInfo.GrabTime > 0{
		arg["grab_time"] = landlordInfo.GrabTime
	}
	if landlordInfo.MingTime > 0{
		arg["ming_time"] = landlordInfo.MingTime
	}
	if len(landlordInfo.TurnUser)>0 {
		arg["turn_user"] = landlordInfo.TurnUser
	}
	if landlordInfo.OperateNum >-1 {
		arg["operate_num"] = landlordInfo.OperateNum
	}
	if len(landlordInfo.NotePopCards)>0 {
		arg["note_pop_cards"] = landlordInfo.NotePopCards
	}

	if len(and) < 1 {
		return util.Err("at last one arguments must be setted on landlord struct, but lid(%v)",landlordInfo.Id)
	}

	query:=bson.M{"$and":and}
	arg["last"] = util.Now()
	fargs["$set"] = arg

	if SHOW_LOG_DB {
		log.D(TAG_LANDLORD_DB+"ChangeLandlordInfo query:%v ,args: %v",util.S2Json(query),util.S2Json(fargs))
	}
	err:=db.C(CN_LANDLORD_INFO).Update(query,fargs)
	if err != nil {
		log.E(TAG_LANDLORD_DB+"ChangeLandlordInfo err: %v",err)
		return err
	}
	return nil

}


/**
初始化房间
 */
func InitRoom(uid string,categoryId int) (*LandlordInfo,error) {
	set,err:=FindLandlordsHall(bson.M{"categories":1})
	if err!=nil {
		return nil,err
	}
	category:=util.Map{}
	for i := 0; i<len(set.Categories); i++ {
		if (int)(set.Categories[i].IntVal("cid")) == categoryId{
			category = set.Categories[i]
		}
	}

	players:=[]util.Map{util.Map{"uid":uid,"cards":"","pop_cards":LC_PASS,"pop_times":0,"status":1,"host":"127.0.0.1"}}
	landlord:=&LandlordInfo{
		Id:bson.NewObjectId().Hex(),
		LastPopCards:LC_PASS,
		NotePopCards:LC_PASS,
		Players:players,
		Status:LS_QUEUE,
		Category:category,
		Last:util.Now(),
		Time:util.Now(),
		QueueTime:util.Now(),
		Size:1,
	}
	log.D("landlord room:%v",util.S2Json(landlord))
	return landlord,nil
}
func InitRoomContinue(uid string,info *LandlordInfo) *LandlordInfo {
	//若添加成功则判断是否满员，满员则开始，明牌选择
	if  info.Size+1 == PLAYER_NUM {
		log.D(TAG_LANDLORD_DB+ " GetRoomRandom full players,start choose ming cards")
		info.MingTime = util.Now()
		info.Status = LS_MING
		info.Size+=1
	}else {
		//不满人
		info.Size+=1
	}
	return info
}