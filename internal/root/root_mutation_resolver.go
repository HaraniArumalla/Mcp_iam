package root

import (
	"context"
	"iam_services_main_v1/gql/models"
)

type RootMutationResolver struct{}

// CreateAccount resolver for adding a new Account
func (r *RootMutationResolver) CreateRoot(ctx context.Context, input models.CreateRootInput) (models.OperationResult, error) {
	return nil, nil

}

// CreateAccount resolver for adding a new Account
func (r *RootMutationResolver) UpdateRoot(ctx context.Context, input models.UpdateRootInput) (models.OperationResult, error) {
	return nil, nil

}

// CreateAccount resolver for adding a new Account
func (r *RootMutationResolver) DeleteRoot(ctx context.Context, input models.DeleteInput) (models.OperationResult, error) {
	return nil, nil
}
