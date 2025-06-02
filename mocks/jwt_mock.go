package mocks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MockAuthenticator is a mock implementation of JWT authenticator for testing
type MockAuthenticator struct {
	ShouldSucceed bool
}

// GinMiddleware returns a mock middleware function that simulates JWT authentication
func (m *MockAuthenticator) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m.ShouldSucceed {
			// Simulate successful authentication by setting user claims
			c.Set("userID", "mock-user-id")
			c.Next()
		} else {
			// Simulate authentication failure
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
		}
	}
}

// NewMockAuthenticator creates a new MockAuthenticator
func NewMockAuthenticator(shouldSucceed bool) *MockAuthenticator {
	return &MockAuthenticator{
		ShouldSucceed: shouldSucceed,
	}
}
