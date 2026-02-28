// Package helper...TODO
package helper

import (
	l "github.com/ccthomas/board-game/internal/logger/mock"
	"go.uber.org/mock/gomock"
)

// TODO - generated during AI Iteration Loop. Left be as low priority item to clean up.
// Moving forward to focus on delivering initial working product
func AssertError() error {
	return &MockError{}
}

type MockError struct{}

func (m *MockError) Error() string { return "failed" }

func NewDummyMockedLogger(ctrl *gomock.Controller) *l.MockLogger {
	mockLogger := l.NewMockLogger(ctrl)
	mockLogger.EXPECT().WithFields(gomock.Any()).Return(mockLogger).AnyTimes()
	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Trace(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	return mockLogger
}
