package grazie

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const _grazieHost = "https://stgn.nmt.grazie.iml.aws.intellij.net"
const _apiUrl = _grazieHost + "/service/v3"

type transport struct {
	client httpClient
}

func newTransport(client httpClient) transport {
	return transport{
		client: client,
	}
}

type httpClient interface {
	UnmarshallJSON(method string, url string, body io.Reader, data interface{}) error
}

func (t transport) Translate(request TranslateRequest) (*TranslateResponse, error) {
	u, err := url.Parse(_apiUrl + "/translate")
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

type TranslateRequest struct {
	Texts  []string `json:"texts"`
	ToLang string   `json:"to_lang"`
}

type TranslateResponse struct {
	Translations []string `json:"translations"`
}
