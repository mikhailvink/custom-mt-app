// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package zendeskgo_sell is a generated GoMock package.
package zendeskgo_sell

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Chat mocks base method.
func (m *MockClient) Chat(ctx context.Context, profile string, messages []ChatMessage) (*ChatMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Chat", ctx, profile, messages)
	ret0, _ := ret[0].(*ChatMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Chat indicates an expected call of Chat.
func (mr *MockClientMockRecorder) Chat(ctx, profile, messages interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Chat", reflect.TypeOf((*MockClient)(nil).Chat), ctx, profile, messages)
}

// GetQuota mocks base method.
func (m *MockClient) GetQuota() (string, string) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuota")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	return ret0, ret1
}

// GetQuota indicates an expected call of GetQuota.
func (mr *MockClientMockRecorder) GetQuota() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuota", reflect.TypeOf((*MockClient)(nil).GetQuota))
}

// QuestionAnswering mocks base method.
func (m *MockClient) QuestionAnswering(ctx context.Context, llmProfile, dataSource, query, context string, docsSize int64) (*QuestionAnsweringResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QuestionAnswering", ctx, llmProfile, dataSource, query, context, docsSize)
	ret0, _ := ret[0].(*QuestionAnsweringResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QuestionAnswering indicates an expected call of QuestionAnswering.
func (mr *MockClientMockRecorder) QuestionAnswering(ctx, llmProfile, dataSource, query, context, docsSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QuestionAnswering", reflect.TypeOf((*MockClient)(nil).QuestionAnswering), ctx, llmProfile, dataSource, query, context, docsSize)
}

// Translate mocks base method.
func (m *MockClient) Translate(ctx context.Context, taskTag, langTo, text string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Translate", ctx, taskTag, langTo, text)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Translate indicates an expected call of Translate.
func (mr *MockClientMockRecorder) Translate(ctx, taskTag, langTo, text interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Translate", reflect.TypeOf((*MockClient)(nil).Translate), ctx, taskTag, langTo, text)
}

// TranslateWithoutAI mocks base method.
func (m *MockClient) TranslateWithoutAI(ctx context.Context, langFrom, langTo string, strings []string) (*TranslateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TranslateWithoutAI", ctx, langFrom, langTo, strings)
	ret0, _ := ret[0].(*TranslateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TranslateWithoutAI indicates an expected call of TranslateWithoutAI.
func (mr *MockClientMockRecorder) TranslateWithoutAI(ctx, langFrom, langTo, strings interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TranslateWithoutAI", reflect.TypeOf((*MockClient)(nil).TranslateWithoutAI), ctx, langFrom, langTo, strings)
}
