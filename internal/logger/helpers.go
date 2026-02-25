package logger

var logger *AppLogger

// InitLogger initializes a global logger with a service name.
func InitLogger(service string) {
	logger = NewAppLogger(service)
}

func Error(msg string) {
	logger.Error(msg)
}

func Errorf(msg string, values ...interface{}) {
	logger.Errorf(msg, values...)
}

func ErrorWithPayload(msg string, payload any, values ...interface{}) {
	logger.ErrorWithPayload(msg, payload, values...)
}

func Warn(msg string) {
	logger.Warn(msg)
}

func Warnf(msg string, values ...interface{}) {
	logger.Warnf(msg, values...)
}

func WarnWithPayload(msg string, payload any, values ...interface{}) {
	logger.WarnWithPayload(msg, payload, values...)
}

func Info(msg string) {
	logger.Info(msg)
}

func Infof(msg string, values ...interface{}) {
	logger.Infof(msg, values...)
}

func InfoWithPayload(msg string, payload any, values ...interface{}) {
	logger.InfoWithPayload(msg, payload, values...)
}

func Fatal(msg string) {
	logger.Fatal(msg)
}

func Fatalf(msg string, values ...interface{}) {
	logger.Fatalf(msg, values...)
}

func FatalWithPayload(msg string, payload any, values ...interface{}) {
	logger.FatalWithPayload(msg, payload, values...)
}

func Panic(msg string) {
	logger.Panic(msg)
}

func Panicf(msg string, values ...interface{}) {
	logger.Panicf(msg, values...)
}

func PanicWithPayload(msg string, payload any, values ...interface{}) {
	logger.PanicWithPayload(msg, payload, values...)
}

func Debug(msg string) {
	logger.Debug(msg)
}

func Debugf(msg string, values ...interface{}) {
	logger.Debugf(msg, values...)
}

func DebugWithPayload(msg string, payload any, values ...interface{}) {
	logger.DebugWithPayload(msg, payload, values...)
}
