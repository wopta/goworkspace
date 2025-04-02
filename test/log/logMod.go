package log

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
	"time"
)

type severityType string

const (
	DEFAULT   = "DEFAULT"
	DEBUG     = "DEBUG"
	INFO      = "INFO"
	NOTICE    = "NOTICE"
	WARNING   = "WARNING"
	ERROR     = "ERROR"
	CRITICAL  = "CRITICAL"
	ALERT     = "ALERT"
	EMERGENCY = "EMERGENCY"
)

type ParserMessage func(string, severityType, []string) ([]byte, error)

func formatDate(t time.Time) string {
	location, _ := time.LoadLocation("Europe/Rome")
	localTime := t.In(location)
	return localTime.In(location).Format(time.DateTime)
}

type Log interface {
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

func NewLog() *LoggerWopta {
	var parser ParserMessage = parserMessageGoogleCloud
	// Create a logger to print structured logs formatted as a single line Json to stdout
	if os.Getenv("env") == "local" {
		parser = parserMessageLocal
	}
	return &LoggerWopta{
		prefix:        []string{},
		writer:        log.Writer(),
		parserMessage: parser,
	}
}

func (l *LoggerWopta) AddPrefix(prefix string) {
	l.prefix = append(l.prefix, prefix)
}

func (l *LoggerWopta) RemovePrefix() {
	l.prefix = slices.Delete(l.prefix, len(l.prefix)-1, 1)
}

func (l *LoggerWopta) SetLog(writer io.Writer) {
	l.writer = writer
}

func (l *LoggerWopta) customLog(message string, severuty severityType) {
	str, err := l.parserMessage(message, severuty, l.prefix)
	if err != nil {
		return
	}
	l.writer.Write(str)

}

func (l *LoggerWopta) Warning(message string) {
	l.customLog(message, WARNING)
}

func (l *LoggerWopta) Error(message string) {
	l.customLog(message, ERROR)
}

func parserMessageLocal(message string, severity severityType, prefix []string) ([]byte, error) {
	conPrefix := strings.Join(prefix, "|")
	if slices.Contains([]severityType{ERROR}, severity) {
		message = "\x1b[49;31m" + message + "\x1b[39;49m"
	} else if slices.Contains([]severityType{WARNING}, severity) {
		message = "\x1b[49;33m" + message + "\x1b[39;49m"
	}
	return fmt.Appendf(nil, "%v \x1b[;32m [ %v ] \x1b[39;49m %v \n", formatDate(time.Now()), conPrefix, message), nil
}

func parserMessageGoogleCloud(message string, severity severityType, prefix []string) ([]byte, error) {
	conPrefix := strings.Join(prefix, "|")
	entry := struct {
		Message  string `json:"message"`
		Severity string `json:"severity,omitempty"`
	}{
		fmt.Sprintf("%v [ %v ] %v", formatDate(time.Now()), conPrefix, message),
		string(severity),
	}
	out, err := json.Marshal(entry)
	if err != nil {
		return []byte{}, err
	}
	return append(out, '\n'), nil
}
