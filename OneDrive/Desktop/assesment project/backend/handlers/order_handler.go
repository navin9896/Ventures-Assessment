package handlers

import (
	"net/http"
	"time"

	"shopping-cart/database"
	"shopping-cart/models"

	"github.com/gin-gonic/gin"
)

type CreateOrderRequest struct {
	CartID uint `json:"cart_id" binding:"required"`
}

func CreateOrder(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	currentUser := user.(*models.User)

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify cart belongs to user
	var cart models.Cart
	if err := database.DB.Where("id = ? AND user_id = ?", req.CartID, currentUser.ID).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found or does not belong to user"})
		return
	}

	// Check if cart is already checked out
	if cart.Status == "checked_out" || cart.Status == "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart already checked out"})
		return
	}

	// Create order
	order := models.Order{
		CartID:    req.CartID,
		UserID:    currentUser.ID,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Mark cart as checked out
	cart.Status = "checked_out"
	database.DB.Save(&cart)

	// Clear user's cart_id
	currentUser.CartID = nil
	database.DB.Save(currentUser)

	// Reload order with relationships
	database.DB.Where("id = ?", order.ID).Preload("Cart").Preload("User").First(&order)

	c.JSON(http.StatusCreated, order)
}

func ListOrders(c *gin.Context) {
	var orders []models.Order
	query := database.DB.Preload("Cart").Preload("Cart.CartItems").Preload("Cart.CartItems.Item").Preload("User")

	// Optional filtering by user_id
	userID := c.Query("user_id")
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

