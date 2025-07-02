package middlewares

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPermitSdkService is a mock of the PermitSdkService
type MockPermitSdkService struct {
	mock.Mock
}

func (m *MockPermitSdkService) Check(ctx context.Context, user, action, resourceType, resourceID, tenant string) (bool, error) {
	args := m.Called(ctx, user, action, resourceType, resourceID, tenant)
	return args.Bool(0), args.Error(1)
}

func (m *MockPermitSdkService) RequirePermission(action, resourceType string) gin.HandlerFunc {
	args := m.Called(action, resourceType)
	return args.Get(0).(gin.HandlerFunc)
}

func (m *MockPermitSdkService) RequirePermissionWithParamsAndContext(ctx context.Context, action, resourceType, resourceIDParam, tenantParam string) gin.HandlerFunc {
	args := m.Called(ctx, action, resourceType, resourceIDParam, tenantParam)
	return args.Get(0).(gin.HandlerFunc)
}

func (m *MockPermitSdkService) RequirePermissionWithContext(ctx context.Context, action, resourceType string) gin.HandlerFunc {
	args := m.Called(ctx, action, resourceType)
	return args.Get(0).(gin.HandlerFunc)
}

func (m *MockPermitSdkService) RequirePermissionWithParams(action, resourceType, resourceIDParam, tenantParam string) gin.HandlerFunc {
	args := m.Called(action, resourceType, resourceIDParam, tenantParam)
	return args.Get(0).(gin.HandlerFunc)
}

func (m *MockPermitSdkService) SomeMissingMethod(ctx context.Context, param string) error {
	args := m.Called(ctx, param)
	return args.Error(0)
}

func TestExtractAction(t *testing.T) {
	tests := []struct {
		name     string
		request  graphQLRequest
		expected string
	}{
		{
			name: "Extract from operationName",
			request: graphQLRequest{
				OperationName: "createAccount",
				Query:         "mutation createAccount { ... }",
			},
			expected: "createAccount",
		},
		{
			name: "Extract from query with mutation",
			request: graphQLRequest{
				OperationName: "",
				Query: `mutation {
					createTenant(input: {name: "Test"}) {
						id
						name
					}
				}`,
			},
			expected: "createTenant",
		},
		{
			name: "Extract from query with query keyword",
			request: graphQLRequest{
				OperationName: "",
				Query: `query {
					tenant(id: "123") {
						id
						name
					}
				}`,
			},
			expected: "tenant",
		},
		{
			name: "Empty query",
			request: graphQLRequest{
				OperationName: "",
				Query:         "",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractAction(tt.request)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDeriveResourceType(t *testing.T) {
	tests := []struct {
		name     string
		action   string
		expected string
	}{
		{
			name:     "Tenant create action",
			action:   "createTenant",
			expected: "ed113bda-bbda-11ef-87ea-c03c5946f955",
		},
		{
			name:     "Tenant query action",
			action:   "tenant",
			expected: "ed113bda-bbda-11ef-87ea-c03c5946f955",
		},
		{
			name:     "Account create action",
			action:   "createAccount",
			expected: "ed113f30-bbda-11ef-87ea-c03c5946f955",
		},
		{
			name:     "Account query action",
			action:   "account",
			expected: "ed113f30-bbda-11ef-87ea-c03c5946f955",
		},
		{
			name:     "Unknown action",
			action:   "unknownAction",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deriveResourceType(tt.action)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		val      string
		expected bool
	}{
		{
			name:     "Value exists in slice",
			slice:    []string{"createTenant", "updateTenant", "deleteTenant"},
			val:      "createTenant",
			expected: true,
		},
		{
			name:     "Value exists in slice case insensitive",
			slice:    []string{"createTenant", "updateTenant", "deleteTenant"},
			val:      "CREATETENANT",
			expected: true,
		},
		{
			name:     "Value does not exist in slice",
			slice:    []string{"createTenant", "updateTenant", "deleteTenant"},
			val:      "createAccount",
			expected: false,
		},
		{
			name:     "Empty slice",
			slice:    []string{},
			val:      "createTenant",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.slice, tt.val)
			assert.Equal(t, tt.expected, result)
		})
	}
}
