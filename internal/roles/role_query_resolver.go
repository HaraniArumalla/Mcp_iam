package roles

import (
	"context"
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/utils"
	"iam_services_main_v1/pkg/logger"
	"net/http"

	"github.com/google/uuid"
)

// RoleQueryResolver handles role-related queries.
type RoleQueryResolver struct {
	PC permit.PermitService
}

// Role retrieves a role by its unique identifier from the permit system.
//
// Parameters:
//   - ctx: The context.Context for the request
//   - id: UUID of the role to retrieve
//
// Returns:
//   - models.OperationResult: Contains either the retrieved role data or error details
//   - error: Any error that occurred during the operation
//
// The function performs the following steps:
//  1. Validates the provided role ID
//  2. Fetches the role data from the permit system
//  3. Maps the retrieved data to a Role struct
//  4. Formats and returns the response
//
// If any error occurs during these steps, it will be logged and returned in the OperationResult.
func (r *RoleQueryResolver) Role(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	logger.LogInfo("Fetching role by ID", "id", id)
	// Validate ID
	if id == uuid.Nil {
		err := fmt.Errorf("invalid role ID: %s", id)
		return utils.FormatErrorResponse(http.StatusBadRequest, "invalid role ID", err.Error()), nil
	}

	// Fetch role from permit system
	data, err := r.PC.SendRequest(ctx, "GET", "resources?include_total_count=true", nil)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Error retrieving roles from permit system", err.Error()), nil
	}

	// Map role data to Role struct
	role, err := MapRoleResponseToStruct(data, id)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Error mapping role data", err.Error()), nil
	}

	// Format and return success response
	successResponse, err := utils.FormatSuccess(role)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to format success response", err.Error()), nil

	}
	return successResponse, nil

}

// Roles retrieves all roles from the permit system.
// It makes a GET request to fetch role resources, maps the response to Role structs,
// and returns the result in a standardized OperationResult format.
//
// The function performs the following steps:
// 1. Fetches role data from the permit system
// 2. Maps the role response data to internal Role structs
// 3. Formats the response as a success operation
//
// Returns:
//   - models.OperationResult: Contains the roles data or error information
//   - error: Returns nil unless there's a critical error preventing operation
func (r *RoleQueryResolver) Roles(ctx context.Context) (models.OperationResult, error) {
	logger.LogInfo("Fetching all roles")

	// Fetch roles from permit system
	roleResources, err := r.PC.SendRequest(ctx, "GET", "resources?include_total_count=true", nil)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Error retrieving roles from permit system", err.Error()), nil
	}

	// Map role data to Role struct
	roles, err := MapRolesResponseToStruct(roleResources)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to map tenant resources to struct", err.Error()), nil

	}

	// Format and return success response
	successResponse, err := utils.FormatSuccess(roles)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to format success response", err.Error()), nil
	}
	return successResponse, nil
}
