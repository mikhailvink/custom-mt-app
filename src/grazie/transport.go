package grazie

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const _grazieHost = "prod.nmt.grazie.iml.aws.intellij.net"
const _apiUrlPrefix = "/service/v3"

type transport struct {
	host   string
	client httpClient
}

func newTransport(client httpClient, host string) transport {
	return transport{
		host:   host,
		client: client,
	}
}

type httpClient interface {
	UnmarshallJSON(method string, url string, body io.Reader, data interface{}) error
}

func (t transport) Translate(request TranslateRequest) (*TranslateResponse, error) {
	u, err := t.apiUrl("/translate")
	if err != nil {
		return nil, fmt.Errorf("cannot parse request url: %v", err)
	}

	var marshalledRequest []byte
	marshalledRequest, err = json.Marshal(&request)

	var translationResponse TranslateResponse
	err = t.client.UnmarshallJSON(http.MethodPost, u.String(), bytes.NewReader(marshalledRequest), &translationResponse)
	if err != nil {
		return nil, fmt.Errorf("cannot do request to %q: %v", u.String(), err)
	}

	return &translationResponse, nil
}

func (t transport) apiUrl(urlPart string) (*url.URL, error) {
	var host string
	if t.host != "" {
		host = t.host
	} else {
		host = _grazieHost
	}
	return url.Parse("https://" + host + _apiUrlPrefix + urlPart)
}

type TranslateRequest struct {
	Texts  []string `json:"texts"`
	ToLang string   `json:"to_lang"`
}

type TranslateResponse struct {
	Translations []string `json:"translations"`
}
