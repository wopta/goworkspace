package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"time"

	env "gitlab.dev.wopta.it/goworkspace/lib/environment"
)

type SeverityType string

const (
	DEFAULT   SeverityType = "DEFAULT"
	DEBUG                  = "DEBUG"
	INFO                   = "INFO"
	NOTICE                 = "NOTICE"
	WARNING                = "WARNING"
	ERROR                  = "ERROR"
	CRITICAL               = "CRITICAL"
	ALERT                  = "ALERT"
	EMERGENCY              = "EMERGENCY"
)

type ParserMessage func(string, SeverityType, []string) ([]byte, error)

type MessageInformation struct {
	Message      string `json:"message"`
	Severity     string `json:"severity,omitempty"`
	ExecutiondId string `json:"executiondId"`
}

func formatDate(t time.Time) string {
	location, _ := time.LoadLocation("Europe/Rome")
	localTime := t.In(location)
	return localTime.In(location).Format(time.DateTime)
}

type Logger interface {
	AddPrefix(string)
	Warning(string)
	Error(string)
	SetLog(io.Writer)
	SetParserMessage(ParserMessage)
	SetInstanceUid(string)
	customLog(string, SeverityType)
}

type LoggerWopta struct {
	prefix        []string
	writer        io.Writer
	parserMessage ParserMessage
}

// Create a new log, check if use the local parser or the google cloud one
func newLog() *LoggerWopta {
	var parser ParserMessage
	if !isLocal() {
		parser = parserMessageGoogleCloud
	} else {
		parser = parserMessageLocal
	}
	return &LoggerWopta{
		prefix:        []string{},
		writer:        os.Stdout,
		parserMessage: parser,
	}
}

// Create a new log used for test
func _newLog(isLocal bool) *LoggerWopta {
	log := newLog()
	if isLocal {
		log.parserMessage = parserMessageLocal
	} else {
		log.parserMessage = parserMessageGoogleCloud
	}
	return log
}

// Append the prefix, ex: [prefix1] -> [prefix1|prefix2]
// Remember to use PopPrefix to remove eventually
func (l *LoggerWopta) addPrefix(prefix string) {
	l.prefix = append(l.prefix, prefix)
}

// Remove the younger prefix, ex: [prefix1|prefix2] -> [prefix1]
func (l *LoggerWopta) popPrefix() {
	if len(l.prefix) == 0 {
		return
	}
	l.prefix = slices.Delete(l.prefix, len(l.prefix)-1, len(l.prefix)) //pop of a stack
}

// Remove all prefixs, ex: [prefix1|prefix2] -> <None>
func (l *LoggerWopta) resetPrefix() {
	l.prefix = []string{}
}

// Set the writer to use for logging
func (l *LoggerWopta) SetWriter(writer io.Writer) {
	l.writer = writer
}

// Log a message with the chosen severity
func (l *LoggerWopta) CustomLog(message string, severity SeverityType) {
	str, err := l.parserMessage(message, severity, l.prefix)
	if err != nil {
		return
	}
	l.writer.Write(str)
}

// Log a formatted message with severity 'DEFAULT'
func (l *LoggerWopta) printf(format string, a ...any) {
	l.CustomLog(fmt.Sprintf(format, a...), DEFAULT)
}

// Log a message with severity equal 'DEFAULT'
func (l *LoggerWopta) println(message string) {
	l.CustomLog(message, DEFAULT)
}

// Log a formatted message with severity 'INFO'
func (l *LoggerWopta) infoF(format string, a ...any) {
	l.CustomLog(fmt.Sprintf(format, a...), INFO)
}

// Log a formatted message with severity 'WARNING'
func (l *LoggerWopta) warningF(format string, a ...any) {
	l.CustomLog(fmt.Sprintf(format, a...), WARNING)
}

// Log a error, with struct : 'error: <err>'
func (l *LoggerWopta) error(err error) {
	if err == nil {
		return
	}
	l.CustomLog("Error: "+err.Error(), ERROR)
}

// Log a formatted message with severity 'ERROR'
func (l *LoggerWopta) errorF(format string, a ...any) {
	l.CustomLog(fmt.Sprintf(format, a...), ERROR)
}

// Compose the final message with the passed parameters for local debugging
func parserMessageLocal(message string, severity SeverityType, prefix []string) ([]byte, error) {
	conPrefix := strings.Join(prefix, "|")
	if slices.Contains([]SeverityType{ERROR, CRITICAL, ALERT, EMERGENCY}, severity) {
		message = "\x1b[49;31m" + message + "\x1b[39;49m"
	} else if slices.Contains([]SeverityType{WARNING}, severity) {
		message = "\x1b[49;33m" + message + "\x1b[39;49m"
	} else if slices.Contains([]SeverityType{INFO}, severity) {
		message = "\x1b[49;34m" + message + "\x1b[39;49m"
	}

	conPrefix = fmt.Sprintf("\x1b[;32m [%v] \x1b[39;49m", conPrefix)
	if len(prefix) == 0 {
		conPrefix = " "
	}
	return fmt.Append(nil, formatDate(time.Now()), conPrefix, message, "\n"), nil
}

// Compose the final message using the given parameters to send to Google Cloud
func parserMessageGoogleCloud(message string, severity SeverityType, prefix []string) ([]byte, error) {
	conPrefix := strings.Join(prefix, "|")

	conPrefix = fmt.Sprintf(" [%v] ", conPrefix)
	if len(prefix) == 0 {
		conPrefix = " "
	}

	entry := MessageInformation{
		Message:      fmt.Sprint(conPrefix, message),
		Severity:     string(severity),
		ExecutiondId: env.GetExecutionId(),
	}
	out, err := json.Marshal(entry)
	if err != nil {
		return []byte{}, err
	}
	return append(out, '\n'), nil
}
func isLocal() bool {
	return slices.Contains([]string{"local", ""}, os.Getenv("env"))
}

var logger *LoggerWopta

// Singleton implementation to get the logger
func Log() *LoggerWopta {
	if logger != nil {
		return logger
	}
	logger = newLog()
	logger.println("INIT LOGGER")
	return logger
}
