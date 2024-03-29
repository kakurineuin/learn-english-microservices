// Code generated by mockery v2.40.1. DO NOT EDIT.

package repository

import (
	context "context"

	model "github.com/kakurineuin/learn-english-microservices/word-service/pkg/model"
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

// ConnectDB provides a mock function with given fields: ctx, uri
func (_m *MockDatabaseRepository) ConnectDB(ctx context.Context, uri string) error {
	ret := _m.Called(ctx, uri)

	if len(ret) == 0 {
		panic("no return value specified for ConnectDB")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, uri)
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
//   - ctx context.Context
//   - uri string
func (_e *MockDatabaseRepository_Expecter) ConnectDB(ctx interface{}, uri interface{}) *MockDatabaseRepository_ConnectDB_Call {
	return &MockDatabaseRepository_ConnectDB_Call{Call: _e.mock.On("ConnectDB", ctx, uri)}
}

func (_c *MockDatabaseRepository_ConnectDB_Call) Run(run func(ctx context.Context, uri string)) *MockDatabaseRepository_ConnectDB_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockDatabaseRepository_ConnectDB_Call) Return(_a0 error) *MockDatabaseRepository_ConnectDB_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabaseRepository_ConnectDB_Call) RunAndReturn(run func(context.Context, string) error) *MockDatabaseRepository_ConnectDB_Call {
	_c.Call.Return(run)
	return _c
}

// CountFavoriteWordMeaningsByUserIdAndWord provides a mock function with given fields: ctx, userId, word
func (_m *MockDatabaseRepository) CountFavoriteWordMeaningsByUserIdAndWord(ctx context.Context, userId string, word string) (int32, error) {
	ret := _m.Called(ctx, userId, word)

	if len(ret) == 0 {
		panic("no return value specified for CountFavoriteWordMeaningsByUserIdAndWord")
	}

	var r0 int32
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (int32, error)); ok {
		return rf(ctx, userId, word)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) int32); ok {
		r0 = rf(ctx, userId, word)
	} else {
		r0 = ret.Get(0).(int32)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, userId, word)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabaseRepository_CountFavoriteWordMeaningsByUserIdAndWord_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CountFavoriteWordMeaningsByUserIdAndWord'
type MockDatabaseRepository_CountFavoriteWordMeaningsByUserIdAndWord_Call struct {
	*mock.Call
}

// CountFavoriteWordMeaningsByUserIdAndWord is a helper method to define mock.On call
//   - ctx context.Context
//   - userId string
//   - word string
func (_e *MockDatabaseRepository_Expecter) CountFavoriteWordMeaningsByUserIdAndWord(ctx interface{}, userId interface{}, word interface{}) *MockDatabaseRepository_CountFavoriteWordMeaningsByUserIdAndWord_Call {
	return &MockDatabaseRepository_CountFavoriteWordMeaningsByUserIdAndWord_Call{Call: _e.mock.On("CountFavoriteWordMeaningsByUserIdAndWord", ctx, userId, word)}
}

func (_c *MockDatabaseRepository_CountFavoriteWordMeaningsByUserIdAndWord_Call) Run(run func(ctx context.Context, userId string, word string)) *MockDatabaseRepository_CountFavoriteWordMeaningsByUserIdAndWord_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockDatabaseRepository_CountFavoriteWordMeaningsByUserIdAndWord_Call) Return(count int32, err error) *MockDatabaseRepository_CountFavoriteWordMeaningsByUserIdAndWord_Call {
	_c.Call.Return(count, err)
	return _c
}

func (_c *MockDatabaseRepository_CountFavoriteWordMeaningsByUserIdAndWord_Call) RunAndReturn(run func(context.Context, string, string) (int32, error)) *MockDatabaseRepository_CountFavoriteWordMeaningsByUserIdAndWord_Call {
	_c.Call.Return(run)
	return _c
}

// CreateFavoriteWordMeaning provides a mock function with given fields: ctx, userId, wordMeaningId
func (_m *MockDatabaseRepository) CreateFavoriteWordMeaning(ctx context.Context, userId string, wordMeaningId string) (string, error) {
	ret := _m.Called(ctx, userId, wordMeaningId)

	if len(ret) == 0 {
		panic("no return value specified for CreateFavoriteWordMeaning")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (string, error)); ok {
		return rf(ctx, userId, wordMeaningId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, userId, wordMeaningId)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, userId, wordMeaningId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabaseRepository_CreateFavoriteWordMeaning_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateFavoriteWordMeaning'
type MockDatabaseRepository_CreateFavoriteWordMeaning_Call struct {
	*mock.Call
}

// CreateFavoriteWordMeaning is a helper method to define mock.On call
//   - ctx context.Context
//   - userId string
//   - wordMeaningId string
func (_e *MockDatabaseRepository_Expecter) CreateFavoriteWordMeaning(ctx interface{}, userId interface{}, wordMeaningId interface{}) *MockDatabaseRepository_CreateFavoriteWordMeaning_Call {
	return &MockDatabaseRepository_CreateFavoriteWordMeaning_Call{Call: _e.mock.On("CreateFavoriteWordMeaning", ctx, userId, wordMeaningId)}
}

func (_c *MockDatabaseRepository_CreateFavoriteWordMeaning_Call) Run(run func(ctx context.Context, userId string, wordMeaningId string)) *MockDatabaseRepository_CreateFavoriteWordMeaning_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockDatabaseRepository_CreateFavoriteWordMeaning_Call) Return(favoriteWordMeaningId string, err error) *MockDatabaseRepository_CreateFavoriteWordMeaning_Call {
	_c.Call.Return(favoriteWordMeaningId, err)
	return _c
}

func (_c *MockDatabaseRepository_CreateFavoriteWordMeaning_Call) RunAndReturn(run func(context.Context, string, string) (string, error)) *MockDatabaseRepository_CreateFavoriteWordMeaning_Call {
	_c.Call.Return(run)
	return _c
}

// CreateWordMeanings provides a mock function with given fields: ctx, wordMeanings
func (_m *MockDatabaseRepository) CreateWordMeanings(ctx context.Context, wordMeanings []model.WordMeaning) ([]string, error) {
	ret := _m.Called(ctx, wordMeanings)

	if len(ret) == 0 {
		panic("no return value specified for CreateWordMeanings")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []model.WordMeaning) ([]string, error)); ok {
		return rf(ctx, wordMeanings)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []model.WordMeaning) []string); ok {
		r0 = rf(ctx, wordMeanings)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []model.WordMeaning) error); ok {
		r1 = rf(ctx, wordMeanings)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabaseRepository_CreateWordMeanings_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateWordMeanings'
type MockDatabaseRepository_CreateWordMeanings_Call struct {
	*mock.Call
}

// CreateWordMeanings is a helper method to define mock.On call
//   - ctx context.Context
//   - wordMeanings []model.WordMeaning
func (_e *MockDatabaseRepository_Expecter) CreateWordMeanings(ctx interface{}, wordMeanings interface{}) *MockDatabaseRepository_CreateWordMeanings_Call {
	return &MockDatabaseRepository_CreateWordMeanings_Call{Call: _e.mock.On("CreateWordMeanings", ctx, wordMeanings)}
}

func (_c *MockDatabaseRepository_CreateWordMeanings_Call) Run(run func(ctx context.Context, wordMeanings []model.WordMeaning)) *MockDatabaseRepository_CreateWordMeanings_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]model.WordMeaning))
	})
	return _c
}

func (_c *MockDatabaseRepository_CreateWordMeanings_Call) Return(wordMeaningIds []string, err error) *MockDatabaseRepository_CreateWordMeanings_Call {
	_c.Call.Return(wordMeaningIds, err)
	return _c
}

func (_c *MockDatabaseRepository_CreateWordMeanings_Call) RunAndReturn(run func(context.Context, []model.WordMeaning) ([]string, error)) *MockDatabaseRepository_CreateWordMeanings_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteFavoriteWordMeaningById provides a mock function with given fields: ctx, favoriteWordMeaningId
func (_m *MockDatabaseRepository) DeleteFavoriteWordMeaningById(ctx context.Context, favoriteWordMeaningId string) (int32, error) {
	ret := _m.Called(ctx, favoriteWordMeaningId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteFavoriteWordMeaningById")
	}

	var r0 int32
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (int32, error)); ok {
		return rf(ctx, favoriteWordMeaningId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) int32); ok {
		r0 = rf(ctx, favoriteWordMeaningId)
	} else {
		r0 = ret.Get(0).(int32)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, favoriteWordMeaningId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabaseRepository_DeleteFavoriteWordMeaningById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteFavoriteWordMeaningById'
type MockDatabaseRepository_DeleteFavoriteWordMeaningById_Call struct {
	*mock.Call
}

// DeleteFavoriteWordMeaningById is a helper method to define mock.On call
//   - ctx context.Context
//   - favoriteWordMeaningId string
func (_e *MockDatabaseRepository_Expecter) DeleteFavoriteWordMeaningById(ctx interface{}, favoriteWordMeaningId interface{}) *MockDatabaseRepository_DeleteFavoriteWordMeaningById_Call {
	return &MockDatabaseRepository_DeleteFavoriteWordMeaningById_Call{Call: _e.mock.On("DeleteFavoriteWordMeaningById", ctx, favoriteWordMeaningId)}
}

func (_c *MockDatabaseRepository_DeleteFavoriteWordMeaningById_Call) Run(run func(ctx context.Context, favoriteWordMeaningId string)) *MockDatabaseRepository_DeleteFavoriteWordMeaningById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockDatabaseRepository_DeleteFavoriteWordMeaningById_Call) Return(deletedCount int32, err error) *MockDatabaseRepository_DeleteFavoriteWordMeaningById_Call {
	_c.Call.Return(deletedCount, err)
	return _c
}

func (_c *MockDatabaseRepository_DeleteFavoriteWordMeaningById_Call) RunAndReturn(run func(context.Context, string) (int32, error)) *MockDatabaseRepository_DeleteFavoriteWordMeaningById_Call {
	_c.Call.Return(run)
	return _c
}

// DisconnectDB provides a mock function with given fields: ctx
func (_m *MockDatabaseRepository) DisconnectDB(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for DisconnectDB")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
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
//   - ctx context.Context
func (_e *MockDatabaseRepository_Expecter) DisconnectDB(ctx interface{}) *MockDatabaseRepository_DisconnectDB_Call {
	return &MockDatabaseRepository_DisconnectDB_Call{Call: _e.mock.On("DisconnectDB", ctx)}
}

func (_c *MockDatabaseRepository_DisconnectDB_Call) Run(run func(ctx context.Context)) *MockDatabaseRepository_DisconnectDB_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockDatabaseRepository_DisconnectDB_Call) Return(_a0 error) *MockDatabaseRepository_DisconnectDB_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabaseRepository_DisconnectDB_Call) RunAndReturn(run func(context.Context) error) *MockDatabaseRepository_DisconnectDB_Call {
	_c.Call.Return(run)
	return _c
}

// FindFavoriteWordMeaningsByUserIdAndWord provides a mock function with given fields: ctx, userId, word, skip, limit
func (_m *MockDatabaseRepository) FindFavoriteWordMeaningsByUserIdAndWord(ctx context.Context, userId string, word string, skip int32, limit int32) ([]model.WordMeaning, error) {
	ret := _m.Called(ctx, userId, word, skip, limit)

	if len(ret) == 0 {
		panic("no return value specified for FindFavoriteWordMeaningsByUserIdAndWord")
	}

	var r0 []model.WordMeaning
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int32, int32) ([]model.WordMeaning, error)); ok {
		return rf(ctx, userId, word, skip, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int32, int32) []model.WordMeaning); ok {
		r0 = rf(ctx, userId, word, skip, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.WordMeaning)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, int32, int32) error); ok {
		r1 = rf(ctx, userId, word, skip, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabaseRepository_FindFavoriteWordMeaningsByUserIdAndWord_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindFavoriteWordMeaningsByUserIdAndWord'
type MockDatabaseRepository_FindFavoriteWordMeaningsByUserIdAndWord_Call struct {
	*mock.Call
}

// FindFavoriteWordMeaningsByUserIdAndWord is a helper method to define mock.On call
//   - ctx context.Context
//   - userId string
//   - word string
//   - skip int32
//   - limit int32
func (_e *MockDatabaseRepository_Expecter) FindFavoriteWordMeaningsByUserIdAndWord(ctx interface{}, userId interface{}, word interface{}, skip interface{}, limit interface{}) *MockDatabaseRepository_FindFavoriteWordMeaningsByUserIdAndWord_Call {
	return &MockDatabaseRepository_FindFavoriteWordMeaningsByUserIdAndWord_Call{Call: _e.mock.On("FindFavoriteWordMeaningsByUserIdAndWord", ctx, userId, word, skip, limit)}
}

func (_c *MockDatabaseRepository_FindFavoriteWordMeaningsByUserIdAndWord_Call) Run(run func(ctx context.Context, userId string, word string, skip int32, limit int32)) *MockDatabaseRepository_FindFavoriteWordMeaningsByUserIdAndWord_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(int32), args[4].(int32))
	})
	return _c
}

func (_c *MockDatabaseRepository_FindFavoriteWordMeaningsByUserIdAndWord_Call) Return(wordMeanings []model.WordMeaning, err error) *MockDatabaseRepository_FindFavoriteWordMeaningsByUserIdAndWord_Call {
	_c.Call.Return(wordMeanings, err)
	return _c
}

func (_c *MockDatabaseRepository_FindFavoriteWordMeaningsByUserIdAndWord_Call) RunAndReturn(run func(context.Context, string, string, int32, int32) ([]model.WordMeaning, error)) *MockDatabaseRepository_FindFavoriteWordMeaningsByUserIdAndWord_Call {
	_c.Call.Return(run)
	return _c
}

// FindWordMeaningsByWordAndUserId provides a mock function with given fields: ctx, word, userId
func (_m *MockDatabaseRepository) FindWordMeaningsByWordAndUserId(ctx context.Context, word string, userId string) ([]model.WordMeaning, error) {
	ret := _m.Called(ctx, word, userId)

	if len(ret) == 0 {
		panic("no return value specified for FindWordMeaningsByWordAndUserId")
	}

	var r0 []model.WordMeaning
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) ([]model.WordMeaning, error)); ok {
		return rf(ctx, word, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []model.WordMeaning); ok {
		r0 = rf(ctx, word, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.WordMeaning)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, word, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabaseRepository_FindWordMeaningsByWordAndUserId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindWordMeaningsByWordAndUserId'
type MockDatabaseRepository_FindWordMeaningsByWordAndUserId_Call struct {
	*mock.Call
}

// FindWordMeaningsByWordAndUserId is a helper method to define mock.On call
//   - ctx context.Context
//   - word string
//   - userId string
func (_e *MockDatabaseRepository_Expecter) FindWordMeaningsByWordAndUserId(ctx interface{}, word interface{}, userId interface{}) *MockDatabaseRepository_FindWordMeaningsByWordAndUserId_Call {
	return &MockDatabaseRepository_FindWordMeaningsByWordAndUserId_Call{Call: _e.mock.On("FindWordMeaningsByWordAndUserId", ctx, word, userId)}
}

func (_c *MockDatabaseRepository_FindWordMeaningsByWordAndUserId_Call) Run(run func(ctx context.Context, word string, userId string)) *MockDatabaseRepository_FindWordMeaningsByWordAndUserId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockDatabaseRepository_FindWordMeaningsByWordAndUserId_Call) Return(wordMeanings []model.WordMeaning, err error) *MockDatabaseRepository_FindWordMeaningsByWordAndUserId_Call {
	_c.Call.Return(wordMeanings, err)
	return _c
}

func (_c *MockDatabaseRepository_FindWordMeaningsByWordAndUserId_Call) RunAndReturn(run func(context.Context, string, string) ([]model.WordMeaning, error)) *MockDatabaseRepository_FindWordMeaningsByWordAndUserId_Call {
	_c.Call.Return(run)
	return _c
}

// GetFavoriteWordMeaningById provides a mock function with given fields: ctx, favoriteWordMeaningId
func (_m *MockDatabaseRepository) GetFavoriteWordMeaningById(ctx context.Context, favoriteWordMeaningId string) (*model.FavoriteWordMeaning, error) {
	ret := _m.Called(ctx, favoriteWordMeaningId)

	if len(ret) == 0 {
		panic("no return value specified for GetFavoriteWordMeaningById")
	}

	var r0 *model.FavoriteWordMeaning
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.FavoriteWordMeaning, error)); ok {
		return rf(ctx, favoriteWordMeaningId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.FavoriteWordMeaning); ok {
		r0 = rf(ctx, favoriteWordMeaningId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FavoriteWordMeaning)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, favoriteWordMeaningId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabaseRepository_GetFavoriteWordMeaningById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetFavoriteWordMeaningById'
type MockDatabaseRepository_GetFavoriteWordMeaningById_Call struct {
	*mock.Call
}

// GetFavoriteWordMeaningById is a helper method to define mock.On call
//   - ctx context.Context
//   - favoriteWordMeaningId string
func (_e *MockDatabaseRepository_Expecter) GetFavoriteWordMeaningById(ctx interface{}, favoriteWordMeaningId interface{}) *MockDatabaseRepository_GetFavoriteWordMeaningById_Call {
	return &MockDatabaseRepository_GetFavoriteWordMeaningById_Call{Call: _e.mock.On("GetFavoriteWordMeaningById", ctx, favoriteWordMeaningId)}
}

func (_c *MockDatabaseRepository_GetFavoriteWordMeaningById_Call) Run(run func(ctx context.Context, favoriteWordMeaningId string)) *MockDatabaseRepository_GetFavoriteWordMeaningById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockDatabaseRepository_GetFavoriteWordMeaningById_Call) Return(favoriteWordMeaning *model.FavoriteWordMeaning, err error) *MockDatabaseRepository_GetFavoriteWordMeaningById_Call {
	_c.Call.Return(favoriteWordMeaning, err)
	return _c
}

func (_c *MockDatabaseRepository_GetFavoriteWordMeaningById_Call) RunAndReturn(run func(context.Context, string) (*model.FavoriteWordMeaning, error)) *MockDatabaseRepository_GetFavoriteWordMeaningById_Call {
	_c.Call.Return(run)
	return _c
}

// WithTransaction provides a mock function with given fields: ctx, transactoinFunc
func (_m *MockDatabaseRepository) WithTransaction(ctx context.Context, transactoinFunc transactionFunc) (interface{}, error) {
	ret := _m.Called(ctx, transactoinFunc)

	if len(ret) == 0 {
		panic("no return value specified for WithTransaction")
	}

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, transactionFunc) (interface{}, error)); ok {
		return rf(ctx, transactoinFunc)
	}
	if rf, ok := ret.Get(0).(func(context.Context, transactionFunc) interface{}); ok {
		r0 = rf(ctx, transactoinFunc)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, transactionFunc) error); ok {
		r1 = rf(ctx, transactoinFunc)
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
//   - ctx context.Context
//   - transactoinFunc transactionFunc
func (_e *MockDatabaseRepository_Expecter) WithTransaction(ctx interface{}, transactoinFunc interface{}) *MockDatabaseRepository_WithTransaction_Call {
	return &MockDatabaseRepository_WithTransaction_Call{Call: _e.mock.On("WithTransaction", ctx, transactoinFunc)}
}

func (_c *MockDatabaseRepository_WithTransaction_Call) Run(run func(ctx context.Context, transactoinFunc transactionFunc)) *MockDatabaseRepository_WithTransaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(transactionFunc))
	})
	return _c
}

func (_c *MockDatabaseRepository_WithTransaction_Call) Return(_a0 interface{}, _a1 error) *MockDatabaseRepository_WithTransaction_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabaseRepository_WithTransaction_Call) RunAndReturn(run func(context.Context, transactionFunc) (interface{}, error)) *MockDatabaseRepository_WithTransaction_Call {
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
