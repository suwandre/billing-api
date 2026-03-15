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
	r.GET("/ping", h.Ping)
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong!"})
}
