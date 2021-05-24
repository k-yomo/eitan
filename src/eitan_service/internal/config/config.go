package config

import (
	"fmt"
	"github.com/k-yomo/eitan/src/pkg/appenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type AppConfig struct {
	Env            appenv.Env `default:"local" envconfig:"APP_ENV"`
	Port           int        `default:"5000"`
	GCPProjectID   string     `default:"local" envconfig:"GCP_PROJECT_ID"`
	AllowedOrigins []string   `default:"http://local.eitan-flash.com:3000" envconfig:"ALLOWED_ORIGINS"`
	RedisURL       string     `default:"localhost:6379" envconfig:"REDIS_URL"`
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

type APIConfig struct {
	AccountServiceGRPCURL string `default:"localhost:4040" envconfig:"ACCOUNT_SERVICE_GRPC_URL"` // domain:port
}

func NewAPIConfig() (*APIConfig, error) {
	apiConfig := &APIConfig{}
	if err := envconfig.Process("", apiConfig); err != nil {
		return nil, err
	}
	return apiConfig, nil
}

func (a AppConfig) IsDeployedEnv() bool {
	return a.Env.IsDeployed()
}

type DBConfig struct {
	Driver   string `default:"mysql"`
	DBName   string `default:"eitandb" envconfig:"DB_NAME"`
	User     string `default:"mysql" envconfig:"DB_USER"`
	Password string `default:"mysql" envconfig:"DB_PASSWORD"`
	URL      string `default:"localhost:13306" envconfig:"DB_URL"`
}

func NewDBConfig() (*DBConfig, error) {
	dbConfig := &DBConfig{}
	if err := envconfig.Process("", dbConfig); err != nil {
		return nil, err
	}
	return dbConfig, nil
}

func (d *DBConfig) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&loc=Asia%%2FTokyo", d.User, d.Password, d.URL, d.DBName)
}
