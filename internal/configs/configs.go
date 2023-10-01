package configs

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Env          string
	DB           string
	SystemUserID string
	Debug        bool
	Metrics      []MetricConfig
}

type MetricConfig struct {
	Name     string
	ID       string
	Schedule string
	Provider string
}

func ReadConfigs(path string) (Config, error) {
	var cfg Config
	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
