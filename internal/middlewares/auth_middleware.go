package middlewares

import (
	"fmt"
	"iam_services_main_v1/pkg/auth/jwt"
	"iam_services_main_v1/pkg/logger"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a Gin middleware for JWT authentication with the default OpenID endpoint
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := os.Getenv("AZURE_TENANT_ID")
		userClaimField := os.Getenv("JWT_USER_CLAIM_FIELD")
		cacheDurationStr := os.Getenv("JWKS_CACHE_DURATION")
		cacheDuration, err := time.ParseDuration(cacheDurationStr)
		if err != nil {
			logger.LogError("Invalid JWKS_CACHE_DURATION format", "error", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Invalid JWKS_CACHE_DURATION format"})
			return
		}
		if tenantID == "" {
			logger.LogError("OpenID config URL not set")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "OpenID config URL not set"})
			return
		}

		// Create authenticator with default options but custom OpenID config URL
		options := jwt.DefaultOptions()
		options.OpenIDConfigURL = fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0/.well-known/openid-configuration", tenantID)
		options.UserClaimField = userClaimField
		options.CacheDuration = cacheDuration
		authenticator, err := jwt.NewAuthenticator(options)
		if err != nil {
			logger.LogError("Failed to initialize JWT authenticator", "error", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to initialize authenticator"})
			return
		}

		// Use the JWT middleware
		authenticator.GinMiddleware()(c)
	}
}
