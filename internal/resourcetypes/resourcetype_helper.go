package resourcetypes

import (
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/internal/tags"
	"iam_services_main_v1/pkg/logger"

	"github.com/google/uuid"
)

// Error definitions
type RoleError string

func (e RoleError) Error() string { return string(e) }

const (
	ErrInvalidDataField    = RoleError("missing or invalid data field")
	ErrInvalidRoleFormat   = RoleError("invalid role data format")
	ErrInvalidActionFormat = RoleError("invalid action data format")
	ErrRoleNotFound        = RoleError("role not found")
)

// RoleResponse represents the structure of the raw response
type RoleResponse struct {
	Data []map[string]interface{} `json:"data"`
}

// MapRolesResponseToStruct processes multiple roles from the response
func MapRolesResponseToStruct(resourcesResponse map[string]interface{}) ([]models.Data, error) {
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
	var roleResponse RoleResponse
	if data, ok := response["data"]; ok {
		if dataSlice, ok := data.([]interface{}); ok {
			roleResponse.Data = make([]map[string]interface{}, len(dataSlice))
			for i, item := range dataSlice {
				if mapItem, ok := item.(map[string]interface{}); ok {
					roleResponse.Data[i] = mapItem
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
		return nil, ErrInvalidRoleFormat
	}

	actionsData, err := extractActions(resourceData)
	if err != nil {
		return nil, err
	}

	return processResources(actionsData, resourceData)
}

// MapToRoleData converts role data to a Role model
func MapToRoleData(roleData, actionsData, resourceData map[string]interface{}) (*models.Role, error) {
	id, err := helpers.GetUUID(roleData, "key")
	if err != nil {
		return nil, fmt.Errorf("invalid role ID: %w", err)
	}

	attributes, err := helpers.GetMap(roleData, "attributes")
	if err != nil {
		return nil, fmt.Errorf("invalid attributes: %w", err)
	}
	createdBy, _ := helpers.GetUUID(attributes, "createdBy")
	updatedBy, _ := helpers.GetUUID(attributes, "updatedBy")

	//assignableScopeRef, _ := helpers.GetUUID(attributes, "assignableScopeRef")

	permissions, err := MapToPermissionsData(roleData, actionsData)
	if err != nil {
		return nil, fmt.Errorf("invalid permissions: %w", err)
	}
	resourceType, err := MapToResourceTypeData(resourceData)
	if err != nil {
		return nil, fmt.Errorf("invalid resource type: %w", err)
	}
	description := helpers.GetString(roleData, "description")
	version := helpers.GetString(attributes, "version")
	roleType := helpers.GetString(attributes, "roleType")
	tags := tags.GetResourceTags(attributes, "tags")
	role := &models.Role{
		ID:              id,
		Name:            helpers.GetString(roleData, "name"),
		Description:     &description,
		Version:         version,
		CreatedAt:       helpers.GetString(attributes, "createdAt"),
		UpdatedAt:       helpers.GetString(attributes, "updatedAt"),
		CreatedBy:       createdBy,
		UpdatedBy:       updatedBy,
		Permissions:     permissions,
		AssignableScope: resourceType,
		Tags:            tags,
	}
	if roleType == "" || roleType == "DEFAULT" {
		role.RoleType = models.RoleTypeEnumDefault
	} else {
		role.RoleType = models.RoleTypeEnumCustom
	}
	return role, nil
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
			ID:              id,
			Name:            key,
			AssignableScope: resourceID,
		})
	}
	return actions
}

// extractActions extracts actions data from the resource data
func extractActions(resourceData map[string]interface{}) (map[string]interface{}, error) {
	actionsData, err := helpers.GetMap(resourceData, "actions")
	if err != nil {
		logger.LogError("failed to get actions map", "error", err)
		return nil, fmt.Errorf("failed to extract actions: %w", err)
	}

	return actionsData, nil
}
