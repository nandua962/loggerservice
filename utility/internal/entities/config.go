package entities

// EnvConfig represents the configuration structure for the application.
type EnvConfig struct {
	Debug                  bool     `default:"true" split_words:"true"`  // Indicates whether the application is in debug mode (default: true)
	Port                   int      `default:"8080" split_words:"true"`  // The port on which the server listens (default: 8080)
	Db                     Database `split_words:"true"`                 // Database configuration
	AcceptedVersions       []string `required:"true" split_words:"true"` // List of accepted API versions (required)
	LocalisationServiceURL string   `split_words:"true"`                 // URL for the localization service
	LoggerServiceURL       string   `split_words:"true"`                 // URL for the logger service
	LoggerSecret           string   `split_words:"true"`                 // Secret key for logging
	EndpointURL            string   `split_words:"true"`                 // URL for the localization endpoint
	ErrorHelpLink          string   `split_words:"true"`
}

// Database represents the database configuration for the application.
type Database struct {
	User      string
	Password  string
	Port      int
	Host      string
	DATABASE  string
	Schema    string
	MaxActive int
	MaxIdle   int
}
