package menu

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateMenuItemHandler(c *gin.Context) {
	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	sqlDB := db.(*sql.DB)

	var menuItem struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Category string  `json:"category"`
	}

	if err := c.ShouldBindJSON(&menuItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	result, err := sqlDB.Exec("INSERT INTO menu_items (name, price, category) VALUES ($1, $2, $3)", menuItem.Name, menuItem.Price, menuItem.Category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create menu item"})
		return
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve last insert ID"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"menu_item": map[string]interface{}{
		"id":       lastInsertId,
		"name":     menuItem.Name,
		"price":    menuItem.Price,
		"category": menuItem.Category,
	}})
}

func ListMenuItemsHandler(c *gin.Context) {
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

	rows, err := sqlDB.Query("SELECT id, name, price, category FROM menu_items")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch menu items"})
		return
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var item struct {
			ID       int
			Name     string
			Price    float64
			Category string
		}

		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Category); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse menu item"})
			return
		}

		items = append(items, map[string]interface{}{
			"id":       item.ID,
			"name":     item.Name,
			"price":    item.Price,
			"category": item.Category,
		})
	}

	c.JSON(http.StatusOK, items)
}

func DeleteMenuItemHandler(c *gin.Context) {
	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	sqlDB := db.(*sql.DB)

	menuItemID := c.Param("id")

	_, err := sqlDB.Exec("DELETE FROM menu_items WHERE id = $1", menuItemID) //
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete menu item"}) //
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu item deleted successfully"}) //
}
