package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	GCPProjectID string `default:"local" envconfig:"GCP_PROJECT_ID"`
	DBDriver     string `default:"mysql"`
	DBName       string `envconfig:"DB_NAME" required:"true"`
	DBUser       string `default:"mysql" envconfig:"DB_USER"`
	DBPassword   string `default:"mysql" envconfig:"DB_PASSWORD"`
	DBURL        string `default:"localhost:13306" envconfig:"DB_URL"`
}

func (c *config) dbDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&loc=Asia%%2FTokyo", c.DBUser, c.DBPassword, c.DBURL, c.DBName)
}

func newConfig() (*config, error) {
	config := config{}
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}