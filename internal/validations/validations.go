package validations

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	validator "github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Define a struct to represent each validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func init() {
	validate = validator.New()
}

// Validate validates a struct using the validator package
func Validate(s interface{}) error {
	return validate.Struct(s)
}

func ValidateStruct(u interface{}) error {
	err := validate.Struct(u)
	if err != nil {
		// Check if it's a ValidationErrors type
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			var validationErrors []ValidationError

			for _, err := range validationErrs {
				validationErrors = append(validationErrors, ValidationError{
					Field: err.Namespace(),
					Message: fmt.Sprintf("Validation failed for field '%s'. Got value: '%v'. Error: failed on '%s' validation",
						err.Namespace(),
						err.Value(),
						err.Tag()),
				})
			}

			validationErrorsJSON, _ := json.Marshal(validationErrors)
			return fmt.Errorf("validation failed: %s", string(validationErrorsJSON))
		}

		// For InvalidValidationError or other error types, return the error as is
		return fmt.Errorf("validation error: %w", err)
	}
	return nil
}

func UpdateDeletedMap() map[string]interface{} {
	return map[string]interface{}{
		"row_status": 0,
	}
}

// ValidateName validates that the input string matches the regex "^[A-Za-z0-9\\-_]+$".
func ValidateName(name string) error {
	// Define the regex pattern
	pattern := `^[A-Za-z0-9\-_]+$`
	// Compile the regex
	re := regexp.MustCompile(pattern)
	// Check if the name matches the regex
	if !re.MatchString(name) {
		return errors.New("invalid name: must contain only alphanumeric characters, hyphens, or underscores")
	}
	return nil
}

func CreateActionMap(store map[string]interface{}, actions []string) map[string]interface{} {
	for _, action := range actions {
		store[action] = map[string]interface{}{
			"name": action,
		}
	}
	return store
}

func GetActionMap(data map[string]interface{}, key string) map[string]interface{} {
	actionMap := data["actions"].(map[string]interface{})
	for _, value := range actionMap {
		value := value.(map[string]interface{})
		for key1 := range value {
			if key1 != "name" {
				delete(value, key1)
			}
		}
	}
	return actionMap
}
