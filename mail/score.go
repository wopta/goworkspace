package mail

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/google/uuid"
	"gitlab.dev.wopta.it/goworkspace/lib"
)

func scoreFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var result map[string]string

	log.AddPrefix("ScoreFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	req := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	json.Unmarshal(req, &result)
	ScoreFido(result["email"])

	log.Println("Handler end -------------------------------------------------")

	return `{"message":"Success send "}`, nil, nil
}

func ScoreFido(data string) <-chan EmailFidoResp {
	r := make(chan EmailFidoResp)
	go func() {
		defer close(r)
		log.Println("ValidateFido")

		var urlstring = "https://api.fido.id/1.0/email"
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		log.Println(getFidoEmailRequest(data))

		req, _ := http.NewRequest(http.MethodPost, urlstring, strings.NewReader(getFidoEmailRequest(data)))
		req.Header.Set("x-api-key", os.Getenv("FIDO_TOKEN_API"))
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		lib.CheckError(err)

		if res != nil {
			body := lib.ErrorByte(io.ReadAll(res.Body))
			defer res.Body.Close()

			var result EmailFidoResp
			json.Unmarshal(body, &result)
			log.Println("body:", string(body))

			r <- result
		}
	}()
	return r
}

func getFidoEmailRequest(data string) string {
	id := uuid.New()
	return `{
		"customer_id": "` + id.String() + `",
		"claims": [
			 "email"
		],
		"email": "` + data + `"
   }`
}

func UnmarshalEmailFidoResp(data []byte) (EmailFidoResp, error) {
	var r EmailFidoResp
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *EmailFidoResp) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type EmailFidoResp struct {
	ResolutionID      string   `json:"resolution_id"`
	CustomerID        string   `json:"customer_id"`
	DeviceRequestTime string   `json:"device_request_time"`
	WebhookURL        string   `json:"webhook_url"`
	Claims            []string `json:"claims"`
	Email             Email    `json:"email"`
}

type Email struct {
	Value                  string `json:"value"`
	Status                 string `json:"status"`
	IsFree                 bool   `json:"is_free"`
	Domain                 string `json:"domain"`
	HasTwitter             bool   `json:"has_twitter"`
	HasAmazon              bool   `json:"has_amazon"`
	HasLinkedin            bool   `json:"has_linkedin"`
	HasAirbnb              bool   `json:"has_airbnb"`
	HasInstagram           bool   `json:"has_instagram"`
	DomainWebsiteExists    bool   `json:"domain_website_exists"`
	FoundOnSerp            bool   `json:"found_on_serp"`
	FirstName              string `json:"first_name"`
	LastName               string `json:"last_name"`
	Education              string `json:"education"`
	Avatar                 string `json:"avatar"`
	AccountLength          string `json:"account_length"`
	AccountDotsCount       string `json:"account_dots_count"`
	AccountNumbersCount    string `json:"account_numbers_count"`
	AccountLettersCount    string `json:"account_letters_count"`
	AccountSymbolsCount    string `json:"account_symbols_count"`
	AccountVowelsCount     string `json:"account_vowels_count"`
	AccountConsonantsCount string `json:"account_consonants_count"`
	AccountDotsRatio       string `json:"account_dots_ratio"`
	AccountNumbersRatio    string `json:"account_numbers_ratio"`
	AccountLettersRatio    string `json:"account_letters_ratio"`
	AccountSymbolsRatio    string `json:"account_symbols_ratio"`
	AccountVowelsRatio     string `json:"account_vowels_ratio"`
	AccountConsonantsRatio string `json:"account_consonants_ratio"`
	Score                  int64  `json:"score"`
	ScoreCluster           string `json:"score_cluster"`
}
