package healthchecks

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

// SimpleStatus returns a basic OK response without timestamp
func (h *HealthHandler) SimpleStatus(c *gin.Context) {
	c.JSON(http.StatusOK, StatusResponse{
		Status: "ok",
	})
}

// LivenessCheck checks if service is alive
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "up",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// ReadinessCheck checks if service is ready
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "up",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}
