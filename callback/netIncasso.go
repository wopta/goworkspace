package callback

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/models/catnat"
)

func incassoNetFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	log.AddPrefix("IncassoNetFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")
	policies, e := getPolicyToCallNetIncasso()
	if e != nil {
		return "", nil, e
	}
	log.Println("Policies to incassare in net-insurance system: ", len(policies))
	catnatClient := catnat.NewNetClient()
	var errors []string
	for i := range policies {
		if e := catnatClient.Incasso(policies[i]); e != nil {
			log.ErrorF("For policy %v there is the error %v", policies[i], e.Error())
			errors = append(errors, "Errore incasso: "+policies[i].CodeCompany)
			continue
		}
	}
	updateAllPolices(policies)
	if len(errors) == 0 {
		log.Println("All policy incassate in net-insurance system")
		return "{}", "", nil
	}
	log.Println("Policies with errors: ", len(errors))
	sendEmailErrorIncasso(len(policies), errors)
	return "{}", "", nil
}

func updateAllPolices(policies []models.Policy) {
	log.Println("Setting CompanyEmitted=true")
	wait := sync.WaitGroup{}
	wait.Add(len(policies))
	for i := range policies {
		go func() {
			policies[i].CompanyEmitted = true
			policies[i].Updated = time.Now().UTC()
			log.Println("saving to firestore...")
			err := lib.SetFirestoreErr(lib.PolicyCollection, policies[i].Uid, &policies[i])
			if err != nil {
				log.Error(err)
			}
			log.Println("firestore saved!")

			policies[i].BigquerySave()
			log.Printf("Policy %v saved", policies[i].CodeCompany)
			wait.Done()
		}()
	}
	wait.Wait()
}

func sendEmailErrorIncasso(nPolicy int, errors []string) {
	var mailRequest mail.MailRequest
	mailRequest.IsHtml = true
	mailRequest.FromAddress = mail.AddressAnna
	mailRequest.To = []string{"processi@wopta.it"}
	mailRequest.Subject = "Polizze net-insurance non incassate!"
	lines := []string{
		"Ciao,",
		fmt.Sprintf("Nella data del %v sono state quietanzate %v polizze cat-nat, %v delle quali hanno avuto problemi, gli errori rilevati sono:<br><br>", time.Now().Format("2006-01-02"), nPolicy, len(errors)),
		strings.Join(errors, "<br>"),
	}
	for _, line := range lines {
		mailRequest.Message = mailRequest.Message + `<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px">` + line + `</p>`
	}
	mailRequest.Message += "<br>Si prega di incassare manualemnte le polizze sopra riportate<br>"
	mailRequest.Message = mailRequest.Message + ` 
		<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px"><br></p><p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px">A presto,</p>
		<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#e50075;font-size:14px"><strong>Anna</strong> di Wopta Assicurazioni</p> `
	mail.SendMail(mailRequest)
}

func getPolicyToCallNetIncasso() ([]models.Policy, error) {
	catnatPolicyToEmit := lib.Firequeries{
		Queries: []lib.Firequery{
			{
				Field:      "companyEmitted",
				Operator:   "==",
				QueryValue: false,
			},
			{
				Field:      "isDeleted",
				Operator:   "==",
				QueryValue: false,
			},
			{
				Field:      "isPay",
				Operator:   "==",
				QueryValue: true,
			},
			{
				Field:      "status",
				Operator:   "==",
				QueryValue: models.PolicyStatusPay,
			},
			{
				Field:      "name",
				Operator:   "==",
				QueryValue: models.CatNatProduct,
			},
		},
	}
	docsnap, e := catnatPolicyToEmit.FirestoreWherefields(lib.PolicyCollection)
	if e != nil {
		return nil, e
	}
	return models.PolicyToListData(docsnap), nil
}
