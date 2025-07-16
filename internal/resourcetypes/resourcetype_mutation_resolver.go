package resourcetypes

import (
	"context"
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/utils"
	"iam_services_main_v1/pkg/logger"
	"log"
	"net/http"
	"strings"
)

type ResourceTypeMutationResolver struct {
	PC permit.PermitService
}

func (r *ResourceTypeMutationResolver) CreateResource(ctx context.Context, input models.CreateResourceInput) (models.OperationResult, error) {
	logger.LogInfo("Started the create resource operation")

	// Prepare resource actions for the resource from the input data
	actions, err := r.prepareResourceActions(ctx, input)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to prepare metadata in create account", err.Error()), nil
	}

	// Create the resource instances in the permit system with provided metadata
	err = r.createResource(ctx, input, actions)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to create the resource instances", err.Error()), nil
	}

	return nil, nil
	// format the success response
	//return r.formatSuccessResponse(ctx, input.ID)
}

// prepareResourceActions converts CreateResourceInput into actions map for resource creation
func (r *ResourceTypeMutationResolver) prepareResourceActions(ctx context.Context, resource models.CreateResourceInput) (map[string]interface{}, error) {
	// Example implementation: return an empty map and nil error
	actions := make(map[string]interface{})
	if resource.Permissions != nil {
		for _, permission := range resource.Permissions {
			name := strings.ToLower(permission.Name)
			namekey := strings.ReplaceAll(name, " ", "_")
			actions[namekey] = map[string]interface{}{
				"key":         namekey,
				"action":      permission.Name,
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
