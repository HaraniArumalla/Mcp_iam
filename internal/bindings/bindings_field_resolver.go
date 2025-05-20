package bindings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/roles"
	"log"

	"github.com/google/uuid"
)

// AccountFieldResolver provides database operations for resolving Account fields
type BindingsResolver struct {
	PC permit.PermitService
}

// ParentOrg resolves the ParentOrg field on the Account type
func (r *BindingsResolver) Principal(ctx context.Context, obj *models.Binding) (models.Principal, error) {
	url := fmt.Sprintf("users/%s", obj.Principal.GetID())
	user, err := r.PC.GetSingleResource(ctx, "GET", url)
	if err != nil {
		return nil, err
	}

	log.Printf("user is: %v", user)
	tenantMap := user["associated_tenants"].([]interface{})

	fmt.Println(tenantMap)
	userDetails := models.User{
		ID:        uuid.MustParse(user["key"].(string)),
		FirstName: user["first_name"].(string),
		LastName:  user["last_name"].(string),
		Name:      user["first_name"].(string) + " " + user["last_name"].(string),
		UpdatedAt: user["updated_at"].(string),
		CreatedAt: user["created_at"].(string),
		Tenant: &models.Tenant{
			ID: uuid.MustParse(tenantMap[0].(map[string]interface{})["tenant"].(string)),
		},
	}
	// clientOrg := BuildOrgUnit(resourceResponse)
	return userDetails, nil
}

// ParentOrg resolves the ParentOrg field on the Account type
func (r *BindingsResolver) Role(ctx context.Context, obj *models.Binding) (*models.Role, error) {
	roleQuery := roles.RoleQueryResolver{PC: r.PC}
	result, err := roleQuery.Role(ctx, obj.Role.ID)
	if err != nil {
		return nil, errors.New("unable to fetch role details")
	}

	return fetchRole(ctx, result)
}
func fetchRole(ctx context.Context, roleQuery models.OperationResult) (*models.Role, error) {
	roleMap := make(map[string]interface{}, 0)
	jsonStr, err := json.Marshal(roleQuery)
	if err != nil {
		return nil, errors.New("error while fetching role")
	}

	err = json.Unmarshal(jsonStr, &roleMap)
	if err != nil {
		return nil, errors.New("errored out")
	}

	isSuccess := roleMap["isSuccess"].(bool)
	dataArray := roleMap["data"]
	if isSuccess && dataArray != nil {
		if data, ok := dataArray.([]interface{}); ok {
			roleMap = data[0].(map[string]interface{})
			description := roleMap["description"].(string)

			role := &models.Role{
				ID:          uuid.MustParse(roleMap["id"].(string)),
				Name:        roleMap["name"].(string),
				Description: &description,
				CreatedAt:   roleMap["createdAt"].(string),
				CreatedBy:   uuid.MustParse(roleMap["createdBy"].(string)),
			}
			roleType := roleMap["roleMap"]
			if roleType == "" {
				role.RoleType = models.RoleTypeEnumDefault
			} else {
				role.RoleType = models.RoleTypeEnumCustom
			}
			return role, nil
		}

	}
	return nil, errors.New("unable to find role information in permit")
}
