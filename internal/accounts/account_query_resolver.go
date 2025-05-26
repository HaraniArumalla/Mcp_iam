package accounts

import (
	"context"
	"fmt"
	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/internal/middlewares"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/utils"
	"iam_services_main_v1/pkg/logger"
	"net/http"

	"github.com/google/uuid"
)

// AccountQueryResolver handles database queries and permission checks for account-related operations using GORM and Permit.io client
type AccountQueryResolver struct {
	PC  permit.PermitService
	PSC *permit.PermitSdkService
}

// Accounts retrieves all account resources for a given tenant.
// It performs the following operations:
// 1. Extracts the tenant ID from the context
// 2. Queries the permit service for account resources associated with the tenant
// 3. Maps the response data to account structures
// 4. Returns formatted operation result containing the accounts
//
// Parameters:
//   - ctx: The context.Context for the request
//
// Returns:
//   - models.OperationResult: Contains either the accounts data or error details
//   - error: Any error encountered during processing
func (r *AccountQueryResolver) Accounts(ctx context.Context) (models.OperationResult, error) {
	logger.LogInfo("Fetching all accounts")

	// Get tenant ID from context
	tenantID, err := helpers.GetTenantID(ctx)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to get tenant ID", err.Error()), nil
	}

	// Fetch account resources from permit
	accounts, err := r.fetchAccounts(ctx, tenantID)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to get all accounts from permit", err.Error()), nil
	}

	// Format success response
	response, _ := utils.FormatSuccessResponse(accounts)
	return response, nil

}

// Account retrieves account details by UUID from the resource instance.
//
// Parameters:
//   - ctx: The context.Context for the request
//   - id: UUID of the account to fetch
//
// Returns:
//   - models.OperationResult: Contains either the account details on success or error details on failure
//   - error: Returns nil as errors are wrapped in OperationResult
//
// The function performs the following operations:
//  1. Sends GET request to fetch account resources using the provided ID
//  2. Maps the response to account structure
//  3. Formats the response as a success operation result
//
// Any errors during these operations are logged and returned as formatted error responses.
func (r *AccountQueryResolver) Account(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	_, err := middlewares.AuthorizationMiddleware(ctx, r.PSC, "getbyid", config.AccountResourceTypeID, id.String())

	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "User is not authorized to account", err.Error()), nil
	}
	logger.LogInfo("Fetching account by ID", "id", id)

	// Get tenant ID from context
	_, err = helpers.GetTenantID(ctx)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to get tenant ID", err.Error()), nil
	}

	// Fetch account resources from permit
	account, err := r.fetchAccount(ctx, id)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to get account resources from permit", err.Error()), nil
	}

	// Format success response
	response, _ := utils.FormatSuccessResponse(account)
	return response, nil
}

// fetchAccounts retrieves all account resources for a given tenant from the resource instance.
func (r *AccountQueryResolver) fetchAccounts(ctx context.Context, tenantID *uuid.UUID) ([]models.Data, error) {
	// Fetch all account resources from permit
	resourceURL := fmt.Sprintf("resource_instances/detailed?tenant=%s&resource=%s",
		tenantID.String(), config.AccountResourceTypeID)
	accountResources, err := r.PC.SendRequest(ctx, "GET", resourceURL, nil)
	if err != nil {
		logger.LogError("Failed to get account resources from permit", "error", err)
		return nil, err
	}

	// Map account resources to struct
	accounts, err := MapAccountsResponseToStruct(accountResources)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// fetchAccount retrieves account details by UUID from the resource instance.
func (r *AccountQueryResolver) fetchAccount(ctx context.Context, id uuid.UUID) ([]models.Data, error) {
	// Fetch account resources from permit
	resourceURL := fmt.Sprintf("resource_instances/%s", id)
	accountResource, err := r.PC.SendRequest(ctx, "GET", resourceURL, nil)
	if err != nil {
		return nil, err
	}

	// Map account resources to struct
	account, err := MapAccountResponseToStruct(accountResource)
	if err != nil {
		return nil, err
	}
	return account, nil
}
