// Code generated by MockGen. DO NOT EDIT.
// Source: store.go

// Package mock_store is a generated GoMock package.
package mock_store

import (
	context "context"
	reflect "reflect"

	store "github.com/Nastez/shortener/internal/store"
	gomock "github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Bootstrap mocks base method.
func (m *MockStore) Bootstrap(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bootstrap", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Bootstrap indicates an expected call of Bootstrap.
func (mr *MockStoreMockRecorder) Bootstrap(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bootstrap", reflect.TypeOf((*MockStore)(nil).Bootstrap), ctx)
}

// Get mocks base method.
func (m *MockStore) Get(ctx context.Context, id string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockStoreMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStore)(nil).Get), ctx, id)
}

// Save mocks base method.
func (m *MockStore) Save(ctx context.Context, url store.URL) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, url)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockStoreMockRecorder) Save(ctx, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockStore)(nil).Save), ctx, url)
}
