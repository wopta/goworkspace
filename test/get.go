package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/lib/log"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"strings"

	"cloud.google.com/go/logging"
)

// func TestGetFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
// 	log.AddPrefix("[TestGetFx] ")
// 	defer log.PopPrefix()
// 	log.Println("Handler start -----------------------------------------------")

// 	operation := chi.URLParam(r, "operation")

// 	if operation == "error" {
// 		return "", nil, GetErrorJson(401, "Bad Request", "Testing error POST")
// 	}
// 	if operation == "lead" {
// 		p, _ := policy.GetPolicy("DAxF495mw4LH9wwFqL9p", "")
// 		mail.SendMailLead(
// 			p,
// 			mail.AddressAnna,
// 			mail.Address{Address: "diogo.carvalho+test@wopta.it"},
// 			mail.Address{},
// 			models.ECommerceFlow,
// 			[]string{},
// 		)
// 	}
// 	if operation == "sign" {
// 		p, _ := policy.GetPolicy("DAxF495mw4LH9wwFqL9p", "")
// 		mail.SendMailSign(
// 			p,
// 			mail.AddressAnna,
// 			mail.Address{Address: "diogo.carvalho+test@wopta.it"},
// 			mail.Address{},
// 			models.RemittanceMgaFlow,
// 		)
// 	}
// 	if operation == "pay" {
// 		// p, _ := policy.GetPolicy("DAxF495mw4LH9wwFqL9p", "")
// 		p, _ := policy.GetPolicy("DdNLAbEsySpMiDDo07jV", "")
// 		mail.SendMailPay(
// 			p,
// 			mail.AddressAnna,
// 			mail.Address{Address: "diogo.carvalho+test@wopta.it"},
// 			mail.Address{},
// 			models.ProviderMgaFlow,
// 		)
// 	}
// 	if operation == "contract" {
// 		p, _ := policy.GetPolicy("DAxF495mw4LH9wwFqL9p", "")
// 		mail.SendMailContract(
// 			p,
// 			nil,
// 			mail.AddressAnna,
// 			mail.Address{Address: "diogo.carvalho+test@wopta.it"},
// 			mail.Address{},
// 			models.ProviderMgaFlow,
// 		)
// 	}
// 	if operation == "proposal" {
// 		p, _ := policy.GetPolicy("6dk9J1gwIlx9fWKMIufu", "")
// 		mail.SendMailProposal(
// 			p,
// 			mail.AddressAnna,
// 			mail.Address{Address: "diogo.carvalho+test@wopta.it"},
// 			mail.Address{Address: "diogo.carvalho+emittent@wopta.it"},
// 			models.ProviderMgaFlow,
// 			[]string{models.ProposalAttachmentName},
// 		)
// 	}
// 	if operation == "reserved" {
// 		p, _ := policy.GetPolicy("FFjvpy7rgqDw3vu02JzF", "")
// 		mail.SendMailReserved(
// 			p,
// 			mail.AddressAnna,
// 			mail.Address{Address: "diogo.carvalho+test@wopta.it"},
// 			mail.Address{},
// 			models.ProviderMgaFlow,
// 			[]string{models.InformationSetAttachmentName, models.ProposalAttachmentName},
// 		)
// 	}
// 	if operation == "approved" {
// 		p, _ := policy.GetPolicy("FFjvpy7rgqDw3vu02JzF", "")
// 		mail.SendMailReservedResult(
// 			p,
// 			mail.AddressAnna,
// 			mail.Address{Address: "diogo.carvalho+test@wopta.it"},
// 			mail.Address{},
// 			models.ProviderMgaFlow,
// 		)
// 	}
// 	if operation == "rejected" {
// 		p, _ := policy.GetPolicy("FFjvpy7rgqDw3vu02JzF", "")
// 		mail.SendMailReservedResult(
// 			p,
// 			mail.AddressAnna,
// 			mail.Address{Address: "diogo.carvalho+test@wopta.it"},
// 			mail.Address{},
// 			models.ProviderMgaFlow,
// 		)
// 	}

// 	log.Println("Handler end -------------------------------------------------")

// 	return "{}", nil, nil
// }

func TestGetFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	log.Println(Entry{
		Severity: logging.Notice.String(),
		Message:  "This is a notice message",
	})

	log.AddPrefix("[TestGetFx] ")

	log.Println(Entry{
		Severity: logging.Warning.String(),
		Message:  "This is a warning message",
	})

	log.Println(Entry{
		Severity: logging.Error.String(),
		Message:  "This is an error message",
	})

	logger := slog.Default().With(
		slog.String("handler", "TestGetFx"),
	)

	logger.Debug("This is a debug log", "env", os.Getenv("env"))
	logger.Info("This is an info log", "env", os.Getenv("env"))
	logger.Warn("This is a warn log", "env", os.Getenv("env"))
	logger.Error("This is an error log", "env", os.Getenv("env"))

	log.PopPrefix()

	log.Println(Entry{
		Severity: logging.Notice.String(),
		Message:  "This is a notice message",
	})

	log.Println(Entry{
		Severity: logging.Warning.String(),
		Message:  "This is a warning message",
	})

	log.Println(Entry{
		Severity: logging.Error.String(),
		Message:  "This is an error message",
	})

	jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx := context.Background()

	jsonLogger.Log(ctx, slog.Level(logging.Notice), "This is a notice message from slog")
	jsonLogger.Log(ctx, slog.Level(logging.Warning), "This is a Warning message from slog")
	jsonLogger.Log(ctx, slog.Level(logging.Error), "This is a Error message from slog")
	jsonLogger.Log(ctx, slog.Level(logging.Critical), "This is a Critical message from slog")

	myLogger := NewLogger("TestGetFx")

	myLogger.Debug("This is a debug message from my custom logger")
	myLogger.Info("This is an info message from my custom logger")
	myLogger.Warn("This is a warning message from my custom logger")
	myLogger.Error("This is an error message from my custom logger")

	return "{}", nil, nil
}

// Entry defines a log entry.
type Entry struct {
	Message  string `json:"message"`
	Severity string `json:"severity,omitempty"`
	Trace    string `json:"logging.googleapis.com/trace,omitempty"`

	// Logs Explorer allows filtering and display of this as `jsonPayload.component`.
	Component string `json:"component,omitempty"`
}

// String renders an entry structure to the JSON format expected by Cloud Logging.
func (e Entry) String() string {
	if e.Severity == "" {
		e.Severity = "INFO"
	}
	out, err := json.Marshal(e)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}
	return string(out)
}

// /
type Logger interface {
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
}

type DefaultLogger struct {
	Prefix string
}

func (l *DefaultLogger) parse(severity, message string) string {
	_, file, line, _ := runtime.Caller(2)
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	prefix := strings.Join([]string{fmt.Sprintf("%s:%d", file, line), l.Prefix}, " ")

	entry := struct {
		Message  string `json:"message"`
		Severity string `json:"severity,omitempty"`
	}{
		strings.Join([]string{prefix, message}, " | "),
		severity,
	}
	out, err := json.Marshal(entry)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}
	return string(out)
}

func (l *DefaultLogger) Debug(msg string) {
	if msg = l.parse("Debug", msg); msg != "" {
		log.Println(msg)
	}
}

func (l *DefaultLogger) Info(msg string) {
	if msg = l.parse("Info", msg); msg != "" {
		log.Println(msg)
	}
}

func (l *DefaultLogger) Warn(msg string) {
	if msg = l.parse("Warning", msg); msg != "" {
		log.Println(msg)
	}
}

func (l *DefaultLogger) Error(msg string) {
	if msg = l.parse("Error", msg); msg != "" {
		log.Println(msg)
	}
}

func NewLogger(prefix string) Logger {
	return &DefaultLogger{
		Prefix: prefix,
	}
}
