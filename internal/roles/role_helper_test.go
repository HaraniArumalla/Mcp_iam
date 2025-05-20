package roles

import (
	"iam_services_main_v1/gql/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// createMockResourceResponse creates a mock resource response for testing
func createMockResourceResponse() map[string]interface{} {
	role1ID := uuid.New()
	role2ID := uuid.New()
	permission1ID := uuid.New()
	permission2ID := uuid.New()
	userID := uuid.New().String()
	resourceID := uuid.New().String()

	return map[string]interface{}{
		"data": []interface{}{
			map[string]interface{}{
				"key":       resourceID,
				"name":      "Resource 1",
				"createdAt": "2025-03-20T19:08:06-05:00",
				"updatedAt": "2025-03-20T19:08:06-05:00",
				"roles": map[string]interface{}{
					"admin": map[string]interface{}{
						"name":        "Admin",
						"description": "Administrator role",
						"key":         role1ID.String(),
						"permissions": []interface{}{"read", "write"},
						"attributes": map[string]interface{}{
							"createdBy":          userID,
							"updatedBy":          userID,
							"createdAt":          "2025-03-20T19:08:06-05:00",
							"updatedAt":          "2025-03-20T19:08:06-05:00",
							"assignableScopeRef": uuid.New().String(),
						},
					},
					"user": map[string]interface{}{
						"name":        "User",
						"description": "Regular user role",
						"key":         role2ID.String(),
						"permissions": []interface{}{"read"},
						"attributes": map[string]interface{}{
							"createdBy":          userID,
							"updatedBy":          userID,
							"createdAt":          "2025-03-20T19:08:06-05:00",
							"updatedAt":          "2025-03-20T19:08:06-05:00",
							"assignableScopeRef": uuid.New().String(),
						},
					},
				},
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

func TestMapRolesResponseToStruct(t *testing.T) {
	tests := []struct {
		name           string
		resourceData   map[string]interface{}
		expectedCount  int
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:          "Valid response with roles",
			resourceData:  createMockResourceResponse(),
			expectedCount: 2,
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
			name: "Invalid response - invalid roles format",
			resourceData: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"key":     "resource-1",
						"name":    "Resource 1",
						"roles":   "invalid-roles", // Invalid roles format (not a map)
						"actions": map[string]interface{}{},
					},
				},
			},
			expectedCount:  0,
			expectedError:  true,
			expectedErrMsg: "failed to extract roles",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roles, err := MapRolesResponseToStruct(tt.resourceData)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Nil(t, roles)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, roles)
				assert.Equal(t, tt.expectedCount, len(roles))

				// Verify role fields
				for _, role := range roles {
					typedRole, ok := role.(*models.Role)
					assert.True(t, ok, "Expected *models.Role type")

					assert.NotEqual(t, uuid.Nil, typedRole.ID)
					assert.NotEmpty(t, typedRole.Name)
					assert.NotNil(t, typedRole.Description)
				}
			}
		})
	}
}

func TestMapRoleResponseToStruct(t *testing.T) {
	mockResponse := createMockResourceResponse()

	// Extract a role ID for testing
	var roleID uuid.UUID
	data := mockResponse["data"].([]interface{})[0].(map[string]interface{})
	roles := data["roles"].(map[string]interface{})
	adminRole := roles["admin"].(map[string]interface{})
	id, _ := uuid.Parse(adminRole["key"].(string))
	roleID = id

	nonExistentID := uuid.New()

	tests := []struct {
		name           string
		resourceData   map[string]interface{}
		roleID         uuid.UUID
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:          "Valid role ID",
			resourceData:  mockResponse,
			roleID:        roleID,
			expectedError: false,
		},
		{
			name:           "Non-existent role ID",
			resourceData:   mockResponse,
			roleID:         nonExistentID,
			expectedError:  true,
			expectedErrMsg: "role not found",
		},
		{
			name:           "Invalid response - missing data field",
			resourceData:   createMockInvalidResponse(),
			roleID:         roleID,
			expectedError:  true,
			expectedErrMsg: "missing or invalid data field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roles, err := MapRoleResponseToStruct(tt.resourceData, tt.roleID)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Nil(t, roles)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, roles)
				assert.Equal(t, 1, len(roles))

				// Verify the role has the correct ID
				typedRole, ok := roles[0].(*models.Role)
				assert.True(t, ok, "Expected *models.Role type")
				assert.Equal(t, tt.roleID, typedRole.ID)
			}
		})
	}
}

func TestMapToRoleData(t *testing.T) {
	userID := uuid.New().String()
	roleID := uuid.New().String()
	permission1ID := uuid.New()
	permission2ID := uuid.New()
	resourceID := uuid.New().String()

	validRoleData := map[string]interface{}{
		"name":        "Admin",
		"description": "Administrator role",
		"key":         roleID,
		"permissions": []interface{}{"read", "write"},
		"attributes": map[string]interface{}{
			"createdBy":          userID,
			"updatedBy":          userID,
			"createdAt":          "2025-03-20T19:08:06-05:00",
			"updatedAt":          "2025-03-20T19:08:06-05:00",
			"assignableScopeRef": uuid.New().String(),
		},
	}

	actionsData := map[string]interface{}{
		"read": map[string]interface{}{
			"id":   permission1ID.String(),
			"name": "Read",
		},
		"write": map[string]interface{}{
			"id":   permission2ID.String(),
			"name": "Write",
		},
	}

	resourceData := map[string]interface{}{
		"key":       resourceID,
		"name":      "Resource 1",
		"createdAt": "2025-03-20T19:08:06-05:00",
		"updatedAt": "2025-03-20T19:08:06-05:00",
		"actions": map[string]interface{}{
			"read": map[string]interface{}{
				"id":   permission1ID.String(),
				"name": "Read",
			},
		},
	}

	t.Run("Valid role data", func(t *testing.T) {
		role, err := MapToRoleData(validRoleData, actionsData, resourceData)

		assert.NoError(t, err)
		assert.NotNil(t, role)
		assert.Equal(t, "Admin", role.Name)
		assert.Equal(t, "Administrator role", *role.Description)
		assert.NotNil(t, role.Permissions)
		assert.Equal(t, 2, len(role.Permissions))
		assert.NotNil(t, role.AssignableScope)
		assert.Equal(t, resourceID, role.AssignableScope.ID.String())
	})

	t.Run("Invalid role ID", func(t *testing.T) {
		invalidRoleData := map[string]interface{}{
			"name":        "Admin",
			"description": "Administrator role",
			"key":         "invalid-uuid", // Invalid UUID
			"permissions": []interface{}{"read", "write"},
			"attributes": map[string]interface{}{
				"createdBy": userID,
				"updatedBy": userID,
				"createdAt": "2025-03-20T19:08:06-05:00",
				"updatedAt": "2025-03-20T19:08:06-05:00",
			},
		}

		role, err := MapToRoleData(invalidRoleData, actionsData, resourceData)

		assert.Error(t, err)
		assert.Nil(t, role)
		assert.Contains(t, err.Error(), "invalid role ID")
	})

	t.Run("Missing attributes", func(t *testing.T) {
		invalidRoleData := map[string]interface{}{
			"name":        "Admin",
			"description": "Administrator role",
			"key":         roleID,
			"permissions": []interface{}{"read", "write"},
			// Missing attributes
		}

		role, err := MapToRoleData(invalidRoleData, actionsData, resourceData)

		assert.Error(t, err)
		assert.Nil(t, role)
		assert.Contains(t, err.Error(), "invalid attributes")
	})

	t.Run("Invalid resource type", func(t *testing.T) {
		invalidResourceData := map[string]interface{}{
			"key": "invalid-uuid", // Invalid resource ID
		}

		role, err := MapToRoleData(validRoleData, actionsData, invalidResourceData)

		assert.Error(t, err)
		assert.Nil(t, role)
		assert.Contains(t, err.Error(), "invalid resource type")
	})
}

func TestMapToResourceTypeData(t *testing.T) {
	permission1ID := uuid.New()
	resourceID := uuid.New()

	validResourceData := map[string]interface{}{
		"key":       resourceID.String(),
		"name":      "Resource 1",
		"createdAt": "2025-03-20T19:08:06-05:00",
		"updatedAt": "2025-03-20T19:08:06-05:00",
		"actions": map[string]interface{}{
			"read": map[string]interface{}{
				"id":   permission1ID.String(),
				"name": "Read",
			},
		},
	}

	t.Run("Valid resource data", func(t *testing.T) {
		resourceType, err := MapToResourceTypeData(validResourceData)

		assert.NoError(t, err)
		assert.NotNil(t, resourceType)
		assert.Equal(t, resourceID, resourceType.ID)
		assert.Equal(t, "Resource 1", resourceType.Name)
		assert.Equal(t, "2025-03-20T19:08:06-05:00", resourceType.CreatedAt)
		assert.Equal(t, "2025-03-20T19:08:06-05:00", resourceType.UpdatedAt)
		assert.NotEmpty(t, resourceType.Permissions)
		assert.Equal(t, 1, len(resourceType.Permissions))
	})

	t.Run("Invalid resource ID", func(t *testing.T) {
		invalidResourceData := map[string]interface{}{
			"key":     "invalid-uuid",
			"name":    "Invalid Resource",
			"actions": map[string]interface{}{},
		}

		resourceType, err := MapToResourceTypeData(invalidResourceData)

		assert.Error(t, err)
		assert.Nil(t, resourceType)
		assert.Contains(t, err.Error(), "invalid resource type ID")
	})

	t.Run("Missing actions", func(t *testing.T) {
		invalidResourceData := map[string]interface{}{
			"key":  resourceID.String(),
			"name": "Resource Without Actions",
			// Missing actions field
		}

		resourceType, err := MapToResourceTypeData(invalidResourceData)

		assert.Error(t, err)
		assert.Nil(t, resourceType)
		assert.Contains(t, err.Error(), "failed to extract actions")
	})
}

func TestMapToResourceActions(t *testing.T) {
	resourceID := uuid.New()
	permission1ID := uuid.New()
	permission2ID := uuid.New()

	actionsData := map[string]interface{}{
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
		permissions := MapToResourceActions(resourceID, actionsData)

		assert.NotNil(t, permissions)
		assert.Equal(t, 2, len(permissions))

		// Check for both permissions
		permissionNames := []string{permissions[0].Name, permissions[1].Name}
		assert.Contains(t, permissionNames, "read")
		assert.Contains(t, permissionNames, "write")

		// Check that resourceID is properly set as assignable scope
		for _, p := range permissions {
			assert.Equal(t, resourceID, p.AssignableScope)
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

		// Should return empty slice - no valid permissions
		assert.Equal(t, 0, len(permissions))
	})

	t.Run("Empty actions data", func(t *testing.T) {
		emptyActionsData := map[string]interface{}{}

		permissions := MapToResourceActions(resourceID, emptyActionsData)
		assert.Equal(t, 0, len(permissions))
	})
}

func TestMapToPermissionsData(t *testing.T) {
	permission1ID := uuid.New()
	permission2ID := uuid.New()

	validRoleData := map[string]interface{}{
		"permissions": []interface{}{"read", "write"},
	}

	validActionsData := map[string]interface{}{
		"read": map[string]interface{}{
			"id":   permission1ID.String(),
			"name": "Read",
		},
		"write": map[string]interface{}{
			"id":   permission2ID.String(),
			"name": "Write",
		},
	}

	t.Run("Valid permissions data", func(t *testing.T) {
		permissions, err := MapToPermissionsData(validRoleData, validActionsData)

		assert.NoError(t, err)
		assert.NotNil(t, permissions)
		assert.Equal(t, 2, len(permissions))

		// Check the permissions data
		permissionNames := []string{permissions[0].Name, permissions[1].Name}
		assert.Contains(t, permissionNames, "Read")
		assert.Contains(t, permissionNames, "Write")
	})

	t.Run("Missing permissions in role data", func(t *testing.T) {
		invalidRoleData := map[string]interface{}{
			// Missing permissions field
		}

		permissions, err := MapToPermissionsData(invalidRoleData, validActionsData)

		assert.Error(t, err)
		assert.Nil(t, permissions)
	})

	t.Run("Invalid permission format in actions data", func(t *testing.T) {
		invalidActionsData := map[string]interface{}{
			"read": "invalid-action-data", // Not a map
		}

		roleData := map[string]interface{}{
			"permissions": []interface{}{"read"},
		}

		permissions, err := MapToPermissionsData(roleData, invalidActionsData)

		// This should not error but return nil permission for the invalid action
		assert.NoError(t, err)
		assert.Empty(t, permissions)
	})

	t.Run("Non-string permission key", func(t *testing.T) {
		roleData := map[string]interface{}{
			"permissions": []interface{}{123}, // Non-string permission
		}

		permissions, err := MapToPermissionsData(roleData, validActionsData)

		assert.NoError(t, err)
		assert.Empty(t, permissions) // Should filter out invalid permission
	})

	t.Run("Unknown permission key", func(t *testing.T) {
		roleData := map[string]interface{}{
			"permissions": []interface{}{"unknown"}, // Not in actionsData
		}

		permissions, err := MapToPermissionsData(roleData, validActionsData)

		assert.NoError(t, err)
		assert.Empty(t, permissions) // Should filter out unknown permission
	})
}

func TestExtractRolesAndActions(t *testing.T) {
	validResourceData := map[string]interface{}{
		"roles": map[string]interface{}{
			"admin": map[string]interface{}{
				"name": "Admin",
			},
		},
		"actions": map[string]interface{}{
			"read": map[string]interface{}{
				"name": "Read",
			},
		},
	}

	t.Run("Valid resource data", func(t *testing.T) {
		roles, actions, err := extractRolesAndActions(validResourceData)

		assert.NoError(t, err)
		assert.NotNil(t, roles)
		assert.NotNil(t, actions)
		assert.Contains(t, roles, "admin")
		assert.Contains(t, actions, "read")
	})

	t.Run("Missing roles", func(t *testing.T) {
		invalidData := map[string]interface{}{
			"actions": map[string]interface{}{},
		}

		roles, actions, err := extractRolesAndActions(invalidData)

		assert.Error(t, err)
		assert.Nil(t, roles)
		assert.Nil(t, actions)
		assert.Contains(t, err.Error(), "failed to extract roles")
	})

	t.Run("Missing actions", func(t *testing.T) {
		invalidData := map[string]interface{}{
			"roles": map[string]interface{}{},
		}

		roles, actions, err := extractRolesAndActions(invalidData)

		assert.Error(t, err)
		assert.Nil(t, roles)
		assert.Nil(t, actions)
		assert.Contains(t, err.Error(), "failed to extract actions")
	})

	t.Run("Invalid roles format", func(t *testing.T) {
		invalidData := map[string]interface{}{
			"roles":   "not-a-map",
			"actions": map[string]interface{}{},
		}

		roles, actions, err := extractRolesAndActions(invalidData)

		assert.Error(t, err)
		assert.Nil(t, roles)
		assert.Nil(t, actions)
		assert.Contains(t, err.Error(), "failed to extract roles")
	})

	t.Run("Invalid actions format", func(t *testing.T) {
		invalidData := map[string]interface{}{
			"roles":   map[string]interface{}{},
			"actions": "not-a-map",
		}

		roles, actions, err := extractRolesAndActions(invalidData)

		assert.Error(t, err)
		assert.Nil(t, roles)
		assert.Nil(t, actions)
		assert.Contains(t, err.Error(), "failed to extract actions")
	})
}

func TestFindRoleByID(t *testing.T) {
	roleID := uuid.New()
	nonExistentID := uuid.New()
	userID := uuid.New().String()
	resourceID := uuid.New().String()

	rolesData := map[string]interface{}{
		"admin": map[string]interface{}{
			"name":        "Admin",
			"description": "Administrator role",
			"key":         roleID.String(),
			"permissions": []interface{}{"read", "write"},
			"attributes": map[string]interface{}{
				"createdBy": userID,
				"updatedBy": userID,
				"createdAt": "2025-03-20T19:08:06-05:00",
				"updatedAt": "2025-03-20T19:08:06-05:00",
			},
		},
	}

	actionsData := map[string]interface{}{
		"read": map[string]interface{}{
			"id":   uuid.New().String(),
			"name": "Read",
		},
		"write": map[string]interface{}{
			"id":   uuid.New().String(),
			"name": "Write",
		},
	}

	resourceData := map[string]interface{}{
		"key":       resourceID,
		"name":      "Resource 1",
		"createdAt": "2025-03-20T19:08:06-05:00",
		"updatedAt": "2025-03-20T19:08:06-05:00",
		"actions": map[string]interface{}{
			"read": map[string]interface{}{
				"id":   uuid.New().String(),
				"name": "Read",
			},
		},
	}

	t.Run("Existing role ID", func(t *testing.T) {
		roles, err := findRoleByID(rolesData, actionsData, resourceData, roleID)

		assert.NoError(t, err)
		assert.NotNil(t, roles)
		assert.Equal(t, 1, len(roles))

		typedRole, ok := roles[0].(*models.Role)
		assert.True(t, ok)
		assert.Equal(t, roleID, typedRole.ID)
	})

	t.Run("Non-existent role ID", func(t *testing.T) {
		roles, err := findRoleByID(rolesData, actionsData, resourceData, nonExistentID)

		assert.Error(t, err)
		assert.Nil(t, roles)
		assert.Contains(t, err.Error(), "role not found")
	})

	t.Run("Invalid role data format", func(t *testing.T) {
		invalidRolesData := map[string]interface{}{
			"admin": "invalid-role-data", // Not a map
		}

		roles, err := findRoleByID(invalidRolesData, actionsData, resourceData, roleID)

		assert.Error(t, err)
		assert.Nil(t, roles)
		assert.Contains(t, err.Error(), "role not found")
	})

	t.Run("Role with invalid ID format", func(t *testing.T) {
		invalidRolesData := map[string]interface{}{
			"admin": map[string]interface{}{
				"name":        "Admin",
				"description": "Administrator role",
				"key":         "invalid-uuid", // Invalid UUID
				"permissions": []interface{}{"read", "write"},
			},
		}

		roles, err := findRoleByID(invalidRolesData, actionsData, resourceData, roleID)

		assert.Error(t, err)
		assert.Nil(t, roles)
		assert.Contains(t, err.Error(), "role not found")
	})
}

// Added missing test cases for edge scenarios and untested branches in helper functions.
func TestExtractDataFromResponse_AdditionalCases(t *testing.T) {
	invalidResponse := map[string]interface{}{
		"data": "invalid-data-format",
	}

	t.Run("Invalid data format", func(t *testing.T) {
		data, ok := extractDataFromResponse(invalidResponse)
		assert.False(t, ok)
		assert.Nil(t, data)
	})
}

func TestProcessPermission_AdditionalCases(t *testing.T) {
	invalidPermission := 123 // Non-string permission key
	validActionsData := map[string]interface{}{
		"read": map[string]interface{}{
			"id":   uuid.New().String(),
			"name": "Read",
		},
	}

	t.Run("Non-string permission key", func(t *testing.T) {
		permission, err := processPermission(invalidPermission, validActionsData)
		assert.Error(t, err)
		assert.Nil(t, permission)
	})
}
