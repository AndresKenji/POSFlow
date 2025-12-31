package routes

import (
	"POSFlowBackend/internal/infrastructure/http/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all HTTP routes
func RegisterRoutes(
	router *gin.Engine,
	productHandler *handlers.ProductHandler,
	orderHandler *handlers.OrderHandler,
) {
	// Health check endpoint
	router.GET("/health", healthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Product routes
		registerProductRoutes(v1, productHandler)

		// Order routes
		registerOrderRoutes(v1, orderHandler)
	}
}

// healthCheck handles health check requests
func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "healthy",
		"service": "POSFlow Backend",
		"version": "1.0.0",
	})
}

// registerProductRoutes registers all product-related routes
func registerProductRoutes(rg *gin.RouterGroup, handler *handlers.ProductHandler) {
	products := rg.Group("/products")
	{
		// Special route: must be before /:id to avoid conflict
		products.GET("/low-stock", handler.GetLowStockProducts)

		// Standard CRUD operations
		products.POST("", handler.CreateProduct)
		products.GET("", handler.ListProducts)
		products.GET("/:id", handler.GetProduct)
		products.PUT("/:id", handler.UpdateProduct)
		products.DELETE("/:id", handler.DeleteProduct)

		// Stock management
		products.POST("/:id/stock", handler.UpdateStock)
	}
}

// registerOrderRoutes registers all order-related routes
func registerOrderRoutes(rg *gin.RouterGroup, handler *handlers.OrderHandler) {
	orders := rg.Group("/orders")
	{
		// Special route: must be before /:id to avoid conflict
		orders.GET("/pending", handler.GetPendingOrders)

		// Standard CRUD operations
		orders.POST("", handler.CreateOrder)
		orders.GET("", handler.ListOrders)
		orders.GET("/:id", handler.GetOrder)

		// Order status management
		orders.PATCH("/:id/status", handler.UpdateOrderStatus)
	}
}
