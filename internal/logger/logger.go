package logger

import (
	"encoding/json"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var TestLogger *AppLogger

func InitTestLogger() {
	TestLogger = &AppLogger{
		logger: log.With().Str("service", "test").Logger(),
	}
}


type Logger interface {
	Error(msg string)
	Errorf(msg string, values ...interface{})
	ErrorWithPayload(msg string, payload any, values ...interface{})
	Panic(msg string)
	Panicf(msg string, values ...interface{})
	PanicWithPayload(msg string, payload any, values ...interface{})
	Fatal(msg string)
	Fatalf(msg string, values ...interface{})
	FatalWithPayload(msg string, payload any, values ...interface{})
	Info(msg string)
	Infof(msg string, values ...interface{})
	InfoWithPayload(msg string, payload any, values ...interface{})
	Warn(msg string)
	Warnf(msg string, values ...interface{})
	WarnWithPayload(msg string, payload any, values ...interface{})
	Debug(msg string)
	Debugf(msg string, values ...interface{})
	DebugWithPayload(msg string, payload any, values ...interface{})
}

type AppLogger struct {
	logger zerolog.Logger
}

// NewAppLogger creates a new logger instance (no hooks).
func NewAppLogger(service string) *AppLogger {
	l := log.With().Str("service", service).Logger()
	return &AppLogger{logger: l}
}

func (l *AppLogger) Error(msg string) {
	l.logger.Error().Msg(msg)
}

func (l *AppLogger) Errorf(msg string, values ...interface{}) {
	l.logger.Error().Msgf(msg, values...)
}

func (l *AppLogger) ErrorWithPayload(msg string, payload any, values ...interface{}) {
	l.logger.Error().
		Bytes("body", l.getBytes(payload)).
		Msgf(msg, values...)
}

func (l *AppLogger) Panic(msg string) {
	l.logger.Panic().Msg(msg)
}

func (l *AppLogger) Panicf(msg string, values ...interface{}) {
	l.logger.Panic().Msgf(msg, values...)
}

func (l *AppLogger) PanicWithPayload(msg string, payload any, values ...interface{}) {
	l.logger.Panic().
		Bytes("body", l.getBytes(payload)).
		Msgf(msg, values...)
}

func (l *AppLogger) Fatal(msg string) {
	l.logger.Fatal().Msg(msg)
}

func (l *AppLogger) Fatalf(msg string, values ...interface{}) {
	l.logger.Fatal().Msgf(msg, values...)
}

func (l *AppLogger) FatalWithPayload(msg string, payload any, values ...interface{}) {
	l.logger.Fatal().
		Bytes("body", l.getBytes(payload)).
		Msgf(msg, values...)
}

func (l *AppLogger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

func (l *AppLogger) Infof(msg string, values ...interface{}) {
	l.logger.Info().Msgf(msg, values...)
}

func (l *AppLogger) InfoWithPayload(msg string, payload any, values ...interface{}) {
	l.logger.Info().
		Bytes("body", l.getBytes(payload)).
		Msgf(msg, values...)
}

func (l *AppLogger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

func (l *AppLogger) Warnf(msg string, values ...interface{}) {
	l.logger.Warn().Msgf(msg, values...)
}

func (l *AppLogger) WarnWithPayload(msg string, payload any, values ...interface{}) {
	l.logger.Warn().
		Bytes("body", l.getBytes(payload)).
		Msgf(msg, values...)
}

func (l *AppLogger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

func (l *AppLogger) Debugf(msg string, values ...interface{}) {
	l.logger.Debug().Msgf(msg, values...)
}

func (l *AppLogger) DebugWithPayload(msg string, payload any, values ...interface{}) {
	l.logger.Debug().
		Bytes("body", l.getBytes(payload)).
		Msgf(msg, values...)
}

func (l *AppLogger) getBytes(data any) []byte {
	body, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Printf("failed to get bytes for data: [%+v]", data)
		return nil
	}
	return body
}
