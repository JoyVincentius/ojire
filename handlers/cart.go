package handlers

import (
	"net/http"

	"ojire/model"

	"github.com/gin-gonic/gin"
)

type addCartPayload struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,min=1"`
}

func GetCartHandler(c *gin.Context) {
	userID := c.GetInt64("userID")
	items, err := model.GetCartItems(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch cart"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func AddToCartHandler(c *gin.Context) {
	var payload addCartPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt64("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	if err := model.AddItemToCart(userID, payload.ProductID, payload.Quantity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "success add product",
	})
}

func RemoveFromCartHandler(c *gin.Context) {
	var payload struct {
		ProductID int64 `json:"product_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetInt64("userID")
	if err := model.RemoveItemFromCart(userID, payload.ProductID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success remove product",
	})
}

func ClearCartHandler(c *gin.Context) {
	userID := c.GetInt64("userID")
	if err := model.ClearCart(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot clear cart"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success clear cart",
	})
}
