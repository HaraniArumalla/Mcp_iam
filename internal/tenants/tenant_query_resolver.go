package tenants

import (
	"context"
	"fmt"
	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/middlewares"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/utils"
	"iam_services_main_v1/pkg/logger"
	"net/http"

	"github.com/google/uuid"
)

// TenantQueryResolver handles database queries and permission checks for Tenant-related operations using GORM and Permit.io client
type TenantQueryResolver struct {
	PC  permit.PermitService
	PSC *permit.PermitSdkService
}

// Tenants retrieves all Tenant resources for a given tenant.
// It performs the following operations:
// 1. Extracts the tenant ID from the context
// 2. Queries the permit service for Tenant resources associated with the tenant
// 3. Maps the response data to Tenant structures
// 4. Returns formatted operation result containing the Tenants
//
// Parameters:
//   - ctx: The context.Context for the request
//
// Returns:
//   - models.OperationResult: Contains either the Tenants data or error details
//   - error: Any error encountered during processing
func (r *TenantQueryResolver) Tenants(ctx context.Context) (models.OperationResult, error) {
	logger.LogInfo("Fetching all tenants")

	// Fetch tenants from permit system
	resourceURL := "tenants?include_total_count=true"
	tenantResources, err := r.PC.SendRequest(ctx, "GET", resourceURL, nil)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to get tenant resources from permit", err.Error()), nil
	}

	// Map tenant data to Tenant struct
	tenants, err := MapTenantsResponseToStruct(tenantResources)

	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to map tenant resources to struct", err.Error()), nil
	}

	// Format and return success response
	successResponse, err := utils.FormatSuccess(tenants)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to format success response", err.Error()), nil
	}
	return successResponse, nil

}

// Tenant retrieves Tenant details by UUID from the resource instance.
//
// Parameters:
//   - ctx: The context.Context for the request
//   - id: UUID of the Tenant to fetch
//
// Returns:
//   - models.OperationResult: Contains either the Tenant details on success or error details on failure
//   - error: Returns nil as errors are wrapped in OperationResult
//
// The function performs the following operations:
//  1. Sends GET request to fetch Tenant resources using the provided ID
//  2. Maps the response to Tenant structure
//  3. Formats the response as a success operation result
//
// Any errors during these operations are logged and returned as formatted error responses.
func (r *TenantQueryResolver) Tenant(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	_, err := middlewares.AuthorizationMiddleware(ctx, r.PSC, "getbyid", config.TenantResourceTypeID, id.String())

	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "User is not authorized to get the tenant details by id", err.Error()), nil
	}
	logger.LogInfo("Fetching tenant by ID", "id", id)

	// Validate ID
	if id == uuid.Nil {
		err := fmt.Errorf("invalid tenant ID: %s", id)
		return utils.FormatErrorResponse(http.StatusBadRequest, "invalid tenant ID", err.Error()), nil
	}

	// Fetch tenant from permit system
	tenantResource, err := r.FetchTenant(ctx, id)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to get tenant resources from permit", err.Error()), nil
	}

	//	Map tenant data to Tenant struct
	tenant, err := MapTenantResponseToStruct(tenantResource)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to map tenant resources to struct", err.Error()), nil
	}

	// Format and return success response
	successResponse, err := utils.FormatSuccess(tenant)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to format success response", err.Error()), nil
	}
	return successResponse, nil
}

// fetchAccount retrieves account details by UUID from the resource instance.
func (r *TenantQueryResolver) FetchTenant(ctx context.Context, id uuid.UUID) (map[string]interface{}, error) {
	// Fetch account resources from permit
	resourceURL := fmt.Sprintf("tenants/%s", id)
	accountResource, err := r.PC.SendRequest(ctx, "GET", resourceURL, nil)
	if err != nil {
		return nil, err
	}

	return accountResource, nil
}
