package users

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v3/json"
	"github.com/rs/zerolog/log"
	"stock-agent.io/internal/types"
)

func (h *Handler) handleClerkWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read webhook body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	var event types.ClerkWebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Error().Err(err).Msg("failed to parse webhook event")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse event"})
		return
	}

	switch event.Type {
	case "user.created":
		h.handleUserCreated(c, event)
	case "user.updated":
		h.handleUserUpdated(c, event)
	case "user.deleted":
		h.handleUserDeleted(c, event)
	default:
		log.Info().Str("type", event.Type).Msg("unhandled webhook event type")
		c.JSON(http.StatusOK, gin.H{"message": "event type not handled"})
	}
}
