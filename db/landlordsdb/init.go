package landlordsdb

import "yule/db"

var ShowLog bool = true

const (
	CN_LANDLORD_SET              = "landlord_set"
	CN_LANDLORD_INFO        = "landlord_info"
	CN_LANDLORD_RECORD           = "landlords_record"
)
const (
	SHOW_LOG_DB = true
	SHOW_LOG_DB_TEST = true
	TAG_LANDLORD_DB = "landlord_db"
)

func Remove()  {
	db.C(CN_LANDLORD_INFO).RemoveAll(nil)
}