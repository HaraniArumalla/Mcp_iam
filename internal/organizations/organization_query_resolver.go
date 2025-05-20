package organizations

import (
	"context"
	"iam_services_main_v1/gql/models"

	"github.com/google/uuid"
)

type OrganizationQueryResolver struct{}

// Organizations resolver for fetching all organizations
func (r *OrganizationQueryResolver) Organizations(ctx context.Context) (models.OperationResult, error) {
	// var organizations []*dto.Organization
	// if err := r.DB.Find(&organizations).Error; err != nil {
	// 	return nil, err
	// }
	// return organizations, nil
	return nil, nil
}

// GetOrganization resolver for fetching a single organization by ID
func (r *OrganizationQueryResolver) Organization(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	// if id == uuid.Nil {
	// 	return nil, errors.New("id cannot be nil")
	// }

	// var organization dto.Organization
	// if err := r.DB.First(&organization, "organization_id = ?", id).Error; err != nil {
	// 	return nil, err
	// }
	// return &organization, nil
	return nil, nil
}
