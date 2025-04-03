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

func NewLog() *LoggerWopta {
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

// used for test
func _newLog(isLocal bool) *LoggerWopta {
	log := NewLog()
	if isLocal {
		log.parserMessage = parserMessageLocal
	} else {
		log.parserMessage = parserMessageGoogleCloud
	}
	return log
}

func (l *LoggerWopta) AddPrefix(prefix string) {
	l.prefix = append(l.prefix, prefix)
}

func (l *LoggerWopta) PopPrefix() {
	if len(l.prefix) == 0 {
		return
	}
	l.prefix = slices.Delete(l.prefix, len(l.prefix)-1, len(l.prefix)) //pop of a stack
}

func (l *LoggerWopta) ResetPrefix() {
	l.prefix = []string{}
}

func (l *LoggerWopta) SetLog(writer io.Writer) {
	l.writer = writer
}

func (l *LoggerWopta) CustomLog(message string, severuty severityType) {
	str, err := l.parserMessage(message, severuty, l.prefix)
	if err != nil {
		return
	}
	l.writer.Write(str)

}
func (l *LoggerWopta) Printf(format string, a ...any) {
	l.CustomLog(fmt.Sprintf(format, a...), DEFAULT)
}

func (l *LoggerWopta) Println(format string, a ...any) {
	l.CustomLog(fmt.Sprintf(format, a...), DEFAULT)
}

func (l *LoggerWopta) InfoF(format string, a ...any) {
	l.CustomLog(fmt.Sprintf(format, a...), INFO)
}

func (l *LoggerWopta) WarningF(format string, a ...any) {
	l.CustomLog(fmt.Sprintf(format, a...), WARNING)
}

func (l *LoggerWopta) Error(err error) {
	if err == nil {
		return
	}
	l.CustomLog("Error: "+err.Error(), ERROR)
}

func (l *LoggerWopta) ErrorF(format string, a ...any) {
	l.CustomLog(fmt.Sprintf(format, a...), ERROR)
}

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

func Log() *LoggerWopta {
	if logger != nil {
		return logger
	}
	logger = NewLog()
	logger.WarningF("INIT LOGGER")
	return logger
}
