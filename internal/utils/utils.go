package utils

import (
	"fmt"
	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/pkg/logger"
	"net/http"
)

// formatSuccess formats a successful response in the `OperationResult` union
func FormatSuccess(data interface{}) (models.OperationResult, error) {
	// Type assertion: Convert 'data' from 'interface{}' to '[]models.Data'
	if typedData, ok := data.([]models.Data); ok {
		// If the assertion succeeds, return the result in OperationResult
		successResponse := &models.SuccessResponse{
			IsSuccess: true,
			Message:   "Operation successful",
			Data:      typedData, // Now, typedData is of type []models.Data
		}
		var opResult models.OperationResult = successResponse
		return opResult, nil
	}
	// If the type assertion fails, return an error
	return nil, fmt.Errorf("expected data to be of type []models.Data, but got %T", data)
}

// FormatSuccessResponse formats a successful response in the `OperationResult` union
func FormatSuccessResponse(data []models.Data) (models.OperationResult, error) {
	successResponse, err := FormatSuccess(data)
	if err != nil {
		logger.LogError("Failed to format success response", "error", err)
		return FormatErrorResponse(http.StatusBadRequest, "Failed to format success response", err.Error()), nil
	}
	return successResponse, nil
}

// FormatErrorResponse formats an error response in the `OperationResult` union
func FormatErrorResponse(errorCode int, message, errDetails string) models.OperationResult {
	var details *string
	if errDetails != "" {
		details = &errDetails
	}
	logger.LogError(message, "status", fmt.Sprint(errorCode), "error", errDetails)
	return &models.ResponseError{
		IsSuccess:     false,
		Message:       config.GenericErrorMessage,
		SystemMessage: message,
		ErrorCode:     fmt.Sprint(errorCode),
		ErrorDetails:  details,
	}
}
