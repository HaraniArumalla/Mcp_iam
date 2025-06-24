package tenants

import (
	"context"
	"errors"
	"fmt"
	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/utils"
	"iam_services_main_v1/pkg/logger"
	"net/http"

	"github.com/google/uuid"
)

type TenantMutationResolver struct {
	PC  permit.PermitService
	PSC *permit.PermitSdkService
}

// CreateTenant resolver for adding a new Tenant
func (t *TenantMutationResolver) CreateTenant(ctx context.Context, input models.CreateTenantInput) (models.OperationResult, error) {
	//Map and validate the tenant creation input data
	_, err := ValidateCreateTenantInput(input)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Invalid input data", err.Error()), nil
	}

	// Get user and tenant context
	userID, tenantID, err := helpers.GetUserAndTenantID(ctx)
	logger.LogInfo("User ID and Tenant ID", "userID", userID, "tenantID", tenantID)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "User ID & Tenant ID not found in context", err.Error()), nil
	}

	// Check permission
	_, err = t.PSC.Check(ctx, userID.String(), "createTenant", config.TenantResourceTypeID, input.ID.String(), tenantID.String())
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "User is not authorized to update the tenant", err.Error()), nil
	}

	// Prepare metadata for the tenant from the input data
	metadata, err := t.prepareMetadata(ctx, input)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to prepare metadata in create tenant", err.Error()), nil
	}

	if _, err = t.PC.SendRequest(ctx, "POST", "tenants", map[string]interface{}{
		"name":       input.Name,
		"key":        input.ID,
		"attributes": metadata,
	}); err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to create tenant in permit system", err.Error()), nil
	}

	// Create resource instance
	if _, err = t.PC.SendRequest(ctx, "POST", "resource_instances", map[string]interface{}{
		"key":        input.ID,
		"resource":   config.TenantResourceTypeID,
		"tenant":     input.ID,
		"attributes": metadata,
	}); err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to create resource instance in permit system", err.Error()), nil
	}

	// Fetch and return created tenant
	return t.getCreatedTenant(ctx, input.ID)
}

// UpdateTenant resolver for updating a Tenant
func (t *TenantMutationResolver) UpdateTenant(ctx context.Context, input models.UpdateTenantInput) (models.OperationResult, error) {
	if input.ID == uuid.Nil {
		err := errors.New("Tenant ID is required")
		return utils.FormatErrorResponse(http.StatusBadRequest, "Tenant ID is required", err.Error()), nil
	}

	// Fetch existing tenant data
	existingTenant, err := t.getExistingTenant(ctx, input.ID)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to get existing tenant data", err.Error()), nil
	}

	// Merge existing data with updates
	updatedMetadata, err := t.mergeTenantData(ctx, existingTenant, input)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to prepare metadata in update tenant", err.Error()), nil
	}

	// Update tenant in permit
	if _, err := t.PC.SendRequest(ctx, "PATCH", fmt.Sprintf("tenants/%s", input.ID), map[string]interface{}{
		"name":       input.Name,
		"attributes": updatedMetadata,
	}); err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to update tenant in permit system", err.Error()), nil
	}

	// Fetch and return created tenant
	return t.getCreatedTenant(ctx, input.ID)
}

// DeleteTenant resolver for deleting a Tenant
func (t *TenantMutationResolver) DeleteTenant(ctx context.Context, input models.DeleteInput) (models.OperationResult, error) {
	// Check if ID is provided
	if input.ID == uuid.Nil {
		err := errors.New("Tenant ID is required")
		return utils.FormatErrorResponse(http.StatusBadRequest, "Tenant ID is required", err.Error()), nil
	}
	// Get user and tenant context
	userID, tenantID, err := helpers.GetUserAndTenantID(ctx)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "User ID & Tenant ID not found in context", err.Error()), nil
	}

	// Check permission
	_, err = t.PSC.Check(ctx, userID.String(), "updateTenant", config.TenantResourceTypeID, input.ID.String(), tenantID.String())
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "User is not authorized to delete the tenant", err.Error()), nil
	}
	// Delete from permit
	if _, err := t.PC.SendRequest(ctx, "DELETE", fmt.Sprintf("tenants/%s", input.ID), nil); err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to delete tenant from permit", err.Error()), nil
	}

	return utils.FormatSuccess([]models.Data{})
}

// prepareMetadata converts CreateTenantInput into metadata map for tenant creation
func (t *TenantMutationResolver) prepareMetadata(ctx context.Context, input models.CreateTenantInput) (map[string]interface{}, error) {
	userID, err := helpers.GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	metadata := map[string]interface{}{
		"id":             input.ID,
		"type":           config.Tenant,
		"name":           input.Name,
		"description":    input.Description,
		"createdBy":      userID,
		"updatedBy":      userID,
		"accountOwnerId": input.AccountOwnerID,
		"status":         "ACTIVE", // Set default status to ACTIVE for new tenants
		"contactInfo": map[string]interface{}{
			"email":       input.ContactInfo.Email,
			"phoneNumber": input.ContactInfo.PhoneNumber,
			"address": map[string]interface{}{
				"street":  input.ContactInfo.Address.Street,
				"city":    input.ContactInfo.Address.City,
				"state":   input.ContactInfo.Address.State,
				"country": input.ContactInfo.Address.Country,
				"zipcode": input.ContactInfo.Address.Zipcode,
			},
		},
		"tags": input.Tags,
	}

	return metadata, nil
}

func (t *TenantMutationResolver) getCreatedTenant(ctx context.Context, tenantID uuid.UUID) (models.OperationResult, error) {
	tenantResolver := &TenantQueryResolver{PC: t.PC}
	data, err := tenantResolver.Tenant(ctx, tenantID)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to fetch created tenant", err.Error()), nil
	}
	return data, nil
}

func (t *TenantMutationResolver) getExistingTenant(ctx context.Context, id uuid.UUID) (map[string]interface{}, error) {
	resourceURL := fmt.Sprintf("tenants/%s", id)
	tenantResource, err := t.PC.SendRequest(ctx, "GET", resourceURL, nil)
	if err != nil {
		logger.LogError("Failed to fetch tenant from Permit", "error", err)
		return nil, fmt.Errorf("failed to fetch tenant: %w", err)
	}

	attributes, err := helpers.GetMap(tenantResource, "attributes")
	if err != nil {
		logger.LogError("Invalid tenant data structure", "error", err)
		return nil, fmt.Errorf("invalid tenant data structure: %s", err)
	}
	return attributes, nil
}

func (t *TenantMutationResolver) mergeTenantData(ctx context.Context, existing map[string]interface{}, input models.UpdateTenantInput) (map[string]interface{}, error) {
	userID, err := helpers.GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	updated := make(map[string]interface{})
	for k, v := range existing {
		updated[k] = v
	}
	updated["updated_by"] = userID

	if input.Name != nil {
		updated["name"] = *input.Name
	}
	if input.Description != nil {
		updated["description"] = *input.Description
	}
	if input.Status != nil {
		updated["status"] = string(*input.Status)
	}

	if input.ContactInfo != nil {
		contactInfo, _ := existing["contactInfo"].(map[string]interface{})
		updated["contactInfo"] = t.mergeContactInfo(
			contactInfo,
			input.ContactInfo,
		)
	}

	return updated, nil
}

func (t *TenantMutationResolver) mergeContactInfo(existing map[string]interface{}, updates *models.ContactInfoInput) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range existing {
		result[k] = v
	}

	if updates.Email != nil {
		result["email"] = *updates.Email
	}
	if updates.PhoneNumber != nil {
		result["phoneNumber"] = *updates.PhoneNumber
	}

	if updates.Address != nil {
		address, _ := existing["address"].(map[string]interface{})
		result["address"] = t.mergeAddress(
			address,
			updates.Address,
		)
	}

	return result
}

func (t *TenantMutationResolver) mergeAddress(existing map[string]interface{}, updates *models.AddressInput) map[string]interface{} {
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
