package grazie

import (
	"crowdin-grazie/http_client"
	"net/http"
)

type Grazie interface {
	Translate(request TranslateRequest) (*TranslateResponse, error)
}

type Config struct {
	Host string
}

type grazieImpl struct {
	transport transport
}

func New(jwtToken string, config Config) Grazie {
	var httpClient = http_client.New(http.DefaultClient)
	httpClient.AddPreprocessFunc(http_client.CreateHeaderSetterPreprocessor("Content-Type", "application/json"))
	httpClient.AddPreprocessFunc(http_client.CreateHeaderSetterPreprocessor("Grazie-Authenticate-JWT", jwtToken))

	return grazieImpl{
		transport: newTransport(httpClient, config.Host),
	}
}

func (grazie grazieImpl) Translate(request TranslateRequest) (*TranslateResponse, error) {
	return grazie.transport.Translate(request)
}
