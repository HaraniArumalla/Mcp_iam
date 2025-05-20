package accounts

import (
	"context"
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/clientorganizationunits"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/tenants"
	"iam_services_main_v1/internal/users"
	"iam_services_main_v1/pkg/logger"
	"log"
)

// AccountFieldResolver provides database operations for resolving Account fields
type AccountFieldResolver struct {
	PC permit.PermitService
}

// ParentOrg resolves the ParentOrg field on the Account type
func (r *AccountFieldResolver) ParentOrg(ctx context.Context, account *models.Account) (models.Organization, error) {
	url := fmt.Sprintf("resource_instances/%s", account.ParentOrg.GetID())
	resourceResponse, err := r.PC.GetSingleResource(ctx, "GET", url)
	if err != nil {
		logger.LogError("error fetching parent org", "error", err)
		return nil, err
	}
	log.Printf("ParentOrg: %v", resourceResponse)
	clientOrg := clientorganizationunits.BuildOrgUnit(resourceResponse)
	return clientOrg, nil
}

// ParentOrg resolves the ParentOrg field on the Account type
func (r *AccountFieldResolver) Tenant(ctx context.Context, account *models.Account) (*models.Tenant, error) {
	tenantResolver := &tenants.TenantQueryResolver{PC: r.PC}
	resourceResponse, err := tenantResolver.FetchTenant(ctx, account.Tenant.GetID())
	if err != nil {
		return nil, err
	}
	tenant, err := tenants.MapTenantData(resourceResponse)
	if err != nil {
		logger.LogError("error mapping Tenant data in MapTenantResponseToStruct", "error", err)
		return nil, err
	}
	return tenant, nil
}

// AccountOwner resolves the AccountOwner field on the Account type
func (r *AccountFieldResolver) AccountOwner(ctx context.Context, account *models.Account) (*models.User, error) {
	userResolver := &users.UserResolver{PC: r.PC}
	return userResolver.GetUser(ctx, account.AccountOwner.GetID())
}
