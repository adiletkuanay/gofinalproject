package orders

type Order struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	MenuItemID int    `json:"menu_item_id"`
	Quantity   int    `json:"quantity"`
	Status     string `json:"status"`
}
