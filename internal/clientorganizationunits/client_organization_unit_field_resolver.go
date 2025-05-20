package clientorganizationunits

import (
	"context"
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/tenants"
	"iam_services_main_v1/internal/users"
	"iam_services_main_v1/pkg/logger"
	"log"
)

// AccountFieldResolver provides database operations for resolving Account fields
type ClientOrganizationUnitResolver struct {
	PC permit.PermitService
}

// ParentOrg resolves the ParentOrg field on the Account type
func (r *ClientOrganizationUnitResolver) ParentOrg(ctx context.Context, corg *models.ClientOrganizationUnit) (models.Organization, error) {
	url := fmt.Sprintf("resource_instances/%s", corg.ParentOrg.GetID())
	resourceResponse, err := r.PC.GetSingleResource(ctx, "GET", url)
	if err != nil {
		logger.LogError("error fetching parent org", "error", err)
		return nil, err
	}
	log.Printf("ParentOrg: %v", resourceResponse)
	clientOrg := BuildOrgUnit(resourceResponse)
	return clientOrg, nil
}

// ParentOrg resolves the ParentOrg field on the Account type
func (r *ClientOrganizationUnitResolver) Tenant(ctx context.Context, corg *models.ClientOrganizationUnit) (*models.Tenant, error) {
	tenantResolver := &tenants.TenantQueryResolver{PC: r.PC}
	resourceResponse, err := tenantResolver.FetchTenant(ctx, corg.Tenant.ID)
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

func (r *ClientOrganizationUnitResolver) Accounts(ctx context.Context, corg *models.ClientOrganizationUnit) ([]*models.Account, error) {
	// accountResolver := accounts.AccountQueryResolver{PC: r.PC}
	// accounts := make([]*models.Account, 0)
	// for _, accountId := range corg.Accounts {
	// 	accountID := accountId.GetID()
	// 	account, err := accountResolver.Account(ctx, accountID)
	// 	if err != nil {
	// 		logger.LogError("error fetching account", "error", err)
	// 		return nil, err
	// 	}
	// 	if account == nil {
	// 		continue
	// 	}

	// }
	// return accounts, nil
	return nil, nil
}

func (r *ClientOrganizationUnitResolver) AccountOwner(ctx context.Context, corg *models.ClientOrganizationUnit) (*models.User, error) {
	userResolver := &users.UserResolver{PC: r.PC}
	return userResolver.GetUser(ctx, corg.AccountOwner.GetID())
}
