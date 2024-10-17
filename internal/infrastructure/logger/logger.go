package logger

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	logger := logrus.New()

	// Enable reportCaller to include file and line number in logs
	logger.SetReportCaller(true)

	// Set the custom formatter
	logger.SetFormatter(&CustomFormatter{})

	return logger
}

// CustomFormatter to format log output with uppercase log level, colors, and custom spacing
type CustomFormatter struct {
}

// Format defines how to format each log entry
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buffer bytes.Buffer

	// Log level in uppercase and colored
	levelColor := getColorByLevel(entry.Level)
	level := fmt.Sprintf("%-7s", strings.ToUpper(entry.Level.String()))
	buffer.WriteString(fmt.Sprintf("%s%s\033[0m", levelColor, level))

	// Timestamp in brackets
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	buffer.WriteString(fmt.Sprintf("[%s] ", timestamp))

	// Add caller information (file and line number)
	if entry.HasCaller() {
		filename := filepath.Base(entry.Caller.File)
		buffer.WriteString(fmt.Sprintf("%s:%d: ", filename, entry.Caller.Line))
	}

	// Log message
	buffer.WriteString(entry.Message)

	// End the log entry with a newline
	buffer.WriteString("\n")

	return buffer.Bytes(), nil
}

// getColorByLevel returns a color based on the log level
func getColorByLevel(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return "\033[36m" // Cyan
	case logrus.InfoLevel:
		return "\033[32m" // Green
	case logrus.WarnLevel:
		return "\033[33m" // Yellow
	case logrus.ErrorLevel:
		return "\033[31m" // Red
	case logrus.FatalLevel, logrus.PanicLevel:
		return "\033[35m" // Magenta
	default:
		return "\033[37m" // White (default)
	}
}
