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
	"iam_services_main_v1/pkg/logger"
	"log"
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
	return r.getPermissionDetailsById(ctx, input.AssignableScopeRef)
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

func (r *PermissionMutationResolver) getPermissionDetailsById(ctx context.Context, resourceID uuid.UUID) (models.OperationResult, error) {
	logger.LogInfo("Fetching all roles")

	// Fetch roles from permit system
	resourceURL := fmt.Sprintf("resources/%s", resourceID)
	resourceResources, err := r.PC.SendRequest(ctx, "GET", resourceURL, nil)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Error retrieving roles from permit system", err.Error()), nil
	}

	// Map resource data to resource struct
	resources, err := MapToPermissionData(resourceResources)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to map tenant resources to struct", err.Error()), nil

	}

	// Format and return success response
	successResponse, err := utils.FormatSuccess([]models.Data{resources})
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to format success response", err.Error()), nil
	}
	return successResponse, nil
}

// MapToPermissionData processes resources types for a role
func MapToPermissionData(resourceData map[string]interface{}) (*models.ResourceType, error) {
	id, err := helpers.GetUUID(resourceData, "key")
	if err != nil {
		return nil, fmt.Errorf("invalid resource type ID: %w", err)
	}
	actionsData, err := helpers.GetMap(resourceData, "actions")
	if err != nil {
		logger.LogError("failed to get actions map", "error", err)
		return nil, fmt.Errorf("failed to extract actions: %w", err)
	}
	actions := MapToPermissionActions(id, actionsData)
	log.Printf("actions: %v", actions)
	return &models.ResourceType{
		ID:          id,
		Name:        helpers.GetString(resourceData, "name"),
		CreatedAt:   helpers.GetString(resourceData, "createdAt"),
		UpdatedAt:   helpers.GetString(resourceData, "updatedAt"),
		Permissions: actions,
	}, nil
}

func MapToPermissionActions(resourceID uuid.UUID, actionsData map[string]interface{}) []*models.Permission {
	var actions []*models.Permission
	for key, action := range actionsData {
		actionData, ok := action.(map[string]interface{})
		if !ok {
			continue
		}
		id, err := helpers.GetUUID(actionData, "id")
		if err != nil {
			logger.LogError("failed to get action ID", "error", err)
			continue
		}
		actions = append(actions, &models.Permission{
			ID:          id,
			Name:        key,
			Description: helpers.StringPtr(helpers.GetString(actionData, "description")),
			CreatedAt:   helpers.GetString(actionData, "createdAt"),
			UpdatedAt:   helpers.GetString(actionData, "updatedAt"),
		})
	}
	return actions
}
