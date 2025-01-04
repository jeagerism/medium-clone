package logger

import (
	"log"
	"os"
)

// Define loggers for different log levels
var (
	ErrorLogger = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger  = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	DebugLogger = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
)

// LogError logs error messages
func LogError(err error) {
	if err != nil {
		ErrorLogger.Println(err)
	}
}

// LogInfo logs informational messages
func LogInfo(message string) {
	InfoLogger.Println(message)
}

// LogDebug logs debug messages
func LogDebug(message string) {
	DebugLogger.Println(message)
}
