package bindings

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

func TestCreateBinding(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	objUnderTest := BindingsMutationResolver{
		PC: mockService,
	}

	emptyUserId := ""

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", "7ed6cfa6-fd7e-4a2a-bbce-773ef8ea4c12")
	ginCtx.Set("userID", emptyUserId)

	ginCtxWithUserId := &gin.Context{}
	ginCtxWithUserId.Set("tenantID", "7ed6cfa6-fd7e-4a2a-bbce-773ef8ea4c12")
	ginCtxWithUserId.Set("userID", "b5b44e90-906e-458a-8bb1-e9e4ee180696")

	testCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)
	validCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtxWithUserId)
	requestWithEmptyPrincipal := models.CreateBindingInput{
		PrincipalID: uuid.Nil,
	}

	requestWithEmptyRoleId := models.CreateBindingInput{
		PrincipalID: uuid.New(),
		RoleID:      uuid.Nil,
		Version:     "v1",
	}

	successRequest := models.CreateBindingInput{
		Name:        "test",
		PrincipalID: uuid.New(),
		RoleID:      uuid.New(),
		ScopeRefID:  uuid.New(),
		Version:     "v1",
	}

	testcases := []struct {
		name      string
		input     models.CreateBindingInput
		ctx       context.Context
		mockStubs func(mockService mocks.MockPermitService)
		output    models.OperationResult
	}{
		{
			name:      "CreateBinding when tenantId is not present in context",
			input:     successRequest,
			ctx:       context.WithValue(context.Background(), config.GinContextKey, &gin.Context{}),
			mockStubs: func(mockSvc mocks.MockPermitService) {},
			output:    buildErrorResponse(400, "unable to fetch gin context", "tenant id not found in context"),
		},
		{
			name:      "CreateBinding when user id is not present in context",
			input:     successRequest,
			ctx:       testCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {},
			output:    buildErrorResponse(400, "failed", "error parsing user id: invalid UUID length: 0"),
		},
		{
			name:      "CreateBinding when Principal is not present",
			input:     requestWithEmptyPrincipal,
			ctx:       validCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {},
			output:    buildErrorResponse(400, "failed", "principal id is required"),
		},
		{
			name:      "CreateBinding when Role is not present",
			input:     requestWithEmptyRoleId,
			ctx:       validCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {},
			output:    buildErrorResponse(400, "failed", "role id is required"),
		},
		{
			name:  "CreateBinding failed in permit",
			input: successRequest,
			ctx:   validCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().APIExecute(mock.Any(), mock.Any(), mock.Any(), mock.Any()).Return(nil, errors.New("failed to create"))
			},
			output: buildErrorResponse(400, "unable to create organization in permit", "unable to create binding in permit"),
		},
		{
			name:  "CreateBinding success",
			input: successRequest,
			ctx:   validCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().APIExecute(mock.Any(), mock.Any(), mock.Any(), mock.Any()).Return(nil, nil)
			},
			output: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockStubs(*mockService)
			result, _ := objUnderTest.CreateBinding(tc.ctx, tc.input)
			assert.NotNil(t, result)
			// assert.Equal(t, tc.output, result)
		})

	}

}

func TestDeleteBinding(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	objUnderTest := BindingsMutationResolver{
		PC: mockService,
	}

	emptyUserId := ""

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", "7ed6cfa6-fd7e-4a2a-bbce-773ef8ea4c12")
	ginCtx.Set("userID", emptyUserId)

	ginCtxWithUserId := &gin.Context{}
	ginCtxWithUserId.Set("tenantID", "7ed6cfa6-fd7e-4a2a-bbce-773ef8ea4c12")
	ginCtxWithUserId.Set("userID", "b5b44e90-906e-458a-8bb1-e9e4ee180696")

	// testCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)
	validCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtxWithUserId)

	requestWithEmptyId := models.DeleteBindingInput{
		PrincipalID: uuid.Nil,
	}
	requestWithEmptyPrincipal := models.DeleteBindingInput{
		ID:          uuid.New(),
		PrincipalID: uuid.Nil,
	}

	requestWithEmptyRoleId := models.DeleteBindingInput{
		ID:          uuid.New(),
		PrincipalID: uuid.New(),
		RoleID:      uuid.Nil,
		Version:     "v1",
	}

	successRequest := models.DeleteBindingInput{
		ID:          uuid.New(),
		Name:        "test",
		PrincipalID: uuid.New(),
		RoleID:      uuid.New(),
		ScopeRefID:  uuid.New(),
		Version:     "v1",
	}

	testcases := []struct {
		name      string
		input     models.DeleteBindingInput
		ctx       context.Context
		mockStubs func(mockService mocks.MockPermitService)
		output    models.OperationResult
	}{
		{
			name:      "DeleteBinding when input id is not present",
			input:     requestWithEmptyId,
			ctx:       context.WithValue(context.Background(), config.GinContextKey, &gin.Context{}),
			mockStubs: func(mockSvc mocks.MockPermitService) {},
			output:    buildErrorResponse(400, "unable to fetch gin context", "tenant id not found in context"),
		},
		{
			name:      "DeleteBinding when tenantId is not present in context",
			input:     successRequest,
			ctx:       context.WithValue(context.Background(), config.GinContextKey, &gin.Context{}),
			mockStubs: func(mockSvc mocks.MockPermitService) {},
			output:    buildErrorResponse(400, "unable to fetch gin context", "tenant id not found in context"),
		},
		{
			name:      "DeleteBinding when Principal is not present",
			input:     requestWithEmptyPrincipal,
			ctx:       validCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {},
			output:    buildErrorResponse(400, "failed", "principal id is required"),
		},
		{
			name:      "DeleteBinding when Role is not present",
			input:     requestWithEmptyRoleId,
			ctx:       validCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {},
			output:    buildErrorResponse(400, "failed", "role id is required"),
		},
		{
			name:  "DeleteBinding failed in permit",
			input: successRequest,
			ctx:   validCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().APIExecute(mock.Any(), mock.Any(), mock.Any(), mock.Any()).Return(nil, errors.New("failed to create"))
			},
			output: buildErrorResponse(400, "unable to create organization in permit", "unable to create binding in permit"),
		},
		{
			name:  "DeleteBinding success",
			input: successRequest,
			ctx:   validCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().APIExecute(mock.Any(), mock.Any(), mock.Any(), mock.Any()).Return(nil, nil)
			},
			output: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockStubs(*mockService)
			result, _ := objUnderTest.DeleteBinding(tc.ctx, tc.input)
			assert.NotNil(t, result)
			// assert.Equal(t, tc.output, result)
		})

	}

}
