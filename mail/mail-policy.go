package mail

import "github.com/wopta/goworkspace/models"

func GetMailPolicy(policy models.Policy, subject string, lines ...string) MailRequest {
	var name string
	//var linkForm string
	if policy.Name == "pmi" {
		name = "Artigiani & Imprese"
		//linkForm = "https://www.wopta.it/it/multi-rischio/"
	}
	var obj MailRequest

	obj.From = "noreply@wopta.it"
	obj.To = []string{policy.Contractor.Mail}
	obj.Message = `<p></p><p>Gentile ` + policy.Contractor.Name + ` ` + policy.Contractor.Surname + obj.Message
	for _, line := range lines {
		obj.Message = obj.Message + line
	}

	obj.Message = obj.Message + `<p>Un saluto.</p><p>ll Team Wopta. Proteggiamo chi sei</p>`
	obj.Subject = "Wopta per te. " + name + " " + subject + " n° " + policy.CodeCompany
	obj.IsHtml = true
	obj.IsAttachment = false

	return obj
}
