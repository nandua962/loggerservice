package entities

// EnvConfig represents the configuration structure for the application.
type EnvConfig struct {
	Debug            bool     `default:"true" split_words:"true"`
	Port             int      `default:"8000" split_words:"true"`
	Db               Database `split_words:"true"`
	AcceptedVersions []string `required:"true" split_words:"true"`
	MigrationPath    string   `split_words:"true"`
	ResetLink        string   `split_words:"true"`
	Redis            Redis
}

// Database represents the configuration for the database connection.
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
type Redis struct {
	Host     string
	UserName string `split_words:"true"`
	Password string
	DB       int `json:"db" default:"0"`
}
