package tenants

import (
	"iam_services_main_v1/gql/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func createValidTenantResponse() map[string]interface{} {
	tenantID := uuid.New()
	userID := uuid.New()

	return map[string]interface{}{
		"key":        tenantID.String(),
		"created_at": "2023-01-01T00:00:00Z",
		"updated_at": "2023-01-02T00:00:00Z",
		"attributes": map[string]interface{}{
			"id":          tenantID.String(),
			"name":        "Test Tenant",
			"description": "A test tenant",
			"createdBy":   userID.String(),
			"updatedBy":   userID.String(),
			"contactInfo": map[string]interface{}{
				"email":       "test@example.com",
				"phoneNumber": "1234567890",
				"address": map[string]interface{}{
					"street":  "123 Test St",
					"city":    "Test City",
					"state":   "TS",
					"zipcode": "12345",
					"country": "US",
				},
			},
			"tags": []interface{}{
				map[string]interface{}{
					"key":   "environment",
					"value": "test",
				},
				map[string]interface{}{
					"key":   "department",
					"value": "engineering",
				},
			},
		},
	}
}

func createTenantsListResponse() map[string]interface{} {
	tenant1 := createValidTenantResponse()
	tenant2 := createValidTenantResponse()

	// Modify tenant2 to have different values
	tenant2Attrs := tenant2["attributes"].(map[string]interface{})
	tenant2Attrs["name"] = "Second Tenant"

	return map[string]interface{}{
		"data": []interface{}{
			tenant1,
			tenant2,
		},
	}
}

func TestMapTenantsResponseToStruct(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		response := createTenantsListResponse()

		// Execute
		tenants, err := MapTenantsResponseToStruct(response)

		// Verify
		assert.NoError(t, err)
		assert.Len(t, tenants, 2)
		assert.Equal(t, "Test Tenant", tenants[0].(*models.Tenant).Name)
		assert.Equal(t, "Second Tenant", tenants[1].(*models.Tenant).Name)
		assert.Len(t, tenants[0].(*models.Tenant).Tags, 2)
	})

	t.Run("Invalid data field", func(t *testing.T) {
		// Setup - missing data field
		response := map[string]interface{}{
			"not_data": []interface{}{},
		}

		// Execute
		tenants, err := MapTenantsResponseToStruct(response)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, tenants)
		assert.Contains(t, err.Error(), "missing or invalid data field")
	})

	t.Run("Invalid tenant data format", func(t *testing.T) {
		// Setup - data contains non-map item
		response := map[string]interface{}{
			"data": []interface{}{
				"not a map",
			},
		}

		// Execute
		tenants, err := MapTenantsResponseToStruct(response)

		// Verify
		assert.NoError(t, err) // This won't error out, it will just log and continue
		assert.Empty(t, tenants)
	})

	t.Run("Error mapping tenant data", func(t *testing.T) {
		// Setup - tenant missing key field
		response := map[string]interface{}{
			"data": []interface{}{
				map[string]interface{}{
					"attributes": map[string]interface{}{
						"name": "Invalid Tenant",
					},
				},
			},
		}

		// Execute
		tenants, err := MapTenantsResponseToStruct(response)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, tenants)
	})
}

func TestMapTenantResponseToStruct(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		response := createValidTenantResponse()

		// Execute
		tenants, err := MapTenantResponseToStruct(response)

		// Verify
		assert.NoError(t, err)
		assert.Len(t, tenants, 1)
		tenant := tenants[0].(*models.Tenant)
		assert.Equal(t, "Test Tenant", tenant.Name)
		assert.Equal(t, "A test tenant", *tenant.Description)
		assert.Equal(t, "test@example.com", *tenant.ContactInfo.Email)
		assert.Equal(t, "123 Test St", *tenant.ContactInfo.Address.Street)
		assert.Len(t, tenant.Tags, 2)
	})

	t.Run("Error mapping tenant data", func(t *testing.T) {
		// Setup - tenant missing key field
		response := map[string]interface{}{
			"attributes": map[string]interface{}{
				"name": "Invalid Tenant",
			},
		}

		// Execute
		tenants, err := MapTenantResponseToStruct(response)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, tenants)
		assert.Contains(t, err.Error(), "invalid or missing UUID")
	})
}

func TestMapTenantData(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		tenantData := createValidTenantResponse()

		// Execute
		tenant, err := MapTenantData(tenantData)

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, tenant)

		// Verify all fields are correctly mapped
		assert.Equal(t, tenantData["key"], tenant.ID.String())

		attributes := tenantData["attributes"].(map[string]interface{})
		assert.Equal(t, attributes["name"], tenant.Name)
		assert.Equal(t, attributes["description"], *tenant.Description)
		assert.Equal(t, tenantData["created_at"], tenant.CreatedAt)
		assert.Equal(t, tenantData["updated_at"], tenant.UpdatedAt)
		assert.Equal(t, attributes["createdBy"], tenant.CreatedBy.String())
		assert.Equal(t, attributes["updatedBy"], tenant.UpdatedBy.String())

		// Verify tags
		assert.Len(t, tenant.Tags, 2)
		assert.Equal(t, "environment", tenant.Tags[0].Key)
		assert.Equal(t, "test", tenant.Tags[0].Value)
		assert.Equal(t, "department", tenant.Tags[1].Key)
		assert.Equal(t, "engineering", tenant.Tags[1].Value)
	})

	t.Run("Missing key field", func(t *testing.T) {
		// Setup
		tenantData := map[string]interface{}{
			"not_key":    uuid.New().String(),
			"attributes": map[string]interface{}{},
		}

		// Execute
		tenant, err := MapTenantData(tenantData)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, tenant)
		assert.Contains(t, err.Error(), "invalid or missing UUID")
	})

	t.Run("Missing attributes", func(t *testing.T) {
		// Setup
		tenantData := map[string]interface{}{
			"key":            uuid.New().String(),
			"not_attributes": map[string]interface{}{},
		}

		// Execute
		tenant, err := MapTenantData(tenantData)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, tenant)
		assert.Contains(t, err.Error(), "missing or invalid map for key: attributes")
	})

	t.Run("Missing contact info", func(t *testing.T) {
		// Setup
		tenantData := map[string]interface{}{
			"key": uuid.New().String(),
			"attributes": map[string]interface{}{
				"name": "Test Tenant",
				// Missing contactInfo
			},
		}

		// Execute
		tenant, err := MapTenantData(tenantData)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, tenant)
		assert.Contains(t, err.Error(), "missing or invalid map for key 'contactInfo'")
	})
}

func TestMapConcatInfo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		attributes := map[string]interface{}{
			"contactInfo": map[string]interface{}{
				"email":       "test@example.com",
				"phoneNumber": "1234567890",
				"address": map[string]interface{}{
					"street":  "123 Test St",
					"city":    "Test City",
					"state":   "TS",
					"zipcode": "12345",
					"country": "US",
				},
			},
		}

		// Execute
		contactInfo, err := mapContactInfo(attributes)

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, contactInfo)
		assert.Equal(t, "test@example.com", *contactInfo.Email)
		assert.Equal(t, "1234567890", *contactInfo.PhoneNumber)
		assert.NotNil(t, contactInfo.Address)
		assert.Equal(t, "123 Test St", *contactInfo.Address.Street)
	})

	t.Run("Missing contact info", func(t *testing.T) {
		// Setup
		attributes := map[string]interface{}{
			"not_contactInfo": map[string]interface{}{},
		}

		// Execute
		contactInfo, err := mapContactInfo(attributes)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, contactInfo)
	})

	t.Run("Missing address", func(t *testing.T) {
		// Setup
		attributes := map[string]interface{}{
			"contactInfo": map[string]interface{}{
				"email":       "test@example.com",
				"phoneNumber": "1234567890",
				// Missing address
			},
		}

		// Execute
		contactInfo, err := mapContactInfo(attributes)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, contactInfo)
	})
}

func TestMapContactAddress(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		contactInfoData := map[string]interface{}{
			"address": map[string]interface{}{
				"street":  "123 Test St",
				"city":    "Test City",
				"state":   "TS",
				"zipcode": "12345",
				"country": "US",
			},
		}

		// Execute
		address, err := mapContactAddress(contactInfoData)

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, address)
		assert.Equal(t, "123 Test St", *address.Street)
		assert.Equal(t, "Test City", *address.City)
		assert.Equal(t, "TS", *address.State)
		assert.Equal(t, "12345", *address.Zipcode)
		assert.Equal(t, "US", *address.Country)
	})

	t.Run("Missing address", func(t *testing.T) {
		// Setup
		contactInfoData := map[string]interface{}{
			"not_address": map[string]interface{}{},
		}

		// Execute
		address, err := mapContactAddress(contactInfoData)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, address)
		assert.Contains(t, err.Error(), "missing or invalid map for key")
	})
}
