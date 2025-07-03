package resourcetypes

import (
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/pkg/logger"
	"log"

	"github.com/google/uuid"
)

// Error definitions
type ResourceError string

func (e ResourceError) Error() string { return string(e) }

const (
	ErrInvalidDataField      = ResourceError("missing or invalid data field")
	ErrInvalidResourceFormat = ResourceError("invalid resource data format")
	ErrInvalidActionFormat   = ResourceError("invalid action data format")
	ErrResourceNotFound      = ResourceError("resource not found")
)

// RoleResponse represents the structure of the raw response
type ResourceResponse struct {
	Data []map[string]interface{} `json:"data"`
}

// MapResourceResponseToStruct processes multiple resources from the response
func MapResourceResponseToStruct(resourcesResponse map[string]interface{}) ([]models.Data, error) {
	rawData, ok := extractDataFromResponse(resourcesResponse)
	if !ok {
		return nil, ErrInvalidDataField
	}

	var roles []models.Data
	for _, item := range rawData {
		processedResources, err := processResourceItem(item)
		if err != nil {
			return nil, err
		}
		roles = append(roles, processedResources...)
	}
	return roles, nil
}

// Helper functions

func extractDataFromResponse(response map[string]interface{}) ([]interface{}, bool) {
	// Convert the raw response to our structured type
	var ResourceResponse ResourceResponse
	if data, ok := response["data"]; ok {
		if dataSlice, ok := data.([]interface{}); ok {
			ResourceResponse.Data = make([]map[string]interface{}, len(dataSlice))
			for i, item := range dataSlice {
				if mapItem, ok := item.(map[string]interface{}); ok {
					ResourceResponse.Data[i] = mapItem
				}
			}
			return dataSlice, true
		}
	}
	logger.LogError("invalid data field in response")
	return nil, false
}

func processResourceItem(item interface{}) ([]models.Data, error) {
	resourceData, ok := item.(map[string]interface{})
	if !ok {
		return nil, ErrInvalidResourceFormat
	}

	resourceType, err := MapToResourceTypeData(resourceData)
	if err != nil {
		return nil, fmt.Errorf("invalid resource type: %w", err)
	}

	log.Println(resourceType)

	// Return the processed resource as a slice
	return []models.Data{*resourceType}, nil
}

// MapToResourceTypeData processes resources types for a role
func MapToResourceTypeData(resourceData map[string]interface{}) (*models.ResourceType, error) {
	id, err := helpers.GetUUID(resourceData, "key")
	if err != nil {
		return nil, fmt.Errorf("invalid resource type ID: %w", err)
	}
	actionsData, err := helpers.GetMap(resourceData, "actions")
	if err != nil {
		logger.LogError("failed to get actions map", "error", err)
		return nil, fmt.Errorf("failed to extract actions: %w", err)
	}
	actions := MapToResourceActions(id, actionsData)
	return &models.ResourceType{
		ID:          id,
		Name:        helpers.GetString(resourceData, "name"),
		CreatedAt:   helpers.GetString(resourceData, "createdAt"),
		UpdatedAt:   helpers.GetString(resourceData, "updatedAt"),
		Permissions: actions,
	}, nil
}

func MapToResourceActions(resourceID uuid.UUID, actionsData map[string]interface{}) []*models.Permission {
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
