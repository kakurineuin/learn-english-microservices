// Code generated by mockery v2.39.1. DO NOT EDIT.

package repository

import (
	context "context"

	model "github.com/kakurineuin/learn-english-microservices/web-service/pkg/model"
	mock "github.com/stretchr/testify/mock"
)

// MockDatabaseRepository is an autogenerated mock type for the DatabaseRepository type
type MockDatabaseRepository struct {
	mock.Mock
}

type MockDatabaseRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDatabaseRepository) EXPECT() *MockDatabaseRepository_Expecter {
	return &MockDatabaseRepository_Expecter{mock: &_m.Mock}
}

// ConnectDB provides a mock function with given fields: uri
func (_m *MockDatabaseRepository) ConnectDB(uri string) error {
	ret := _m.Called(uri)

	if len(ret) == 0 {
		panic("no return value specified for ConnectDB")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(uri)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabaseRepository_ConnectDB_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConnectDB'
type MockDatabaseRepository_ConnectDB_Call struct {
	*mock.Call
}

// ConnectDB is a helper method to define mock.On call
//   - uri string
func (_e *MockDatabaseRepository_Expecter) ConnectDB(uri interface{}) *MockDatabaseRepository_ConnectDB_Call {
	return &MockDatabaseRepository_ConnectDB_Call{Call: _e.mock.On("ConnectDB", uri)}
}

func (_c *MockDatabaseRepository_ConnectDB_Call) Run(run func(uri string)) *MockDatabaseRepository_ConnectDB_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockDatabaseRepository_ConnectDB_Call) Return(_a0 error) *MockDatabaseRepository_ConnectDB_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabaseRepository_ConnectDB_Call) RunAndReturn(run func(string) error) *MockDatabaseRepository_ConnectDB_Call {
	_c.Call.Return(run)
	return _c
}

// CreateUser provides a mock function with given fields: ctx, user
func (_m *MockDatabaseRepository) CreateUser(ctx context.Context, user model.User) (string, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.User) (string, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.User) string); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabaseRepository_CreateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUser'
type MockDatabaseRepository_CreateUser_Call struct {
	*mock.Call
}

// CreateUser is a helper method to define mock.On call
//   - ctx context.Context
//   - user model.User
func (_e *MockDatabaseRepository_Expecter) CreateUser(ctx interface{}, user interface{}) *MockDatabaseRepository_CreateUser_Call {
	return &MockDatabaseRepository_CreateUser_Call{Call: _e.mock.On("CreateUser", ctx, user)}
}

func (_c *MockDatabaseRepository_CreateUser_Call) Run(run func(ctx context.Context, user model.User)) *MockDatabaseRepository_CreateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.User))
	})
	return _c
}

func (_c *MockDatabaseRepository_CreateUser_Call) Return(userId string, err error) *MockDatabaseRepository_CreateUser_Call {
	_c.Call.Return(userId, err)
	return _c
}

func (_c *MockDatabaseRepository_CreateUser_Call) RunAndReturn(run func(context.Context, model.User) (string, error)) *MockDatabaseRepository_CreateUser_Call {
	_c.Call.Return(run)
	return _c
}

// DisconnectDB provides a mock function with given fields:
func (_m *MockDatabaseRepository) DisconnectDB() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for DisconnectDB")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabaseRepository_DisconnectDB_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DisconnectDB'
type MockDatabaseRepository_DisconnectDB_Call struct {
	*mock.Call
}

// DisconnectDB is a helper method to define mock.On call
func (_e *MockDatabaseRepository_Expecter) DisconnectDB() *MockDatabaseRepository_DisconnectDB_Call {
	return &MockDatabaseRepository_DisconnectDB_Call{Call: _e.mock.On("DisconnectDB")}
}

func (_c *MockDatabaseRepository_DisconnectDB_Call) Run(run func()) *MockDatabaseRepository_DisconnectDB_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockDatabaseRepository_DisconnectDB_Call) Return(_a0 error) *MockDatabaseRepository_DisconnectDB_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabaseRepository_DisconnectDB_Call) RunAndReturn(run func() error) *MockDatabaseRepository_DisconnectDB_Call {
	_c.Call.Return(run)
	return _c
}

// GetAdminUser provides a mock function with given fields: ctx
func (_m *MockDatabaseRepository) GetAdminUser(ctx context.Context) (*model.User, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAdminUser")
	}

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*model.User, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *model.User); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabaseRepository_GetAdminUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAdminUser'
type MockDatabaseRepository_GetAdminUser_Call struct {
	*mock.Call
}

// GetAdminUser is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockDatabaseRepository_Expecter) GetAdminUser(ctx interface{}) *MockDatabaseRepository_GetAdminUser_Call {
	return &MockDatabaseRepository_GetAdminUser_Call{Call: _e.mock.On("GetAdminUser", ctx)}
}

func (_c *MockDatabaseRepository_GetAdminUser_Call) Run(run func(ctx context.Context)) *MockDatabaseRepository_GetAdminUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockDatabaseRepository_GetAdminUser_Call) Return(user *model.User, err error) *MockDatabaseRepository_GetAdminUser_Call {
	_c.Call.Return(user, err)
	return _c
}

func (_c *MockDatabaseRepository_GetAdminUser_Call) RunAndReturn(run func(context.Context) (*model.User, error)) *MockDatabaseRepository_GetAdminUser_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserById provides a mock function with given fields: ctx, userId
func (_m *MockDatabaseRepository) GetUserById(ctx context.Context, userId string) (*model.User, error) {
	ret := _m.Called(ctx, userId)

	if len(ret) == 0 {
		panic("no return value specified for GetUserById")
	}

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.User, error)); ok {
		return rf(ctx, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.User); ok {
		r0 = rf(ctx, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabaseRepository_GetUserById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserById'
type MockDatabaseRepository_GetUserById_Call struct {
	*mock.Call
}

// GetUserById is a helper method to define mock.On call
//   - ctx context.Context
//   - userId string
func (_e *MockDatabaseRepository_Expecter) GetUserById(ctx interface{}, userId interface{}) *MockDatabaseRepository_GetUserById_Call {
	return &MockDatabaseRepository_GetUserById_Call{Call: _e.mock.On("GetUserById", ctx, userId)}
}

func (_c *MockDatabaseRepository_GetUserById_Call) Run(run func(ctx context.Context, userId string)) *MockDatabaseRepository_GetUserById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockDatabaseRepository_GetUserById_Call) Return(user *model.User, err error) *MockDatabaseRepository_GetUserById_Call {
	_c.Call.Return(user, err)
	return _c
}

func (_c *MockDatabaseRepository_GetUserById_Call) RunAndReturn(run func(context.Context, string) (*model.User, error)) *MockDatabaseRepository_GetUserById_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserByUsername provides a mock function with given fields: ctx, username
func (_m *MockDatabaseRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	ret := _m.Called(ctx, username)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByUsername")
	}

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.User, error)); ok {
		return rf(ctx, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.User); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabaseRepository_GetUserByUsername_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserByUsername'
type MockDatabaseRepository_GetUserByUsername_Call struct {
	*mock.Call
}

// GetUserByUsername is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
func (_e *MockDatabaseRepository_Expecter) GetUserByUsername(ctx interface{}, username interface{}) *MockDatabaseRepository_GetUserByUsername_Call {
	return &MockDatabaseRepository_GetUserByUsername_Call{Call: _e.mock.On("GetUserByUsername", ctx, username)}
}

func (_c *MockDatabaseRepository_GetUserByUsername_Call) Run(run func(ctx context.Context, username string)) *MockDatabaseRepository_GetUserByUsername_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockDatabaseRepository_GetUserByUsername_Call) Return(user *model.User, err error) *MockDatabaseRepository_GetUserByUsername_Call {
	_c.Call.Return(user, err)
	return _c
}

func (_c *MockDatabaseRepository_GetUserByUsername_Call) RunAndReturn(run func(context.Context, string) (*model.User, error)) *MockDatabaseRepository_GetUserByUsername_Call {
	_c.Call.Return(run)
	return _c
}

// WithTransaction provides a mock function with given fields: transactoinFunc
func (_m *MockDatabaseRepository) WithTransaction(transactoinFunc transactionFunc) (interface{}, error) {
	ret := _m.Called(transactoinFunc)

	if len(ret) == 0 {
		panic("no return value specified for WithTransaction")
	}

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(transactionFunc) (interface{}, error)); ok {
		return rf(transactoinFunc)
	}
	if rf, ok := ret.Get(0).(func(transactionFunc) interface{}); ok {
		r0 = rf(transactoinFunc)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(transactionFunc) error); ok {
		r1 = rf(transactoinFunc)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabaseRepository_WithTransaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithTransaction'
type MockDatabaseRepository_WithTransaction_Call struct {
	*mock.Call
}

// WithTransaction is a helper method to define mock.On call
//   - transactoinFunc transactionFunc
func (_e *MockDatabaseRepository_Expecter) WithTransaction(transactoinFunc interface{}) *MockDatabaseRepository_WithTransaction_Call {
	return &MockDatabaseRepository_WithTransaction_Call{Call: _e.mock.On("WithTransaction", transactoinFunc)}
}

func (_c *MockDatabaseRepository_WithTransaction_Call) Run(run func(transactoinFunc transactionFunc)) *MockDatabaseRepository_WithTransaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(transactionFunc))
	})
	return _c
}

func (_c *MockDatabaseRepository_WithTransaction_Call) Return(_a0 interface{}, _a1 error) *MockDatabaseRepository_WithTransaction_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabaseRepository_WithTransaction_Call) RunAndReturn(run func(transactionFunc) (interface{}, error)) *MockDatabaseRepository_WithTransaction_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDatabaseRepository creates a new instance of MockDatabaseRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDatabaseRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDatabaseRepository {
	mock := &MockDatabaseRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
