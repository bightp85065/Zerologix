package mysql

import "time"

const (
	OrderStatusPending = iota + 1
	OrderStatusFailure
	OrderStatusCancelled
	OrderStatusCompleted

	OrderActionBuy = iota + 1
	OrderActionSell

	OrderPriceTypeMarket = iota + 1
	OrderPriceTypeLimit
)

// TODO ADD DESCRIPTION
type Order struct {
	ID        int       `gorm:"column:id;primary_key" json:"id"`                      //
	Product   int       `gorm:"column:product" json:"product"`                        //
	Quantity  int       `gorm:"column:quantity" json:"quantity"`                      //
	Creator   string    `gorm:"column:creator" json:"creator"`                        //
	Action    int       `gorm:"column:action" json:"action"`                          // 1:buy, 2:sell
	PriceType int       `gorm:"column:price_type" json:"price_type"`                  // 1:market, 2:limit
	Price     int       `gorm:"column:price" json:"price"`                            //
	Status    int       `gorm:"column:status" json:"status"`                          // 1:"pending", 2:"failure", 3:"cancelled", 4:"completed"
	Created   time.Time `gorm:"column:db_create_time;default:" json:"db_create_time"` //
	Updated   time.Time `gorm:"column:db_modify_time;default:" json:"db_modify_time"` //
}

// TODO ADD FIELD
var OrderFields = struct {
	ID string
}{
	ID: "id",
}

func (o Order) TableName() string {
	return "order"
}
