package payment

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	tr "github.com/wopta/goworkspace/transaction"
)

func DeleteTransactionFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var err error

	log.SetPrefix("[DeleteTransactionFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	origin := r.Header.Get("Origin")
	uid := r.Header.Get("uid")
	log.Printf("getting from firestore transaction '%s'", uid)

	transaction := tr.GetTransactionByUid(uid, origin)
	if transaction == nil {
		log.Printf("transaction '%s' not found", uid)
		return "", nil, fmt.Errorf("transaction '%s' not found", uid)
	}
	bytes, _ := json.Marshal(transaction)
	log.Printf("found transaction: %s", string(bytes))

	/*
		switch transaction.ProviderName {
		case models.FabrickPaymentProvider:
			err = fabrickExpireBill(transaction.ProviderId)
		default:
			err = fmt.Errorf("payment provider not implemented: %s", transaction.ProviderName)
		}
		if err != nil {
			log.Printf(">>>>>> error deleting transaction on provider: %s", err.Error())
		}
	*/
	log.Printf("deleting transaction on DBs...")
	if err = tr.DeleteTransaction(transaction, origin, "Cancellata manualmente"); err != nil {
		log.Printf("error deleting transaction on DBs: %s", err.Error())
	} else {
		log.Printf("transaction deleted!")
	}

	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, err
}
