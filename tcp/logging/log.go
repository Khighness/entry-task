package logging

import (
	"github.com/sirupsen/logrus"

	"github.com/Khighness/entry-task/pkg/log"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-22

var Log *logrus.Logger = log.NewLogger(logrus.InfoLevel, "tcp", true)
