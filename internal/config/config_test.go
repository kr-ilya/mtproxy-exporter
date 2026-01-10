package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Note: Load() tests are commented out because they cause flag redefinition issues
// in Go's testing framework. The getEnv and getDurationEnv functions are tested instead.

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		want         string
	}{
		{
			name:         "env var set",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "custom",
			want:         "custom",
		},
		{
			name:         "env var not set",
			key:          "TEST_KEY_NOT_SET",
			defaultValue: "default",
			envValue:     "",
			want:         "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			got := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetDurationEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue time.Duration
		envValue     string
		want         time.Duration
	}{
		{
			name:         "valid duration",
			key:          "TEST_DURATION",
			defaultValue: 10 * time.Second,
			envValue:     "30s",
			want:         30 * time.Second,
		},
		{
			name:         "invalid duration",
			key:          "TEST_DURATION_INVALID",
			defaultValue: 10 * time.Second,
			envValue:     "invalid",
			want:         10 * time.Second,
		},
		{
			name:         "not set",
			key:          "TEST_DURATION_NOT_SET",
			defaultValue: 10 * time.Second,
			envValue:     "",
			want:         10 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			got := getDurationEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.want, got)
		})
	}
}
