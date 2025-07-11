package permissions

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
	"strings"

	"github.com/google/uuid"
)

type PermissionMutationResolver struct {
	PC permit.PermitService
}

func (r *PermissionMutationResolver) CreatePermission(ctx context.Context, input models.CreatePermissionInput) (models.OperationResult, error) {
	// Get user and tenant context
	userID, tenantID, err := helpers.GetUserAndTenantID(ctx)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "failed to fetch user and tenant IDs", err.Error()), nil
	}

	// Validate input
	if err := validateCreatePermissionInput(input); err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "input validation failed", err.Error()), nil

	}

	// Create role in permit system
	if err := r.createPermissionInPermit(ctx, input, userID, tenantID); err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "permit role creation failed", err.Error()), nil

	}

	// Fetch and return created role
	//return r.getCreatedRole(ctx, input.ID)
}

func (r *PermissionMutationResolver) DeletePermission(ctx context.Context, input models.DeleteInput) (models.OperationResult, error) {

	return nil, nil
}

func (r *PermissionMutationResolver) UpdatePermission(ctx context.Context, input models.UpdatePermissionInput) (models.OperationResult, error) {
	return nil, nil
}

// validateCreatePermissionInput checks if the provided CreatePermissionInput is valid by ensuring required fields are not empty
// and the role name follows naming conventions. Returns an error if validation fails, nil otherwise.
func validateCreatePermissionInput(input models.CreatePermissionInput) error {

	// Validate input
	if input.Name == "" || input.AssignableScopeRef == uuid.Nil {
		return errors.New("invalid input recieved")
	}

	if err := validations.ValidateName(input.Name); err != nil {
		return fmt.Errorf("invalid role name: %w", err)
	}

	return nil
}

// createPermissionInPermit function creates a new role in Permit.io service with the provided input parameters and returns any error encountered
func (r *PermissionMutationResolver) createPermissionInPermit(ctx context.Context, input models.CreatePermissionInput, userID, tenantID *uuid.UUID) error {
	name := strings.ToLower(input.Name)
	namekey := strings.ReplaceAll(name, " ", "_")

	permitMap := map[string]interface{}{
		"description": input.Description,
		"key":         namekey,
		"name":        input.Name,
	}

	_, err := r.PC.SendRequest(ctx, "POST", fmt.Sprintf("resources/%s/actions", input.AssignableScopeRef), permitMap)
	return err
}
