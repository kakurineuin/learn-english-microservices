// Code generated by mockery v2.39.1. DO NOT EDIT.

package examservice

import (
	pb "github.com/kakurineuin/learn-english-microservices/web-service/pb"
	mock "github.com/stretchr/testify/mock"
)

// MockExamService is an autogenerated mock type for the ExamService type
type MockExamService struct {
	mock.Mock
}

type MockExamService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockExamService) EXPECT() *MockExamService_Expecter {
	return &MockExamService_Expecter{mock: &_m.Mock}
}

// Connect provides a mock function with given fields:
func (_m *MockExamService) Connect() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Connect")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockExamService_Connect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Connect'
type MockExamService_Connect_Call struct {
	*mock.Call
}

// Connect is a helper method to define mock.On call
func (_e *MockExamService_Expecter) Connect() *MockExamService_Connect_Call {
	return &MockExamService_Connect_Call{Call: _e.mock.On("Connect")}
}

func (_c *MockExamService_Connect_Call) Run(run func()) *MockExamService_Connect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockExamService_Connect_Call) Return(_a0 error) *MockExamService_Connect_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockExamService_Connect_Call) RunAndReturn(run func() error) *MockExamService_Connect_Call {
	_c.Call.Return(run)
	return _c
}

// CreateExam provides a mock function with given fields: topic, description, isPublic, userId
func (_m *MockExamService) CreateExam(topic string, description string, isPublic bool, userId string) (*pb.CreateExamResponse, error) {
	ret := _m.Called(topic, description, isPublic, userId)

	if len(ret) == 0 {
		panic("no return value specified for CreateExam")
	}

	var r0 *pb.CreateExamResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, bool, string) (*pb.CreateExamResponse, error)); ok {
		return rf(topic, description, isPublic, userId)
	}
	if rf, ok := ret.Get(0).(func(string, string, bool, string) *pb.CreateExamResponse); ok {
		r0 = rf(topic, description, isPublic, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pb.CreateExamResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, bool, string) error); ok {
		r1 = rf(topic, description, isPublic, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockExamService_CreateExam_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateExam'
type MockExamService_CreateExam_Call struct {
	*mock.Call
}

// CreateExam is a helper method to define mock.On call
//   - topic string
//   - description string
//   - isPublic bool
//   - userId string
func (_e *MockExamService_Expecter) CreateExam(topic interface{}, description interface{}, isPublic interface{}, userId interface{}) *MockExamService_CreateExam_Call {
	return &MockExamService_CreateExam_Call{Call: _e.mock.On("CreateExam", topic, description, isPublic, userId)}
}

func (_c *MockExamService_CreateExam_Call) Run(run func(topic string, description string, isPublic bool, userId string)) *MockExamService_CreateExam_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(bool), args[3].(string))
	})
	return _c
}

func (_c *MockExamService_CreateExam_Call) Return(_a0 *pb.CreateExamResponse, _a1 error) *MockExamService_CreateExam_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockExamService_CreateExam_Call) RunAndReturn(run func(string, string, bool, string) (*pb.CreateExamResponse, error)) *MockExamService_CreateExam_Call {
	_c.Call.Return(run)
	return _c
}

// CreateQuestion provides a mock function with given fields: examId, ask, answers, userId
func (_m *MockExamService) CreateQuestion(examId string, ask string, answers []string, userId string) (*pb.CreateQuestionResponse, error) {
	ret := _m.Called(examId, ask, answers, userId)

	if len(ret) == 0 {
		panic("no return value specified for CreateQuestion")
	}

	var r0 *pb.CreateQuestionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, []string, string) (*pb.CreateQuestionResponse, error)); ok {
		return rf(examId, ask, answers, userId)
	}
	if rf, ok := ret.Get(0).(func(string, string, []string, string) *pb.CreateQuestionResponse); ok {
		r0 = rf(examId, ask, answers, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pb.CreateQuestionResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, []string, string) error); ok {
		r1 = rf(examId, ask, answers, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockExamService_CreateQuestion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateQuestion'
type MockExamService_CreateQuestion_Call struct {
	*mock.Call
}

// CreateQuestion is a helper method to define mock.On call
//   - examId string
//   - ask string
//   - answers []string
//   - userId string
func (_e *MockExamService_Expecter) CreateQuestion(examId interface{}, ask interface{}, answers interface{}, userId interface{}) *MockExamService_CreateQuestion_Call {
	return &MockExamService_CreateQuestion_Call{Call: _e.mock.On("CreateQuestion", examId, ask, answers, userId)}
}

func (_c *MockExamService_CreateQuestion_Call) Run(run func(examId string, ask string, answers []string, userId string)) *MockExamService_CreateQuestion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].([]string), args[3].(string))
	})
	return _c
}

func (_c *MockExamService_CreateQuestion_Call) Return(_a0 *pb.CreateQuestionResponse, _a1 error) *MockExamService_CreateQuestion_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockExamService_CreateQuestion_Call) RunAndReturn(run func(string, string, []string, string) (*pb.CreateQuestionResponse, error)) *MockExamService_CreateQuestion_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteExam provides a mock function with given fields: examId, userId
func (_m *MockExamService) DeleteExam(examId string, userId string) (*pb.DeleteExamResponse, error) {
	ret := _m.Called(examId, userId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteExam")
	}

	var r0 *pb.DeleteExamResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*pb.DeleteExamResponse, error)); ok {
		return rf(examId, userId)
	}
	if rf, ok := ret.Get(0).(func(string, string) *pb.DeleteExamResponse); ok {
		r0 = rf(examId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pb.DeleteExamResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(examId, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockExamService_DeleteExam_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteExam'
type MockExamService_DeleteExam_Call struct {
	*mock.Call
}

// DeleteExam is a helper method to define mock.On call
//   - examId string
//   - userId string
func (_e *MockExamService_Expecter) DeleteExam(examId interface{}, userId interface{}) *MockExamService_DeleteExam_Call {
	return &MockExamService_DeleteExam_Call{Call: _e.mock.On("DeleteExam", examId, userId)}
}

func (_c *MockExamService_DeleteExam_Call) Run(run func(examId string, userId string)) *MockExamService_DeleteExam_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockExamService_DeleteExam_Call) Return(_a0 *pb.DeleteExamResponse, _a1 error) *MockExamService_DeleteExam_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockExamService_DeleteExam_Call) RunAndReturn(run func(string, string) (*pb.DeleteExamResponse, error)) *MockExamService_DeleteExam_Call {
	_c.Call.Return(run)
	return _c
}

// Disconnect provides a mock function with given fields:
func (_m *MockExamService) Disconnect() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Disconnect")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockExamService_Disconnect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Disconnect'
type MockExamService_Disconnect_Call struct {
	*mock.Call
}

// Disconnect is a helper method to define mock.On call
func (_e *MockExamService_Expecter) Disconnect() *MockExamService_Disconnect_Call {
	return &MockExamService_Disconnect_Call{Call: _e.mock.On("Disconnect")}
}

func (_c *MockExamService_Disconnect_Call) Run(run func()) *MockExamService_Disconnect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockExamService_Disconnect_Call) Return(_a0 error) *MockExamService_Disconnect_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockExamService_Disconnect_Call) RunAndReturn(run func() error) *MockExamService_Disconnect_Call {
	_c.Call.Return(run)
	return _c
}

// FindExams provides a mock function with given fields: pageIndex, pageSize, userId
func (_m *MockExamService) FindExams(pageIndex int32, pageSize int32, userId string) (*pb.FindExamsResponse, error) {
	ret := _m.Called(pageIndex, pageSize, userId)

	if len(ret) == 0 {
		panic("no return value specified for FindExams")
	}

	var r0 *pb.FindExamsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(int32, int32, string) (*pb.FindExamsResponse, error)); ok {
		return rf(pageIndex, pageSize, userId)
	}
	if rf, ok := ret.Get(0).(func(int32, int32, string) *pb.FindExamsResponse); ok {
		r0 = rf(pageIndex, pageSize, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pb.FindExamsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(int32, int32, string) error); ok {
		r1 = rf(pageIndex, pageSize, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockExamService_FindExams_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindExams'
type MockExamService_FindExams_Call struct {
	*mock.Call
}

// FindExams is a helper method to define mock.On call
//   - pageIndex int32
//   - pageSize int32
//   - userId string
func (_e *MockExamService_Expecter) FindExams(pageIndex interface{}, pageSize interface{}, userId interface{}) *MockExamService_FindExams_Call {
	return &MockExamService_FindExams_Call{Call: _e.mock.On("FindExams", pageIndex, pageSize, userId)}
}

func (_c *MockExamService_FindExams_Call) Run(run func(pageIndex int32, pageSize int32, userId string)) *MockExamService_FindExams_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int32), args[1].(int32), args[2].(string))
	})
	return _c
}

func (_c *MockExamService_FindExams_Call) Return(_a0 *pb.FindExamsResponse, _a1 error) *MockExamService_FindExams_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockExamService_FindExams_Call) RunAndReturn(run func(int32, int32, string) (*pb.FindExamsResponse, error)) *MockExamService_FindExams_Call {
	_c.Call.Return(run)
	return _c
}

// FindQuestions provides a mock function with given fields: pageIndex, pageSize, examId, userId
func (_m *MockExamService) FindQuestions(pageIndex int32, pageSize int32, examId string, userId string) (*pb.FindQuestionsResponse, error) {
	ret := _m.Called(pageIndex, pageSize, examId, userId)

	if len(ret) == 0 {
		panic("no return value specified for FindQuestions")
	}

	var r0 *pb.FindQuestionsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(int32, int32, string, string) (*pb.FindQuestionsResponse, error)); ok {
		return rf(pageIndex, pageSize, examId, userId)
	}
	if rf, ok := ret.Get(0).(func(int32, int32, string, string) *pb.FindQuestionsResponse); ok {
		r0 = rf(pageIndex, pageSize, examId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pb.FindQuestionsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(int32, int32, string, string) error); ok {
		r1 = rf(pageIndex, pageSize, examId, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockExamService_FindQuestions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindQuestions'
type MockExamService_FindQuestions_Call struct {
	*mock.Call
}

// FindQuestions is a helper method to define mock.On call
//   - pageIndex int32
//   - pageSize int32
//   - examId string
//   - userId string
func (_e *MockExamService_Expecter) FindQuestions(pageIndex interface{}, pageSize interface{}, examId interface{}, userId interface{}) *MockExamService_FindQuestions_Call {
	return &MockExamService_FindQuestions_Call{Call: _e.mock.On("FindQuestions", pageIndex, pageSize, examId, userId)}
}

func (_c *MockExamService_FindQuestions_Call) Run(run func(pageIndex int32, pageSize int32, examId string, userId string)) *MockExamService_FindQuestions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int32), args[1].(int32), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *MockExamService_FindQuestions_Call) Return(_a0 *pb.FindQuestionsResponse, _a1 error) *MockExamService_FindQuestions_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockExamService_FindQuestions_Call) RunAndReturn(run func(int32, int32, string, string) (*pb.FindQuestionsResponse, error)) *MockExamService_FindQuestions_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateExam provides a mock function with given fields: examId, topic, description, isPublic, userId
func (_m *MockExamService) UpdateExam(examId string, topic string, description string, isPublic bool, userId string) (*pb.UpdateExamResponse, error) {
	ret := _m.Called(examId, topic, description, isPublic, userId)

	if len(ret) == 0 {
		panic("no return value specified for UpdateExam")
	}

	var r0 *pb.UpdateExamResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, string, bool, string) (*pb.UpdateExamResponse, error)); ok {
		return rf(examId, topic, description, isPublic, userId)
	}
	if rf, ok := ret.Get(0).(func(string, string, string, bool, string) *pb.UpdateExamResponse); ok {
		r0 = rf(examId, topic, description, isPublic, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pb.UpdateExamResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, string, bool, string) error); ok {
		r1 = rf(examId, topic, description, isPublic, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockExamService_UpdateExam_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateExam'
type MockExamService_UpdateExam_Call struct {
	*mock.Call
}

// UpdateExam is a helper method to define mock.On call
//   - examId string
//   - topic string
//   - description string
//   - isPublic bool
//   - userId string
func (_e *MockExamService_Expecter) UpdateExam(examId interface{}, topic interface{}, description interface{}, isPublic interface{}, userId interface{}) *MockExamService_UpdateExam_Call {
	return &MockExamService_UpdateExam_Call{Call: _e.mock.On("UpdateExam", examId, topic, description, isPublic, userId)}
}

func (_c *MockExamService_UpdateExam_Call) Run(run func(examId string, topic string, description string, isPublic bool, userId string)) *MockExamService_UpdateExam_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string), args[3].(bool), args[4].(string))
	})
	return _c
}

func (_c *MockExamService_UpdateExam_Call) Return(_a0 *pb.UpdateExamResponse, _a1 error) *MockExamService_UpdateExam_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockExamService_UpdateExam_Call) RunAndReturn(run func(string, string, string, bool, string) (*pb.UpdateExamResponse, error)) *MockExamService_UpdateExam_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateQuestion provides a mock function with given fields: questionId, ask, answers, userId
func (_m *MockExamService) UpdateQuestion(questionId string, ask string, answers []string, userId string) (*pb.UpdateQuestionResponse, error) {
	ret := _m.Called(questionId, ask, answers, userId)

	if len(ret) == 0 {
		panic("no return value specified for UpdateQuestion")
	}

	var r0 *pb.UpdateQuestionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, []string, string) (*pb.UpdateQuestionResponse, error)); ok {
		return rf(questionId, ask, answers, userId)
	}
	if rf, ok := ret.Get(0).(func(string, string, []string, string) *pb.UpdateQuestionResponse); ok {
		r0 = rf(questionId, ask, answers, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pb.UpdateQuestionResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, []string, string) error); ok {
		r1 = rf(questionId, ask, answers, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockExamService_UpdateQuestion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateQuestion'
type MockExamService_UpdateQuestion_Call struct {
	*mock.Call
}

// UpdateQuestion is a helper method to define mock.On call
//   - questionId string
//   - ask string
//   - answers []string
//   - userId string
func (_e *MockExamService_Expecter) UpdateQuestion(questionId interface{}, ask interface{}, answers interface{}, userId interface{}) *MockExamService_UpdateQuestion_Call {
	return &MockExamService_UpdateQuestion_Call{Call: _e.mock.On("UpdateQuestion", questionId, ask, answers, userId)}
}

func (_c *MockExamService_UpdateQuestion_Call) Run(run func(questionId string, ask string, answers []string, userId string)) *MockExamService_UpdateQuestion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].([]string), args[3].(string))
	})
	return _c
}

func (_c *MockExamService_UpdateQuestion_Call) Return(_a0 *pb.UpdateQuestionResponse, _a1 error) *MockExamService_UpdateQuestion_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockExamService_UpdateQuestion_Call) RunAndReturn(run func(string, string, []string, string) (*pb.UpdateQuestionResponse, error)) *MockExamService_UpdateQuestion_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockExamService creates a new instance of MockExamService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockExamService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockExamService {
	mock := &MockExamService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
