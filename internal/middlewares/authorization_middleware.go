package middlewares

import (
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/pkg/logger"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a Gin middleware for JWT authentication with the default OpenID endpoint
func AuthorizationMiddleware(psc *permit.PermitSdkService, action, resourceType, resourceId string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user and tenant context
		ctx := c.Request.Context()
		userID, tenantID, err := helpers.GetUserAndTenantID(ctx)
		logger.LogInfo("User ID and Tenant ID", "userID", userID, "tenantID", tenantID)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		// Check permission
		_, err = psc.Check(ctx, userID.String(), action, resourceType, resourceId, tenantID.String())
		if err != nil {
			c.AbortWithError(404, err)
			return
		}
		c.Next()
		logger.LogInfo("Authorization check passed", "userId", userID, "tenantId", tenantID, "action", action, "resourceType", resourceType, "resourceId", resourceId)
	}
}
