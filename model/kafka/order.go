package kafka

import "time"

// TODO ADD DESCRIPTION
type OrderRawMessage struct {
	ID        int       `json:"id"`
	Product   int       `json:"product"`
	Quantity  int       `json:"quantity"`
	Creator   string    `json:"creator"`
	Action    int       `json:"action"`     // 1:buy, 2:sell
	PriceType int       `json:"price_type"` // 1:market, 2:limit
	Price     int       `json:"price"`
	Status    int       `json:"status"` // 1:"pending", 2:"failure", 3:"cancelled", 4:"completed"
	Created   time.Time `json:"db_create_time"`
	Updated   time.Time `json:"db_modify_time"`
}
