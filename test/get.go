package test

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/policy"
)

func TestGetFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.SetPrefix("[TestGetFx] ")
	defer log.SetPrefix("")
	log.Println("Handler start -----------------------------------------------")

	operation := chi.URLParam(r, "operation")

	if operation == "error" {
		return "", nil, GetErrorJson(401, "Bad Request", "Testing error POST")
	}
	if operation == "lead" {
		p, _ := policy.GetPolicy("DAxF495mw4LH9wwFqL9p", "")
		mail.SendMailLead(
			p,
			mail.AddressAnna,
			mail.Address{Address: "diogo.carvalho+test@wopta.it"},
			mail.Address{},
			models.ECommerceFlow,
			[]string{},
		)
	}
	if operation == "sign" {
		p, _ := policy.GetPolicy("DAxF495mw4LH9wwFqL9p", "")
		mail.SendMailSign(
			p,
			mail.AddressAnna,
			mail.Address{Address: "diogo.carvalho+test@wopta.it"},
			mail.Address{},
			models.RemittanceMgaFlow,
		)
	}
	if operation == "pay" {
		// p, _ := policy.GetPolicy("DAxF495mw4LH9wwFqL9p", "")
		p, _ := policy.GetPolicy("DdNLAbEsySpMiDDo07jV", "")
		mail.SendMailPay(
			p,
			mail.AddressAnna,
			mail.Address{Address: "diogo.carvalho+test@wopta.it"},
			mail.Address{},
			models.ProviderMgaFlow,
		)
	}
	if operation == "contract" {
		p, _ := policy.GetPolicy("DAxF495mw4LH9wwFqL9p", "")
		mail.SendMailContract(
			p,
			nil,
			mail.AddressAnna,
			mail.Address{Address: "diogo.carvalho+test@wopta.it"},
			mail.Address{},
			models.ProviderMgaFlow,
		)
	}
	if operation == "proposal" {
		p, _ := policy.GetPolicy("6dk9J1gwIlx9fWKMIufu", "")
		mail.SendMailProposal(
			p,
			mail.AddressAnna,
			mail.Address{Address: "diogo.carvalho+test@wopta.it"},
			mail.Address{Address: "diogo.carvalho+emittent@wopta.it"},
			models.ProviderMgaFlow,
			[]string{models.ProposalAttachmentName},
		)
	}
	if operation == "reserved" {
		p, _ := policy.GetPolicy("FFjvpy7rgqDw3vu02JzF", "")
		mail.SendMailReserved(
			p,
			mail.AddressAnna,
			mail.Address{Address: "diogo.carvalho+test@wopta.it"},
			mail.Address{},
			models.ProviderMgaFlow,
			[]string{models.InformationSetAttachmentName, models.ProposalAttachmentName},
		)
	}
	if operation == "approved" {
		p, _ := policy.GetPolicy("FFjvpy7rgqDw3vu02JzF", "")
		mail.SendMailReservedResult(
			p,
			mail.AddressAnna,
			mail.Address{Address: "diogo.carvalho+test@wopta.it"},
			mail.Address{},
			models.ProviderMgaFlow,
		)
	}
	if operation == "rejected" {
		p, _ := policy.GetPolicy("FFjvpy7rgqDw3vu02JzF", "")
		mail.SendMailReservedResult(
			p,
			mail.AddressAnna,
			mail.Address{Address: "diogo.carvalho+test@wopta.it"},
			mail.Address{},
			models.ProviderMgaFlow,
		)
	}

	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, nil
}
