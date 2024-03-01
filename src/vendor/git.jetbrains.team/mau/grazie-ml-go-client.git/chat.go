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
	RoleSystem    = "System"
	RoleAssistant = "Assistant"
	RoleUser      = "User"

	ProfileOpenaiGpt4    = "openai-gpt-4"
	ProfileOpenaiChatGpt = "openai-chat-gpt"
)

type ChatMessage struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type chatRequest struct {
	Chat    chat   `json:"chat"`
	Profile string `json:"profile"`
}

type chat struct {
	Messages []ChatMessage `json:"messages"`
}

func (c *client) Chat(ctx context.Context, profile string, messages []ChatMessage) (*ChatMessage, error) {
	request := chatRequest{
		Chat: chat{
			Messages: messages,
		},
		Profile: profile,
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "cannot marshal request")
	}

	response, err := c.request(ctx, http.MethodPost, c.buildUrl("/v5/llm/chat/stream/v3"), bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "cannot complete chat messages")
	}

	parsedResponse, err := parseChatResponse(response)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse response body")
	}

	return &ChatMessage{
		Role: RoleAssistant,
		Text: strings.Join(parsedResponse, ""),
	}, nil
}
