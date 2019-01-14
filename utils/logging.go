package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/rollbar/rollbar-go"
)

const (
	infoLevel     = "INFO"
	warnLevel     = "WARN"
	debugLevel    = "DEBUG"
	errorLevel    = "ERROR"
	criticalLevel = "CRITICAL"
)

var (
	// color pallete map
	colorRed    = "\033[0;31m"
	colorGreen  = "\033[0;32m"
	colorOrange = "\033[0;33m"
	colorBlue   = "\033[0;34m"
	colorPurple = "\033[0;35m"
	colorCyan   = "\033[0;36m"
	colorGray   = "\033[0;37m"
	resetColor  = "\x1b[0m"
)

// Logger a struct for logging
type Logger struct {
	Name      string
	color     bool
	timestamp bool
	reportErr bool
	buf       *bytes.Buffer
	mu        sync.RWMutex
	outWriter io.Writer
	errWriter io.Writer
}

// NewLogger creates new logger
func NewLogger(name string) *Logger {
	var buf bytes.Buffer
	log := &Logger{
		Name:      name,
		color:     true,
		timestamp: true,
		reportErr: false,
		buf:       &buf,
		mu:        sync.RWMutex{},
		outWriter: os.Stdout,
		errWriter: os.Stderr,
	}

	if IsDeployment() {
		log.WithoutColor()
		log.WithoutTimestamp()
		log.WithReport()
	}
	return log
}

// SetWriter sets output writer
func (l *Logger) SetWriter(w io.Writer) {
	l.outWriter = w
	l.errWriter = w
}

// WithColor explicitly turn on colorful features on the log
func (l *Logger) WithColor() *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.color = true
	return l
}

// WithoutColor explicitly turn off colorful features on the log
func (l *Logger) WithoutColor() *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.color = false
	return l
}

// WithoutTimestamp explicitly turn off timestamp features on the log
func (l *Logger) WithoutTimestamp() *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.timestamp = false
	return l
}

// WithReport explicitly turn on error reporting to rollbar
func (l *Logger) WithReport() *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.reportErr = true
	return l
}

func (l *Logger) setPrefix(prefix string) {
	l.buf.WriteString(fmt.Sprintf("[%s - %s", l.Name, prefix))
}

func (l *Logger) getColor(lv string) string {
	switch lv {
	case infoLevel:
		return colorCyan
	case warnLevel:
		return colorOrange
	case debugLevel:
		return colorPurple
	case errorLevel:
		return colorRed
	default:
		return colorBlue
	}
}

func (l *Logger) setColor(lv string) {
	l.buf.WriteString(l.getColor(lv))
}

func (l *Logger) setTimestamp() {
	now := time.Now()

	year, month, day := now.Date()
	l.buf.WriteString(fmt.Sprintf(" %d/%d/%d ", day, int(month), year))
	l.buf.WriteString(fmt.Sprintf("%d:%d:%d", now.Hour(), now.Minute(), now.Second()))
}

// Output write logs
func (l *Logger) Output(lv string, format string, v ...interface{}) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.buf.Reset()

	if l.color {
		l.setColor(lv)
	}

	l.setPrefix(lv)

	if l.timestamp {
		l.setTimestamp()
	}

	l.buf.WriteString("] ")

	var data string
	if format == "" {
		data = fmt.Sprint(v...)
	} else {
		data = fmt.Sprintf(format, v...)
	}
	l.buf.WriteString(data)
	l.buf.WriteString("\n")

	if l.color {
		l.buf.WriteString(string(resetColor))
	}

	writer := l.outWriter
	if lv == errorLevel || lv == criticalLevel {
		writer = l.errWriter
	}
	// Flush buffer to output
	_, err := writer.Write(l.buf.Bytes())
	return err
}

// Info logs messages with INFO level
func (l *Logger) Info(v ...interface{}) {
	l.Output(infoLevel, "", v...)
}

// Infof logs messages with INFO level
func (l *Logger) Infof(format string, v ...interface{}) {
	l.Output(infoLevel, format, v...)
}

// Warn logs messages with WARN level
func (l *Logger) Warn(v ...interface{}) {
	l.Output(warnLevel, "", v...)
}

// Warnf logs messages with WARN level
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Output(warnLevel, format, v...)
}

// Debug logs messages with DEBUG level
func (l *Logger) Debug(v ...interface{}) {
	l.Output(debugLevel, "", v...)
}

// Debugf logs messages with DEBUG level
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Output(debugLevel, format, v...)
}

// Error logs messages with ERROR level
func (l *Logger) Error(v ...interface{}) {
	l.Output(errorLevel, "", v...)
	if l.reportErr {
		rollbar.Error(v...)
	}
}

// Errorf logs messages with ERROR level
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Output(errorLevel, format, v...)

	if l.reportErr {
		rollbar.Errorf(string(errorLevel), format, v...)
	}
}

// Panic logs messages with ERROR level and calls panic
func (l *Logger) Panic(v interface{}) {
	l.Output(criticalLevel, "", v)
	if l.reportErr {
		rollbar.Error(v)
	}
	panic(v)
}

// Critical logs messages with CRITICAL level and calls os.Exit
func (l *Logger) Critical(v ...interface{}) {
	l.Output(criticalLevel, "", v...)
	if l.reportErr {
		rollbar.Critical(v...)
	}
	os.Exit(1)
}

// Criticalf logs messages with CRITICAL level and calls os.Exit
func (l *Logger) Criticalf(format string, v ...interface{}) {
	l.Output(criticalLevel, format, v...)
	if l.reportErr {
		rollbar.Errorf(criticalLevel, format, v...)
	}
	os.Exit(1)
}
