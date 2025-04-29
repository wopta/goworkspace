package log

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
	m.sended = p
	return len(p), nil
}

func TestLocalSendMultipleMessages(t *testing.T) {
	log := _newLog(true)
	mockW := mockWriter{}
	log.SetWriter(&mockW)

	for i := 0; i < 100; i++ {
		message := randStringBytes(100)
		log.Printf(message)
		if !isLocalMessageCorrect(mockW, message) {
			t.Fatal("message isnt correct,got:", string(mockW.sended))
		}
	}

}

func TestGoogleCloudSendCustomMessage(t *testing.T) {
	log := _newLog(false)
	mockW := mockWriter{}
	log.SetWriter(&mockW)
	var randI int
	for i := 0; i < 100; i++ {
		message := randStringBytes(100)
		randI = rand.Intn(len(sev))
		log.CustomLog(message, sev[randI])
		messageS := MessageInformation{
			Message:  message,
			Severity: string(sev[randI]),
		}
		if !isGoogleMessageCorrect(mockW, messageS) {
			t.Fatal("message isnt correct,got:", string(mockW.sended))
		}
	}

}

func TestSendGoogleMessage(t *testing.T) {
	mockW := mockWriter{}
	log := _newLog(false)
	log.SetWriter(&mockW)
	for i := 0; i < 100; i++ {
		message := randStringBytes(100)
		function, s := randFunctionToLog(log)

		function(message)
		messagesS := MessageInformation{
			Message:  message,
			Severity: string(s),
		}
		if !isGoogleMessageCorrect(mockW, messagesS) {
			t.Fatal("message isnt correct,got:", string(mockW.sended))
		}
	}

}
func TestAddPopPrefix(t *testing.T) {
	mockW := mockWriter{}
	log := _newLog(false)
	log.SetWriter(&mockW)
	log.AddPrefix("prefix1")
	log.AddPrefix("prefix2")

	log.Println("fjklsd")
	if ok, _ := regexp.MatchString(fmt.Sprintf("%v", "prefix1|prefix2"), string(mockW.sended)); !ok {
		t.Fatal("no prefix founded")
	}

	log.PopPrefix()
	log.Println("fjklsd")
	if ok, _ := regexp.MatchString(fmt.Sprintf("%v", `prefix1`), string(mockW.sended)); !ok {
		t.Fatal("no prefix founded")
	}
	if ok, _ := regexp.MatchString(fmt.Sprintf("%v", `prefix2`), string(mockW.sended)); ok {
		t.Fatal("prefix founded when it shouldn't")
	}

	log.PopPrefix()
	log.Println("fjklsd")
	if ok, _ := regexp.MatchString(fmt.Sprintf("%v", `prefix1`), string(mockW.sended)); ok {
		t.Fatal("prefix founded when it shouldn't")
	}

	log.PopPrefix()
	log.PopPrefix()
	log.PopPrefix()
	log.PopPrefix()
	log.Println("fjklsd")
	if ok, _ := regexp.MatchString(fmt.Sprintf("%v", `prefix1`), string(mockW.sended)); ok {
		t.Fatal("prefix founded when it shouldn't")
	}
	if ok, _ := regexp.MatchString(fmt.Sprintf("%v", `prefix2`), string(mockW.sended)); ok {
		t.Fatal("prefix founded when it shouldn't")
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
func isGoogleMessageCorrect(writer mockWriter, message ...MessageInformation) bool {
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

var sev = []SeverityType{DEFAULT, DEBUG, INFO, NOTICE, WARNING, ERROR, CRITICAL, ALERT, EMERGENCY}

// Get a random function and his relative Severity type
func randFunctionToLog(logger *LoggerWopta) (func(string, ...any), SeverityType) {

	makeStruct := func(fun func(string, ...any), sever SeverityType) struct {
		f func(string, ...any)
		s SeverityType
	} {
		return struct {
			f func(string, ...any)
			s SeverityType
		}{
			fun,
			sever,
		}
	}
	sev := []struct {
		f func(string, ...any)
		s SeverityType
	}{
		makeStruct(logger.ErrorF, ERROR),
		makeStruct(logger.InfoF, INFO),
		makeStruct(logger.Printf, DEFAULT),
		makeStruct(logger.WarningF, WARNING),
	}
	i := rand.Intn(len(sev))
	return sev[i].f, sev[i].s
}
