package landlordsdb

import (
	"github.com/Centny/gwf/util"
)

type LandlordInfo struct {
	Id            string   `bson:"_id" json:"lid"`                                 //the landlords id
	RoomId		int    `bson:"room_id" json:"room_id,omitempty"`
	Players		[]util.Map `bson:"players" json:"players,omitempty"`                   //the landlords players list
	//uid:用户id,status:Online,Offline,Collocation托管,cards：牌,popCards：出牌,ming：明牌方式,host：ip地址，popTimes出牌次数,grab_score抢地主分数
	TurnUser      string   `bson:"turn_user" json:"turn_user,omitempty"`           //the landlords user turn
	OperateNum      int      `bson:"operate_num" json:"operate_num,omitempty"` //当前回合的操作人数
	LandlordCards string   `bson:"landlord_cards" json:"landlord_cards,omitempty"` //the landlord cards attribute list
	LastPopCards  string   `bson:"last_pop_cards" json:"last_pop_cards,omitempty"` //the last pop cards attribute list
	NotePopCards  string   `bson:"note_pop_cards" json:"note_pop_cards,omitempty"` //提示出牌
	Status        int      `bson:"status" json:"status,omitempty"`                 //the landlords status//等待，抢地主，加倍，出牌，结束，无效
	Category      util.Map `bson:"category" json:"category,omitempty"`             //the landlords category 场别
	LandlordUser  string   `bson:"landlord_user" json:"landlord_user,omitempty"`   //the  landlords user
	Multiple      util.Map `bson:"multiple" json:"multiple,omitempty"`             //the landlords 倍数
	Last          int64    `bson:"last" json:"last,omitempty"`                     //the last updated time
	Time          int64    `bson:"time" json:"time,omitempty"`                     //the create time
	QueueTime     int64    `bson:"queue_time" json:"queue_time,omitempty"`                  //排队开始时间
	GrabTime      int64    `bson:"grab_time" json:"grab_time,omitempty"`                   //抢地主开始时间
	DoubleTime    int64    `bson:"double_time" json:"double_time,omitempty"`                 //加倍开始时间
	PopCardTime   int64    `bson:"pop_card_time" json:"pop_card_time,omitempty"`               //出牌开始时间
	OverGameTime  int64    `bson:"over_game_time" json:"over_game_time,omitempty"`              //游戏结束结果开始显示时间
	MingTime  int64    `bson:"ming_time" json:"ming_time,omitempty"`              //明牌选择开始时间
	Size          int      `bson:"size" json:"size,omitempty"`                     //the landlords size
	Record 	*LandlordRecord	`bson:"-" json:"record,omitempty"`                    //记录
}

//当前用户的游戏状态
const (
	UGS_OFF       = 2 //不在游戏中
	UGS_LANDLORDS = 1  //斗地主中
)

//游戏类型
const (
	UGT_LANDLORDS           = 1 //斗地主
	UGT_FRIED_GOLDEN_FLOWER = 2 //炸金花
	UGT_PUSH_BOBBIN         = 3 //推筒子
	UGT_BULLFIGHT           = 4 //斗牛
)

//斗地主的状态
const (
	//用户状态
	LUS_ONLINE = 1
	LUS_OFFLINE = 2
	LUS_COLLOCATION = 3//托管


	LS_INVALID     = -1
	LS_QUEUE       = 1
	LS_MING = 2
	LS_GRABBING    = 3
	LS_DOUBLING    = 4
	LS_FIGHTING    = 5
	LS_SHOW_RESULT = 6
)

//倍数,底分
const (
	LW_LANDLORD = "landlord"
	LW_FARMER = "farmer"

	LD_GRAB_SCORE = "grab_score"
	LD_DOUBLE_USERS = "double_users"
	LD_BOMBS = "bombs"
	LD_SPRING = "spring"
	LD_ANTI_SPRING = "anti_spring"
	LD_MING = "ming"
)

//type  Multiple struct {
//	Double string[]
//	Bomb  string[]
//	Spring int//农民不出牌 默认0
//	AntiSpring int//地主只出了一次就不再出 默认0
//	Ming	string//明牌用户
//}

//进入游戏后的结果
type EntryGameResult struct {
	Result        int         `bson:"-" json:"result,omitempty"`         //结果：继续游戏，进入斗地主游戏大厅
	LandlordsInfo LandlordInfo   `bson:"-" json:"landlords_info,omitempty"` //继续游戏的游戏信息
	LandlordHall  LandlordSet `bson:"-" json:"landlord_hall,omitempty"`  //斗地主大厅
}

const (
	EGR_LANDLORD_HALL = 1 //进入斗地主游戏大厅
	EGR_CONTINUE      = 2 //继续游戏
	//开始新游戏
	START_HALL     = 1 //从游戏大厅选择指定场别开始新游戏
	START_QUICKLY  = 2 //从游戏大厅快速进入游戏
	START_CONTINUE = 3 //在游戏结束后继续开始新游戏
	//游戏操作
	OP_GRAB        = 1 //抢地主
	OP_DOUBLE      = 2 //加倍
	OP_POP_CARD = 3 //出牌
	OP_GET_NOTE    = 4 //获取提示
	OP_PASS_CARD   = 5 //不出
	OP_MING = 6 //明牌
	OP_DEAL = 7	//发牌
	OP_CONTINUE = 8//继续
	OP_RETURN_HALL  = 9//返回大厅

	//进入操作
	OP_ENTRY_MING = 8
	//不出牌
	LC_PASS = "PASS"

)

//记录
type LandlordRecord struct {
	Id                 string   `bson:"_id" json:"id"`                                              //the landlords record id
	Users              []string `bson:"users" json:"users,omitempty"`                               //the landlords users list
	LandlordUser       string   `bson:"landlord_user" json:"landlord_user,omitempty"`               //the  landlords user
	WinRole            string   `bson:"win_role" json:"win_role"`                                   //赢家身份
	Category           util.Map `bson:"category" json:"category,omitempty"`                         //the landlords category 场别
	LandlordMultiple   int64      `bson:"landlord_multiple" json:"landlord_multiple,omitempty"`       //the 地主倍数
	DoubleMultiple     int64      `bson:"double_multiple" json:"double_multiple,omitempty"`           //加倍倍数
	BombMultiple       int64      `bson:"bomb_multiple" json:"bomb_multiple,omitempty"`               //炸弹倍数
	SpringMultiple     int64      `bson:"spring_multiple" json:"spring_multiple,omitempty"`           //春天倍数
	AntiSpringMultiple int64      `bson:"anti_spring_multiple" json:"anti_spring_multiple,omitempty"` //反春天倍数
	SumMultiple        int64      `bson:"sum_multiple" json:"sum_multiple,omitempty"`                 //总倍数
	Time               int64    `bson:"time" json:"time,omitempty"`                                 //the create time
}

type LandlordSet struct {
	Id	     string	`bson:"_id" json:"-"`
	SoundEffect  int      `bson:"sound_effect" json:"sound_effect,omitempty"`             //音效
	Music        int      `bson:"music" json:"music,omitempty"`                    //音乐
	Categories   []util.Map `bson:"categories" json:"categories,omitempty"`               //场别分类 title level endScore minGold
	QueueTime    int64      `bson:"queue_time" json:"queue_time,omitempty"`     //排队时间
	MingTime    int64      `bson:"ming_time" json:"ming_time,omitempty"`     //明牌时间
	GrabTime     int64      `bson:"grab_time" json:"grab_time,omitempty"`      //抢地主时间
	DoubleTime   int64      `bson:"double_time" json:"double_time,omitempty"`    //加倍时间
	PopCardTime  int64      `bson:"pop_card_time" json:"pop_card_time,omitempty"`  //出牌时间
	OverGameTime int64      `bson:"over_game_time" json:"over_game_time,omitempty"` //游戏结束结果显示时间
}

/**
常量
*/
const (
	SOUND_OFF = false
	SOUND_ON  = true
	//斗地主分类
	LC_TITLE = "title"
	LC_LEVEL = "level"
	LC_SCORE_END = "score_end"
	LC_GOLD_MIN = "gold_min"

	LC_PRIMER_TITLE       = "入门级"
	LC_PRIMER_LEVEL       = "初级场"
	LC_PRIMER_SCORE_END   = 0.01
	LC_PRIMER_GOLD_MIN    = 10
	LC_SUPERIOR_TITLE     = "高手级"
	LC_SUPERIOR_LEVEL     = "中级场"
	LC_SUPERIOR_SCORE_END = 0.5
	LC_SUPERIOR_GOLD_MIN  = 14
	LC_MASTER_TITLE       = "大师级"
	LC_MASTER_LEVEL       = "高级场"
	LC_MASTER_SCORE_END   = 1
	LC_MASTER_GOLD_MIN    = 30
	LC_TOP_TITLE          = "巅峰级"
	LC_TOP_LEVEL          = "土豪场"
	LC_TOP_SCORE_END      = 3
	LC_TOP_GOLD_MIN       = 1000
	//时间限制
	LT_QUEUE              = 60*TIME_INT64
	LT_GRAB               = 10*TIME_INT64
	LT_MING			= 5*TIME_INT64
	LT_DOUBLE             = 10*TIME_INT64
	LT_POP_CARD           = 30*TIME_INT64
	LT_GAME_OVER_SHOW	= 10*TIME_INT64

	TIME_INT64 = 100

	//明牌类型倍数
	LM_TD_BEFORE_DEAL = 5//发牌前
	LM_TD_DEALING_FRONT = 4//发牌前半部分
	LM_TD_DEALING_REAR = 3//发牌后半部分
	LM_TD_BEFORE_POP = 2//出牌前
	LM_TD_NONE = 0//不明牌

	LC_PRIMER_ID   = 1
	LC_SUPERIOR_ID = 2
	LC_MASTER_ID   = 3
	LC_TOP_ID      = 4
)
