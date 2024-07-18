package test

import (
	"log"
	"log/slog"
	"net/http"
	"os"
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
	log.SetPrefix("[TestGetFx] ")

	logger := slog.Default().With(
		slog.String("handler", "TestGetFx"),
	)

	logger.Debug("This is a debug log", "env", os.Getenv("env"))
	logger.Info("This is an info log", "env", os.Getenv("env"))
	logger.Warn("This is a warn log", "env", os.Getenv("env"))
	logger.Error("This is an error log", "env", os.Getenv("env"))

	return "{}", nil, nil
}
