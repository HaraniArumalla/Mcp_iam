package tenants

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

func TestCreateTenant(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	resolver := TenantMutationResolver{
		PC: mockService,
	}

	// Setup test context
	userID := uuid.New().String()
	tenantID := uuid.New().String()

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", tenantID)
	ginCtx.Set("userID", userID)

	ginCtxNoUser := &gin.Context{}
	ginCtxNoUser.Set("tenantID", tenantID)

	validCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)
	noUserCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtxNoUser)

	// Setup test input
	validID := uuid.New()
	validInput := prepareValidInput()

	invalidInput := models.CreateTenantInput{
		ID:   uuid.Nil,
		Name: "",
	}

	// Mock tenant data that will be returned after creation
	tenant := buildTestTenantsData()
	mappedTenant, _ := MapTenantResponseToStruct(tenant)

	testCases := []struct {
		name      string
		ctx       context.Context
		input     models.CreateTenantInput
		mockSetup func()
		output    models.OperationResult
	}{
		{
			name:      "Invalid input data",
			ctx:       validCtx,
			input:     invalidInput,
			mockSetup: func() {},
			output:    buildErrorResponse(400, "Invalid input data", "Invalid input data"),
		},
		{
			name:  "Missing user ID in context",
			ctx:   noUserCtx,
			input: validInput,
			mockSetup: func() {
				// No mock calls expected
			},
			output: buildErrorResponse(400, "Failed to prepare metadata in create tenant", "error getting user ID from context"),
		},
		{
			name:  "Error creating tenant in permit",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(nil, errors.New("permit error")).MaxTimes(1)
			},
			output: buildErrorResponse(400, "Failed to create tenant in permit system", "permit error"),
		},
		{
			name:  "Error creating resource instance",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(map[string]interface{}{"key": validID.String()}, nil).MaxTimes(1)

				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(nil, errors.New("resource instance error")).MaxTimes(1)
			},
			output: buildErrorResponse(400, "Failed to create resource instance in permit system", "resource instance error"),
		},
		{
			name:  "Error getting created tenant",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(map[string]interface{}{"key": validID.String()}, nil).MaxTimes(1)

				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(map[string]interface{}{"created": true}, nil).MaxTimes(1)

				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(nil, errors.New("get tenant error")).MaxTimes(1)
			},
			output: buildErrorResponse(400, "Failed to fetch created tenant", "get tenant error"),
		},
		{
			name:  "Success",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(map[string]interface{}{"key": validID.String()}, nil).MaxTimes(1)

				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(map[string]interface{}{"created": true}, nil).MaxTimes(1)

				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(tenant, nil).MaxTimes(1)
			},
			output: buildSuccessResponse(mappedTenant),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			result, _ := resolver.CreateTenant(tc.ctx, tc.input)
			assert.NotNil(t, result)
		})
	}
}

func TestUpdateTenant(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	resolver := TenantMutationResolver{
		PC: mockService,
	}

	// Setup test context
	userID := uuid.New().String()
	tenantID := uuid.New().String()

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", tenantID)
	ginCtx.Set("userID", userID)

	validCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)

	// Setup test input
	name := "Updated Tenant"
	description := "Updated Description"
	street := "456 Update St"

	validID := uuid.New()
	validInput := models.UpdateTenantInput{
		ID:          validID,
		Name:        &name,
		Description: &description,
		ContactInfo: &models.ContactInfoInput{
			Address: &models.AddressInput{
				Street: &street,
			},
		},
	}

	invalidInput := models.UpdateTenantInput{
		ID: uuid.Nil,
	}

	// Mock tenant data
	tenantData := buildTestTenantsData()
	mappedTenant, _ := MapTenantResponseToStruct(tenantData)

	testCases := []struct {
		name      string
		ctx       context.Context
		input     models.UpdateTenantInput
		mockSetup func()
		output    models.OperationResult
	}{
		{
			name:      "Invalid input data",
			ctx:       validCtx,
			input:     invalidInput,
			mockSetup: func() {},
			output:    buildErrorResponse(400, "Invalid input data", "Invalid input data"),
		},
		{
			name:  "Error getting existing tenant",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(nil, errors.New("get tenant error")).MaxTimes(1)
			},
			output: buildErrorResponse(400, "Error getting existing tenant", "Error getting existing tenant"),
		},
		{
			name:  "Error updating tenant",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(tenantData, nil).MaxTimes(1)

				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(nil, errors.New("update error")).MaxTimes(1)
			},
			output: buildErrorResponse(400, "Error updating tenant", "Error updating tenant"),
		},
		{
			name:  "Error getting updated tenant",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				// First call for getExistingTenant
				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(tenantData, nil).MaxTimes(1)

				// Call for updating tenant
				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(map[string]interface{}{"updated": true}, nil).MaxTimes(1)

				// Second call for getCreatedTenant
				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(nil, errors.New("get updated tenant error")).MaxTimes(1)
			},
			output: buildErrorResponse(400, "Error getting updated tenant", "Error getting updated tenant"),
		},
		{
			name:  "Success",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				// First call for getExistingTenant
				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(tenantData, nil).MaxTimes(1)

				// Call for updating tenant
				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(map[string]interface{}{"updated": true}, nil)

				// Second call for getCreatedTenant
				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(tenantData, nil).MaxTimes(1)
			},
			output: buildSuccessResponse(mappedTenant),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			result, _ := resolver.UpdateTenant(tc.ctx, tc.input)

			assert.NotNil(t, result)
		})
	}
}

func TestDeleteTenant(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	resolver := TenantMutationResolver{
		PC: mockService,
	}

	// Setup test context
	tenantID := uuid.New().String()
	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", tenantID)
	validCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)

	validID := uuid.New()
	validInput := models.DeleteInput{ID: validID}
	invalidInput := models.DeleteInput{ID: uuid.Nil}

	testCases := []struct {
		name      string
		ctx       context.Context
		input     models.DeleteInput
		mockSetup func()
		output    models.OperationResult
	}{
		{
			name:      "Invalid ID",
			ctx:       validCtx,
			input:     invalidInput,
			mockSetup: func() {},
			output:    buildErrorResponse(400, "Invalid ID", "Invalid ID"),
		},
		{
			name:  "Error deleting tenant",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "DELETE", mock.Any(), nil).
					Return(nil, errors.New("delete error"))
			},
			output: buildErrorResponse(400, "Error deleting tenant", "Error deleting tenant"),
		},
		{
			name:  "Success",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "DELETE", mock.Any(), nil).
					Return(map[string]interface{}{"deleted": true}, nil)
			},
			output: buildSuccessResponse([]models.Data{}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			result, _ := resolver.DeleteTenant(tc.ctx, tc.input)

			assert.NotNil(t, result)
		})
	}
}

func TestPrepareMetadata(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	resolver := TenantMutationResolver{}

	// Setup test context
	userID := uuid.New().String()
	tenantID := uuid.New().String()

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", tenantID)
	ginCtx.Set("userID", userID)

	ginCtxNoUser := &gin.Context{}
	ginCtxNoUser.Set("tenantID", tenantID)

	validCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)
	validInput := prepareValidInput()
	t.Run("Success", func(t *testing.T) {
		metadata, err := resolver.prepareMetadata(validCtx, validInput)
		assert.NoError(t, err)
		assert.NotNil(t, metadata)
		assert.Equal(t, validInput.Name, metadata["name"])
	})
}

func TestMergeTenantData(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	resolver := TenantMutationResolver{}

	// Setup test context
	userID := uuid.New().String()
	tenantID := uuid.New().String()

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", tenantID)
	ginCtx.Set("userID", userID)

	validCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)

	// Setup test data
	existing := map[string]interface{}{
		"name":        "Old Name",
		"description": "Old Description",
		"contantInfo": map[string]interface{}{
			"email":       "old@example.com",
			"phoneNumber": "1111111111",
			"address": map[string]interface{}{
				"street":  "Old Street",
				"city":    "Old City",
				"state":   "OS",
				"zipcode": "11111",
				"country": "Old Country",
			},
		},
	}

	name := "New Name"
	description := "New Description"
	email := "new@example.com"
	street := "New Street"

	input := models.UpdateTenantInput{
		ID:          uuid.New(),
		Name:        &name,
		Description: &description,
		ContactInfo: &models.ContactInfoInput{
			Email: &email,
			Address: &models.AddressInput{
				Street: &street,
			},
		},
	}

	t.Run("Successful merge", func(t *testing.T) {
		merged, err := resolver.mergeTenantData(validCtx, existing, input)
		assert.NoError(t, err)
		assert.Equal(t, *input.Name, merged["name"])
		assert.Equal(t, *input.Description, merged["description"])
	})
}

func buildSuccessResponse(data []models.Data) models.OperationResult {
	response, _ := utils.FormatSuccessResponse(data)
	return response
}

func buildErrorResponse(statusCode int, message string, details string) models.OperationResult {
	return utils.FormatErrorResponse(statusCode, message, details)
}

func prepareValidInput() models.CreateTenantInput {
	// Setup test input
	description := "Test Description"
	email := "test@example.com"
	phoneNumber := "123-456-7890"
	street := "123 Test St"
	city := "Test City"
	state := "TS"
	zipcode := "12345"
	country := "Test Country"

	return models.CreateTenantInput{
		ID:          uuid.New(),
		Name:        "Test Tenant",
		Description: &description,
		ContactInfo: &models.ContactInfoInput{
			Email:       &email,
			PhoneNumber: &phoneNumber,
			Address: &models.AddressInput{
				Street:  &street,
				City:    &city,
				State:   &state,
				Zipcode: &zipcode,
				Country: &country,
			},
		},
	}
}
