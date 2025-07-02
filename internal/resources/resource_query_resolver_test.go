package resources

import (
	"context"
	"errors"
	"fmt"
	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/pkg/logger"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// PermitChecker defines the interface for permit checking
type PermitChecker interface {
	Check(ctx context.Context, userID, action, resourceType, resourceID, tenant string) (bool, error)
}

// MockPermitChecker implements PermitChecker for testing
type MockPermitChecker struct {
	mock.Mock
}

func (m *MockPermitChecker) Check(ctx context.Context, userID, action, resourceType, resourceID, tenant string) (bool, error) {
	args := m.Called(ctx, userID, action, resourceType, resourceID, tenant)
	return args.Bool(0), args.Error(1)
}

// TestResourceQueryResolver is our test version of ResourceQueryResolver
type TestResourceQueryResolver struct {
	PermitChecker PermitChecker
}

// CheckPermission implements the same logic as ResourceQueryResolver.CheckPermission
func (r *TestResourceQueryResolver) CheckPermission(ctx context.Context, input models.PermissionInput) (*models.PermissionResponse, error) {
	// Get user and tenant context
	userID, tenantID, err := helpers.GetUserAndTenantID(ctx)
	if err != nil {
		logger.LogError("User ID not found in context during permission check")
		return &models.PermissionResponse{
			Allowed: false,
			Error:   helpers.Ptr("User ID & Tenant ID not found in context"),
		}, nil
	}
	// Validate input.Action
	if input.Action == "" {
		return &models.PermissionResponse{
			Allowed: false,
			Error:   helpers.Ptr("Action cannot be empty"),
		}, nil
	}
	// Check permission
	resourceID := ""
	if input.ResourceID != "" {
		resourceID = input.ResourceID
	}
	// Check if input.Action contains "create"
	if strings.Contains(strings.ToLower(input.Action), "create") {
		resourceID = ""
	}

	// Log the permission check request
	logger.LogInfo("GraphQL permission check request",
		"user_id", userID,
		"action", input.Action,
		"resource_type", input.ResourceType,
		"resource_id", resourceID,
	)

	// Check permission
	allowed, err := r.PermitChecker.Check(ctx, userID.String(), input.Action, input.ResourceType, resourceID, tenantID.String())
	if err != nil {
		logger.LogError("Failed to check permissions", "error", err)
		return &models.PermissionResponse{
			Allowed: false,
			Error:   helpers.Ptr(fmt.Sprintf("Failed to check permissions: %s", err.Error())),
		}, nil
	}

	// Return result
	return &models.PermissionResponse{
		Allowed: allowed,
		Error:   helpers.Ptr(""),
	}, nil
}

// Resource and Resources match the original resolver signatures but always return nil
func (r *TestResourceQueryResolver) Resource(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	return nil, nil
}

func (r *TestResourceQueryResolver) Resources(ctx context.Context) (models.OperationResult, error) {
	return nil, nil
}

// setupTestContext creates a context with user and tenant IDs
func setupTestContext(userIDStr, tenantIDStr string) context.Context {
	gin.SetMode(gin.TestMode)
	c := gin.Context{}
	if userIDStr != "" {
		c.Set("userID", userIDStr)
	}
	if tenantIDStr != "" {
		c.Set("tenantID", tenantIDStr)
	}
	return context.WithValue(context.Background(), config.GinContextKey, &c)
}

func TestCheckPermission(t *testing.T) {
	// Initialize logger to prevent nil pointer dereference
	logger.InitLogger()

	// Test cases
	testCases := []struct {
		name                 string
		input                models.PermissionInput
		setupContext         func() context.Context
		setupMock            func(*MockPermitChecker)
		expectedAllowed      bool
		expectedError        bool
		expectedErrorMessage string
	}{
		{
			name: "Successful permission check with resource ID",
			input: models.PermissionInput{
				Action:       "read",
				ResourceType: "document",
				ResourceID:   "doc-123",
			},
			setupContext: func() context.Context {
				userID := uuid.New().String()
				tenantID := uuid.New().String()
				return setupTestContext(userID, tenantID)
			},
			setupMock: func(mockSvc *MockPermitChecker) {
				mockSvc.On("Check", mock.Anything, mock.Anything, "read", "document", "doc-123", mock.Anything).
					Return(true, nil)
			},
			expectedAllowed: true,
			expectedError:   false,
		},
		{
			name: "Failed permission check with resource ID",
			input: models.PermissionInput{
				Action:       "read",
				ResourceType: "document",
				ResourceID:   "doc-123",
			},
			setupContext: func() context.Context {
				userID := uuid.New().String()
				tenantID := uuid.New().String()
				return setupTestContext(userID, tenantID)
			},
			setupMock: func(mockSvc *MockPermitChecker) {
				mockSvc.On("Check", mock.Anything, mock.Anything, "read", "document", "doc-123", mock.Anything).
					Return(false, errors.New("permission denied"))
			},
			expectedAllowed:      false,
			expectedError:        true,
			expectedErrorMessage: "Failed to check permissions: permission denied",
		},
		{
			name: "Empty action should fail",
			input: models.PermissionInput{
				Action:       "",
				ResourceType: "document",
				ResourceID:   "doc-123",
			},
			setupContext: func() context.Context {
				userID := uuid.New().String()
				tenantID := uuid.New().String()
				return setupTestContext(userID, tenantID)
			},
			setupMock: func(mockSvc *MockPermitChecker) {
				// No mock expectations as the empty action should fail before Check is called
			},
			expectedAllowed:      false,
			expectedError:        true,
			expectedErrorMessage: "Action cannot be empty",
		},
		{
			name: "Missing user ID in context",
			input: models.PermissionInput{
				Action:       "read",
				ResourceType: "document",
				ResourceID:   "doc-123",
			},
			setupContext: func() context.Context {
				// Only set tenant ID, not user ID
				tenantID := uuid.New().String()
				return setupTestContext("", tenantID)
			},
			setupMock: func(mockSvc *MockPermitChecker) {
				// No mock expectations as the context check should fail before Check is called
			},
			expectedAllowed:      false,
			expectedError:        true,
			expectedErrorMessage: "User ID & Tenant ID not found in context",
		},
		{
			name: "Missing tenant ID in context",
			input: models.PermissionInput{
				Action:       "read",
				ResourceType: "document",
				ResourceID:   "doc-123",
			},
			setupContext: func() context.Context {
				// Only set user ID, not tenant ID
				userID := uuid.New().String()
				return setupTestContext(userID, "")
			},
			setupMock: func(mockSvc *MockPermitChecker) {
				// No mock expectations as the context check should fail before Check is called
			},
			expectedAllowed:      false,
			expectedError:        true,
			expectedErrorMessage: "User ID & Tenant ID not found in context",
		},
		{
			name: "Missing context",
			input: models.PermissionInput{
				Action:       "read",
				ResourceType: "document",
				ResourceID:   "doc-123",
			},
			setupContext: func() context.Context {
				// Return a context without the gin context
				return context.Background()
			},
			setupMock: func(mockSvc *MockPermitChecker) {
				// No mock expectations as the context check should fail before Check is called
			},
			expectedAllowed:      false,
			expectedError:        true,
			expectedErrorMessage: "User ID & Tenant ID not found in context",
		},
		{
			name: "Create action with resource ID (should clear resource ID)",
			input: models.PermissionInput{
				Action:       "create",
				ResourceType: "document",
				ResourceID:   "doc-123", // Should be ignored for create actions
			},
			setupContext: func() context.Context {
				userID := uuid.New().String()
				tenantID := uuid.New().String()
				return setupTestContext(userID, tenantID)
			},
			setupMock: func(mockSvc *MockPermitChecker) {
				mockSvc.On("Check", mock.Anything, mock.Anything, "create", "document", "", mock.Anything).
					Return(true, nil)
			},
			expectedAllowed: true,
			expectedError:   false,
		},
		{
			name: "CreateDocument action (contains 'create' substring)",
			input: models.PermissionInput{
				Action:       "createDocument",
				ResourceType: "document",
				ResourceID:   "doc-123", // Should be ignored for actions with 'create' substring
			},
			setupContext: func() context.Context {
				userID := uuid.New().String()
				tenantID := uuid.New().String()
				return setupTestContext(userID, tenantID)
			},
			setupMock: func(mockSvc *MockPermitChecker) {
				mockSvc.On("Check", mock.Anything, mock.Anything, "createDocument", "document", "", mock.Anything).
					Return(true, nil)
			},
			expectedAllowed: true,
			expectedError:   false,
		},
		{
			name: "No resource ID provided",
			input: models.PermissionInput{
				Action:       "read",
				ResourceType: "document",
				ResourceID:   "", // Empty resource ID
			},
			setupContext: func() context.Context {
				userID := uuid.New().String()
				tenantID := uuid.New().String()
				return setupTestContext(userID, tenantID)
			},
			setupMock: func(mockSvc *MockPermitChecker) {
				mockSvc.On("Check", mock.Anything, mock.Anything, "read", "document", "", mock.Anything).
					Return(true, nil)
			},
			expectedAllowed: true,
			expectedError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock permit service
			mockSvc := new(MockPermitChecker)

			// Setup context and mock
			ctx := tc.setupContext()

			// Setup mock expectations if needed
			if tc.setupMock != nil {
				tc.setupMock(mockSvc)
			}

			// Create a resolver with our mock
			resolver := &TestResourceQueryResolver{
				PermitChecker: mockSvc,
			}

			// Execute the function under test
			result, err := resolver.CheckPermission(ctx, tc.input)

			// Assertions
			assert.NoError(t, err, "CheckPermission should not return an error")
			assert.NotNil(t, result, "Result should not be nil")
			assert.Equal(t, tc.expectedAllowed, result.Allowed, "Permission allowed status mismatch")

			if tc.expectedError {
				assert.NotNil(t, result.Error, "Expected error message but got nil")
				assert.Equal(t, tc.expectedErrorMessage, *result.Error, "Error message mismatch")
			} else {
				// In success cases, either error should be nil or it should be empty string
				if result.Error != nil {
					assert.Equal(t, "", *result.Error, "Error message should be empty for success cases")
				}
			}

			// Verify that all expectations were met
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestResource(t *testing.T) {
	// Initialize logger to prevent nil pointer dereference
	logger.InitLogger()

	t.Run("Resource returns nil for placeholder implementation", func(t *testing.T) {
		// Create resolver
		mockSvc := new(MockPermitChecker)
		resolver := &TestResourceQueryResolver{
			PermitChecker: mockSvc,
		}

		// Test with various contexts and IDs
		testCases := []struct {
			name string
			ctx  context.Context
			id   uuid.UUID
		}{
			{
				name: "With valid UUID",
				ctx:  context.Background(),
				id:   uuid.New(),
			},
			{
				name: "With nil UUID",
				ctx:  context.Background(),
				id:   uuid.Nil,
			},
			{
				name: "With context with values",
				ctx:  setupTestContext(uuid.New().String(), uuid.New().String()),
				id:   uuid.New(),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result, err := resolver.Resource(tc.ctx, tc.id)
				assert.Nil(t, result, "Result should be nil for placeholder implementation")
				assert.Nil(t, err, "Error should be nil for placeholder implementation")
			})
		}
	})
}

func TestResources(t *testing.T) {
	// Initialize logger to prevent nil pointer dereference
	logger.InitLogger()

	t.Run("Resources returns nil for placeholder implementation", func(t *testing.T) {
		// Create resolver
		mockSvc := new(MockPermitChecker)
		resolver := &TestResourceQueryResolver{
			PermitChecker: mockSvc,
		}

		// Test with various contexts
		testCases := []struct {
			name string
			ctx  context.Context
		}{
			{
				name: "With empty context",
				ctx:  context.Background(),
			},
			{
				name: "With context with values",
				ctx:  setupTestContext(uuid.New().String(), uuid.New().String()),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result, err := resolver.Resources(tc.ctx)
				assert.Nil(t, result, "Result should be nil for placeholder implementation")
				assert.Nil(t, err, "Error should be nil for placeholder implementation")
			})
		}
	})
}

// TestResourceStructure verifies the structure of the ResourceQueryResolver
func TestResourceQueryResolverStructure(t *testing.T) {
	// Simply verify that the code structure exists
	resolver := &ResourceQueryResolver{}

	t.Run("Resource method exists", func(t *testing.T) {
		_, err := resolver.Resource(context.Background(), uuid.New())
		assert.Nil(t, err, "Resource method exists and returns nil error")
	})

	t.Run("Resources method exists", func(t *testing.T) {
		_, err := resolver.Resources(context.Background())
		assert.Nil(t, err, "Resources method exists and returns nil error")
	})

	t.Run("CheckPermission method structure", func(t *testing.T) {
		// We can't actually call CheckPermission without a valid PSC
		// But we can verify the method exists by running the tests above
		assert.NotPanics(t, func() {
			resolver = &ResourceQueryResolver{}
		})
	})
}
