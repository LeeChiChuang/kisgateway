package logx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kisgateway/serverlib/conf"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	infoLevel int = iota
	ErrLevel
	ServerLevel
)

const (
	timeFormat = "2006-01-02T15:04:05.000Z07"
	logMode string = "log.mode"
	lv string = "log.level"

	levelAlert  = "alert"
	levelInfo   = "info"
	levelError  = "error"
	levelSevere = "severe"
	levelFatal  = "fatal"
	levelSlow   = "slow"
	levelStat   = "stat"

	console string = "console"
)

var (
	logLevel     uint32
	writeConsole bool
	infoLog      io.WriteCloser
	errLog       io.WriteCloser
	statLog      io.WriteCloser
	serverLog      io.WriteCloser

	once        sync.Once
	confErr = errors.New("config log not set")
)

type logEntry struct {
	Timestamp string `json:"@timestamp"`
	Level     string `json:"level"`
	Duration  string `json:"duration,omitempty"`
	Content   string `json:"content"`
}

func SetUp() {
	mode := conf.GetConf(logMode)
	if mode == "" {
		msg := formatWithCaller(confErr, 3)
		log.Print(msg)
		os.Exit(1)
	}

	setLogLevel(conf.GetConf(lv))
	if mode == console {
		setupConsole()
	}
}

func Info(msg string, v ...interface{})  {
	if !shouldLog(infoLevel) {
		return
	}
	outPut(infoLog, levelInfo, fmt.Sprintf(msg, v...))
}

func outPut(w io.WriteCloser, level string, msg string) {
	logInfo := logEntry{
		Timestamp: getTimestamp(),
		Level:     level,
		Content:   msg,
	}
	outputJson(w, logInfo)
}

func getTimestamp() string {
	return time.Now().Format(timeFormat)
}

func outputJson(writer io.WriteCloser, info logEntry) {
	content, err := json.Marshal(info)
	if err != nil {
		log.Println(err)
	} else {
		_, _ = writer.Write(append(content, '\n'))
	}
}

func shouldLog(level int) bool {
	return true
	//return atomic.LoadUint32(&logLevel) > uint32(level)
}

func setupConsole() {
	once.Do(func() {
		infoLog = NewLoggerWriter(log.New(os.Stdin, "", 0))
		errLog = NewLoggerWriter(log.New(os.Stdin, "", 0))
		statLog = NewLoggerWriter(log.New(os.Stdin, "", 0))
		serverLog = NewLoggerWriter(log.New(os.Stdin, "", 0))
	})
}

func setLogLevel(s string) {
	level, _ := strconv.Atoi(s)
	atomic.StoreUint32(&logLevel, uint32(level))
}

func formatWithCaller(err error, depth int) string {
	var buf strings.Builder

	caller := getCaller(depth)
	if len(caller) > 0 {
		buf.WriteString(caller)
		buf.WriteByte(' ')
	}
	buf.WriteString(err.Error())

	return buf.String()
}

func getCaller(depth int) string {
	var buf strings.Builder

	_, file, line, ok := runtime.Caller(depth)
	if ok {
		buf.WriteString(file)
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(line))
	}

	return buf.String()
}
