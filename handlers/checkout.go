package handlers

import (
	"net/http"

	"ojire/model"

	"github.com/gin-gonic/gin"
)

func CheckoutHandler(c *gin.Context) {
	userID := c.GetInt64("userID")

	cartItems, err := model.GetCartItems(userID)
	if err != nil || len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart empty or cannot fetch"})
		return
	}

	var total int64
	for _, it := range cartItems {
		total += it.TotalCents
	}

	tx, err := model.BeginTx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot start transaction"})
		return
	}
	defer tx.Rollback()

	var orderID int64
	res, err := tx.Exec("INSERT INTO orders (user_id, total_cents) VALUES (?, ?)", userID, total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "order creation failed"})
		return
	}
	orderID, _ = res.LastInsertId()

	for _, it := range cartItems {
		_, err = tx.Exec(`
            INSERT INTO order_items (order_id, product_id, quantity, price_cents)
            VALUES (?, ?, ?, ?)
        `, orderID, it.ProductID, it.Quantity, it.PriceCents)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "order items insert failed"})
			return
		}
		_, err = tx.Exec(`
            UPDATE products SET stock = stock - ?
            WHERE id = ? AND stock >= ?
        `, it.Quantity, it.ProductID, it.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "stock update failed"})
			return
		}
	}

	if err = model.ClearCartTx(tx, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot clear cart"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "checkout commit failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order_id":    orderID,
		"total_cents": total,
		"message":     "checkout successful",
	})
}
