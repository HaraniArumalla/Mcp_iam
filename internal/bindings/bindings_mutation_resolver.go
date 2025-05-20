package bindings

import (
	"context"
	"fmt"
	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/internal/constants"
	"iam_services_main_v1/internal/permit"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

const (
	TENANT_RESOURCE_TYPE_ID = "ed113bda-bbda-11ef-87ea-c03c5946f955"
)

type BindingsMutationResolver struct {
	PC permit.PermitService
}

func (r *BindingsMutationResolver) CreateBinding(ctx context.Context, input models.CreateBindingInput) (models.OperationResult, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"class":       "bindings_mutation_resolver",
		"method":      "CreateBinding",
		"bindingName": input.Name,
	})
	logger.Info("create binding request received")

	tenantId, err := helpers.GetTenantID(ctx)
	if err != nil {
		logger.Info("unable to find tenant id in context")
		return buildErrorResponse(http.StatusBadRequest, "unable to find tenant id in context", err.Error()), nil
	}

	userId, err := helpers.GetUserID(ctx)
	if err != nil {
		logger.Info("unable to find user id in context")
		return buildErrorResponse(http.StatusBadRequest, "unable to find user id in context", err.Error()), nil
	}

	if input.PrincipalID == uuid.Nil {
		logger.Info("unable to find principal id in request")
		return buildErrorResponse(http.StatusBadRequest, "unable to find principal id in input", "principal id is required"), nil
	}

	if input.RoleID == uuid.Nil {
		logger.Info("unable to find role id in request")
		return buildErrorResponse(http.StatusBadRequest, "unable to find role id in input", "role id is required"), nil
	}

	roleID := input.RoleID
	resourceId := uuid.New()
	currentDate := time.Now().String()

	// Create binding in Permit
	reqBody := map[string]interface{}{
		constants.ROLE:              input.RoleID.String(),
		constants.TENANT:            tenantId.String(),
		constants.USER:              input.PrincipalID.String(),
		constants.RESOURCE_INSTANCE: input.ScopeRefID.String() + ":" + tenantId.String(),
	}
	_, err = r.PC.APIExecute(ctx, constants.POST, constants.PERMIT_ROLE_ASSIGNMENTS, reqBody)

	if err != nil {
		return buildErrorResponse(http.StatusBadRequest, err.Error(), "unable to create binding in permit"), nil
	}

	logger.Info("Binding created successfully")

	data := &models.Binding{
		ID:        resourceId,
		Name:      input.Name,
		CreatedAt: currentDate,
		UpdatedAt: currentDate,
		Principal: &models.User{
			ID: input.PrincipalID,
		},
		Role:      &models.Role{ID: roleID},
		Version:   "V1",
		CreatedBy: *userId,
		UpdatedBy: *userId,
	}

	dataArray := make([]models.Data, 0)
	dataArray = append(dataArray, data)

	operationResult := models.SuccessResponse{
		Data:      dataArray,
		IsSuccess: true,
		Message:   "Binding created successfully",
	}
	return operationResult, nil

}

func (r *BindingsMutationResolver) UpdateBinding(ctx context.Context, input models.UpdateBindingInput) (models.OperationResult, error) {
	return nil, nil
	// logger := log.WithContext(ctx).WithFields(log.Fields{
	// 	"class":     "bindings_mutation_resolver",
	// 	"method":    "UpdateBinding",
	// 	"bindingID": input.ID,
	// })
	// logger.Info("update binding request received")
	//return nil, nil
	// logger := log.WithContext(ctx).WithFields(log.Fields{
	// 	"class":     "bindings_mutation_resolver",
	// 	"method":    "UpdateBinding",
	// 	"bindingID": input.ID,
	// })
	// logger.Info("update binding request received")

	// if input.ID == uuid.Nil {
	// 	return buildErrorResponse(http.StatusBadRequest, "invalid id provided", "invalid binding id provided"), nil
	// }
	// if input.Name == "" {
	// 	return buildErrorResponse(http.StatusBadRequest, "name is invalid", "name should not be empty"), nil
	// }
	// tenantId, err := helpers.GetTenantID(ctx)
	// if err != nil {
	// 	return buildErrorResponse(http.StatusBadRequest, "tenant id is missing", "tenant id is missing"), nil
	// }
	// if input.ID == uuid.Nil {
	// 	return buildErrorResponse(http.StatusBadRequest, "invalid id provided", "invalid binding id provided"), nil
	// }
	// if input.Name == "" {
	// 	return buildErrorResponse(http.StatusBadRequest, "name is invalid", "name should not be empty"), nil
	// }
	// tenantId, err := helpers.GetTenantID(ctx)
	// if err != nil {
	// 	return buildErrorResponse(http.StatusBadRequest, "tenant id is missing", "tenant id is missing"), nil
	// }

	// userId, err := helpers.GetUserID(ctx)
	// if err != nil {
	// 	return buildErrorResponse(http.StatusBadRequest, "user id is missing", "user id is missing"), nil
	// }
	// userId, err := helpers.GetUserID(ctx)
	// if err != nil {
	// 	return buildErrorResponse(http.StatusBadRequest, "user id is missing", "user id is missing"), nil
	// }

	// _, err = r.PC.APIExecute(ctx, constants.POST, constants.PERMIT_ROLE_ASSIGNMENTS, map[string]interface{}{
	// 	constants.USER:              input.PrincipalID,
	// 	constants.ROLE:              input.RoleID,
	// 	constants.TENANT:            tenantId,
	// 	constants.RESOURCE_INSTANCE: TENANT_RESOURCE_TYPE_ID + ":" + tenantId.String(),
	// })
	// _, err = r.PC.APIExecute(ctx, constants.POST, constants.PERMIT_ROLE_ASSIGNMENTS, map[string]interface{}{
	// 	constants.USER:              input.PrincipalID,
	// 	constants.ROLE:              input.RoleID,
	// 	constants.TENANT:            tenantId,
	// 	constants.RESOURCE_INSTANCE: TENANT_RESOURCE_TYPE_ID + ":" + tenantId.String(),
	// })

	// if err != nil {
	// 	return buildErrorResponse(http.StatusBadRequest, err.Error(), "unable to create binding in permit"), nil
	// }
	// if err != nil {
	// 	return buildErrorResponse(http.StatusBadRequest, err.Error(), "unable to create binding in permit"), nil
	// }

	// createdAt := time.Now().String()
	// assignment := &models.Binding{
	// 	ID:        input.ID,
	// 	Name:      input.Name,
	// 	CreatedAt: createdAt,
	// 	UpdatedAt: createdAt,
	// 	Role:      &models.Role{ID: input.RoleID},
	// 	Version:   "V1",
	// 	UpdatedBy: *userId,
	// }
	// createdAt := time.Now().String()
	// assignment := &models.Binding{
	// 	ID:        input.ID,
	// 	Name:      input.Name,
	// 	CreatedAt: createdAt,
	// 	UpdatedAt: createdAt,
	// 	Role:      &models.Role{ID: input.RoleID},
	// 	Version:   "V1",
	// 	UpdatedBy: *userId,
	// }

	// if err != nil {
	// 	return buildErrorResponse(http.StatusBadRequest, err.Error(), "cannot delete binding"), nil
	// }
	// if err != nil {
	// 	return buildErrorResponse(http.StatusBadRequest, err.Error(), "cannot delete binding"), nil
	// }

	// dataArray := make([]models.Data, 0)
	// dataArray = append(dataArray, assignment)
	// dataArray := make([]models.Data, 0)
	// dataArray = append(dataArray, assignment)

	// operationResult := models.SuccessResponse{
	// 	Data:      dataArray,
	// 	IsSuccess: true,
	// 	Message:   "Binding updated successfully",
	// }
	// return operationResult, nil
	// operationResult := models.SuccessResponse{
	// 	Data:      dataArray,
	// 	IsSuccess: true,
	// 	Message:   "Binding updated successfully",
	// }
	// return operationResult, nil
}

// DeleteBinding is the resolver for the deleteBinding field.
func (r *BindingsMutationResolver) DeleteBinding(ctx context.Context, input models.DeleteBindingInput) (models.OperationResult, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"class":     "bindings_mutation_resolver",
		"method":    "DeleteBinding",
		"bindingId": input.ID,
	})

	logger.Info("delete binding request received")
	if input.ID == uuid.Nil {
		logger.Error("invalid id provided. Please provide valid binding id")
		return buildErrorResponse(http.StatusBadRequest, "id is not present in request", "id is mandatory"), nil
	}
	tenantId, err := helpers.GetTenantID(ctx)
	if err != nil {
		return buildErrorResponse(http.StatusBadRequest, "unable to find tenantId in context", err.Error()), nil
	}

	if input.PrincipalID == uuid.Nil {
		logger.Info("unable to find principal id in request")
		return buildErrorResponse(http.StatusBadRequest, "unable to find principal id in input", "principal id is required"), nil
	}

	if input.RoleID == uuid.Nil {
		logger.Info("unable to find role id in request")
		return buildErrorResponse(http.StatusBadRequest, "unable to find role id in input", "role id is required"), nil
	}

	deleteRequestBody := map[string]interface{}{
		constants.ROLE:              input.RoleID,
		constants.TENANT:            tenantId,
		constants.USER:              input.PrincipalID,
		constants.RESOURCE_INSTANCE: input.ScopeRefID.String() + ":" + tenantId.String(),
	}
	_, err = r.PC.APIExecute(ctx, constants.DELETE, constants.PERMIT_ROLE_ASSIGNMENTS, deleteRequestBody)

	if err != nil {
		return buildErrorResponse(http.StatusNotFound, err.Error(), "unable to delete binding in permit"), nil
	}

	logger.Info("binding deleted successfully")
	result := &models.SuccessResponse{
		Data:      make([]models.Data, 0),
		IsSuccess: true,
		Message:   "Binding deleted successfully",
	}
	return result, nil
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
