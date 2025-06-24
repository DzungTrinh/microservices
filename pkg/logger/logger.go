package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"microservices/user-management/pkg/constants"
	"microservices/user-management/pkg/logger/config"
	"os"
	"runtime"
	"strings"
	"time"
)

// For mapping config LogService to email_service LogService levels
var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

// ILogger methods interface
type ILogger interface {
	InitLogger()
	Sync() error
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	WarnMsg(msg string, err error)
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Err(msg string, err error)
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Printf(template string, args ...interface{})
	WithName(name string)
	HttpMiddlewareAccessLogger(method string, uri string, status int, size int64, time time.Duration)
	GrpcMiddlewareAccessLogger(method string, time time.Duration, metaData map[string][]string, err error)
	GrpcClientInterceptorLogger(method string, req interface{}, reply interface{}, time time.Duration, metaData map[string][]string, err error)
	KafkaProcessMessage(topic string, partition int, message string, workerID int, offset int64, time time.Time)
	KafkaLogCommittedMessage(topic string, partition int, offset int64)
}

// LoggerService methods struct
type LoggerService struct {
	level       string
	devMode     bool
	encoding    string
	path        string
	sugarLogger *zap.SugaredLogger
	logger      *zap.Logger
}

// NewLoggerService App ILogger constructor
func NewLoggerService(cfg config.Config) *LoggerService {
	loggerService := &LoggerService{level: cfg.LogLevel, devMode: cfg.DevMode, encoding: cfg.Encoder, path: cfg.Path}
	loggerService.InitLogger()
	return loggerService
}

func getPath(pathLog string) string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	path := dir + "/" + strings.TrimLeft(pathLog, "/")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	if strings.Contains(runtime.GOOS, "windows") {
		path = path + "\\"
	} else {
		path = path + "/"
	}
	now := time.Now()
	filePath := fmt.Sprintf("%d-%d-%d.log", now.Day(), now.Month(), now.Year())
	return path + filePath
}

func configure(pathLog string) zapcore.WriteSyncer {
	path := getPath(pathLog)

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    1, // megabytes
		MaxBackups: 4,
	})
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stderr),
		zapcore.AddSync(w),
	)
}

// InitLogger Init logger
func (l *LoggerService) InitLogger() {
	logLevel, exist := loggerLevelMap[l.level]
	if !exist {
		logLevel = zapcore.DebugLevel
	}
	logWriter := configure(l.path)

	var encoderCfg zapcore.EncoderConfig
	if l.devMode {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	var encoder zapcore.Encoder
	encoderCfg.NameKey = "[SERVICE]"
	encoderCfg.TimeKey = "[TIME]"
	encoderCfg.LevelKey = "[LEVEL]"
	//encoderCfg.FunctionKey = "[CALLER]"
	encoderCfg.CallerKey = "[LINE]"
	encoderCfg.MessageKey = "[MESSAGE]"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	encoderCfg.EncodeName = zapcore.FullNameEncoder
	encoderCfg.EncodeDuration = zapcore.StringDurationEncoder

	if l.encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.logger = logger
	l.sugarLogger = logger.Sugar()
}

// WithName add logger microservice name
func (l *LoggerService) WithName(name string) {
	l.logger = l.logger.Named(name)
	l.sugarLogger = l.sugarLogger.Named(name)
}

// Debug uses fmt.Sprint to construct and log a message.
func (l *LoggerService) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

// Debugf uses fmt.Sprintf to log a templated message
func (l *LoggerService) Debugf(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

// Info uses fmt.Sprint to construct and log a message
func (l *LoggerService) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (l *LoggerService) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

// Printf uses fmt.Sprintf to log a templated message
func (l *LoggerService) Printf(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func (l *LoggerService) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l *LoggerService) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

// Error uses fmt.Sprint to construct and log a message.
func (l *LoggerService) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l *LoggerService) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the logger then panics. (See DPanicLevel for details.)
func (l *LoggerService) DPanic(args ...interface{}) {
	l.sugarLogger.DPanic(args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the logger then panics. (See DPanicLevel for details.)
func (l *LoggerService) DPanicf(template string, args ...interface{}) {
	l.sugarLogger.DPanicf(template, args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (l *LoggerService) Panic(args ...interface{}) {
	l.sugarLogger.Panic(args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics
func (l *LoggerService) Panicf(template string, args ...interface{}) {
	l.sugarLogger.Panicf(template, args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l *LoggerService) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *LoggerService) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}

// Sync flushes any buffered log entries
func (l *LoggerService) Sync() error {
	go l.logger.Sync() // nolint: errcheck
	return l.sugarLogger.Sync()
}

func (l *LoggerService) HttpMiddlewareAccessLogger(method, uri string, status int, size int64, time time.Duration) {
	l.logger.Info(
		constants.HTTP,
		zap.String(constants.METHOD, method),
		zap.String(constants.URI, uri),
		zap.Int(constants.STATUS, status),
		zap.Int64(constants.SIZE, size),
		zap.Duration(constants.TIME, time),
	)
}

func (l *LoggerService) GrpcMiddlewareAccessLogger(method string, time time.Duration, metaData map[string][]string, err error) {
	l.logger.Info(
		constants.GRPC,
		zap.String(constants.METHOD, method),
		zap.Duration(constants.TIME, time),
		zap.Any(constants.METADATA, metaData),
		zap.Error(err),
	)
}

func (l *LoggerService) GrpcClientInterceptorLogger(method string, req, reply interface{}, time time.Duration, metaData map[string][]string, err error) {
	l.logger.Info(
		constants.GRPC,
		zap.String(constants.METHOD, method),
		zap.Any(constants.REQUEST, req),
		zap.Any(constants.REPLY, reply),
		zap.Duration(constants.TIME, time),
		zap.Any(constants.METADATA, metaData),
		zap.Error(err),
	)
}

func (l *LoggerService) KafkaProcessMessage(topic string, partition int, message string, workerID int, offset int64, time time.Time) {
	l.logger.Debug(
		"Processing Kafka message",
		zap.String(constants.Topic, topic),
		zap.Int(constants.Partition, partition),
		zap.String(constants.Message, message),
		zap.Int(constants.WorkerID, workerID),
		zap.Int64(constants.Offset, offset),
		zap.Time(constants.Time, time),
	)
}

func (l *LoggerService) KafkaLogCommittedMessage(topic string, partition int, offset int64) {
	l.logger.Info(
		"Committed Kafka message",
		zap.String(constants.Topic, topic),
		zap.Int(constants.Partition, partition),
		zap.Int64(constants.Offset, offset),
	)
}
