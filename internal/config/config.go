package config

import (
	"os"
	"time"
	"errors"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port int `yaml:"port"`
	Strategy string `yaml:"strategy"`
	HealthCheckInterval time.Duration `yaml:"health_check_interval"`
	Backends []*BackConf `yaml:"backends"`
}

type BackConf struct {
	URL string `yaml:"url"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	if config.Port <= 0 ||
		config.Strategy == "" ||
		config.HealthCheckInterval == 0 ||
		len(config.Backends) == 0 {
			return nil, errors.New("Invalid data")
	}
	return &config, nil
}
