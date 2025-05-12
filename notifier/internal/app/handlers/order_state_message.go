package handlers

type OrderStateMessage struct {
	OrderId    int64  `json:"order_id"`
	UserId     int64  `json:"user_id"`
	FromStatus string `json:"from_status"`
	Status     string `json:"to_status"`
}
