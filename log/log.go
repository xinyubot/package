package log

import (
	"errors"
	"io"
	"reflect"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// newLogger creates and returns a pointer to a new zap logger ^ ^
func newLogger(scenario, logminlevel string, outstream io.Writer) (zaplogger *zap.Logger, sync func() error, err error) {
	// if any of the following steps panics in an unforeseen way, deferred recovery will catch it
	defer func() {
		if err := recover(); err != nil {
			err = errors.New(reflect.ValueOf(err).String())
		}
	}()

	_ = scenario

	// New a zap Logger and set the output dest to passed in `outstream`
	t := zapcore.AddSync(outstream)
	syncers := []zapcore.WriteSyncer{t}
	var zapLevel zapcore.Level
	if err := zapLevel.Set(logminlevel); err != nil {
		return nil, nil, err
	}
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapLevel)
	core := zapcore.NewCore(
		// 配置 时间，等级，caller，stacktrace 键值
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),

		// 多写，文件和 stdout
		zapcore.NewMultiWriteSyncer(syncers...),
		atomicLevel,
	)

	// create a new zap logger
	zaplogger = zap.New(core,
		zap.AddCaller(),
		// 0 是 Helper 能获取正确的 caller 的值
		zap.AddCallerSkip(0),
		// 自动加上 stacktrace 最小等级
		zap.AddStacktrace(zapcore.FatalLevel),
	)
	sync = zaplogger.Sync

	// RET
	return
}

// Register creates new zap loggers and regists them to global var, then returns zap.Sync() for graceful shutdown
func Register(logpath, logminlevel string) (syncers []func() error, err error) {
	// graceful shutdown
	errFlag := false
	defer func() {
		if errFlag {
			for _, syncer := range syncers {
				_ = syncer() // 关闭已注册的 Logger
			}
		}
	}()

	// 注册 Error Logger
	erl, err := getRotateLog(logpath + ErrorLogFile)
	if err != nil {
		errFlag = true
		return
	}
	el, es, err := newLogger("Error", logminlevel, erl)
	if err != nil {
		errFlag = true
		return
	}
	syncers = append(syncers, es)

	// 注册 Request Logger
	rrl, err := getRotateLog(logpath + RequestLogFile)
	if err != nil {
		errFlag = true
		return
	}
	rl, rs, err := newLogger("Request", logminlevel, rrl)
	if err != nil {
		errFlag = true
		return
	}
	syncers = append(syncers, rs)

	// 注册 Call Logger
	crl, err := getRotateLog(logpath + CallLogFile)
	if err != nil {
		errFlag = true
		return
	}
	cl, cs, err := newLogger("Call", logminlevel, crl)
	if err != nil {
		errFlag = true
		return
	}
	syncers = append(syncers, cs)

	// 注册 Debug Logger
	drl, err := getRotateLog(logpath + DebugLogFile)
	if err != nil {
		errFlag = true
		return
	}
	dl, ds, err := newLogger("Error", "debug", drl)
	if err != nil {
		errFlag = true
		return
	}
	syncers = append(syncers, ds)

	// 注册 Logger 们到全局变量中
	ErrorLogger = el
	RequestLogger = rl
	CallLogger = cl
	DebugLogger = dl

	return
}

// getRotateLog returns a hook created by NewRotateLog
func getRotateLog(filename string) (*RotateLog, error) {
	return NewRoteteLog(
		filename+"_%d_%d_%d.log",
		WithLinkPath(filename+".log"),
		WithRotateTime(time.Minute),
	)
}
