package environment

import (
	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	Postgres Postgres
}

func Load() (*Configuration, error) {
	configuration := &Configuration{}

	if err := envconfig.Process("", configuration); err != nil {
		return nil, err
	}
	if err := validator.New().Struct(configuration); err != nil {
		return nil, err
	}
	return configuration, nil
}
