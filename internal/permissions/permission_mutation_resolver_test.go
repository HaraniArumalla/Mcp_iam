package permissions

import (
	"context"
	"errors"
	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	mocks "iam_services_main_v1/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	mock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreatePermission(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	resolver := PermissionMutationResolver{
		PC: mockService,
	}

	// Setup test context with proper values
	userID := uuid.New()
	tenantID := uuid.New()

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", tenantID)
	ginCtx.Set("userID", userID)

	ginCtxNoUser := &gin.Context{}
	ginCtxNoUser.Set("tenantID", tenantID)

	ginCtxNoTenant := &gin.Context{}
	ginCtxNoTenant.Set("userID", userID)

	validCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)

	// Setup test input
	validInput := models.CreatePermissionInput{
		Name:               "test_permission",
		Description:        "Test description",
		AssignableScopeRef: uuid.New(),
	}

	invalidInputNoName := models.CreatePermissionInput{
		Name:               "",
		Description:        "Test description",
		AssignableScopeRef: uuid.New(),
	}

	invalidInputNoScope := models.CreatePermissionInput{
		Name:               "test_permission",
		Description:        "Test description",
		AssignableScopeRef: uuid.Nil,
	}

	// Mock permission data that will be returned after creation
	permissionData := buildTestPermissionData()

	testCases := []struct {
		name      string
		ctx       context.Context
		input     models.CreatePermissionInput
		mockSetup func()
		wantErr   bool
	}{
		{
			name:  "Invalid input - no name",
			ctx:   validCtx,
			input: invalidInputNoName,
			mockSetup: func() {
				// No mock calls expected - validation fails first
			},
			wantErr: true,
		},
		{
			name:  "Invalid input - no assignable scope",
			ctx:   validCtx,
			input: invalidInputNoScope,
			mockSetup: func() {
				// No mock calls expected - validation fails first
			},
			wantErr: true,
		},
		{
			name:  "Error creating permission in permit",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				expectedEndpoint := "resources/" + validInput.AssignableScopeRef.String() + "/actions"
				mockService.EXPECT().
					SendRequest(mock.Any(), "POST", expectedEndpoint, mock.Any()).
					Return(nil, errors.New("permit service error")).Times(1)
			},
			wantErr: true,
		},
		{
			name:  "Error getting permission details after creation",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				expectedEndpoint := "resources/" + validInput.AssignableScopeRef.String() + "/actions"
				mockService.EXPECT().
					SendRequest(mock.Any(), "POST", expectedEndpoint, mock.Any()).
					Return(map[string]interface{}{"key": validInput.AssignableScopeRef.String()}, nil).Times(1)

				getEndpoint := "resources/" + validInput.AssignableScopeRef.String()
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", getEndpoint, nil).
					Return(nil, errors.New("get permission error")).Times(1)
			},
			wantErr: true,
		},
		{
			name:  "Successful permission creation",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				expectedEndpoint := "resources/" + validInput.AssignableScopeRef.String() + "/actions"
				mockService.EXPECT().
					SendRequest(mock.Any(), "POST", expectedEndpoint, mock.Any()).
					Return(map[string]interface{}{"key": validInput.AssignableScopeRef.String()}, nil).Times(1)

				getEndpoint := "resources/" + validInput.AssignableScopeRef.String()
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", getEndpoint, nil).
					Return(permissionData, nil).Times(1)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			result, err := resolver.CreatePermission(tc.ctx, tc.input)

			// Check for errors based on the expected outcome
			if tc.wantErr {
				// For error cases, we expect either an error or a ResponseError (error response)
				if err == nil {
					// Check if result is a ResponseError (error response)
					if responseError, ok := result.(*models.ResponseError); ok {
						assert.False(t, responseError.IsSuccess)
					} else {
						t.Errorf("Expected error response but got: %T", result)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				// For success cases, we expect a SuccessResponse
				if successResponse, ok := result.(*models.SuccessResponse); ok {
					assert.True(t, successResponse.IsSuccess)
				}
			}

			if result == nil {
				t.Error("Result should not be nil")
			}
		})
	}
}

func TestValidateCreatePermissionInput(t *testing.T) {
	testCases := []struct {
		name    string
		input   models.CreatePermissionInput
		wantErr bool
	}{
		{
			name: "Valid input",
			input: models.CreatePermissionInput{
				Name:               "ValidPermission",
				Description:        "Valid description",
				AssignableScopeRef: uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "Empty name",
			input: models.CreatePermissionInput{
				Name:               "",
				Description:        "Valid description",
				AssignableScopeRef: uuid.New(),
			},
			wantErr: true,
		},
		{
			name: "Nil assignable scope",
			input: models.CreatePermissionInput{
				Name:               "ValidPermission",
				Description:        "Valid description",
				AssignableScopeRef: uuid.Nil,
			},
			wantErr: true,
		},
		{
			name: "Invalid name with special characters",
			input: models.CreatePermissionInput{
				Name:               "Invalid@Permission#",
				Description:        "Valid description",
				AssignableScopeRef: uuid.New(),
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateCreatePermissionInput(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreatePermissionInPermit(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	resolver := PermissionMutationResolver{
		PC: mockService,
	}

	ctx := context.Background()
	userID := uuid.New()
	tenantID := uuid.New()
	input := models.CreatePermissionInput{
		Name:               "Test Permission",
		Description:        "Test description",
		AssignableScopeRef: uuid.New(),
	}

	testCases := []struct {
		name      string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Successful creation",
			mockSetup: func() {
				expectedEndpoint := "resources/" + input.AssignableScopeRef.String() + "/actions"
				expectedPayload := map[string]interface{}{
					"description": input.Description,
					"key":         "test_permission",
					"name":        input.Name,
				}
				mockService.EXPECT().
					SendRequest(ctx, "POST", expectedEndpoint, expectedPayload).
					Return(map[string]interface{}{"key": "test_permission"}, nil).Times(1)
			},
			wantErr: false,
		},
		{
			name: "Permit service error",
			mockSetup: func() {
				expectedEndpoint := "resources/" + input.AssignableScopeRef.String() + "/actions"
				mockService.EXPECT().
					SendRequest(ctx, "POST", expectedEndpoint, mock.Any()).
					Return(nil, errors.New("permit service error")).Times(1)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			err := resolver.createPermissionInPermit(ctx, input, &userID, &tenantID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetPermissionDetailsById(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	resolver := PermissionMutationResolver{
		PC: mockService,
	}

	ctx := context.Background()
	resourceID := uuid.New()

	testCases := []struct {
		name      string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Successful retrieval",
			mockSetup: func() {
				expectedEndpoint := "resources/" + resourceID.String()
				permissionData := buildTestPermissionData()
				mockService.EXPECT().
					SendRequest(ctx, "GET", expectedEndpoint, nil).
					Return(permissionData, nil).Times(1)
			},
			wantErr: false,
		},
		{
			name: "Permit service error",
			mockSetup: func() {
				expectedEndpoint := "resources/" + resourceID.String()
				mockService.EXPECT().
					SendRequest(ctx, "GET", expectedEndpoint, nil).
					Return(nil, errors.New("permit service error")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Invalid permission data format",
			mockSetup: func() {
				expectedEndpoint := "resources/" + resourceID.String()
				invalidData := map[string]interface{}{
					"invalid": "data",
				}
				mockService.EXPECT().
					SendRequest(ctx, "GET", expectedEndpoint, nil).
					Return(invalidData, nil).Times(1)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			result, _ := resolver.getPermissionDetailsById(ctx, resourceID)

			assert.NotNil(t, result)
		})
	}
}

func TestMapToPermissionData(t *testing.T) {
	permissionID := uuid.New()
	actionID := uuid.New()

	testCases := []struct {
		name     string
		input    map[string]interface{}
		wantErr  bool
		expected *models.ResourceType
	}{
		{
			name: "Valid permission data",
			input: map[string]interface{}{
				"key":       permissionID.String(),
				"name":      "Test Resource Type",
				"createdAt": "2025-07-14T10:00:00Z",
				"updatedAt": "2025-07-14T10:00:00Z",
				"actions": map[string]interface{}{
					"read": map[string]interface{}{
						"id":          actionID.String(),
						"description": "Read permission",
						"createdAt":   "2025-07-14T10:00:00Z",
						"updatedAt":   "2025-07-14T10:00:00Z",
					},
				},
			},
			wantErr: false,
			expected: &models.ResourceType{
				ID:        permissionID,
				Name:      "Test Resource Type",
				CreatedAt: "2025-07-14T10:00:00Z",
				UpdatedAt: "2025-07-14T10:00:00Z",
			},
		},
		{
			name: "Invalid key format",
			input: map[string]interface{}{
				"key":  "invalid-uuid",
				"name": "Test Resource Type",
				"actions": map[string]interface{}{
					"read": map[string]interface{}{
						"id": actionID.String(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Missing actions",
			input: map[string]interface{}{
				"key":  permissionID.String(),
				"name": "Test Resource Type",
			},
			wantErr: true,
		},
		{
			name: "Invalid actions format",
			input: map[string]interface{}{
				"key":     permissionID.String(),
				"name":    "Test Resource Type",
				"actions": "invalid-actions-format",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MapToPermissionData(tc.input)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expected.ID, result.ID)
				assert.Equal(t, tc.expected.Name, result.Name)
				assert.Equal(t, tc.expected.CreatedAt, result.CreatedAt)
				assert.Equal(t, tc.expected.UpdatedAt, result.UpdatedAt)
				assert.NotNil(t, result.Permissions)
			}
		})
	}
}

func TestMapToPermissionActions(t *testing.T) {
	resourceID := uuid.New()
	actionID1 := uuid.New()
	actionID2 := uuid.New()

	testCases := []struct {
		name        string
		resourceID  uuid.UUID
		actionsData map[string]interface{}
		expected    int // expected number of permissions
	}{
		{
			name:       "Valid actions data",
			resourceID: resourceID,
			actionsData: map[string]interface{}{
				"read": map[string]interface{}{
					"id":          actionID1.String(),
					"description": "Read permission",
					"createdAt":   "2025-07-14T10:00:00Z",
					"updatedAt":   "2025-07-14T10:00:00Z",
				},
				"write": map[string]interface{}{
					"id":          actionID2.String(),
					"description": "Write permission",
					"createdAt":   "2025-07-14T10:00:00Z",
					"updatedAt":   "2025-07-14T10:00:00Z",
				},
			},
			expected: 2,
		},
		{
			name:       "Invalid action data format",
			resourceID: resourceID,
			actionsData: map[string]interface{}{
				"read":  "invalid-format",
				"write": map[string]interface{}{"id": actionID2.String()},
			},
			expected: 1,
		},
		{
			name:       "Invalid action ID",
			resourceID: resourceID,
			actionsData: map[string]interface{}{
				"read": map[string]interface{}{
					"id": "invalid-uuid",
				},
			},
			expected: 0,
		},
		{
			name:        "Empty actions data",
			resourceID:  resourceID,
			actionsData: map[string]interface{}{},
			expected:    0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := MapToPermissionActions(tc.resourceID, tc.actionsData)
			assert.Len(t, result, tc.expected)

			for _, permission := range result {
				assert.NotEqual(t, uuid.Nil, permission.ID)
				assert.NotEmpty(t, permission.Name)
			}
		})
	}
}

// Helper function to build test permission data
func buildTestPermissionData() map[string]interface{} {
	permissionID := uuid.New()
	actionID := uuid.New()

	return map[string]interface{}{
		"key":       permissionID.String(),
		"name":      "Test Resource Type",
		"createdAt": "2025-07-14T10:00:00Z",
		"updatedAt": "2025-07-14T10:00:00Z",
		"actions": map[string]interface{}{
			"read": map[string]interface{}{
				"id":          actionID.String(),
				"description": "Read permission",
				"createdAt":   "2025-07-14T10:00:00Z",
				"updatedAt":   "2025-07-14T10:00:00Z",
			},
		},
	}
}
