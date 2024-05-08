package ports

type Config interface {
	//LoadConfig(filePth string) error
	GetDatabaseConfig() DatabaseConfig
	GetHTTPConfig() HTTPConfig
	GetConstants() Constants
	GetStatics() Statics
	GetLoggerConfig() LoggerConfig
	GetBrokerConfig() BrokerConfig
}

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

// Add HTTPConfig struct to hold HTTP server configuration
type HTTPConfig struct {
	Port int `yaml:"port"`
}

// !!
type Exchanges struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Durable bool   `json:"durable"`
}

type Queues struct {
	Name       string
	Durable    bool
	AutoDelete bool
}

type Binding struct {
	Queue      string
	Exchange   []Exchanges
	RoutingKey string
}

type BrokerConfig struct {
	User      string
	Password  string
	Host      string
	Port      string
	Exchanges []Exchanges
	Queues    []Queues
	Bindings  []Binding
}

//!!

// Constants struct holds constant configuration values
type Constants struct {
	MaxItemsPerPage  int `mapstructure:"maxItemsPerPage"`
	MaxRetryAttempts int `mapstructure:"maxRetryAttempts"`
}

// Statics struct holds static configuration values
type Statics struct {
	WelcomeMessage  string `mapstructure:"welcomeMessage"`
	DefaultUserRole string `mapstructure:"defaultUserRole"`
}

type LoggerConfig struct {
	Filename   string `yaml:"filename"` // Path to log file (optional)
	LocalTime  bool   `yaml:"localTime"`
	MaxSize    int    `yaml:"maxSize"`    // Max log file size in megabytes
	MaxBackups int    `yaml:"maxBackups"` // Max number of archived log files
	MaxAge     int    `yaml:"maxAge"`     // Max age of archived logs in days
	Compress   bool   `yaml:"compress"`   // Compress archived log files
	LogLevel   string `yaml:"logLevel"`   // Default log level (optional)
}

//------------------->
