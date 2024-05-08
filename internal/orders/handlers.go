package orders

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateOrderHandler(c *gin.Context) {
	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	sqlDB, ok := db.(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid database connection"})
		return
	}

	var order struct {
		UserID     int `json:"user_id"`
		MenuItemID int `json:"menu_item_id"`
		Quantity   int `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if order.UserID == 0 || order.MenuItemID == 0 || order.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order data"})
		return
	}

	_, err := sqlDB.Exec("INSERT INTO orders (user_id, menu_item_id, quantity) VALUES ($1, $2, $3)", order.UserID, order.MenuItemID, order.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Order created"})
}
func ListOrdersHandler(c *gin.Context) {
	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	sqlDB := db.(*sql.DB)

	query := `
        SELECT 
            o.id AS order_id, 
            u.username AS user_name, 
            m.name AS menu_item_name, 
            m.price AS menu_item_price, 
            o.quantity 
        FROM 
            orders o 
        INNER JOIN 
            users u ON o.user_id = u.id 
        INNER JOIN 
            menu_items m ON o.menu_item_id = m.id
    `

	rows, err := sqlDB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	defer rows.Close()

	orders := []gin.H{}
	for rows.Next() {
		var order struct {
			OrderID       int
			UserName      string
			MenuItemName  string
			MenuItemPrice float64
			Quantity      int
		}

		if err := rows.Scan(&order.OrderID, &order.UserName, &order.MenuItemName, &order.MenuItemPrice, &order.Quantity); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse order"})
			return
		}

		orders = append(orders, gin.H{
			"order_id":        order.OrderID,
			"user_name":       order.UserName,
			"menu_item_name":  order.MenuItemName,
			"menu_item_price": order.MenuItemPrice,
			"quantity":        order.Quantity,
		})
	}

	c.JSON(http.StatusOK, orders)
}

func UpdateOrderHandler(c *gin.Context) {
	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	sqlDB := db.(*sql.DB)

	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := sqlDB.Exec("UPDATE orders SET menu_item_id = ?, quantity = ?, status = ? WHERE id = ?", order.MenuItemID, order.Quantity, order.Status, order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

func DeleteOrderHandler(c *gin.Context) {
	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	sqlDB := db.(*sql.DB)

	id := c.Param("id")

	_, err := sqlDB.Exec("DELETE FROM orders WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
}
