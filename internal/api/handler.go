package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suwandre/billing-api/internal"
)

type Handler struct {
	store internal.Store
}

func NewHandler(store internal.Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	v1.GET("/ping", h.Ping)
	h.RegisterCustomerRoutes(v1)
	h.RegisterPlanRoutes(v1)
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong!"})
}
