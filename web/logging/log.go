package logging

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-22

const (
	dateFormat string = "2006-01-02"
	timeFormat string = "2006-01-02 15:04:05"
	fileSuffix string = ".log"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	Log.SetLevel(logrus.WarnLevel)
	Log.SetFormatter(&LogFormatter{})
	Log.SetOutput(io.MultiWriter(LogFile(), os.Stdout))
	Log.SetReportCaller(true)
}

type LogFormatter struct{}

// Format 日志输出格式
func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buf *bytes.Buffer
	if entry.Buffer != nil {
		buf = entry.Buffer
	} else {
		buf = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format(timeFormat)
	var newLog string

	if entry.HasCaller() {
		newLog = fmt.Sprintf("%s [%s] (%s:%d | %s) %s\n",
			timestamp, entry.Level, filepath.Base(entry.Caller.File), entry.Caller.Line, entry.Caller.Function, entry.Message)
	} else {
		newLog = fmt.Sprintf("%s [%s] %s\n", timestamp, entry.Level, entry.Message)
	}

	buf.WriteString(newLog)
	return buf.Bytes(), nil
}

// LogFile 日志输出文件
func LogFile() *os.File {
	dir, _ := os.Getwd()
	logDirPath := dir + "/log/web"
	_, err := os.Stat(logDirPath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(logDirPath, 0777); err != nil {
			log.Fatalln("Create log dir failed")
			return nil
		}
	}
	logFileName := time.Now().Format(dateFormat) + fileSuffix
	fileName := path.Join(logDirPath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			log.Fatalln("Create log file failed")
			return nil
		}
	}
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatalln("Open log file failed")
		return nil
	}
	return logFile
}
