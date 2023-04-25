package slack

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"

	"crowdin-grazie/config"
)

type Client interface {
	Error(msg string)
}

type client struct {
	slackClient *slack.Client
	cfg         *config.Config
}

func New(cfg *config.Config) Client {
	return &client{
		slackClient: slack.New(cfg.SlackToken),
		cfg:         cfg,
	}
}

func (c *client) Error(msg string) {
	msg = fmt.Sprintf("*Crowdin Grazie Integration*\n*Error:* %s", msg)

	_, _, err := c.slackClient.PostMessage(c.cfg.SlackAlertsChannelID, slack.MsgOptionText(msg, false))
	if err != nil {
		logrus.WithError(err).WithField("message", msg).Error("cannot send error message")
	}
}
