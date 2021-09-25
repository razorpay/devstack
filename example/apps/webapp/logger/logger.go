//Package logger provides initializes logger and external logging service hooks
package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

var (
	L = logrus.New()
)

type JaegerLogger struct{}

func (l *JaegerLogger) Error(msg string) {
	L.Printf("ERROR: %s", msg)
}

// Infof logs a message at info priority
func (l *JaegerLogger) Infof(msg string, args ...interface{}) {
	L.Printf(msg, args...)
}

// Debugf logs a message at debug priority
func (l *JaegerLogger) Debugf(msg string, args ...interface{}) {
	L.Printf(fmt.Sprintf("DEBUG: %s", msg), args...)
}
