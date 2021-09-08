package entrypoints

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

func httpSuccess(writer http.ResponseWriter, payload []byte) {
	_, err := writer.Write(payload)
	if err != nil {
		httpErrorAndLog(writer, fmt.Errorf("cannot write to http response: %v", err), http.StatusInternalServerError)
	}
}

func httpErrorAndLog(writer http.ResponseWriter, err error, httpStatus int) {
	logrus.WithError(err).WithField("status", httpStatus).Error("error")
	http.Error(writer, err.Error(), httpStatus)
}
