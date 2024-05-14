package zendeskgo_sell

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

const (
	LangEN      = "en"
	LangDE      = "de"
	LangFR      = "fr"
	LangES      = "es"
	LangRU      = "ru"
	LangKO      = "ko"
	LangZH      = "zh"
	LangJA      = "ja"
	LangUnknown = "unknown"
)

const (
	GrazieTranslateTag  = "llm-grazie-default"
	CrowdinTranslateTag = "llm-crowdin-default"
)

func (c *client) TranslateWithoutAI(ctx context.Context, langFrom string, langTo string, strings []string) (*TranslateResponse, error) {
	request := TranslateRequest{
		Texts:    strings,
		FromLang: langFrom,
		ToLang:   langTo,
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "cannot marshal request")
	}

	response, err := c.request(ctx, http.MethodPost, c.buildUrl("/v5/trf/nmt/translate"), bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "cannot execute translate request")
	}

	translateResponse := &TranslateResponse{}
	err = json.Unmarshal([]byte(response), translateResponse)
	if err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal response")
	}

	return translateResponse, nil
}

type TranslateRequest struct {
	Texts    []string `json:"texts"`
	FromLang string   `json:"fromLang"`
	ToLang   string   `json:"toLang"`
}

type TranslateResponse struct {
	Translations []Translation `json:"translations"`
}

type Translation struct {
	Text        string `json:"text"`
	Translation string `json:"translation"`
	FromLang    string `json:"fromLang"`
	ToLang      string `json:"toLang"`
}

func (c *client) Translate(ctx context.Context, taskTag string, langTo string, text string) (string, error) {
	request := taskRequest{
		Task: "text-translate:" + taskTag,
		Parameters: taskParameters{
			Data: append(make([]interface{}, 0),
				taskParameterKey{
					Type: "text",
					Fqdn: "text",
				},
				taskParameterValue{
					Type:  "text",
					Value: text,
				},
				taskParameterKey{
					Type: "text",
					Fqdn: "lang",
				},
				taskParameterValue{
					Type:  "text",
					Value: langTo,
				},
			),
		},
	}
	data, err := json.Marshal(request)
	if err != nil {
		return "", errors.Wrap(err, "cannot marshal request")
	}

	response, err := c.requestStream(ctx, http.MethodPost, c.buildUrl("/v5/task/stream/v2"), bytes.NewReader(data))
	if err != nil {
		return "", errors.Wrap(err, "cannot execute translate request")
	}

	return response, nil
}

type taskRequest struct {
	Task       string         `json:"task"`
	Parameters taskParameters `json:"parameters"`
}

type taskParameters struct {
	Data []interface{} `json:"data"`
}

type taskParameterKey struct {
	Type string `json:"type"`
	Fqdn string `json:"fqdn"`
}

type taskParameterValue struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
