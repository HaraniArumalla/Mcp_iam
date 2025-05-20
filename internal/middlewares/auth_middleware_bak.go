package middlewares

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"iam_services_main_v1/pkg/logger"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// JWTClaims structure
type JWTClaims struct {
	UserID string `json:"user_id"`
	Exp    int64  `json:"exp"`
}

// DecodeToken decodes a JWT token without validation (similar to jwt.io)
func DecodeToken(tokenString string) (*JWTClaims, error) {
	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Split the token into parts
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		logger.LogError("Invalid token format")
		return nil, fmt.Errorf("invalid token format")
	}

	// Decode the claims part (second part)
	claimsPart, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		logger.LogError("Failed to decode claims", "error", err)
		return nil, fmt.Errorf("failed to decode claims: %v", err)
	}

	// Parse the claims
	var claims JWTClaims
	if err := json.Unmarshal(claimsPart, &claims); err != nil {
		logger.LogError("Failed to parse claims", "error", err)
		return nil, fmt.Errorf("failed to parse claims: %v", err)
	}

	// Check if token is expired
	if claims.Exp > 0 && time.Unix(claims.Exp, 0).Before(time.Now()) {
		logger.LogError("Token is expired")
		return nil, fmt.Errorf("token is expired")
	}

	return &claims, nil
}

// AuthMiddleware validates JWT in the Authorization header
func AuthMiddleware1() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		// authHeader := c.GetHeader("Authorization")
		// if authHeader == "" {
		// 	logger.LogError("Authorization header missing")
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		// 	c.Abort()
		// 	return
		// }

		// // Check if the token is prefixed with "Bearer"
		// tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// if tokenString == authHeader {
		// 	logger.LogError("Bearer token missing")
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token missing"})
		// 	c.Abort()
		// 	return
		// }

		// // Debug: Decode token without validation
		// decodedClaims, err := DecodeToken(authHeader)
		// if err != nil {
		// 	logger.LogError("Failed to decode token", "error", err)
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
		// 	c.Abort()
		// 	return
		// }

		// logger.LogInfo("JWT token decoded successfully", "claims", decodedClaims)
		// c.Set("userID", decodedClaims.UserID)
		c.Set("userID", "ebec5b8d-d656-452e-9fe8-0bab337c59b3")
		tenantID := c.GetHeader("X-Tenant-ID")

		if tenantID != "" {
			c.Set("tenantID", tenantID)
		}

		c.Next()
	}
}
