package provider

import (
	"encoding/json"
	"errors"
	"math"
	"math/rand"
	"time"
)

type RandomProvider struct {
	Min float64
	Max float64
}

func NewRandomProvider(options map[string]interface{}) (*RandomProvider, error) {
	rp := &RandomProvider{}
	if err := rp.ValidateOptions(options); err != nil {
		return nil, err
	}
	return rp, nil
}

func (rp *RandomProvider) Get(date time.Time) (float64, string, error) {
	startTime := time.Now()
	if rp.Min >= rp.Max {
		return 0, "", errors.New("min should be less than max")
	}

	rand.Seed(time.Now().UnixNano())
	randomValue := rp.Min + rand.Float64()*(rp.Max-rp.Min)
	roundedValue := math.Round(randomValue)

	meta := Meta{
		Provider: "RandomProvider",
		Elapsed:  time.Since(startTime).String(),
		Retries:  0, // In this example, we don't have retries, so it's 0.
	}

	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return 0, "", errors.New("could not serialize meta data")
	}

	return roundedValue, string(metaJSON), nil
}

func (rp *RandomProvider) ValidateOptions(options map[string]interface{}) error {
	min, okMin := options["min"].(float64)
	max, okMax := options["max"].(float64)

	if !okMin || !okMax {
		return errors.New("both min and max should be provided and must be float64")
	}

	if min >= max {
		return errors.New("min should be less than max")
	}

	rp.Min = min
	rp.Max = max

	return nil
}
