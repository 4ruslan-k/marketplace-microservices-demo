// Code generated by MockGen. DO NOT EDIT.
// Source: internal/services/user.go

// Package mock_applicationservices is a generated GoMock package.
package mock_applicationservices

import (
	user "authentication/internal/domain/entities/user"
	domainservices "authentication/internal/domain/services"
	dto "authentication/internal/services/dto"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserApplicationService is a mock of UserApplicationService interface.
type MockUserApplicationService struct {
	ctrl     *gomock.Controller
	recorder *MockUserApplicationServiceMockRecorder
}

// MockUserApplicationServiceMockRecorder is the mock recorder for MockUserApplicationService.
type MockUserApplicationServiceMockRecorder struct {
	mock *MockUserApplicationService
}

// NewMockUserApplicationService creates a new mock instance.
func NewMockUserApplicationService(ctrl *gomock.Controller) *MockUserApplicationService {
	mock := &MockUserApplicationService{ctrl: ctrl}
	mock.recorder = &MockUserApplicationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserApplicationService) EXPECT() *MockUserApplicationServiceMockRecorder {
	return m.recorder
}

// ChangeCurrentPassword mocks base method.
func (m *MockUserApplicationService) ChangeCurrentPassword(ctx context.Context, socialAccount dto.ChangeCurrentPasswordInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeCurrentPassword", ctx, socialAccount)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeCurrentPassword indicates an expected call of ChangeCurrentPassword.
func (mr *MockUserApplicationServiceMockRecorder) ChangeCurrentPassword(ctx, socialAccount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeCurrentPassword", reflect.TypeOf((*MockUserApplicationService)(nil).ChangeCurrentPassword), ctx, socialAccount)
}

// CreateUser mocks base method.
func (m *MockUserApplicationService) CreateUser(ctx context.Context, createUser user.CreateUserParams) (*dto.UserOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, createUser)
	ret0, _ := ret[0].(*dto.UserOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserApplicationServiceMockRecorder) CreateUser(ctx, createUser interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserApplicationService)(nil).CreateUser), ctx, createUser)
}

// DeleteUser mocks base method.
func (m *MockUserApplicationService) DeleteUser(ctx context.Context, deleteUser dto.DeleteUserInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, deleteUser)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockUserApplicationServiceMockRecorder) DeleteUser(ctx, deleteUser interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserApplicationService)(nil).DeleteUser), ctx, deleteUser)
}

// DisableTotp mocks base method.
func (m *MockUserApplicationService) DisableTotp(ctx context.Context, userID, otp string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisableTotp", ctx, userID, otp)
	ret0, _ := ret[0].(error)
	return ret0
}

// DisableTotp indicates an expected call of DisableTotp.
func (mr *MockUserApplicationServiceMockRecorder) DisableTotp(ctx, userID, otp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisableTotp", reflect.TypeOf((*MockUserApplicationService)(nil).DisableTotp), ctx, userID, otp)
}

// EnableTotp mocks base method.
func (m *MockUserApplicationService) EnableTotp(ctx context.Context, userID, otp string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnableTotp", ctx, userID, otp)
	ret0, _ := ret[0].(error)
	return ret0
}

// EnableTotp indicates an expected call of EnableTotp.
func (mr *MockUserApplicationServiceMockRecorder) EnableTotp(ctx, userID, otp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableTotp", reflect.TypeOf((*MockUserApplicationService)(nil).EnableTotp), ctx, userID, otp)
}

// GenerateTotpSetup mocks base method.
func (m *MockUserApplicationService) GenerateTotpSetup(ctx context.Context, userID string) (domainservices.TotpSetupInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateTotpSetup", ctx, userID)
	ret0, _ := ret[0].(domainservices.TotpSetupInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateTotpSetup indicates an expected call of GenerateTotpSetup.
func (mr *MockUserApplicationServiceMockRecorder) GenerateTotpSetup(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateTotpSetup", reflect.TypeOf((*MockUserApplicationService)(nil).GenerateTotpSetup), ctx, userID)
}

// GetUserByID mocks base method.
func (m *MockUserApplicationService) GetUserByID(ctx context.Context, ID string) (*dto.UserOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, ID)
	ret0, _ := ret[0].(*dto.UserOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockUserApplicationServiceMockRecorder) GetUserByID(ctx, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockUserApplicationService)(nil).GetUserByID), ctx, ID)
}

// LoginWithEmailAndPassword mocks base method.
func (m *MockUserApplicationService) LoginWithEmailAndPassword(ctx context.Context, email, password string) (dto.LoginOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginWithEmailAndPassword", ctx, email, password)
	ret0, _ := ret[0].(dto.LoginOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginWithEmailAndPassword indicates an expected call of LoginWithEmailAndPassword.
func (mr *MockUserApplicationServiceMockRecorder) LoginWithEmailAndPassword(ctx, email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginWithEmailAndPassword", reflect.TypeOf((*MockUserApplicationService)(nil).LoginWithEmailAndPassword), ctx, email, password)
}

// LoginWithTotpCode mocks base method.
func (m *MockUserApplicationService) LoginWithTotpCode(ctx context.Context, passwordVerificationTokenID, code string) (dto.UserOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginWithTotpCode", ctx, passwordVerificationTokenID, code)
	ret0, _ := ret[0].(dto.UserOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginWithTotpCode indicates an expected call of LoginWithTotpCode.
func (mr *MockUserApplicationServiceMockRecorder) LoginWithTotpCode(ctx, passwordVerificationTokenID, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginWithTotpCode", reflect.TypeOf((*MockUserApplicationService)(nil).LoginWithTotpCode), ctx, passwordVerificationTokenID, code)
}

// SocialLogin mocks base method.
func (m *MockUserApplicationService) SocialLogin(ctx context.Context, socialAccount dto.SocialLoginInput) (*dto.LoginOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SocialLogin", ctx, socialAccount)
	ret0, _ := ret[0].(*dto.LoginOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SocialLogin indicates an expected call of SocialLogin.
func (mr *MockUserApplicationServiceMockRecorder) SocialLogin(ctx, socialAccount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SocialLogin", reflect.TypeOf((*MockUserApplicationService)(nil).SocialLogin), ctx, socialAccount)
}

// UpdateUser mocks base method.
func (m *MockUserApplicationService) UpdateUser(ctx context.Context, updateUser dto.UpdateUserInput) (*dto.UserOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, updateUser)
	ret0, _ := ret[0].(*dto.UserOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserApplicationServiceMockRecorder) UpdateUser(ctx, updateUser interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserApplicationService)(nil).UpdateUser), ctx, updateUser)
}