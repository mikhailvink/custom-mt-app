package zendeskgo_sell

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

const (
	IdeaDataSource  = "helpdesk"
	RiderDataSource = "rider"

	ChatGpt4LlmProfile = "openai-gpt-4"
)

func (c *client) QuestionAnswering(ctx context.Context, llmProfile string, dataSource string, query string, docsSize int64) (*QuestionAnsweringResponse, error) {
	request := QuestionAnsweringRequest{
		Query:      query,
		Size:       docsSize,
		DataSource: dataSource,
		LlmProfile: llmProfile,
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "cannot marshal request")
	}

	response, err := c.request(ctx, http.MethodPost, c.buildUrl("/v5/meta/qa/answer/v1"), bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "cannot execute question answering request")
	}

	qaResponse := &QuestionAnsweringResponse{}
	err = json.Unmarshal([]byte(response), qaResponse)
	if err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal response")
	}

	return qaResponse, nil
}

type QuestionAnsweringRequest struct {
	Query      string `json:"query"`
	Size       int64  `json:"size"`
	DataSource string `json:"dataSource"`
	LlmProfile string `json:"llmProfile,omitempty"`
}

type QuestionAnsweringResponse struct {
	Answers Answers `json:"answers"`
}

type Answers struct {
	Documents       []Document      `json:"documents"`
	SummaryWithRefs SummaryWithRefs `json:"summaryWithRefs"`
}

type SummaryWithRefs struct {
	Text       string   `json:"text"`
	References []string `json:"references"`
}

type DocID struct {
	Internal string `json:"internal"`
	Readable string `json:"readable"`
}

type Selection struct {
	Start        int64 `json:"start"`
	EndExclusive int64 `json:"endExclusive"`
}

type Document struct {
	ID        DocID     `json:"id"`
	Content   string    `json:"content"`
	Selection Selection `json:"selection"`
}
