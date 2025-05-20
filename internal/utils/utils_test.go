package utils

import (
	"fmt"
	"testing"

	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Test data helper function
func createTestData() []models.Data {
	return []models.Data{
		&models.Account{
			ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Name: "test",
		},
	}
}

func TestFormatSuccess(t *testing.T) {
	testcases := []struct {
		name        string
		input       interface{}
		wantSuccess bool
		wantErr     bool
		wantMsg     string
	}{
		{
			name:        "Valid data with items",
			input:       createTestData(),
			wantSuccess: true,
			wantErr:     false,
			wantMsg:     "Operation successful",
		},
		{
			name:        "Empty data slice",
			input:       []models.Data{},
			wantSuccess: true,
			wantErr:     false,
			wantMsg:     "Operation successful",
		},
		{
			name:        "Nil slice",
			input:       []models.Data(nil),
			wantSuccess: true,
			wantErr:     false,
			wantMsg:     "Operation successful",
		},
		{
			name:        "Invalid type",
			input:       "invalid",
			wantSuccess: false,
			wantErr:     true,
			wantMsg:     "expected data to be of type []models.Data",
		},
		{
			name:        "Nil input",
			input:       nil,
			wantSuccess: false,
			wantErr:     true,
			wantMsg:     "expected data to be of type []models.Data",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := FormatSuccess(tc.input)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantMsg)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)

			successResp, ok := result.(*models.SuccessResponse)
			assert.True(t, ok)
			assert.True(t, successResp.IsSuccess)
			assert.Equal(t, tc.wantMsg, successResp.Message)

			if tc.input != nil {
				assert.Equal(t, tc.input, successResp.Data)
			}
		})
	}
}

func TestFormatErrorResponse(t *testing.T) {
	testcases := []struct {
		name        string
		code        int
		msg         string
		details     string
		wantMsg     string
		wantDetails *string
	}{
		{
			name:        "Basic error",
			code:        400,
			msg:         "Test error",
			details:     "Error details",
			wantMsg:     config.GenericErrorMessage,
			wantDetails: ptr("Error details"),
		},
		{
			name:        "Empty details",
			code:        500,
			msg:         "Server error",
			details:     "",
			wantMsg:     config.GenericErrorMessage,
			wantDetails: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatErrorResponse(tc.code, tc.msg, tc.details)
			assert.NotNil(t, result)

			errorResp, ok := result.(*models.ResponseError)
			assert.True(t, ok)

			assert.Equal(t, tc.wantMsg, errorResp.Message)
			assert.Equal(t, tc.msg, errorResp.SystemMessage)
			assert.Equal(t, fmt.Sprint(tc.code), errorResp.ErrorCode)
			assert.Equal(t, tc.wantDetails, errorResp.ErrorDetails)
			assert.False(t, errorResp.IsSuccess)
		})
	}
}

func TestFormatSuccessResponse(t *testing.T) {
	testcases := []struct {
		name     string
		input    []models.Data
		wantErr  bool
		wantCode int
	}{
		{
			name:     "Valid data",
			input:    createTestData(),
			wantErr:  false,
			wantCode: 200,
		},
		{
			name:     "Empty slice",
			input:    []models.Data{},
			wantErr:  false,
			wantCode: 200,
		},
		{
			name:     "Nil input",
			input:    nil,
			wantErr:  false,
			wantCode: 200,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := FormatSuccessResponse(tc.input)

			assert.Nil(t, err)
			assert.NotNil(t, result)

			if !tc.wantErr {
				successResp, ok := result.(*models.SuccessResponse)
				assert.True(t, ok)
				assert.True(t, successResp.IsSuccess)
				assert.Equal(t, tc.input, successResp.Data)
			}
		})
	}
}

// Helper function to create string pointer
func ptr(s string) *string {
	return &s
}
