// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mock_executor is a generated GoMock package.
package mock_executor

import (
	gomock "github.com/golang/mock/gomock"
	pb "github.com/meshplus/bitxhub-model/pb"
	reflect "reflect"
)

// MockExecutor is a mock of Executor interface
type MockExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockExecutorMockRecorder
}

// MockExecutorMockRecorder is the mock recorder for MockExecutor
type MockExecutorMockRecorder struct {
	mock *MockExecutor
}

// NewMockExecutor creates a new mock instance
func NewMockExecutor(ctrl *gomock.Controller) *MockExecutor {
	mock := &MockExecutor{ctrl: ctrl}
	mock.recorder = &MockExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockExecutor) EXPECT() *MockExecutorMockRecorder {
	return m.recorder
}

// Start mocks base method
func (m *MockExecutor) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockExecutorMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockExecutor)(nil).Start))
}

// Stop mocks base method
func (m *MockExecutor) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop
func (mr *MockExecutorMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockExecutor)(nil).Stop))
}

// HandleIBTP mocks base method
func (m *MockExecutor) HandleIBTP(ibtp *pb.IBTP) *pb.IBTP {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleIBTP", ibtp)
	ret0, _ := ret[0].(*pb.IBTP)
	return ret0
}

// HandleIBTP indicates an expected call of HandleIBTP
func (mr *MockExecutorMockRecorder) HandleIBTP(ibtp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleIBTP", reflect.TypeOf((*MockExecutor)(nil).HandleIBTP), ibtp)
}

// Rollback mocks base method
func (m *MockExecutor) Rollback(ibtp *pb.IBTP, isSrcChain bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Rollback", ibtp, isSrcChain)
}

// Rollback indicates an expected call of Rollback
func (mr *MockExecutorMockRecorder) Rollback(ibtp, isSrcChain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rollback", reflect.TypeOf((*MockExecutor)(nil).Rollback), ibtp, isSrcChain)
}

// QueryLatestMeta mocks base method
func (m *MockExecutor) QueryLatestMeta() map[string]uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryLatestMeta")
	ret0, _ := ret[0].(map[string]uint64)
	return ret0
}

// QueryLatestMeta indicates an expected call of QueryLatestMeta
func (mr *MockExecutorMockRecorder) QueryLatestMeta() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryLatestMeta", reflect.TypeOf((*MockExecutor)(nil).QueryLatestMeta))
}

// QueryReceipt mocks base method
func (m *MockExecutor) QueryReceipt(from string, idx uint64, originalIBTP *pb.IBTP) (*pb.IBTP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryReceipt", from, idx, originalIBTP)
	ret0, _ := ret[0].(*pb.IBTP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryReceipt indicates an expected call of QueryReceipt
func (mr *MockExecutorMockRecorder) QueryReceipt(from, idx, originalIBTP interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryReceipt", reflect.TypeOf((*MockExecutor)(nil).QueryReceipt), from, idx, originalIBTP)
}

// QueryDstRollbackMeta mocks base method
func (m *MockExecutor) QueryDstRollbackMeta() map[string]uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryDstRollbackMeta")
	ret0, _ := ret[0].(map[string]uint64)
	return ret0
}

// QueryDstRollbackMeta indicates an expected call of QueryDstRollbackMeta
func (mr *MockExecutorMockRecorder) QueryDstRollbackMeta() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryDstRollbackMeta", reflect.TypeOf((*MockExecutor)(nil).QueryDstRollbackMeta))
}
