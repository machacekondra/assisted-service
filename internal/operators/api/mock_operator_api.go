// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/openshift/assisted-service/internal/operators/api (interfaces: Operator)

// Package api is a generated GoMock package.
package api

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	common "github.com/openshift/assisted-service/internal/common"
	models "github.com/openshift/assisted-service/models"
	reflect "reflect"
)

// MockOperator is a mock of Operator interface
type MockOperator struct {
	ctrl     *gomock.Controller
	recorder *MockOperatorMockRecorder
}

// MockOperatorMockRecorder is the mock recorder for MockOperator
type MockOperatorMockRecorder struct {
	mock *MockOperator
}

// NewMockOperator creates a new mock instance
func NewMockOperator(ctrl *gomock.Controller) *MockOperator {
	mock := &MockOperator{ctrl: ctrl}
	mock.recorder = &MockOperatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOperator) EXPECT() *MockOperatorMockRecorder {
	return m.recorder
}

// GenerateManifests mocks base method
func (m *MockOperator) GenerateManifests(arg0 *common.Cluster) (map[string][]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateManifests", arg0)
	ret0, _ := ret[0].(map[string][]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateManifests indicates an expected call of GenerateManifests
func (mr *MockOperatorMockRecorder) GenerateManifests(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateManifests", reflect.TypeOf((*MockOperator)(nil).GenerateManifests), arg0)
}

// GetClusterValidationID mocks base method
func (m *MockOperator) GetClusterValidationID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClusterValidationID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetClusterValidationID indicates an expected call of GetClusterValidationID
func (mr *MockOperatorMockRecorder) GetClusterValidationID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClusterValidationID", reflect.TypeOf((*MockOperator)(nil).GetClusterValidationID))
}

// GetDependencies mocks base method
func (m *MockOperator) GetDependencies() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDependencies")
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetDependencies indicates an expected call of GetDependencies
func (mr *MockOperatorMockRecorder) GetDependencies() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDependencies", reflect.TypeOf((*MockOperator)(nil).GetDependencies))
}

// GetHostRequirements mocks base method
func (m *MockOperator) GetHostRequirements(arg0 context.Context, arg1 *common.Cluster, arg2 *models.Host) (*models.ClusterHostRequirementsDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHostRequirements", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.ClusterHostRequirementsDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHostRequirements indicates an expected call of GetHostRequirements
func (mr *MockOperatorMockRecorder) GetHostRequirements(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHostRequirements", reflect.TypeOf((*MockOperator)(nil).GetHostRequirements), arg0, arg1, arg2)
}

// GetHostValidationID mocks base method
func (m *MockOperator) GetHostValidationID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHostValidationID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetHostValidationID indicates an expected call of GetHostValidationID
func (mr *MockOperatorMockRecorder) GetHostValidationID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHostValidationID", reflect.TypeOf((*MockOperator)(nil).GetHostValidationID))
}

// GetMonitoredOperator mocks base method
func (m *MockOperator) GetMonitoredOperator() *models.MonitoredOperator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMonitoredOperator")
	ret0, _ := ret[0].(*models.MonitoredOperator)
	return ret0
}

// GetMonitoredOperator indicates an expected call of GetMonitoredOperator
func (mr *MockOperatorMockRecorder) GetMonitoredOperator() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMonitoredOperator", reflect.TypeOf((*MockOperator)(nil).GetMonitoredOperator))
}

// GetName mocks base method
func (m *MockOperator) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName
func (mr *MockOperatorMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockOperator)(nil).GetName))
}

// GetPreflightRequirements mocks base method
func (m *MockOperator) GetPreflightRequirements(arg0 context.Context, arg1 *common.Cluster) (*models.OperatorHardwareRequirements, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPreflightRequirements", arg0, arg1)
	ret0, _ := ret[0].(*models.OperatorHardwareRequirements)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPreflightRequirements indicates an expected call of GetPreflightRequirements
func (mr *MockOperatorMockRecorder) GetPreflightRequirements(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPreflightRequirements", reflect.TypeOf((*MockOperator)(nil).GetPreflightRequirements), arg0, arg1)
}

// GetProperties mocks base method
func (m *MockOperator) GetProperties() models.OperatorProperties {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProperties")
	ret0, _ := ret[0].(models.OperatorProperties)
	return ret0
}

// GetProperties indicates an expected call of GetProperties
func (mr *MockOperatorMockRecorder) GetProperties() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProperties", reflect.TypeOf((*MockOperator)(nil).GetProperties))
}

// ValidateCluster mocks base method
func (m *MockOperator) ValidateCluster(arg0 context.Context, arg1 *common.Cluster) (ValidationResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateCluster", arg0, arg1)
	ret0, _ := ret[0].(ValidationResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateCluster indicates an expected call of ValidateCluster
func (mr *MockOperatorMockRecorder) ValidateCluster(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateCluster", reflect.TypeOf((*MockOperator)(nil).ValidateCluster), arg0, arg1)
}

// ValidateHost mocks base method
func (m *MockOperator) ValidateHost(arg0 context.Context, arg1 *common.Cluster, arg2 *models.Host) (ValidationResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateHost", arg0, arg1, arg2)
	ret0, _ := ret[0].(ValidationResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateHost indicates an expected call of ValidateHost
func (mr *MockOperatorMockRecorder) ValidateHost(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateHost", reflect.TypeOf((*MockOperator)(nil).ValidateHost), arg0, arg1, arg2)
}
