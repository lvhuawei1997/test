package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
)

type Logger struct {
	file *os.File
	mu   sync.RWMutex
}

func NewLogger(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	return &Logger{file: file}, nil
}

func (l *Logger) Close() {
	l.mu.RLock()
	defer l.mu.RUnlock()

	l.file.Close()
}

func (l *Logger) readLog() ([]string, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var logEntries []string
	file, err := os.Open(l.file.Name())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		logEntries = append(logEntries, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return logEntries, nil
}

func processLog(logText string) (string, string) {
	// 定义匹配的日期和错误信息的正则表达式
	dateRegex := `\[(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\]`
	errorRegex := `(\w+): (.*)`

	// 编译正则表达式
	datePattern := regexp.MustCompile(dateRegex)
	errorPattern := regexp.MustCompile(errorRegex)

	// 使用正则表达式匹配日志文本
	dateMatch := datePattern.FindStringSubmatch(logText)
	errorMatch := errorPattern.FindStringSubmatch(logText)

	// 提取匹配到的信息
	var date, errorMessage string
	if len(dateMatch) > 1 {
		date = dateMatch[1]
	}
	if len(errorMatch) > 1 {
		errorMessage = errorMatch[1]
	}
	return date, errorMessage
}

func showLog(logEntries []string) {
	for _, entry := range logEntries {
		date, errorMessage := processLog(entry)
		fmt.Printf("Log Entry: %s\n", entry)
		fmt.Printf("Date: %s, Level: %s\n", date, errorMessage)
		fmt.Println()
		fmt.Println("-----------")
	}
}

func main() {
	logger, err := NewLogger("/Users/ke/test/test/main/example.log")
	if err != nil {
		log.Fatalf("Error opening log file: %v\n", err)
	}
	defer logger.Close()

	// 读取日志并展示
	logEntries, err := logger.readLog()
	if err != nil {
		log.Printf("Error reading log file: %v\n", err)
	} else {
		showLog(logEntries)
	}
}
