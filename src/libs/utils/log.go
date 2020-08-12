package utils

import (
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

const (
	LOG_LEVEL_DEBUG = iota
	LOG_LEVEL_INFO
	LOG_LEVEL_WARN
	LOG_LEVEL_ERR
	LOG_LEVEL_PANIC
	LOG_LEVEL_FATAL

	TIME_FORMAT_LOG = "06-01-02"
)

var (
	logLevel = LOG_LEVEL_DEBUG
	logDir   = "log"
	rootDir  = ""
)

func InitLog(level int, dir string) {
	logLevel = level
	logDir = dir
	rootDir, _ = os.Getwd()
}

// LogDebug 测试信息(只在debug时需要)
func LogDebug(tag string, v interface{}) {
	if logLevel <= LOG_LEVEL_DEBUG && v != nil {
		go outLog2File(LOG_LEVEL_DEBUG, tag, v)
	}
}

// LogInfo 正常信息(上线之后也需要查看)
func LogInfo(tag string, v interface{}) {
	if logLevel <= LOG_LEVEL_INFO && v != nil {
		go outLog2File(LOG_LEVEL_INFO, tag, v)
	}
}

// LogWarn 警告(警告类错误，但程序可以继续运行)
func LogWarn(tag string, v interface{}) {
	if logLevel <= LOG_LEVEL_WARN && v != nil {
		go outLog2File(LOG_LEVEL_WARN, tag, v)
	}
}

// LogErr 打印err(严重错误，程序不能往下进行)
func LogErr(tag string, v interface{}) {
	if logLevel <= LOG_LEVEL_ERR && v != nil {
		go outLog2File(LOG_LEVEL_ERR, tag, v, string(debug.Stack()))
	}
}

// LogPanic 断掉当前的go(严重错误，直接踢出服务器)
func LogPanic(tag string, v interface{}) {
	if logLevel <= LOG_LEVEL_PANIC && v != nil {
		outLog2File(LOG_LEVEL_PANIC, tag, v, string(debug.Stack()))
	}
}

// Fatal 停止项目，谨慎使用
func LogFatal(tag string, v interface{}) {
	if logLevel <= LOG_LEVEL_FATAL && v != nil {
		outLog2File(LOG_LEVEL_FATAL, tag, v, string(debug.Stack()))
	}
}

// getLogPath
func outLog2File(level int, v ...interface{}) {
	log.Println(v)
	// 还没有初始化
	if len(rootDir) <= 0 || rootDir == "" {
		return
	}
	// 获取logFile
	var file *os.File
	switch level {
	case LOG_LEVEL_DEBUG:
		file = getLogFile("debug")
	case LOG_LEVEL_INFO:
		file = getLogFile("info")
	case LOG_LEVEL_WARN:
		file = getLogFile("warn")
	case LOG_LEVEL_ERR:
		file = getLogFile("err")
	case LOG_LEVEL_PANIC:
		file = getLogFile("panic")
	case LOG_LEVEL_FATAL:
		file = getLogFile("fatal")
	}
	if file == nil {
		return
	}
	defer file.Close()
	// 开始log输出
	logger := log.New(file, "\r\n", log.Ldate|log.Ltime|log.Lshortfile)
	switch level {
	case LOG_LEVEL_DEBUG:
		logger.Println(v)
	case LOG_LEVEL_INFO:
		logger.Println(v)
	case LOG_LEVEL_WARN:
		logger.Println(v)
	case LOG_LEVEL_ERR:
		logger.Println(v)
	case LOG_LEVEL_PANIC:
		logger.Panicln(v)
	case LOG_LEVEL_FATAL:
		logger.Fatalln(v)
	}
}

// getLogFile
func getLogFile(levelName string) *os.File {
	tFormat := time.Now().Format(TIME_FORMAT_LOG)
	fileName := tFormat + "." + levelName + ".log"
	logPath := filepath.Join(rootDir, logDir, fileName)
	// 可读，没有则创建，追加
	file, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	return file
}
