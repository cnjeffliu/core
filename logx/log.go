/*
 * @Author: cnzf1
 * @Date: 2019-01-26 15:20:02
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-17 18:36:58
 * @Description: 日志组件封装
 */
package logx

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger
var aLevel zap.AtomicLevel

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

type config struct {
	filename   string
	level      string
	maxSize    int
	maxBackups int
	maxAge     int
	compress   bool
}

type OptionFunc func(c *config)

// WithPath alters the saved log full path, for example "./log/debug.log"
func WithPath(filename string) OptionFunc {
	return func(c *config) {
		c.filename = filename
	}
}

// WithLevel alters the log output level, options: debug/info/warn/error
func WithLevel(level string) OptionFunc {
	return func(c *config) {
		c.level = strings.ToLower(level)
	}
}

// WithMaxSize alters the max size of a single log file, default 100Mbit
func WithMaxSize(maxSize int) OptionFunc {
	return func(c *config) {
		c.maxSize = maxSize
	}
}

// WithBackups alters the max num of old log files to retain, default 10
func WithBackups(maxBackups int) OptionFunc {
	return func(c *config) {
		c.maxBackups = maxBackups
	}
}

// WithMaxAge alters the max num of days to retain old log files, default 7 days
func WithMaxAge(maxAge int) OptionFunc {
	return func(c *config) {
		c.maxAge = maxAge
	}
}

// WithCompress determines if the rotated log files should be compressed, default false
func WithCompress(compress bool) OptionFunc {
	return func(c *config) {
		c.compress = compress
	}
}

func Init(opts ...OptionFunc) {
	c := &config{
		filename:   "./log/debug.log",
		level:      "info",
		maxSize:    100,
		maxBackups: 10,
		maxAge:     7,
		compress:   false,
	}

	for _, opt := range opts {
		opt(c)
	}

	aLevel = zap.NewAtomicLevel()
	aLevel.SetLevel(getLogLevel(c.level))

	hook := lumberjack.Logger{
		Filename:   c.filename,
		MaxSize:    c.maxSize,
		MaxBackups: c.maxBackups,
		MaxAge:     c.maxAge,   //days
		Compress:   c.compress, // disabled by default
	}

	fileWriter := zapcore.AddSync(&hook)

	//consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	//core := zapcore.NewTee(
	//	zapcore.NewCore(consoleEncoder, fileWriter, highPriority),
	//)
	//logger = zap.New(core)

	// config := zap.NewProductionEncoderConfig()
	config := zap.NewDevelopmentEncoderConfig()
	config.ConsoleSeparator = " | "

	config.EncodeLevel = levelEncoder
	config.EncodeTime = timeEncoder
	config.EncodeCaller = callerEncoder

	//encoder := zapcore.NewJSONEncoder(config)
	encoder := zapcore.NewConsoleEncoder(config)

	core := zapcore.NewCore(encoder, fileWriter, aLevel)
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func toStr(format string, fmtArgs []interface{}) string {
	msg := format
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(format, fmtArgs...)
	}
	return msg
}

func getLogLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func levelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%5s", l.CapitalString()))
}

func callerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	fullPath := caller.FullPath()

	fullPath, fileName := filepath.Split(fullPath)
	idx := strings.LastIndexByte(fullPath[:len(fullPath)-1], '/')
	if idx >= 0 {
		fullPath = fullPath[idx+1:]
	}

	fullPath = fullPath + fileName
	enc.AppendString(fmt.Sprintf("%20s", fullPath))
}

func SetLogLevel(lvl string) {
	aLevel.SetLevel(getLogLevel(lvl))
}

func Debug(args ...interface{}) {
	logger.Debug(fmt.Sprint(args...))
}

func Debugf(fmt string, args ...interface{}) {
	logger.Debug(toStr(fmt, args))
}

func Info(args ...interface{}) {
	logger.Info(fmt.Sprint(args...))
}

func Infof(fmt string, args ...interface{}) {
	logger.Info(toStr(fmt, args))
}

func Warn(args ...interface{}) {
	logger.Warn(fmt.Sprint(args...))
}

func Warnf(fmt string, args ...interface{}) {
	logger.Warn(toStr(fmt, args))
}

func Error(args ...interface{}) {
	logger.Error(fmt.Sprint(args...))
}
func Errorf(fmt string, args ...interface{}) {
	logger.Error(toStr(fmt, args))
}

func Fatal(args ...interface{}) {
	logger.Fatal(fmt.Sprint(args...))
}
func Fatalf(fmt string, args ...interface{}) {
	logger.Fatal(toStr(fmt, args))
}

func Panic(args ...interface{}) {
	logger.Fatal(fmt.Sprint(args...))
	panic(fmt.Sprint(args...))
}
func Panicf(fmt string, args ...interface{}) {
	logger.Fatal(toStr(fmt, args))
	panic(toStr(fmt, args))
}
