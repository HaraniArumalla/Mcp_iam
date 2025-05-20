package tenants

import (
	"context"
	"errors"
	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/utils"
	mocks "iam_services_main_v1/mocks"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	mock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) (*TenantQueryResolver, *mocks.MockPermitService, context.Context) {
	ctrl := mock.NewController(t)
	mockService := mocks.NewMockPermitService(ctrl)
	resolver := &TenantQueryResolver{PC: mockService}

	gin.SetMode(gin.TestMode)
	ctx := &gin.Context{}
	ctx.Set("tenantID", uuid.New().String())
	testCtx := context.WithValue(context.Background(), config.GinContextKey, ctx)

	return resolver, mockService, testCtx
}

func TestTenant(t *testing.T) {
	resolver, mockService, ctx := setupTest(t)
	validID := uuid.New()

	// Valid tenant data
	tenantData := map[string]interface{}{
		"key": validID.String(),
		"attributes": map[string]interface{}{
			"name": "test tenant",
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
		},
	}

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
					Return(tenantData, nil)
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
			result, err := resolver.Tenant(ctx, validID)
			if tt.wantErr {
				assert.NotNil(t, result) // Error response should still be returned
			} else {
				assert.NotNil(t, result)
				assert.Nil(t, err)
			}
		})
	}
}

func TestTenants(t *testing.T) {
	resolver, mockService, ctx := setupTest(t)

	// Valid tenants data
	tenantsData := map[string]interface{}{
		"data": []interface{}{
			map[string]interface{}{
				"key": uuid.New().String(),
				"attributes": map[string]interface{}{
					"name": "test tenant",
					"contactInfo": map[string]interface{}{
						"email":       "test@example.com",
						"phoneNumber": "1234567890",
						"address": map[string]interface{}{
							"street":  "123 Test St",
							"city":    "Test City",
							"state":   "TS",
							"zipcode": "12345",
							"country": "Test Country",
						},
					},
				},
			},
		},
	}

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
					Return(tenantsData, nil)
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
			result, err := resolver.Tenants(ctx)
			if tt.wantErr {
				assert.NotNil(t, result) // Error response should still be returned
			} else {
				assert.NotNil(t, result)
				assert.Nil(t, err)
			}
		})
	}
}

func TestTenantQueryResolver_Tenant(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	objUnderTest := TenantQueryResolver{
		PC: mockService,
	}

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", "7ed6cfa6-fd7e-4a2a-bbce-773ef8ea4c12")
	testCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)

	tenant := buildTestTenantsData()
	mappedTenant, _ := MapTenantResponseToStruct(tenant)

	testcases := []struct {
		name      string
		input     uuid.UUID
		ctx       context.Context
		mockStubs func(mockService mocks.MockPermitService)
		output    models.OperationResult
	}{
		{
			name:  "Test GetTenant with nil UUID",
			input: uuid.Nil,
			ctx:   testCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				// No mock calls expected
			},
			output: utils.FormatErrorResponse(400, "invalid tenant ID", "invalid tenant ID: 00000000-0000-0000-0000-000000000000"),
		},
		{
			name:  "Test GetTenant when permit service returns error",
			input: uuid.New(),
			ctx:   testCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(nil, errors.New("permit service error"))
			},
			output: utils.FormatErrorResponse(400, "Failed to get tenant resources from permit", "permit service error"),
		},
		{
			name:  "Test GetTenant success",
			input: uuid.New(),
			ctx:   testCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(tenant, nil)
			},
			output: func() models.OperationResult {
				result, _ := utils.FormatSuccessResponse(mappedTenant)
				return result
			}(),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockStubs(*mockService)
			result, _ := objUnderTest.Tenant(tc.ctx, tc.input)
			assert.NotNil(t, result)
		})
	}
}

func TestTenantQueryResolver_Tenants(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	objUnderTest := TenantQueryResolver{
		PC: mockService,
	}

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", "7ed6cfa6-fd7e-4a2a-bbce-773ef8ea4c12")
	testCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)

	tenant := buildTestTenantsData()
	tenants := map[string]interface{}{
		"data": []interface{}{tenant},
	}
	mappedTenants, _ := MapTenantsResponseToStruct(tenants)

	testcases := []struct {
		name      string
		ctx       context.Context
		mockStubs func(mockService mocks.MockPermitService)
		output    models.OperationResult
	}{
		{
			name: "Test GetTenants when permit service returns error",
			ctx:  testCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(nil, errors.New("permit service error"))
			},
			output: utils.FormatErrorResponse(400, "Failed to get tenant resources from permit", "permit service error"),
		},
		{
			name: "Test GetTenants success",
			ctx:  testCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(tenants, nil)
			},
			output: func() models.OperationResult {
				result, _ := utils.FormatSuccessResponse(mappedTenants)
				return result
			}(),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockStubs(*mockService)
			result, _ := objUnderTest.Tenants(tc.ctx)
			assert.NotNil(t, result)
		})
	}
}

func TestTenantQueryResolver_FetchTenant(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	resolver := TenantQueryResolver{PC: mockService}

	testID := uuid.New()
	testCtx := createTestContext()
	testTenant := buildTestTenantsData()

	tests := []struct {
		name      string
		ctx       context.Context
		id        uuid.UUID
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Success",
			ctx:  testCtx,
			id:   testID,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(testTenant, nil)
			},
			wantErr: false,
		},
		{
			name: "Error - Service Failed",
			ctx:  testCtx,
			id:   testID,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(nil, errors.New("service error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			result, err := resolver.FetchTenant(tt.ctx, tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func createTestContext() context.Context {
	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", uuid.New().String())
	return context.WithValue(context.Background(), config.GinContextKey, ginCtx)
}

func buildTestTenantsData() map[string]interface{} {
	id := uuid.New()
	return map[string]interface{}{
		"key":        id.String(),
		"created_at": time.Now().String(),
		"updated_at": time.Now().String(),
		"attributes": map[string]interface{}{
			"name":        "test-tenant",
			"description": "Test Tenant",
			"createdBy":   uuid.New().String(),
			"updatedBy":   uuid.New().String(),
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
		},
	}
}
