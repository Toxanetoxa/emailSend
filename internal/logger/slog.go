package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
)

type File struct {
	file *os.File
}

func NewLogger(filename string) (*File, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &File{file: file}, nil
}

func (l *File) logToFile(level string, methodName string, message string) error {
	log.SetOutput(l.file)
	_, err := fmt.Fprintf(l.file, "[%s] %s - %s: %s\n", time.Now().Format(time.RFC3339), level, methodName, message)
	return err
}

func (l *File) Error(methodName string, err error) {
	_, printErr := color.New(color.FgRed).Printf("[ERROR] %s - %s: %s\n", time.Now().Format(time.RFC3339), methodName, err.Error())
	if printErr != nil {
		return
	}
	if fileErr := l.logToFile("ERROR", methodName, err.Error()); fileErr != nil {
		_, err2 := color.New(color.FgRed).Printf("[ERROR] Failed to write to log file: %s\n", fileErr.Error())
		if err2 != nil {
			return
		}
	}
}

func (l *File) Warning(methodName string, message string) {
	_, printErr := color.New(color.FgYellow).Printf("[WARNING] %s - %s: %s\n", time.Now().Format(time.RFC3339), methodName, message)
	if printErr != nil {
		return
	}
	if fileErr := l.logToFile("WARNING", methodName, message); fileErr != nil {
		_, err := color.New(color.FgRed).Printf("[ERROR] Failed to write to log file: %s\n", fileErr.Error())
		if err != nil {
			return
		}
	}
}

func (l *File) Success(methodName string, message string) {
	_, printErr := color.New(color.FgGreen).Printf("[SUCCESS] %s - %s: %s\n", time.Now().Format(time.RFC3339), methodName, message)
	if printErr != nil {
		return
	}
	if fileErr := l.logToFile("SUCCESS", methodName, message); fileErr != nil {
		_, err := color.New(color.FgRed).Printf("[ERROR] Failed to write to log file: %s\n", fileErr.Error())
		if err != nil {
			return
		}
	}
}

func (l *File) Close() {
	err := l.file.Close()
	if err != nil {
		_ = fmt.Errorf("cannot close Logger")
		return
	}
}
