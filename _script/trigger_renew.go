package _script

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	env "gitlab.dev.wopta.it/goworkspace/lib/environment"
)

func TriggerRenewDraft(dryRun bool, appCheckToken string) {
	baseUrl := "https://api.dev.wopta.it/"
	if env.IsProduction() {
		baseUrl = "https://api.prod.wopta.it/"
	}

	endpointUrl := baseUrl + "renew/v1/draft?policyType=multiYear&quoteType=fixed"
	for _, date := range []string{} {
		parsedDate, _ := time.Parse(time.DateOnly, date)
		targetDate := parsedDate.AddDate(1, 0, -45).Format(time.DateOnly)

		fmt.Printf("startDate: %s - targetDate: %s\n", date, targetDate)
		body := strings.NewReader(fmt.Sprintf(`{"dryRun": %v, "policyUid":"%s"}`, dryRun, "4kjHg19CnIQfgrLux3QA"))

		req, err := http.NewRequest(http.MethodPost, endpointUrl, body)
		if err != nil {
			panic(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Firebase-Appcheck", appCheckToken)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}

		if res.StatusCode != http.StatusOK {
			log.Printf("Failed to renew draft in date %s due to status code %d", date, res.StatusCode)
		}

		fmt.Println("Press the Enter Key to continue")
		fmt.Scanln() // wait for Enter Key
	}

}

func TriggerRenewPromote(dryRun bool, appCheckToken string) {
	baseUrl := "https://api.dev.wopta.it/"
	if env.IsProduction() {
		baseUrl = "https://api.prod.wopta.it/"
	}

	endpointUrl := baseUrl + "renew/v1/promote"
	for _, date := range []string{} {
		fmt.Printf("date: %s\n", date)
		body := strings.NewReader(fmt.Sprintf(`{"dryRun": %v, "date": "%s"}`, dryRun, date))

		req, err := http.NewRequest(http.MethodPost, endpointUrl, body)
		if err != nil {
			panic(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Firebase-Appcheck", appCheckToken)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}

		if res.StatusCode != http.StatusOK {
			log.Printf("Failed to renew draft in date %s due to status code %d", date, res.StatusCode)
		}

		fmt.Println("Press the Enter Key to continue")
		fmt.Scanln() // wait for Enter Key
	}

}
