// groups/group_mutation_resolver_test.go
package groups

// import (
// 	"context"
// 	"fmt"
// 	"iam_services_main_v1/gql/models"
// 	"log"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// Mock generateUniqueID function for testing
// var generateUniqueID = func() string {
// 	return "12345" // Return a fixed ID for testing purposes
// }

// Test CreateGroup function
// func TestCreateGroup(t *testing.T) {
// 	err1 := "Validation failed: field1 cannot be empty"
// 	err2 := "Validation failed: field2 must be greater than 0"
// 	tests := []struct {
// 		name          string
// 		input         models.GroupInput
// 		expectedError error
// 		expectedID    string
// 	}{
// 		{
// 			name: "Valid Input",
// 			input: models.GroupInput{
// 				Name:     "Valid Group",
// 				TenantID: "10",
// 			},
// 			expectedError: nil,
// 			expectedID:    "12345", // Mocked ID
// 		},
// 		{
// 			name: "Invalid Name (empty)",
// 			input: models.GroupInput{
// 				Name:     "",
// 				TenantID: "10",
// 			},
// 			expectedError: fmt.Errorf("validation error %v", err1),
// 			expectedID:    "",
// 		},
// 		{
// 			name: "Invalid TenantID (too small)",
// 			input: models.GroupInput{
// 				Name:     "Valid Group",
// 				TenantID: "0",
// 			},
// 			expectedError: fmt.Errorf("validation error %v", err2),
// 			expectedID:    "",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			log.Println("Creating group with input 1:", tt.input)
// 			// Create a mock resolver
// 			resolver := &GroupMutationResolver{}

// 			// Call the CreateGroup function
// 			groupEntity, err := resolver.CreateGroup(context.Background(), tt.input)
// 			log.Println("Creating group with ouput 1:", groupEntity, err)
// 			// Check if the error matches the expected error
// 			if tt.expectedError != nil {
// 				assert.Equal(t, tt.expectedError, err.Error())
// 				assert.Nil(t, groupEntity) // Ensure no entity is returned in case of error
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.expectedID, groupEntity.ID) // Ensure the correct ID is assigned
// 				assert.Equal(t, tt.input.Name, groupEntity.Name)
// 				assert.Equal(t, tt.input.TenantID, groupEntity.TenantID)
// 			}
// 		})
// 	}
// }
