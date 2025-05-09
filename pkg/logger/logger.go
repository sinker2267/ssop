package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
	debugLogger *log.Logger
)

// LogLevel 日志级别
type LogLevel int

const (
	// DEBUG 调试级别
	DEBUG LogLevel = iota
	// INFO 信息级别
	INFO
	// WARN 警告级别
	WARN
	// ERROR 错误级别
	ERROR
	// FATAL 致命错误级别
	FATAL
)

var currentLevel LogLevel

// InitLogger 初始化日志
func InitLogger(level string) {
	// 设置日志级别
	switch level {
	case "debug":
		currentLevel = DEBUG
	case "info":
		currentLevel = INFO
	case "warn":
		currentLevel = WARN
	case "error":
		currentLevel = ERROR
	case "fatal":
		currentLevel = FATAL
	default:
		currentLevel = INFO
	}

	// 初始化日志记录器
	debugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	warnLogger = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	fatalLogger = log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Debug 记录调试日志
func Debug(msg string, keyvals ...interface{}) {
	if currentLevel <= DEBUG {
		debugLogger.Println(formatLog(msg, keyvals...))
	}
}

// Info 记录信息日志
func Info(msg string, keyvals ...interface{}) {
	if currentLevel <= INFO {
		infoLogger.Println(formatLog(msg, keyvals...))
	}
}

// Warn 记录警告日志
func Warn(msg string, keyvals ...interface{}) {
	if currentLevel <= WARN {
		warnLogger.Println(formatLog(msg, keyvals...))
	}
}

// Error 记录错误日志
func Error(msg string, keyvals ...interface{}) {
	if currentLevel <= ERROR {
		errorLogger.Println(formatLog(msg, keyvals...))
	}
}

// Fatal 记录致命错误日志并退出程序
func Fatal(msg string, keyvals ...interface{}) {
	if currentLevel <= FATAL {
		fatalLogger.Println(formatLog(msg, keyvals...))
		os.Exit(1)
	}
}

// formatLog 格式化日志
func formatLog(msg string, keyvals ...interface{}) string {
	if len(keyvals) == 0 {
		return msg
	}

	formatted := msg
	for i := 0; i < len(keyvals); i += 2 {
		key := keyvals[i]
		var value interface{} = "MISSING"
		if i+1 < len(keyvals) {
			value = keyvals[i+1]
		}
		formatted += " " + key.(string) + "=" + stringify(value)
	}
	return formatted
}

// stringify 将任意值转为字符串
func stringify(v interface{}) string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprint(v)
} 