package middlewares

import (
	"context"
	"encoding/json"
	"iam_services_main_v1/config"
	"iam_services_main_v1/helpers"
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
		// Allow introspection queries without permission checks
		if action == "IntrospectionQuery" || strings.Contains(req.Query, "__schema") {
			logger.LogInfo("Allowing GraphQL introspection query")
			c.Next()
			return
		}
		resourceType := deriveResourceType(action)
		if resourceType == "" {
			logger.LogError("Unknown resourceType for action", "action", action)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unknown action"})
			return
		}

		logger.LogInfo("GraphQL request", "action", action, "resourceType", resourceType)

		ctx := c.Request.Context()
		authorized, err := AuthorizationMiddleware(ctx, psc, action, resourceType, "*")
		if err != nil || !authorized {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User does not have permission to perform this action"})
			return
		}

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
	if strings.Contains(strings.ToLower(action), "tenant") {
		return config.TenantResourceTypeID
	}
	if strings.Contains(strings.ToLower(action), "clientorganizationunit") {
		return config.ClientOrgUnitResourceTypeID
	}
	if strings.Contains(strings.ToLower(action), "account") {
		return config.AccountResourceTypeID
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

// AuthMiddleware creates a Gin middleware for JWT authentication with the default OpenID endpoint
func AuthorizationMiddleware(ctx context.Context, psc *permit.PermitSdkService, action, resourceType, resourceId string) (bool, error) {

	logger.LogInfo("Authorization middleware invoked", "action", action, "resourceType", resourceType, "resourceId", resourceId)
	// Get user and tenant context
	userID, tenantID, err := helpers.GetUserAndTenantID(ctx)
	logger.LogInfo("User ID and Tenant ID", "userID", userID, "tenantID", tenantID)
	if err != nil {
		return false, err
	}
	logger.LogInfo("User ID and Tenant ID extracted", "userID", userID, "tenantID", tenantID)
	logger.LogInfo("the action is and resourceType is", "action", action, "resourceType", resourceType)
	// Check permission
	_, err = psc.Check(ctx, userID.String(), strings.ToLower(action), resourceType, resourceId, tenantID.String())
	if err != nil {
		return false, err
	}
	logger.LogInfo("Authorization check passed", "userId", userID, "tenantId", tenantID, "action", action, "resourceType", resourceType, "resourceId", resourceId)
	return true, nil
}
