// Code generated by MockGen. DO NOT EDIT.
// Source: internal/application/nats.go

// Package mock_app is a generated GoMock package.
package mock_app

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	nats "github.com/nats-io/nats.go"
)

// MockNatsClient is a mock of NatsClient interface.
type MockNatsClient struct {
	ctrl     *gomock.Controller
	recorder *MockNatsClientMockRecorder
}

// MockNatsClientMockRecorder is the mock recorder for MockNatsClient.
type MockNatsClientMockRecorder struct {
	mock *MockNatsClient
}

// NewMockNatsClient creates a new mock instance.
func NewMockNatsClient(ctrl *gomock.Controller) *MockNatsClient {
	mock := &MockNatsClient{ctrl: ctrl}
	mock.recorder = &MockNatsClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNatsClient) EXPECT() *MockNatsClientMockRecorder {
	return m.recorder
}

// AddConsumer mocks base method.
func (m *MockNatsClient) AddConsumer(streamName, consumerName, subject string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddConsumer", streamName, consumerName, subject)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddConsumer indicates an expected call of AddConsumer.
func (mr *MockNatsClientMockRecorder) AddConsumer(streamName, consumerName, subject interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddConsumer", reflect.TypeOf((*MockNatsClient)(nil).AddConsumer), streamName, consumerName, subject)
}

// CreateStream mocks base method.
func (m *MockNatsClient) CreateStream(streamName, streamSubjects string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateStream", streamName, streamSubjects)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateStream indicates an expected call of CreateStream.
func (mr *MockNatsClientMockRecorder) CreateStream(streamName, streamSubjects interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStream", reflect.TypeOf((*MockNatsClient)(nil).CreateStream), streamName, streamSubjects)
}

// PublishMessage mocks base method.
func (m *MockNatsClient) PublishMessage(subject, message string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PublishMessage", subject, message)
}

// PublishMessage indicates an expected call of PublishMessage.
func (mr *MockNatsClientMockRecorder) PublishMessage(subject, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishMessage", reflect.TypeOf((*MockNatsClient)(nil).PublishMessage), subject, message)
}

// PublishMessageEphemeral mocks base method.
func (m *MockNatsClient) PublishMessageEphemeral(subject, message string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PublishMessageEphemeral", subject, message)
}

// PublishMessageEphemeral indicates an expected call of PublishMessageEphemeral.
func (mr *MockNatsClientMockRecorder) PublishMessageEphemeral(subject, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishMessageEphemeral", reflect.TypeOf((*MockNatsClient)(nil).PublishMessageEphemeral), subject, message)
}

// SubscribeDurable mocks base method.
func (m *MockNatsClient) SubscribeDurable(subject, streamName, consumerName string, handler func(*nats.Msg) error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SubscribeDurable", subject, streamName, consumerName, handler)
}

// SubscribeDurable indicates an expected call of SubscribeDurable.
func (mr *MockNatsClientMockRecorder) SubscribeDurable(subject, streamName, consumerName, handler interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeDurable", reflect.TypeOf((*MockNatsClient)(nil).SubscribeDurable), subject, streamName, consumerName, handler)
}

// SubscribeEphemeral mocks base method.
func (m *MockNatsClient) SubscribeEphemeral(subject string, handler func(*nats.Msg) error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SubscribeEphemeral", subject, handler)
}

// SubscribeEphemeral indicates an expected call of SubscribeEphemeral.
func (mr *MockNatsClientMockRecorder) SubscribeEphemeral(subject, handler interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeEphemeral", reflect.TypeOf((*MockNatsClient)(nil).SubscribeEphemeral), subject, handler)
}