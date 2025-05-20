package tenants

import (
	"context"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/users"
)

// needs to be implemented
type TenantFieldResolver struct {
	PC permit.PermitService
}

// ParentOrg resolves the ParentOrg field on the Account type
func (r *TenantFieldResolver) ClientOrganizationUnits(ctx context.Context, tenant *models.Tenant) ([]*models.ClientOrganizationUnit, error) {
	return nil, nil
}

// AccountOwner resolves the AccountOwner field on the Account type
func (r *TenantFieldResolver) AccountOwner(ctx context.Context, tenant *models.Tenant) (*models.User, error) {
	userResolver := &users.UserResolver{PC: r.PC}
	return userResolver.GetUser(ctx, tenant.AccountOwner.GetID())
}
