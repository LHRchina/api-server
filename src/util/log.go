package util

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	infoLog, warnLog, errLog = logrus.New(), logrus.New(), logrus.New()
)

func initLog(l *logrus.Logger, path string, level logrus.Level) {
	out, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("open "+path+" err:", err)
		os.Exit(-1)
	}
	l.Out = out
	l.SetLevel(level)
}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	initLog(infoLog, "log/info.log", logrus.InfoLevel)
	initLog(warnLog, "log/warn.log", logrus.WarnLevel)
	initLog(errLog, "log/err.log", logrus.ErrorLevel)
}

func Info(v ...interface{}) {
	infoLog.Println(v)
}

func Err(v ...interface{}) {
	errLog.Println(v)
}

func Warn(v ...interface{}) {
	warnLog.Println(v)
}
