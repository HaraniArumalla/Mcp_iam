package permit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"iam_services_main_v1/pkg/logger"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
)

// PermitClient is a client for interacting with the Permit API.
type PermitClient struct {
	BaseURL string
	Headers map[string]string
	Client  *http.Client
}

type PermitServiceImpl struct {
	PermitClient *PermitClient
}

// NewPermitClient creates a new PermitClient.
func NewPermitServiceImpl(permitClient *PermitClient) PermitService {
	return &PermitServiceImpl{
		PermitClient: permitClient,
	}
}

// SendRequest sends an HTTP request without retry logic.
func (pc *PermitServiceImpl) SendRequest(ctx context.Context, method, endpoint string, payload interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	// Serialize payload to JSON
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			logger.LogError("Failed to marshal payload", "error", err)
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	// Change base URL for specific endpoints (roles/resources)
	if strings.Contains(endpoint, "roles") || strings.Contains(endpoint, "resources") {
		pc.PermitClient.BaseURL = strings.Replace(pc.PermitClient.BaseURL, "facts", "schema", 1)
	} else {
		pc.PermitClient.BaseURL = strings.Replace(pc.PermitClient.BaseURL, "schema", "facts", 1)
	}
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", pc.PermitClient.BaseURL, endpoint), body)
	if err != nil {
		logger.LogError("Failed to create HTTP request", "error", err)
		return nil, err
	}

	// Log URL
	logger.LogInfo("Permit request URL", "url", req.URL.String())

	// Add headers
	for key, value := range pc.PermitClient.Headers {
		req.Header.Set(key, value)
	}

	// Send the request
	resp, err := pc.PermitClient.Client.Do(req)
	if err != nil {
		logger.LogError("HTTP request failed", "error", err)
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logger.LogError("Failed to close response body", "error", err)
		}
	}()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("response body ", string(body))
		logger.LogError("error occurred when calling Permit API", "status", resp.StatusCode, "body", string(body))
		return nil, fmt.Errorf("http error: %d", resp.StatusCode)
	}

	// Parse response body
	respBody, err := io.ReadAll(resp.Body)
	fmt.Println("response body is ", string(respBody))
	if err != nil {
		logger.LogError("Failed to read response body", "error", err)
		return nil, err
	}

	if len(respBody) == 0 {
		logger.LogInfo("Empty response body got from permit API")
		return nil, nil
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		logger.LogError("Failed to unmarshal response", "error", err)
		return nil, err
	}

	return result, nil
}

type PermitBinding struct {
	Id             string `json:"id"`
	User           string `json:"user"`
	Role           string `json:"role"`
	Tenant         string `json:"tenant"`
	UserId         string `json:"user_id"`
	RoleId         string `json:"role_id"`
	TenantId       string `json:"tenant_id"`
	OrganizationId string `json:"organization_id"`
	ProjectId      string `json:"project_id"`
	EnvironmentId  string `json:"environment_id"`
	CreatedAt      string `json:"created_at"`
}

func (pc *PermitServiceImpl) sendGetRequest(ctx context.Context, method, endpoint string) ([]map[string]interface{}, error) {
	var response []map[string]interface{}

	operation := func() error {
		req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", pc.PermitClient.BaseURL, endpoint), nil)
		if err != nil {
			log.Printf("Failed to create HTTP request: %v", err)
			return backoff.Permanent(err)
		}

		// Add headers
		for key, value := range pc.PermitClient.Headers {
			req.Header.Set(key, value)
		}

		// Send the request
		resp, err := pc.PermitClient.Client.Do(req)
		if err != nil {
			log.Printf("HTTP request failed: %v", err)
			return err
		}
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				logger.LogError("Failed to close response body", "error", err)
			}
		}()

		// Check response status
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("HTTP error: %d - %s", resp.StatusCode, string(body))
			return fmt.Errorf("HTTP error: %d", resp.StatusCode)
		}

		// Parse response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			return backoff.Permanent(err)
		}

		if method == "GET" && len(respBody) == 0 {
			log.Printf("Empty response body")
			return nil
		}

		if err := json.Unmarshal(respBody, &response); err != nil {
			log.Printf("Failed to unmarshal response: %v", err)
			return backoff.Permanent(err)
		}

		return nil
	}

	// Use exponential backoff for retries
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 30 * time.Second

	if err := backoff.Retry(operation, bo); err != nil {
		log.Printf("Request failed after retries: %v", err)
		return response, err
	}

	return response, nil
}

func (pc *PermitServiceImpl) fetchSingleResource(ctx context.Context, method, endpoint string) (map[string]interface{}, error) {
	var response map[string]interface{}

	operation := func() error {
		req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", pc.PermitClient.BaseURL, endpoint), nil)
		if err != nil {
			log.Printf("Failed to create HTTP request: %v", err)
			return backoff.Permanent(err)
		}

		// Add headers
		for key, value := range pc.PermitClient.Headers {
			req.Header.Set(key, value)
		}

		// Send the request
		resp, err := pc.PermitClient.Client.Do(req)
		if err != nil {
			log.Printf("HTTP request failed: %v", err)
			return err
		}
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				logger.LogError("Failed to close response body", "error", err)
			}
		}()

		// Check response status
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("HTTP error: %d - %s", resp.StatusCode, string(body))
			return fmt.Errorf("HTTP error: %d", resp.StatusCode)
		}

		// Parse response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			return backoff.Permanent(err)
		}

		if method == "GET" && len(respBody) == 0 {
			log.Printf("Empty response body")
			return nil
		}

		if err := json.Unmarshal(respBody, &response); err != nil {
			log.Printf("Failed to unmarshal response: %v", err)
			return backoff.Permanent(err)
		}

		return nil
	}

	// Use exponential backoff for retries
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 30 * time.Second

	if err := backoff.Retry(operation, bo); err != nil {
		log.Printf("Request failed after retries: %v", err)
		return response, err
	}

	return response, nil
}

func (pc *PermitServiceImpl) APIExecute(ctx context.Context, method, endpoint string, payload interface{}) (interface{}, error) {
	return pc.SendRequest(ctx, method, endpoint, payload)
}

func (pc *PermitServiceImpl) ExecuteGetAPI(ctx context.Context, method, endpoint string) ([]map[string]interface{}, error) {
	return pc.sendGetRequest(ctx, method, endpoint)
}

func (pc *PermitServiceImpl) GetSingleResource(ctx context.Context, method, endpoint string) (map[string]interface{}, error) {
	return pc.fetchSingleResource(ctx, method, endpoint)
}
