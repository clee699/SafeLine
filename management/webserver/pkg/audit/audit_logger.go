package audit

import (
    "log"
    "os"
    "sync"
    "time"
)

type Logger struct {
    mu    sync.Mutex
    file  *os.File
}

func NewLogger(filename string) (*Logger, error) {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        return nil, err
    }

    return &Logger{file: file}, nil
}

func (l *Logger) Log(action, userID string) {
    l.mu.Lock()
    defer l.mu.Unlock()

    logEntry := time.Now().UTC().Format("2006-01-02 15:04:05") + " | User: " + userID + " | Action: " + action + "\n"
    if _, err := l.file.WriteString(logEntry); err != nil {
        log.Println("Error writing to log file:", err)
    }
}

func (l *Logger) Close() {
    l.file.Close()
}