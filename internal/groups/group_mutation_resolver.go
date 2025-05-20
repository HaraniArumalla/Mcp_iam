package groups

// type GroupMutationResolver struct {
// 	DB *gorm.DB
// }

// // CreateGroup resolver for adding a new Group
// func (r *GroupMutationResolver) CreateGroup(ctx context.Context, input models.GroupInput) (*dto.GroupEntity, error) {
// 	log.Println("Creating group with input:", input)

// 	dtoInput := dto.GroupInput{
// 		Name:     input.Name,
// 		TenantID: input.TenantID,
// 	}
// 	if err := helpers.ValidateStruct(dtoInput); err != nil {
// 		return nil, err
// 	}
// 	Group := &dto.GroupEntity{Name: dtoInput.Name, TenantID: dtoInput.TenantID}

// 	if err := r.DB.Create(Group).Error; err != nil {
// 		return nil, err
// 	}
// 	log.Println("Creating group with output 1:", Group)
// 	return Group, nil
// }

// // CreateGroup resolver for adding a new Group
// func (r *GroupMutationResolver) UpdateGroup(ctx context.Context, id string, input models.GroupInput) (*dto.GroupEntity, error) {
// 	var Group dto.GroupEntity
// 	if err := r.DB.First(&Group, id).Error; err != nil {
// 		return nil, err
// 	}

// 	Group.Name = input.Name

// 	Group.TenantID = input.TenantID

// 	if err := r.DB.Save(&Group).Error; err != nil {
// 		return nil, err
// 	}
// 	return &Group, nil
// }

// // CreateGroup resolver for adding a new Group
// func (r *GroupMutationResolver) DeleteGroup(ctx context.Context, id string) (bool, error) {
// 	var Group dto.GroupEntity
// 	if err := r.DB.First(&Group, id).Error; err != nil {
// 		return false, err
// 	}
// 	/*
// 		if input.Name != nil {
// 			Group.Name = input.Name
// 		}

// 		if input.TenantID != nil {
// 			Group.TenantId = input.TenantId
// 		}*/

// 	if err := r.DB.Delete(&Group).Error; err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }
