package clientorganizationunits

import (
	"context"
	"fmt"
	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
	constants "iam_services_main_v1/internal/constants"
	"iam_services_main_v1/internal/permit"
	tag_helper "iam_services_main_v1/internal/tags"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

type ClientOrganizationUnitQueryResolver struct {
	PC permit.PermitService
}

func (r *ClientOrganizationUnitQueryResolver) ClientOrganizationUnit(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"className":  "organization_query_resolver",
		"methodName": "AllOrganizations",
	})

	_, err := helpers.GetTenantID(ctx)
	if err != nil {
		logger.Error("unable to fetch tenantID")
		return buildErrorResponse(400, "unable to find tenantid", "tenantId is not present in header"), nil
	}

	url := fmt.Sprintf(constants.PERMIT_RESOURCE_INSTANCES+"/%s", id)
	res, err := r.PC.GetSingleResource(ctx, constants.GET, url)
	if err != nil {
		return buildErrorResponse(http.StatusNotFound, err.Error(), "unable to fetch resource from permit"), nil
	}

	return buildSuccessResponse(BuildOrgUnit(res)), nil
}

func (r *ClientOrganizationUnitQueryResolver) ClientOrganizationUnits(ctx context.Context) (models.OperationResult, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"className":  "organization_query_resolver",
		"methodName": "ClientOrganizationUnits",
	})
	tenantId, err := helpers.GetTenantID(ctx)
	if err != nil {
		return buildErrorResponse(http.StatusBadRequest, "unable to find tenant id", "unable to find tenant id in headers"), nil
	}

	endpoint := fmt.Sprintf(constants.PERMIT_RESOURCE_INSTANCES+"/detailed?tenant=%s&resource=%s", tenantId, constants.CORG_RESOURCE_TYPE_ID)
	permitResult, err := r.PC.SendRequest(ctx, constants.GET, endpoint, nil)
	if err != nil {
		logger.Error("unable to fetch client orgs from permit")
		message := "unable to fetch client organization units from permit"
		return buildErrorResponse(http.StatusNotFound, err.Error(), message), nil
	}

	clientOrgs := permitResult["data"].([]interface{})
	var orgs []models.Data
	for _, corg := range clientOrgs {
		clientOrgUnit := corg.(map[string]interface{})

		if clientOrgUnit != nil {
			unit := BuildOrgUnit(corg.(map[string]interface{}))
			orgs = append(orgs, unit)
		}
	}
	result := models.SuccessResponse{
		Data:      orgs,
		IsSuccess: true,
		Message:   "All Client Organizations retrieved successfully",
	}

	return result, nil
}

func BuildOrgUnit(result map[string]interface{}) *models.ClientOrganizationUnit {
	attributes := result[constants.ATTRIBUTES].(map[string]interface{})
	description := attributes[constants.DESCRIPTION].(string)
	tenantId := uuid.MustParse(attributes[constants.TENANT_ID].(string))
	key := uuid.MustParse(attributes[constants.KEY].(string))
	createdBy := uuid.MustParse(attributes[constants.CREATED_BY].(string))
	updatedBy := uuid.MustParse(attributes[constants.UPDATED_BY].(string))
	relationType := models.RelationTypeEnum(attributes[constants.CORG_RELATION_TYPE].(string))
	status := models.StatusTypeEnum(attributes[constants.CORG_STATUS].(string))
	accountOwner := &models.User{ID: uuid.MustParse(attributes[constants.CORG_ACCOUNT_OWNER_ID].(string))}

	// if len(tags) == 0 {
	// 	tagList = make([]*models.Tags, 0)
	// } else {
	// 	for _, v := range tags {
	// 		tag := &models.Tags{}
	// 		if tagMap, ok := v.(map[string]interface{}); ok {
	// 			tag.Key = helpers.GetString(tagMap, "key")
	// 			tag.Value = helpers.GetString(tagMap, "value")
	// 			tagList = append(tagList, tag)
	// 		}
	// 	}
	// }
	unit := &models.ClientOrganizationUnit{
		ID:           key,
		Name:         attributes[constants.NAME].(string),
		Description:  &description,
		Tenant:       &models.Tenant{ID: tenantId},
		CreatedAt:    attributes[constants.CREATED_AT].(string),
		UpdatedAt:    attributes[constants.UPDATED_AT].(string),
		CreatedBy:    createdBy,
		UpdatedBy:    updatedBy,
		RelationType: relationType,
		AccountOwner: accountOwner,
		Status:       status,
		Tags:         tag_helper.GetResourceTags(attributes, constants.CORG_TAGS),
	}
	return unit
}

func buildSuccessResponse(unit *models.ClientOrganizationUnit) models.OperationResult {
	clientOrgs := make([]models.Data, 0)
	clientOrgs = append(clientOrgs, unit)
	result := models.SuccessResponse{
		Data:      clientOrgs,
		IsSuccess: true,
		Message:   "Client Organization Unit retrieved successfully",
	}
	return result
}
func buildErrorResponse(statusCode int, errorDetail, message string) models.ResponseError {
	return models.ResponseError{
		ErrorCode:     fmt.Sprint(statusCode),
		ErrorDetails:  &errorDetail,
		IsSuccess:     false,
		Message:       config.GenericErrorMessage,
		SystemMessage: message,
	}
}
