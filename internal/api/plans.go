package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suwandre/billing-api/internal/db/plans"
)

func (h *Handler) RegisterPlanRoutes(r *gin.RouterGroup) {
	r.POST("/plans", h.CreatePlan)
	r.POST("/plans/pricing", h.CreatePlanPricing)
	r.GET("/plans", h.ListPlans)
}

func (h *Handler) CreatePlan(c *gin.Context) {
	type planRequest struct {
		Name string `json:"name" binding:"required"`
	}

	body := planRequest{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error parsing request": err.Error()})
		return
	}

	plan, err := h.store.Plans().Create(c.Request.Context(), &plans.Plan{
		Name: body.Name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error creating plan": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

func (h *Handler) CreatePlanPricing(c *gin.Context) {
	type planPricingRequest struct {
		PlanID uuid.UUID `json:"plan_id" binding:"required"`
		Type   uint8     `json:"type"`
		Price  float64   `json:"price" binding:"required,gt=0"`
	}

	body := planPricingRequest{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error parsing request": err.Error()})
		return
	}

	planPricing, err := h.store.Plans().CreatePricing(c.Request.Context(), &plans.PlanPricing{
		PlanID: body.PlanID,
		Type:   plans.PricingType(body.Type),
		Price:  body.Price,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error creating plan pricing": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, planPricing)
}

func (h *Handler) ListPlans(c *gin.Context) {
	plans, err := h.store.Plans().List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error listing plans": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plans)
}
