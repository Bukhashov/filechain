package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"github.com/sirupsen/logrus"
)

type writerHook struct {
	Writer []io.Writer
	LogLevels []logrus.Level
}
func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String(); if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() Logger {
	return Logger{e}
}

func (l *Logger) GetLoggerWithField(k string, v interface{}) *Logger {
	return &Logger{l.WithField(k, v)}
}

func init(){
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	if err := os.MkdirAll("logs", 0644); err != nil {
		panic(err)
	}
	
	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0777); if err != nil {
		panic(err)
	}

	l.SetOutput(io.Discard)
	l.AddHook(&writerHook{
		Writer: []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})
	l.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(l)
}