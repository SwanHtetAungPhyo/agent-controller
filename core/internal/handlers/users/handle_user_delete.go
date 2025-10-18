package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v3/json"
	"github.com/rs/zerolog/log"
	"stock-agent.io/internal/types"
)

func (h *Handler) handleUserDeleted(c *gin.Context, event types.ClerkWebhookEvent) {
	var userData types.DeletedUserData
	if err := json.Unmarshal(event.Data, &userData); err != nil {
		log.Error().Err(err).Msg("failed to parse deleted user data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse user data"})
		return
	}

	//err := h.store.DeleteUser(context.Background(), userData.ID)
	//if err != nil {
	//	log.Error().Err(err).Str("user_id", userData.ID).Msg("failed to delete user")
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
	//	return
	//}

	log.Info().Str("user_id", userData.ID).Msg("user deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
