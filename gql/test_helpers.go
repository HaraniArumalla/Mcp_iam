package gql

import (
	"iam_services_main_v1/mocks"
	"testing"

	"github.com/golang/mock/gomock"
)

type testSetup struct {
	ctrl       *gomock.Controller
	mockPermit *mocks.MockPermitService
	resolver   *Resolver
}

func setupTest(t *testing.T) *testSetup {
	ctrl := gomock.NewController(t)
	mockPermit := mocks.NewMockPermitService(ctrl)

	resolver := &Resolver{
		PC: mockPermit,
	}

	return &testSetup{
		ctrl:       ctrl,
		mockPermit: mockPermit,
		resolver:   resolver,
	}
}
