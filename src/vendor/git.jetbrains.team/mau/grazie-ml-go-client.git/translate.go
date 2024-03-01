package zendeskgo_sell

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

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

func (c *client) Translate(ctx context.Context, taskTag string, langTo string, text string) (string, error) {
	request := TaskRequest{
		Task: "text-translate:" + taskTag,
		Parameters: TaskParameters{
			Data: append(make([]interface{}, 0),
				TaskParameterKey{
					Type: "text",
					Fqdn: "text",
				},
				TaskParameterValue{
					Type:  "text",
					Value: text,
				},
				TaskParameterKey{
					Type: "text",
					Fqdn: "lang",
				},
				TaskParameterValue{
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

	response, err := c.request(ctx, http.MethodPost, c.buildUrl("/v5/task/stream/v2"), bytes.NewReader(data))
	if err != nil {
		return "", errors.Wrap(err, "cannot execute translate request")
	}

	parsedResponse, err := parseResponse(response)
	if err != nil {
		return "", errors.Wrap(err, "cannot parse response body")
	}

	translations := make([]string, 0, len(parsedResponse))
	for _, element := range parsedResponse {
		responsePart := &TaskResponse{}
		err := json.Unmarshal([]byte(element), responsePart)
		if err != nil {
			return "", err
		}
		translations = append(translations, responsePart.Content)
	}

	return strings.Join(translations, ""), nil
}

type TaskRequest struct {
	Task       string         `json:"task"`
	Parameters TaskParameters `json:"parameters"`
}

type TaskParameters struct {
	Data []interface{} `json:"data"`
}

type TaskParameterKey struct {
	Type string `json:"type"`
	Fqdn string `json:"fqdn"`
}

type TaskParameterValue struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type TaskResponse struct {
	Content string `json:"content"`
}
