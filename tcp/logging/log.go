package logging

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-22

const (
	dateFormat        = "2006-01-02"
	timeFormat        = "2006-01-02 15:04:05.999"
	maxFunctionLength = 30
	fileSuffix        = ".log"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	Log.SetLevel(logrus.InfoLevel)
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

	datetime := entry.Time.Format(timeFormat)
	if len(datetime) < len(timeFormat) {
		for i := 0; i < len(timeFormat)-len(datetime); i++ {
			datetime = datetime + "0"
		}
	}
	logLevel := strings.ToUpper(entry.Level.String())
	if len(logLevel) < 5 {
		logLevel = logLevel + " "
	}
	function := entry.Caller.Function[strings.LastIndex(entry.Caller.Function, "/")+1:]
	funcLen := len(function)
	if funcLen < maxFunctionLength {
		for i := 0; i < maxFunctionLength-funcLen; i++ {
			function = function + " "
		}
	} else if funcLen > maxFunctionLength {
		function = function[len(function)-maxFunctionLength:]
	}
	logStr := fmt.Sprintf("%s %s [%s] - %s\n", datetime, logLevel, function, entry.Message)

	buf.WriteString(logStr)
	return buf.Bytes(), nil
}

// LogFile 日志输出文件
func LogFile() *os.File {
	dir, _ := os.Getwd()
	logDirPath := dir + "/log/tcp"
	_, err := os.Stat(logDirPath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(logDirPath, 0777); err != nil {
			log.Fatalln("Create log dir failed")
			return nil
		}
	}
	logFileName := time.Now().Format(dateFormat) + fileSuffix
	fileName := path.Join(logDirPath, logFileName)
	if _, err = os.Stat(fileName); err != nil {
		if _, err = os.Create(fileName); err != nil {
			log.Println("Create log file failed")
			return nil
		}
	}
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Println("Open log file failed")
		return nil
	}
	return logFile
}
