package Sakura

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strconv"
	"strings"
)

var LogLevel = map[string]logrus.Level{
	"PANIC": logrus.PanicLevel,
	"FATAL": logrus.FatalLevel,
	"ERROR": logrus.ErrorLevel,
	"WARN":  logrus.WarnLevel,
	"INFO":  logrus.InfoLevel,
	"DEBUG": logrus.DebugLevel,
	"TRACE": logrus.TraceLevel,
}

var LevelArray = []string{
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	"P",
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	"F",
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	"E",
	// WarnLevel level. Non-critical entries that deserve eyes.
	"W",
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	"I",
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	"D",
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	"T",
}

var Log = logrus.New()
var buildMode = "DEV"

func InitLogger(__buildMode string, f func()) {
	buildMode = __buildMode
	Log.Out = os.Stdout
	if GetConfig(false).LoggerLevel == nil {
		Log.Level = logrus.DebugLevel
	} else {
		Log.Level = LogLevel[strings.ToUpper(*GetConfig(false).LoggerLevel)]
	}
	formatter := &Formatter{}
	Log.SetReportCaller(true)
	Log.SetFormatter(formatter)
	Log.SetOutput(formatter)
	f()
}

func fieldParse(obj interface{}) string {
	var ret string
	switch v := obj.(type) {
	case string:
		ret = v
	case float64:
		ret = strconv.FormatFloat(v, 'E', -1, 64)
	case int:
		ret = strconv.Itoa(v)
	case int64:
		ret = strconv.FormatInt(v, 0x0a)
	case nil:
		ret = "null"
	default:
		ret = "Unsupported"
	}
	return ret
}

// Formatter 自定义 formatter
type Formatter struct{}

// Formatter.Write implement the Output Writer interface
func (f *Formatter) Write(p []byte) (n int, err error) {
	if buildMode == "REL" {
		switch p[1] {
		case 'P':
			color.Red(string(p))
		case 'F':
			color.Red(string(p))
		case 'E':
			color.Red(string(p))
		case 'W':
			color.Yellow(string(p))
		case 'I':
			color.Blue(string(p))
		case 'D':
			color.Magenta(string(p))
		case 'T':
			color.White(string(p))
		}
	} else {
		switch p[1] {
		case 'P':
			fmt.Printf("\x1b[%dm"+string(p)+" \x1b[0m\n", 31)
		case 'F':
			fmt.Printf("\x1b[%dm"+string(p)+" \x1b[0m\n", 31)
		case 'E':
			fmt.Printf("\x1b[%dm"+string(p)+" \x1b[0m\n", 31)
		case 'W':
			fmt.Printf("\x1b[%dm"+string(p)+" \x1b[0m\n", 33)
		case 'I':
			fmt.Printf("\x1b[%dm"+string(p)+" \x1b[0m\n", 32)
		case 'D':
			fmt.Printf("\x1b[%dm"+string(p)+" \x1b[0m\n", 35)
		case 'T':
			fmt.Printf("\x1b[%dm"+string(p)+" \x1b[0m\n", 29)
		}
	}
	return 0, nil
}

// Format implement the Formatter interface
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	remained := len(entry.Data)
	if remained > 0 {
		entry.Message += " ["
	}
	for k, v := range entry.Data {
		entry.Message += k + ":" + fieldParse(v)
		remained--
		if remained != 0 {
			entry.Message += ", "
		} else {
			entry.Message += "]"
		}
	}
	filename := path.Base(entry.Caller.File)
	b.WriteString(fmt.Sprintf("[%s] %s [%s:%d] %s", LevelArray[entry.Level], entry.Time.Format("2006-01-02 15:04:05"), filename, entry.Caller.Line, entry.Message))
	return b.Bytes(), nil
}

func Blue(msg string) {
	if buildMode == "REL" {
		color.Blue(msg)
	} else {
		fmt.Printf("\x1b[%dm"+msg+" \x1b[0m\n", 34)
	}
}
