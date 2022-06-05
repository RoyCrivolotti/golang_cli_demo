// Code generated by MockGen. DO NOT EDIT.
// Source: ./http.go

// Package mock_clients is a generated GoMock package.
package mock_clients

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIHttpClient is a mock of IHttpClient interface.
type MockIHttpClient struct {
	ctrl     *gomock.Controller
	recorder *MockIHttpClientMockRecorder
}

// MockIHttpClientMockRecorder is the mock recorder for MockIHttpClient.
type MockIHttpClientMockRecorder struct {
	mock *MockIHttpClient
}

// NewMockIHttpClient creates a new mock instance.
func NewMockIHttpClient(ctrl *gomock.Controller) *MockIHttpClient {
	mock := &MockIHttpClient{ctrl: ctrl}
	mock.recorder = &MockIHttpClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIHttpClient) EXPECT() *MockIHttpClientMockRecorder {
	return m.recorder
}

// Post mocks base method.
func (m *MockIHttpClient) Post(url string, data interface{}) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Post", url, data)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Post indicates an expected call of Post.
func (mr *MockIHttpClientMockRecorder) Post(url, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Post", reflect.TypeOf((*MockIHttpClient)(nil).Post), url, data)
}
