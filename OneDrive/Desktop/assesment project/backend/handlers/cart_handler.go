package handlers

import (
	"net/http"
	"time"

	"shopping-cart/database"
	"shopping-cart/models"

	"github.com/gin-gonic/gin"
)

type CreateCartRequest struct {
	ItemIDs []uint `json:"item_ids"`
}

func CreateCart(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	currentUser := user.(*models.User)

	var req CreateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already has an active cart
	var cart models.Cart
	if currentUser.CartID != nil {
		// User has a cart, use it
		if err := database.DB.Where("id = ?", *currentUser.CartID).Preload("CartItems").Preload("CartItems.Item").First(&cart).Error; err != nil {
			// Cart doesn't exist, create new one
			cart = models.Cart{
				UserID:    currentUser.ID,
				Name:      "My Cart",
				Status:    "active",
				CreatedAt: time.Now(),
			}
			if err := database.DB.Create(&cart).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
				return
			}
			currentUser.CartID = &cart.ID
			database.DB.Save(currentUser)
		}
	} else {
		// User has no cart, create new one
		cart = models.Cart{
			UserID:    currentUser.ID,
			Name:      "My Cart",
			Status:    "active",
			CreatedAt: time.Now(),
		}
		if err := database.DB.Create(&cart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
			return
		}
		currentUser.CartID = &cart.ID
		database.DB.Save(currentUser)
	}

	// Add items to cart
	for _, itemID := range req.ItemIDs {
		// Check if item exists
		var item models.Item
		if err := database.DB.First(&item, itemID).Error; err != nil {
			continue // Skip invalid items
		}

		// Check if item already in cart
		var existingCartItem models.CartItem
		if err := database.DB.Where("cart_id = ? AND item_id = ?", cart.ID, itemID).First(&existingCartItem).Error; err != nil {
			// Item not in cart, add it
			cartItem := models.CartItem{
				CartID: cart.ID,
				ItemID: itemID,
			}
			database.DB.Create(&cartItem)
		}
	}

	// Reload cart with items
	database.DB.Where("id = ?", cart.ID).Preload("CartItems").Preload("CartItems.Item").First(&cart)

	c.JSON(http.StatusOK, cart)
}

func ListCarts(c *gin.Context) {
	var carts []models.Cart
	query := database.DB.Preload("CartItems").Preload("CartItems.Item").Preload("User")

	// Optional filtering by user_id
	userID := c.Query("user_id")
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&carts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch carts"})
		return
	}

	c.JSON(http.StatusOK, carts)
}

func GetUserCart(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	currentUser := user.(*models.User)

	if currentUser.CartID == nil {
		c.JSON(http.StatusOK, gin.H{"message": "No active cart", "cart": nil})
		return
	}

	var cart models.Cart
	if err := database.DB.Where("id = ?", *currentUser.CartID).Preload("CartItems").Preload("CartItems.Item").First(&cart).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Cart not found", "cart": nil})
		return
	}

	c.JSON(http.StatusOK, cart)
}

