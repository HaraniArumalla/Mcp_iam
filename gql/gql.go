package gql

import (
	"iam_services_main_v1/gql/generated"
	"iam_services_main_v1/internal/accounts"
	"iam_services_main_v1/internal/bindings"
	"iam_services_main_v1/internal/clientorganizationunits"
	"iam_services_main_v1/internal/groups"
	"iam_services_main_v1/internal/organizations"
	"iam_services_main_v1/internal/permissions"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/resources"
	role "iam_services_main_v1/internal/roles"
	"iam_services_main_v1/internal/root"
	"iam_services_main_v1/internal/tenants"
)

// Resolver holds references to the DB and acts as a central resolver
type Resolver struct {
	PC  permit.PermitService
	PSC *permit.PermitSdkService
}

// Query returns the root query resolvers, delegating to feature-based resolvers
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{

		TenantQueryResolver:                 &tenants.TenantQueryResolver{PC: r.PC},
		AccountQueryResolver:                &accounts.AccountQueryResolver{PC: r.PC},
		ClientOrganizationUnitQueryResolver: &clientorganizationunits.ClientOrganizationUnitQueryResolver{PC: r.PC},
		RoleQueryResolver:                   &role.RoleQueryResolver{PC: r.PC},
		PermissionQueryResolver:             &permissions.PermissionQueryResolver{},
		BindingsQueryResolver:               &bindings.BindingsQueryResolver{PC: r.PC},
		ResourceQueryResolver:               &resources.ResourceQueryResolver{PSC: r.PSC},
		GroupQueryResolver:                  &groups.GroupQueryResolver{},
		OrganizationQueryResolver:           &organizations.OrganizationQueryResolver{},
		RootQueryResolver:                   &root.RootQueryResolver{},
	}
}

// Mutation returns the root mutation resolvers, delegating to feature-based resolvers
func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{

		AccountMutationResolver:                &accounts.AccountMutationResolver{PC: r.PC},
		TenantMutationResolver:                 &tenants.TenantMutationResolver{PC: r.PC, PSC: r.PSC},
		ClientOrganizationUnitMutationResolver: &clientorganizationunits.ClientOrganizationUnitMutationResolver{PC: r.PC},
		RoleMutationResolver:                   &role.RoleMutationResolver{PC: r.PC},
		PermissionMutationResolver:             &permissions.PermissionMutationResolver{},
		BindingsMutationResolver:               &bindings.BindingsMutationResolver{PC: r.PC},
		RootMutationResolver:                   &root.RootMutationResolver{},
	}
}

// Root resolvers for Query and Mutation
type queryResolver struct {
	*tenants.TenantQueryResolver
	*accounts.AccountQueryResolver
	*role.RoleQueryResolver
	*clientorganizationunits.ClientOrganizationUnitQueryResolver
	*permissions.PermissionQueryResolver
	*bindings.BindingsQueryResolver
	*resources.ResourceQueryResolver
	*groups.GroupQueryResolver
	*organizations.OrganizationQueryResolver
	*root.RootQueryResolver
}

type mutationResolver struct {
	*tenants.TenantMutationResolver
	*accounts.AccountMutationResolver
	*clientorganizationunits.ClientOrganizationUnitMutationResolver
	*role.RoleMutationResolver
	*permissions.PermissionMutationResolver
	*bindings.BindingsMutationResolver
	*root.RootMutationResolver
}

// Account resolves fields for the Account type
func (r *Resolver) Account() generated.AccountResolver {
	return &accounts.AccountFieldResolver{PC: r.PC}
}

// Account resolves fields for the Account type
func (r *Resolver) ClientOrganizationUnit() generated.ClientOrganizationUnitResolver {
	return &clientorganizationunits.ClientOrganizationUnitResolver{PC: r.PC}
}

// Account resolves fields for the Account type
func (r *Resolver) Binding() generated.BindingResolver {
	return &bindings.BindingsResolver{PC: r.PC}
}

// Account resolves fields for the Account type
func (r *Resolver) Tenant() generated.TenantResolver {
	return &tenants.TenantFieldResolver{PC: r.PC}
}
