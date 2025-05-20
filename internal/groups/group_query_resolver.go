package groups

import (
	"context"
	"iam_services_main_v1/gql/models"
)

type GroupQueryResolver struct{}

// Tenants resolver for fetching all Tenants
func (r *GroupQueryResolver) Groups(ctx context.Context) (models.OperationResult, error) {
	// var Groups []*dto.GroupEntity
	// if err := r.DB.Find(&Groups).Error; err != nil {
	// 	return nil, err
	// }
	// return Groups, nil
	return nil, nil
}

// GetTenant resolver for fetching a single Tenant by ID
func (r *GroupQueryResolver) Group(ctx context.Context, id string) (models.OperationResult, error) {
	// if id == "" {
	// 	return nil, errors.New("id cannot be nil")
	// }

	// var Group dto.GroupEntity
	// if err := r.DB.First(&Group, id).Error; err != nil {
	// 	return nil, err
	// }
	// return &Group, nil
	return nil, nil
}
