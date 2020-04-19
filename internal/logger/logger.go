package logger

import (
	"log"
	"os"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

// Debug logs a debug message with the given fields
func Debug(message string) {
	log.Printf(DebugColor, message)
}

// Info logs a debug message with the given fields
func Info(message string) {
	log.Printf(InfoColor, message)
}

// Notice logs a debug message with the given fields
func Notice(message string) {
	log.Printf(NoticeColor, message)
}

// Warn logs a debug message with the given fields
func Warn(message string) {
	log.Printf(WarningColor, message)
}

// Error logs a debug message with the given fields
func Error(message string) {
	log.Printf(ErrorColor, message)
}

// Fatal logs a message than calls os.Exit(1)
func Fatal(message string) {
	Error(message)
	os.Exit(1)
}
