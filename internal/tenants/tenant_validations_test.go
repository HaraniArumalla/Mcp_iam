package tenants

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

func createValidTenantInput() models.CreateTenantInput {
	return models.CreateTenantInput{
		ID:             uuid.New(),
		Name:           "Test Tenant",
		Description:    ptr("Test Description"),
		AccountOwnerID: uuid.New(), // Changed from AccountOwnerId to AccountOwnerID to match the struct field
		ContactInfo: &models.ContactInfoInput{
			Email:       ptr("test@example.com"),
			PhoneNumber: ptr("1234567890"),
			Address: &models.AddressInput{
				Street:  ptr("123 Test St"),
				City:    ptr("Test City"),
				State:   ptr("TS"),
				Zipcode: ptr("12345"),
				Country: ptr("US"),
			},
		},
		Tags: []*models.TagInput{
			{
				Key:   "test",
				Value: "test",
			},
		},
	}
}

func TestValidateCreateTenantInput(t *testing.T) {
	tests := []struct {
		name          string
		input         models.CreateTenantInput
		wantErr       bool
		expectedError string
	}{
		{
			name:    "Valid Input",
			input:   createValidTenantInput(),
			wantErr: false,
		},
		{
			name: "Missing Contact Info",
			input: func() models.CreateTenantInput {
				input := createValidTenantInput()
				input.ContactInfo = nil
				return input
			}(),
			wantErr:       true,
			expectedError: "contact info is required",
		},
		{
			name: "Missing Address",
			input: func() models.CreateTenantInput {
				input := createValidTenantInput()
				input.ContactInfo.Address = nil
				return input
			}(),
			wantErr:       true,
			expectedError: "address is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateCreateTenantInput(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
				assert.Equal(t, dto.CreateTenantInput{}, result) // Zero value for error case
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)

				// Verify mapping
				assert.Equal(t, tt.input.ID, result.ID)
				assert.Equal(t, tt.input.Name, result.Name)
				assert.Equal(t, tt.input.Description, result.Description)

				// Verify contact info mapping
				assert.Equal(t, *tt.input.ContactInfo.Email, result.ContactInfo.Email)
				assert.Equal(t, *tt.input.ContactInfo.PhoneNumber, result.ContactInfo.PhoneNumber)
				assert.Equal(t, *tt.input.ContactInfo.Address.Street, result.ContactInfo.Address.Street)
			}
		})
	}
}

func TestValidateRequiredFields(t *testing.T) {
	tests := []struct {
		name    string
		input   models.CreateTenantInput
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid Input",
			input:   createValidTenantInput(),
			wantErr: false,
		},
		{
			name: "Missing Contact Info",
			input: func() models.CreateTenantInput {
				input := createValidTenantInput()
				input.ContactInfo = nil
				return input
			}(),
			wantErr: true,
			errMsg:  "contact info is required",
		},
		{
			name: "Missing Address",
			input: func() models.CreateTenantInput {
				input := createValidTenantInput()
				input.ContactInfo.Address = nil
				return input
			}(),
			wantErr: true,
			errMsg:  "address is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRequiredFields(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMapToTenant(t *testing.T) {
	input := createValidTenantInput()
	result := MapToTenant(input)

	// Verify all fields are correctly mapped
	assert.Equal(t, input.ID, result.ID)
	assert.Equal(t, input.Name, result.Name)
	assert.Equal(t, input.Description, result.Description)

	// Verify contact info mapping
	assert.Equal(t, *input.ContactInfo.Email, result.ContactInfo.Email)
	assert.Equal(t, *input.ContactInfo.PhoneNumber, result.ContactInfo.PhoneNumber)

	// Verify address mapping
	assert.Equal(t, *input.ContactInfo.Address.Street, result.ContactInfo.Address.Street)
	assert.Equal(t, *input.ContactInfo.Address.City, result.ContactInfo.Address.City)
	assert.Equal(t, *input.ContactInfo.Address.State, result.ContactInfo.Address.State)
	assert.Equal(t, *input.ContactInfo.Address.Country, result.ContactInfo.Address.Country)
	assert.Equal(t, *input.ContactInfo.Address.Zipcode, result.ContactInfo.Address.Zipcode)
}

func TestValidateComponents(t *testing.T) {
	tests := []struct {
		name    string
		input   dto.CreateTenantInput
		wantErr bool
	}{
		{
			name:    "Valid Input",
			input:   MapToTenant(createValidTenantInput()),
			wantErr: false,
		},
		{
			name: "Invalid Address - Missing Street",
			input: func() dto.CreateTenantInput {
				input := MapToTenant(createValidTenantInput())
				input.ContactInfo.Address.Street = ""
				return input
			}(),
			wantErr: true,
		},
		{
			name: "Invalid Address - Missing City",
			input: func() dto.CreateTenantInput {
				input := MapToTenant(createValidTenantInput())
				input.ContactInfo.Address.City = ""
				return input
			}(),
			wantErr: true,
		},
		{
			name: "Invalid Contact Info - Missing Email",
			input: func() dto.CreateTenantInput {
				input := MapToTenant(createValidTenantInput())
				input.ContactInfo.Email = ""
				return input
			}(),
			wantErr: true,
		},
		{
			name: "Invalid Input - Missing Name",
			input: func() dto.CreateTenantInput {
				input := MapToTenant(createValidTenantInput())
				input.Name = ""
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
