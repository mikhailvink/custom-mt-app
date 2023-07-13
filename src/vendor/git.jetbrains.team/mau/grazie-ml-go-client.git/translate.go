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

func (c *client) Translate(ctx context.Context, langFrom string, langTo string, strings []string) (*TranslateResponse, error) {
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
