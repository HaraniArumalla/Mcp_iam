package resources

import (
	"context"
	"iam_services_main_v1/gql/models"

	"github.com/google/uuid"
)

type ResourceQueryResolver struct{}

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
