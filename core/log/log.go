package log

import (
	"encoding/json"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var std = &Logger{logger: logrus.StandardLogger()}

type Config struct {
	FileLog string
	Level   string
	UseJSON bool
}

type Logger struct {
	logger *logrus.Logger
}

func newLogger(conf Config) (*Logger, error) {
	logger := logrus.New()

	if conf.Level == "" {
		conf.Level = "info"
	}
	level, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		level = logrus.ErrorLevel
	}
	logger.SetLevel(level)

	var formatter logrus.Formatter
	formatter = &logrus.TextFormatter{}
	if conf.UseJSON {
		formatter = &logrus.JSONFormatter{}
	}
	logger.SetFormatter(formatter)

	if conf.FileLog != "" {
		logger.Hooks.Add(lfshook.NewHook(conf.FileLog, formatter))
	}

	return &Logger{logger: logger}, nil
}

func NewLogger(conf Config) (*Logger, error) {
	return newLogger(conf)
}

func Init(conf Config) error {
	logger, err := newLogger(conf)
	if err != nil {
		return err
	}

	std = logger
	return nil
}

func Println(msg string) {
	std.Println(msg)
}
func Printf(format string, args ...interface{}) {
	std.Printf(format, args...)
}
func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}
func DebugWithField(param map[string]interface{}, message string, args ...interface{}) {
	std.DebugWithField(param, message, args...)
}
func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}
func InfoWithField(param map[string]interface{}, message string, args ...interface{}) {
	std.InfoWithField(param, message, args...)
}
func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}
func ErrorWithField(param map[string]interface{}, message string, args ...interface{}) {
	std.ErrorWithField(param, message, args...)
}
func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}
func FatalWithField(param map[string]interface{}, message string, args ...interface{}) {
	std.FatalWithField(param, message, args...)
}
func PrettyJSON(input interface{}) {
	out, _ := json.MarshalIndent(input, "", "  ")
	std.Println(string(out))
}

func (l *Logger) Println(msg string) {
	l.logger.Println(msg)
}
func (l *Logger) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}
func (l *Logger) DebugWithField(param map[string]interface{}, message string, args ...interface{}) {
	l.logger.WithFields(param).Debugf(message, args...)
}
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}
func (l *Logger) InfoWithField(param map[string]interface{}, message string, args ...interface{}) {
	l.logger.WithFields(param).Infof(message, args...)
}
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}
func (l *Logger) ErrorWithField(param map[string]interface{}, message string, args ...interface{}) {
	l.logger.WithFields(param).Errorf(message, args...)
}
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}
func (l *Logger) FatalWithField(param map[string]interface{}, message string, args ...interface{}) {
	l.logger.WithFields(param).Fatalf(message, args...)
}
