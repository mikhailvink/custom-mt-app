package grazie

import (
	"crowdin-grazie/http_client"
	"net/http"
)

type Grazie interface {
	Translate(request TranslateRequest) (*TranslateResponse, error)
}

type grazieImpl struct {
	transport transport
}

func New(jwtToken string) Grazie {
	var httpClient = http_client.New(http.DefaultClient)
	httpClient.AddPreprocessFunc(http_client.CreateHeaderSetterPreprocessor("Content-Type", "application/json"))
	httpClient.AddPreprocessFunc(http_client.CreateHeaderSetterPreprocessor("Grazie-Authenticate-JWT", jwtToken))
	return grazieImpl{
		transport: newTransport(httpClient),
	}
}

func (grazie grazieImpl) Translate(request TranslateRequest) (*TranslateResponse, error) {
	return grazie.transport.Translate(request)
}
