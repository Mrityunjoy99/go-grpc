package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfig defines the structure for testing configuration.
type TestConfig struct {
	Server TestServerConfig `validate:"required"`
}

// TestServerConfig defines server-specific configuration for tests.
type TestServerConfig struct {
	Port string `validate:"required,numeric" default:"8080"`
	Host string `validate:"required" default:"localhost"`
}

func TestLoadConfig_Success(t *testing.T) {
	// Set environment variables for the test
	err := os.Setenv("SERVER_PORT", "9090")
	if err != nil {
		t.Fatalf("failed to set environment variable: %v", err)
	}

	err = os.Setenv("SERVER_HOST", "testhost")
	if err != nil {
		t.Fatalf("failed to set environment variable: %v", err)
	}

	// Clean up environment variables after the test
	defer func() {
		if err := os.Unsetenv("SERVER_PORT"); err != nil {
			t.Logf("failed to unset environment variable: %v", err)
		}
	}()
	defer func() {
		if err := os.Unsetenv("SERVER_HOST"); err != nil {
			t.Logf("failed to unset environment variable: %v", err)
		}
	}()

	// Create a new instance of the test config
	cfg := &TestConfig{}

	// Load the configuration
	loadedCfg, err := LoadConfig(cfg)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the configuration was loaded correctly
	assert.Equal(t, "9090", loadedCfg.Server.Port)
	assert.Equal(t, "testhost", loadedCfg.Server.Host)
}

func TestLoadConfig_DefaultValues(t *testing.T) {
	// Create a new instance of the test config
	cfg := &TestConfig{}

	// Load the configuration without setting any environment variables
	loadedCfg, err := LoadConfig(cfg)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the default values were loaded correctly
	assert.Equal(t, "8080", loadedCfg.Server.Port)
	assert.Equal(t, "localhost", loadedCfg.Server.Host)
}

func TestLoadConfig_ValidationError(t *testing.T) {
	// Set an invalid environment variable for the test
	err := os.Setenv("SERVER_PORT", "invalid")
	if err != nil {
		t.Fatalf("failed to set environment variable: %v", err)
	}

	// Clean up environment variables after the test
	defer func() {
		if err := os.Unsetenv("SERVER_PORT"); err != nil {
			t.Logf("failed to unset environment variable: %v", err)
		}
	}()

	// Create a new instance of the test config
	cfg := &TestConfig{}

	// Load the configuration and expect an error
	_, err = LoadConfig(cfg)
	
	// Assert that there was a validation error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'Port' failed on the 'numeric' tag")
}
