package environment

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	GrazieHost     = "GRAZIE_HOST"
	EnvGrazieToken = "GRAZIE_TOKEN"
	ClientId       = "CLIENT_ID"
	ClientSecret   = "CLIENT_SECRET"
)

func MustGetEnv(varname string) string {
	value := os.Getenv(varname)
	if value == "" {
		logrus.WithField("envvar", varname).WithError(fmt.Errorf("environment variable is empty")).Fatal("cannot get value")
	}

	return value
}
