package accounts

import (
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/dto"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func ptr(s string) *string {
	return &s
}

func createValidAccountInput() models.CreateAccountInput {
	cardNumber := "4111111111111111"
	cardType := "visa"
	expirationDate := "12/25"
	cvv := "123"
	street := "123 Test St"
	city := "Test City"
	state := "TS"
	zipcode := "12345"
	country := "US"
	return models.CreateAccountInput{
		ID:             uuid.New(),
		AccountOwnerID: uuid.New(),
		Name:           "Test Account",
		Description:    ptr("Test Description"),
		TenantID:       uuid.New(),
		ParentID:       uuid.New(),
		RelationType:   "PARENT",
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

func TestValidateCreateAccountInput(t *testing.T) {
	tests := []struct {
		name          string
		input         models.CreateAccountInput
		wantErr       bool
		expectedError string
	}{
		{
			name:    "Valid Input",
			input:   createValidAccountInput(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateCreateAccountInput(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
				assert.Equal(t, dto.CreateAccountInput{}, result) // Zero value for error case
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)

				// Verify mapping
				assert.Equal(t, tt.input.ID, result.ID)
				assert.Equal(t, tt.input.Name, result.Name)
				assert.Equal(t, tt.input.Description, result.Description)
				assert.Equal(t, tt.input.TenantID, result.TenantID)
				assert.Equal(t, tt.input.ParentID, result.ParentID)
				assert.Equal(t, tt.input.AccountOwnerID, result.AccountOwnerID)
				assert.Equal(t, tt.input.RelationType, models.RelationTypeEnum(string(result.RelationType)))

			}
		})
	}
}

func TestMapToAccount(t *testing.T) {
	input := createValidAccountInput()
	result := MapToAccount(input)

	// Verify all fields are correctly mapped
	assert.Equal(t, input.ID, result.ID)
	assert.Equal(t, input.Name, result.Name)
	assert.Equal(t, input.Description, result.Description)
	assert.Equal(t, input.TenantID, result.TenantID)
	assert.Equal(t, input.ParentID, result.ParentID)

	// Verify billing info mapping
	assert.Equal(t, input.BillingInfo.CreditCardNumber, result.BillingInfo.CreditCardNumber)
	assert.Equal(t, input.BillingInfo.CreditCardType, result.BillingInfo.CreditCardType)
	assert.Equal(t, input.BillingInfo.ExpirationDate, result.BillingInfo.ExpirationDate)
	assert.Equal(t, input.BillingInfo.Cvv, result.BillingInfo.Cvv)

	// Verify billing address mapping
	assert.Equal(t, input.BillingInfo.BillingAddress.Street, result.BillingInfo.BillingAddress.Street)
	assert.Equal(t, input.BillingInfo.BillingAddress.City, result.BillingInfo.BillingAddress.City)
	assert.Equal(t, input.BillingInfo.BillingAddress.State, result.BillingInfo.BillingAddress.State)
	assert.Equal(t, input.BillingInfo.BillingAddress.Country, result.BillingInfo.BillingAddress.Country)
	assert.Equal(t, input.BillingInfo.BillingAddress.Zipcode, result.BillingInfo.BillingAddress.Zipcode)
}

func TestValidateComponents(t *testing.T) {
	tests := []struct {
		name    string
		input   dto.CreateAccountInput
		wantErr bool
	}{
		{
			name:    "Valid Input",
			input:   MapToAccount(createValidAccountInput()),
			wantErr: false,
		},
		{
			name: "Invalid Billing Address - Missing Street",
			input: func() dto.CreateAccountInput {
				input := MapToAccount(createValidAccountInput())
				input.BillingInfo.BillingAddress.Street = ""
				return input
			}(),
			wantErr: true,
		},
		{
			name: "Invalid Billing Address - Missing City",
			input: func() dto.CreateAccountInput {
				input := MapToAccount(createValidAccountInput())
				input.BillingInfo.BillingAddress.City = ""
				return input
			}(),
			wantErr: true,
		},
		{
			name: "Invalid Billing Info - Missing Credit Card Number",
			input: func() dto.CreateAccountInput {
				input := MapToAccount(createValidAccountInput())
				input.BillingInfo.CreditCardNumber = ""
				return input
			}(),
			wantErr: true,
		},
		{
			name: "Invalid Billing Info - Missing Credit Card Type",
			input: func() dto.CreateAccountInput {
				input := MapToAccount(createValidAccountInput())
				input.BillingInfo.CreditCardType = ""
				return input
			}(),
			wantErr: true,
		},
		{
			name: "Invalid Billing Info - Missing Expiration Date",
			input: func() dto.CreateAccountInput {
				input := MapToAccount(createValidAccountInput())
				input.BillingInfo.ExpirationDate = ""
				return input
			}(),
			wantErr: true,
		},
		{
			name: "Invalid Billing Info - Missing CVV",
			input: func() dto.CreateAccountInput {
				input := MapToAccount(createValidAccountInput())
				input.BillingInfo.Cvv = ""
				return input
			}(),
			wantErr: true,
		},
		{
			name: "Invalid Input - Missing Name",
			input: func() dto.CreateAccountInput {
				input := MapToAccount(createValidAccountInput())
				input.Name = ""
				return input
			}(),
			wantErr: true,
		},
		{
			name: "Invalid Input - Invalid TenantID",
			input: func() dto.CreateAccountInput {
				input := MapToAccount(createValidAccountInput())
				input.TenantID = uuid.Nil
				return input
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateComponents(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
