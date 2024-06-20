package _script

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func TriggerRenew(dryRun bool, appCheckToken string) {
	baseUrl := "https://api.dev.wopta.it/"
	//filePrefix := "dev"
	if os.Getenv("env") == "prod" {
		baseUrl = "https://api.prod.wopta.it/"
		//filePrefix = "prod"
	}

	/*var dates []string
	b, err := os.ReadFile(fmt.Sprintf("./_script/%s-dates.json", filePrefix))
	lib.CheckError(err)
	err = json.Unmarshal(b, &dates)
	lib.CheckError(err)*/

	endpointUrl := baseUrl + "renew/v1/draft?policyType=multiYear&quoteType=fixed"
	for _, date := range []string{"2024-07-03"} {
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
			continue
		}

		fmt.Println("Press the Enter Key to continue")
		fmt.Scanln() // wait for Enter Key
	}

}
