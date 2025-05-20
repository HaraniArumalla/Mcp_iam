package middlewares

import (
	"iam_services_main_v1/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGinContextToContextMiddleware(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create a test recorder and router
	w := httptest.NewRecorder()
	router := gin.New()

	// Apply our middleware
	router.Use(GinContextToContextMiddleware())

	// Define a handler that checks for the Gin context in the request context
	router.GET("/test", func(c *gin.Context) {
		// Try to retrieve the Gin context from the request context
		ginCtx, exists := c.Request.Context().Value(config.GinContextKey).(*gin.Context)

		// Verify that the Gin context exists and is the same as our current context
		if !exists || ginCtx != c {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gin context not found or incorrect"})
			return
		}

		// Set a value in the Gin context to verify we can use it
		ginCtx.Set("testKey", "testValue")

		// Verify we can retrieve the value from the current context
		value, exists := c.Get("testKey")
		if !exists || value != "testValue" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve value from Gin context"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	// Verify the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestGinContextToContextMiddlewareMultipleHandlers(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create a test recorder and router
	w := httptest.NewRecorder()
	router := gin.New()

	// Apply our middleware
	router.Use(GinContextToContextMiddleware())

	// Define a middleware that sets a value in the Gin context
	setValueMiddleware := func(c *gin.Context) {
		c.Set("contextValue", "middlewareValue")
		c.Next()
	}

	// Apply the middleware
	router.Use(setValueMiddleware)

	// Define a handler that checks if the value is accessible
	router.GET("/test", func(c *gin.Context) {
		// Try to retrieve the Gin context from the request context
		ginCtx, exists := c.Request.Context().Value(config.GinContextKey).(*gin.Context)

		// Verify that the Gin context exists
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gin context not found"})
			return
		}

		// Try to get the value set by the previous middleware
		value, exists := ginCtx.Get("contextValue")
		if !exists || value != "middlewareValue" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Context value not found or incorrect"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "value": value})
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	// Verify the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
	assert.Contains(t, w.Body.String(), "middlewareValue")
}

func TestGinContextToContextMiddlewareWithCustomContextKey(t *testing.T) {
	// This test verifies that the middleware uses the correct key to store the Gin context

	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create a test recorder and router
	w := httptest.NewRecorder()
	router := gin.New()

	// Apply our middleware
	router.Use(GinContextToContextMiddleware())

	// Define a handler that verifies the exact context key used
	router.GET("/test", func(c *gin.Context) {
		// The context should be stored with key config.GinContextKey specifically
		// Try to retrieve with the exact key constant
		ginCtx, exists := c.Request.Context().Value(config.GinContextKey).(*gin.Context)

		// Verify that the Gin context exists and is the correct one
		if !exists || ginCtx != c {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gin context not found or incorrect using GinContextKey"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	// Verify the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestGinContextToContextMiddlewareWithNextCall(t *testing.T) {
	// This test verifies that the middleware properly calls the next handler in the chain

	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create a test recorder and router
	w := httptest.NewRecorder()
	router := gin.New()

	// Track if the middleware calls Next()
	middlewareCalled := false
	nextHandlerCalled := false

	// Custom middleware to track the call chain
	router.Use(func(c *gin.Context) {
		middlewareCalled = true
		c.Next()
		// This should execute after all other handlers
		assert.True(t, nextHandlerCalled, "Next handler should have been called")
	})

	// Apply our middleware
	router.Use(GinContextToContextMiddleware())

	// Define a handler that verifies the middleware chain
	router.GET("/test", func(c *gin.Context) {
		// Verify that previous middleware was called
		assert.True(t, middlewareCalled, "Previous middleware should have been called")

		nextHandlerCalled = true
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	// Verify the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
	assert.True(t, middlewareCalled, "Middleware should have been called")
	assert.True(t, nextHandlerCalled, "Next handler should have been called")
}
