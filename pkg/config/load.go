package config

import "github.com/kelseyhightower/envconfig"

// Load loads the configuration from the environment.
func LoadMainAppConfig() (MainAppConfig, error) {
	config := MainAppConfig{}
	err := envconfig.Process("", &config)
	if err != nil {
		return config, err
	}

	return config, err
}
