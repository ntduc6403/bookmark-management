package api

import (
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

// Config holds the configuration values for the API server, including the application port, service name, and instance ID.
type Config struct {
	AppPort     string `default:"8080" envconfig:"APP_PORT"`
	ServiceName string `default:"bookmark_service" envconfig:"SERVICE_NAME"`
	InstanceID  string `envconfig:"INSTANCE_ID"`
}

// NewConfig creates a new instance of the Config struct and populates it with values from environment variables using the envconfig package. If there is an error during processing, it returns the error.
func NewConfig() (*Config, error) {

	cfg := &Config{}
	err := envconfig.Process("api", cfg)
	if err != nil {
		return nil, err
	}

	if cfg.InstanceID == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}

		cfg.InstanceID = id.String()
	}

	return cfg, err
}
