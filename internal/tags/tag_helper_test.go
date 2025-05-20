package tags

import (
	"iam_services_main_v1/gql/models"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetResourceTags(t *testing.T) {
	testCases := []struct {
		name     string
		data     map[string]interface{}
		key      string
		expected []*models.Tags
	}{
		{
			name: "Valid input with multiple tags",
			data: map[string]interface{}{
				"tags": []interface{}{
					map[string]interface{}{
						"key":   "environment",
						"value": "production",
					},
					map[string]interface{}{
						"key":   "owner",
						"value": "team-a",
					},
				},
			},
			key: "tags",
			expected: []*models.Tags{
				{Key: "environment", Value: "production"},
				{Key: "owner", Value: "team-a"},
			},
		},
		{
			name: "Valid input with a single tag",
			data: map[string]interface{}{
				"tags": []interface{}{
					map[string]interface{}{
						"key":   "environment",
						"value": "development",
					},
				},
			},
			key: "tags",
			expected: []*models.Tags{
				{Key: "environment", Value: "development"},
			},
		},
		{
			name: "Empty tags slice",
			data: map[string]interface{}{
				"tags": []interface{}{},
			},
			key:      "tags",
			expected: nil,
		},
		{
			name: "Key doesn't exist in the map",
			data: map[string]interface{}{
				"other": "value",
			},
			key:      "tags",
			expected: nil,
		},
		{
			name: "Key exists but value is not a slice",
			data: map[string]interface{}{
				"tags": "not a slice",
			},
			key:      "tags",
			expected: nil,
		},
		{
			name: "Slice contains non-map elements",
			data: map[string]interface{}{
				"tags": []interface{}{
					"not a map",
					123,
					true,
				},
			},
			key:      "tags",
			expected: nil,
		},
		{
			name: "Slice contains maps without key and/or value fields",
			data: map[string]interface{}{
				"tags": []interface{}{
					map[string]interface{}{
						"other": "field",
					},
					map[string]interface{}{
						"key": "only-key",
					},
					map[string]interface{}{
						"value": "only-value",
					},
				},
			},
			key: "tags",
			expected: []*models.Tags{
				{Key: "", Value: ""},
				{Key: "only-key", Value: ""},
				{Key: "", Value: "only-value"},
			},
		},
		{
			name: "Mixed valid and invalid elements in slice",
			data: map[string]interface{}{
				"tags": []interface{}{
					map[string]interface{}{
						"key":   "valid",
						"value": "tag",
					},
					"not a map",
					map[string]interface{}{
						"key":   "another",
						"value": "valid-tag",
					},
				},
			},
			key: "tags",
			expected: []*models.Tags{
				{Key: "valid", Value: "tag"},
				{Key: "another", Value: "valid-tag"},
			},
		},
		{
			name: "Test with different key name",
			data: map[string]interface{}{
				"labels": []interface{}{
					map[string]interface{}{
						"key":   "label1",
						"value": "value1",
					},
				},
			},
			key: "labels",
			expected: []*models.Tags{
				{Key: "label1", Value: "value1"},
			},
		},
		{
			name:     "Test with nil map",
			data:     nil,
			key:      "tags",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetResourceTags(tc.data, tc.key)

			if tc.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, len(tc.expected), len(result))

				// Check each tag matches the expected value
				for i, expectedTag := range tc.expected {
					if i < len(result) {
						assert.Equal(t, expectedTag.Key, result[i].Key)
						assert.Equal(t, expectedTag.Value, result[i].Value)
					}
				}
			}
		})
	}
}

func TestGetResourceTagsDeepEqual(t *testing.T) {
	// Create test data
	data := map[string]interface{}{
		"tags": []interface{}{
			map[string]interface{}{
				"key":   "environment",
				"value": "production",
			},
			map[string]interface{}{
				"key":   "owner",
				"value": "team-a",
			},
		},
	}

	// Create expected result
	expected := []*models.Tags{
		{Key: "environment", Value: "production"},
		{Key: "owner", Value: "team-a"},
	}

	// Get actual result
	result := GetResourceTags(data, "tags")

	// Compare using DeepEqual (alternative way to verify complex structures)
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}
