package payment

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/wopta/goworkspace/models"
)

func fabrickExpireBill(providerId string) error {
	log.Println("starting fabrick expire bill request...")
	var urlstring = os.Getenv("FABRICK_BASEURL") + "api/fabrick/pace/v4.0/mods/back/v1.0/payments/expirationDate"
	const expirationTimeSuffix = "00:00:00"

	expirationDate := fmt.Sprintf(
		"%s %s",
		time.Now().UTC().AddDate(0, 0, -1).Format(models.TimeDateOnly),
		expirationTimeSuffix,
	)
	requestBody := fmt.Sprintf(`{"id":"%s","newExpirationDate":"%s"}`, providerId, expirationDate)
	log.Printf("fabrick expire bill request body: %s", requestBody)

	req, err := http.NewRequest(http.MethodPut, urlstring, strings.NewReader(requestBody))
	if err != nil {
		log.Printf("error creating request: %s", err.Error())
		return err
	}
	res, err := getFabrickClient(urlstring, req)
	if err != nil {
		log.Printf("error getting response: %s", err.Error())
		return err
	}

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("fabrick expire bill response error: %s", err.Error())
		return err
	}
	log.Println("fabrick expire bill response body: ", string(respBody))
	if res.StatusCode != http.StatusOK {
		log.Printf("fabrick expire bill error status %s", res.Status)
		return fmt.Errorf("fabrick expire bill error status %s", res.Status)
	}

	log.Println("fabrick expire bill completed!")

	return nil
}
