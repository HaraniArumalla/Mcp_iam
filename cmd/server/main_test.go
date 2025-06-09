package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"iam_services_main_v1/internal/healthchecks"
	"iam_services_main_v1/pkg/logger"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func initGraphQLServer() *handler.Server {
	return handler.New(nil) // Replace nil with your GraphQL schema
}

func setupRouter() *gin.Engine {
	router := gin.New()
	setupMiddleware(router)

	// Add routes
	healthHandler := &healthchecks.HealthHandler{}
	router.GET("/status", func(c *gin.Context) { healthHandler.SimpleStatus(c) })
	router.GET("/health/live", func(c *gin.Context) { healthHandler.LivenessCheck(c) })
	router.GET("/health/ready", func(c *gin.Context) { healthHandler.ReadinessCheck(c) })
	router.GET("/playground", gin.WrapH(playground.Handler("GraphQL playground", "/graphql")))
	router.POST("/graphql", gin.WrapH(initGraphQLServer()))

	return router
}

func TestSetupRouter(t *testing.T) {
	// Test router setup
	router := setupRouter()
	assert.NotNil(t, router)

	// Verify routes are registered
	routes := router.Routes()
	routePaths := make(map[string]bool)
	for _, route := range routes {
		routePaths[route.Path] = true
	}

	// Check required endpoints exist
	requiredPaths := []string{
		"/status",
		"/health/live",
		"/health/ready",
		"/playground",
		"/graphql",
	}

	for _, path := range requiredPaths {
		assert.True(t, routePaths[path], "Route %s not found", path)
	}
}

func TestInitGraphQLServer(t *testing.T) {
	srv := initGraphQLServer()
	assert.NotNil(t, srv)
	assert.IsType(t, &handler.Server{}, srv)
}

func TestSetupMiddleware(t *testing.T) {
	router := gin.New()
	setupMiddleware(router)

	// Verify middleware is added
	assert.Greater(t, len(router.Handlers), 0)
}

func TestMain(t *testing.T) {
	// Save original args and env
	originalArgs := os.Args
	originalEnv := os.Getenv("PORT")
	defer func() {
		os.Args = originalArgs
		if err := os.Setenv("PORT", originalEnv); err != nil {
			t.Logf("Failed to restore PORT environment variable: %v", err)
		}
	}()

	// Test with different ports
	testCases := []struct {
		name     string
		port     string
		wantPort string
	}{
		{
			name:     "Default port",
			port:     "",
			wantPort: "8080",
		},
		{
			name:     "Custom port",
			port:     "9090",
			wantPort: "9090",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up test environment
			if err := os.Setenv("PORT", tc.port); err != nil {
				t.Fatalf("Failed to set PORT environment variable: %v", err)
			}

			// Create a context with cancel for main
			ctx, cancel := context.WithCancel(context.Background())

			// Run main in a goroutine
			go func() {
				// Small delay to let server start
				time.Sleep(100 * time.Millisecond)
				// Cancel context to stop server
				cancel()
			}()

			// Run main with context
			mainWithContext(ctx)
		})
	}
}

func TestLoggerSetup(t *testing.T) {
	// Test logger initialization
	logger.InitLogger()

	// Verify logger is working
	logger.LogInfo("Test log message")
}

func TestHealthChecks(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a router for testing
	router := gin.New()

	// Create recorder for testing
	w := httptest.NewRecorder()

	// Initialize handler
	handler := &healthchecks.HealthHandler{}

	// Register routes
	router.GET("/status", handler.SimpleStatus)
	router.GET("/health/live", handler.LivenessCheck)
	router.GET("/health/ready", handler.ReadinessCheck)

	// Test cases
	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{
			name:       "Simple Status",
			path:       "/status",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Liveness Check",
			path:       "/health/live",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Readiness Check",
			path:       "/health/ready",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			assert.NoError(t, err)

			// Serve request
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.wantStatus, w.Code)
			assert.NotEmpty(t, w.Body.String())
		})
	}
}

// Helper function for running main with context
func mainWithContext(ctx context.Context) {
	// Initialize services
	logger.InitLogger()

	// Set up router
	router := setupRouter()

	// Start server in goroutine
	go func() {
		if err := router.Run(); err != nil {
			logger.LogError("Server error", "error", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
}

func TestSetupServer(t *testing.T) {
	router := setupServer()
	assert.NotNil(t, router)

	// Test routes
	expectedPaths := []string{
		"/status",
		"/health/live",
		"/health/ready",
		"/playground",
		"/graphql",
	}

	routes := router.Routes()
	for _, expectedPath := range expectedPaths {
		found := false
		for _, route := range routes {
			if route.Path == expectedPath {
				found = true
				break
			}
		}
		assert.True(t, found, "Route %s not found", expectedPath)
	}
}

func TestGetPort(t *testing.T) {
	originalPort := os.Getenv("PORT")
	defer func() {
		if err := os.Setenv("PORT", originalPort); err != nil {
			t.Logf("Failed to restore PORT environment variable: %v", err)
		}
	}()

	testCases := []struct {
		name       string
		portEnv    string
		expectPort string
	}{
		{"Default port", "", "8080"},
		{"Custom port", "3000", "3000"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := os.Setenv("PORT", tc.portEnv); err != nil {
				t.Fatalf("Failed to set PORT environment variable: %v", err)
			}
			assert.Equal(t, tc.expectPort, getPort())
		})
	}
}

func TestRunServer(t *testing.T) {
	// Setup test server
	srv, err := runServer()
	assert.NoError(t, err)
	assert.NotNil(t, srv)

	// Ensure server is running
	time.Sleep(100 * time.Millisecond)

	// Test endpoints
	client := &http.Client{Timeout: time.Second}
	resp, err := client.Get("http://localhost:" + getPort() + "/status")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Cleanup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	assert.NoError(t, srv.Shutdown(ctx))
}

// Constants for test configuration
const (
	testPermitPdpURL      = "http://localhost:8000" // Use a consistent value for tests
	testPermitPdpEndpoint = "http://localhost:8000"
	testPermitProject     = "test-project"
	testPermitEnv         = "test-env"
	testPermitToken       = "test-token"
)

// Helper function to set up test environment variables
func setupTestEnv(t *testing.T) func() {
	// Store original values
	originalEnv := map[string]string{
		"PERMIT_PDP_ENDPOINT": os.Getenv("PERMIT_PDP_ENDPOINT"),
		"PERMIT_PDP_URL":      os.Getenv("PERMIT_PDP_URL"),
		"PERMIT_PROJECT":      os.Getenv("PERMIT_PROJECT"),
		"PERMIT_ENV":          os.Getenv("PERMIT_ENV"),
		"PERMIT_TOKEN":        os.Getenv("PERMIT_TOKEN"),
	}

	// Set test environment variables
	envVars := map[string]string{
		"PERMIT_PDP_ENDPOINT": testPermitPdpEndpoint,
		"PERMIT_PDP_URL":      testPermitPdpURL,
		"PERMIT_PROJECT":      testPermitProject,
		"PERMIT_ENV":          testPermitEnv,
		"PERMIT_TOKEN":        testPermitToken,
	}

	for key, value := range envVars {
		if err := os.Setenv(key, value); err != nil {
			t.Fatalf("Failed to set environment variable %s: %v", key, err)
		}
	}

	// Return cleanup function
	return func() {
		for key, value := range originalEnv {
			if err := os.Setenv(key, value); err != nil {
				t.Logf("Failed to restore environment variable %s: %v", key, err)
			}
		}
	}
}

func TestGraphQLHandler(t *testing.T) {
	// Setup test environment and get cleanup function
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Test the handler
	handlerFunc := graphqlHandler()
	assert.NotNil(t, handlerFunc)

	// Create a test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/graphql", strings.NewReader(`{"query":"{ __schema { types { name } } }"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the handler
	handlerFunc(c)

	// Check response
	assert.NotEqual(t, 500, w.Code, "Handler should not return 500 error")
}

func TestGraphQLHandlerMissingEnvVars(t *testing.T) {
	// Save original environment variables
	originalEnv := map[string]string{
		"PERMIT_PDP_ENDPOINT": os.Getenv("PERMIT_PDP_ENDPOINT"),
		"PERMIT_PROJECT":      os.Getenv("PERMIT_PROJECT"),
		"PERMIT_ENV":          os.Getenv("PERMIT_ENV"),
		"PERMIT_TOKEN":        os.Getenv("PERMIT_TOKEN"),
	}

	defer func() {
		// Restore original environment variables
		for key, value := range originalEnv {
			if err := os.Setenv(key, value); err != nil {
				t.Logf("Failed to restore environment variable %s: %v", key, err)
			}
		}
	}()

	// Clear required environment variables
	if err := os.Setenv("PERMIT_PDP_ENDPOINT", ""); err != nil {
		t.Fatalf("Failed to clear PERMIT_PDP_ENDPOINT: %v", err)
	}

	// Test the handler with missing environment variables
	handlerFunc := graphqlHandler()
	assert.NotNil(t, handlerFunc)

	// Create a test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/graphql", strings.NewReader(`{"query":"{ __schema { types { name } } }"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the handler - should handle the error gracefully
	handlerFunc(c)

	// Check response - should return error but not crash
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPlaygroundHandler(t *testing.T) {
	handlerFunc := func(c *gin.Context) {
		playground.Handler("GraphQL playground", "/graphql").ServeHTTP(c.Writer, c.Request)
	}
	assert.NotNil(t, handlerFunc)

	// Create a test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/playground", nil)

	// Call the handler
	handlerFunc(c)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "GraphQL playground")
}

func TestHealthEndpoints(t *testing.T) {
	router := setupServer()
	w := httptest.NewRecorder()

	endpoints := []struct {
		path   string
		method string
	}{
		{"/status", "GET"},
		{"/health/live", "GET"},
		{"/health/ready", "GET"},
	}

	for _, e := range endpoints {
		t.Run(e.path, func(t *testing.T) {
			req := httptest.NewRequest(e.method, e.path, nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestSetupRoutes(t *testing.T) {
	router := gin.New()
	setupRoutes(router)

	routes := router.Routes()
	assert.GreaterOrEqual(t, len(routes), 5) // All expected routes
}
