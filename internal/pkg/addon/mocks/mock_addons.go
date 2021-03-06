// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/pkg/addon/addons.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockworkspaceReader is a mock of workspaceReader interface
type MockworkspaceReader struct {
	ctrl     *gomock.Controller
	recorder *MockworkspaceReaderMockRecorder
}

// MockworkspaceReaderMockRecorder is the mock recorder for MockworkspaceReader
type MockworkspaceReaderMockRecorder struct {
	mock *MockworkspaceReader
}

// NewMockworkspaceReader creates a new mock instance
func NewMockworkspaceReader(ctrl *gomock.Controller) *MockworkspaceReader {
	mock := &MockworkspaceReader{ctrl: ctrl}
	mock.recorder = &MockworkspaceReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockworkspaceReader) EXPECT() *MockworkspaceReaderMockRecorder {
	return m.recorder
}

// ReadAddonsDir mocks base method
func (m *MockworkspaceReader) ReadAddonsDir(svcName string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadAddonsDir", svcName)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadAddonsDir indicates an expected call of ReadAddonsDir
func (mr *MockworkspaceReaderMockRecorder) ReadAddonsDir(svcName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadAddonsDir", reflect.TypeOf((*MockworkspaceReader)(nil).ReadAddonsDir), svcName)
}

// ReadAddon mocks base method
func (m *MockworkspaceReader) ReadAddon(svcName, fileName string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadAddon", svcName, fileName)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadAddon indicates an expected call of ReadAddon
func (mr *MockworkspaceReaderMockRecorder) ReadAddon(svcName, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadAddon", reflect.TypeOf((*MockworkspaceReader)(nil).ReadAddon), svcName, fileName)
}
