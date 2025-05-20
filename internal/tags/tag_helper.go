package tags

import (
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
)

// GetResourceTags extracts and returns a slice of tags from the given data map using the specified key.
// It expects the value to be a slice of maps with "key" and "value" fields, returning an error if the key
func GetResourceTags(data map[string]interface{}, key string) []*models.Tags {
	var tags []*models.Tags
	if val, ok := data[key].([]interface{}); ok {
		for _, v := range val {
			tag := &models.Tags{}
			if tagMap, ok := v.(map[string]interface{}); ok {
				tag.Key = helpers.GetString(tagMap, "key")
				tag.Value = helpers.GetString(tagMap, "value")
				tags = append(tags, tag)
			}
		}
		return tags
	}
	return nil
}
