package middleware

import (
	"net/http"
	"strings"

	"github.com/Doris-Mwito5/ginja-ai/internal/jwt"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	authorizationBearer = "Bearer"
	authPayloadKey      = "auth_payload"
)

// AuthMiddleware validates the JWT Bearer token on every request.
func AuthMiddleware(jwtMaker jwt.JWTToken) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Get the Authorization header
		header := c.GetHeader(authorizationHeader)
		if len(header) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "MISSING_TOKEN",
				"message": "authorization header is required",
			})
			return
		}

		// 2. Expect "Bearer <token>"
		parts := strings.Fields(header)
		if len(parts) != 2 || !strings.EqualFold(parts[0], authorizationBearer) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "INVALID_TOKEN_FORMAT",
				"message": "authorization header must be in format: Bearer <token>",
			})
			return
		}

		// 3. Verify the token using your JWTToken implementation
		tokenStr := parts[1]
		payload, err := jwtMaker.VerifyToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "INVALID_TOKEN",
				"message": err.Error(),
			})
			return
		}

		// 4. Stash the payload on context for downstream handlers
		c.Set(authPayloadKey, payload)
		c.Next()
	}
}

// GetAuthPayload retrieves the JWT payload from the Gin context.
// Call this in any handler that needs the logged-in user's info.
func GetAuthPayload(c *gin.Context) *jwt.Payload {
	payload, _ := c.Get(authPayloadKey)
	return payload.(*jwt.Payload)
}