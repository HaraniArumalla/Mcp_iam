package helpers

import (
	"context"
	"iam_services_main_v1/config"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetGinContext(t *testing.T) {
	// Setup
	ginCtx := &gin.Context{}
	ctx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)
	emptyCtx := context.Background()

	// Test cases
	tests := []struct {
		name    string
		ctx     context.Context
		want    *gin.Context
		wantErr bool
	}{
		{
			name:    "Valid gin context",
			ctx:     ctx,
			want:    ginCtx,
			wantErr: false,
		},
		{
			name:    "Missing gin context",
			ctx:     emptyCtx,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGinContext(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGinContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetGinContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckValueExists(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		fallback string
		want     string
	}{
		{
			name:     "Field exists",
			field:    "value",
			fallback: "default",
			want:     "value",
		},
		{
			name:     "Field is empty",
			field:    "",
			fallback: "default",
			want:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckValueExists(tt.field, tt.fallback); got != tt.want {
				t.Errorf("CheckValueExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTenantID(t *testing.T) {
	// Setup
	validUUID := uuid.New()

	tests := []struct {
		name     string
		setupCtx func() context.Context
		want     *uuid.UUID
		wantErr  bool
	}{
		{
			name: "Valid string tenant ID",
			setupCtx: func() context.Context {
				ginCtx := &gin.Context{}
				ginCtx.Set("tenantID", validUUID.String())
				return context.WithValue(context.Background(), config.GinContextKey, ginCtx)
			},
			want:    &validUUID,
			wantErr: false,
		},
		{
			name: "Valid UUID tenant ID",
			setupCtx: func() context.Context {
				ginCtx := &gin.Context{}
				ginCtx.Set("tenantID", validUUID)
				return context.WithValue(context.Background(), config.GinContextKey, ginCtx)
			},
			want:    &validUUID,
			wantErr: false,
		},
		{
			name: "Invalid UUID string",
			setupCtx: func() context.Context {
				ginCtx := &gin.Context{}
				ginCtx.Set("tenantID", "not-a-uuid")
				return context.WithValue(context.Background(), config.GinContextKey, ginCtx)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Missing tenant ID",
			setupCtx: func() context.Context {
				ginCtx := &gin.Context{}
				return context.WithValue(context.Background(), config.GinContextKey, ginCtx)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid context",
			setupCtx: func() context.Context {
				return context.Background()
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupCtx()
			got, err := GetTenantID(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTenantID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && got != nil {
				if *tt.want != *got {
					t.Errorf("GetTenantID() = %v, want %v", *got, *tt.want)
				}
			} else if (tt.want == nil && got != nil) || (tt.want != nil && got == nil) {
				t.Errorf("GetTenantID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUserID(t *testing.T) {
	// Similar structure to TestGetTenantID
	validUUID := uuid.New()

	tests := []struct {
		name     string
		setupCtx func() context.Context
		want     *uuid.UUID
		wantErr  bool
	}{
		{
			name: "Valid string user ID",
			setupCtx: func() context.Context {
				ginCtx := &gin.Context{}
				ginCtx.Set("userID", validUUID.String())
				return context.WithValue(context.Background(), config.GinContextKey, ginCtx)
			},
			want:    &validUUID,
			wantErr: false,
		},
		{
			name: "Valid UUID user ID",
			setupCtx: func() context.Context {
				ginCtx := &gin.Context{}
				ginCtx.Set("userID", validUUID)
				return context.WithValue(context.Background(), config.GinContextKey, ginCtx)
			},
			want:    &validUUID,
			wantErr: false,
		},
		{
			name: "Invalid UUID string",
			setupCtx: func() context.Context {
				ginCtx := &gin.Context{}
				ginCtx.Set("userID", "not-a-uuid")
				return context.WithValue(context.Background(), config.GinContextKey, ginCtx)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Missing user ID",
			setupCtx: func() context.Context {
				ginCtx := &gin.Context{}
				return context.WithValue(context.Background(), config.GinContextKey, ginCtx)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupCtx()
			got, err := GetUserID(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && got != nil {
				if *tt.want != *got {
					t.Errorf("GetUserID() = %v, want %v", *got, *tt.want)
				}
			} else if (tt.want == nil && got != nil) || (tt.want != nil && got == nil) {
				t.Errorf("GetUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

type TestStruct struct {
	Name         string
	Age          int
	Active       bool
	Description  *string
	NestedStruct struct {
		Field1 string
		Field2 int
	}
	NilField *string
}

func TestStructToMap(t *testing.T) {
	// Setup
	desc := "Test Description"
	ts := TestStruct{
		Name:        "Test",
		Age:         30,
		Active:      true,
		Description: &desc,
	}
	ts.NestedStruct.Field1 = "Nested Field"
	ts.NestedStruct.Field2 = 42

	tests := []struct {
		name  string
		input interface{}
		want  map[string]interface{}
	}{
		{
			name:  "Valid struct",
			input: ts,
			want: map[string]interface{}{
				"Name":        "Test",
				"Age":         30,
				"Active":      true,
				"Description": &desc,
				"NestedStruct": map[string]interface{}{
					"Field1": "Nested Field",
					"Field2": 42,
				},
				"NilField": nil,
			},
		},
		{
			name:  "Pointer to struct",
			input: &ts,
			want: map[string]interface{}{
				"Name":        "Test",
				"Age":         30,
				"Active":      true,
				"Description": &desc,
				"NestedStruct": map[string]interface{}{
					"Field1": "Nested Field",
					"Field2": 42,
				},
				"NilField": nil,
			},
		},
		{
			name:  "Non-struct input",
			input: "not a struct",
			want:  map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StructToMap(tt.input)
			// For complex struct comparison, we'll check key by key
			for k, v := range tt.want {
				if k == "NilField" {
					continue // Skip nil field comparison
				}
				if !reflect.DeepEqual(got[k], v) {
					t.Errorf("StructToMap() key %s = %v, want %v", k, got[k], v)
				}
			}
		})
	}
}

type SourceStruct struct {
	Name    string
	Age     int
	Address string
}

type DestStruct struct {
	Name    string
	Age     int
	Email   string // Different field, should not be copied
	Address string
}

func TestMapStruct(t *testing.T) {
	// Setup
	src := SourceStruct{
		Name:    "John Doe",
		Age:     30,
		Address: "123 Main St",
	}
	dest := DestStruct{
		Email: "john@example.com", // Should remain unchanged
	}

	tests := []struct {
		name    string
		src     interface{}
		dst     interface{}
		want    *DestStruct
		wantErr bool
	}{
		{
			name: "Valid struct mapping",
			src:  &src,
			dst:  &dest,
			want: &DestStruct{
				Name:    "John Doe",
				Age:     30,
				Email:   "john@example.com", // Should remain unchanged
				Address: "123 Main St",
			},
			wantErr: false,
		},
		{
			name:    "Non-pointer source",
			src:     src,
			dst:     &dest,
			want:    &dest, // Should remain unchanged
			wantErr: true,
		},
		{
			name:    "Non-pointer destination",
			src:     &src,
			dst:     dest,
			want:    &dest, // Should remain unchanged
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Valid struct mapping" {
				// For valid case, reset dest to initial state
				dest = DestStruct{
					Email: "john@example.com",
				}
			}

			err := MapStruct(tt.src, tt.dst)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.name == "Valid struct mapping" {
				// Check fields were copied correctly for valid case
				destPtr := tt.dst.(*DestStruct)
				if destPtr.Name != tt.want.Name ||
					destPtr.Age != tt.want.Age ||
					destPtr.Email != tt.want.Email ||
					destPtr.Address != tt.want.Address {
					t.Errorf("MapStruct() = %v, want %v", *destPtr, *tt.want)
				}
			}
		})
	}
}

func TestMergeMaps(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		existing map[string]interface{}
		updates  map[string]interface{}
		want     map[string]interface{}
	}{
		{
			name: "Basic merge",
			existing: map[string]interface{}{
				"name": "John",
				"age":  30,
			},
			updates: map[string]interface{}{
				"age":     31,
				"address": "123 Main St",
			},
			want: map[string]interface{}{
				"name":    "John",
				"age":     31,
				"address": "123 Main St",
			},
		},
		{
			name: "Nested map merge",
			existing: map[string]interface{}{
				"name": "John",
				"address": map[string]interface{}{
					"street": "First St",
					"city":   "Old City",
				},
			},
			updates: map[string]interface{}{
				"address": map[string]interface{}{
					"city":  "New City",
					"state": "State",
				},
			},
			want: map[string]interface{}{
				"name": "John",
				"address": map[string]interface{}{
					"street": "First St",
					"city":   "New City",
					"state":  "State",
				},
			},
		},
		{
			name: "Handle nil values",
			existing: map[string]interface{}{
				"name": "John",
				"age":  30,
			},
			updates: map[string]interface{}{
				"age": nil,
				"new": "value",
			},
			want: map[string]interface{}{
				"name": "John",
				"age":  30, // Should not be changed by nil
				"new":  "value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeMaps(tt.existing, tt.updates)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetString(t *testing.T) {
	// Test data
	data := map[string]interface{}{
		"string_key":  "string_value",
		"int_key":     42,
		"nonexistent": nil,
	}

	tests := []struct {
		name string
		data map[string]interface{}
		key  string
		want string
	}{
		{
			name: "String value",
			data: data,
			key:  "string_key",
			want: "string_value",
		},
		{
			name: "Non-string value",
			data: data,
			key:  "int_key",
			want: "", // Should return empty string for non-string value
		},
		{
			name: "Nonexistent key",
			data: data,
			key:  "nonexistent_key",
			want: "", // Should return empty string for missing key
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetString(tt.data, tt.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetUUID(t *testing.T) {
	validUUID := uuid.New()

	// Test data
	data := map[string]interface{}{
		"valid_uuid":   validUUID.String(),
		"invalid_uuid": "not-a-uuid",
		"non_string":   42,
	}

	tests := []struct {
		name    string
		data    map[string]interface{}
		key     string
		want    uuid.UUID
		wantErr bool
	}{
		{
			name:    "Valid UUID string",
			data:    data,
			key:     "valid_uuid",
			want:    validUUID,
			wantErr: false,
		},
		{
			name:    "Invalid UUID string",
			data:    data,
			key:     "invalid_uuid",
			want:    uuid.UUID{},
			wantErr: true,
		},
		{
			name:    "Non-string value",
			data:    data,
			key:     "non_string",
			want:    uuid.UUID{},
			wantErr: true,
		},
		{
			name:    "Nonexistent key",
			data:    data,
			key:     "nonexistent_key",
			want:    uuid.UUID{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUUID(tt.data, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("GetUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMap(t *testing.T) {
	// Test data
	nestedMap := map[string]interface{}{
		"nested_key": "nested_value",
	}

	data := map[string]interface{}{
		"map_key":     nestedMap,
		"string_key":  "not_a_map",
		"nonexistent": nil,
	}

	tests := []struct {
		name    string
		data    map[string]interface{}
		key     string
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:    "Valid map",
			data:    data,
			key:     "map_key",
			want:    nestedMap,
			wantErr: false,
		},
		{
			name:    "Non-map value",
			data:    data,
			key:     "string_key",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nonexistent key",
			data:    data,
			key:     "nonexistent_key",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMap(tt.data, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringPtr(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Empty string",
			input: "",
			want:  "",
		},
		{
			name:  "Non-empty string",
			input: "test",
			want:  "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := StringPtr(tt.input)
			assert.NotNil(t, ptr)
			assert.Equal(t, tt.want, *ptr)
		})
	}
}

func TestGetSlice(t *testing.T) {
	// Test data
	slice := []interface{}{"value1", "value2"}
	typedSlice := []string{"typed1", "typed2"}

	data := map[string]interface{}{
		"slice_key":   slice,
		"typed_slice": typedSlice,
		"non_slice":   "not_a_slice",
		"nil_value":   nil,
		"nonexistent": nil,
	}

	tests := []struct {
		name    string
		data    map[string]interface{}
		key     string
		want    []interface{}
		wantErr bool
	}{
		{
			name:    "Interface slice",
			data:    data,
			key:     "slice_key",
			want:    slice,
			wantErr: false,
		},
		{
			name:    "Typed slice",
			data:    data,
			key:     "typed_slice",
			want:    []interface{}{"typed1", "typed2"}, // Should convert to []interface{}
			wantErr: false,
		},
		{
			name:    "Non-slice value",
			data:    data,
			key:     "non_slice",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nil value",
			data:    data,
			key:     "nil_value",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nonexistent key",
			data:    data,
			key:     "nonexistent_key",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSlice(tt.data, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				assert.Equal(t, len(tt.want), len(got))
				for i, v := range tt.want {
					assert.Equal(t, v, got[i])
				}
			}
		})
	}
}

func TestGetUserAndTenantID(t *testing.T) {
	// Setup

	tests := []struct {
		name       string
		setupCtx   func() context.Context
		wantUser   *uuid.UUID
		wantTenant *uuid.UUID
		wantErr    bool
	}{

		{
			name: "Invalid context",
			setupCtx: func() context.Context {
				return context.Background()
			},
			wantUser:   nil,
			wantTenant: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupCtx()
			gotUser, gotTenant, err := GetUserAndTenantID(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserAndTenantID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantUser != nil && gotUser != nil {
				if *tt.wantUser != *gotUser {
					t.Errorf("GetUserAndTenantID() got user = %v, want %v", *gotUser, *tt.wantUser)
				}
			} else if (tt.wantUser == nil && gotUser != nil) || (tt.wantUser != nil && gotUser == nil) {
				t.Errorf("GetUserAndTenantID() got user = %v, want %v", gotUser, tt.wantUser)
			}

			if tt.wantTenant != nil && gotTenant != nil {
				if *tt.wantTenant != *gotTenant {
					t.Errorf("GetUserAndTenantID() got tenant = %v, want %v", *gotTenant, *tt.wantTenant)
				}
			} else if (tt.wantTenant == nil && gotTenant != nil) || (tt.wantTenant != nil && gotTenant == nil) {
				t.Errorf("GetUserAndTenantID() got tenant = %v, want %v", gotTenant, tt.wantTenant)
			}
		})
	}
}
