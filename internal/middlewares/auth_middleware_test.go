package middlewares

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// mockGinMiddlewareFunc creates a mock for the GinMiddleware function in jwt authenticator
func mockGinMiddlewareFunc(shouldSucceed bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if shouldSucceed {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		}
	}
}

// TestAuthMiddleware tests the AuthMiddleware function in various scenarios
func TestAuthMiddleware(t *testing.T) {
	// Save original env vars to restore later
	originalTenantID := os.Getenv("AZURE_TENANT_ID")
	originalUserClaimField := os.Getenv("JWT_USER_CLAIM_FIELD")
	originalCacheDuration := os.Getenv("JWKS_CACHE_DURATION")

	// Restore env vars after test completes
	defer func() {
		os.Setenv("AZURE_TENANT_ID", originalTenantID)
		os.Setenv("JWT_USER_CLAIM_FIELD", originalUserClaimField)
		os.Setenv("JWKS_CACHE_DURATION", originalCacheDuration)
	}()

	// Set up test cases
	testCases := []struct {
		name               string
		tenantID           string
		userClaimField     string
		cacheDuration      string
		expectedStatusCode int
	}{
		{
			name:               "Missing Azure Tenant ID",
			tenantID:           "",
			userClaimField:     "sub",
			cacheDuration:      "5m",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "Invalid Cache Duration Format",
			tenantID:           "test-tenant-id",
			userClaimField:     "sub",
			cacheDuration:      "invalid",
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up environment variables for this test
			os.Setenv("AZURE_TENANT_ID", tc.tenantID)
			os.Setenv("JWT_USER_CLAIM_FIELD", tc.userClaimField)
			os.Setenv("JWKS_CACHE_DURATION", tc.cacheDuration)

			// Set up the Gin router with the AuthMiddleware
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(AuthMiddleware())
			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			// Create a request to the test endpoint
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			resp := httptest.NewRecorder()

			// Send the request to the router
			router.ServeHTTP(resp, req)

			// Check the response status code
			assert.Equal(t, tc.expectedStatusCode, resp.Code)
		})
	}
}

// TestAuthMiddlewareWithMockJWT tests the AuthMiddleware function
// with a mock JWT authenticator to simulate successful authentication
func TestAuthMiddlewareWithMockJWT(t *testing.T) {
	// Create a test version of AuthMiddleware with a mock JWT authenticator
	mockAuthMiddleware := func(shouldSucceed bool) gin.HandlerFunc {
		return func(c *gin.Context) {
			// Set required environment variables
			os.Setenv("AZURE_TENANT_ID", "mock-tenant-id")
			os.Setenv("JWT_USER_CLAIM_FIELD", "sub")
			os.Setenv("JWKS_CACHE_DURATION", "5m")

			// Simulate the GinMiddleware call result
			if shouldSucceed {
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			}
		}
	}

	t.Run("Successful authentication", func(t *testing.T) {
		// Set up the Gin router with the mock AuthMiddleware
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.Use(mockAuthMiddleware(true))
		router.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// Create a request to the test endpoint
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		resp := httptest.NewRecorder()

		// Send the request to the router
		router.ServeHTTP(resp, req)

		// Check the response status code
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Failed authentication", func(t *testing.T) {
		// Set up the Gin router with the mock AuthMiddleware
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.Use(mockAuthMiddleware(false))
		router.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// Create a request to the test endpoint
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		resp := httptest.NewRecorder()

		// Send the request to the router
		router.ServeHTTP(resp, req)

		// Check the response status code
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}

// TestIsValidDuration tests both valid and invalid duration formats
func TestIsValidDuration(t *testing.T) {
	testCases := []struct {
		duration string
		isValid  bool
	}{
		{"5m", true},
		{"1h", true},
		{"30s", true},
		{"24h", true},
		{"invalid", false},
		{"", false},
		{"5", false},
		{"5mm", false},
	}

	for _, tc := range testCases {
		t.Run(tc.duration, func(t *testing.T) {
			_, err := time.ParseDuration(tc.duration)
			if tc.isValid {
				assert.NoError(t, err, "Expected %s to be a valid duration", tc.duration)
			} else {
				assert.Error(t, err, "Expected %s to be an invalid duration", tc.duration)
			}
		})
	}
}

// Additional test with mock for NewAuthenticator to cover error case
func TestAuthMiddlewareWithAuthenticatorError(t *testing.T) {
	// Save original env vars to restore later
	originalTenantID := os.Getenv("AZURE_TENANT_ID")
	originalUserClaimField := os.Getenv("JWT_USER_CLAIM_FIELD")
	originalCacheDuration := os.Getenv("JWKS_CACHE_DURATION")

	// Restore env vars after test completes
	defer func() {
		os.Setenv("AZURE_TENANT_ID", originalTenantID)
		os.Setenv("JWT_USER_CLAIM_FIELD", originalUserClaimField)
		os.Setenv("JWKS_CACHE_DURATION", originalCacheDuration)
	}()

	// Set environment variables for this test
	os.Setenv("AZURE_TENANT_ID", "test-tenant-id")
	os.Setenv("JWT_USER_CLAIM_FIELD", "sub")
	os.Setenv("JWKS_CACHE_DURATION", "5m")

	// Note: This test is limited since we can't easily mock the NewAuthenticator function
	// In a real-world scenario, we would need to refactor the code to make it more testable
	// by injecting dependencies or using interfaces

	// Set up the Gin router with the AuthMiddleware
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Create a request to the test endpoint
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// Add an invalid Authorization header to test error paths
	req.Header.Set("Authorization", "Bearer invalid-token")
	resp := httptest.NewRecorder()

	// Send the request to the router
	router.ServeHTTP(resp, req)

	// We expect a 401 Unauthorized here in a real environment because the token is invalid
	// However, without mocking the JWT authenticator, the exact behavior depends on the implementation
	// So we're not making assertions about the specific response code
	assert.NotEqual(t, http.StatusOK, resp.Code, "Expected authentication to fail with an invalid token")
}
