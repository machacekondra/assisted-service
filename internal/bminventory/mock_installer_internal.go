// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/openshift/assisted-service/internal/bminventory (interfaces: InstallerInternals)

// Package bminventory is a generated GoMock package.
package bminventory

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	common "github.com/openshift/assisted-service/internal/common"
	installer "github.com/openshift/assisted-service/restapi/operations/installer"
	types "k8s.io/apimachinery/pkg/types"
	reflect "reflect"
)

// MockInstallerInternals is a mock of InstallerInternals interface
type MockInstallerInternals struct {
	ctrl     *gomock.Controller
	recorder *MockInstallerInternalsMockRecorder
}

// MockInstallerInternalsMockRecorder is the mock recorder for MockInstallerInternals
type MockInstallerInternalsMockRecorder struct {
	mock *MockInstallerInternals
}

// NewMockInstallerInternals creates a new mock instance
func NewMockInstallerInternals(ctrl *gomock.Controller) *MockInstallerInternals {
	mock := &MockInstallerInternals{ctrl: ctrl}
	mock.recorder = &MockInstallerInternalsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockInstallerInternals) EXPECT() *MockInstallerInternalsMockRecorder {
	return m.recorder
}

// DeregisterClusterInternal mocks base method
func (m *MockInstallerInternals) DeregisterClusterInternal(arg0 context.Context, arg1 installer.DeregisterClusterParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeregisterClusterInternal", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeregisterClusterInternal indicates an expected call of DeregisterClusterInternal
func (mr *MockInstallerInternalsMockRecorder) DeregisterClusterInternal(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeregisterClusterInternal", reflect.TypeOf((*MockInstallerInternals)(nil).DeregisterClusterInternal), arg0, arg1)
}

// GenerateClusterISOInternal mocks base method
func (m *MockInstallerInternals) GenerateClusterISOInternal(arg0 context.Context, arg1 installer.GenerateClusterISOParams) (*common.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateClusterISOInternal", arg0, arg1)
	ret0, _ := ret[0].(*common.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateClusterISOInternal indicates an expected call of GenerateClusterISOInternal
func (mr *MockInstallerInternalsMockRecorder) GenerateClusterISOInternal(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateClusterISOInternal", reflect.TypeOf((*MockInstallerInternals)(nil).GenerateClusterISOInternal), arg0, arg1)
}

// GetClusterByKubeKey mocks base method
func (m *MockInstallerInternals) GetClusterByKubeKey(arg0 types.NamespacedName) (*common.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClusterByKubeKey", arg0)
	ret0, _ := ret[0].(*common.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClusterByKubeKey indicates an expected call of GetClusterByKubeKey
func (mr *MockInstallerInternalsMockRecorder) GetClusterByKubeKey(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClusterByKubeKey", reflect.TypeOf((*MockInstallerInternals)(nil).GetClusterByKubeKey), arg0)
}

// GetClusterInternal mocks base method
func (m *MockInstallerInternals) GetClusterInternal(arg0 context.Context, arg1 installer.GetClusterParams) (*common.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClusterInternal", arg0, arg1)
	ret0, _ := ret[0].(*common.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClusterInternal indicates an expected call of GetClusterInternal
func (mr *MockInstallerInternalsMockRecorder) GetClusterInternal(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClusterInternal", reflect.TypeOf((*MockInstallerInternals)(nil).GetClusterInternal), arg0, arg1)
}

// GetCommonHostInternal mocks base method
func (m *MockInstallerInternals) GetCommonHostInternal(arg0 context.Context, arg1, arg2 string) (*common.Host, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommonHostInternal", arg0, arg1, arg2)
	ret0, _ := ret[0].(*common.Host)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommonHostInternal indicates an expected call of GetCommonHostInternal
func (mr *MockInstallerInternalsMockRecorder) GetCommonHostInternal(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommonHostInternal", reflect.TypeOf((*MockInstallerInternals)(nil).GetCommonHostInternal), arg0, arg1, arg2)
}

// InstallClusterInternal mocks base method
func (m *MockInstallerInternals) InstallClusterInternal(arg0 context.Context, arg1 installer.InstallClusterParams) (*common.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstallClusterInternal", arg0, arg1)
	ret0, _ := ret[0].(*common.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InstallClusterInternal indicates an expected call of InstallClusterInternal
func (mr *MockInstallerInternalsMockRecorder) InstallClusterInternal(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstallClusterInternal", reflect.TypeOf((*MockInstallerInternals)(nil).InstallClusterInternal), arg0, arg1)
}

// RegisterClusterInternal mocks base method
func (m *MockInstallerInternals) RegisterClusterInternal(arg0 context.Context, arg1 *types.NamespacedName, arg2 installer.RegisterClusterParams) (*common.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterClusterInternal", arg0, arg1, arg2)
	ret0, _ := ret[0].(*common.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterClusterInternal indicates an expected call of RegisterClusterInternal
func (mr *MockInstallerInternalsMockRecorder) RegisterClusterInternal(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterClusterInternal", reflect.TypeOf((*MockInstallerInternals)(nil).RegisterClusterInternal), arg0, arg1, arg2)
}

// UpdateClusterInternal mocks base method
func (m *MockInstallerInternals) UpdateClusterInternal(arg0 context.Context, arg1 installer.UpdateClusterParams) (*common.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateClusterInternal", arg0, arg1)
	ret0, _ := ret[0].(*common.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateClusterInternal indicates an expected call of UpdateClusterInternal
func (mr *MockInstallerInternalsMockRecorder) UpdateClusterInternal(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateClusterInternal", reflect.TypeOf((*MockInstallerInternals)(nil).UpdateClusterInternal), arg0, arg1)
}

// UpdateHostApprovedInternal mocks base method
func (m *MockInstallerInternals) UpdateHostApprovedInternal(arg0 context.Context, arg1, arg2 string, arg3 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateHostApprovedInternal", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateHostApprovedInternal indicates an expected call of UpdateHostApprovedInternal
func (mr *MockInstallerInternalsMockRecorder) UpdateHostApprovedInternal(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateHostApprovedInternal", reflect.TypeOf((*MockInstallerInternals)(nil).UpdateHostApprovedInternal), arg0, arg1, arg2, arg3)
}
