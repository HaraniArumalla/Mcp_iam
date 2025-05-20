package groups

// type GroupFieldResolver struct {
// 	DB *gorm.DB
// }

// // ID resolves the ID field on the Article type
// func (r *GroupFieldResolver) ID(ctx context.Context, obj *dto.GroupEntity) (string, error) {
// 	return obj.ID, nil
// }

// // Implement the CreatedAt method for Group
// func (r *GroupFieldResolver) CreatedAt(ctx context.Context, obj *dto.GroupEntity) (*string, error) {
// 	// Assuming the "createdAt" field is a time.Time object, format it as a string
// 	createdAtStr := obj.CreatedAt.Format(time.RFC3339)
// 	return &createdAtStr, nil
// }

// // Implement the CreatedAt method for Group
// func (r *GroupFieldResolver) UpdatedAt(ctx context.Context, obj *dto.GroupEntity) (*string, error) {
// 	// Assuming the "createdAt" field is a time.Time object, format it as a string
// 	UpdatedAtStr := obj.UpdatedAt.Format(time.RFC3339)
// 	return &UpdatedAtStr, nil
// }

// func (r *GroupFieldResolver) Tenant(ctx context.Context, obj *dto.GroupEntity) (*dto.TenantResource, error) {
// 	var tenant dto.TenantResource
// 	if err := r.DB.First(&tenant, obj.TenantID).Error; err != nil {
// 		return nil, err
// 	}
// 	return &tenant, nil
// }
