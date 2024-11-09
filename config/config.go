package config

import (
	"errors"
	"fmt"
	"io/ioutil" //nolint:staticcheck
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		URI             string `yaml:"uri"`
		Username        string `yaml:"username"`
		Password        string `yaml:"password"`
		DBName          string `yaml:"db_name"`
		MaxPoolSize     uint64 `yaml:"max_pool_size"`
		MaxConnIdleTime uint64 `yaml:"max_conn_idle_time"`
		RetryWrites     bool   `yaml:"retry_writes"`
	} `yaml:"database"`

	Server struct {
		Port        int    `yaml:"port"`
		Environment string `yaml:"environment"`
		AppHost     string `yaml:"app_host"`
		AppDomain   string `yaml:"app_domain"`
	} `yaml:"server"`

	Auth struct {
		JWTSecret string `yaml:"jwt_secret"`
	} `yaml:"auth"`
}

var (
	config                 *Config
	configfileUnmarshalErr = func(err error) error {
		return errors.New(fmt.Sprintf("error unmarshalling config data: %w", err)) //nolint:gosimple,govet
	}
	configFileReadingError = func(filename string, err error) error {
		return errors.New(fmt.Sprintf("error reading config file %s: %w", filename, err)) //nolint:gosimple,govet
	}
)

// LoadConfig reads the specified YAML config file based on the environment.
func LoadConfig(env string) error {
	filename := os.Getenv("CONFIG_FILE")

	if filename == "" {
		filename = fmt.Sprintf("config/config.%s.yml", env)
	}

	config = &Config{}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return configFileReadingError(filename, err)
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return configfileUnmarshalErr(err)
	}

	log.Printf("Loaded config file: %#v", config)

	return nil
}

func GetConfig() *Config {
	if config == nil {
		LoadConfig(os.Getenv("ENV")) // nolint:errcheck,nolintlint
	}
	return config
}
