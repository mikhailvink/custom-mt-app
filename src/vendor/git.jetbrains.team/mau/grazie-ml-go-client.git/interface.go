//go:generate mockgen -package ${GOPACKAGE} -destination mock_client.go -source interface.go
package zendeskgo_sell

import "context"

type Client interface {
	Chat(ctx context.Context, profile string, messages []ChatMessage) (*ChatMessage, error)
	QuestionAnswering(ctx context.Context, llmProfile string, dataSource string, query string, context string, docsSize int64) (*QuestionAnsweringResponse, error)
	Translate(ctx context.Context, taskTag string, langTo string, text string) (string, error)
	TranslateWithoutAI(ctx context.Context, langFrom string, langTo string, strings []string) (*TranslateResponse, error)
	GetQuota() (string, string)
}
