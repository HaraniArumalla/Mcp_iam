package helpers

import (
	"context"
	"fmt"
	"iam_services_main_v1/config"
	"iam_services_main_v1/pkg/logger"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetGinContext extracts the Gin context from a regular context
func GetGinContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(config.GinContextKey)
	if ginContext == nil {
		return nil, fmt.Errorf("could not retrieve gin context")
	}
	return ginContext.(*gin.Context), nil
}

// CheckValueExists returns fallback value if field is empty, otherwise returns field value
func CheckValueExists(field string, fallback string) string {
	if field == "" {
		return fallback
	}
	return field
}

// GetTenantID extracts and validates the tenant ID from the gin context, returning it as a UUID pointer
func GetTenantID(ctx context.Context) (*uuid.UUID, error) {
	ginCtx, ok := ctx.Value(config.GinContextKey).(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("gin context not found in the GetTenantID request")
	}
	tenantID, exists := ginCtx.Get("tenantID")
	if !exists {
		return nil, fmt.Errorf("tenant id not found in context")
	}

	switch tenantID := tenantID.(type) {
	case string:
		parsedTenantID, err := uuid.Parse(tenantID)
		if err != nil {
			return nil, fmt.Errorf("error parsing tenant id: %w", err)
		}
		return &parsedTenantID, nil
	case uuid.UUID:
		return &tenantID, nil
	default:
		return nil, fmt.Errorf("invalid tenant id type")
	}
}

// GetUserID extracts and validates the user ID from the Gin context, supporting both string and UUID formats.
func GetUserID(ctx context.Context) (*uuid.UUID, error) {
	ginCtx, ok := ctx.Value(config.GinContextKey).(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("gin context not found in the GetUserID request")
	}
	userID, exists := ginCtx.Get("userID")
	if !exists {
		return nil, fmt.Errorf("user id not found in context")
	}

	switch userID := userID.(type) {
	case string:
		parsedUserID, err := uuid.Parse(userID)
		if err != nil {
			return nil, fmt.Errorf("error parsing user id: %w", err)
		}
		return &parsedUserID, nil
	case uuid.UUID:
		return &userID, nil
	default:
		return nil, fmt.Errorf("invalid user id type")
	}
}

// StructToMap converts a struct to a map[string]interface{}, handling nested structs and pointer fields.
// Nil pointer fields are skipped. The resulting map uses struct field names as keys.
func StructToMap(input interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(input)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return result
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		//Skip if field is nil
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			continue
		}

		//convert nested struct
		if fieldValue.Kind() == reflect.Struct {
			result[field.Name] = StructToMap(fieldValue.Interface())
		} else {
			result[field.Name] = fieldValue.Interface()
		}

	}

	return result
}

// Automate struct mapping using reflection
func MapStruct(src interface{}, dst interface{}) error {
	// Get the value of the source struct and destination struct
	srcValue := reflect.ValueOf(src)
	dstValue := reflect.ValueOf(dst)

	// Ensure both are pointers (to modify the destination)
	if srcValue.Kind() != reflect.Ptr || dstValue.Kind() != reflect.Ptr {
		return fmt.Errorf("both source and destination must be pointers")
	}

	// Dereference the pointers to work with the actual values
	srcValue = srcValue.Elem()
	dstValue = dstValue.Elem()

	// Ensure both are structs
	if srcValue.Kind() != reflect.Struct || dstValue.Kind() != reflect.Struct {
		return fmt.Errorf("both source and destination must be structs")
	}

	// Iterate over the fields of the source struct
	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		dstField := dstValue.FieldByName(srcValue.Type().Field(i).Name)

		// If the destination field is valid and can be set, copy the value
		if dstField.IsValid() && dstField.CanSet() {
			// Only copy the value if the field types match
			if srcField.Type() == dstField.Type() {
				dstField.Set(srcField)
			}
		}
	}
	return nil
}

// MergeMaps merges two maps by recursively combining nested maps and overwriting non-map values from updates into existing map
func MergeMaps(existing, updates map[string]interface{}) map[string]interface{} {
	for key, newValue := range updates {
		if newValue == nil {
			// skip nil values in updates
			continue
		}
		// check if the value is a map and recursively merge
		existingValue, exists := existing[key]

		if exists {
			if existingMap, ok := existingValue.(map[string]interface{}); ok {
				if newMap, ok := newValue.(map[string]interface{}); ok {
					existing[key] = MergeMaps(existingMap, newMap)
					continue
				}
			}
		}
		//overwrite or add the value from updates
		existing[key] = newValue
	}
	return existing
}

// GetString retrieves a string value from a map by key, returning empty string if not found or not a string type
func GetString(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

// GetUUID extracts and parses a UUID string from a map for the given key, returning the UUID or an error if invalid/missing
func GetUUID(data map[string]interface{}, key string) (uuid.UUID, error) {
	if val, ok := data[key].(string); ok {
		return uuid.Parse(val)
	}
	return uuid.UUID{}, fmt.Errorf("invalid or missing UUID for key: %s", key)
}

// GetMap extracts and type-asserts a nested map[string]interface{} from the given data using the specified key
func GetMap(data map[string]interface{}, key string) (map[string]interface{}, error) {
	if val, ok := data[key].(map[string]interface{}); ok {
		return val, nil
	}
	return nil, fmt.Errorf("missing or invalid map for key: %s", key)
}

// StringPtr returns a pointer to the given string.
func StringPtr(s string) *string {
	return &s
}

// GetSlice extracts and type-asserts a slice of interfaces from the given data using the specified key.
// It supports extracting both []interface{} and typed slices, returning an error if the key doesn't exist
// or if the value is not a slice type.
func GetSlice(data map[string]interface{}, key string) ([]interface{}, error) {
	val, exists := data[key]
	if !exists {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	switch v := val.(type) {
	case []interface{}:
		return v, nil
	case nil:
		return nil, fmt.Errorf("nil value for key: %s", key)
	default:
		// Handle case where value might be a slice but not []interface{}
		valRef := reflect.ValueOf(val)
		if valRef.Kind() == reflect.Slice {
			result := make([]interface{}, valRef.Len())
			for i := 0; i < valRef.Len(); i++ {
				result[i] = valRef.Index(i).Interface()
			}
			return result, nil
		}
		return nil, fmt.Errorf("value for key %s is not a slice type", key)
	}
}

// GetUserAndTenantID extracts user ID and tenant ID from context, returning both UUIDs or an error if extraction fails
func GetUserAndTenantID(ctx context.Context) (*uuid.UUID, *uuid.UUID, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		logger.LogError("error occurred when fetching the user id")
		return nil, nil, err
	}

	tenantID, err := GetTenantID(ctx)
	if err != nil {
		logger.LogError("error occurred when fetching the tenant id", "error", err)
		return nil, nil, err
	}

	return userID, tenantID, nil
}

func Ptr(s string) *string {
	return &s
}

// ConvertToZFormat converts date string with timezone offset to UTC Z format
func ConvertToZFormat(input string) (string, error) {
	t, err := time.Parse(time.RFC3339, input)
	if err != nil {
		logger.LogError("error parsing time", "input", input, "error", err)
		return "", err
	}
	return t.UTC().Format(time.RFC3339), nil
}
