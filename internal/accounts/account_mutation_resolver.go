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
	"strings"
	"time"

	"github.com/google/uuid"
)

// AccountMutationResolver handles GraphQL mutations for account-related operations using GORM DB and Permit client
type AccountMutationResolver struct {
	PC  permit.PermitService
	PSC *permit.PermitSdkService
}

// CreateAccount creates a new account based on the provided input.
// It fetches the user ID and tenant ID from the context, validates the input,
// creates the account resource and metadata, and stores them in the database.
// Parameters:
// - ctx: The context for the request, containing authentication and other request-scoped values.
// - input: The input data for creating the account, including account details and billing information.
// Returns:
// - models.OperationResult: The result of the operation, including the created account data.
// - error: An error if the operation fails.
func (r *AccountMutationResolver) CreateAccount(ctx context.Context, input models.CreateAccountInput) (models.OperationResult, error) {
	logger.LogInfo("Started the create account operation")
	// Get the user ID and tenant ID from the context, which are required for account creation
	tenantID, err := helpers.GetTenantID(ctx)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to get tenant ID", err.Error()), nil
	}

	//Map and validate the account creation input data
	_, err = ValidateCreateAccountInput(input)
	if err != nil {
		logger.LogError("error occurred during input mapping and validation", "error", err)
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to do input mapping and validation", err.Error()), nil
	}

	// Prepare metadata for the account from the input data
	metadata, err := r.prepareMetadata(ctx, input)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to prepare metadata in create account", err.Error()), nil
	}

	// Create the resource instances in the permit system with provided metadata
	err = r.createResourceInstances(ctx, input, tenantID, metadata)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to create the resource instances", err.Error()), nil
	}

	// Create relationship tuples in the permit system to establish parent-child relationships
	err = r.createRelationshipTuples(ctx, input, tenantID)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to create the relationship tuples", err.Error()), nil
	}

	// format the success response
	return r.formatSuccessResponse(ctx, input.ID)
}

// UpdateAccount processes an account update operation by:
// 1. Retrieving the existing account data
// 2. Merging it with the provided updates
// 3. Updating the record in both Permit system and database
//
// Parameters:
//   - ctx: Context for the operation
//   - input: UpdateAccountInput containing the account updates
//
// Returns:
//   - models.OperationResult: Result of the operation
//   - error: Any error that occurred during the process
//
// The function ensures atomicity by validating existing data before applying updates
// and maintains consistency across both Permit system and local database.
// Returns a formatted success response with empty data array if successful,
// or a formatted error if any step fails.
func (r *AccountMutationResolver) UpdateAccount(ctx context.Context, input models.UpdateAccountInput) (models.OperationResult, error) {
	_, err := middlewares.AuthorizationMiddleware(ctx, r.PSC, "update", config.AccountResourceTypeID, input.ID.String())

	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "User is not authorized to update the account", err.Error()), nil
	}
	logger.LogInfo("Started the update account operation")

	// Fetch existing account data
	existingAccount, err := r.getExistingAccount(ctx, input.ID)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to get existing account data", err.Error()), nil
	}

	// Merge existing data with updates
	updatedMetadata, err := r.mergeAccountData(ctx, existingAccount, input)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to prepare metadata in update account", err.Error()), nil
	}

	// Update in Permit system
	if err := r.updatePermitResource(ctx, input.ID, updatedMetadata); err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to update account in Permit", err.Error()), nil
	}

	// format the success response
	return r.formatSuccessResponse(ctx, input.ID)
}

// DeleteAccount performs a deletion operation for an account.
// It removes the account from both the external resource service and updates related database records.
//
// The function performs the following operations:
// 1. Validates the user context and retrieves the user ID
// 2. Sends a DELETE request to the resource service
// 3. Updates the account metadata and resource records in the database by setting their RowStatus to 0 (soft delete)
//
// Parameters:
//   - ctx: The context.Context for the request
//   - input: models.DeleteInput containing the ID of the account to be deleted
//
// Returns:
//   - models.OperationResult: Contains the operation result, either success or error details
//   - error: Any error that occurred during the operation
//
// Note: This function performs a soft delete by updating the RowStatus rather than removing records
func (r *AccountMutationResolver) DeleteAccount(ctx context.Context, input models.DeleteInput) (models.OperationResult, error) {
	_, err := middlewares.AuthorizationMiddleware(ctx, r.PSC, "delete", config.AccountResourceTypeID, input.ID.String())

	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "User is not authorized to delete the account", err.Error()), nil
	}
	logger.LogInfo("Started the delete account operation")

	// Delete the resource instance from Permit
	_, err = r.PC.SendRequest(ctx, "DELETE", fmt.Sprintf("resource_instances/%s", input.ID), map[string]interface{}{})
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to delete account in Permit", err.Error()), nil
	}

	// Format success response
	response, _ := utils.FormatSuccessResponse([]models.Data{})
	return response, nil
}

// prepareMetadata converts CreateAccountInput into metadata map for account creation
func (r *AccountMutationResolver) prepareMetadata(ctx context.Context, input models.CreateAccountInput) (map[string]interface{}, error) {
	userID, err := helpers.GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	metadata := map[string]interface{}{
		"id":             input.ID,
		"parentId":       input.ParentID,
		"relationType":   input.RelationType,
		"status":         "ACTIVE",
		"tenantId":       input.TenantID,
		"accountOwnerId": input.AccountOwnerID,
		"name":           input.Name,
		"description":    input.Description,
		"createdBy":      userID,
		"updatedBy":      userID,
		"createdAt":      time.Now().UTC().Format(time.RFC3339),
		"updatedAt":      time.Now().UTC().Format(time.RFC3339),
	}
	if input.Tags != nil {
		metadata["tags"] = input.Tags
	}

	// Only add billing info if it exists
	if input.BillingInfo != nil {
		billingInfo := map[string]interface{}{
			"creditCardNumber": input.BillingInfo.CreditCardNumber,
			"creditCardType":   input.BillingInfo.CreditCardType,
			"expirationDate":   input.BillingInfo.ExpirationDate,
			"cvv":              input.BillingInfo.Cvv,
		}

		// Add billing address if it exists
		if input.BillingInfo.BillingAddress != nil {
			billingInfo["billingAddress"] = map[string]interface{}{
				"street":  input.BillingInfo.BillingAddress.Street,
				"city":    input.BillingInfo.BillingAddress.City,
				"state":   input.BillingInfo.BillingAddress.State,
				"country": input.BillingInfo.BillingAddress.Country,
				"zipcode": input.BillingInfo.BillingAddress.Zipcode,
			}
		}

		metadata["billingInfo"] = billingInfo
	}

	return metadata, nil
}

// createResourceInstances creates resource instances for a given account by making a POST request to the resource_instances endpoint
func (r *AccountMutationResolver) createResourceInstances(ctx context.Context, input models.CreateAccountInput, tenantID *uuid.UUID, metadata map[string]interface{}) error {
	AccountResourceTypeID, _ := uuid.Parse(config.AccountResourceTypeID)
	_, err := r.PC.SendRequest(ctx, "POST", "resource_instances", map[string]interface{}{
		"key":        input.ID,
		"resource":   AccountResourceTypeID.String(),
		"tenant":     tenantID.String(),
		"attributes": metadata,
	})

	if err != nil {
		logger.LogError("error occurred when creating the resource_instances", "error", err)
		return err
	}

	return nil
}

// createRelationshipTuples creates a parent-child relationship tuple for the account in the permission system, linking the account ID with its resource type
func (r *AccountMutationResolver) createRelationshipTuples(ctx context.Context, input models.CreateAccountInput, tenantID *uuid.UUID) error {
	_, err := r.PC.SendRequest(ctx, "POST", "relationship_tuples", map[string]interface{}{
		"object":   config.ClientOrgUnitResourceTypeID + ":" + input.ParentID.String(),
		"relation": strings.ToLower(string(input.RelationType)),
		"subject":  config.AccountResourceTypeID + ":" + input.ID.String(),
		"tenant":   tenantID.String(),
	})

	if err != nil {
		logger.LogError("error occurred when creating the relationship tuple", "error", err)
		return err
	}

	return nil
}

// formatSuccessResponse creates an operation result from account creation input and formats it as a success response
func (r *AccountMutationResolver) formatSuccessResponse(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	account, err := r.getAccountById(ctx, id)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to get the account details by id", err.Error()), nil
	}
	return account, nil
}

// getExistingAccount fetches an existing account's attributes from Permit using the provided ID and returns them as a map
func (r *AccountMutationResolver) getExistingAccount(ctx context.Context, id uuid.UUID) (map[string]interface{}, error) {
	resourceURL := fmt.Sprintf("resource_instances/%s", id)
	accountResource, err := r.PC.SendRequest(ctx, "GET", resourceURL, nil)
	if err != nil {
		logger.LogError("Failed to fetch account from Permit", "error", err)
		return nil, fmt.Errorf("failed to fetch account: %w", err)
	}

	attributes, err := helpers.GetMap(accountResource, "attributes")
	if err != nil {
		logger.LogError("Invalid account data structure", "error", err)
		return nil, fmt.Errorf("invalid account data structure: %s", err)
	}
	return attributes, nil
}

// mergeAccountData combines existing account data with new input data, updating only provided fields
func (r *AccountMutationResolver) mergeAccountData(ctx context.Context, existing map[string]interface{}, input models.UpdateAccountInput) (map[string]interface{}, error) {
	userID, err := helpers.GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	updated := make(map[string]interface{})
	for k, v := range existing {
		updated[k] = v
	}
	updated["updatedBy"] = userID
	updated["updatedAt"] = time.Now().UTC().Format(time.RFC3339)

	if input.Name != nil {
		updated["name"] = *input.Name
	}
	if input.ParentID != nil {
		updated["parentId"] = input.ParentID
	}
	if input.TenantID != nil {
		updated["tenantId"] = input.TenantID
	}
	if input.Description != nil {
		updated["description"] = *input.Description
	}
	if input.Tags != nil {
		updated["tags"] = input.Tags
	}
	if input.RelationType != nil && *input.RelationType != models.RelationTypeEnum("") {
		updated["relationType"] = input.RelationType
	}
	if input.AccountOwnerID != nil && *input.AccountOwnerID != uuid.Nil {
		updated["accountOwnerId"] = input.AccountOwnerID
	}

	if input.BillingInfo != nil {
		billingInfo, _ := existing["billingInfo"].(map[string]interface{})
		updated["billingInfo"] = r.mergeBillingInfo(
			billingInfo,
			input.BillingInfo,
		)
	}

	return updated, nil
}

// mergeBillingInfo merges billing information updates into existing billing info, returning a new map with updated values while preserving unchanged fields
func (r *AccountMutationResolver) mergeBillingInfo(existing map[string]interface{}, updates *models.UpdateBillingInfoInput) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range existing {
		result[k] = v
	}

	if updates.CreditCardNumber != nil {
		result["creditCardNumber"] = *updates.CreditCardNumber
	}
	if updates.CreditCardType != nil {
		result["creditCardType"] = *updates.CreditCardType
	}
	if updates.ExpirationDate != nil {
		result["expirationDate"] = *updates.ExpirationDate
	}
	if updates.Cvv != nil {
		result["cvv"] = *updates.Cvv
	}

	if updates.BillingAddress != nil {
		billingAddress, _ := existing["billingAddress"].(map[string]interface{})
		result["billingAddress"] = r.mergeBillingAddress(
			billingAddress,
			updates.BillingAddress,
		)
	}

	return result
}

// mergeBillingAddress merges updates into existing billing address fields and returns the updated map of address fields
func (r *AccountMutationResolver) mergeBillingAddress(existing map[string]interface{}, updates *models.UpdateBillingAddressInput) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range existing {
		result[k] = v
	}

	if updates.Street != nil {
		result["street"] = *updates.Street
	}
	if updates.City != nil {
		result["city"] = *updates.City
	}
	if updates.State != nil {
		result["state"] = *updates.State
	}
	if updates.Country != nil {
		result["country"] = *updates.Country
	}
	if updates.Zipcode != nil {
		result["zipcode"] = *updates.Zipcode
	}

	return result
}

// updatePermitResource updates metadata attributes for a resource instance in Permit with the given ID
func (r *AccountMutationResolver) updatePermitResource(ctx context.Context, id uuid.UUID, metadata map[string]interface{}) error {
	resourceURL := fmt.Sprintf("resource_instances/%s", id)
	_, err := r.PC.SendRequest(ctx, "PATCH", resourceURL, map[string]interface{}{
		"attributes": metadata,
	})
	if err != nil {
		logger.LogError("Failed to update resource in Permit", "error", err)
		return fmt.Errorf("failed to update permit resource: %w", err)
	}
	return nil
}

// getAccountById fetches a newly created account by its ID using the account query resolver.
// It returns the account data as an OperationResult or an error if the fetch fails.
func (r *AccountMutationResolver) getAccountById(ctx context.Context, accountID uuid.UUID) (models.OperationResult, error) {
	accountResolver := &AccountQueryResolver{PC: r.PC}
	data, err := accountResolver.Account(ctx, accountID)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to get the account details by id from permit", err.Error()), nil
	}
	return data, nil
}
