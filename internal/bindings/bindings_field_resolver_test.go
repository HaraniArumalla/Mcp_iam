package bindings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	mocks "iam_services_main_v1/mocks"
	"iam_services_main_v1/pkg/logger"

	"github.com/gin-gonic/gin"
	mock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestFieldResolverBindings(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	objUnderTest := BindingsResolver{
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

	binding := buildBinding()
	principal := buildPrincipal()
	user := models.User{
		ID:   uuid.New(),
		Name: "tenant",
	}
	testcases := []struct {
		name      string
		input     *models.Binding
		ctx       context.Context
		mockStubs func(mockService mocks.MockPermitService)
		output    models.Principal
	}{

		{
			name:  "Fetch Tenant failure",
			input: binding,
			ctx:   validCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().GetSingleResource(mock.Any(), mock.Any(), mock.Any()).Return(nil, errors.New("failed to fetch"))
			},
			output: nil,
		},
		{
			name:  "Fetch Tenant success",
			input: binding,
			ctx:   validCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().GetSingleResource(mock.Any(), mock.Any(), mock.Any()).Return(principal, nil)
			},
			output: user,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockStubs(*mockService)
			_, _ = objUnderTest.Principal(tc.ctx, tc.input)
			// assert.Equal(t, result, tc.output)
		})

	}

}

func TestFieldResolverBindingsRole(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	objUnderTest := BindingsResolver{
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

	binding := buildBinding()
	bindingData := fetchBindingData()
	user := models.User{
		ID:   uuid.New(),
		Name: "tenant",
	}
	testcases := []struct {
		name      string
		input     *models.Binding
		ctx       context.Context
		mockStubs func(mockService mocks.MockPermitService)
		output    models.Principal
	}{

		{
			name:  "Fetch Tenant failure",
			input: binding,
			ctx:   validCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).Return(nil, errors.New("failed to fetch"))
			},
			output: nil,
		},
		{
			name:  "Fetch Tenant success",
			input: binding,
			ctx:   validCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().SendRequest(mock.Any(), mock.Any(), mock.Any(), mock.Any()).Return(bindingData, nil)
			},
			output: user,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockStubs(*mockService)
			_, _ = objUnderTest.Role(tc.ctx, tc.input)
			// assert.Equal(t, result, tc.output)
		})

	}

}

func buildBinding() *models.Binding {
	return &models.Binding{
		ID: uuid.New(),
		Principal: &models.User{
			ID: uuid.New(),
		},
		Role: &models.Role{
			ID: uuid.MustParse("a3cfb093-280f-5b18-8c7d-a96df392e443"),
		},
	}
}

func buildPrincipal() map[string]interface{} {
	principal := make(map[string]interface{}, 0)
	principal["key"] = uuid.New().String()
	principal["first_name"] = "first"
	principal["last_name"] = "last"
	principal["created_at"] = time.Now().String()
	principal["updated_at"] = time.Now().String()
	associated_tenants := make([]interface{}, 0)
	associated_tenant := make(map[string]interface{}, 0)
	associated_tenant["tenant"] = uuid.New().String()
	associated_tenants = append(associated_tenants, associated_tenant)
	principal["associated_tenants"] = associated_tenants
	return principal
}

func fetchBindingData() map[string]interface{} {
	file, err := os.Open("test_files/bindings.json")
	if err != nil {
		fmt.Println("error")
	}
	defer func() {
		err := file.Close()
		if err != nil {
			logger.LogError("Failed to close the file", "error", err)
		}
	}()

	// Read the file contents into a byte slice
	fileContents, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Declare a variable to hold the unmarshalled map
	var result map[string]interface{}

	// Unmarshal the JSON into the map
	err = json.Unmarshal(fileContents, &result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}
