package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suwandre/billing-api/internal/db/plans"
)

func (h *Handler) RegisterPlanRoutes(r *gin.RouterGroup) {
	r.POST("/subscriptions", h.CreateSubscription)
	r.POST("/subscriptions/pricing", h.CreateSubscriptionPricing)
	r.GET("/subscriptions", h.List)
}

func (h *Handler) CreateSubscription(c *gin.Context) {
	type subscriptionRequest struct {
		Name string `json:"name" binding:"required"`
	}

	body := subscriptionRequest{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error parsing request": err.Error()})
		return
	}

	subscription, err := h.store.Subscriptions().CreateSubscription(c.Request.Context(), &plans.Subscription{
		Name: body.Name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error creating subscription": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

func (h *Handler) CreateSubscriptionPricing(c *gin.Context) {
	type subscriptionPricingRequest struct {
		SubscriptionID uuid.UUID `json:"subscription_id" binding:"required"`
		Type           uint8     `json:"type"`
		Price          float64   `json:"price" binding:"required" gt=0"`
	}

	body := subscriptionPricingRequest{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error parsing request": err.Error()})
		return
	}

	subscriptionPricing, err := h.store.Subscriptions().CreateSubscriptionPricing(c.Request.Context(), &plans.SubscriptionPricing{
		SubscriptionID: body.SubscriptionID,
		Type:           plans.PricingType(body.Type),
		Price:          body.Price,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error creating subscription pricing": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subscriptionPricing)
}

func (h *Handler) List(c *gin.Context) {
	subscriptions, err := h.store.Subscriptions().List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error listing subscriptions": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}
