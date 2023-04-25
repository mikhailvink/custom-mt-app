package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	GrazieHost  string `envconfig:"GRAZIE_HOST" required:"true"`
	GrazieToken string `envconfig:"GRAZIE_TOKEN" required:"true"`

	ClientID     string `envconfig:"CLIENT_ID" required:"true"`
	ClientSecret string `envconfig:"CLIENT_SECRET" required:"true"`

	SlackToken           string `envconfig:"SLACK_TOKEN" required:"true"`
	SlackAlertsChannelID string `envconfig:"SLACK_ALERTS_CHANNEL_ID" required:"true"`
}

func Parse() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	return cfg, err
}
