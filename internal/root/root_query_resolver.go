package root

import (
	"context"
	"iam_services_main_v1/gql/models"

	"github.com/google/uuid"
)

type RootQueryResolver struct{}

func (r *RootQueryResolver) Root(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	return nil, nil
}
