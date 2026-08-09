package middleware

import (
	"ams-backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	KeyUserID = "userID"
	KeyRole   = "role"
)

// Auth validates the JWT and sets userID + role in context.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false, "message": "authorization header required",
			})
			return
		}

		claims, err := utils.ParseJWT(strings.TrimPrefix(auth, "Bearer "))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false, "message": "invalid or expired token",
			})
			return
		}

		c.Set(KeyUserID, claims.UserID)
		c.Set(KeyRole, claims.Role)
		c.Next()
	}
}

// RequireRole aborts with 403 if the user does not have the specified role.
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString(KeyRole) != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false, "message": "insufficient permissions",
			})
			return
		}
		c.Next()
	}
}

func GetUserID(c *gin.Context) int64 {
	v, _ := c.Get(KeyUserID)
	id, _ := v.(int64)
	return id
}

func GetRole(c *gin.Context) string {
	return c.GetString(KeyRole)
}
