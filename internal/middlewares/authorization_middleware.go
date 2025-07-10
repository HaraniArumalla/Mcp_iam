package middlewares

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"

	"iam_services_main_v1/config"
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/pkg/logger"
)

type graphQLRequest struct {
	OperationName string `json:"operationName"`
	Query         string `json:"query"`
}

func GraphQLAuthMiddleware(psc *permit.PermitSdkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPost {
			c.Next()
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		c.Request.Body = io.NopCloser(strings.NewReader(string(body))) // Restore body for downstream

		var req graphQLRequest
		if err := json.Unmarshal(body, &req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid GraphQL request"})
			return
		}

		action := extractAction(req)
		logger.LogInfo("Extracted action from GraphQL request", "action", action)
		// Allow introspection queries
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

// Uses the graphql-go parser to extract the first field name from the query
func extractAction(req graphQLRequest) string {
	if req.OperationName != "" {
		return req.OperationName
	}

	src := source.NewSource(&source.Source{
		Body: []byte(req.Query),
		Name: "GraphQL Request",
	})

	doc, err := parser.Parse(parser.ParseParams{Source: src})
	if err != nil {
		logger.LogError("Failed to parse GraphQL query", "error", err)
		return ""
	}

	for _, def := range doc.Definitions {
		if op, ok := def.(*ast.OperationDefinition); ok {
			if len(op.SelectionSet.Selections) > 0 {
				if field, ok := op.SelectionSet.Selections[0].(*ast.Field); ok {
					fullName := field.Name.Value
					logger.LogInfo("Extracted full action name from GraphQL query", "fullName", fullName)
					// Strip known subgraph suffix like _mcp_iam_o
					// if strings.HasSuffix(fullName, config.SubgraphName) {
					// 	return strings.TrimSuffix(fullName, config.SubgraphName)
					// }
					// Find last occurrence of double underscores (used by supergraph/subgraph)
					if idx := strings.LastIndex(fullName, "__"); idx != -1 {
						logger.LogInfo("Stripping subgraph suffix from action name", "fullName", fullName, "index", idx)
						logger.LogInfo("Returning action name without subgraph suffix", "actionName", fullName[:idx])
						return fullName[:idx]
					}

					// Default: return full name
					return fullName
				}
			}
		}
	}
	return ""
}

// Returns resource type based on action naming conventions
func deriveResourceType(action string) string {
	lower := strings.ToLower(action)

	switch {
	case strings.Contains(lower, "tenant"):
		return config.TenantResourceTypeID
	case strings.Contains(lower, "clientorganizationunit"):
		return config.ClientOrgUnitResourceTypeID
	case strings.Contains(lower, "account"):
		return config.AccountResourceTypeID
	case strings.Contains(lower, "role"):
		return config.RoleResourceTypeID
	case strings.Contains(lower, "permission"):
		return config.PermissionResourceTypeID
	case strings.Contains(lower, "binding"):
		return config.BindingResourceTypeID
	default:
		return ""
	}
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, val) {
			return true
		}
	}
	return false
}

// Middleware wrapper that checks authorization via Permit.io
func AuthorizationMiddleware(ctx context.Context, psc *permit.PermitSdkService, action, resourceType, resourceId string) (bool, error) {
	logger.LogInfo("Authorization middleware invoked", "action", action, "resourceType", resourceType, "resourceId", resourceId)

	userID, tenantID, err := helpers.GetUserAndTenantID(ctx)
	if err != nil {
		logger.LogError("Failed to extract user and tenant IDs", "error", err)
		return false, err
	}

	logger.LogInfo("User ID and Tenant ID extracted", "userID", userID, "tenantID", tenantID)

	_, err = psc.Check(ctx, userID.String(), strings.ToLower(action), resourceType, resourceId, tenantID.String())
	if err != nil {
		logger.LogError("Authorization failed", "error", err)
		return false, err
	}

	logger.LogInfo("Authorization check passed", "userId", userID, "tenantId", tenantID, "action", action, "resourceType", resourceType, "resourceId", resourceId)
	return true, nil
}
