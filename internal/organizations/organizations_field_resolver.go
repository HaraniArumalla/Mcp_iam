package organizations

// type OrganizationFieldResolver struct {
// 	DB *gorm.DB
// }

// // ID resolves the ID field on the Article type
// func (r *OrganizationFieldResolver) ID(ctx context.Context, obj *dto.Organization) (uuid.UUID, error) {
// 	return obj.OrganizationID, nil
// }

// // Implement the CreatedAt method for Organization
// func (r *OrganizationFieldResolver) CreatedAt(ctx context.Context, obj *dto.Organization) (*string, error) {
// 	// Assuming the "createdAt" field is a time.Time object, format it as a string
// 	createdAtStr := obj.CreatedAt.Format(time.RFC3339)
// 	return &createdAtStr, nil
// }

// // Implement the CreatedAt method for Organization
// func (r *OrganizationFieldResolver) UpdatedAt(ctx context.Context, obj *dto.Organization) (*string, error) {
// 	// Assuming the "createdAt" field is a time.Time object, format it as a string
// 	UpdatedAtStr := obj.UpdatedAt.Format(time.RFC3339)
// 	return &UpdatedAtStr, nil
// }

// func (r *OrganizationFieldResolver) ParentOrganization(ctx context.Context, obj *dto.Organization) (*dto.Organization, error) {
// 	if obj.ParentOrgId == uuid.Nil {
// 		return nil, nil // or return nil if there's no parent
// 	}
// 	var organization dto.Organization
// 	if err := r.DB.First(&organization, "organization_id = ? ", obj.ParentOrgId).Error; err != nil {
// 		return nil, err
// 	}
// 	return &organization, nil
// }
