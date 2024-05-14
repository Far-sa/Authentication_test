package ports

type Config interface {
	//LoadConfig(filePth string) error
	GetDatabaseConfig() DatabaseConfig
	GetHTTPConfig() HTTPConfig
	// GetConstants() Constants
	// GetStatics() Statics
	// GetLoggerConfig() LoggerConfig
	GetBrokerConfig() BrokerConfig
}

type DatabaseConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbName"`
}

// Add HTTPConfig struct to hold HTTP server configuration
type HTTPConfig struct {
	Port int `yaml:"port"`
}

// !!
type Exchanges struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Durable bool   `yaml:"durable"`
}

type Queues struct {
	Name       string `yaml:"name"`
	Durable    bool   `yaml:"durable"`
	AutoDelete bool   `yaml:"autodelete"`
}

type Binding struct {
	Queue      string      `yaml:"queue"`
	Exchange   []Exchanges `yaml:"exchange"`
	RoutingKey string      `yaml:"routing_key"`
}

type BrokerConfig struct {
	User      string      `yaml:"user"`
	Password  string      `yaml:"password"`
	Host      string      `yaml:"host"`
	Port      string      `yaml:"port"`
	Exchanges []Exchanges `yaml:"exchanges"`
	Queues    []Queues    `yaml:"queues"`
	Bindings  []Binding   `yaml:"bindings"`
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
