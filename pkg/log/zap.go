package log

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strconv"
)

var _ Logger = (*ZapLogger)(nil)

const (
	LogLevel    = "LOG_LEVEL"
	LogFile     = "LOG_FILE"
	LogFileSize = "LOG_FILE_SIZE"
	LogFileNum  = "LOG_FILE_NUM"
	LogFileAge  = "LOG_FILE_AGE"
	LogStdout   = "LOG_STDOUT"
)

type ZapLogger struct {
	log  *zap.Logger
	Sync func() error
}

func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// NewDefaultJsonLogger 配置zap日志,将zap日志库引入
func NewDefaultJsonLogger() Logger {
	//配置zap日志库的编码器
	encoder := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "",
		StacktraceKey:  "stack",
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	logLevel, err := zapcore.ParseLevel(env(LogLevel, "info"))
	if err != nil {
		logLevel = zap.InfoLevel
	}
	l := NewZapLogger(
		encoder,
		zap.NewAtomicLevelAt(logLevel),
		zap.AddStacktrace(
			zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		//zap.AddCaller(),
		//zap.AddCallerSkip(2),
		zap.Development(),
	)
	return With(l, "caller", DefaultCaller)
}

// 日志自动切割，采用 lumberjack 实现的
func getLogWriter() zapcore.WriteSyncer {
	maxSize, err := strconv.Atoi(env(LogFileSize, "1024"))
	if err != nil {
		panic(err)
	}
	maxNum, err := strconv.Atoi(env(LogFileNum, "3"))
	if err != nil {
		panic(err)
	}
	maxAge, err := strconv.Atoi(env(LogFileAge, "30"))
	if err != nil {
		panic(err)
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename:   env(LogFile, "/var/log/agile-cloud/app.log"), //指定日志存储位置
		MaxSize:    maxSize,                                      //日志的最大大小（M）
		MaxBackups: maxNum,                                       //日志的最大保存数量
		MaxAge:     maxAge,                                       //日志文件存储最大天数
		Compress:   false,                                        //是否执行压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

// NewZapLogger return a zap logger.
func NewZapLogger(encoder zapcore.EncoderConfig, level zap.AtomicLevel, opts ...zap.Option) *ZapLogger {
	//日志切割
	writeSyncer := getLogWriter()
	//设置日志级别
	var core zapcore.Core
	writers := []zapcore.WriteSyncer{zapcore.AddSync(writeSyncer)}
	if env(LogStdout, "false") != "false" {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}
	core = zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder),
		zapcore.NewMultiWriteSyncer(writers...),
		level,
	)
	zapLogger := zap.New(core, opts...)
	return &ZapLogger{log: zapLogger, Sync: zapLogger.Sync}
}

// Log 实现log接口
func (l *ZapLogger) Log(level Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	switch level {
	case LevelDebug:
		l.log.Debug("", data...)
	case LevelInfo:
		l.log.Info("", data...)
	case LevelWarn:
		l.log.Warn("", data...)
	case LevelError:
		l.log.Error("", data...)
	case LevelFatal:
		l.log.Fatal("", data...)
	}
	return nil
}
