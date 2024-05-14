package zendeskgo_sell

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

const (
	TypeSystemMessage    = "system_message"
	TypeAssistantMessage = "assistant_message"
	TypeUserMessage      = "user_message"

	ProfileOpenaiGpt4    = "openai-gpt-4"
	ProfileOpenaiChatGpt = "openai-chat-gpt"
)

type ChatMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
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

	response, err := c.requestStream(ctx, http.MethodPost, c.buildUrl("/v5/llm/chat/stream/v6"), bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "cannot complete chat messages")
	}

	return &ChatMessage{
		Type:    TypeAssistantMessage,
		Content: response,
	}, nil
}
