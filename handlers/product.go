package handlers

import (
	"net/http"

	"ojire/model"

	"github.com/gin-gonic/gin"
)

func ListProductsHandler(c *gin.Context) {
	products, err := model.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch products"})
		return
	}
	c.JSON(http.StatusOK, products)
}
