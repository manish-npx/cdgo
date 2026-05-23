package logging

import (
"go.uber.org/zap"
"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger for structured logging
type Logger struct {
*zap.Logger
}

// New creates a new logger instance
func New() *Logger {
config := zap.Config{
Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
Development: false,
Encoding:    "console",
EncoderConfig: zapcore.EncoderConfig{
TimeKey:        "time",
LevelKey:       "level",
NameKey:        "name",
CallerKey:      "caller",
FunctionKey:    zapcore.OmitKey,
MessageKey:     "msg",
StacktraceKey:  "stack",
LineEnding:     zapcore.DefaultLineEnding,
EncodeLevel:    zapcore.CapitalColorLevelEncoder,
EncodeTime:     zapcore.ISO8601TimeEncoder,
EncodeDuration: zapcore.StringDurationEncoder,
EncodeCaller:   zapcore.ShortCallerEncoder,
},
OutputPaths:      []string{"stdout"},
ErrorOutputPaths: []string{"stderr"},
}

logger, err := config.Build()
if err != nil {
panic(err)
}

return &Logger{Logger: logger}
}

// NewDevelopment creates a logger for development
func NewDevelopment() *Logger {
logger, _ := zap.NewDevelopment()
return &Logger{Logger: logger}
}

// WithFields returns a logger with additional fields
func (l *Logger) WithFields(fields ...zap.Field) *Logger {
return &Logger{Logger: l.Logger.With(fields...)}
}

// Info logs an info message with structured fields
func (l *Logger) Info(msg string, fields ...zap.Field) {
l.Logger.Info(msg, fields...)
}

// Error logs an error message with structured fields
func (l *Logger) Error(msg string, fields ...zap.Field) {
l.Logger.Error(msg, fields...)
}

// Debug logs a debug message with structured fields
func (l *Logger) Debug(msg string, fields ...zap.Field) {
l.Logger.Debug(msg, fields...)
}

// Warn logs a warning message with structured fields
func (l *Logger) Warn(msg string, fields ...zap.Field) {
l.Logger.Warn(msg, fields...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
l.Logger.Fatal(msg, fields...)
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() {
l.Logger.Sync()
}
