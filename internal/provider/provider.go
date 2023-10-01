package provider

import (
	"correlateme-data-provider/internal/configs"
	"errors"
	"fmt"
	"time"
)

type Provider interface {
	Get(date time.Time) (value float64, meta string, err error)
	ValidateOptions(options map[string]interface{}) error
}

type Meta struct {
	Provider string `json:"provider"`
	Elapsed  string `json:"elapsed"`
	Retries  int    `json:"retries"`
}

func InitProviders(config configs.Config) (map[string]Provider, error) {
	providerMap := make(map[string]Provider)

	for _, metricConfig := range config.Metrics {
		var provider Provider
		var err error

		switch metricConfig.Provider {
		case "RandomProvider":
			provider, err = NewRandomProvider(metricConfig.Options)
		default:
			return nil, errors.New("unknown provider")
		}

		if err != nil {
			return nil, fmt.Errorf("error creating provider for metric '%s': %w", metricConfig.Name, err)
		}

		providerMap[metricConfig.ID] = provider
	}

	return providerMap, nil
}
