// audit_logger.go

package audit

import (
	"fmt"
	"log"
	"os"
	"time"
)

// AuditLog describes the structure of an audit log entry.
type AuditLog struct {
	Timestamp time.Time `json:"timestamp"`
	User      string    `json:"user"`
	Action    string    `json:"action"`
	Details   string    `json:"details"`
}

// Logger is responsible for logging audit entries.
type Logger struct {
	file *os.File
}

// NewLogger creates a new Logger that logs to a specified file.
func NewLogger(filePath string) (*Logger, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return &Logger{file: file}, nil
}

// Log logs a new audit entry.
func (l *Logger) Log(user, action, details string) {
	logEntry := AuditLog{
		Timestamp: time.Now(),
		User:      user,
		Action:    action,
		Details:   details,
	}
	// Log the entry to the file
	if _, err := l.file.WriteString(fmt.Sprintf("%v: %s performed %s - %s\n", logEntry.Timestamp, logEntry.User, logEntry.Action, logEntry.Details)); err != nil {
		log.Panic(err)
	}
}

// Close closes the log file.
func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		log.Panic(err)
	}
}

// Example usage:
// func main() {
//	logger, err := NewLogger("audit.log")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer logger.Close()
//
//	logger.Log("john_doe", "DELETE", "Deleted a resource")
//}