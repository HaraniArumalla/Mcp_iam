package clientorganizationunits

import (
	"context"
	"errors"
	"testing"
	"time"

	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	mocks "iam_services_main_v1/mocks"

	"github.com/gin-gonic/gin"
	mock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestClientOrganizationUnitQueryResolver(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	objUnderTest := ClientOrganizationUnitQueryResolver{
		PC: mockService,
	}

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", "7ed6cfa6-fd7e-4a2a-bbce-773ef8ea4c12")
	testCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)

	corg := buildClientOrganization()
	testcases := []struct {
		name      string
		input     uuid.UUID
		ctx       context.Context
		mockStubs func(mockService mocks.MockPermitService)
		output    models.OperationResult
	}{
		{
			name:  "Test GetClientOrganizationUnitByID when client returns error",
			input: uuid.New(),
			ctx:   testCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().GetSingleResource(mock.Any(), mock.Any(), mock.Any()).Return(nil, errors.New("failed"))
			},
			output: buildErrorResponse(400, "failed", "unable to fetch resource from permit"),
		},

		{
			name:  "Test GetClientOrganizationUnitByID success",
			input: uuid.New(),
			ctx:   testCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().GetSingleResource(mock.Any(), mock.Any(), mock.Any()).Return(corg, nil)
			},
			output: buildSuccessResponse(BuildOrgUnit(corg)),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockStubs(*mockService)
			result, _ := objUnderTest.ClientOrganizationUnit(tc.ctx, tc.input)
			assert.NotNil(t, result)
		})

	}

}

func TestClientOrganizationUnitsQueryResolver(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	objUnderTest := ClientOrganizationUnitQueryResolver{
		PC: mockService,
	}

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", "7ed6cfa6-fd7e-4a2a-bbce-773ef8ea4c12")
	testCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)

	corg := buildClientOrganization()
	corgs := make([]interface{}, 0)
	corgs = append(corgs, corg)
	corgResult := make(map[string]interface{}, 0)
	corgResult["data"] = corgs
	testcases := []struct {
		name      string
		ctx       context.Context
		mockStubs func(mockService mocks.MockPermitService)
		output    models.OperationResult
	}{
		{
			name: "Test GetClientOrganizationUnitByID when client returns error",
			ctx:  testCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().SendRequest(mock.Any(), mock.Any(), mock.Any(), nil).Return(nil, errors.New("failed"))
			},
			output: buildErrorResponse(400, "failed", "unable to fetch resource from permit"),
		},

		{
			name: "Test GetClientOrganizationUnitByID success",
			ctx:  testCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).Return(corgResult, nil)
			},
			output: buildSuccessResponse(BuildOrgUnit(corg)),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockStubs(*mockService)
			result, _ := objUnderTest.ClientOrganizationUnits(tc.ctx)
			assert.NotNil(t, result)
		})

	}

}

func buildClientOrganizationData() *models.ClientOrganizationUnit {
	return &models.ClientOrganizationUnit{
		ParentOrg: nil,
		Tenant: &models.Tenant{
			ID: uuid.New(),
		},
	}
}

func buildClientOrganization() map[string]interface{} {
	id := uuid.New()
	result := make(map[string]interface{}, 0)
	attributes := make(map[string]interface{}, 0)
	attributes["tenantId"] = id.String()
	attributes["parentOrgId"] = "parentOrgId"
	attributes["created_at"] = time.Now().String()
	attributes["updated_at"] = time.Now().String()
	attributes["created_by"] = uuid.New().String()
	attributes["updated_by"] = uuid.New().String()
	attributes["description"] = "description"
	attributes["name"] = "clientorg"
	attributes["key"] = id.String()
	attributes["relation_type"] = "PARENT"
	attributes["status"] = "ACTIVE"
	attributes["account_owner_id"] = uuid.New().String()
	result["key"] = id.String()
	result["name"] = "clientorg"
	result["attributes"] = attributes

	return result
}
