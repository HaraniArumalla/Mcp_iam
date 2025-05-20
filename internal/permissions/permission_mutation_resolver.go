package permissions

import (
	"context"
	"iam_services_main_v1/gql/models"
)

type PermissionMutationResolver struct{}

func (r *PermissionMutationResolver) CreatePermission(ctx context.Context, input models.CreatePermissionInput) (models.OperationResult, error) {

	return nil, nil
}

func (r *PermissionMutationResolver) DeletePermission(ctx context.Context, input models.DeleteInput) (models.OperationResult, error) {

	return nil, nil
}

func (r *PermissionMutationResolver) UpdatePermission(ctx context.Context, input models.UpdatePermissionInput) (models.OperationResult, error) {
	return nil, nil
}
