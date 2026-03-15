package routes

import (
	"net/http"

	"school-vault/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to School Nutrition Vault API!",
		})
	})

	// Endpoint Barang (Items) - Sekarang FULL CRUD!
	r.POST("/items", controllers.CreateItem)       // Create
	r.GET("/items", controllers.GetItems)          // Read (All)
	r.PUT("/items/:id", controllers.UpdateItem)    // Update
	r.DELETE("/items/:id", controllers.DeleteItem) // Delete

	r.GET("/items/critical", controllers.GetCriticalStocks)

	// Endpoint Transaksi (Mutasi Stok)
	r.POST("/transactions", controllers.CreateTransaction)

	return r
}
