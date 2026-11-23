package redis

import "github.com/kelseyhightower/envconfig"

type config struct {
	Address  string `default:"localhost:6379" envconfig:"REDIS_ADDRESS"`
	Password string `default:"" envconfig:"REDIS_PASSWORD"`
	DB       int    `default:"0" envconfig:"REDIS_DB"`
}

// NewConfig creates a new instance of the Config struct and populates it with values from environment variables using the envconfig package. If there is an error during processing, it returns the error.
func newConfig(envPrefix string) (*config, error) {

	cfg := &config{}
	err := envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, err
}
