package logger

import (
	"io"
	"os"

	"github.com/PostService/infrastructure/config"
	log "github.com/sirupsen/logrus"
)

// Logger interface for logging info
type Logger interface {
	Print(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})
	Info(args ...interface{})
	Printf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Errorln(args ...interface{})
}

type loggerImpl struct{}

// New is func for initializing logger
func New(conf config.LoggerConfig) (Logger, error) {
	f, err := os.OpenFile(conf.FileName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	log.SetOutput(f)
	log.SetLevel(log.Level(conf.Level))
	mw := io.MultiWriter(f, os.Stderr)
	log.SetOutput(mw)

	return &loggerImpl{}, err
}

// Print logs a message at level Info on the standard logger.
func (*loggerImpl) Print(args ...interface{}) {
	log.Print(args...)
}

// Info logs a message at level Info on the standard logger.
func (*loggerImpl) Info(args ...interface{}) {
	log.Info(args...)
}

// Error logs a message at level Error on the standard logger.
func (*loggerImpl) Error(args ...interface{}) {
	log.Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func (*loggerImpl) Panic(args ...interface{}) {
	log.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func (*loggerImpl) Fatal(args ...interface{}) {
	log.Fatal(args...)
}

// Printf logs a message at level Info on the standard logger.
func (*loggerImpl) Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func (*loggerImpl) Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func (*loggerImpl) Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func (*loggerImpl) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Errorln logs a message at level Error on the standard logger.
func (*loggerImpl) Errorln(args ...interface{}) {
	log.Errorln(args...)
}
