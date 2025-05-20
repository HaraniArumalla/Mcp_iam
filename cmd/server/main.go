package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"iam_services_main_v1/config"
	"iam_services_main_v1/gql"
	"iam_services_main_v1/gql/generated"
	"iam_services_main_v1/internal/healthchecks"
	"iam_services_main_v1/internal/middlewares"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

// setupServer initializes and returns the configured gin router
func setupServer() *gin.Engine {
	// Initialize services
	logger.InitLogger()
	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		logger.LogFatal("failed to load environment variables", "error", err)
	}

	// Set up router
	router := gin.New()
	setupMiddleware(router)
	setupRoutes(router)

	return router
}

// setupMiddleware adds middleware to the router
func setupMiddleware(router *gin.Engine) {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.CORSMiddleware())

}

// setupRoutes configures all routes for the server
func setupRoutes(router *gin.Engine) {
	healthHandler := &healthchecks.HealthHandler{}
	router.GET("/status", healthHandler.SimpleStatus)
	router.GET("/health/live", healthHandler.LivenessCheck)
	router.GET("/health/ready", healthHandler.ReadinessCheck)
	router.Use(middlewares.RequestLogger())
	router.Use(middlewares.GinContextToContextMiddleware())
	router.GET("/playground", gin.WrapH(playground.Handler("GraphQL playground", "/graphql")))
	router.Use(middlewares.AuthMiddleware())
	router.POST("/graphql", graphqlHandler())
}

// graphqlHandler creates and returns the GraphQL handler
func graphqlHandler() gin.HandlerFunc {
	// Create permit service with panic recovery
	permitclint := NewPermitClient()
	if permitclint == nil {
		logger.LogError("Failed to create Permit client")
		return func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
		}
	}

	permitService := permit.NewPermitServiceImpl(permitclint)
	config := generated.Config{
		Resolvers: &gql.Resolver{PC: permitService},
	}

	// Create handler using preferred constructor
	h := handler.New(generated.NewExecutableSchema(config))

	// Enable introspection
	h.Use(extension.Introspection{})

	// Add all supported transports
	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})
	h.AddTransport(transport.MultipartForm{})

	// Wrap with recovery to prevent application crashes
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.LogError("GraphQL handler panic", "error", r)
				debug.PrintStack() // Print stack trace to logs
				c.JSON(http.StatusInternalServerError, gin.H{
					"errors": []map[string]interface{}{{"message": "Internal server error"}},
				})
			}
		}()

		h.ServeHTTP(c.Writer, c.Request)
	}
}

// getPort returns the server port from environment or default
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

// runServer starts the HTTP server and returns a shutdown function
func runServer() (*http.Server, error) {
	router := setupServer()
	srv := &http.Server{
		Addr:    ":" + getPort(),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.LogError("Server failed", "error", err)
		}
	}()

	return srv, nil
}

func main() {
	// Set up panic recovery for the main function
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("FATAL: Application panic: %v\n", r)
			// Print stack trace
			debug.PrintStack()
			os.Exit(1)
		}
	}()

	// Initialize logger
	logger.InitLogger()
	logger.LogInfo("Starting server initialization")

	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		logger.LogError("Failed to load environment variables", "error", err)
		os.Exit(1)
	}
	logger.LogInfo("Environment variables loaded successfully")

	// Start the server
	srv, err := runServer()
	if err != nil {
		logger.LogFatal("Failed to start server", "error", err)
	}
	logger.LogInfo("Server started successfully on port " + getPort())

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.LogError("Server forced to shutdown", "error", err)
	}
}

func NewPermitClient() *permit.PermitClient {
	baseURL := os.Getenv("PERMIT_PDP_ENDPOINT")
	projectID := os.Getenv("PERMIT_PROJECT")
	envID := os.Getenv("PERMIT_ENV")
	apiKey := os.Getenv("PERMIT_TOKEN")

	if baseURL == "" || projectID == "" || envID == "" || apiKey == "" {
		logger.LogFatal("One or more required environment variables are not set")
		return nil
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Disable certificate verification
		},
	}

	return &permit.PermitClient{
		BaseURL: fmt.Sprintf("%s/v2/facts/%s/%s", baseURL, projectID, envID),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", apiKey),
			"Content-Type":  "application/json",
		},
		Client: &http.Client{Transport: transport, Timeout: 30 * time.Second},
	}
}
