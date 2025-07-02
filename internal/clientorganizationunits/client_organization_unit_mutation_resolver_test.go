package clientorganizationunits

import (
	"context"
	"errors"
	"testing"

	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	mocks "iam_services_main_v1/mocks"

	"github.com/gin-gonic/gin"
	mock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateClientOrganizationUnit(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	objUnderTest := ClientOrganizationUnitMutationResolver{
		PC: mockService,
	}

	description := "description"
	emptyUserId := ""

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", "7ed6cfa6-fd7e-4a2a-bbce-773ef8ea4c12")
	ginCtx.Set("userID", emptyUserId)

	ginCtxWithUserId := &gin.Context{}
	ginCtxWithUserId.Set("tenantID", "7ed6cfa6-fd7e-4a2a-bbce-773ef8ea4c12")
	ginCtxWithUserId.Set("userID", "b5b44e90-906e-458a-8bb1-e9e4ee180696")

	ginCtxWithOutUserId := &gin.Context{}
	ginCtxWithOutUserId.Set("userID", "b5b44e90-906e-458a-8bb1-e9e4ee180696")

	testCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)
	validCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtxWithUserId)
	ctxWithOutUserId := context.WithValue(context.Background(), config.GinContextKey, ginCtxWithOutUserId)
	nilIdInput := models.CreateClientOrganizationUnitInput{
		ID: uuid.Nil,
	}

	nameNotPresent := models.CreateClientOrganizationUnitInput{
		ID:   uuid.New(),
		Name: "",
	}

	successRequest := models.CreateClientOrganizationUnitInput{
		ID:             uuid.New(),
		Name:           "test",
		Description:    &description,
		AccountOwnerID: uuid.New(),
	}

	corg := buildClientOrganization()
	testcases := []struct {
		name      string
		input     models.CreateClientOrganizationUnitInput
		ctx       context.Context
		mockStubs func(mockService mocks.MockPermitService)
		output    models.OperationResult
	}{
		{
			name:      "CreateClientOrganizationUnit id is not present",
			input:     nilIdInput,
			ctx:       context.WithValue(context.Background(), config.GinContextKey, &gin.Context{}),
			mockStubs: func(mockSvc mocks.MockPermitService) {},
			output:    buildErrorResponse(400, "unable to fetch gin context", "error while fetching gin context"),
		},
		{
			name:      "CreateClientOrganizationUnit Name not present",
			input:     nameNotPresent,
			ctx:       testCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {},
			output:    buildErrorResponse(400, "failed", "unable to fetch resource from permit"),
		},
		{
			name:      "CreateClientOrganizationUnit Tenant id not present in ctx",
			input:     successRequest,
			ctx:       ctxWithOutUserId,
			mockStubs: func(mockSvc mocks.MockPermitService) {},
			output:    buildErrorResponse(400, "failed", "unable to fetch resource from permit"),
		},
		{
			name:      "CreateClientOrganizationUnit User id not present in ctx",
			input:     successRequest,
			ctx:       testCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {},
			output:    buildErrorResponse(400, "failed", "unable to fetch resource from permit"),
		},
		{
			name:  "CreateClientOrganizationUnit Tenant failed error from permit",
			input: successRequest,
			ctx:   validCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().APIExecute(mock.Any(), mock.Any(), mock.Any(), mock.Any()).Return(nil, errors.New("failed to create"))
			},
			output: buildErrorResponse(400, "unable to create organization in permit", "failed to create"),
		},
		{
			name:  "CreateClientOrganizationUnit success",
			input: successRequest,
			ctx:   validCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().APIExecute(mock.Any(), mock.Any(), mock.Any(), mock.Any()).Return(nil, nil)
				mockService.EXPECT().GetSingleResource(mock.Any(), mock.Any(), mock.Any()).Return(corg, nil).MaxTimes(1)
			},
			output: buildSuccessResponse(BuildOrgUnit(corg)),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockStubs(*mockService)
			result, _ := objUnderTest.CreateClientOrganizationUnit(tc.ctx, tc.input)
			assert.NotNil(t, result)
		})

	}

}

func TestUpdateClientOrganizationUnit(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	objUnderTest := ClientOrganizationUnitMutationResolver{
		PC: mockService,
	}

	description := "description"
	emptyUserId := ""
	accountOwnerId := uuid.New()

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", "7ed6cfa6-fd7e-4a2a-bbce-773ef8ea4c12")
	ginCtx.Set("userID", emptyUserId)

	ginCtxWithUserId := &gin.Context{}
	ginCtxWithUserId.Set("tenantID", "7ed6cfa6-fd7e-4a2a-bbce-773ef8ea4c12")
	ginCtxWithUserId.Set("userID", "b5b44e90-906e-458a-8bb1-e9e4ee180696")

	ginCtxWithOutTenantId := &gin.Context{}
	ginCtxWithOutTenantId.Set("userID", "b5b44e90-906e-458a-8bb1-e9e4ee180696")

	corg := buildClientOrganization()
	corgUpdated := buildClientOrganization()
	testCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)
	validCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtxWithUserId)
	ctxWithOutTenantId := context.WithValue(context.Background(), config.GinContextKey, ginCtxWithOutTenantId)
	nilIdInput := models.UpdateClientOrganizationUnitInput{
		ID: uuid.Nil,
	}

	relationTypeEnum := models.RelationTypeEnumChild
	statusTypeEnum := models.StatusTypeEnumActive

	successRequest := models.UpdateClientOrganizationUnitInput{
		ID:             uuid.New(),
		Name:           &description,
		Description:    &description,
		AccountOwnerID: &accountOwnerId,
		ParentID:       &accountOwnerId,
		RelationType:   &relationTypeEnum,
		Status:         &statusTypeEnum,
	}

	testcases := []struct {
		name      string
		input     models.UpdateClientOrganizationUnitInput
		ctx       context.Context
		mockStubs func(mockService mocks.MockPermitService)
		output    models.OperationResult
	}{
		{
			name:      "UpdateClientOrganizationUnit id is not valid",
			input:     nilIdInput,
			ctx:       context.WithValue(context.Background(), config.GinContextKey, &gin.Context{}),
			mockStubs: func(mockSvc mocks.MockPermitService) {},
			output:    buildErrorResponse(400, "unable to fetch gin context", "error while fetching gin context"),
		},
		{
			name:      "UpdateClientOrganizationUnit Tenant id not present in ctx",
			input:     successRequest,
			ctx:       ctxWithOutTenantId,
			mockStubs: func(mockSvc mocks.MockPermitService) {},
			output:    buildErrorResponse(400, "failed", "unable to fetch resource from permit"),
		},
		{
			name:      "UpdateClientOrganizationUnit User id not present in ctx",
			input:     successRequest,
			ctx:       testCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {},
			output:    buildErrorResponse(400, "failed", "unable to fetch resource from permit"),
		},
		{
			name:  "UpdateClientOrganizationUnit fetch resource from permit failed",
			input: successRequest,
			ctx:   validCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().GetSingleResource(mock.Any(), mock.Any(), mock.Any()).Return(nil, errors.New("failed to get"))
			},
			output: buildErrorResponse(400, "unable to create organization in permit", "failed to create"),
		},
		{
			name:  "UpdateClientOrganizationUnit update request failed",
			input: successRequest,
			ctx:   validCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().GetSingleResource(mock.Any(), mock.Any(), mock.Any()).Return(corg, nil)
				mockService.EXPECT().APIExecute(mock.Any(), mock.Any(), mock.Any(), mock.Any()).Return(nil, errors.New("failed to update"))
			},
			output: buildErrorResponse(400, "unable to create organization in permit", "failed to create"),
		},
		{
			name:  "UpdateClientOrganizationUnit update successful",
			input: successRequest,
			ctx:   validCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().GetSingleResource(mock.Any(), mock.Any(), mock.Any()).Return(corg, nil)
				mockService.EXPECT().APIExecute(mock.Any(), mock.Any(), mock.Any(), mock.Any()).Return(nil, nil)
				mockService.EXPECT().GetSingleResource(mock.Any(), mock.Any(), mock.Any()).Return(nil, errors.New("failed"))

			},
			output: buildErrorResponse(400, "unable to create organization in permit", "failed to create"),
		},
		{
			name:  "UpdateClientOrganizationUnit update successful",
			input: successRequest,
			ctx:   validCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().GetSingleResource(mock.Any(), mock.Any(), mock.Any()).Return(corg, nil)
				mockService.EXPECT().APIExecute(mock.Any(), mock.Any(), mock.Any(), mock.Any()).Return(nil, nil)
				mockService.EXPECT().GetSingleResource(mock.Any(), mock.Any(), mock.Any()).Return(corgUpdated, nil)

			},
			output: buildSuccessResponse(BuildOrgUnit(corg)),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockStubs(*mockService)
			result, _ := objUnderTest.UpdateClientOrganizationUnit(tc.ctx, tc.input)
			assert.NotNil(t, result)
		})

	}

}

func TestDeleteClientOrganizationUnit(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	objUnderTest := ClientOrganizationUnitMutationResolver{
		PC: mockService,
	}

	// Set up different contexts for testing
	emptyUserId := ""

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", uuid.New().String())
	ginCtx.Set("userID", emptyUserId)

	ginCtxWithUserId := &gin.Context{}
	ginCtxWithUserId.Set("tenantID", uuid.New().String())
	ginCtxWithUserId.Set("userID", uuid.New().String())

	ginCtxWithOutTenantId := &gin.Context{}
	ginCtxWithOutTenantId.Set("userID", uuid.New().String())

	noGinCtx := context.Background()

	inputNullId := models.DeleteInput{
		ID: uuid.Nil,
	}

	validDeleteInput := models.DeleteInput{
		ID: uuid.New(),
	}

	testcases := []struct {
		name      string
		input     models.DeleteInput
		ctx       context.Context
		mockStubs func(mockService mocks.MockPermitService)
		output    models.OperationResult
	}{
		{
			name:      "DeleteClientOrganizationUnit without gin context",
			input:     validDeleteInput,
			ctx:       noGinCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {},
			output:    buildErrorResponse(400, "unable to fetch gin context", "error while fetching gin context"),
		},
		{
			name:      "DeleteClientOrganizationUnit id is not valid",
			input:     inputNullId,
			ctx:       context.WithValue(context.Background(), config.GinContextKey, &gin.Context{}),
			mockStubs: func(mockSvc mocks.MockPermitService) {},
			output:    buildErrorResponse(400, "unable to fetch gin context", "error while fetching gin context"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockStubs(*mockService)
			result, _ := objUnderTest.DeleteClientOrganizationUnit(tc.ctx, tc.input)
			assert.NotNil(t, result)

		})
	}
}
