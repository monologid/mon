package mon

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func InitiateLog(level string) {
	logLevel := logrus.InfoLevel

	switch strings.ToLower(level) {
	case "panic":
		logLevel = logrus.PanicLevel
	case "fatal":
		logLevel = logrus.FatalLevel
	case "error":
		logLevel = logrus.ErrorLevel
	case "warn":
		logLevel = logrus.WarnLevel
	case "info":
		logLevel = logrus.InfoLevel
	case "debug":
		logLevel = logrus.DebugLevel
	case "trace":
		logLevel = logrus.TraceLevel
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logLevel)
}

type ILog interface {
	SetMessage(msg string) ILog
	SetError(err error) ILog

	Info()
	Debug()
	Error(err error) error
}

type log struct {
	Filename string
	Method   string
	Message  string
	Data     interface{}
	Err      error
}

func (l *log) SetMessage(msg string) ILog {
	l.Message = msg
	return l
}

func (l *log) SetError(err error) ILog {
	l.Err = err
	return l
}

func (l *log) get() *logrus.Entry {
	fields := logrus.Fields{
		"filename": l.Filename,
		"method":   l.Method,
		"data":     l.Data,
	}

	if l.Err != nil {
		fields["error"] = l.Err.Error()
	}

	return logrus.WithFields(fields)
}

func (l *log) Info() {
	l.get().Info(l.Message)
}

func (l *log) Debug() {
	l.get().Debug(l.Message)
}

func (l *log) Error(errMsg error) error {
	l.get().Error(errMsg.Error())
	return errMsg
}

func NewLog(filename string, method string, data interface{}) ILog {
	return &log{
		Filename: filename,
		Method:   method,
		Data:     data,
	}
}
