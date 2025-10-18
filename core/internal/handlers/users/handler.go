package users

import (
	"github.com/gin-gonic/gin"
	db "stock-agent.io/db/sqlc"
)

type Handler struct {
	router     *gin.Engine
	store      db.Store
	webhookKey string
}

func NewHandler(router *gin.Engine, store db.Store, webhookKey string) *Handler {
	return &Handler{
		router:     router,
		store:      store,
		webhookKey: webhookKey,
	}
}

func (h *Handler) RegisterRoutes() {
	h.router.POST("/webhooks/clerk", h.handleClerkWebhook)

}
