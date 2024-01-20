// Code generated by mockery v2.39.1. DO NOT EDIT.

package crawler

import (
	model "github.com/kakurineuin/learn-english-microservices/word-service/pkg/model"
	mock "github.com/stretchr/testify/mock"
)

// MockSpider is an autogenerated mock type for the Spider type
type MockSpider struct {
	mock.Mock
}

type MockSpider_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSpider) EXPECT() *MockSpider_Expecter {
	return &MockSpider_Expecter{mock: &_m.Mock}
}

// FindWordMeaningsFromDictionary provides a mock function with given fields: word
func (_m *MockSpider) FindWordMeaningsFromDictionary(word string) ([]model.WordMeaning, error) {
	ret := _m.Called(word)

	if len(ret) == 0 {
		panic("no return value specified for FindWordMeaningsFromDictionary")
	}

	var r0 []model.WordMeaning
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]model.WordMeaning, error)); ok {
		return rf(word)
	}
	if rf, ok := ret.Get(0).(func(string) []model.WordMeaning); ok {
		r0 = rf(word)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.WordMeaning)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(word)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSpider_FindWordMeaningsFromDictionary_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindWordMeaningsFromDictionary'
type MockSpider_FindWordMeaningsFromDictionary_Call struct {
	*mock.Call
}

// FindWordMeaningsFromDictionary is a helper method to define mock.On call
//   - word string
func (_e *MockSpider_Expecter) FindWordMeaningsFromDictionary(word interface{}) *MockSpider_FindWordMeaningsFromDictionary_Call {
	return &MockSpider_FindWordMeaningsFromDictionary_Call{Call: _e.mock.On("FindWordMeaningsFromDictionary", word)}
}

func (_c *MockSpider_FindWordMeaningsFromDictionary_Call) Run(run func(word string)) *MockSpider_FindWordMeaningsFromDictionary_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockSpider_FindWordMeaningsFromDictionary_Call) Return(_a0 []model.WordMeaning, _a1 error) *MockSpider_FindWordMeaningsFromDictionary_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSpider_FindWordMeaningsFromDictionary_Call) RunAndReturn(run func(string) ([]model.WordMeaning, error)) *MockSpider_FindWordMeaningsFromDictionary_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSpider creates a new instance of MockSpider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSpider(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSpider {
	mock := &MockSpider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
