package callback

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib"
	env "gitlab.dev.wopta.it/goworkspace/lib/environment"
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
			log.WarningF("For policy %v there is the error %v", policies[i], e.Error())
			errors = append(errors, "Errore incasso: "+policies[i].CodeCompany)
			continue
		} else {
			policies[i].AddSystemNote(func(p *models.Policy) models.PolicyNote {
				return models.PolicyNote{
					Text:               "Incasso eseguito nei sistemi di net-insurance",
					ReadableByProducer: false,
				}

			})
		}
	}
	updateAllPolices(policies)
	if len(errors) == 0 {
		log.Println("All policy incassate in net-insurance system")
		return "{}", "", nil
	}
	sendEmailErrorIncasso(len(policies), errors)
	return "{}", "", fmt.Errorf("Policies with errors: %v", len(errors))
}

func updateAllPolices(policies []models.Policy) {
	log.Println("Setting CompanyEmitted=true")
	wait := sync.WaitGroup{}
	logString := ""
	wait.Add(len(policies))
	for i := range policies {
		go func() {
			policies[i].CompanyEmitted = true
			policies[i].Updated = time.Now().UTC()
			err := lib.SetFirestoreErr(lib.PolicyCollection, policies[i].Uid, &policies[i])
			if err != nil {
				log.Error(err)
			}

			policies[i].BigquerySave()
			logString += fmt.Sprintf("Policy %v saved \n", policies[i].CodeCompany)
			wait.Done()
		}()
	}
	wait.Wait()
	log.Println(logString)
}

func sendEmailErrorIncasso(nPolicy int, errors []string) {

	subject := "Polizze net-insurance non incassate!"
	lines := []string{
		"Ciao,",
		fmt.Sprintf("Nella data del %v sono state quietanzate %v polizze cat-nat, %v delle quali hanno avuto problemi, gli errori rilevati sono:<br><br>", time.Now().Format("2006-01-02"), nPolicy, len(errors)),
		strings.Join(errors, "<br>"),
	}
	var body string
	for _, line := range lines {
		body += body + `<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px">` + line + `</p>`
	}
	body += "<br>Si prega di incassare manualemnte le polizze sopra riportate<br>"
	body += fmt.Sprintf("<h6>Execution Id: %v </h6>", env.GetExecutionId())
	body += ` 
		<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px"><br></p><p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px">A presto,</p>
		<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#e50075;font-size:14px"><strong>Anna</strong> di Wopta Assicurazioni</p> `

	mail.SendBaseEmail(body, subject, lib.GetMailProcessi("cat-nat"))
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
