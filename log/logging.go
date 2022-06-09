package log

// RegisterLogger regists all loggers and returns the closer functions of the loggers
func RegisterLogger(logPath string) (syncers []func() error, err error) {
	return Register(logPath, "Info") // 统一设置成 Info 完事，反正 Debug 的信息也走的是 DebugLogger.Info()...
}
