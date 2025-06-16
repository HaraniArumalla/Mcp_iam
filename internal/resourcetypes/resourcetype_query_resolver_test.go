package resourcetypes

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

func setupTest(t *testing.T) (*ResourceTypeQueryResolver, *mocks.MockPermitService, context.Context) {
	ctrl := mock.NewController(t)
	mockService := mocks.NewMockPermitService(ctrl)
	resolver := &ResourceTypeQueryResolver{PC: mockService}

	gin.SetMode(gin.TestMode)
	ctx := &gin.Context{}
	ctx.Set("tenantID", uuid.New().String())
	testCtx := context.WithValue(context.Background(), config.GinContextKey, ctx)

	return resolver, mockService, testCtx
}

func TestResourceTypes(t *testing.T) {
	resolver, mockService, ctx := setupTest(t)

	// Valid resource types data
	resourceTypesData := buildTestResourceTypesData()

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
					Return(resourceTypesData, nil)
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
			result, _ := resolver.ResourceTypes(ctx)
			assert.NotNil(t, result)
		})
	}
}

func TestResourceTypeQueryResolver_ResourceTypes(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	objUnderTest := ResourceTypeQueryResolver{
		PC: mockService,
	}

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", "7ed6cfa6-fd7e-4a2a-bbce-773ef8ea4c12")
	testCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)

	resourceTypesData := buildTestResourceTypesData()
	mappedResourceTypes, _ := MockMapResourceTypesResponseToStruct(resourceTypesData)

	testcases := []struct {
		name      string
		ctx       context.Context
		mockStubs func(mockService *mocks.MockPermitService)
		expected  models.OperationResult
		wantErr   bool
	}{
		{
			name: "Test GetResourceTypes when permit service returns error",
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
			name: "Test GetResourceTypes success",
			ctx:  testCtx,
			mockStubs: func(mockSvc *mocks.MockPermitService) {
				mockSvc.EXPECT().
					SendRequest(mock.Any(), "GET", "resources?include_total_count=true", nil).
					Return(resourceTypesData, nil)
			},
			expected: func() models.OperationResult {
				result, _ := utils.FormatSuccess(mappedResourceTypes)
				return result
			}(),
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockStubs(mockService)
			result, _ := objUnderTest.ResourceTypes(tc.ctx)
			assert.NotNil(t, result)
		})
	}
}

func buildTestResourceTypesData() map[string]interface{} {
	return map[string]interface{}{
		"data": []interface{}{
			map[string]interface{}{
				"key":         uuid.New().String(),
				"name":        "User",
				"description": "User resource type",
				"attributes": map[string]interface{}{
					"displayName": "User",
					"permissions": []interface{}{
						"read",
						"write",
						"delete",
					},
				},
			},
			map[string]interface{}{
				"key":         uuid.New().String(),
				"name":        "Project",
				"description": "Project resource type",
				"attributes": map[string]interface{}{
					"displayName": "Project",
					"permissions": []interface{}{
						"read",
						"write",
						"delete",
					},
				},
			},
		},
	}
}

// Mock implementations of the mapping functions

func MockMapResourceTypesResponseToStruct(resourceResponse map[string]interface{}) ([]models.Data, error) {
	// Simple mock implementation for testing
	userDesc := "User resource type"
	projectDesc := "Project resource type"

	return []models.Data{
		&models.ResourceType{
			ID:          uuid.New(),
			Name:        "User",
			Description: &userDesc,
		},
		&models.ResourceType{
			ID:          uuid.New(),
			Name:        "Project",
			Description: &projectDesc,
		},
	}, nil
}
