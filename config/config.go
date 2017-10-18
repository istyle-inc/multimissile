package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

const (
	// DefaultPort port number served multimissile by default
	DefaultPort = "29300"
	// DefaultLogLevel default log level multimissile use error level
	DefaultLogLevel = "error"
	// DefaultTimeout default time request will be timeout
	DefaultTimeout = 5
	// DefaultMaxIdleConnsPerHost max value of idle connection per each host
	DefaultMaxIdleConnsPerHost = 100
	// DefaultIdleConnTimeout default time idle connection will be expired
	DefaultIdleConnTimeout = 30
	// DefaultProxyReadTimeout default time proxy read will be timeout
	DefaultProxyReadTimeout = 60
	// DefaultShutdownTimeout this might not be used any where?
	DefaultShutdownTimeout = 10
)

// Config struct of configure
type Config struct {
	Port                string
	LogLevel            string
	Timeout             int
	MaxIdleConnsPerHost int
	DisableCompression  bool
	IdleConnTimeout     int
	ProxyReadTimeout    int
	ShutdownTimeout     int
	Endpoints           []EndPoint
}

// EndPoint struct of one of Endpoints
type EndPoint struct {
	// Endpoint Name
	Name string
	// Endpoint URL
	Ep string
	// Headers to set http-headers
	ProxySetHeaders [][]string
	// Headers to pass from origin-request
	ProxyPassHeaders [][]string
	// Threshold values of http status code, would be recognize as success
	AcceptableHTTPStatuses []int
	// Threshold values of http status code, would be recognize as failure
	ExceptableHTTPStatuses []int
}

// LoadBytes load config file and unmarshal to config struct
func LoadBytes(bytes []byte) (Config, error) {
	var config Config
	err := toml.Unmarshal(bytes, &config)
	return config, err
}

// Load load config from file path
func Load(confPath string) (Config, error) {
	var config Config
	bytes, err := ioutil.ReadFile(confPath)
	if err != nil {
		return config, err
	}

	config, err = LoadBytes(bytes)
	if err != nil {
		return config, err
	}

	if config.Port == "" {
		config.Port = DefaultPort
	}

	if config.LogLevel == "" {
		config.LogLevel = DefaultLogLevel
	}

	if config.Timeout <= 0 {
		config.Timeout = DefaultTimeout
	}

	if config.MaxIdleConnsPerHost <= 0 {
		config.MaxIdleConnsPerHost = DefaultMaxIdleConnsPerHost
	}

	if config.IdleConnTimeout <= 0 {
		config.IdleConnTimeout = DefaultIdleConnTimeout
	}

	if config.ProxyReadTimeout <= 0 {
		config.ProxyReadTimeout = DefaultProxyReadTimeout
	}

	if config.ShutdownTimeout <= 0 {
		config.ShutdownTimeout = DefaultShutdownTimeout
	}

	if len(config.Endpoints) == 0 {
		return config, errors.New("empty Endpoints")
	}

	for _, ep := range config.Endpoints {
		if ep.Name == "" {
			return config, errors.New("empty Endpoint name")
		}
		if ep.Ep == "" {
			return config, errors.New("empty Endpoint URL")
		}
	}

	return config, nil
}

// FindEp search endpoint using name
func FindEp(conf Config, name string) (EndPoint, error) {
	for _, ep := range conf.Endpoints {
		if ep.Name == name {
			return ep, nil
		}
	}

	return EndPoint{}, fmt.Errorf("ep:%s is not found", name)
}
