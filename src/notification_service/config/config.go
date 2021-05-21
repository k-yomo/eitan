package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type AppEnv string

const (
	Local AppEnv = "local"
	Test  AppEnv = "test"
	Dev   AppEnv = "dev"
	Prod  AppEnv = "prod"
)

type AppConfig struct {
	Env            AppEnv `default:"local" envconfig:"APP_ENV"`
	Port           int    `default:"6000"`
	GCPProjectID   string `default:"local" envconfig:"GCP_PROJECT_ID"`
	SendGridAPIKey string `envconfig:"SEND_GRID_API_KEY"`
}

func NewAppConfig() (*AppConfig, error) {
	appConfig := &AppConfig{}
	if err := envconfig.Process("", appConfig); err != nil {
		return nil, err
	}
	if !map[AppEnv]bool{Local: true, Dev: true, Prod: true}[appConfig.Env] {
		return nil, errors.Errorf("%s is invalid for env", appConfig.Env)
	}
	return appConfig, nil
}

func (a AppConfig) IsDeployedEnv() bool {
	return a.Env != Local && a.Env != Test
}
