package config

type Logger struct {
	DisableStackTrace    bool `default:"false" envconfig:"LOGGER_DISABLE_STACK_TRACE"`
	ProductionFormatting bool `default:"true" envconfig:"LOGGER_PRODUCTION_FORMATTING"`
}

type (
	// Main app configurations
	MainAppConfig struct {
		Logger
		CustomCSS string `default:"" envconfig:"CUSTOM_CSS"`
	}
)
