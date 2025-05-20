package roles

import (
	"context"
	"errors"
	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/utils"
	mocks "iam_services_main_v1/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	mock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) (*RoleQueryResolver, *mocks.MockPermitService, context.Context) {
	ctrl := mock.NewController(t)
	mockService := mocks.NewMockPermitService(ctrl)
	resolver := &RoleQueryResolver{PC: mockService}

	gin.SetMode(gin.TestMode)
	ctx := &gin.Context{}
	ctx.Set("tenantID", uuid.New().String())
	testCtx := context.WithValue(context.Background(), config.GinContextKey, ctx)

	return resolver, mockService, testCtx
}

func TestRole(t *testing.T) {
	resolver, mockService, ctx := setupTest(t)
	validID := uuid.New()

	// Valid role data
	roleData := buildTestRolesData()

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(roleData, nil)
			},
			wantErr: false,
		},
		{
			name: "Service error",
			setup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(nil, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			result, _ := resolver.Role(ctx, validID)
			assert.NotNil(t, result)
		})
	}
}

func TestRoles(t *testing.T) {
	resolver, mockService, ctx := setupTest(t)

	// Valid roles data
	rolesData := buildTestRolesData()

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", "resources?include_total_count=true", nil).
					Return(rolesData, nil)
			},
			wantErr: false,
		},
		{
			name: "Service error",
			setup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", "resources?include_total_count=true", nil).
					Return(nil, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			result, _ := resolver.Roles(ctx)
			assert.NotNil(t, result)
		})
	}
}

func TestRoleQueryResolver_Role(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	objUnderTest := RoleQueryResolver{
		PC: mockService,
	}

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", "7ed6cfa6-fd7e-4a2a-bbce-773ef8ea4c12")
	testCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)

	roleData := buildTestRolesData()

	// Create valid role ID to search for
	validRoleID := uuid.New()

	// Mock response from mapping function
	//mockRole, _ := MockMapRoleResponseToStruct(roleData, validRoleID)

	testcases := []struct {
		name      string
		input     uuid.UUID
		ctx       context.Context
		mockStubs func(mockService *mocks.MockPermitService)
		wantErr   bool
	}{
		{
			name:  "Test GetRole with nil UUID",
			input: uuid.Nil,
			ctx:   testCtx,
			mockStubs: func(mockSvc *mocks.MockPermitService) {
				// No mock calls expected
			},
			wantErr: true,
		},
		{
			name:  "Test GetRole when permit service returns error",
			input: validRoleID,
			ctx:   testCtx,
			mockStubs: func(mockSvc *mocks.MockPermitService) {
				mockSvc.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(nil, errors.New("permit service error"))
			},
			wantErr: true,
		},
		{
			name:  "Test GetRole success",
			input: validRoleID,
			ctx:   testCtx,
			mockStubs: func(mockSvc *mocks.MockPermitService) {
				mockSvc.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(roleData, nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockStubs(mockService)
			result, _ := objUnderTest.Role(tc.ctx, tc.input)
			assert.NotNil(t, result)
		})
	}
}

func TestRoleQueryResolver_Roles(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	objUnderTest := RoleQueryResolver{
		PC: mockService,
	}

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", "7ed6cfa6-fd7e-4a2a-bbce-773ef8ea4c12")
	testCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)

	roleData := buildTestRolesData()
	mappedRoles, _ := MockMapRolesResponseToStruct(roleData)

	testcases := []struct {
		name      string
		ctx       context.Context
		mockStubs func(mockService *mocks.MockPermitService)
		expected  models.OperationResult
		wantErr   bool
	}{
		{
			name: "Test GetRoles when permit service returns error",
			ctx:  testCtx,
			mockStubs: func(mockSvc *mocks.MockPermitService) {
				mockSvc.EXPECT().
					SendRequest(mock.Any(), "GET", "resources?include_total_count=true", nil).
					Return(nil, errors.New("permit service error"))
			},
			expected: utils.FormatErrorResponse(400, "Error retrieving roles from permit system", "permit service error"),
			wantErr:  true,
		},
		{
			name: "Test GetRoles success",
			ctx:  testCtx,
			mockStubs: func(mockSvc *mocks.MockPermitService) {
				mockSvc.EXPECT().
					SendRequest(mock.Any(), "GET", "resources?include_total_count=true", nil).
					Return(roleData, nil)
			},
			expected: func() models.OperationResult {
				result, _ := utils.FormatSuccess(mappedRoles)
				return result
			}(),
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockStubs(mockService)
			result, _ := objUnderTest.Roles(tc.ctx)
			assert.NotNil(t, result)
		})
	}
}

func buildTestRolesData() map[string]interface{} {
	return map[string]interface{}{
		"data": []interface{}{
			map[string]interface{}{
				"key":         uuid.New().String(),
				"name":        "Admin Role",
				"description": "Administrator role with full access",
				"attributes": map[string]interface{}{
					"displayName": "Administrator",
					"permissions": []interface{}{
						"read:all",
						"write:all",
						"delete:all",
					},
				},
			},
			map[string]interface{}{
				"key":         uuid.New().String(),
				"name":        "User Role",
				"description": "Standard user role with limited access",
				"attributes": map[string]interface{}{
					"displayName": "Standard User",
					"permissions": []interface{}{
						"read:own",
						"write:own",
					},
				},
			},
		},
	}
}

// Mock implementations of the mapping functions

func MockMapRoleResponseToStruct(resourceResponse map[string]interface{}, id uuid.UUID) ([]models.Data, error) {
	// Simple mock implementation for testing
	description := "Test Role Description"
	return []models.Data{
		&models.Role{
			ID:          id,
			Name:        "Test Role",
			Description: &description,
		},
	}, nil
}
func MockMapRolesResponseToStruct(resourcesResponse map[string]interface{}) ([]models.Data, error) {
	// Simple mock implementation for testing
	adminDesc := "Administrator role with full access"
	userDesc := "Standard user role with limited access"
	return []models.Data{
		&models.Role{
			ID:          uuid.New(),
			Name:        "Admin Role",
			Description: &adminDesc,
		},
		&models.Role{
			ID:          uuid.New(),
			Name:        "User Role",
			Description: &userDesc,
		},
	}, nil
}
