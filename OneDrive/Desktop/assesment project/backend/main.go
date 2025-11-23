package main

import (
	"log"

	"shopping-cart/database"
	"shopping-cart/handlers"
	"shopping-cart/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.InitDB()
	defer database.DB.Close()

	// Seed data
	database.SeedData()

	// Setup router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// User routes
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("", handlers.CreateUser)
		userRoutes.GET("", handlers.ListUsers)
		userRoutes.POST("/login", handlers.Login)
	}

	// Item routes
	itemRoutes := r.Group("/items")
	{
		itemRoutes.POST("", handlers.CreateItem)
		itemRoutes.GET("", handlers.ListItems)
	}

	// Cart routes (require authentication)
	cartRoutes := r.Group("/carts")
	cartRoutes.Use(middleware.AuthMiddleware())
	{
		cartRoutes.POST("", handlers.CreateCart)
		cartRoutes.GET("", handlers.ListCarts)
		cartRoutes.GET("/me", handlers.GetUserCart)
	}

	// Order routes (require authentication)
	orderRoutes := r.Group("/orders")
	orderRoutes.Use(middleware.AuthMiddleware())
	{
		orderRoutes.POST("", handlers.CreateOrder)
		orderRoutes.GET("", handlers.ListOrders)
	}

	// Start server
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

