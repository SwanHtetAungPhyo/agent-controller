package users

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	db "stock-agent.io/db/sqlc"
	"stock-agent.io/internal/events"
	"stock-agent.io/internal/types"
)

type Handler struct {
	store          db.Store
	webhookKey     string
	eventPublisher *events.Publisher
}

func NewHandler(store db.Store, webhookKey string, eventPublisher *events.Publisher) *Handler {
	return &Handler{
		store:          store,
		webhookKey:     webhookKey,
		eventPublisher: eventPublisher,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/webhooks/clerk", h.handleClerkWebhook)

	// Test endpoints
	api := router.Group("/api/v1")
	{
		api.POST("/test-user-event", h.handleTestUserEvent)
	}
}

func (h *Handler) handleClerkWebhook(c *gin.Context) {
	var webhookEvent types.ClerkWebhookEvent
	if err := c.ShouldBindJSON(&webhookEvent); err != nil {
		log.Error().Err(err).Msg("Failed to bind webhook event")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	log.Info().
		Str("event_type", webhookEvent.Type).
		Int64("timestamp", webhookEvent.Timestamp).
		Msg("Received Clerk webhook")

	switch webhookEvent.Type {
	case "user.created":
		h.handleUserCreated(c, webhookEvent.Data)
	case "user.updated":
		h.handleUserUpdated(c, webhookEvent.Data)
	case "user.deleted":
		h.handleUserDeleted(c, webhookEvent.Data)
	default:
		log.Warn().Str("event_type", webhookEvent.Type).Msg("Unhandled webhook event type")
		c.JSON(http.StatusOK, gin.H{"message": "Event type not handled"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}

func (h *Handler) handleUserCreated(c *gin.Context, data json.RawMessage) {
	var userData types.UserData
	if err := json.Unmarshal(data, &userData); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal user data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		return
	}

	email := ""
	if len(userData.EmailAddresses) > 0 {
		email = userData.EmailAddresses[0].EmailAddress
	}

	log.Info().
		Str("user_id", userData.ID).
		Str("email", email).
		Str("first_name", userData.FirstName).
		Str("last_name", userData.LastName).
		Msg("Processing user created event")

	if err := h.eventPublisher.PublishUserCreated(
		userData.ID,
		email,
		userData.FirstName,
		userData.LastName,
	); err != nil {
		log.Error().Err(err).Msg("Failed to publish user created event")
	}
}

func (h *Handler) handleUserUpdated(c *gin.Context, data json.RawMessage) {
	var userData types.UserData
	if err := json.Unmarshal(data, &userData); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal user data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		return
	}

	email := ""
	if len(userData.EmailAddresses) > 0 {
		email = userData.EmailAddresses[0].EmailAddress
	}

	log.Info().
		Str("user_id", userData.ID).
		Str("email", email).
		Msg("Processing user updated event")

	if err := h.eventPublisher.PublishUserUpdated(
		userData.ID,
		email,
		userData.FirstName,
		userData.LastName,
	); err != nil {
		log.Error().Err(err).Msg("Failed to publish user updated event")
	}
}

func (h *Handler) handleUserDeleted(c *gin.Context, data json.RawMessage) {
	var deletedData types.DeletedUserData
	if err := json.Unmarshal(data, &deletedData); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal deleted user data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deleted user data"})
		return
	}

	log.Info().
		Str("user_id", deletedData.ID).
		Msg("Processing user deleted event")

	if err := h.eventPublisher.PublishUserDeleted(deletedData.ID); err != nil {
		log.Error().Err(err).Msg("Failed to publish user deleted event")
	}
}

func (h *Handler) handleTestUserEvent(c *gin.Context) {
	var request struct {
		Email     string `json:"email" binding:"required,email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		EventType string `json:"event_type"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if request.FirstName == "" {
		request.FirstName = "Test"
	}
	if request.LastName == "" {
		request.LastName = "User"
	}
	if request.EventType == "" {
		request.EventType = "user.created"
	}

	userID := "test-user-" + request.Email

	log.Info().
		Str("user_id", userID).
		Str("email", request.Email).
		Str("event_type", request.EventType).
		Msg("Processing test user event")

	var err error
	switch request.EventType {
	case "user.created":
		err = h.eventPublisher.PublishUserCreated(
			userID,
			request.Email,
			request.FirstName,
			request.LastName,
		)
	case "user.updated":
		err = h.eventPublisher.PublishUserUpdated(
			userID,
			request.Email,
			request.FirstName,
			request.LastName,
		)
	case "user.deleted":
		err = h.eventPublisher.PublishUserDeleted(userID)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event type. Use: user.created, user.updated, or user.deleted"})
		return
	}

	if err != nil {
		log.Error().Err(err).Msg("Failed to publish test user event")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to publish event",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Test user event published successfully",
		"user_id":    userID,
		"email":      request.Email,
		"event_type": request.EventType,
	})
}
