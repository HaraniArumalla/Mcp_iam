package permit_sdk

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockPermitSdkService is a mock implementation of PermitSdkService
type MockPermitSdkService struct {
	mock.Mock
}

// Check mocks the Check method of PermitSdkService
func (m *MockPermitSdkService) Check(ctx context.Context, userID, action, resourceType, resourceID, tenant string) (bool, error) {
	args := m.Called(ctx, userID, action, resourceType, resourceID, tenant)
	return args.Bool(0), args.Error(1)
}
