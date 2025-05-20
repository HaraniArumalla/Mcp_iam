package permissions

import (
	"context"
	"iam_services_main_v1/gql/models"

	"github.com/google/uuid"
)

// PermissionQueryResolver handles permission-related queries.
type PermissionQueryResolver struct{}

// Permissions resolves the list of all permissions.
func (r *PermissionQueryResolver) Permissions(ctx context.Context) (models.OperationResult, error) {

	return nil, nil
}

// GetPermission resolves a single permission by ID.
func (r *PermissionQueryResolver) Permission(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {

	return nil, nil
}
