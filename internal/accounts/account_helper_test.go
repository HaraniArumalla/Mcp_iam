package accounts

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMapAccountsResponseToStruct(t *testing.T) {
	// Create valid test data
	validID := uuid.New()
	validData := map[string]interface{}{
		"data": []interface{}{
			buildTestAccountData(validID),
		},
	}

	tests := []struct {
		name    string
		input   map[string]interface{}
		wantErr bool
	}{
		{
			name:    "Valid input",
			input:   validData,
			wantErr: false,
		},
		{
			name:    "Nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name: "Missing data field",
			input: map[string]interface{}{
				"wrong": "data",
			},
			wantErr: true,
		},
		{
			name: "Invalid data format",
			input: map[string]interface{}{
				"data": "not an array",
			},
			wantErr: true,
		},
		{
			name: "Empty data array",
			input: map[string]interface{}{
				"data": []interface{}{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MapAccountsResponseToStruct(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.input != nil && len(tt.input["data"].([]interface{})) > 0 {
					assert.NotNil(t, result)
					assert.Len(t, result, len(tt.input["data"].([]interface{})))
				}
			}
		})
	}
}

func TestMapAccountResponseToStruct(t *testing.T) {
	// Create valid test data
	validID := uuid.New()
	validData := buildTestAccountData(validID)

	tests := []struct {
		name    string
		input   map[string]interface{}
		wantErr bool
	}{
		{
			name:    "Valid input",
			input:   validData,
			wantErr: false,
		},
		{
			name:    "Nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name: "Invalid key field",
			input: map[string]interface{}{
				"key": "not-a-uuid",
			},
			wantErr: true,
		},
		{
			name: "Missing attributes",
			input: map[string]interface{}{
				"key": validID.String(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MapAccountResponseToStruct(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, 1)
			}
		})
	}
}

func TestMapAccountData(t *testing.T) {
	validID := uuid.New()
	validData := buildTestAccountData(validID)

	tests := []struct {
		name    string
		input   map[string]interface{}
		wantErr bool
	}{
		{
			name:    "Valid input",
			input:   validData,
			wantErr: false,
		},
		{
			name:    "Nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name: "Invalid key field",
			input: map[string]interface{}{
				"key": "not-a-uuid",
			},
			wantErr: true,
		},
		{
			name: "Missing attributes",
			input: map[string]interface{}{
				"key": validID.String(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mapAccountData(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, validID, result.ID)
				assert.NotNil(t, result.BillingInfo)
				assert.NotNil(t, result.BillingInfo.BillingAddress)
			}
		})
	}
}

func TestMapBillingInfo(t *testing.T) {
	validAttributes := map[string]interface{}{
		"billingInfo": map[string]interface{}{
			"creditCardNumber": "1234-5678-9012-3456",
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
	}

	tests := []struct {
		name    string
		input   map[string]interface{}
		wantErr bool
	}{
		{
			name:    "Valid input",
			input:   validAttributes,
			wantErr: false,
		},
		{
			name:    "Missing billing info",
			input:   map[string]interface{}{},
			wantErr: true,
		},
		{
			name: "Invalid billing info type",
			input: map[string]interface{}{
				"billingInfo": "not a map",
			},
			wantErr: true,
		},
		{
			name: "Missing billing address",
			input: map[string]interface{}{
				"billingInfo": map[string]interface{}{
					"creditCardNumber": "1234-5678-9012-3456",
					"creditCardType":   "Visa",
					"expirationDate":   "12/25",
					"cvv":              "123",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mapBillingInfo(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "1234-5678-9012-3456", result.CreditCardNumber)
				assert.Equal(t, "Visa", result.CreditCardType)
				assert.Equal(t, "12/25", result.ExpirationDate)
				assert.Equal(t, "123", result.Cvv)
				assert.NotNil(t, result.BillingAddress)
			}
		})
	}
}

func TestMapBillingAddress(t *testing.T) {
	validBillingInfo := map[string]interface{}{
		"billingAddress": map[string]interface{}{
			"street":  "123 Test St",
			"city":    "Test City",
			"state":   "TS",
			"zipcode": "12345",
			"country": "US",
		},
	}

	tests := []struct {
		name    string
		input   map[string]interface{}
		wantErr bool
	}{
		{
			name:    "Valid input",
			input:   validBillingInfo,
			wantErr: false,
		},
		{
			name:    "Missing billing address",
			input:   map[string]interface{}{},
			wantErr: true,
		},
		{
			name: "Invalid billing address type",
			input: map[string]interface{}{
				"billingAddress": "not a map",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mapBillingAddress(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "123 Test St", result.Street)
				assert.Equal(t, "Test City", result.City)
				assert.Equal(t, "TS", result.State)
				assert.Equal(t, "12345", result.Zipcode)
				assert.Equal(t, "US", result.Country)
			}
		})
	}
}

// Helper function to build test account data
func buildTestAccountData(id uuid.UUID) map[string]interface{} {
	return map[string]interface{}{
		"key":        id.String(),
		"name":       "Test Account",
		"created_at": time.Now().String(),
		"updated_at": time.Now().String(),
		"attributes": map[string]interface{}{
			"id":          id.String(),
			"name":        "Test Account",
			"description": "Test Account Description",
			"tenantId":    uuid.New().String(),
			"parentId":    uuid.New().String(),
			"createdBy":   uuid.New().String(),
			"updatedBy":   uuid.New().String(),
			"billingInfo": map[string]interface{}{
				"creditCardNumber": "1234-5678-9012-3456",
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
