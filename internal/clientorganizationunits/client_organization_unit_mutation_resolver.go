package clientorganizationunits

import (
	"context"
	"errors"
	"fmt"
	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
	constants "iam_services_main_v1/internal/constants"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/utils"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

type ClientOrganizationUnitMutationResolver struct {
	PC permit.PermitService
}

func (r *ClientOrganizationUnitMutationResolver) CreateClientOrganizationUnit(ctx context.Context, input models.CreateClientOrganizationUnitInput) (models.OperationResult, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"className":  "organization_mutation_resolver",
		"methodName": "CreateClientOrganizationUnit",
	})
	logger.Info("create clientOrganizationUnit")

	if input.ID == uuid.Nil {
		return buildErrorResponse(http.StatusBadRequest, "ID is mandatory", "ID is mandatory"), nil
	}
	if input.Name == "" {
		return buildErrorResponse(http.StatusBadRequest, "Name is mandatory", "Name is mandatory"), nil
	}

	tenantId, err := helpers.GetTenantID(ctx)
	if err != nil {
		return buildErrorResponse(http.StatusBadRequest, "Tenant is mandatory", "Tenant id is not present"), nil
	}
	userId, err := helpers.GetUserID(ctx)
	if err != nil {
		return buildErrorResponse(http.StatusBadRequest, "user id is not present", "user id not present in header"), nil
	}

	if input.AccountOwnerID == uuid.Nil {
		return buildErrorResponse(http.StatusBadRequest, "Account owner id is mandatory", "Account owner id is mandatory"), nil
	}

	resourceTypeId := uuid.MustParse(constants.CORG_RESOURCE_TYPE_ID)
	resourceId := input.ID
	currentDate := time.Now()

	attributes := make(map[string]interface{})
	attributes[constants.KEY] = resourceId
	attributes[constants.NAME] = input.Name
	attributes[constants.DESCRIPTION] = input.Description
	attributes[constants.PARENT_RESOURCE_ID] = input.ParentID
	attributes[constants.RESOURCE_TYPE_ID] = resourceTypeId
	attributes[constants.TENANT_ID] = tenantId
	attributes[constants.ROW_STATUS] = 1
	attributes[constants.CREATED_AT] = currentDate
	attributes[constants.UPDATED_AT] = currentDate
	attributes[constants.CREATED_BY] = *userId
	attributes[constants.UPDATED_BY] = *userId
	attributes[constants.CORG_ACCOUNT_OWNER_ID] = input.AccountOwnerID
	attributes[constants.CORG_RELATION_TYPE] = input.RelationType
	attributes[constants.CORG_STATUS] = "ACTIVE"
	attributes[constants.TYPE] = config.ClientOrganizationUnit
	attributes[constants.CORG_TAGS] = input.Tags

	_, err = r.PC.APIExecute(ctx, constants.POST, constants.PERMIT_RESOURCE_INSTANCES, map[string]interface{}{
		constants.KEY:        resourceId,
		constants.TENANT:     tenantId,
		constants.RESOURCE:   resourceTypeId,
		constants.ATTRIBUTES: attributes,
	})

	if err != nil {
		return buildErrorResponse(http.StatusBadRequest, err.Error(), "unable to create organization in permit"), nil
	}

	return r.FormatSuccessResponse(ctx, resourceId)
}

func (r *ClientOrganizationUnitMutationResolver) UpdateClientOrganizationUnit(ctx context.Context, input models.UpdateClientOrganizationUnitInput) (models.OperationResult, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"className":  "client_organization_mutation_resolver",
		"methodName": "UpdateClientOrganizationUnit",
	})
	var zeroUUID uuid.UUID
	if input.ID == zeroUUID {
		return buildErrorResponse(http.StatusBadRequest, "invalid id provided in the request", "unable to update organization in permit"), errors.New("id is mandatory for update")
	}

	tenantId, err := helpers.GetTenantID(ctx)
	if err != nil {
		logger.Error("unable to find tenant id")
		return buildErrorResponse(http.StatusBadRequest, "tenant id not present", "tenant id not present"), nil
	}

	userId, err := helpers.GetUserID(ctx)
	if err != nil {
		logger.Error("unable to find user id")
		return buildErrorResponse(http.StatusBadRequest, "user id not present", "user id not present"), nil
	}

	if input.ParentID != nil && input.RelationType == nil {
		return buildErrorResponse(http.StatusBadRequest, "relation should not be empty", "relation type should not be empty"), nil
	}

	if !input.RelationType.IsValid() {
		return buildErrorResponse(http.StatusBadRequest, "relation type is invalid", "relation type is invalid"), nil
	}

	url := fmt.Sprintf(constants.PERMIT_RESOURCE_INSTANCES+"/%s", input.ID)
	clientOrg, err := r.PC.GetSingleResource(ctx, "GET", url)
	if err != nil {
		return buildErrorResponse(http.StatusBadRequest, err.Error(), "unable to update organization in permit"), nil
	}

	attributes := clientOrg[constants.ATTRIBUTES].(map[string]interface{})

	updatedAt := time.Now().String()
	if attributes != nil {
		attributes[constants.NAME] = input.Name
		attributes[constants.DESCRIPTION] = input.Description
		attributes[constants.PARENT_RESOURCE_ID] = input.ParentID
		attributes[constants.UPDATED_AT] = updatedAt
		attributes[constants.UPDATED_BY] = *userId
		attributes[constants.TENANT_ID] = tenantId
		attributes[constants.CORG_ACCOUNT_OWNER_ID] = input.AccountOwnerID
		attributes[constants.CORG_RELATION_TYPE] = input.RelationType
		attributes[constants.CORG_STATUS] = input.Status
		attributes[constants.CORG_TAGS] = input.Tags

		logger.Info("Client org in permit is ", clientOrg)
		reqBody := map[string]interface{}{
			constants.ATTRIBUTES: attributes,
		}
		_, err = r.PC.APIExecute(ctx, constants.PATCH, url, reqBody)

		if err != nil {
			return buildErrorResponse(http.StatusBadRequest, err.Error(), "unable to update organization in permit"), nil
		}
	}

	return r.FormatSuccessResponse(ctx, input.ID)
}

func (r *ClientOrganizationUnitMutationResolver) DeleteClientOrganizationUnit(ctx context.Context, input models.DeleteInput) (models.OperationResult, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"className":  "client_organization_mutation_resolver",
		"methodName": "DeleteClientOrganizationUnit",
	})

	var zeroUUID uuid.UUID
	if input.ID == zeroUUID {
		return buildErrorResponse(http.StatusBadRequest, "invalid input id", "id is mandatory"), errors.New("id is mandatory for delete")
	}

	logger.Info("delete client organization request received")

	resourceURL := fmt.Sprintf("resource_instances/%s", input.ID)
	clientOrgResource, err := r.PC.SendRequest(ctx, "GET", resourceURL, nil)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to get existing client org data", err.Error()), nil
	}
	if clientOrgResource == nil {
		return utils.FormatErrorResponse(http.StatusNotFound, "Client org not found", "The client org with the provided ID does not exist"), nil
	}
	resourceType := helpers.GetString(clientOrgResource, "resource")
	if resourceType != config.ClientOrgUnitResourceTypeID {
		return utils.FormatErrorResponse(http.StatusBadRequest, "The provided client org ID does not match the expected resource type for deletion", "The provided client org ID does not match the expected resource type for deletion"), nil
	}

	tenantID, err := helpers.GetTenantID(ctx)
	if err != nil {
		return utils.FormatErrorResponse(http.StatusBadRequest, "Failed to get tenant ID", err.Error()), nil
	}
	subjectType := config.ClientOrgUnitResourceTypeID + ":" + input.ID.String()
	url := fmt.Sprintf("relationship_tuples/detailed?tenant=%s&subject=%s",
		tenantID.String(),
		subjectType,
	)
	clientOrgUnitResponse, err := r.PC.SendRequest(ctx, constants.GET, url, nil)
	if err != nil {
		return buildErrorResponse(http.StatusNotFound, err.Error(), "unable to fetch resource from permit"), nil
	}
	rawData, ok := clientOrgUnitResponse["data"].([]interface{})
	if !ok {
		return buildErrorResponse(http.StatusNotFound, "missing or invalid data field", "missing or invalid data field"), nil
	}
	if len(rawData) > 0 {
		return buildErrorResponse(http.StatusBadRequest, "unable to delete organization due to associated accounts", "organization has accounts associated with it"), nil
	}

	url = fmt.Sprintf(constants.PERMIT_RESOURCE_INSTANCES+"/%s", input.ID)

	_, err = r.PC.APIExecute(ctx, constants.DELETE, url, nil)
	if err != nil {
		// do we need to rever in our database too?
		return buildErrorResponse(http.StatusBadRequest, err.Error(), "unable to delete organization in permit"), nil
	}

	result := models.SuccessResponse{
		Data:      make([]models.Data, 0),
		IsSuccess: true,
		Message:   "client organization Deleted successfully",
	}
	return result, nil
}

func (r *ClientOrganizationUnitMutationResolver) FormatSuccessResponse(ctx context.Context, resourceId uuid.UUID) (models.OperationResult, error) {
	url := fmt.Sprintf(constants.PERMIT_RESOURCE_INSTANCES+"/%s", resourceId)
	res, err := r.PC.GetSingleResource(ctx, constants.GET, url)
	if err != nil {
		return buildErrorResponse(http.StatusInternalServerError, "unable to fetch corg details", err.Error()), nil
	}
	return buildSuccessResponse(BuildOrgUnit(res)), nil

}
