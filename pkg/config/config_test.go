package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("SERVER_HOST", "testhost")

	// Clean up environment variables after the test
	defer os.Unsetenv("SERVER_PORT")
	defer os.Unsetenv("SERVER_HOST")

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
	os.Setenv("SERVER_PORT", "invalid")

	// Clean up environment variables after the test
	defer os.Unsetenv("SERVER_PORT")

	// Create a new instance of the test config
	cfg := &TestConfig{}

	// Load the configuration
	_, err := LoadConfig(cfg)

	// Assert that there was a validation error
	assert.Error(t, err)
}
