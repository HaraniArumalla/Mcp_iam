package resources

import (
	"context"
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/internal/permit"
	"strings"

	"iam_services_main_v1/pkg/logger"

	"github.com/google/uuid"
)

type ResourceQueryResolver struct {
	PSC *permit.PermitSdkService
}

func (r *ResourceQueryResolver) Resource(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	// tenantID, _ := helpers.GetTenantID(ctx)

	// conditions := map[string]interface{}{
	// 	"row_status": 1,
	// }

	// if tenantID != nil {
	// 	conditions["tenant_id"] = *tenantID
	// }
	// resource, err := dao.GetResourceItem(r.DB, conditions)
	// if err != nil {
	// 	return nil, fmt.Errorf("error occured when fetching the resources: %w", err)
	// }

	// resourceObj, err := dto.MapToResource(resource.ResourceTypeID.String(), *resource)
	// if err != nil {
	// 	return nil, fmt.Errorf("error mapping resource: %w", err)
	// }
	// if resource, ok := resourceObj.(models.Resource); ok {
	// 	return resource, nil
	// } else {
	// 	return nil, fmt.Errorf("resource doesn't implement models.resource: %v", resourceObj)
	// }
	return nil, nil
}

func (r *ResourceQueryResolver) Resources(ctx context.Context) (models.OperationResult, error) {
	// tenantID, _ := helpers.GetTenantID(ctx)

	// conditions := map[string]interface{}{
	// 	"row_status": 1,
	// }

	// if tenantID != nil {
	// 	conditions["tenant_id"] = *tenantID
	// }
	// resources, err := dao.GetAllResources(r.DB, conditions)
	// if err != nil {
	// 	return nil, fmt.Errorf("error occured when fetching the resources: %w", err)
	// }

	// allResources := make([]models.Resource, 0, len(resources))
	// for _, resource := range resources {
	// 	resourceObj, err := dto.MapToResource(resource.ResourceTypeID.String(), resource)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error mapping resource: %w", err)
	// 	}
	// 	if resource, ok := resourceObj.(models.Resource); ok {
	// 		allResources = append(allResources, resource)
	// 	} else {
	// 		return nil, fmt.Errorf("resource doesn't implement models.resource: %v", resourceObj)
	// 	}

	// }

	// return allResources, nil
	return nil, nil
}

// CheckPermission is the resolver for the checkPermission field.
func (r *ResourceQueryResolver) CheckPermission(ctx context.Context, input models.PermissionInput) (*models.PermissionResponse, error) {
	// Get user and tenant context
	userID, tenantID, err := helpers.GetUserAndTenantID(ctx)
	if err != nil {
		logger.LogError("User ID not found in context during permission check")
		return &models.PermissionResponse{
			Allowed: false,
			Error:   helpers.Ptr("User ID & Tenant ID not found in context"),
		}, nil
	}
	// Validate input.Action
	if input.Action == "" {
		return &models.PermissionResponse{
			Allowed: false,
			Error:   helpers.Ptr("Action cannot be empty"),
		}, nil
	}
	// Check permission
	resourceID := ""
	if input.ResourceID != nil {
		resourceID = *input.ResourceID
	}
	// Check if input.Action contains "create"
	if strings.Contains(strings.ToLower(input.Action), "create") {
		resourceID = ""
	}

	// Log the permission check request
	logger.LogInfo("GraphQL permission check request",
		"user_id", userID,
		"action", input.Action,
		"resource_type", input.ResourceType,
		"resource_id", resourceID,
	)

	// Check permission
	allowed, err := r.PSC.Check(ctx, userID.String(), input.Action, input.ResourceType, resourceID, tenantID.String())
	if err != nil {
		logger.LogError("Failed to check permissions", "error", err)
		return &models.PermissionResponse{
			Allowed: false,
			Error:   helpers.Ptr(fmt.Sprintf("Failed to check permissions: %s", err.Error())),
		}, nil
	}

	// Return result
	return &models.PermissionResponse{
		Allowed: allowed,
		Error:   helpers.Ptr(""),
	}, nil
}
