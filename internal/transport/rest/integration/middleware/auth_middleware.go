package middleware

import (
	"net/http"
	"strings"

	"github.com/1stpay/1stpay/internal/domain/usecase"
	"github.com/gin-gonic/gin"
)

func APIKeyAuthMiddleware(apiKeyUsecase usecase.MerchantAPIKeyUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Api-key" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}
		apiKey := parts[1]
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			return
		}

		merchant, err := apiKeyUsecase.ValidateAPIKey(apiKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}

		c.Set("merchant", merchant)
		c.Next()
	}
}
