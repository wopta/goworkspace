package lib

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"testing"
)

type mockWriter struct {
	sended []byte
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	m.sended = append(m.sended, p...)
	return len(p), nil
}

func TestLocalSendMultipleMessages(t *testing.T) {
	log := _newLog(true)
	mockW := mockWriter{}
	log.SetLog(&mockW)
	messages := make([]string, 100)

	for i := 0; i < 100; i++ {
		message := randStringBytes(100)
		log.Printf(message)
		messages[i] = message
	}

	if !isLocalMessageCorrect(mockW, messages...) {
		t.Fatal("message isnt correct,got:", string(mockW.sended))
	}
}

func TestGoogleCloudSendCustomMessage(t *testing.T) {
	log := _newLog(false)
	mockW := mockWriter{}
	log.SetLog(&mockW)
	messages := make([]MessageGoogleCloud, 100)
	var randI int
	for i := 0; i < len(messages); i++ {
		message := randStringBytes(100)
		randI = rand.Intn(len(sev))
		log.CustomLog(message, sev[randI])
		messages[i] = MessageGoogleCloud{
			Message:  message,
			Severity: string(sev[randI]),
		}
	}

	if !isGoogleMessageCorrect(mockW, messages...) {
		t.Fatal("message isnt correct,got:", string(mockW.sended))
	}
}

func TestSendGoogleMessage(t *testing.T) {
	log := _newLog(false)
	mockW := mockWriter{}
	log.SetLog(&mockW)
	messages := make([]MessageGoogleCloud, 100)
	for i := 0; i < len(messages); i++ {
		message := randStringBytes(100)
		function, s := randFunctionLogSend(log)

		function(message)
		messages[i] = MessageGoogleCloud{
			Message:  message,
			Severity: string(s),
		}
	}

	if !isGoogleMessageCorrect(mockW, messages...) {
		t.Fatal("message isnt correct,got:", string(mockW.sended))
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Generate random string that has lenght equal to n
func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// Check if the writer's written the right messages with right pattern in local
func isLocalMessageCorrect(writer mockWriter, message ...string) bool {
	chunks := strings.Split(string(writer.sended), "\n")
	chunks = chunks[:len(chunks)-1] //we have a '\n' at the very end, so we dont considere it
	if len(message) != len(chunks) {
		return false
	}
	for i, chunk := range chunks {
		patter := fmt.Sprintf(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} .*%v.*`, message[i])
		if ok, _ := regexp.MatchString(patter, chunk); !ok {
			return false
		}
	}
	return true
}

// Check if the writer's written the right messages with right structure in google cloud
func isGoogleMessageCorrect(writer mockWriter, message ...MessageGoogleCloud) bool {
	chunks := strings.Split(string(writer.sended), "\n")
	chunks = chunks[:len(chunks)-1] //we have a '\n' at the very end, so we dont considere it
	if len(message) != len(chunks) {
		return false
	}
	for i, chunk := range chunks {
		patter := fmt.Sprintf(`"message":"\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} .*%v.*","severity":"%v"`, message[i].Message, message[i].Severity)
		if ok, _ := regexp.MatchString(patter, chunk); !ok {
			return false
		}
	}
	return true
}

var sev = []severityType{DEFAULT, DEBUG, INFO, NOTICE, WARNING, ERROR, CRITICAL, ALERT, EMERGENCY}

// Get a random function and his relative Severity type
func randFunctionLogSend(logger *LoggerWopta) (func(string, ...any), severityType) {

	makeStruct := func(fun func(string, ...any), sever severityType) struct {
		f func(string, ...any)
		s severityType
	} {
		return struct {
			f func(string, ...any)
			s severityType
		}{
			f: fun,
			s: sever,
		}
	}
	sev := []struct {
		f func(string, ...any)
		s severityType
	}{
		makeStruct(logger.ErrorF, ERROR),
		makeStruct(logger.Println, DEFAULT),
		makeStruct(logger.InfoF, INFO),
		makeStruct(logger.Printf, DEFAULT),
		makeStruct(logger.WarningF, WARNING),
	}
	i := rand.Intn(len(sev))
	return sev[i].f, sev[i].s
}
