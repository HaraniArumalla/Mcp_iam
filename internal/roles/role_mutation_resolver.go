package roles

import (
	"context"
	"errors"
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/utils"
	"iam_services_main_v1/internal/validations"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// RoleMutationResolver handles role-related mutations.
type RoleMutationResolver struct {
	PC permit.PermitService
}

// CreateRole creates a new role in the system.
// It performs the following operations:
// 1. Validates user and tenant context
// 2. Validates the input parameters
// 3. Creates the role in the permit system
// 4. Retrieves and returns the created role
//
// Parameters:
//   - ctx: The context carrying user authentication and request scoping
//   - input: CreateRoleInput containing the role details to be created
//
// Returns:
//   - models.OperationResult: The result of the operation including the created role
//   - error: An error if any step fails, nil otherwise
func (r *RoleMutationResolver) CreateRole(ctx context.Context, input models.CreateRoleInput) (models.OperationResult, error) {
	// Get user and tenant context
	userID, tenantID, err := helpers.GetUserAndTenantID(ctx)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "failed to fetch user and tenant IDs", err.Error()), nil
	}

	// Validate input
	if err := validateCreateRoleInput(input); err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "input validation failed", err.Error()), nil

	}

	// Create role in permit system
	if err := r.createRoleInPermit(ctx, input, userID, tenantID); err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "permit role creation failed", err.Error()), nil

	}

	// Fetch and return created role
	return r.getCreatedRole(ctx, input.ID)
}

// UpdateRole updates an existing role in the system with the provided input.
// It performs the following steps:
// 1. Validates user context and tenant information
// 2. Validates the update role input parameters
// 3. Updates the role in the permit system
// 4. Retrieves and returns the updated role
//
// Parameters:
//   - ctx: The context.Context for the request
//   - input: models.UpdateRoleInput containing the role update information
//
// Returns:
//   - models.OperationResult: The result of the update operation
//   - error: Any error that occurred during the process
//
// Possible errors:
//   - User/tenant context validation errors
//   - Input validation errors
//   - Permit system update errors
func (r *RoleMutationResolver) UpdateRole(ctx context.Context, input models.UpdateRoleInput) (models.OperationResult, error) {
	// Get user and tenant context
	userID, tenantID, err := helpers.GetUserAndTenantID(ctx)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "failed to fetch user and tenant IDs", err.Error()), nil
	}
	// Validate Input
	if err := r.validateUpdateRoleInput(input); err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Error validating update role input", err.Error()), nil
	}

	// Update role in permit system
	if err := r.updateRoleInPermit(ctx, input, userID, tenantID); err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "permit role creation failed", err.Error()), nil
	}

	// Fetch and return created role
	return r.getCreatedRole(ctx, input.ID)
}

// DeleteRole handles the deletion of a role from the permit system.
// It performs validation of the input fields and processes the role deletion in the permit system.
//
// Parameters:
//   - ctx: The context for managing timeouts and cancellation
//   - input: DeleteRoleInput containing the role details to be deleted
//
// Returns:
//   - models.OperationResult: Contains the operation status and any relevant data
//   - error: Returns nil if successful, otherwise returns an error describing what went wrong
func (r *RoleMutationResolver) DeleteRole(ctx context.Context, input models.DeleteRoleInput) (models.OperationResult, error) {
	// Validate required input fields
	if err := validateDeleteRoleInput(input); err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Error validating delete role input", err.Error()), nil
	}

	// Delete role from permit system
	if err := r.deleteRoleFromPermit(ctx, input); err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Error deleting role in permit", err.Error()), nil
	}

	// Return success response
	return utils.FormatSuccess([]models.Data{})
}

// createRoleInPermit function creates a new role in Permit.io service with the provided input parameters and returns any error encountered
func (r *RoleMutationResolver) createRoleInPermit(ctx context.Context, input models.CreateRoleInput, userID, tenantID *uuid.UUID) error {
	metadata, err := r.prepareMetadataForCreateInput(input, userID, tenantID)
	if err != nil {
		return fmt.Errorf("metadata preparation failed: %w", err)
	}

	permitMap := map[string]interface{}{
		"name":        input.Name,
		"description": input.Description,
		"key":         input.ID,
		"attributes":  metadata,
		"permissions": input.Permissions,
	}

	_, err = r.PC.SendRequest(ctx, "POST", fmt.Sprintf("resources/%s/roles", input.AssignableScopeRef), permitMap)
	return err
}

// updateRoleInPermit function updates a role in Permit.io service with the provided input parameters and returns any error encountered
func (r *RoleMutationResolver) updateRoleInPermit(ctx context.Context, input models.UpdateRoleInput, userID, tenantID *uuid.UUID) error {
	metadata, err := r.prepareMetadataForUpdateInput(input, userID, tenantID)
	if err != nil {
		return fmt.Errorf("metadata preparation failed: %w", err)
	}

	permitMap := map[string]interface{}{
		"name":        input.Name,
		"description": input.Description,
		"attributes":  metadata,
		"permissions": input.Permissions,
	}

	_, err = r.PC.SendRequest(ctx, "PATCH", fmt.Sprintf("resources/%s/roles/%s", input.AssignableScopeRef, input.ID), permitMap)
	return err
}

// getCreatedRole fetches a newly created role by its ID using the role query resolver.
// It returns the role data as an OperationResult or an error if the fetch fails.
func (r *RoleMutationResolver) getCreatedRole(ctx context.Context, roleID uuid.UUID) (models.OperationResult, error) {
	roleResolver := &RoleQueryResolver{PC: r.PC}
	data, err := roleResolver.Role(ctx, roleID)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "role fetch failed", err.Error()), nil
	}
	return data, nil
}

// Helper function to validate delete role input
func validateDeleteRoleInput(input models.DeleteRoleInput) error {
	if input.ID == uuid.Nil {
		return errors.New("role ID is required")
	}
	if input.AssignableScopeRef == uuid.Nil {
		return errors.New("assignable scope reference is required")
	}
	return nil
}

// Helper function to delete role from permit system
func (r *RoleMutationResolver) deleteRoleFromPermit(ctx context.Context, input models.DeleteRoleInput) error {
	endpoint := fmt.Sprintf("resources/%s/roles/%s", input.AssignableScopeRef, input.ID)
	_, err := r.PC.SendRequest(ctx, "DELETE", endpoint, nil)
	return err
}

// validateUpdateRoleInput checks if the provided UpdateRoleInput is valid by ensuring required fields are not empty
// and the role name follows naming conventions. Returns an error if validation fails, nil otherwise.
func (r *RoleMutationResolver) validateUpdateRoleInput(input models.UpdateRoleInput) error {
	// Check required fields
	if input.ID == uuid.Nil || input.AssignableScopeRef == uuid.Nil {
		return errors.New("id and assignable scope reference are required")
	}

	// Validate name if provided
	if input.Name != "" {
		if err := validations.ValidateName(input.Name); err != nil {
			return fmt.Errorf("invalid role name: %w", err)
		}
	}

	return nil
}

// validateCreateRoleInput checks if the provided CreateRoleInput is valid by ensuring required fields are not empty
// and the role name follows naming conventions. Returns an error if validation fails, nil otherwise.
func validateCreateRoleInput(input models.CreateRoleInput) error {

	// Validate input
	if input.Name == "" || input.AssignableScopeRef == uuid.Nil || input.RoleType == "" || input.RoleType == "DEFAULT" {
		return errors.New("invalid input recieved")
	}

	if err := validations.ValidateName(input.Name); err != nil {
		return fmt.Errorf("invalid role name: %w", err)
	}

	return nil
}

// prepareMetadata converts CreateAccountInput into metadata map for account creation
func (r *RoleMutationResolver) prepareMetadataForCreateInput(input models.CreateRoleInput, userID *uuid.UUID, tenantID *uuid.UUID) (map[string]interface{}, error) {
	metadata := map[string]interface{}{
		"id":                 input.ID,
		"tenantId":           *tenantID,
		"name":               input.Name,
		"description":        input.Description,
		"assignableScopeRef": input.AssignableScopeRef,
		"version":            input.Version,
		"roleType":           input.RoleType,
		"permissions":        input.Permissions,
		"createdBy":          *userID,
		"updatedBy":          *userID,
		"createdAt":          time.Now().Format(time.RFC3339),
		"updatedAt":          time.Now().Format(time.RFC3339),
	}
	if input.Tags != nil {
		metadata["tags"] = input.Tags
	}

	return metadata, nil
}

// prepareMetadata converts UpdateAccountInput into metadata map for account creation
func (r *RoleMutationResolver) prepareMetadataForUpdateInput(input models.UpdateRoleInput, userID *uuid.UUID, tenantID *uuid.UUID) (map[string]interface{}, error) {
	metadata := map[string]interface{}{
		"id":                 input.ID,
		"tenantId":           *tenantID,
		"name":               input.Name,
		"description":        input.Description,
		"assignableScopeRef": input.AssignableScopeRef,
		"version":            input.Version,
		"roleType":           input.RoleType,
		"permissions":        input.Permissions,
		"createdBy":          *userID,
		"updatedBy":          *userID,
		"createdAt":          time.Now().Format(time.RFC3339),
		"updatedAt":          time.Now().Format(time.RFC3339),
	}
	if input.Tags != nil {
		metadata["tags"] = input.Tags
	}

	return metadata, nil
}
