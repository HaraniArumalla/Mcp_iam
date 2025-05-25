package middlewares

import (
	"context"
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/pkg/logger"
)

// AuthMiddleware creates a Gin middleware for JWT authentication with the default OpenID endpoint
func AuthorizationMiddleware(ctx context.Context, psc *permit.PermitSdkService, action, resourceType, resourceId string) (bool, error) {

	logger.LogInfo("Authorization middleware invoked", "action", action, "resourceType", resourceType, "resourceId", resourceId)
	// Get user and tenant context
	userID, tenantID, err := helpers.GetUserAndTenantID(ctx)
	logger.LogInfo("User ID and Tenant ID", "userID", userID, "tenantID", tenantID)
	if err != nil {
		return false, err
	}

	// Check permission
	_, err = psc.Check(ctx, userID.String(), action, resourceType, resourceId, tenantID.String())
	if err != nil {
		return false, err
	}
	logger.LogInfo("Authorization check passed", "userId", userID, "tenantId", tenantID, "action", action, "resourceType", resourceType, "resourceId", resourceId)
	return true, nil
}
