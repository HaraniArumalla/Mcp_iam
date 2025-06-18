package accounts

import (
	"fmt"
	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/internal/tags"
	"iam_services_main_v1/pkg/logger"
)

// MapAccountsResponseToStruct maps a response containing resource data to a slice of models.Data structs.
func MapAccountsResponseToStruct(resourcesResponse map[string]interface{}) ([]models.Data, error) {
	rawData, ok := resourcesResponse["data"].([]interface{})
	if !ok {
		logger.LogError("invalid data field in resourcesResponse")
		return nil, fmt.Errorf("missing or invalid data field")
	}

	var accounts []models.Data
	for _, item := range rawData {
		accountData, ok := item.(map[string]interface{})
		if !ok {
			logger.LogError("invalid account data format in response")
			continue
		}
		account, err := mapAccountData(accountData)
		if err != nil {
			logger.LogError("error mapping account data", "error", err)
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

// MapAccountResponseToStruct maps a resource response represented as a map to a slice of models.Data.
func MapAccountResponseToStruct(resourceResponse map[string]interface{}) ([]models.Data, error) {
	account, err := mapAccountData(resourceResponse)
	if err != nil {
		logger.LogError("error mapping account data in MapAccountResponseToStruct", "error", err)
		return nil, err
	}
	return []models.Data{account}, nil
}

// mapAccountData maps a resource response to an Account model.
func mapAccountData(accountData map[string]interface{}) (*models.Account, error) {
	id, err := helpers.GetUUID(accountData, "key")
	if err != nil {
		logger.LogError("failed to get UUID from account data", "error", err)
		return nil, err
	}

	attributes, err := helpers.GetMap(accountData, "attributes")
	if err != nil {
		logger.LogError("failed to get attributes map from account data", "error", err)
		return nil, err
	}
	accountOwnerId, _ := helpers.GetUUID(attributes, "accountOwnerId")
	parentId, _ := helpers.GetUUID(attributes, "parentId")
	tenantId, _ := helpers.GetUUID(attributes, "tenantId")
	createdBy, _ := helpers.GetUUID(attributes, "createdBy")
	updatedBy, _ := helpers.GetUUID(attributes, "updatedBy")
	description := helpers.GetString(attributes, "description")
	status := helpers.GetString(attributes, "status")
	relationType := helpers.GetString(attributes, "relationType")
	billingInfo, _ := mapBillingInfo(attributes)
	tags := tags.GetResourceTags(attributes, "tags")

	var parentOrg models.Organization
	if relationType == "CHILD" {
		parentOrg = &models.Tenant{ID: parentId}
	} else {
		parentOrg = &models.ClientOrganizationUnit{ID: parentId}
	}

	return &models.Account{
		ID:           id,
		Type:         config.Account,
		ParentOrg:    parentOrg,
		Tenant:       &models.Tenant{ID: tenantId},
		AccountOwner: &models.User{ID: accountOwnerId},
		Status:       models.StatusTypeEnum(status),
		RelationType: models.RelationTypeEnum(relationType),
		Name:         helpers.GetString(attributes, "name"),
		Description:  &description,
		CreatedBy:    createdBy,
		UpdatedBy:    updatedBy,
		CreatedAt:    helpers.GetString(accountData, "created_at"),
		UpdatedAt:    helpers.GetString(accountData, "updated_at"),
		BillingInfo:  billingInfo,
		Tags:         tags,
	}, nil
}

// mapBillingInfo maps the billing information from the provided attributes map to a BillingInfo model.
func mapBillingInfo(attributes map[string]interface{}) (*models.BillingInfo, error) {
	billingInfoData, err := helpers.GetMap(attributes, "billingInfo")
	if err != nil {
		return nil, err
	}

	billingAddress, err := mapBillingAddress(billingInfoData)
	if err != nil {
		return nil, err
	}

	return &models.BillingInfo{
		CreditCardNumber: helpers.GetString(billingInfoData, "creditCardNumber"),
		CreditCardType:   helpers.GetString(billingInfoData, "creditCardType"),
		ExpirationDate:   helpers.GetString(billingInfoData, "expirationDate"),
		Cvv:              helpers.GetString(billingInfoData, "cvv"),
		BillingAddress:   billingAddress,
	}, nil
}

// mapBillingAddress extracts billing address information from a given map and maps it to a BillingAddress model.
func mapBillingAddress(billingInfoData map[string]interface{}) (*models.BillingAddress, error) {
	billingAddressData, err := helpers.GetMap(billingInfoData, "billingAddress")
	if err != nil {
		logger.LogError("failed to get billing address map from billing info data", "error", err)
		return nil, err
	}

	return &models.BillingAddress{
		Street:  helpers.GetString(billingAddressData, "street"),
		City:    helpers.GetString(billingAddressData, "city"),
		State:   helpers.GetString(billingAddressData, "state"),
		Zipcode: helpers.GetString(billingAddressData, "zipcode"),
		Country: helpers.GetString(billingAddressData, "country"),
	}, nil
}
