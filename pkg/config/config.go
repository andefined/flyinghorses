package config

import (
	"fmt"
	"time"
)

// Config Struct strores enviroment variables as configuration settings
type Config struct {
	Env    string `envconfig:"ENV" default:"development"`
	NumCPU int    `envconfig:"NUM_CPU" default:"1"`
	// Logger Settings
	Log struct {
		Level string `envconfig:"LOG_LEVEL" default:"info"`
		Path  string `envconfig:"LOG_PATH" default:"tmp"`
	}
	// gRPC Service Settings
	Service struct {
		Network         string        `envconfig:"SERVICE_NETWORK" default:"tcp"`
		Host            string        `envconfig:"SERVICE_HOST" default:"0.0.0.0"`
		Port            string        `envconfig:"SERVICE_PORT" default:"50050"`
		ReadTimeout     time.Duration `envconfig:"SERVICE_READ_TIMEOUT" default:"10s"`
		WriteTimeout    time.Duration `envconfig:"SERVICE_WRITE_TIMEOUT" default:"20s"`
		ShutdownTimeout time.Duration `envconfig:"SERVICE_SHUTDOWN_TIMEOUT" default:"10s"`
		DomainName      string        `envconfig:"SERVICE_DOMAIN_NAME" default:"flyinghorses.xyz"`
	}
	Elastic struct {
		Host        string        `envconfig:"ELASTIC_URL" default:"http://localhost:9200"`
		User        string        `envconfig:"ELASTIC_USER" default:""`
		Pass        string        `envconfig:"ELASTIC_PASS" default:""`
		Index       string        `envconfig:"ELASTIC_INDEX" default:"the-follower_"`
		DialTimeout time.Duration `envconfig:"DIAL_TIMEOUT" default:"30s"`
	}
	Redis struct {
		Host string `envconfig:"REDIS_HOST" default:"localhost"`
		Port string `envconfig:"REDIS_PORT" default:"6379"`
		Path string `envconfig:"REDIS_PATH" default:"0"`
	}
	SRS struct {
		EARFCN []string `envconfig:"SRS_EARFCN" default:"6400,2850,1700,3050,1451,1426,500,1844,1301,3350,6200,6300"`
		MCC    []string `envconfig:"SRS_MCCS" default:"202"`
		MNC    []string `envconfig:"SRS_MNCS" default:"01,02,03,04,05,06,07,09,10,11,12,13,14,15,16,299,999"`
	}
	Wiggle struct {
		Enabled  bool   `envconfig:"WIGGLE_ENABLED" default:"true"`
		ApiKey   string `envconfig:"WIGGLE_API_KEY" default:""`
		ApiToken string `envconfig:"WIGGLE_API_TOKEN" default:""`
		ApiURL   string `envconfig:"WIGGLE_API_URL" default:""`
	}
	OpenCellId struct {
		Enabled  bool   `envconfig:"OPEN_CELL_ID_ENABLED" default:"true"`
		ApiKey   string `envconfig:"OPEN_CELL_ID_API_KEY" default:""`
		ApiToken string `envconfig:"OPEN_CELL_ID_API_TOKEN" default:""`
		ApiURL   string `envconfig:"OPEN_CELL_ID_API_URL" default:""`
	}
	GPSD struct {
		Enabled bool   `envconfig:"GPSD_ENABLED" default:"true"`
		Host    string `envconfig:"GPSD_HOST" default:"0.0.0.0"`
		Port    string `envconfig:"GPSD_PORT" default:"2947"`
	}
	// TLS Certificates
	TLS struct {
		PEMFile string `envconfig:"TLS_PEM_FILE" default:"server-cert.pem"`
		KEYFile string `envconfig:"TLS_KEY_FILE" default:"server-key.key"`
	}
}

// NewConfig creates a new configuration struct
func NewConfig() *Config {
	return new(Config)
}

// GetServiceURL returns service's URL in host:port foprmat
func (c *Config) GetServiceURL() string {
	return fmt.Sprintf("%s:%s", c.Service.Host, c.Service.Port)
}

// RedisURL returns server `host:port`
func (c *Config) RedisURL() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}
