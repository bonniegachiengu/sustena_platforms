package utils

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	file, err := os.OpenFile("sustena.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

type SustenaError struct {
	Time    time.Time
	Message string
	File    string
	Line    int
}

func (e *SustenaError) Error() string {
	return fmt.Sprintf("[%s] %s in %s:%d", e.Time.Format(time.RFC3339), e.Message, e.File, e.Line)
}

func NewError(message string) error {
	_, file, line, _ := runtime.Caller(1)
	return &SustenaError{
		Time:    time.Now(),
		Message: message,
		File:    file,
		Line:    line,
	}
}

func LogInfo(message string) {
	InfoLogger.Println(message)
}

func LogError(err error) {
	ErrorLogger.Println(err)
}

func LogDebug(message string) {
	log.Printf("DEBUG: %s\n", message)
}

func LogWarning(message string) {
	log.Printf("WARNING: %s\n", message)
}