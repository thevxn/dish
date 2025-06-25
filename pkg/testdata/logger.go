package testhelpers

// MockLogger is a mock implementation of the Logger interface with empty method implementations.
type MockLogger struct{}

func (l *MockLogger) Trace(v ...any)                 {}
func (l *MockLogger) Tracef(format string, v ...any) {}
func (l *MockLogger) Debug(v ...any)                 {}
func (l *MockLogger) Debugf(format string, v ...any) {}
func (l *MockLogger) Info(v ...any)                  {}
func (l *MockLogger) Infof(format string, v ...any)  {}
func (l *MockLogger) Warn(v ...any)                  {}
func (l *MockLogger) Warnf(format string, v ...any)  {}
func (l *MockLogger) Error(v ...any)                 {}
func (l *MockLogger) Errorf(format string, v ...any) {}
func (l *MockLogger) Panic(v ...any)                 {}
func (l *MockLogger) Panicf(format string, v ...any) {}
