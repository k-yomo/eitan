package config

import (
	"github.com/k-yomo/eitan/src/pkg/appenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type AppConfig struct {
	Env            appenv.Env `default:"local" envconfig:"APP_ENV"`
	GCPProjectID   string `default:"local" envconfig:"GCP_PROJECT_ID"`
	SendGridAPIKey string `envconfig:"SEND_GRID_API_KEY"`
}

func NewAppConfig() (*AppConfig, error) {
	appConfig := &AppConfig{}
	if err := envconfig.Process("", appConfig); err != nil {
		return nil, err
	}
	if !appConfig.Env.IsValid() {
		return nil, errors.Errorf("%s is invalid for env", appConfig.Env)
	}
	return appConfig, nil
}
