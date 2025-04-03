package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"slices"
	"strings"
	"time"
)

type severityType string

const (
	DEFAULT   severityType = "DEFAULT"
	DEBUG                  = "DEBUG"
	INFO                   = "INFO"
	NOTICE                 = "NOTICE"
	WARNING                = "WARNING"
	ERROR                  = "ERROR"
	CRITICAL               = "CRITICAL"
	ALERT                  = "ALERT"
	EMERGENCY              = "EMERGENCY"
)

type ParserMessage func(string, severityType, []string) ([]byte, error)
type MessageGoogleCloud struct {
	Message  string `json:"message"`
	Severity string `json:"severity,omitempty"`
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
	customLog(string, severityType)
}

type LoggerWopta struct {
	prefix        []string
	writer        io.Writer
	parserMessage ParserMessage
}

// Create a new log, check if use the local parser or the google cloud one
func newLog() *LoggerWopta {
	var parser ParserMessage
	if !IsLocal() {
		parser = parserMessageGoogleCloud
	} else {
		parser = parserMessageLocal
	}
	return &LoggerWopta{
		prefix:        []string{},
		writer:        log.Writer(),
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
func (l *LoggerWopta) AddPrefix(prefix string) {
	l.prefix = append(l.prefix, prefix)
}

// Remove the younger prefix, ex: [prefix1|prefix2] -> [prefix1]
func (l *LoggerWopta) PopPrefix() {
	if len(l.prefix) == 0 {
		return
	}
	l.prefix = slices.Delete(l.prefix, len(l.prefix)-1, len(l.prefix)) //pop of a stack
}

// Remove all prefixs, ex: [prefix1|prefix2] -> <None>
func (l *LoggerWopta) ResetPrefix() {
	l.prefix = []string{}
}

// Set the writer to use for logging
func (l *LoggerWopta) SetWriter(writer io.Writer) {
	l.writer = writer
}

// Log a message with the chosen severity
func (l *LoggerWopta) CustomLog(message string, severity severityType) {
	str, err := l.parserMessage(message, severity, l.prefix)
	if err != nil {
		return
	}
	l.writer.Write(str)
}

// Log a formatted message with severity 'DEFAULT'
func (l *LoggerWopta) Printf(format string, a ...any) {
	l.CustomLog(fmt.Sprintf(format, a...), DEFAULT)
}

// Log a message with severity equal 'DEFAULT'
func (l *LoggerWopta) Println(message string) {
	l.CustomLog(fmt.Sprintf(message), DEFAULT)
}

// Log a formatted message with severity 'INFO'
func (l *LoggerWopta) InfoF(format string, a ...any) {
	l.CustomLog(fmt.Sprintf(format, a...), INFO)
}

// Log a formatted message with severity 'WARNING'
func (l *LoggerWopta) WarningF(format string, a ...any) {
	l.CustomLog(fmt.Sprintf(format, a...), WARNING)
}

// Log a error, with struct : 'Error: <err>'
func (l *LoggerWopta) Error(err error) {
	if err == nil {
		return
	}
	l.CustomLog("Error: "+err.Error(), ERROR)
}

// Log a formatted message with severity 'ERROR'
func (l *LoggerWopta) ErrorF(format string, a ...any) {
	l.CustomLog(fmt.Sprintf(format, a...), ERROR)
}

// Compose the final message with the passed parameters for local debugging
func parserMessageLocal(message string, severity severityType, prefix []string) ([]byte, error) {
	conPrefix := strings.Join(prefix, "|")
	if slices.Contains([]severityType{ERROR, CRITICAL, ALERT, EMERGENCY}, severity) {
		message = "\x1b[49;31m" + message + "\x1b[39;49m"
	} else if slices.Contains([]severityType{WARNING}, severity) {
		message = "\x1b[49;33m" + message + "\x1b[39;49m"
	} else if slices.Contains([]severityType{INFO}, severity) {
		message = "\x1b[49;34m" + message + "\x1b[39;49m"
	}

	conPrefix = fmt.Sprintf("\x1b[;32m [%v] \x1b[39;49m", conPrefix)
	if len(prefix) == 0 {
		conPrefix = " "
	}
	return fmt.Appendf(nil, "%v%v%v\n", formatDate(time.Now()), conPrefix, message), nil
}

// Compose the final message using the given parameters to send to Google Cloud
func parserMessageGoogleCloud(message string, severity severityType, prefix []string) ([]byte, error) {
	conPrefix := strings.Join(prefix, "|")

	conPrefix = fmt.Sprintf(" [%v] ", conPrefix)
	if len(prefix) == 0 {
		conPrefix = " "
	}

	entry := MessageGoogleCloud{
		fmt.Sprintf("%v%v%v", formatDate(time.Now()), conPrefix, message),
		string(severity),
	}
	out, err := json.Marshal(entry)
	if err != nil {
		return []byte{}, err
	}
	return append(out, '\n'), nil
}

var logger *LoggerWopta

// Singleton implementation to get the logger
func Log() *LoggerWopta {
	if logger != nil {
		return logger
	}
	logger = newLog()
	logger.WarningF("INIT LOGGER")
	return logger
}
