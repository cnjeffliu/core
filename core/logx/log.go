/*
* 日志组件封装
* jeff.liu <zhifeng172@163.com> 2019.01.26
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

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
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

var aLevel zap.AtomicLevel

func InitLog(filename string) {
	aLevel = zap.NewAtomicLevel()
	aLevel.SetLevel(getLogLevel("debug"))

	hook := lumberjack.Logger{
		Filename:   filename,
		MaxSize:    30, // megabytes
		MaxBackups: 10,
		MaxAge:     7,    //days
		Compress:   true, // disabled by default
	}

	fileWriter := zapcore.AddSync(&hook)

	//consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	//core := zapcore.NewTee(
	//	// 打印在文件中
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
