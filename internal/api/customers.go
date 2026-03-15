package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suwandre/billing-api/internal/db/customers"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) RegisterCustomerRoutes(r *gin.RouterGroup) {
	r.POST("/create", h.CreateCustomer)
	r.GET("/getByEmail", h.GetCustomerByEmail)
}

func (h *Handler) CreateCustomer(c *gin.Context) {
	type createCustomerRequest struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	body := createCustomerRequest{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error parsing request": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error hashing password": err.Error()})
		return
	}

	customer, err := h.store.Customers().Create(c.Request.Context(), &customers.Customer{
		Email:        body.Email,
		Username:     body.Username,
		PasswordHash: string(hash),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error creating customer": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

func (h *Handler) GetCustomerByEmail(c *gin.Context) {
	email, ok := c.GetQuery("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	customer, err := h.store.Customers().GetByEmail(c.Request.Context(), email)
	if err != nil {
		if err.Error() == "customer not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error getting customer": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}
