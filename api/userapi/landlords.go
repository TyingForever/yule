package userapi

import (
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"github.com/Centny/gwf/log"
	"yule/db/landlordsdb"
	"strings"
)

//进入游戏
//进入游戏，判断该用户是否还在处于游戏中，若处于游戏中则判断其是在哪个游戏中，并进入到指定游戏中
//@url,需求登录，GET请求
//	~/usr/api/enterLandlords		GET
//@arg,json对象中的参数以及query中参数
//	token		R	用户凭证
//	game_type	R	游戏类型
/*
	样例	~/pub/api/enterLandlords?token=xx&game_type=1
*/
//@ret,返回通用code/data
//	code	I	0：进入成功，1：参数错误，2：json body错误，3：进入游戏失败
//	entry_game_result O	进入游戏的结果
//	entry_game_result.result I 结果：1.进入游戏大厅，继续游戏
//	landlord_hall	O	斗地主大厅设置
//	sound_effect	I 	音效
//	music		I	声音
//	categories	O	场别分类
//	categories.title	S	标题
//	categories.level	S	级别
//	categories.score_end	F	低分
//	categories.gold_min	F	最低金额
//	queue_time	I	排队时间
//	grab_time	I	抢地主时间
//	double_time	I	加倍时间
//	pop_card_time	I	出牌时间
//	over_game_time	I	游戏结束结果显示时间
//	landlords_info	O	正处于斗地主中
//	lid		I	游戏信息id
//	users		A	游戏用户id list
//	turnUser	S	轮到的用户
//	landlord_card	A	地主牌
//	last_pop_card	A	上家出的牌
//	a_cards		A	a玩家的手牌
//	b_cards		A	b玩家的手牌
//	c_cards		C	c玩家的手牌
//	stats		I	当前游戏状态
//	category	O	当前场别
//	landlord_user	S	地主玩家
//	multiple	O	倍数
//	multiple.double A	加倍的用户
//	multiple.bomb	A	炸弹
//	multiple.spring	I	春天
//	multiple.anti_spring	I	反春天
//	last		I	修改时间
//	time		I	创建时间
//	queue_time	I	排队开始时间
//	grab_time	I	抢地主开始时间
//	double_time	I	加倍开始时间
//	pop_card_time	I	出牌开始时间
//	over_game_time	I	游戏结束结果开始显示时间
/*	样例
*/
//@tag,游戏,斗地主
//@author,zhnagyq,2017-08-02
//@case,yule
func EntryLandlords(hs *routing.HTTPSession) routing.HResult {
	var gameType int
	err:=hs.ValidF(`gameType,R|I,R:-1;`,&gameType)
	if err!=nil {
		log.E("arg-err,%v", err)
		return hs.MsgResErr(1, "arg-err", err)
	}
	uid:=hs.StrVal("uid")
	if LOG_API {
		log.D(TAG_LANDLORDS+"uid(%v),gameType(%v)",uid,gameType)
	}
	result,landlord,set,err:=landlordsdb.IsEntryLandlord(uid)
	if err != nil {
		return hs.MsgResErr(2, "服务器错误", err)
	}
	if result == landlordsdb.UGS_OFF{
		return hs.MsgRes(util.Map{"result":result,"landlordSet":set})
	}
	if result == landlordsdb.UGS_LANDLORDS{
		return hs.MsgRes(util.Map{"result":result,"landlordInfo":landlord})
	}
	return hs.MsgResErr(2,"服务器错误", err)
}

func DoEntryLandlords(token string,gameType int) (util.Map, error) {
	res, err := util.HGet2("%v/usr/api/entryLandlords?token=%v&gameType=%v", SrvAddr(), token, gameType)
	if err != nil {
		return nil, err
	}

	if res.IntVal("code") == 0 {
		return res.MapVal("data"), err
	} else {
		return nil, util.Err("DoEntryLandlords by error->%v, %v", err, util.S2Json(res))
	}
}



//开始斗地主新游戏
//选择指定条件开始斗地主新游戏，在游戏大厅选择指定场别进行新游戏，
//在游戏大厅快速开始游戏（根据携带金额快速匹配），在游戏房间内继续游戏，在游戏房间内若指定时间内不开始游戏则选择新房间快速开始游戏
//@url,需求登录，GET请求
//	~/usr/api/startNewLandlords		GET
//@arg,json对象中的参数以及query中参数
//	token		R	用户凭证
//	operate		R	操作 1.START_HALL,2.START_QUICKLY,3.START_CONTINUE
//	level	O	场别 1,2,3,4
/*
	样例	~/pub/api/startNewLandlords?token=xx&operate=1&categoryId=1
*/
//@ret,返回通用code/data
//	code	I	0：开始成功，1：参数错误，2：json body错误，3：开始游戏失败
//	landlords	O	斗地主新游戏
//	lid		I	游戏信息id
//	users		A	游戏用户id list
//	turnUser	S	轮到的用户
//	landlord_card	A	地主牌
//	last_pop_card	A	上家出的牌
//	a_cards		A	a玩家的手牌
//	b_cards		A	b玩家的手牌
//	c_cards		C	c玩家的手牌
//	stats		I	当前游戏状态
//	category	O	当前场别
//	landlord_user	S	地主玩家
//	multiple	O	倍数
//	multiple.double A	加倍的用户
//	multiple.bomb	A	炸弹
//	multiple.spring	I	春天
//	multiple.anti_spring	I	反春天
//	last		I	修改时间
//	time		I	创建时间
//	queue_time	I	排队开始时间
//	grab_time	I	抢地主开始时间
//	double_time	I	加倍开始时间
//	pop_card_time	I	出牌开始时间
//	over_game_time	I	游戏结束结果开始显示时间
/*	样例
*/
//@tag,游戏,斗地主
//@author,zhnagyq,2017-08-02
//@case,yule
func StartNewLandlords(hs *routing.HTTPSession) routing.HResult  {
	var operate,categoryId int
	err:=hs.ValidF(`
	operate,R|I,R:-1;
	categoryId,O|I,R:-1;
	`,&operate,&categoryId)
	if err!=nil {
		log.E("arg-err,%v", err)
		return hs.MsgResErr(1, "arg-err", err)
	}
	var addr = strings.Split(hs.R.Header.Get("X-Real-IP"), ":")[0]
	if len(addr) < 1 {
		addr = strings.Split(hs.R.RemoteAddr, ":")[0]
	}
	uid:=hs.StrVal("uid")
	if LOG_API {
		log.D(TAG_LANDLORDS+"uid(%v),operate(%v),categoryId(%v),addr(%v)",uid,operate,categoryId,addr)
	}
	if operate<1 {
		return hs.MsgResErr(1, "arg-err", err)
	}
	roomInfo:=&landlordsdb.LandlordInfo{}
	switch operate {
	case landlordsdb.START_HALL://从大厅选择分类进入
		roomInfo,err=landlordsdb.GetRoomRandom(uid,categoryId)
		break
	case landlordsdb.START_QUICKLY://从大厅快速进入房间
		break
	}
	if err!=nil {
		log.E("GetRoomRandom error(%v)",err)
		return hs.MsgResErr(2, "srv-err", err)
	}
	return hs.MsgRes(util.Map{"landlordInfo":roomInfo})
}

func DoStartNewLandlords(token string,operate,categoryId int) (util.Map, error) {
	res, err := util.HGet2("%v/usr/api/startNewLandlords?token=%v&operate=%v&categoryId=%v", SrvAddr(), token, operate,categoryId)
	if err != nil {
		return nil, err
	}

	if res.IntVal("code") == 0 {
		return res.MapVal("data"), err
	} else {
		return nil, util.Err("DoStartNewLandlords by error->%v, %v", err, util.S2Json(res))
	}
}

//获取斗地主信息
//在等待排队时需要轮询获取斗地主信息
//@url,需求登录，GET请求
//	~/usr/api/startNewLandlords		GET
//@arg,json对象中的参数以及query中参数
//	token		R	用户凭证
//	lid		R	斗地主id号
/*
	样例	~/pub/api/startNewLandlords?token=xx&lid=1
*/
//@ret,返回通用code/data
//	code	I	0：获取成功，1：参数错误，2：json body错误，3：获取游戏失败
//	landlords	O	斗地主新游戏
//	lid		I	游戏信息id
//	users		A	游戏用户id list
//	turnUser	S	轮到的用户
//	landlord_card	A	地主牌
//	last_pop_card	A	上家出的牌
//	a_cards		A	a玩家的手牌
//	b_cards		A	b玩家的手牌
//	c_cards		C	c玩家的手牌
//	stats		I	当前游戏状态
//	category	O	当前场别
//	landlord_user	S	地主玩家
//	multiple	O	倍数
//	multiple.double A	加倍的用户
//	multiple.bomb	A	炸弹
//	multiple.spring	I	春天
//	multiple.anti_spring	I	反春天
//	last		I	修改时间
//	time		I	创建时间
//	queue_time	I	排队开始时间
//	grab_time	I	抢地主开始时间
//	double_time	I	加倍开始时间
//	pop_card_time	I	出牌开始时间
//	over_game_time	I	游戏结束结果开始显示时间
/*	样例
*/
//@tag,游戏,斗地主
//@author,zhnagyq,2017-08-02
//@case,yule
func GetLandlords(hs *routing.HTTPSession) routing.HResult  {
	var lid string
	err:=hs.ValidF(`
	lid,R|S,L:0;
	`,&lid)
	if err!=nil {
		log.E("arg-err,%v", err)
		return hs.MsgResErr(1, "arg-err", err)
	}
	uid:=hs.StrVal("uid")
	if LOG_API {
		log.D(TAG_LANDLORDS+"uid(%v),lid(%v)",uid,lid)
	}
	if len(lid)<0 {
		return hs.MsgResErr(1, "参数错误", util.Err("lid(%v) error", lid))
	}
	landlordInfo,err:=landlordsdb.FindRoomV(uid,lid,"",0,0,nil)
	if err != nil {
		log.E("FindRoomV lids(%v)  err(%v)", lid, err)
		return hs.MsgResErr(2, "服务器错误", err)
	}
	if LOG_API {
		log.E("FindRoomV success info: %v",util.S2Json(landlordInfo))
	}
	return hs.MsgRes(util.Map{"landlordInfo":landlordInfo})
}

func DoGetLandlords(token,lid string) (util.Map, error) {
	res, err := util.HGet2("%v/usr/api/getLandlords?token=%v&lid=%v", SrvAddr(), token, lid)
	if err != nil {
		return nil, err
	}

	if res.IntVal("code") == 0 {
		return res.MapVal("data"), err
	} else {
		return nil, util.Err("DoGetLandlords by error->%v, %v", err, util.S2Json(res))
	}
}

//操作斗地主
//选择指定条件开始斗地主新游戏，在游戏大厅选择指定场别进行新游戏，
//在游戏大厅快速开始游戏（根据携带金额快速匹配），在游戏房间内继续游戏，在游戏房间内若指定时间内不开始游戏则选择新房间快速开始游戏
//@url,需求登录，GET请求
//	~/usr/api/OperateLandlords		GET
//@arg,json对象中的参数以及query中参数
//	token		R	用户凭证
//	lid		R	当前房间id
//	operate		R	操作 1.OP_GRAB 2.OP_DOUBLE,3.OP_POP_CARD,4.OP_GET_NOTE,5.OP_PASS_CARD
//	pop_cards	O	出牌
//	ming_type	O	明牌方式
//	double		O	加倍倍数
//	grab		O	抢地主分数
/*
	样例	~/pub/api/OperateLandlords?token=xx&lid=xxx&operate=1&pop_cards=1,2
*/
//@ret,返回通用code/data
//	code	I	0：进入成功，1：参数错误，2：json body错误，3：进入游戏失败
//	note_cards A	提示出牌
/*	样例
*/
//@tag,游戏,斗地主
//@author,zhnagyq,2017-08-02
//@case,yule
func OperateLandlords(hs *routing.HTTPSession) routing.HResult  {
	var lid string
	var operate,ming_type,double,grab int
	var pop_cards string
	err:=hs.ValidF(`
	lid,R|S,L:0;
	operate,R|I,R:-1;
	pop_cards,O|S,L:0;
	ming_type,O|I,R:-1;
	double,O|I,R:-1;
	grab,O|I,R:-1;
	`,&lid,&operate,&pop_cards,&ming_type,&double,&grab)
	if err!=nil {
		log.E("arg-err,%v", err)
		return hs.MsgResErr(1, "arg-err", err)
	}
	uid:=hs.StrVal("uid")
	if LOG_API {
		log.D(TAG_LANDLORDS+"uid(%v),lid(%v),operate(%v),pop_cards(%v),ming_type(%v),double(%v),grab(%v)",
			uid,lid,operate,pop_cards,ming_type,double,grab)
	}
	l_u:=&landlordsdb.LandlordInfo{Id:lid}
	if double>0 {
		l_u.Multiple=util.Map{landlordsdb.LD_DOUBLE_USERS:uid}
	}
	if grab>0 {
		l_u.Multiple=util.Map{landlordsdb.LD_GRAB_SCORE:grab}
	}
	landlordInfo,err:=landlordsdb.OperateLandlordInfo(uid,pop_cards,l_u,operate,ming_type)
	if err != nil {
		log.E("OperateLandlordInfo lids(%v),operate(%v),popCards(%v), mingType(%v),double(%v),grab(%v), err(%v)",
			lid,operate,pop_cards,ming_type,double,grab, err)
		return hs.MsgResErr(2, "服务器错误", err)
	}
	if LOG_API {
		log.D("OperateLandlordInfo success info: %v",util.S2Json(landlordInfo))
	}
	return hs.MsgRes(util.Map{"landlordInfo":landlordInfo})
}

func DoOperateLandlords(token,lid string,operate,ming_type,double,grab int ,pop_cards string) (util.Map, error) {
	res, err := util.HGet2("%v/usr/api/operateLandlords?token=%v&lid=%v&operate=%v&pop_cards=%v&ming_type=%v&double=%v&grab=%v",
		SrvAddr(), token, lid,operate,pop_cards,ming_type,double,grab)
	if err != nil {
		return nil, err
	}

	if res.IntVal("code") == 0 {
		return res.MapVal("data"), err
	} else {
		return nil, util.Err("DoOperateLandlords by error->%v, %v", err, util.S2Json(res))
	}
}