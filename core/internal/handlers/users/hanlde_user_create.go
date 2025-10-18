package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v3/json"
	"github.com/rs/zerolog/log"
	"stock-agent.io/internal/types"
)

func (h *Handler) handleUserCreated(c *gin.Context, event types.ClerkWebhookEvent) {
	var userData types.UserData
	if err := json.Unmarshal(event.Data, &userData); err != nil {
		log.Error().Err(err).Msg("failed to parse user data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse user data"})
		return
	}

	//email := ""
	//if len(userData.EmailAddresses) > 0 {
	//	email = userData.EmailAddresses[0].EmailAddress
	//}
	//
	//params := db.CreateUserParams{
	//	ID:        userData.ID,
	//	Email:     email,
	//	FirstName: userData.FirstName,
	//	LastName:  userData.LastName,
	//	ImageUrl:  userData.ImageURL,
	//	Username:  userData.Username,
	//	CreatedAt: time.Unix(userData.CreatedAt/1000, 0),
	//	UpdatedAt: time.Unix(userData.UpdatedAt/1000, 0),
	//}
	//
	//_, err := h.store.CreateUser(context.Background(), params)
	//if err != nil {
	//	log.Error().Err(err).Str("user_id", userData.ID).Msg("failed to create user")
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
	//	return
	//}

	log.Info().Str("user_id", userData.ID).Msg("user created successfully")
	c.JSON(http.StatusOK, gin.H{"message": "user created"})
}
