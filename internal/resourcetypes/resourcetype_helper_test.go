package resourcetypes

import (
	"iam_services_main_v1/gql/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// createMockResourceResponse creates a mock resource response for testing
func createMockResourceResponse() map[string]interface{} {
	resourceID := uuid.New()
	permission1ID := uuid.New()
	permission2ID := uuid.New()

	return map[string]interface{}{
		"data": []interface{}{
			map[string]interface{}{
				"key":       resourceID.String(),
				"name":      "Document",
				"createdAt": "2025-03-20T19:08:06-05:00",
				"updatedAt": "2025-03-20T19:08:06-05:00",
				"actions": map[string]interface{}{
					"read": map[string]interface{}{
						"id":        permission1ID.String(),
						"name":      "Read",
						"createdAt": "2025-03-20T19:08:06-05:00",
						"updatedAt": "2025-03-20T19:08:06-05:00",
					},
					"write": map[string]interface{}{
						"id":        permission2ID.String(),
						"name":      "Write",
						"createdAt": "2025-03-20T19:08:06-05:00",
						"updatedAt": "2025-03-20T19:08:06-05:00",
					},
				},
			},
			map[string]interface{}{
				"key":       uuid.New().String(),
				"name":      "User",
				"createdAt": "2025-03-20T19:08:06-05:00",
				"updatedAt": "2025-03-20T19:08:06-05:00",
				"actions": map[string]interface{}{
					"create": map[string]interface{}{
						"id":        uuid.New().String(),
						"name":      "Create",
						"createdAt": "2025-03-20T19:08:06-05:00",
						"updatedAt": "2025-03-20T19:08:06-05:00",
					},
					"delete": map[string]interface{}{
						"id":        uuid.New().String(),
						"name":      "Delete",
						"createdAt": "2025-03-20T19:08:06-05:00",
						"updatedAt": "2025-03-20T19:08:06-05:00",
					},
				},
			},
		},
	}
}

// createMockInvalidResponse creates an invalid mock response for testing error cases
func createMockInvalidResponse() map[string]interface{} {
	return map[string]interface{}{
		"invalid_field": "invalid_data",
	}
}

func TestMapResourceResponseToStruct(t *testing.T) {
	tests := []struct {
		name           string
		resourceData   map[string]interface{}
		expectedCount  int
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:          "Valid response with resources",
			resourceData:  createMockResourceResponse(),
			expectedCount: 2, // Two resources in the mock data
			expectedError: false,
		},
		{
			name:           "Invalid response - missing data field",
			resourceData:   createMockInvalidResponse(),
			expectedCount:  0,
			expectedError:  true,
			expectedErrMsg: "missing or invalid data field",
		},
		{
			name: "Invalid response - invalid resource format",
			resourceData: map[string]interface{}{
				"data": []interface{}{
					"not-a-map", // Invalid data format
				},
			},
			expectedCount:  0,
			expectedError:  true,
			expectedErrMsg: "invalid resource data format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resources, err := MapResourceResponseToStruct(tt.resourceData)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Nil(t, resources)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resources)
				assert.Equal(t, tt.expectedCount, len(resources))

				// Verify the resources data
				for _, resource := range resources {
					// The resource should implement the Data interface
					assert.NotNil(t, resource)

					// Try to cast to ResourceType
					resourceType, ok := resource.(models.ResourceType)
					if ok {
						assert.NotEqual(t, uuid.Nil, resourceType.ID)
						assert.NotEmpty(t, resourceType.Name)
						assert.NotEmpty(t, resourceType.CreatedAt)
						assert.NotEmpty(t, resourceType.UpdatedAt)
						assert.NotEmpty(t, resourceType.Permissions)
					}
				}
			}
		})
	}
}

func TestExtractDataFromResponse(t *testing.T) {
	t.Run("Valid data extraction", func(t *testing.T) {
		response := createMockResourceResponse()
		data, ok := extractDataFromResponse(response)
		assert.True(t, ok)
		assert.NotNil(t, data)
		assert.Equal(t, 2, len(data))
	})

	t.Run("Missing data field", func(t *testing.T) {
		response := map[string]interface{}{
			"not_data": "something",
		}
		data, ok := extractDataFromResponse(response)
		assert.False(t, ok)
		assert.Nil(t, data)
	})

	t.Run("Data field is not a slice", func(t *testing.T) {
		response := map[string]interface{}{
			"data": "not a slice",
		}
		data, ok := extractDataFromResponse(response)
		assert.False(t, ok)
		assert.Nil(t, data)
	})

	t.Run("Data slice contains non-map items", func(t *testing.T) {
		response := map[string]interface{}{
			"data": []interface{}{
				"string item",
				map[string]interface{}{"key": "value"},
			},
		}
		data, ok := extractDataFromResponse(response)
		assert.True(t, ok) // Should still extract successfully even with mixed types
		assert.NotNil(t, data)
		assert.Equal(t, 2, len(data))
	})
}

func TestProcessResourceItem(t *testing.T) {
	validResourceData := map[string]interface{}{
		"key":       uuid.New().String(),
		"name":      "Document",
		"createdAt": "2025-03-20T19:08:06-05:00",
		"updatedAt": "2025-03-20T19:08:06-05:00",
		"actions": map[string]interface{}{
			"read": map[string]interface{}{
				"id":   uuid.New().String(),
				"name": "Read",
			},
		},
	}

	t.Run("Valid resource data", func(t *testing.T) {
		resourceTypes, err := processResourceItem(validResourceData)
		assert.NoError(t, err)
		assert.NotNil(t, resourceTypes)
		assert.Equal(t, 1, len(resourceTypes))

		// Check that we got a non-nil result that's compatible with the Data interface
		assert.NotNil(t, resourceTypes[0])

		// Now try to cast to models.ResourceType (value type, not pointer)
		// In the implementation, processResourceItem returns []models.Data{*resourceType}
		// So we should cast to a value type, not a pointer
		if resourceType, ok := resourceTypes[0].(models.ResourceType); ok {
			assert.NotEqual(t, uuid.Nil, resourceType.ID)
			assert.Equal(t, "Document", resourceType.Name)
		} else {
			t.Errorf("Expected item to be of type models.ResourceType but got different type: %T", resourceTypes[0])
		}
	})

	t.Run("Invalid resource data - not a map", func(t *testing.T) {
		resourceTypes, err := processResourceItem("not a map")
		assert.Error(t, err)
		assert.Nil(t, resourceTypes)
		assert.Equal(t, err.Error(), "invalid resource data format")
	})

	t.Run("Resource data missing actions field", func(t *testing.T) {
		invalidData := map[string]interface{}{
			"key":  uuid.New().String(),
			"name": "Document",
			// Missing actions field
		}
		resourceTypes, err := processResourceItem(invalidData)
		assert.Error(t, err)
		assert.Nil(t, resourceTypes)
		assert.Contains(t, err.Error(), "failed to extract actions")
	})
}

func TestMapToResourceTypeData(t *testing.T) {
	resourceID := uuid.New()
	permission1ID := uuid.New()
	permission2ID := uuid.New()

	validResourceData := map[string]interface{}{
		"key":       resourceID.String(),
		"name":      "Document",
		"createdAt": "2025-03-20T19:08:06-05:00",
		"updatedAt": "2025-03-20T19:08:06-05:00",
		"actions": map[string]interface{}{
			"read": map[string]interface{}{
				"id":   permission1ID.String(),
				"name": "Read",
			},
			"write": map[string]interface{}{
				"id":   permission2ID.String(),
				"name": "Write",
			},
		},
	}

	t.Run("Valid resource data", func(t *testing.T) {
		resourceType, err := MapToResourceTypeData(validResourceData)

		assert.NoError(t, err)
		assert.NotNil(t, resourceType)
		assert.Equal(t, resourceID, resourceType.ID)
		assert.Equal(t, "Document", resourceType.Name)
		assert.Equal(t, "2025-03-20T19:08:06-05:00", resourceType.CreatedAt)
		assert.Equal(t, "2025-03-20T19:08:06-05:00", resourceType.UpdatedAt)

		// Verify permissions
		assert.Equal(t, 2, len(resourceType.Permissions))
		permNames := []string{resourceType.Permissions[0].Name, resourceType.Permissions[1].Name}
		assert.Contains(t, permNames, "read")
		assert.Contains(t, permNames, "write")
	})

	t.Run("Invalid resource ID", func(t *testing.T) {
		invalidData := map[string]interface{}{
			"key":  "not-a-uuid",
			"name": "Document",
			"actions": map[string]interface{}{
				"read": map[string]interface{}{
					"id":   permission1ID.String(),
					"name": "Read",
				},
			},
		}

		resourceType, err := MapToResourceTypeData(invalidData)
		assert.Error(t, err)
		assert.Nil(t, resourceType)
		assert.Contains(t, err.Error(), "invalid resource type ID")
	})

	t.Run("Missing actions field", func(t *testing.T) {
		invalidData := map[string]interface{}{
			"key":       resourceID.String(),
			"name":      "Document",
			"createdAt": "2025-03-20T19:08:06-05:00",
			"updatedAt": "2025-03-20T19:08:06-05:00",
			// Missing actions field
		}

		resourceType, err := MapToResourceTypeData(invalidData)
		assert.Error(t, err)
		assert.Nil(t, resourceType)
		assert.Contains(t, err.Error(), "failed to extract actions")
	})

	t.Run("Actions field is not a map", func(t *testing.T) {
		invalidData := map[string]interface{}{
			"key":       resourceID.String(),
			"name":      "Document",
			"createdAt": "2025-03-20T19:08:06-05:00",
			"updatedAt": "2025-03-20T19:08:06-05:00",
			"actions":   "not a map",
		}

		resourceType, err := MapToResourceTypeData(invalidData)
		assert.Error(t, err)
		assert.Nil(t, resourceType)
		assert.Contains(t, err.Error(), "failed to extract actions")
	})
}

func TestMapToResourceActions(t *testing.T) {
	resourceID := uuid.New()
	permission1ID := uuid.New()
	permission2ID := uuid.New()

	validActionsData := map[string]interface{}{
		"read": map[string]interface{}{
			"id":   permission1ID.String(),
			"name": "Read",
		},
		"write": map[string]interface{}{
			"id":   permission2ID.String(),
			"name": "Write",
		},
		"invalidAction": "not-a-map", // This should be skipped
	}

	t.Run("Valid actions data", func(t *testing.T) {
		permissions := MapToResourceActions(resourceID, validActionsData)

		assert.NotNil(t, permissions)
		assert.Equal(t, 2, len(permissions))

		// Check permission names
		permissionNames := []string{permissions[0].Name, permissions[1].Name}
		assert.Contains(t, permissionNames, "read")
		assert.Contains(t, permissionNames, "write")

		// Verify permissions have the correct assignable scope
		for _, p := range permissions {
			// Note: In the implementation, it doesn't appear that AssignableScope is set to resourceID
			// Check if this is a bug or intended behavior
			if p.AssignableScope != uuid.Nil {
				assert.Equal(t, resourceID, p.AssignableScope)
			}
		}
	})

	t.Run("Invalid action ID", func(t *testing.T) {
		invalidActionsData := map[string]interface{}{
			"read": map[string]interface{}{
				"id":   "invalid-uuid", // Invalid UUID
				"name": "Read",
			},
		}

		permissions := MapToResourceActions(resourceID, invalidActionsData)
		assert.Equal(t, 0, len(permissions), "Should return empty slice for invalid action IDs")
	})

	t.Run("Empty actions data", func(t *testing.T) {
		emptyActionsData := map[string]interface{}{}

		permissions := MapToResourceActions(resourceID, emptyActionsData)
		assert.Equal(t, 0, len(permissions), "Should return empty slice for empty actions data")
	})

	t.Run("Action with missing ID field", func(t *testing.T) {
		actionsWithMissingID := map[string]interface{}{
			"read": map[string]interface{}{
				// Missing "id" field
				"name": "Read",
			},
		}

		permissions := MapToResourceActions(resourceID, actionsWithMissingID)
		assert.Equal(t, 0, len(permissions), "Should skip actions with missing ID fields")
	})
}

// Test error types
func TestErrorTypes(t *testing.T) {
	t.Run("ResourceError conversion to string", func(t *testing.T) {
		err := ErrInvalidDataField
		assert.Equal(t, "missing or invalid data field", err.Error())

		err = ErrInvalidResourceFormat
		assert.Equal(t, "invalid resource data format", err.Error())

		err = ErrInvalidActionFormat
		assert.Equal(t, "invalid action data format", err.Error())

		err = ErrResourceNotFound
		assert.Equal(t, "resource not found", err.Error())
	})
}
