package config

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"gopkg.in/go-playground/validator.v9"
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
	Port                string     `validate:"required"`
	LogLevel            string     `validate:"required"`
	Timeout             int        `validate:"required"`
	MaxIdleConnsPerHost int        `validate:"required"`
	DisableCompression  bool       `validate:"required"`
	IdleConnTimeout     int        `validate:"required"`
	ProxyReadTimeout    int        `validate:"required"`
	ShutdownTimeout     int        `validate:"required"`
	Endpoints           []EndPoint `validate:"required"`
}

// EndPoint struct of one of Endpoints
type EndPoint struct {
	Name                   string `validate:"required"`
	To                     string `validate:"required"`
	ProxySetHeaders        [][]string
	ProxyPassHeaders       [][]string
	AcceptableHTTPStatuses []int
	ExceptableHTTPStatuses []int
}

func initialize() Config {
	return Config{
		Port:                DefaultPort,
		LogLevel:            DefaultLogLevel,
		Timeout:             DefaultTimeout,
		MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
		IdleConnTimeout:     DefaultIdleConnTimeout,
		ProxyReadTimeout:    DefaultProxyReadTimeout,
		ShutdownTimeout:     DefaultShutdownTimeout,
	}
}

// LoadBytes load config file and unmarshal to config struct
func LoadBytes(bytes []byte) (config Config, err error) {
	config = initialize()
	err = toml.Unmarshal(bytes, &config)

	return config, err
}

// Load load config from file path
func Load(confPath string) (Config, error) {
	var config Config
	bytes, err := ioutil.ReadFile(confPath)
	if err != nil {
		return config, err
	}

	if config, err = LoadBytes(bytes); err != nil {
		return config, err
	}

	validate := validator.New()
	err = validate.Struct(config)

	return config, err
}

// FindTo search endpoint using name
func FindTo(conf Config, name string) (EndPoint, error) {
	for _, ep := range conf.Endpoints {
		if ep.Name == name {
			return ep, nil
		}
	}

	return EndPoint{}, fmt.Errorf("ep:%s is not found", name)
}
