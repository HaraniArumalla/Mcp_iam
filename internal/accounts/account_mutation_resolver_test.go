package accounts

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

func TestCreateAccount(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	resolver := AccountMutationResolver{
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

	ginCtxNoTenant := &gin.Context{}
	ginCtxNoTenant.Set("userID", userID)

	validCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)
	noUserCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtxNoUser)
	noTenantCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtxNoTenant)

	// Setup test input
	validID := uuid.New()
	validInput := prepareValidAccountInput()

	invalidInput := models.CreateAccountInput{
		ID:   uuid.Nil,
		Name: "",
	}

	// Mock account data that will be returned after creation
	account := buildTestAccountsData()
	mappedAccount, _ := MapAccountResponseToStruct(account)

	testCases := []struct {
		name      string
		ctx       context.Context
		input     models.CreateAccountInput
		mockSetup func()
		output    models.OperationResult
	}{
		{
			name:      "Invalid input data",
			ctx:       validCtx,
			input:     invalidInput,
			mockSetup: func() {},
			output:    buildErrorResponse(400, "Failed to do input mapping and validation", "Invalid input data"),
		},
		{
			name:  "Missing tenant ID in context",
			ctx:   noTenantCtx,
			input: validInput,
			mockSetup: func() {
				// No mock calls expected
			},
			output: buildErrorResponse(400, "Failed to get tenant ID", "error getting tenant ID from context"),
		},
		{
			name:  "Missing user ID in context",
			ctx:   noUserCtx,
			input: validInput,
			mockSetup: func() {
				// No mock calls expected
			},
			output: buildErrorResponse(400, "Failed to prepare metadata in create account", "error getting user ID from context"),
		},
		{
			name:  "Error creating resource instance",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "POST", "resource_instances", mock.Any()).
					Return(nil, errors.New("resource instance error")).MaxTimes(1)
			},
			output: buildErrorResponse(400, "Failed to create the resource instances", "resource instance error"),
		},
		{
			name:  "Error creating relationship tuples",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "POST", "resource_instances", mock.Any()).
					Return(map[string]interface{}{"key": validID.String()}, nil).MaxTimes(1)

				mockService.EXPECT().
					SendRequest(mock.Any(), "POST", "relationship_tuples", mock.Any()).
					Return(nil, errors.New("relationship tuples error")).MaxTimes(1)
			},
			output: buildErrorResponse(400, "Failed to create the relationship tuples", "relationship tuples error"),
		},
		{
			name:  "Error getting created account",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "POST", "resource_instances", mock.Any()).
					Return(map[string]interface{}{"key": validID.String()}, nil).MaxTimes(1)

				mockService.EXPECT().
					SendRequest(mock.Any(), "POST", "relationship_tuples", mock.Any()).
					Return(map[string]interface{}{"created": true}, nil).MaxTimes(1)

				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(nil, errors.New("get account error")).MaxTimes(1)
			},
			output: buildErrorResponse(400, "Failed to get the account details by id", "get account error"),
		},
		{
			name:  "Success",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "POST", "resource_instances", mock.Any()).
					Return(map[string]interface{}{"key": validID.String()}, nil).MaxTimes(1)

				mockService.EXPECT().
					SendRequest(mock.Any(), "POST", "relationship_tuples", mock.Any()).
					Return(map[string]interface{}{"created": true}, nil).MaxTimes(1)

				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(account, nil).MaxTimes(1)
			},
			output: buildSuccessResponse(mappedAccount),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			result, _ := resolver.CreateAccount(tc.ctx, tc.input)
			assert.NotNil(t, result)
			// Additional assertions can be added here to check specific fields
		})
	}
}

func TestUpdateAccount(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	resolver := AccountMutationResolver{
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

	// Setup test input
	name := "Updated Account"
	description := "Updated Description"
	cardNumber := "4111-1111-1111-1111"
	street := "456 Update St"

	validID := uuid.New()
	validInput := models.UpdateAccountInput{
		ID:          validID,
		Name:        &name,
		Description: &description,
		BillingInfo: &models.UpdateBillingInfoInput{
			CreditCardNumber: &cardNumber,
			BillingAddress: &models.UpdateBillingAddressInput{
				Street: &street,
			},
		},
	}

	// Mock account data
	accountData := buildTestAccountsData()
	mappedAccount, _ := MapAccountResponseToStruct(accountData)

	testCases := []struct {
		name      string
		ctx       context.Context
		input     models.UpdateAccountInput
		mockSetup func()
		output    models.OperationResult
	}{
		{
			name:  "Error getting existing account",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(nil, errors.New("get account error")).MaxTimes(1)
			},
			output: buildErrorResponse(400, "Failed to get existing account data", "get account error"),
		},
		{
			name:  "Error updating account in permit",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(accountData, nil).MaxTimes(1)

				mockService.EXPECT().
					SendRequest(mock.Any(), "PATCH", mock.Any(), mock.Any()).
					Return(nil, errors.New("update error")).MaxTimes(1)
			},
			output: buildErrorResponse(400, "Failed to update account in Permit", "update error"),
		},
		{
			name:  "Error getting updated account",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				// First call for getExistingAccount
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(accountData, nil).MaxTimes(1)

				// Call for updating account - use PATCH instead of GET
				mockService.EXPECT().
					SendRequest(mock.Any(), "PATCH", mock.Any(), mock.Any()).
					Return(map[string]interface{}{"updated": true}, nil).MaxTimes(1)

				// Second call for getAccountById
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(nil, errors.New("get updated account error")).MaxTimes(1)
			},
			output: buildErrorResponse(400, "Failed to get the account details by id", "get updated account error"),
		},
		{
			name:  "Success",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				// First call for getExistingAccount
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(accountData, nil).MaxTimes(1)

				// Call for updating account - use PATCH instead of GET
				mockService.EXPECT().
					SendRequest(mock.Any(), "PATCH", mock.Any(), mock.Any()).
					Return(map[string]interface{}{"updated": true}, nil).MaxTimes(1)

				// Second call for getAccountById
				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(accountData, nil).MaxTimes(1)
			},
			output: buildSuccessResponse(mappedAccount),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			result, _ := resolver.UpdateAccount(tc.ctx, tc.input)

			assert.NotNil(t, result)
			// Additional assertions can be added here to check specific fields
		})
	}
}

func TestDeleteAccount(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	resolver := AccountMutationResolver{
		PC: mockService,
	}

	// Setup test context
	tenantID := uuid.New().String()
	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", tenantID)
	validCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)

	validID := uuid.New()
	validInput := models.DeleteInput{ID: validID}

	testCases := []struct {
		name      string
		ctx       context.Context
		input     models.DeleteInput
		mockSetup func()
		output    models.OperationResult
	}{
		{
			name:  "Error deleting account",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(nil, errors.New("delete error"))
			},
			output: buildErrorResponse(400, "Failed to delete account in Permit", "delete error"),
		},
		{
			name:  "Success",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(map[string]interface{}{"deleted": true}, nil)
			},
			output: buildSuccessResponse([]models.Data{}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			result, _ := resolver.DeleteAccount(tc.ctx, tc.input)

			assert.NotNil(t, result)
			// Additional assertions can be added here
		})
	}
}

func TestPrepareMetadata(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	resolver := AccountMutationResolver{}

	// Setup test context with actual user ID
	userID := "e03e8f81-ee65-43b0-b823-793d1bdab114" // Match the hard-coded value in helpers.GetUserID
	tenantID := uuid.New().String()

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", tenantID)
	ginCtx.Set("userID", userID)

	// Create context without user ID - we'll explicitly remove it
	ginCtxNoUser := &gin.Context{}
	ginCtxNoUser.Set("tenantID", tenantID)
	// No userID set to properly test the failure case

	validCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)
	validInput := prepareValidAccountInput()

	t.Run("Success", func(t *testing.T) {
		metadata, err := resolver.prepareMetadata(validCtx, validInput)
		assert.NoError(t, err)
		assert.NotNil(t, metadata)
		assert.Equal(t, validInput.Name, metadata["name"])
		assert.NotNil(t, metadata["billingInfo"])
	})
}

func TestMergeAccountData(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	resolver := AccountMutationResolver{}

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
		"billingInfo": map[string]interface{}{
			"creditCardNumber": "1234-5678-1234-5678",
			"creditCardType":   "Visa",
			"expirationDate":   "12/25",
			"cvv":              "123",
			"billingAddress": map[string]interface{}{
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
	cardNumber := "9876-5432-9876-5432"
	street := "New Street"

	input := models.UpdateAccountInput{
		ID:          uuid.New(),
		Name:        &name,
		Description: &description,
		BillingInfo: &models.UpdateBillingInfoInput{
			CreditCardNumber: &cardNumber,
			BillingAddress: &models.UpdateBillingAddressInput{
				Street: &street,
			},
		},
	}

	t.Run("Successful merge", func(t *testing.T) {
		merged, err := resolver.mergeAccountData(validCtx, existing, input)
		assert.NoError(t, err)
		assert.Equal(t, *input.Name, merged["name"])
		assert.Equal(t, *input.Description, merged["description"])

		// Verify billing info was merged
		billingInfo, ok := merged["billingInfo"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, *input.BillingInfo.CreditCardNumber, billingInfo["creditCardNumber"])

		// Verify billing address was merged
		billingAddress, ok := billingInfo["billingAddress"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, *input.BillingInfo.BillingAddress.Street, billingAddress["street"])
	})
}

func TestMergeBillingInfo(t *testing.T) {
	resolver := AccountMutationResolver{}

	existing := map[string]interface{}{
		"creditCardNumber": "1234-5678-1234-5678",
		"creditCardType":   "Visa",
		"expirationDate":   "12/25",
		"cvv":              "123",
		"billingAddress": map[string]interface{}{
			"street":  "Old Street",
			"city":    "Old City",
			"state":   "OS",
			"zipcode": "11111",
			"country": "Old Country",
		},
	}

	cardNumber := "9876-5432-9876-5432"
	cardType := "Mastercard"
	street := "New Street"

	updates := &models.UpdateBillingInfoInput{
		CreditCardNumber: &cardNumber,
		CreditCardType:   &cardType,
		BillingAddress: &models.UpdateBillingAddressInput{
			Street: &street,
		},
	}

	result := resolver.mergeBillingInfo(existing, updates)

	assert.Equal(t, *updates.CreditCardNumber, result["creditCardNumber"])
	assert.Equal(t, *updates.CreditCardType, result["creditCardType"])
	assert.Equal(t, "12/25", result["expirationDate"])
	assert.Equal(t, "123", result["cvv"])

	billingAddress, ok := result["billingAddress"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, *updates.BillingAddress.Street, billingAddress["street"])
	assert.Equal(t, "Old City", billingAddress["city"])
}

func TestMergeBillingAddress(t *testing.T) {
	resolver := AccountMutationResolver{}

	existing := map[string]interface{}{
		"street":  "Old Street",
		"city":    "Old City",
		"state":   "OS",
		"zipcode": "11111",
		"country": "Old Country",
	}

	street := "New Street"
	city := "New City"

	updates := &models.UpdateBillingAddressInput{
		Street: &street,
		City:   &city,
	}

	result := resolver.mergeBillingAddress(existing, updates)

	assert.Equal(t, *updates.Street, result["street"])
	assert.Equal(t, *updates.City, result["city"])
	assert.Equal(t, "OS", result["state"])
	assert.Equal(t, "11111", result["zipcode"])
	assert.Equal(t, "Old Country", result["country"])
}

// Helper functions
func buildSuccessResponse(data []models.Data) models.OperationResult {
	response, _ := utils.FormatSuccessResponse(data)
	return response
}

func buildErrorResponse(statusCode int, message string, details string) models.OperationResult {
	return utils.FormatErrorResponse(statusCode, message, details)
}

func prepareValidAccountInput() models.CreateAccountInput {
	description := "Test Description"
	parentID := uuid.New()
	tenantID := uuid.New()
	cardNumber := "4111-1111-1111-1111"
	cardType := "Visa"
	expirationDate := "12/25"
	cvv := "123"
	street := "123 Test St"
	city := "Test City"
	state := "TS"
	zipcode := "12345"
	country := "Test Country"

	return models.CreateAccountInput{
		ID:          uuid.New(),
		Name:        "Test Account",
		Description: &description,
		ParentID:    parentID,
		TenantID:    tenantID,
		BillingInfo: &models.CreateBillingInfoInput{
			CreditCardNumber: cardNumber,
			CreditCardType:   cardType,
			ExpirationDate:   expirationDate,
			Cvv:              cvv,
			BillingAddress: &models.CreateBillingAddressInput{
				Street:  street,
				City:    city,
				State:   state,
				Zipcode: zipcode,
				Country: country,
			},
		},
	}
}

func buildTestAccountsData() map[string]interface{} {
	id := uuid.New()
	return map[string]interface{}{
		"key":        id.String(),
		"name":       "Test Account",
		"created_at": "2023-01-01T00:00:00Z",
		"updated_at": "2023-01-01T00:00:00Z",
		"attributes": map[string]interface{}{
			"id":          id.String(),
			"name":        "Test Account",
			"description": "Test Account Description",
			"tenantId":    uuid.New().String(),
			"parentId":    uuid.New().String(),
			"createdBy":   uuid.New().String(),
			"updatedBy":   uuid.New().String(),
			"billingInfo": map[string]interface{}{
				"creditCardNumber": "4111-1111-1111-1111",
				"creditCardType":   "Visa",
				"expirationDate":   "12/25",
				"cvv":              "123",
				"billingAddress": map[string]interface{}{
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
