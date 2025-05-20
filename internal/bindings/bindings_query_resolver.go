package bindings

import (
	"context"
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/internal/constants"
	"iam_services_main_v1/internal/permit"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type BindingsQueryResolver struct {
	PC permit.PermitService
}

func (r *BindingsQueryResolver) Binding(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"className": "binding_query_resolver",
		"method":    "GetBinding",
		"bindingId": id,
	})

	logger.Infof("getbinding request received with Id %v", id)
	if id == uuid.Nil {
		logger.Error("invalid binding id provided")
		return buildErrorResponse(http.StatusBadRequest, "invalid id provided", "user id is invalid"), nil
	}

	url := fmt.Sprintf(constants.PERMIT_ROLE_ASSIGNMENTS+"?user=%s", id)
	bindingsFromPermit, err := r.PC.ExecuteGetAPI(ctx, constants.GET, url)

	if err != nil {
		logger.Error("unable to fetch bindings from permit ")
		return buildErrorResponse(http.StatusInternalServerError, "unable to fetch assignments from permit", err.Error()), nil
	}
	fmt.Println(bindingsFromPermit)

	assignments := make([]models.Data, 0)

	for _, binding := range bindingsFromPermit {
		assignment := models.Binding{
			Role:      &models.Role{ID: uuid.MustParse(binding["role"].(string))},
			ID:        uuid.MustParse(binding["id"].(string)),
			Principal: &models.User{ID: uuid.MustParse(binding["user"].(string))},
			CreatedAt: binding["created_at"].(string),
		}
		assignments = append(assignments, assignment)
	}
	fmt.Println(bindingsFromPermit)

	result := &models.SuccessResponse{
		Data:      assignments,
		IsSuccess: true,
		Message:   "Bindings retrieved successfully",
	}

	return result, nil
}

// AllBindings is the resolver for the allBindings field.
func (r *BindingsQueryResolver) Bindings(ctx context.Context) (models.OperationResult, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"className": "binding_query_resolver",
		"method":    "AllBindings",
	})

	tenantId, err := helpers.GetTenantID(ctx)
	if err != nil {
		logger.Error("unable to find tenant id from context")
		return buildErrorResponse(http.StatusBadRequest, "failed to fetch tenant id", "tenant id is missing in headers"), nil
	}

	// resourceInstance := fmt.Sprintf("?tenant="+"%s", tenantId)
	url := fmt.Sprintf(constants.PERMIT_ROLE_ASSIGNMENTS+"?tenant=%s", tenantId)
	res, err := r.PC.ExecuteGetAPI(ctx, constants.GET, url)

	if err != nil {
		logger.Error("unable to fetch resources from permit")
		return buildErrorResponse(http.StatusInternalServerError, "failed to fetch resources from permit", "failed to fetch bindings from permit"), nil
	}

	logger.Infof("fetch allBindings request received")
	var bindings []models.Data

	for _, binding := range res {
		createdAt := binding["created_at"].(string)
		updatedAt := binding["created_at"].(string)
		// principalType := r.FetchPrincipalBasedOnPrincipalId(ctx, uuid.MustParse(binding.UserId))

		bindingData := &models.Binding{
			ID:        uuid.MustParse(binding["id"].(string)),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Role:      &models.Role{ID: uuid.MustParse(binding["role"].(string))},
			Principal: &models.User{
				ID: uuid.MustParse(binding["user"].(string)),
			},
		}

		bindings = append(bindings, bindingData)
	}

	result := &models.SuccessResponse{
		Data:      bindings,
		IsSuccess: true,
		Message:   "Bindings retrieved successfully",
	}

	return result, nil
}
