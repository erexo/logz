package logz

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
)

var (
	NotInitializedErr = errors.New("Logz is not initialized")
	InitializedErr    = errors.New("Logz is already initialized")
)

var (
	initialized bool
	lTrace      *log.Logger
	lInfo       *log.Logger
	lWarning    *log.Logger
	lError      *log.Logger
	lCritical   *log.Logger
	logStack    LogLevel
	write       io.Writer
)

type LogLevel byte

const (
	LogLevelTrace LogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
	LogLevelCritical
)

func Init(out io.Writer, logStdLevel, logOutLevel, stackLevel LogLevel, logFileName bool) error {
	if initialized {
		return InitializedErr
	}

	flags := log.LstdFlags
	if logFileName {
		flags |= log.Lshortfile
	}
	lTrace = newLogger(LogLevelTrace, out, logStdLevel, logOutLevel, flags)
	lInfo = newLogger(LogLevelInfo, out, logStdLevel, logOutLevel, flags)
	lWarning = newLogger(LogLevelWarning, out, logStdLevel, logOutLevel, flags)
	lError = newLogger(LogLevelError, out, logStdLevel, logOutLevel, flags)
	lCritical = newLogger(LogLevelCritical, out, logStdLevel, logOutLevel, flags)
	logStack = stackLevel

	write = out
	initialized = true
	return nil
}

func Close() error {
	if !initialized {
		return NotInitializedErr
	}

	r := recover()
	if r != nil {
		var ok bool
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}
		Log(LogLevelCritical, err)
	}
	initialized = false
	if close, ok := write.(io.Closer); ok {
		close.Close()
	}
	if r != nil {
		panic(r) //rethrow from here
	}
	return nil
}

func Log(level LogLevel, v ...interface{}) {
	if !initialized {
		log.Print(NotInitializedErr)
		return
	}
	log := getLogger(level)
	if log != nil {
		log.Output(2, fmt.Sprint(v...))
	}
	if level >= logStack {
		fmt.Fprintln(write, string(debug.Stack()))
	}
}

func Logf(level LogLevel, format string, v ...interface{}) {
	if !initialized {
		log.Print(NotInitializedErr)
		return
	}
	log := getLogger(level)
	if log != nil {
		log.Output(2, fmt.Sprintf(format, v...))
	}
	if level >= logStack {
		fmt.Fprintln(write, string(debug.Stack()))
	}
}

func Trace(v ...interface{}) {
	if !initialized {
		log.Print(NotInitializedErr)
		return
	}
	if lTrace != nil {
		lTrace.Output(2, fmt.Sprint(v...))
	}
	if LogLevelTrace >= logStack {
		fmt.Fprintln(write, string(debug.Stack()))
	}
}

func Tracef(format string, v ...interface{}) {
	if !initialized {
		log.Print(NotInitializedErr)
		return
	}
	if lTrace != nil {
		lTrace.Output(2, fmt.Sprintf(format, v...))
	}
	if LogLevelTrace >= logStack {
		fmt.Fprintln(write, string(debug.Stack()))
	}
}

func Info(v ...interface{}) {
	if !initialized {
		log.Print(NotInitializedErr)
		return
	}
	if lInfo != nil {
		lInfo.Output(2, fmt.Sprint(v...))
	}
	if LogLevelInfo >= logStack {
		fmt.Fprintln(write, string(debug.Stack()))
	}
}

func Infof(format string, v ...interface{}) {
	if !initialized {
		log.Print(NotInitializedErr)
		return
	}
	if lInfo != nil {
		lInfo.Output(2, fmt.Sprintf(format, v...))
	}
	if LogLevelInfo >= logStack {
		fmt.Fprintln(write, string(debug.Stack()))
	}
}

func Warning(v ...interface{}) {
	if !initialized {
		log.Print(NotInitializedErr)
		return
	}
	if lWarning != nil {
		lWarning.Output(2, fmt.Sprint(v...))
	}
	if LogLevelWarning >= logStack {
		fmt.Fprintln(write, string(debug.Stack()))
	}
}

func Warningf(format string, v ...interface{}) {
	if !initialized {
		log.Print(NotInitializedErr)
		return
	}
	if lWarning != nil {
		lWarning.Output(2, fmt.Sprintf(format, v...))
	}
	if LogLevelWarning >= logStack {
		fmt.Fprintln(write, string(debug.Stack()))
	}
}

func Error(v ...interface{}) {
	if !initialized {
		log.Print(NotInitializedErr)
		return
	}
	if lError != nil {
		lError.Output(2, fmt.Sprint(v...))
	}
	if LogLevelError >= logStack {
		fmt.Fprintln(write, string(debug.Stack()))
	}
}

func Errorf(format string, v ...interface{}) {
	if !initialized {
		log.Print(NotInitializedErr)
		return
	}
	if lError != nil {
		lError.Output(2, fmt.Sprintf(format, v...))
	}
	if LogLevelError >= logStack {
		fmt.Fprintln(write, string(debug.Stack()))
	}
}

func Critical(v ...interface{}) {
	if !initialized {
		log.Print(NotInitializedErr)
		return
	}
	if lCritical != nil {
		lCritical.Output(2, fmt.Sprint(v...))
	}
	if LogLevelCritical >= logStack {
		fmt.Fprintln(write, string(debug.Stack()))
	}
	os.Exit(1)
}

func Criticalf(format string, v ...interface{}) {
	if !initialized {
		log.Print(NotInitializedErr)
		return
	}
	if lCritical != nil {
		lCritical.Output(2, fmt.Sprintf(format, v...))
	}
	if LogLevelCritical >= logStack {
		fmt.Fprintln(write, string(debug.Stack()))
	}
	os.Exit(1)
}

func GetLogLevel(str string) LogLevel {
	switch str {
	case "trace":
		return LogLevelTrace
	case "info",
		"information":
		return LogLevelInfo
	case "warning",
		"warn":
		return LogLevelWarning
	case "error":
		return LogLevelError
	case "fatal":
		return LogLevelCritical
	default:
		log.Println("Invalid LogLevel", str)
		return LogLevelTrace
	}
}

func getLogPrefix(level LogLevel) string {
	switch level {
	case LogLevelTrace:
		return "TRACE|"
	case LogLevelInfo:
		return " INFO|"
	case LogLevelWarning:
		return " WARN|"
	case LogLevelError:
		return "ERROR|"
	case LogLevelCritical:
		return "FATAL|"
	default:
		return ""
	}
}

func newLogger(level LogLevel, out io.Writer, logStdLevel, logOutLevel LogLevel, flags int) *log.Logger {
	var w io.Writer
	if level >= logStdLevel {
		w = os.Stdout
	}
	if level >= logOutLevel && out != w {
		if w != nil {
			w = io.MultiWriter(w, out)
		} else {
			w = out
		}
	}
	if w == nil {
		return nil
	}
	return log.New(w, getLogPrefix(level), flags)
}

func getLogger(level LogLevel) *log.Logger {
	switch level {
	case LogLevelTrace:
		return lTrace
	case LogLevelInfo:
		return lInfo
	case LogLevelWarning:
		return lWarning
	case LogLevelError:
		return lError
	case LogLevelCritical:
		return lCritical
	default:
		panic(fmt.Sprintf("Invalid log level %v", level))
	}
}
