package test

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"

	"cloud.google.com/go/logging"
)

// func TestGetFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
// 	log.SetPrefix("[TestGetFx] ")
// 	defer log.SetPrefix("")
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

	log.SetPrefix("[TestGetFx] ")

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

	log.SetFlags(0)
	log.SetPrefix("")

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
