package accounts

import (
	"context"
	"errors"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAccountFieldResolver_ParentOrg(t *testing.T) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock PermitService
	mockPermitService := mocks.NewMockPermitService(ctrl)

	// Create AccountFieldResolver with mock
	resolver := &AccountFieldResolver{
		PC: mockPermitService,
	}

	// Test context
	ctx := context.Background()

	// Test data
	account := &models.Account{
		ID:        uuid.New(),
		ParentOrg: &models.ClientOrganizationUnit{ID: uuid.New()},
	}

	t.Run("Success_With_ClientOrg", func(t *testing.T) {
		// Mock response data for ClientOrganizationUnit
		mockResponse := buildTestClientOrganizationData(account.ParentOrg.GetID())

		// Set up expected behavior
		mockPermitService.EXPECT().
			GetSingleResource(ctx, "GET", gomock.Any()).
			Return(mockResponse, nil).MaxTimes(1)

		// Call method being tested
		result, err := resolver.ParentOrg(ctx, account)

		// Assert results
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Test Org", result.GetName())
	})

	t.Run("Success_With_Tenant", func(t *testing.T) {
		// Mock response data with tenant type
		mockResponse := buildTestTenantData(account.ParentOrg.GetID())
		mockResponse["attributes"].(map[string]interface{})["type"] = "Tenant"

		// Set up expected behavior
		mockPermitService.EXPECT().
			GetSingleResource(ctx, "GET", gomock.Any()).
			Return(mockResponse, nil).MaxTimes(1)

		// Call method being tested
		result, err := resolver.ParentOrg(ctx, account)

		// Assert results
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Test Tenant", result.GetName())
	})

	t.Run("Error_Fetching_ParentOrg", func(t *testing.T) {
		// Set up expected behavior for error case
		expectedError := errors.New("failed to fetch resource")
		mockPermitService.EXPECT().
			GetSingleResource(ctx, "GET", gomock.Any()).
			Return(nil, expectedError).MaxTimes(1)

		// Call method being tested
		result, err := resolver.ParentOrg(ctx, account)

		// Assert error
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
	})

	t.Run("Error_Invalid_Attributes", func(t *testing.T) {
		// Mock response data without attributes
		mockResponse := map[string]interface{}{
			"key":  account.ParentOrg.GetID().String(),
			"name": "Test Org",
		}

		// Set up expected behavior
		mockPermitService.EXPECT().
			GetSingleResource(ctx, "GET", gomock.Any()).
			Return(mockResponse, nil).MaxTimes(1)

		// Call method being tested
		result, err := resolver.ParentOrg(ctx, account)

		// Assert error
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "missing or invalid map for key: attributes")
	})
}

func TestAccountFieldResolver_Tenant(t *testing.T) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock PermitService
	mockPermitService := mocks.NewMockPermitService(ctrl)

	// Create AccountFieldResolver with mock
	resolver := &AccountFieldResolver{
		PC: mockPermitService,
	}

	// Test context
	ctx := context.Background()

	// Test data
	tenantID := uuid.New()
	account := &models.Account{
		ID:     uuid.New(),
		Tenant: &models.Tenant{ID: tenantID},
	}

	t.Run("Success", func(t *testing.T) {
		// Mock response data for tenant
		mockTenantResponse := buildTestTenantData(tenantID)

		// Set up expected behavior
		mockPermitService.EXPECT().
			SendRequest(gomock.Any(), "GET", "tenants/"+tenantID.String(), nil).
			Return(mockTenantResponse, nil).MaxTimes(1)

		// Call method being tested
		result, err := resolver.Tenant(ctx, account)

		// Assert results
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Test Tenant", result.Name)
		assert.Equal(t, tenantID, result.ID)
		assert.Equal(t, "Test Tenant Description", *result.Description)
	})

	t.Run("Error_Fetching_Tenant", func(t *testing.T) {
		// Set up expected behavior for error case
		expectedError := errors.New("failed to fetch tenant")
		mockPermitService.EXPECT().
			SendRequest(gomock.Any(), "GET", "tenants/"+tenantID.String(), nil).
			Return(nil, expectedError).MaxTimes(1)

		// Call method being tested
		result, err := resolver.Tenant(ctx, account)

		// Assert error
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
	})

	t.Run("Error_Mapping_Tenant", func(t *testing.T) {
		// Mock response data with invalid format to cause mapping error
		mockInvalidResponse := map[string]interface{}{
			"id":   "not-a-valid-uuid", // This will cause mapping error
			"name": "Invalid Tenant",
		}

		// Set up expected behavior
		mockPermitService.EXPECT().
			SendRequest(gomock.Any(), "GET", "tenants/"+tenantID.String(), nil).
			Return(mockInvalidResponse, nil).MaxTimes(1)

		// Call method being tested
		result, err := resolver.Tenant(ctx, account)

		// Assert error
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid or missing UUID for key") // This would be part of the error message
	})
}

func TestAccountFieldResolver_AccountOwner(t *testing.T) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock PermitService
	mockPermitService := mocks.NewMockPermitService(ctrl)

	// Create AccountFieldResolver with mock
	resolver := &AccountFieldResolver{
		PC: mockPermitService,
	}

	// Test context
	ctx := context.Background()

	// Test data
	ownerID := uuid.New()
	account := &models.Account{
		ID:           uuid.New(),
		AccountOwner: &models.User{ID: ownerID},
	}

	t.Run("Success", func(t *testing.T) {
		// Mock response data for user
		mockUserResponse := buildTestUserData(ownerID)

		// Set up expected behavior - use GetSingleResource instead of SendRequest
		mockPermitService.EXPECT().
			GetSingleResource(ctx, "GET", "users/"+ownerID.String()).
			Return(mockUserResponse, nil).MaxTimes(1)

		// Call method being tested
		result, err := resolver.AccountOwner(ctx, account)

		// Assert results
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ownerID, result.ID)
		assert.Equal(t, "Test User", result.Name)
		assert.Equal(t, "test@example.com", result.Email)
	})

	t.Run("Error_Fetching_User", func(t *testing.T) {
		// Set up expected behavior for error case - use GetSingleResource instead of SendRequest
		expectedError := errors.New("failed to fetch user")
		mockPermitService.EXPECT().
			GetSingleResource(ctx, "GET", "users/"+ownerID.String()).
			Return(nil, expectedError).MaxTimes(1)

		// Call method being tested
		result, err := resolver.AccountOwner(ctx, account)

		// Assert error
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
	})
}

// Helper function to build test tenant data
func buildTestTenantData(id uuid.UUID) map[string]interface{} {
	return map[string]interface{}{
		"key":        id.String(),
		"name":       "Test Tenant",
		"created_at": time.Now().String(),
		"updated_at": time.Now().String(),
		"attributes": map[string]interface{}{
			"id":          id.String(),
			"name":        "Test Tenant",
			"description": "Test Tenant Description",
			"tenantId":    id.String(),
			"parentId":    uuid.New().String(),
			"createdBy":   uuid.New().String(),
			"updatedBy":   uuid.New().String(),
			"contactInfo": map[string]interface{}{
				"email":       "test@example.com",
				"phoneNumber": "123-456-7890",
				"address": map[string]interface{}{
					"street":  "123 Test St",
					"city":    "Test City",
					"state":   "TS",
					"zipcode": "12345",
					"country": "US",
				},
			},
		},
	}
}

func buildTestClientOrganizationData(id uuid.UUID) map[string]interface{} {
	result := make(map[string]interface{}, 0)
	attributes := make(map[string]interface{}, 0)
	attributes["tenantId"] = id.String()
	attributes["parentOrgId"] = "parentId"
	attributes["created_at"] = time.Now().String()
	attributes["updated_at"] = time.Now().String()
	attributes["created_by"] = uuid.New().String()
	attributes["updated_by"] = uuid.New().String()
	attributes["description"] = "description"
	attributes["account_owner_id"] = uuid.New().String()
	attributes["status"] = "ACTIVE"
	attributes["relation_type"] = "SELF"
	attributes["name"] = "Test Org"
	attributes["type"] = "ClientOrganizationUnit"
	attributes["key"] = id.String()
	result["key"] = id.String()
	result["name"] = "Test Org"
	result["attributes"] = attributes

	return result
}

// Helper function to build test user data
func buildTestUserData(id uuid.UUID) map[string]interface{} {
	return map[string]interface{}{
		"key":        id.String(),
		"first_name": "Test",
		"last_name":  "User",
		"email":      "test@example.com",
		"created_at": time.Now().String(),
		"updated_at": time.Now().String(),
		"associated_tenants": []interface{}{
			map[string]interface{}{
				"tenant": uuid.New().String(),
			},
		},
	}
}
