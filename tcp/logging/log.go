package logging

import (
	"github.com/sirupsen/logrus"

	"github.com/Khighness/entry-task/pkg/logger"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-22

var Log *logrus.Logger = logger.NewLogger(logrus.InfoLevel, "tcp", true)
