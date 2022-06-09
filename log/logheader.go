package log

import (
	"go.uber.org/zap"
)

var (
	// 日志 Logger
	ErrorLogger   *zap.Logger
	RequestLogger *zap.Logger
	CallLogger    *zap.Logger
	DebugLogger   *zap.Logger

	// 文件名
	ErrorLogFile   string
	RequestLogFile string
	CallLogFile    string
	DebugLogFile   string
)
