package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvOrDefault(t *testing.T) {
	testCases := []struct {
		name         string
		key          string
		defaultVal   string
		envValue     string
		shouldSetEnv bool
		want         string
	}{
		{
			name:         "Environment variable exists",
			key:          "TEST_VAR",
			defaultVal:   "default",
			envValue:     "custom_value",
			shouldSetEnv: true,
			want:         "custom_value",
		},
		{
			name:         "Environment variable empty",
			key:          "TEST_VAR",
			defaultVal:   "default",
			envValue:     "",
			shouldSetEnv: true,
			want:         "default",
		},
		{
			name:         "Environment variable not set",
			key:          "TEST_VAR",
			defaultVal:   "default",
			shouldSetEnv: false,
			want:         "default",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear the environment variable
			if err := os.Unsetenv(tc.key); err != nil {
				t.Errorf("failed to unset environment variable %s: %v", tc.key, err)
			}

			// Set environment variable if required
			if tc.shouldSetEnv {
				if err := os.Setenv(tc.key, tc.envValue); err != nil {
					t.Errorf("failed to set environment variable %s: %v", tc.key, err)
				}
			}

			got := GetEnvOrDefault(tc.key, tc.defaultVal)
			assert.Equal(t, tc.want, got)
		})
	}
}

func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func TestLoadEnv(t *testing.T) {
	// Save original environment
	origEnv := map[string]string{
		"PERMIT_PDP_ENDPOINT": os.Getenv("PERMIT_PDP_ENDPOINT"),
		"PERMIT_PROJECT":      os.Getenv("PERMIT_PROJECT"),
		"PERMIT_ENV":          os.Getenv("PERMIT_ENV"),
		"PERMIT_TOKEN":        os.Getenv("PERMIT_TOKEN"),
	}

	// Restore environment after tests
	defer func() {
		for k, v := range origEnv {
			if v != "" {
				if err := os.Setenv(k, v); err != nil {
					t.Errorf("failed to set environment variable %s: %v", k, err)
				}
			} else {
				if err := os.Unsetenv(k); err != nil {
					t.Errorf("failed to unset environment variable %s: %v", k, err)
				}
			}
		}
	}()

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Missing all variables",
			setup: func() {
				os.Clearenv()
			},
			wantErr: true,
		},
		{
			name: "Empty values",
			setup: func() {
				if err := os.Setenv("PERMIT_PDP_ENDPOINT", ""); err != nil {
					t.Errorf("failed to set environment variable PERMIT_PDP_ENDPOINT: %v", err)
				}
				if err := os.Setenv("PERMIT_PROJECT", ""); err != nil {
					t.Errorf("failed to set environment variable PERMIT_PROJECT: %v", err)
				}
				if err := os.Setenv("PERMIT_ENV", ""); err != nil {
					t.Errorf("failed to set environment variable PERMIT_ENV: %v", err)
				}
				if err := os.Setenv("PERMIT_TOKEN", ""); err != nil {
					t.Errorf("failed to set environment variable PERMIT_TOKEN: %v", err)
				}
			},
			wantErr: true,
		},
		{
			name: "Missing one variable",
			setup: func() {
				if err := os.Setenv("PERMIT_PDP_ENDPOINT", "http://test.com"); err != nil {
					t.Errorf("failed to set environment variable PERMIT_PDP_ENDPOINT: %v", err)
				}
				if err := os.Setenv("PERMIT_PROJECT", "test-project"); err != nil {
					t.Errorf("failed to set environment variable PERMIT_PROJECT: %v", err)
				}
				if err := os.Setenv("PERMIT_ENV", "test-env"); err != nil {
					t.Errorf("failed to set environment variable PERMIT_ENV: %v", err)
				}
				if err := os.Unsetenv("PERMIT_TOKEN"); err != nil {
					t.Errorf("failed to unset environment variable PERMIT_TOKEN: %v", err)
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment before each test
			os.Clearenv()

			// Setup test environment
			tt.setup()

			// Run test
			err := LoadEnv()

			if tt.wantErr {
				assert.Error(t, err, "LoadEnv() should return error")
			} else {
				assert.NoError(t, err, "LoadEnv() should not return error")

				// Verify required environment variables
				assert.NotEmpty(t, os.Getenv("PERMIT_PDP_ENDPOINT"))
				assert.NotEmpty(t, os.Getenv("PERMIT_PROJECT"))
				assert.NotEmpty(t, os.Getenv("PERMIT_ENV"))
				assert.NotEmpty(t, os.Getenv("PERMIT_TOKEN"))

				// Verify config constants are set
				assert.NotEmpty(t, AccountResourceTypeID)
				assert.NotEmpty(t, GenericErrorMessage)
			}
		})
	}
}
