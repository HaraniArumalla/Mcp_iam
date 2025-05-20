package tenants

import (
	"fmt"
	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/internal/tags"
	"iam_services_main_v1/pkg/logger"
)

// MapTenantsResponseToStruct maps a response containing resource data to a slice of models.Data structs.
func MapTenantsResponseToStruct(resourcesResponse map[string]interface{}) ([]models.Data, error) {
	rawData, ok := resourcesResponse["data"].([]interface{})
	if !ok {
		logger.LogError("invalid data field in resourcesResponse")
		return nil, fmt.Errorf("missing or invalid data field")
	}

	var Tenants []models.Data
	for _, item := range rawData {
		TenantData, ok := item.(map[string]interface{})
		if !ok {
			logger.LogError("invalid Tenant data format in response")
			continue
		}
		Tenant, err := MapTenantData(TenantData)
		if err != nil {
			logger.LogError("error mapping Tenant data", "error", err)
			return nil, err
		}
		Tenants = append(Tenants, Tenant)
	}
	return Tenants, nil
}

// MapTenantResponseToStruct maps a resource response represented as a map to a slice of models.Data.
func MapTenantResponseToStruct(resourceResponse map[string]interface{}) ([]models.Data, error) {
	Tenant, err := MapTenantData(resourceResponse)
	if err != nil {
		logger.LogError("error mapping Tenant data in MapTenantResponseToStruct", "error", err)
		return nil, err
	}
	return []models.Data{Tenant}, nil
}

// mapTenantData maps a resource response to an Tenant model.
func MapTenantData(TenantData map[string]interface{}) (*models.Tenant, error) {
	logger.LogInfo("Mapping Tenant data to struct : ", "TenantData", TenantData)
	id, err := helpers.GetUUID(TenantData, "key")
	if err != nil {
		logger.LogError("failed to get UUID from Tenant data", "error", err)
		return nil, err
	}

	attributes, ok := TenantData["attributes"].(map[string]interface{})
	if !ok {
		logger.LogError("failed to get attributes map from Tenant data", "error", "missing or invalid map for key: attributes")
		return nil, fmt.Errorf("missing or invalid map for key: attributes")
	}

	if _, ok := attributes["contactInfo"].(map[string]interface{}); !ok {
		logger.LogError("failed to get contactInfo map from Tenant data", "error", "missing or invalid map for key 'contactInfo'")
		return nil, fmt.Errorf("missing or invalid map for key 'contactInfo'")
	}

	accountOwnerId, _ := helpers.GetUUID(attributes, "accountOwnerId")
	createdBy, _ := helpers.GetUUID(attributes, "createdBy")
	updatedBy, _ := helpers.GetUUID(attributes, "updatedBy")
	description := helpers.GetString(attributes, "description")
	status := helpers.GetString(attributes, "status")
	if status == "" {
		status = "ACTIVE" // Set default status to ACTIVE
	}
	contactInfo, _ := mapContactInfo(attributes)
	tags := tags.GetResourceTags(attributes, "tags")
	return &models.Tenant{
		ID:           id,
		Type:         config.Tenant,
		Name:         helpers.GetString(attributes, "name"),
		Status:       models.StatusTypeEnum(status),
		AccountOwner: &models.User{ID: accountOwnerId},
		Description:  &description,
		CreatedAt:    helpers.GetString(TenantData, "created_at"),
		UpdatedAt:    helpers.GetString(TenantData, "updated_at"),
		CreatedBy:    createdBy,
		UpdatedBy:    updatedBy,
		ContactInfo:  contactInfo,
		Tags:         tags,
	}, nil
}

// mapConcatInfo maps the billing information from the provided attributes map to a ContactInfo model.
func mapContactInfo(attributes map[string]interface{}) (*models.ContactInfo, error) {
	contactInfoData, err := helpers.GetMap(attributes, "contactInfo")
	if err != nil {
		return nil, err
	}

	contactAddress, err := mapContactAddress(contactInfoData)
	if err != nil {
		return nil, err
	}

	email := helpers.GetString(contactInfoData, "email")
	phoneNumber := helpers.GetString(contactInfoData, "phoneNumber")
	return &models.ContactInfo{
		Email:       &email,
		PhoneNumber: &phoneNumber,
		Address:     contactAddress,
	}, nil
}

// function mapContactAddress extracts address information from a given map and maps it to a address model.
func mapContactAddress(contactInfoData map[string]interface{}) (*models.Address, error) {
	contactAddressData, err := helpers.GetMap(contactInfoData, "address")
	if err != nil {
		logger.LogError("failed to get address map from contact info data", "error", err)
		return nil, err
	}

	street := helpers.GetString(contactAddressData, "street")
	city := helpers.GetString(contactAddressData, "city")
	state := helpers.GetString(contactAddressData, "state")
	zipcode := helpers.GetString(contactAddressData, "zipcode")
	country := helpers.GetString(contactAddressData, "country")

	return &models.Address{
		Street:  &street,
		City:    &city,
		State:   &state,
		Zipcode: &zipcode,
		Country: &country,
	}, nil
}
