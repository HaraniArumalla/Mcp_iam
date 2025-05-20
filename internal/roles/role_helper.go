package roles

import (
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/internal/tags"
	"iam_services_main_v1/pkg/logger"
	"log"

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
		processedRoles, err := processResourceItem(item)
		if err != nil {
			return nil, err
		}
		roles = append(roles, processedRoles...)
	}
	return roles, nil
}

// MapRoleResponseToStruct processes a single role by ID
func MapRoleResponseToStruct(resourcesResponse map[string]interface{}, id uuid.UUID) ([]models.Data, error) {
	rawData, ok := extractDataFromResponse(resourcesResponse)
	if !ok {
		return nil, ErrInvalidDataField
	}

	for _, item := range rawData {
		resourceData, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		rolesData, actionsData, err := extractRolesAndActions(resourceData)
		if err != nil {
			continue
		}

		roleData, err := findRoleByID(rolesData, actionsData, resourceData, id)
		if err != nil {
			continue
		}
		return roleData, nil
	}
	return nil, fmt.Errorf("%w: %s", ErrRoleNotFound, id)
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

	rolesData, actionsData, err := extractRolesAndActions(resourceData)
	if err != nil {
		return nil, err
	}

	return processRoles(rolesData, actionsData, resourceData)
}

func processRoles(rolesData, actionsData, resourceData map[string]interface{}) ([]models.Data, error) {
	var roles []models.Data
	for _, role := range rolesData {
		roleData, ok := role.(map[string]interface{})
		if !ok {
			logger.LogError(string(ErrInvalidRoleFormat))
			continue
		}

		roleModel, err := MapToRoleData(roleData, actionsData, resourceData)
		if err != nil {
			continue
		}
		roles = append(roles, roleModel)
	}
	return roles, nil
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

// MapToPermissionsData processes permissions for a role
func MapToPermissionsData(roleData, actionsData map[string]interface{}) ([]*models.Permission, error) {
	permissionsData, err := helpers.GetSlice(roleData, "permissions")
	if err != nil {
		return nil, err
	}

	var permissions []*models.Permission
	for _, perm := range permissionsData {
		permission, err := processPermission(perm, actionsData)
		if err != nil {
			logger.LogError("permission processing error", "error", err)
			continue
		}
		if permission != nil {
			permissions = append(permissions, permission)
		}
	}
	return permissions, nil
}

func processPermission(permKey interface{}, actionsData map[string]interface{}) (*models.Permission, error) {
	permissionStr, ok := permKey.(string)
	if !ok {
		return nil, ErrInvalidActionFormat
	}

	actionData, ok := actionsData[permissionStr].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	id, err := helpers.GetUUID(actionData, "id")
	if err != nil {
		return nil, err
	}

	return &models.Permission{
		ID:   id,
		Name: helpers.GetString(actionData, "name"),
	}, nil
}

// extractRolesAndActions extracts roles and actions data from the resource data
func extractRolesAndActions(resourceData map[string]interface{}) (map[string]interface{}, map[string]interface{}, error) {
	rolesData, err := helpers.GetMap(resourceData, "roles")
	log.Printf("rolesData: %v", rolesData)
	if err != nil {
		logger.LogError("failed to get roles map", "error", err)
		return nil, nil, fmt.Errorf("failed to extract roles: %w", err)
	}

	actionsData, err := helpers.GetMap(resourceData, "actions")
	if err != nil {
		logger.LogError("failed to get actions map", "error", err)
		return nil, nil, fmt.Errorf("failed to extract actions: %w", err)
	}

	return rolesData, actionsData, nil
}

// findRoleByID searches for and returns a specific role by ID
func findRoleByID(rolesData map[string]interface{}, actionsData, resourceData map[string]interface{}, id uuid.UUID) ([]models.Data, error) {
	for _, role := range rolesData {
		roleData, ok := role.(map[string]interface{})
		if !ok {
			continue
		}

		roleID, err := helpers.GetUUID(roleData, "key")
		log.Printf("roleID: %v", roleID)
		if err != nil || roleID != id {
			continue
		}

		roleModel, err := MapToRoleData(roleData, actionsData, resourceData)
		if err != nil {
			return nil, fmt.Errorf("failed to map role data: %w", err)
		}
		return []models.Data{roleModel}, nil
	}
	return nil, fmt.Errorf("%w: %s", ErrRoleNotFound, id)
}
