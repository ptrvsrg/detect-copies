package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// CustomTextFormatter is a custom logrus log formatter.
type CustomTextFormatter struct{}

// Format formats the log entry into a custom text format.
func (f *CustomTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	message := fmt.Sprintf("%s %5.5s %s\n",
		entry.Time.Format("2006-01-02 15:04:05.000"), // Date-time
		strings.ToUpper(entry.Level.String()),        // Log level
		entry.Message,                                // Log message
	)

	return []byte(message), nil
}

// Log is a preconfigured logrus Logger instance.
var Log = &logrus.Logger{
	Out:       os.Stdout,              // Log output destination (stdout in this case)
	Level:     logrus.InfoLevel,       // Log level (InfoLevel in this case)
	Formatter: &CustomTextFormatter{}, // Custom log formatter
}
