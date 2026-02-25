package logger

// ActiveLogger holds the currently active logger instance
var ActiveLogger Logger = NewAppLogger("default")

// SetLogger lets you replace the active logger (useful for tests)
func SetLogger(l Logger) {
    ActiveLogger = l
}

// NoopLogger implements Logger with no-ops
type NoopLogger struct{}

func (NoopLogger) Error(string)                                 {}
func (NoopLogger) Errorf(string, ...interface{})                {}
func (NoopLogger) ErrorWithPayload(string, any, ...interface{}) {}
func (NoopLogger) Panic(string)                                 {}
func (NoopLogger) Panicf(string, ...interface{})                {}
func (NoopLogger) PanicWithPayload(string, any, ...interface{}) {}
func (NoopLogger) Fatal(string)                                 {}
func (NoopLogger) Fatalf(string, ...interface{})                {}
func (NoopLogger) FatalWithPayload(string, any, ...interface{}) {}
func (NoopLogger) Info(string)                                  {}
func (NoopLogger) Infof(string, ...interface{})                 {}
func (NoopLogger) InfoWithPayload(string, any, ...interface{})  {}
func (NoopLogger) Warn(string)                                  {}
func (NoopLogger) Warnf(string, ...interface{})                 {}
func (NoopLogger) WarnWithPayload(string, any, ...interface{})  {}
func (NoopLogger) Debug(string)                                 {}
func (NoopLogger) Debugf(string, ...interface{})                {}
func (NoopLogger) DebugWithPayload(string, any, ...interface{}) {}
