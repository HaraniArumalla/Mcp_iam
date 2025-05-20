package roles

import (
	"context"
	"errors"
	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	mocks "iam_services_main_v1/mocks"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	mock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateRole(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	resolver := RoleMutationResolver{
		PC: mockService,
	}

	// Setup test context
	userID := uuid.New().String()
	tenantID := uuid.New().String()

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", tenantID)
	ginCtx.Set("userID", userID)

	ginCtxNoUser := &gin.Context{}
	ginCtxNoUser.Set("tenantID", tenantID)

	ginCtxNoTenant := &gin.Context{}
	ginCtxNoTenant.Set("userID", userID)
	validCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)
	noUserCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtxNoUser)
	noTenantCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtxNoTenant)

	// Setup test input
	validID := uuid.New()
	scopeRef := uuid.New()
	desc := "Test Role Description"
	validInput := models.CreateRoleInput{
		ID:                 validID,
		Name:               "TestRole",
		Description:        &desc,
		AssignableScopeRef: scopeRef,
		RoleType:           "CUSTOM",
		Version:            "1.0",
		Permissions:        []string{"read", "write"},
	}

	invalidInput := models.CreateRoleInput{
		ID:       uuid.Nil,
		Name:     "",
		RoleType: "DEFAULT",
	}

	// Mock role data that will be returned after creation
	roleData := buildTestMutationRolesData()
	//mappedRole, _ := MapRoleResponseToStruct(roleData, validID)

	testCases := []struct {
		name      string
		ctx       context.Context
		input     models.CreateRoleInput
		mockSetup func()
		wantErr   bool
	}{
		{
			name:      "Missing user ID in context",
			ctx:       noUserCtx,
			input:     validInput,
			mockSetup: func() {},
			wantErr:   true,
		},
		{
			name:      "Missing tenant ID in context",
			ctx:       noTenantCtx,
			input:     validInput,
			mockSetup: func() {},
			wantErr:   true,
		},
		{
			name:      "Invalid input data",
			ctx:       validCtx,
			input:     invalidInput,
			mockSetup: func() {},
			wantErr:   true,
		},
		{
			name:  "Error creating role in permit",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "POST", mock.Any(), mock.Any()).
					Return(nil, errors.New("permit error")).MaxTimes(1)
			},
			wantErr: true,
		},
		{
			name:  "Error getting created role",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "POST", mock.Any(), mock.Any()).
					Return(map[string]interface{}{"created": true}, nil).MaxTimes(1)

				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(nil, errors.New("get role error")).MaxTimes(1)
			},
			wantErr: true,
		},
		{
			name:  "Success",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "POST", mock.Any(), mock.Any()).
					Return(map[string]interface{}{"key": validID.String()}, nil).MaxTimes(1)

				mockService.EXPECT().
					SendRequest(mock.Any(), "GET", mock.Any(), nil).
					Return(roleData, nil).MaxTimes(1)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			result, _ := resolver.CreateRole(tc.ctx, tc.input)
			assert.NotNil(t, result)

		})
	}
}

func TestUpdateRole(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	resolver := RoleMutationResolver{
		PC: mockService,
	}

	// Setup test context
	userID := uuid.New().String()
	tenantID := uuid.New().String()

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", tenantID)
	ginCtx.Set("userID", userID)

	validCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)

	// Setup test input
	validID := uuid.New()
	scopeRef := uuid.New()
	desc := "Updated Role Description"
	validInput := models.UpdateRoleInput{
		ID:                 validID,
		Name:               "Updated Role",
		Description:        &desc,
		AssignableScopeRef: scopeRef,
		RoleType:           "CUSTOM",
		Version:            "1.1",
		Permissions:        []string{"read", "write", "execute"},
	}

	invalidInput := models.UpdateRoleInput{
		ID:                 uuid.Nil,
		AssignableScopeRef: uuid.Nil,
	}

	// Mock role data
	roleData := buildTestMutationRolesData()

	testCases := []struct {
		name      string
		ctx       context.Context
		input     models.UpdateRoleInput
		mockSetup func()
		wantErr   bool
	}{
		{
			name:      "Invalid input data",
			ctx:       validCtx,
			input:     invalidInput,
			mockSetup: func() {},
			wantErr:   true,
		},
		{
			name:  "Error updating role in permit",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(nil, errors.New("update error")).MaxTimes(1)
			},
			wantErr: true,
		},
		{
			name:  "Error getting updated role",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(map[string]interface{}{"updated": true}, nil).MaxTimes(1)

				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), nil).
					Return(nil, errors.New("get updated role error")).MaxTimes(1)
			},
			wantErr: true,
		},
		{
			name:  "Success",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).
					Return(map[string]interface{}{"updated": true}, nil).MaxTimes(1)

				mockService.EXPECT().
					SendRequest(mock.Any(), mock.Any(), mock.Any(), nil).
					Return(roleData, nil).MaxTimes(1)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			result, _ := resolver.UpdateRole(tc.ctx, tc.input)

			assert.NotNil(t, result)
		})
	}
}

func TestDeleteRole(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	resolver := RoleMutationResolver{
		PC: mockService,
	}

	// Setup test context
	userID := uuid.New().String()
	tenantID := uuid.New().String()
	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", tenantID)
	ginCtx.Set("userID", userID)
	validCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)

	validID := uuid.New()
	scopeRef := uuid.New()
	validInput := models.DeleteRoleInput{
		ID:                 validID,
		AssignableScopeRef: scopeRef,
	}

	invalidInput := models.DeleteRoleInput{
		ID:                 uuid.Nil,
		AssignableScopeRef: uuid.Nil,
	}

	testCases := []struct {
		name      string
		ctx       context.Context
		input     models.DeleteRoleInput
		mockSetup func()
		wantErr   bool
	}{
		{
			name:      "Invalid input - nil ID",
			ctx:       validCtx,
			input:     invalidInput,
			mockSetup: func() {},
			wantErr:   true,
		},
		{
			name:  "Error deleting role",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "DELETE", mock.Any(), nil).
					Return(nil, errors.New("delete error")).MaxTimes(1)
			},
			wantErr: true,
		},
		{
			name:  "Success",
			ctx:   validCtx,
			input: validInput,
			mockSetup: func() {
				mockService.EXPECT().
					SendRequest(mock.Any(), "DELETE", mock.Any(), nil).
					Return(map[string]interface{}{"deleted": true}, nil).MaxTimes(1)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			result, _ := resolver.DeleteRole(tc.ctx, tc.input)

			assert.NotNil(t, result)
		})
	}
}

func TestPrepareMetadataForCreateInput(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	resolver := RoleMutationResolver{}

	// Setup test data
	userID := uuid.New()
	tenantID := uuid.New()
	desc := "Test Role Description"
	validInput := models.CreateRoleInput{
		ID:                 uuid.New(),
		Name:               "Test Role",
		Description:        &desc,
		AssignableScopeRef: uuid.New(),
		RoleType:           "CUSTOM",
		Version:            "1.0",
		Permissions:        []string{"read", "write"},
	}

	t.Run("Success", func(t *testing.T) {
		metadata, err := resolver.prepareMetadataForCreateInput(validInput, &userID, &tenantID)

		assert.NoError(t, err)
		assert.NotNil(t, metadata)
		assert.Equal(t, validInput.Name, metadata["name"])
		assert.Equal(t, validInput.Description, metadata["description"])
		assert.Equal(t, validInput.RoleType, metadata["roleType"])
		assert.Equal(t, validInput.Permissions, metadata["permissions"])
		assert.Equal(t, userID, metadata["createdBy"])
		assert.Equal(t, userID, metadata["updatedBy"])
		assert.Equal(t, tenantID, metadata["tenantId"])
	})
}

func TestPrepareMetadataForUpdateInput(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	resolver := RoleMutationResolver{}

	// Setup test data
	userID := uuid.New()
	tenantID := uuid.New()
	desc := "Updated Role Description"
	validInput := models.UpdateRoleInput{
		ID:                 uuid.New(),
		Name:               "Updated Role",
		Description:        &desc,
		AssignableScopeRef: uuid.New(),
		RoleType:           "CUSTOM",
		Version:            "1.1",
		Permissions:        []string{"read", "write", "execute"},
	}

	t.Run("Success", func(t *testing.T) {
		metadata, err := resolver.prepareMetadataForUpdateInput(validInput, &userID, &tenantID)

		assert.NoError(t, err)
		assert.NotNil(t, metadata)
		assert.Equal(t, validInput.Name, metadata["name"])
		assert.Equal(t, validInput.Description, metadata["description"])
		assert.Equal(t, validInput.RoleType, metadata["roleType"])
		assert.Equal(t, validInput.Permissions, metadata["permissions"])
		assert.Equal(t, userID, metadata["updatedBy"])
		assert.Equal(t, tenantID, metadata["tenantId"])
	})
}

func buildTestMutationRolesData() map[string]interface{} {
	return map[string]interface{}{
		"data": []interface{}{
			map[string]interface{}{
				"key":         uuid.New().String(),
				"name":        "Admin Role",
				"description": "Administrator role with full access",
				"attributes": map[string]interface{}{
					"displayName": "Administrator",
					"permissions": []interface{}{
						"read:all",
						"write:all",
						"delete:all",
					},
					"createdAt": time.Now().Format(time.RFC3339),
					"updatedAt": time.Now().Format(time.RFC3339),
					"createdBy": uuid.New().String(),
					"updatedBy": uuid.New().String(),
				},
			},
		},
	}
}
