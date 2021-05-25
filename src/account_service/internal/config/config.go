package config

import (
	"fmt"
	"github.com/k-yomo/eitan/src/pkg/appenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type AppConfig struct {
	Env                 appenv.Env `default:"local" envconfig:"APP_ENV"`
	HTTPPort            int        `default:"4000" envconfig:"HTTP_PORT"`
	GRPCPort            int        `default:"4040" envconfig:"GRPC_PORT"`
	AppRootURL          string     `default:"http://account.local.eitan-flash.com:4000" envconfig:"APP_ROOT_URL"`
	WebAppURL           string     `default:"http://local.eitan-flash.com:3000" envconfig:"WEB_APP_URL"`
	AllowedOrigins      []string   `default:"http://local.eitan-flash.com:3000" envconfig:"ALLOWED_ORIGINS"`
	SessionCookieDomain string     `default:"local.eitan-flash.com" envconfig:"SESSION_COOKIE_DOMAIN"`
	GCPProjectID        string     `default:"local" envconfig:"GCP_PROJECT_ID"`
	GoogleAuthClientKey string     `envconfig:"GOOGLE_AUTH_CLIENT_KEY"`
	GoogleAuthSecret    string     `envconfig:"GOOGLE_AUTH_SECRET"`
	RedisURL            string     `default:"localhost:6379" envconfig:"REDIS_URL"`
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

type DBConfig struct {
	Driver   string `default:"mysql"`
	DBName   string `default:"accountdb" envconfig:"DB_NAME"`
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
