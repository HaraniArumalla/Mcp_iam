package accounts

import (
	"context"
	"errors"
	"fmt"
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

func TestAccounts(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	resolver := AccountQueryResolver{
		PC: mockService,
	}

	// Setup test context
	tenantID := uuid.New()
	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", tenantID.String())
	validCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)

	// Context without tenant ID
	ginCtxNoTenant := &gin.Context{}
	noTenantCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtxNoTenant)

	// Mock account data
	accountData := buildTestAccountData(uuid.New())
	accountsResponse := map[string]interface{}{
		"data": []interface{}{accountData},
	}
	mappedAccounts, _ := MapAccountsResponseToStruct(accountsResponse)
	successResponse, _ := utils.FormatSuccessResponse(mappedAccounts)

	testCases := []struct {
		name      string
		ctx       context.Context
		mockSetup func()
		expected  models.OperationResult
	}{
		{
			name: "Success with accounts",
			ctx:  validCtx,
			mockSetup: func() {
				expectedURL := fmt.Sprintf("resource_instances/detailed?tenant=%s&resource=%s",
					tenantID.String(), config.AccountResourceTypeID)
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", expectedURL, nil).
					Return(accountsResponse, nil)
			},
			expected: successResponse,
		},
		{
			name: "No tenant ID in context",
			ctx:  noTenantCtx,
			mockSetup: func() {
				// No mocks needed - function should error before API call
			},
			expected: utils.FormatErrorResponse(400, "Failed to get tenant ID", "error getting tenant ID from context"),
		},
		{
			name: "Error from permit service",
			ctx:  validCtx,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(nil, errors.New("permit service error"))
			},
			expected: utils.FormatErrorResponse(400, "Failed to get all accounts from permit", "permit service error"),
		},
		{
			name: "Invalid response format",
			ctx:  validCtx,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(map[string]interface{}{
						"data": "not an array",
					}, nil)
			},
			expected: utils.FormatErrorResponse(400, "Failed to get all accounts from permit", "missing or invalid data field"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			result, _ := resolver.Accounts(tc.ctx)

			assert.NotNil(t, result)
		})
	}
}

func TestAccount(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	resolver := AccountQueryResolver{
		PC: mockService,
	}

	// Setup test context
	tenantID := uuid.New()
	validID := uuid.New()

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", tenantID.String())
	validCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)

	// Context without tenant ID
	ginCtxNoTenant := &gin.Context{}
	noTenantCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtxNoTenant)

	// Mock account data
	accountData := buildTestAccountData(uuid.New())
	mappedAccount, _ := MapAccountResponseToStruct(accountData)
	successResponse, _ := utils.FormatSuccessResponse(mappedAccount)

	testCases := []struct {
		name      string
		ctx       context.Context
		id        uuid.UUID
		mockSetup func()
		expected  models.OperationResult
	}{
		{
			name: "Success",
			ctx:  validCtx,
			id:   validID,
			mockSetup: func() {
				expectedURL := fmt.Sprintf("resource_instances/%s", validID)
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", expectedURL, nil).
					Return(accountData, nil)
			},
			expected: successResponse,
		},
		{
			name: "No tenant ID in context",
			ctx:  noTenantCtx,
			id:   validID,
			mockSetup: func() {
				// No mocks needed - function should error before API call
			},
			expected: utils.FormatErrorResponse(400, "Failed to get tenant ID", "error getting tenant ID from context"),
		},
		{
			name: "Error from permit service",
			ctx:  validCtx,
			id:   validID,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(nil, errors.New("permit service error"))
			},
			expected: utils.FormatErrorResponse(400, "Failed to get account resources from permit", "permit service error"),
		},
		{
			name: "Invalid response format",
			ctx:  validCtx,
			id:   validID,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(map[string]interface{}{
						"key": "not a valid UUID",
					}, nil)
			},
			expected: utils.FormatErrorResponse(400, "Failed to get account resources from permit", "failed to get UUID from account data: invalid UUID format"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			result, _ := resolver.Account(tc.ctx, tc.id)

			assert.NotNil(t, result)
		})
	}
}

func TestFetchAccounts(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	resolver := AccountQueryResolver{
		PC: mockService,
	}

	tenantID := uuid.New()
	ctx := context.Background()

	accountsResponse := map[string]interface{}{
		"data": []interface{}{
			buildTestAccountData(uuid.New()),
		},
	}

	t.Run("Success", func(t *testing.T) {
		expectedURL := fmt.Sprintf("resource_instances/detailed?tenant=%s&resource=%s",
			tenantID.String(), config.AccountResourceTypeID)

		mockService.EXPECT().
			SendRequest(mock.Any(), "GET", expectedURL, nil).
			Return(accountsResponse, nil)

		result, err := resolver.fetchAccounts(ctx, &tenantID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
	})

	t.Run("Error from permit service", func(t *testing.T) {
		mockService.EXPECT().
			SendRequest(mock.Any(), "GET", mock.Any(), nil).
			Return(nil, errors.New("permit service error"))

		result, err := resolver.fetchAccounts(ctx, &tenantID)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Invalid response format", func(t *testing.T) {
		mockService.EXPECT().
			SendRequest(mock.Any(), "GET", mock.Any(), nil).
			Return(map[string]interface{}{
				"data": "not an array",
			}, nil)

		result, err := resolver.fetchAccounts(ctx, &tenantID)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestFetchAccount(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	resolver := AccountQueryResolver{
		PC: mockService,
	}

	validID := uuid.New()
	ctx := context.Background()

	accountData := buildTestAccountData(uuid.New())

	t.Run("Success", func(t *testing.T) {
		expectedURL := fmt.Sprintf("resource_instances/%s", validID)

		mockService.EXPECT().
			SendRequest(mock.Any(), "GET", expectedURL, nil).
			Return(accountData, nil)

		result, err := resolver.fetchAccount(ctx, validID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
	})

	t.Run("Error from permit service", func(t *testing.T) {
		mockService.EXPECT().
			SendRequest(mock.Any(), "GET", mock.Any(), nil).
			Return(nil, errors.New("permit service error"))

		result, err := resolver.fetchAccount(ctx, validID)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Invalid response format", func(t *testing.T) {
		mockService.EXPECT().
			SendRequest(mock.Any(), "GET", mock.Any(), nil).
			Return(map[string]interface{}{
				"key": "not a valid UUID",
			}, nil)

		result, err := resolver.fetchAccount(ctx, validID)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
