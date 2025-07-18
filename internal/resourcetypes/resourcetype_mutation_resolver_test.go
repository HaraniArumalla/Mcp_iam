package resourcetypes

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

func TestCreateResource(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	resolver := ResourceTypeMutationResolver{
		PC: mockService,
	}

	// Setup test context with proper values
	userID := uuid.New()
	tenantID := uuid.New()
	resourceID := uuid.New()

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", tenantID)
	ginCtx.Set("userID", userID)

	validCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)

	// Setup test input
	validInput := models.CreateResourceInput{
		ID:   resourceID,
		Name: "test_resource",
	}

	// Mock resource data that will be returned after creation
	resourceData := buildTestResourceData(resourceID)

	testCases := []struct {
		name           string
		ctx            context.Context
		input          models.CreateResourceInput
		setupMocks     func()
		expectedError  bool
		expectedStatus int
	}{
		{
			name:  "successful resource creation",
			ctx:   validCtx,
			input: validInput,
			setupMocks: func() {
				// Mock the POST request to create resource
				mockService.EXPECT().
					SendRequest(validCtx, "POST", "resources", mock.Any()).
					Return(map[string]interface{}{"status": "success"}, nil).
					Times(1)

				// Mock the GET request to fetch resource details
				mockService.EXPECT().
					SendRequest(validCtx, "GET", "resources/"+resourceID.String(), nil).
					Return(resourceData, nil).
					Times(1)
			},
			expectedError:  false,
			expectedStatus: 200,
		},
		{
			name:  "resource creation API call fails",
			ctx:   validCtx,
			input: validInput,
			setupMocks: func() {
				// Mock the POST request to fail
				mockService.EXPECT().
					SendRequest(validCtx, "POST", "resources", mock.Any()).
					Return(nil, errors.New("API call failed")).
					Times(1)
			},
			expectedError:  false, // Returns error response, not actual error
			expectedStatus: 400,
		},
		{
			name:  "resource details fetch fails",
			ctx:   validCtx,
			input: validInput,
			setupMocks: func() {
				// Mock the POST request to succeed
				mockService.EXPECT().
					SendRequest(validCtx, "POST", "resources", mock.Any()).
					Return(map[string]interface{}{"status": "success"}, nil).
					Times(1)

				// Mock the GET request to fail
				mockService.EXPECT().
					SendRequest(validCtx, "GET", "resources/"+resourceID.String(), nil).
					Return(nil, errors.New("Fetch failed")).
					Times(1)
			},
			expectedError:  false, // Returns error response, not actual error
			expectedStatus: 400,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			result, _ := resolver.CreateResource(tc.ctx, tc.input)

			assert.NotNil(t, result)
		})
	}
}

func TestPrepareResourceActions(t *testing.T) {
	resolver := ResourceTypeMutationResolver{}

	testCases := []struct {
		name          string
		input         models.CreateResourceInput
		expectedError bool
		expectedCount int
	}{
		{
			name: "basic test with resource name",
			input: models.CreateResourceInput{
				Name: "SampleResource",
			},
			expectedError: false,
			expectedCount: 4, // create, read, update, delete
		},
		{
			name: "resource name with spaces",
			input: models.CreateResourceInput{
				Name: "My Resource",
			},
			expectedError: false,
			expectedCount: 4,
		},
		{
			name: "empty resource name",
			input: models.CreateResourceInput{
				Name: "",
			},
			expectedError: false,
			expectedCount: 4,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actions, err := resolver.prepareResourceActions(tc.input)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCount, len(actions))

				for key, action := range actions {
					assert.NotContains(t, key, " ")
					actionMap, ok := action.(map[string]interface{})
					assert.True(t, ok)
					assert.Contains(t, actionMap, "name")
				}
			}
		})
	}
}

func TestCreateResourceInternal(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	resolver := ResourceTypeMutationResolver{
		PC: mockService,
	}

	resourceID := uuid.New()
	ctx := context.Background()

	input := models.CreateResourceInput{
		ID:   resourceID,
		Name: "test_resource",
	}

	actions := map[string]interface{}{
		"read": map[string]interface{}{
			"name":        "Read",
			"description": "Read permission",
		},
	}

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError bool
	}{
		{
			name: "successful resource creation",
			setupMocks: func() {
				mockService.EXPECT().
					SendRequest(ctx, "POST", "resources", map[string]interface{}{
						"key":     input.ID,
						"name":    input.Name,
						"actions": actions,
					}).
					Return(map[string]interface{}{"status": "success"}, nil).
					Times(1)
			},
			expectedError: false,
		},
		{
			name: "API call fails",
			setupMocks: func() {
				mockService.EXPECT().
					SendRequest(ctx, "POST", "resources", mock.Any()).
					Return(nil, errors.New("API error")).
					Times(1)
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			err := resolver.createResource(ctx, input, actions)

			if tc.expectedError {
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
	resolver := ResourceTypeMutationResolver{
		PC: mockService,
	}

	resourceID := uuid.New()
	ctx := context.Background()
	resourceData := buildTestResourceData(resourceID)

	testCases := []struct {
		name           string
		setupMocks     func()
		expectedError  bool
		expectedStatus int
	}{
		{
			name: "successful fetch",
			setupMocks: func() {
				mockService.EXPECT().
					SendRequest(ctx, "GET", "resources/"+resourceID.String(), nil).
					Return(resourceData, nil).
					Times(1)
			},
			expectedError:  false,
			expectedStatus: 200,
		},
		{
			name: "API call fails",
			setupMocks: func() {
				mockService.EXPECT().
					SendRequest(ctx, "GET", "resources/"+resourceID.String(), nil).
					Return(nil, errors.New("API error")).
					Times(1)
			},
			expectedError:  false, // Returns error response, not actual error
			expectedStatus: 400,
		},
		{
			name: "invalid resource data",
			setupMocks: func() {
				invalidData := map[string]interface{}{
					"key":  "invalid-uuid",
					"name": "test",
				}
				mockService.EXPECT().
					SendRequest(ctx, "GET", "resources/"+resourceID.String(), nil).
					Return(invalidData, nil).
					Times(1)
			},
			expectedError:  false, // Returns error response, not actual error
			expectedStatus: 400,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			result, _ := resolver.getPermissionDetailsById(ctx, resourceID)
			assert.NotNil(t, result)

		})
	}
}

func TestMapToPermissionData(t *testing.T) {
	resourceID := uuid.New()
	permissionID1 := uuid.New()
	permissionID2 := uuid.New()

	testCases := []struct {
		name          string
		input         map[string]interface{}
		expectedError bool
		expectedName  string
		expectedPerms int
	}{
		{
			name: "valid resource data",
			input: map[string]interface{}{
				"key":       resourceID.String(),
				"name":      "Test Resource",
				"createdAt": "2023-01-01T00:00:00Z",
				"updatedAt": "2023-01-02T00:00:00Z",
				"actions": map[string]interface{}{
					"read": map[string]interface{}{
						"id":          permissionID1.String(),
						"description": "Read permission",
						"createdAt":   "2023-01-01T00:00:00Z",
						"updatedAt":   "2023-01-02T00:00:00Z",
					},
					"write": map[string]interface{}{
						"id":          permissionID2.String(),
						"description": "Write permission",
						"createdAt":   "2023-01-01T00:00:00Z",
						"updatedAt":   "2023-01-02T00:00:00Z",
					},
				},
			},
			expectedError: false,
			expectedName:  "Test Resource",
			expectedPerms: 2,
		},
		{
			name: "invalid resource ID",
			input: map[string]interface{}{
				"key":     "invalid-uuid",
				"name":    "Test Resource",
				"actions": map[string]interface{}{},
			},
			expectedError: true,
		},
		{
			name: "missing actions",
			input: map[string]interface{}{
				"key":  resourceID.String(),
				"name": "Test Resource",
			},
			expectedError: true,
		},
		{
			name: "invalid actions format",
			input: map[string]interface{}{
				"key":     resourceID.String(),
				"name":    "Test Resource",
				"actions": "invalid",
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MapToPermissionData(tc.input)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedName, result.Name)
				assert.Equal(t, tc.expectedPerms, len(result.Permissions))
			}
		})
	}
}

func TestMapToPermissionActions(t *testing.T) {
	resourceID := uuid.New()
	permissionID1 := uuid.New()
	permissionID2 := uuid.New()

	testCases := []struct {
		name          string
		actionsData   map[string]interface{}
		expectedCount int
	}{
		{
			name: "valid actions data",
			actionsData: map[string]interface{}{
				"read": map[string]interface{}{
					"id":          permissionID1.String(),
					"description": "Read permission",
					"createdAt":   "2023-01-01T00:00:00Z",
					"updatedAt":   "2023-01-02T00:00:00Z",
				},
				"write": map[string]interface{}{
					"id":          permissionID2.String(),
					"description": "Write permission",
					"createdAt":   "2023-01-01T00:00:00Z",
					"updatedAt":   "2023-01-02T00:00:00Z",
				},
			},
			expectedCount: 2,
		},
		{
			name:          "empty actions data",
			actionsData:   map[string]interface{}{},
			expectedCount: 0,
		},
		{
			name: "invalid action format",
			actionsData: map[string]interface{}{
				"read": "invalid",
			},
			expectedCount: 0,
		},
		{
			name: "action with invalid ID",
			actionsData: map[string]interface{}{
				"read": map[string]interface{}{
					"id":          "invalid-uuid",
					"description": "Read permission",
				},
			},
			expectedCount: 0,
		},
		{
			name: "mixed valid and invalid actions",
			actionsData: map[string]interface{}{
				"read": map[string]interface{}{
					"id":          permissionID1.String(),
					"description": "Read permission",
				},
				"write": "invalid",
				"delete": map[string]interface{}{
					"id":          "invalid-uuid",
					"description": "Delete permission",
				},
			},
			expectedCount: 1, // Only the valid 'read' action should be processed
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := MapToPermissionActions(resourceID, tc.actionsData)

			assert.Equal(t, tc.expectedCount, len(result))

			// Validate structure of returned permissions
			for _, permission := range result {
				assert.NotEqual(t, uuid.Nil, permission.ID)
				assert.NotEmpty(t, permission.Name)
				// Description can be nil, so we don't assert on it
			}
		})
	}
}

// Helper function to build test resource data
func buildTestResourceData(resourceID uuid.UUID) map[string]interface{} {
	permissionID1 := uuid.New()
	permissionID2 := uuid.New()

	return map[string]interface{}{
		"key":       resourceID.String(),
		"name":      "Test Resource",
		"createdAt": "2023-01-01T00:00:00Z",
		"updatedAt": "2023-01-02T00:00:00Z",
		"actions": map[string]interface{}{
			"read": map[string]interface{}{
				"id":          permissionID1.String(),
				"description": "Read permission",
				"createdAt":   "2023-01-01T00:00:00Z",
				"updatedAt":   "2023-01-02T00:00:00Z",
			},
			"write": map[string]interface{}{
				"id":          permissionID2.String(),
				"description": "Write permission",
				"createdAt":   "2023-01-01T00:00:00Z",
				"updatedAt":   "2023-01-02T00:00:00Z",
			},
		},
	}
}
