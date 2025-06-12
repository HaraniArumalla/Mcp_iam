package resourcetypes

import (
	"context"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/utils"
	"iam_services_main_v1/pkg/logger"
	"net/http"
)

type ResourceTypeQueryResolver struct {
	PC permit.PermitService
}

func (r *ResourceTypeQueryResolver) ResourceTypes(ctx context.Context) (models.OperationResult, error) {
	logger.LogInfo("Fetching all roles")

	// Fetch roles from permit system
	roleResources, err := r.PC.SendRequest(ctx, "GET", "resources?include_total_count=true", nil)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Error retrieving roles from permit system", err.Error()), nil
	}

	logger.LogInfo("Successfully fetched roles from permit system", "count", len(roleResources))

	// // Map role data to Role struct
	// roles, err := MapRolesResponseToStruct(roleResources)
	// if err != nil {
	// 	return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to map tenant resources to struct", err.Error()), nil

	// }

	// // Format and return success response
	// successResponse, err := utils.FormatSuccess(roles)
	// if err != nil {
	// 	return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to format success response", err.Error()), nil
	// }
	return nil, nil
}
