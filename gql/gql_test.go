package gql

import (
	"iam_services_main_v1/gql/generated"
	"iam_services_main_v1/internal/accounts"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolver_Query(t *testing.T) {
	ts := setupTest(t)
	defer ts.ctrl.Finish()

	tests := []struct {
		name          string
		checkResolver func(*queryResolver) bool
		expectedError bool
	}{
		{
			name: "account query resolver initialization",
			checkResolver: func(qr *queryResolver) bool {
				return qr.AccountQueryResolver != nil &&
					qr.AccountQueryResolver.PC == ts.mockPermit
			},
		},
		{
			name: "tenant query resolver initialization",
			checkResolver: func(qr *queryResolver) bool {
				return qr.TenantQueryResolver != nil &&
					qr.TenantQueryResolver.PC == ts.mockPermit
			},
		},
		{
			name: "role query resolver initialization",
			checkResolver: func(qr *queryResolver) bool {
				return qr.RoleQueryResolver != nil &&
					qr.RoleQueryResolver.PC == ts.mockPermit
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryResolver := ts.resolver.Query().(*queryResolver)
			assert.True(t, tt.checkResolver(queryResolver))
		})
	}
}

func TestResolver_Mutation(t *testing.T) {
	ts := setupTest(t)
	defer ts.ctrl.Finish()

	tests := []struct {
		name          string
		checkResolver func(*mutationResolver) bool
		expectedError bool
	}{
		{
			name: "account mutation resolver initialization",
			checkResolver: func(mr *mutationResolver) bool {
				return mr.AccountMutationResolver != nil &&
					mr.AccountMutationResolver.PC == ts.mockPermit
			},
		},
		{
			name: "tenant mutation resolver initialization",
			checkResolver: func(mr *mutationResolver) bool {
				return mr.TenantMutationResolver != nil &&
					mr.TenantMutationResolver.PC == ts.mockPermit
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mutationResolver := ts.resolver.Mutation().(*mutationResolver)
			assert.True(t, tt.checkResolver(mutationResolver))
		})
	}
}

func TestResolver_Account(t *testing.T) {
	ts := setupTest(t)
	defer ts.ctrl.Finish()

	tests := []struct {
		name string
		want generated.AccountResolver
	}{
		{
			name: "account field resolver initialization",
			want: &accounts.AccountFieldResolver{PC: ts.mockPermit},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ts.resolver.Account()
			assert.IsType(t, tt.want, got)
			assert.Equal(t, tt.want.(*accounts.AccountFieldResolver).PC,
				got.(*accounts.AccountFieldResolver).PC)
		})
	}
}

func TestResolver_DependencyInjection(t *testing.T) {
	ts := setupTest(t)
	defer ts.ctrl.Finish()

	t.Run("verify permit service injection", func(t *testing.T) {
		// Query resolvers
		query := ts.resolver.Query()
		queryR := query.(*queryResolver)
		assert.Equal(t, ts.mockPermit, queryR.AccountQueryResolver.PC)
		assert.Equal(t, ts.mockPermit, queryR.TenantQueryResolver.PC)

		// Mutation resolvers
		mutation := ts.resolver.Mutation()
		mutationR := mutation.(*mutationResolver)
		assert.Equal(t, ts.mockPermit, mutationR.AccountMutationResolver.PC)
		assert.Equal(t, ts.mockPermit, mutationR.TenantMutationResolver.PC)

		// Field resolvers
		account := ts.resolver.Account()
		accountR := account.(*accounts.AccountFieldResolver)
		assert.Equal(t, ts.mockPermit, accountR.PC)
	})
}

func TestResolver_InterfaceImplementation(t *testing.T) {
	var _ generated.QueryResolver = (*queryResolver)(nil)
	var _ generated.MutationResolver = (*mutationResolver)(nil)
	var _ generated.AccountResolver = (*accounts.AccountFieldResolver)(nil)
}
