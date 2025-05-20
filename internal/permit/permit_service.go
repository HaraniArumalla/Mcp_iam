package permit

import "context"

type PermitService interface {
	APIExecute(ctx context.Context, method, endpoint string, payload interface{}) (interface{}, error)

	ExecuteGetAPI(ctx context.Context, method, endpoint string) ([]map[string]interface{}, error)

	GetSingleResource(ctx context.Context, method, endpoint string) (map[string]interface{}, error)

	SendRequest(ctx context.Context, method, endpoint string, payload interface{}) (map[string]interface{}, error)
}
