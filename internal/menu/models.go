package menu

// MenuItem represents a menu item in the system
type MenuItem struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"` // Optional, for categorizing items (e.g., beverages, meals)
}
