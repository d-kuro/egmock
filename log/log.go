package log

import (
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewDevelopment()
	Error     = logger.Error
	Info      = logger.Info
)
