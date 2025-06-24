package middlewares

import (
	"encoding/json"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/pkg/logger"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type graphQLRequest struct {
	OperationName string `json:"operationName"`
	Query         string `json:"query"`
}

func GraphQLAuthMiddleware(psc *permit.PermitSdkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only process POST /graphql
		if c.Request.Method != http.MethodPost {
			c.Next()
			return
		}

		// Read body safely
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		c.Request.Body = io.NopCloser(strings.NewReader(string(body))) // Reuse body

		var req graphQLRequest
		if err := json.Unmarshal(body, &req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid GraphQL request"})
			return
		}

		// Extract action from operationName or query body
		action := extractAction(req)
		resourceType := deriveResourceType(action)
		if resourceType == "" {
			logger.LogError("Unknown resourceType for action", "action", action)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unknown action"})
			return
		}

		logger.LogInfo("GraphQL request", "action", action, "resourceType", resourceType)

		c.Next()
	}
}

func extractAction(req graphQLRequest) string {
	if req.OperationName != "" {
		return req.OperationName
	}

	lines := strings.Split(req.Query, "\n")
	inOperation := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "mutation") || strings.HasPrefix(line, "query") {
			inOperation = true
			continue
		}

		if inOperation {
			// skip braces or empty lines
			if strings.HasPrefix(line, "{") || strings.HasPrefix(line, "}") {
				continue
			}

			// Example: "createAccount(input: {...}) {" → extract only "createAccount"
			firstPart := strings.Fields(line)
			if len(firstPart) > 0 {
				action := firstPart[0]
				// strip trailing `(` if it exists
				if idx := strings.Index(action, "("); idx != -1 {
					action = action[:idx]
				}
				return action
			}
		}
	}

	return ""
}

// deriveResourceType returns resource type based on action
func deriveResourceType(action string) string {
	tenantActions := []string{"createTenant", "updateTenant", "deleteTenant", "tenants", "tenant"}
	accountActions := []string{"createAccount", "updateAccount", "deleteAccount", "accounts", "account"}

	if contains(tenantActions, action) {
		return "ed113bda-bbda-11ef-87ea-c03c5946f955"
	}
	if contains(accountActions, action) {
		return "ed113f30-bbda-11ef-87ea-c03c5946f955"
	}
	return ""
}

// Utility functions
func contains(slice []string, val string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, val) {
			return true
		}
	}
	return false
}
