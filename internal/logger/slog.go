package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
)

type LoggerFile struct {
	file *os.File
}

func NewLogger(filename string) (*LoggerFile, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &LoggerFile{file: file}, nil
}

func (l *LoggerFile) logToFile(level string, methodName string, message string) error {
	log.SetOutput(l.file)
	_, err := fmt.Fprintf(l.file, "[%s] %s - %s: %s\n", time.Now().Format(time.RFC3339), level, methodName, message)
	return err
}

func (l *LoggerFile) Error(methodName string, err error) {
	_, err = color.New(color.FgRed).Printf("[ERROR] %s - %s: %s\n", time.Now().Format(time.RFC3339), methodName, err.Error())
	if err != nil {
		return
	}
	if fileErr := l.logToFile("ERROR", methodName, err.Error()); fileErr != nil {
		_, err2 := color.New(color.FgRed).Printf("[ERROR] Failed to write to log file: %s\n", fileErr.Error())
		if err2 != nil {
			return
		}
	}
}

func (l *LoggerFile) Warning(methodName string, message string) {
	_, err := color.New(color.FgYellow).Printf("[WARNING] %s - %s: %s\n", time.Now().Format(time.RFC3339), methodName, message)
	if err != nil {
		return
	}
	if fileErr := l.logToFile("WARNING", methodName, message); fileErr != nil {
		_, err := color.New(color.FgRed).Printf("[ERROR] Failed to write to log file: %s\n", fileErr.Error())
		if err != nil {
			return
		}
	}
}

func (l *LoggerFile) Success(methodName string, message string) {
	_, err := color.New(color.FgGreen).Printf("[SUCCESS] %s - %s: %s\n", time.Now().Format(time.RFC3339), methodName, message)
	if err != nil {
		return
	}
}

func (l *LoggerFile) Close() {
	err := l.file.Close()
	if err != nil {
		_ = fmt.Errorf("cannot closing Logger")
		return
	}
}
