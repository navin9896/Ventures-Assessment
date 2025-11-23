package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"shopping-cart/database"
	"shopping-cart/handlers"
	"shopping-cart/middleware"
	"shopping-cart/models"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestShoppingCart(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shopping Cart Suite")
}

var _ = Describe("Shopping Cart API", func() {
	var router *gin.Engine
	var testUser models.User
	var testToken string

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.New()

		// Initialize test database (using in-memory SQLite for tests)
		database.InitTestDB()
		database.SeedData()

		// Setup routes
		router.POST("/users", handlers.CreateUser)
		router.POST("/users/login", handlers.Login)
		router.POST("/items", handlers.CreateItem)
		router.GET("/items", handlers.ListItems)
		router.POST("/carts", middleware.AuthMiddleware(), handlers.CreateCart)
		router.GET("/carts", middleware.AuthMiddleware(), handlers.ListCarts)
		router.POST("/orders", middleware.AuthMiddleware(), handlers.CreateOrder)
		router.GET("/orders", middleware.AuthMiddleware(), handlers.ListOrders)

		// Create a test user
		userReq := handlers.CreateUserRequest{
			Username: "testuser",
			Password: "testpass123",
		}
		userBody, _ := json.Marshal(userReq)
		req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(userBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Login to get token
		loginReq := handlers.LoginRequest{
			Username: "testuser",
			Password: "testpass123",
		}
		loginBody, _ := json.Marshal(loginReq)
		req = httptest.NewRequest("POST", "/users/login", bytes.NewBuffer(loginBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var loginResp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &loginResp)
		testToken = loginResp["token"].(string)
	})

	Describe("User Creation", func() {
		It("should create a new user", func() {
			userReq := handlers.CreateUserRequest{
				Username: "newuser",
				Password: "password123",
			}
			body, _ := json.Marshal(userReq)
			req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusCreated))
			var user models.User
			json.Unmarshal(w.Body.Bytes(), &user)
			Expect(user.Username).To(Equal("newuser"))
			Expect(user.Password).To(BeEmpty())
		})
	})

	Describe("Login", func() {
		It("should login with valid credentials", func() {
			loginReq := handlers.LoginRequest{
				Username: "testuser",
				Password: "testpass123",
			}
			body, _ := json.Marshal(loginReq)
			req := httptest.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(resp["token"]).ToNot(BeEmpty())
		})

		It("should reject invalid credentials", func() {
			loginReq := handlers.LoginRequest{
				Username: "testuser",
				Password: "wrongpass",
			}
			body, _ := json.Marshal(loginReq)
			req := httptest.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Describe("Cart Creation", func() {
		It("should create a cart with items", func() {
			cartReq := handlers.CreateCartRequest{
				ItemIDs: []uint{1, 2},
			}
			body, _ := json.Marshal(cartReq)
			req := httptest.NewRequest("POST", "/carts", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+testToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
		})
	})

	Describe("Order Creation", func() {
		It("should create an order from a cart", func() {
			// First create a cart
			cartReq := handlers.CreateCartRequest{
				ItemIDs: []uint{1},
			}
			body, _ := json.Marshal(cartReq)
			req := httptest.NewRequest("POST", "/carts", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+testToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var cart models.Cart
			json.Unmarshal(w.Body.Bytes(), &cart)

			// Create order
			orderReq := handlers.CreateOrderRequest{
				CartID: cart.ID,
			}
			orderBody, _ := json.Marshal(orderReq)
			req = httptest.NewRequest("POST", "/orders", bytes.NewBuffer(orderBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+testToken)
			w = httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusCreated))
			var order models.Order
			json.Unmarshal(w.Body.Bytes(), &order)
			Expect(order.CartID).To(Equal(cart.ID))
		})
	})
})

