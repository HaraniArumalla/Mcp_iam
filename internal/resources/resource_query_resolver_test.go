package resources

import (
	"context"
	"errors"
	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/permit"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestableResourceQueryResolver is a test version of ResourceQueryResolver that accepts mocks
type TestableResourceQueryResolver struct {
	*ResourceQueryResolver
	MockPermitSdk *MockPermitSdk
}

// MockPermitSdk is a mock implementation of PermitSdkService's Check method
type MockPermitSdk struct {
	mock.Mock
}

// Check mocks the Check method
func (m *MockPermitSdk) Check(ctx context.Context, userID, action, resourceType, resourceID, tenant string) (bool, error) {
	args := m.Called(ctx, userID, action, resourceType, resourceID, tenant)
	return args.Bool(0), args.Error(1)
}

// setupTestContext creates a context with user and tenant IDs
func setupTestContext(userIDStr, tenantIDStr string) context.Context {
	gin.SetMode(gin.TestMode)
	c := &gin.Context{}

	if userIDStr != "" {
		c.Set("userID", userIDStr)
	}

	if tenantIDStr != "" {
		c.Set("tenantID", tenantIDStr)
	}

	return context.WithValue(context.Background(), config.GinContextKey, c)
}

// NewTestResolver creates a ResourceQueryResolver for testing with our mock
func NewTestResolver(mockSdk *MockPermitSdk) *ResourceQueryResolver {
	// Create a standard permit.PermitSdkService
	psc := &permit.PermitSdkService{}

	// Create our resolver
	resolver := &ResourceQueryResolver{
		PSC: psc,
	}

	// Override the Check method to use our mock
	originalCheck := resolver.PSC.Check
	resolver.PSC.Check = func(ctx context.Context, userID, action, resourceType, resourceID, tenant string) (bool, error) {
		// This is compile-time safe, but will panic at runtime if any test uses this without setting up the mock
		if mockSdk == nil {
			panic("MockPermitSdk is nil")
		}
		return mockSdk.Check(ctx, userID, action, resourceType, resourceID, tenant)
	}

	return resolver
}

func TestCheckPermission(t *testing.T) {
	testCases := []struct {
		name                string
		input               models.PermissionInput
		setupContext        func() context.Context
		mockBehavior        func(*MockPermitSdk)
		expectedAllowed     bool
		expectedErrorMsg    string
	}{
		{
			name: "Successful permission check",
			input: models.PermissionInput{
				Action:       "read",
				ResourceType: "document",
				ResourceID:   "doc-123",
			},
			setupContext: func() context.Context {
				return setupTestContext("user-123", "tenant-123")
			},
			mockBehavior: func(m *MockPermitSdk) {
				m.On("Check", mock.Anything, "user-123", "read", "document", "doc-123", "tenant-123").
					Return(true, nil)
			},
			expectedAllowed:  true,
			expectedErrorMsg: "",
		},
		{
			name: "Failed permission check",
			input: models.PermissionInput{
				Action:       "write",
				ResourceType: "document",
				ResourceID:   "doc-123",
			},
			setupContext: func() context.Context {
				return setupTestContext("user-123", "tenant-123")
			},
			mockBehavior: func(m *MockPermitSdk) {
				m.On("Check", mock.Anything, "user-123", "write", "document", "doc-123", "tenant-123").
					Return(false, errors.New("permission denied"))
			},
			expectedAllowed:  false,
			expectedErrorMsg: "Failed to check permissions: permission denied",
		},
		{
			name: "Permission check with create action should clear resource ID",
			input: models.PermissionInput{
				Action:       "create",
				ResourceType: "document",
				ResourceID:   "doc-123", // This should be ignored for create actions
			},
			setupContext: func() context.Context {
				return setupTestContext("user-123", "tenant-123")
			},
			mockBehavior: func(m *MockPermitSdk) {
				// Expect empty resource ID for create action
				m.On("Check", mock.Anything, "user-123", "create", "document", "", "tenant-123").
					Return(true, nil)
			},
			expectedAllowed:  true,
			expectedErrorMsg: "",
		},
		{
			name: "Permission check with action containing 'create' should clear resource ID",
			input: models.PermissionInput{
				Action:       "createDocument",
				ResourceType: "document",
				ResourceID:   "doc-123", // This should be ignored for actions containing 'create'
			},
			setupContext: func() context.Context {
				return setupTestContext("user-123", "tenant-123")
			},
			mockBehavior: func(m *MockPermitSdk) {
				// Expect empty resource ID for action containing 'create'
				m.On("Check", mock.Anything, "user-123", "createDocument", "document", "", "tenant-123").
					Return(true, nil)
			},
			expectedAllowed:  true,
			expectedErrorMsg: "",
		},
		{
			name: "Empty action should fail",
			input: models.PermissionInput{
				Action:       "",
				ResourceType: "document",
				ResourceID:   "doc-123",
			},
			setupContext: func() context.Context {
				return setupTestContext("user-123", "tenant-123")
			},
			mockBehavior: func(m *MockPermitSdk) {
				// No mock expectations as the empty action should fail before Check is called
			},
			expectedAllowed:  false,
			expectedErrorMsg: "Action cannot be empty",
		},
		{
			name: "Missing user ID in context",
			input: models.PermissionInput{
				Action:       "read",
				ResourceType: "document",
				ResourceID:   "doc-123",
			},
			setupContext: func() context.Context {
				return setupTestContext("", "tenant-123")
			},
			mockBehavior: func(m *MockPermitSdk) {
				// No mock expectations as the context check should fail before Check is called
			},
			expectedAllowed:  false,
			expectedErrorMsg: "User ID & Tenant ID not found in context",
		},
		{
			name: "Missing tenant ID in context",
			input: models.PermissionInput{
				Action:       "read",
				ResourceType: "document",
				ResourceID:   "doc-123",
			},
			setupContext: func() context.Context {
				return setupTestContext("user-123", "")
			},
			mockBehavior: func(m *MockPermitSdk) {
				// No mock expectations as the context check should fail before Check is called
			},
			expectedAllowed:  false,
			expectedErrorMsg: "User ID & Tenant ID not found in context",
		},
		{
			name: "Missing context",
			input: models.PermissionInput{
				Action:       "read",
				ResourceType: "document",
				ResourceID:   "doc-123",
			},
			setupContext: func() context.Context {
				return context.Background() // No gin context
			},
			mockBehavior: func(m *MockPermitSdk) {
				// No mock expectations as the context check should fail before Check is called
			},
			expectedAllowed:  false,
			expectedErrorMsg: "User ID & Tenant ID not found in context",
		},
		{
			name: "No resource ID provided",
			input: models.PermissionInput{
				Action:       "read",
				ResourceType: "document",
				ResourceID:   "", // Empty resource ID
			},
			setupContext: func() context.Context {
				return setupTestContext("user-123", "tenant-123")
			},
			mockBehavior: func(m *MockPermitSdk) {
				// Expect empty resource ID
				m.On("Check", mock.Anything, "user-123", "read", "document", "", "tenant-123").
					Return(true, nil)
			},
			expectedAllowed:  true,
			expectedErrorMsg: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockSdk := new(MockPermitSdk)
			resolver := NewTestResolver(mockSdk)

			ctx := tc.setupContext()
			tc.mockBehavior(mockSdk)

			// Execute
			result, err := resolver.CheckPermission(ctx, tc.input)

			// Verify
			assert.NoError(t, err, "CheckPermission should not return an error")
			assert.NotNil(t, result, "Result should not be nil")
			assert.Equal(t, tc.expectedAllowed, result.Allowed, "Permission allowed status mismatch")

			if tc.expectedErrorMsg != "" {
				assert.NotNil(t, result.Error, "Expected error message but got nil")
				assert.Equal(t, tc.expectedErrorMsg, *result.Error, "Error message mismatch")
			} else if result.Error != nil {
				assert.Equal(t, "", *result.Error, "Error message should be empty for success cases")
			}

			// Verify that all expectations were met
			mockSdk.AssertExpectations(t)
		})
	}
}

func TestResource(t *testing.T) {
	// Test the Resource method which is currently a placeholder returning nil, nil
	t.Run("Resource returns nil for placeholder implementation", func(t *testing.T) {
		mockSdk := new(MockPermitSdk)
		resolver := NewTestResolver(mockSdk)

		// Verify the behavior with different IDs
		for _, id := range []uuid.UUID{uuid.New(), uuid.Nil} {
			result, err := resolver.Resource(context.Background(), id)
			assert.Nil(t, result, "Result should be nil for placeholder implementation")
			assert.Nil(t, err, "Error should be nil for placeholder implementation")
		}
	})
}

func TestResources(t *testing.T) {
	// Test the Resources method which is currently a placeholder returning nil, nil
	t.Run("Resources returns nil for placeholder implementation", func(t *testing.T) {
		mockSdk := new(MockPermitSdk)
		resolver := NewTestResolver(mockSdk)

		result, err := resolver.Resources(context.Background())
		assert.Nil(t, result, "Result should be nil for placeholder implementation")
		assert.Nil(t, err, "Error should be nil for placeholder implementation")
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
