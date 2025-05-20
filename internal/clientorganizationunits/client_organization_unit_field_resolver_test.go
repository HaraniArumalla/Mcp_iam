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
)

func TestFieldResolverTenant(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	objUnderTest := ClientOrganizationUnitResolver{
		PC: mockService,
	}

	emptyUserId := ""

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", "7ed6cfa6-fd7e-4a2a-bbce-773ef8ea4c12")
	ginCtx.Set("userID", emptyUserId)

	ginCtxWithUserId := &gin.Context{}
	ginCtxWithUserId.Set("tenantID", "7ed6cfa6-fd7e-4a2a-bbce-773ef8ea4c12")
	ginCtxWithUserId.Set("userID", "b5b44e90-906e-458a-8bb1-e9e4ee180696")

	ginCtxWithOutUserId := &gin.Context{}
	ginCtxWithOutUserId.Set("userID", "b5b44e90-906e-458a-8bb1-e9e4ee180696")

	validCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtxWithUserId)

	corg := buildClientOrganizationData()
	tenantMap := buildTenant()
	tenant := &models.Tenant{
		ID:   uuid.New(),
		Name: "tenant",
	}
	testcases := []struct {
		name      string
		input     *models.ClientOrganizationUnit
		ctx       context.Context
		mockStubs func(mockService mocks.MockPermitService)
		output    *models.Tenant
	}{

		{
			name:  "Fetch Tenant failure",
			input: corg,
			ctx:   validCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).Return(nil, errors.New("failed to fetch"))
			},
			output: nil,
		},
		{
			name:  "Fetch Tenant success",
			input: corg,
			ctx:   validCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).Return(tenantMap, nil)
			},
			output: tenant,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockStubs(*mockService)
			_, _ = objUnderTest.Tenant(tc.ctx, tc.input)
			// assert.Equal(t, result, tc.output)
		})

	}

}

func buildTenant() map[string]interface{} {
	tenantMap := make(map[string]interface{}, 0)
	attributes := make(map[string]interface{}, 0)
	contactInfo := make(map[string]interface{}, 0)
	address := make(map[string]interface{}, 0)
	tenantMap["key"] = uuid.New().String()
	tenantMap["created_at"] = time.Now().String()
	tenantMap["updated_at"] = time.Now().String()

	contactInfo["email"] = "test@tmobile.com"
	contactInfo["phoneNumber"] = "123-45-6789"
	address["city"] = "newyork"
	address["state"] = "NY"
	address["street"] = "1234 test street"
	address["country"] = "USA"
	address["zipcode"] = "11111"
	contactInfo["address"] = address

	attributes["created_by"] = uuid.New().String()
	attributes["updated_by"] = uuid.New().String()
	attributes["description"] = "description"
	attributes["name"] = "tenant"
	attributes["contactInfo"] = contactInfo
	tenantMap["attributes"] = attributes
	return tenantMap

}
