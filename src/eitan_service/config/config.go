package config

import (
	"fmt"
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
	Env            AppEnv   `default:"local" envconfig:"APP_ENV"`
	Port           int      `default:"5000"`
	GCPProjectID   string   `default:"local" envconfig:"GCP_PROJECT_ID"`
	AllowedOrigins []string `default:"http://local.eitan-flash.com:3000" envconfig:"ALLOWED_ORIGINS"`
	RedisURL       string   `default:"localhost:6379" envconfig:"REDIS_URL"`
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
	return a.Env != Local && a.Env != Test
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
