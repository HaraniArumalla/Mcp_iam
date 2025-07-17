package resourcetypes

import (
	"context"
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/utils"
	"iam_services_main_v1/pkg/logger"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type ResourceTypeMutationResolver struct {
	PC permit.PermitService
}

func (r *ResourceTypeMutationResolver) CreateResource(ctx context.Context, input models.CreateResourceInput) (models.OperationResult, error) {
	logger.LogInfo("Started the create resource operation")

	// Prepare resource actions for the resource from the input data
	actions, err := r.prepareResourceActions(input)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to prepare metadata in create account", err.Error()), nil
	}

	// Create the resource instances in the permit system with provided metadata
	err = r.createResource(ctx, input, actions)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to create the resource instances", err.Error()), nil
	}

	// format the success response
	return r.getPermissionDetailsById(ctx, input.ID)
}

// prepareResourceActions converts CreateResourceInput into actions map for resource creation
func (r *ResourceTypeMutationResolver) prepareResourceActions(resource models.CreateResourceInput) (map[string]interface{}, error) {
	// Example implementation: return an empty map and nil error
	actions := make(map[string]interface{})
	if resource.Permissions != nil {
		for _, permission := range resource.Permissions {
			name := strings.ToLower(permission.Name)
			namekey := strings.ReplaceAll(name, " ", "_")
			actions[namekey] = map[string]interface{}{
				"name":        permission.Name,
				"description": permission.Description,
			}
		}
	} else {
		logger.LogError("Permissions are nil in the resource input")
		return nil, fmt.Errorf("permissions cannot be nil")
	}

	log.Println("Prepared actions for resource:", actions)

	return actions, nil
}

// createResource creates resource  for a given account by making a POST request to the resources endpoint
func (r *ResourceTypeMutationResolver) createResource(ctx context.Context, input models.CreateResourceInput, actions map[string]interface{}) error {
	_, err := r.PC.SendRequest(ctx, "POST", "resources", map[string]interface{}{
		"key":     input.ID,
		"name":    input.Name,
		"actions": actions,
	})

	if err != nil {
		logger.LogError("error occurred when creating the resource_instances", "error", err)
		return err
	}

	return nil
}

func (r *ResourceTypeMutationResolver) getPermissionDetailsById(ctx context.Context, resourceID uuid.UUID) (models.OperationResult, error) {
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
