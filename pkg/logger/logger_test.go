package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	// Test multiple initializations
	for i := 0; i < 3; i++ {
		InitLogger()
		assert.NotNil(t, log)
		assert.IsType(t, &logrus.Logger{}, log)
		assert.IsType(t, &logrus.JSONFormatter{}, log.Formatter)
	}
}

func TestLoggingFunctions(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer
	InitLogger()
	log.SetOutput(&buf)

	tests := []struct {
		name     string
		logFunc  func(string, ...interface{})
		level    string
		message  string
		fields   []interface{}
		validate func(t *testing.T, output map[string]interface{})
	}{
		{
			name:    "Info logging",
			logFunc: LogInfo,
			level:   "info",
			message: "test info",
			fields:  []interface{}{"key", "value"},
			validate: func(t *testing.T, output map[string]interface{}) {
				assert.Equal(t, "test info", output["msg"])
				assert.Equal(t, "value", output["key"])
			},
		},
		{
			name:    "Error logging",
			logFunc: LogError,
			level:   "error",
			message: "test error",
			fields:  []interface{}{"error", errors.New("test error")},
			validate: func(t *testing.T, output map[string]interface{}) {
				assert.Equal(t, "test error", output["msg"])
				assert.Equal(t, "test error", output["error"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc(tt.message, tt.fields...)

			var output map[string]interface{}
			err := json.Unmarshal(buf.Bytes(), &output)
			assert.NoError(t, err)

			assert.Equal(t, tt.level, output["level"])
			tt.validate(t, output)
		})
	}
}

func TestParseFields(t *testing.T) {
	tests := []struct {
		name          string
		fields        []interface{}
		expectedCount int
	}{
		{
			name:          "Empty fields",
			fields:        []interface{}{},
			expectedCount: 0,
		},
		{
			name:          "Valid fields",
			fields:        []interface{}{"key1", "value1", "key2", 123},
			expectedCount: 2,
		},
		{
			name:          "Error field",
			fields:        []interface{}{"error", errors.New("test error")},
			expectedCount: 1,
		},
		{
			name:          "Invalid key type",
			fields:        []interface{}{123, "value"},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := ParseFields(tt.fields...)
			assert.Equal(t, tt.expectedCount, len(fields))
		})
	}
}
