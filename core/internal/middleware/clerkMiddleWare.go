package middleware

import (
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
	"stock-agent.io/internal/types"
)

func (m *Manager) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userID := claims.Subject
		clerkUser, err := m.userClient.Get(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set(types.UserContextKey, clerkUser)
		c.Set(types.UserIDContextKey, userID)

		c.Next()
	}
}
