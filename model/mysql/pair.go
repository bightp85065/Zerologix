package mysql

import "time"

// TODO ADD DESCRIPTION
type Pair struct {
	ID          int       `gorm:"column:id;primary_key" json:"id"`                      //
	BuyOrderID  int       `gorm:"column:buy_order_id" json:"buy_order_id"`              //
	SellOrderID int       `gorm:"column:sell_order_id" json:"sell_order_id"`            //
	Price       int       `gorm:"column:price" json:"price"`                            //
	Status      int       `gorm:"column:status" json:"status"`                          // 1:"pending", 2:"failure", 3:"cancelled", 4:"completed"
	MatchTime   time.Time `gorm:"column:match_time" json:"match_time"`                  //
	Created     time.Time `gorm:"column:db_create_time;default:" json:"db_create_time"` //
	Updated     time.Time `gorm:"column:db_modify_time;default:" json:"db_modify_time"` //
}

// TODO ADD FIELD
var PairFields = struct {
	ID string
}{
	ID: "id",
}
