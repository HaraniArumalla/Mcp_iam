package users

import (
	"context"
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/internal/permit"

	"github.com/google/uuid"
)

// UserResolver handles database queries and permission checks for user-related operations using Permit.io client
type UserResolver struct {
	PC permit.PermitService
}

// GetUser retrieves a user by their UUID from the resource provider.
// It constructs the endpoint URL using the user ID, makes a GET request,
// and maps the response to a models.User object.
// If the request fails, it returns the error.
func (r *UserResolver) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	url := fmt.Sprintf("users/%s", id.String())
	user, err := r.PC.GetSingleResource(ctx, "GET", url)
	if err != nil {
		return nil, err
	}
	firstName := helpers.GetString(user, "first_name")
	lastName := helpers.GetString(user, "last_name")
	email := helpers.GetString(user, "email")
	createdAt := helpers.GetString(user, "created_at")
	updatedAt := helpers.GetString(user, "updated_at")
	return &models.User{
		ID:        uuid.MustParse(user["key"].(string)),
		FirstName: firstName,
		LastName:  lastName,
		Name:      firstName + " " + lastName,
		Email:     email,
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
	}, nil
}
