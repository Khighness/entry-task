package logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-22

const (
	dateFormat        = "2006-01-02"
	timeFormat        = "2006-01-02 15:04:05.999"
	maxFunctionLength = 30
	fileSuffix        = ".logger"
)

// NewLogger 创建Logger
func NewLogger(logLevel logrus.Level, logFile string, enableConsole bool) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logLevel)
	logger.SetFormatter(&LogFormatter{})
	if enableConsole && logFile != "" {
		logger.SetOutput(io.MultiWriter(LogFile(logFile), os.Stdout))
	} else if enableConsole {
		logger.SetOutput(os.Stdout)
	} else if logFile != "" {
		logger.SetOutput(LogFile(logFile))
	} else {
		logger.SetOutput(nil)
	}
	logger.SetReportCaller(true)
	return logger
}

// LogFormatter 自定义格式
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
	if len(logLevel) > 5 {
		logLevel = logLevel[:5]
	} else if len(logLevel) < 5 {
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
func LogFile(file string) *os.File {
	dir, _ := os.Getwd()
	logDirPath := dir + "/logger/" + file
	_, err := os.Stat(logDirPath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(logDirPath, 0777); err != nil {
			log.Fatalln("Create logger dir failed")
			return nil
		}
	}
	logFileName := time.Now().Format(dateFormat) + fileSuffix
	fileName := path.Join(logDirPath, logFileName)
	if _, err = os.Stat(fileName); err != nil {
		if _, err = os.Create(fileName); err != nil {
			log.Println("Create logger file failed")
			return nil
		}
	}
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Println("Open logger file failed")
		return nil
	}
	return logFile
}
