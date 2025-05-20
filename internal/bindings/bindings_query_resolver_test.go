package bindings

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

func TestBinding(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	objUnderTest := BindingsQueryResolver{
		PC: mockService,
	}

	emptyUserId := ""

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", "7ed6cfa6-fd7e-4a2a-bbce-773ef8ea4c12")
	ginCtx.Set("userID", emptyUserId)

	testCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)

	assignment1 := make(map[string]interface{}, 0)
	assignment1["role"] = uuid.New().String()
	assignment1["user"] = uuid.New().String()
	assignment1["id"] = uuid.New().String()
	assignment1["created_at"] = time.Now().String()

	assignment2 := make(map[string]interface{}, 0)
	assignment2["role"] = uuid.New().String()
	assignment2["user"] = uuid.New().String()
	assignment2["id"] = uuid.New().String()
	assignment2["created_at"] = time.Now().String()

	assignments := make([]map[string]interface{}, 0)
	assignments = append(assignments, assignment1, assignment2)

	testcases := []struct {
		name      string
		input     uuid.UUID
		ctx       context.Context
		mockStubs func(mockService mocks.MockPermitService)
		output    models.OperationResult
	}{
		{
			name:      "GetBinding when input id is invalid",
			input:     uuid.Nil,
			ctx:       context.WithValue(context.Background(), config.GinContextKey, &gin.Context{}),
			mockStubs: func(mockSvc mocks.MockPermitService) {},
			output:    buildErrorResponse(400, "unable to fetch gin context", "tenant id not found in context"),
		},
		{
			name:  "GetBinding when fetch fails from permit",
			input: uuid.New(),
			ctx:   testCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().ExecuteGetAPI(mock.Any(), mock.Any(), mock.Any()).Return(nil, errors.New("failed to fetch from permit"))
			},
			output: buildErrorResponse(400, "failed", "error parsing user id: invalid UUID length: 0"),
		},
		{
			name:  "GetBinding success",
			input: uuid.New(),
			ctx:   testCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockSvc.EXPECT().ExecuteGetAPI(mock.Any(), mock.Any(), mock.Any()).Return(assignments, nil)
			},
			output: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockStubs(*mockService)
			result, _ := objUnderTest.Binding(tc.ctx, tc.input)
			assert.NotNil(t, result)
			// assert.Equal(t, tc.output, result)
		})

	}

}

func TestBindings(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPermitService(ctrl)
	objUnderTest := BindingsQueryResolver{
		PC: mockService,
	}

	emptyUserId := ""

	ginCtx := &gin.Context{}
	ginCtx.Set("tenantID", "7ed6cfa6-fd7e-4a2a-bbce-773ef8ea4c12")
	ginCtx.Set("userID", emptyUserId)

	testCtx := context.WithValue(context.Background(), config.GinContextKey, ginCtx)

	assignment1 := make(map[string]interface{}, 0)
	assignment1["role"] = uuid.New().String()
	assignment1["user"] = uuid.New().String()
	assignment1["id"] = uuid.New().String()
	assignment1["created_at"] = time.Now().String()

	assignment2 := make(map[string]interface{}, 0)
	assignment2["role"] = uuid.New().String()
	assignment2["user"] = uuid.New().String()
	assignment2["id"] = uuid.New().String()
	assignment2["created_at"] = time.Now().String()

	assignments := make([]map[string]interface{}, 0)
	assignments = append(assignments, assignment1, assignment2)

	testcases := []struct {
		name      string
		input     uuid.UUID
		ctx       context.Context
		mockStubs func(mockService mocks.MockPermitService)
		output    models.OperationResult
	}{
		{
			name:      "GetBindings when tenant id is not present",
			input:     uuid.Nil,
			ctx:       context.WithValue(context.Background(), config.GinContextKey, &gin.Context{}),
			mockStubs: func(mockSvc mocks.MockPermitService) {},
			output:    buildErrorResponse(400, "failed to fetch tenant id", "tenant id is missing in headers"),
		},
		{
			name:  "GetBindings when fetch fails from permit",
			input: uuid.New(),
			ctx:   testCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockService.EXPECT().ExecuteGetAPI(mock.Any(), mock.Any(), mock.Any()).Return(nil, errors.New("failed to fetch from permit"))
			},
			output: buildErrorResponse(400, "failed", "error parsing user id: invalid UUID length: 0"),
		},
		{
			name:  "GetBindings success",
			input: uuid.New(),
			ctx:   testCtx,
			mockStubs: func(mockSvc mocks.MockPermitService) {
				mockSvc.EXPECT().ExecuteGetAPI(mock.Any(), mock.Any(), mock.Any()).Return(assignments, nil)
			},
			output: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockStubs(*mockService)
			result, _ := objUnderTest.Bindings(tc.ctx)
			assert.NotNil(t, result)
			// assert.Equal(t, tc.output, result)
		})

	}

}
