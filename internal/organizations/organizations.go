package organizations

// type OrganizationQueryResolver struct {
// 	DB *gorm.DB
// }

// // Organizations resolver for fetching all organizations
// func (r *OrganizationQueryResolver) Organizations(ctx context.Context) ([]*dto.Organization, error) {
// 	var organizations []*dto.Organization
// 	if err := r.DB.Find(&organizations).Error; err != nil {
// 		return nil, err
// 	}
// 	return organizations, nil
// }

// // GetOrganization resolver for fetching a single organization by ID
// func (r *OrganizationQueryResolver) GetOrganization(ctx context.Context, id uuid.UUID) (models.Organization, error) {
// 	if id == uuid.Nil {
// 		return nil, errors.New("id cannot be nil")
// 	}

// 	var organization dto.Organization
// 	if err := r.DB.First(&organization, "organization_id = ?", id).Error; err != nil {
// 		return nil, err
// 	}
// 	return nil, nil
// }

// func (r *OrganizationQueryResolver) AllOrganizations(ctx context.Context) ([]models.Organization, error) {
// 	return nil, nil
// }
