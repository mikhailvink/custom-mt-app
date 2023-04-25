package entrypoints

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"crowdin-grazie/slack"
)

type HandlerCreator struct {
	slackClient slack.Client
}

func NewHandlerCreator(slackClient slack.Client) *HandlerCreator {
	return &HandlerCreator{
		slackClient: slackClient,
	}
}

func (hc *HandlerCreator) httpSuccess(writer http.ResponseWriter, payload []byte) {
	_, err := writer.Write(payload)
	if err != nil {
		hc.httpErrorAndLog(writer, fmt.Errorf("cannot write to http response: %v", err), http.StatusInternalServerError)
	}
}

func (hc *HandlerCreator) httpErrorAndLog(writer http.ResponseWriter, err error, httpStatus int) {
	hc.slackClient.Error(err.Error())
	logrus.WithError(err).WithField("status", httpStatus).Error("error")
	http.Error(writer, err.Error(), httpStatus)
}
