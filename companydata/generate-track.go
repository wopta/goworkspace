package companydata

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/policy"
	"gitlab.dev.wopta.it/goworkspace/transaction"
)

type GeneratorTrack interface {
	Emit(policy models.Policy, trans *models.Transaction) [][]string
	Delete(policy models.Policy, trans *models.Transaction) [][]string
	Payment(policy models.Policy, trans *models.Transaction) [][]string
}
type GeneratorAxaTrack struct{}

func (_ GeneratorAxaTrack) Emit(policy models.Policy, trans *models.Transaction) [][]string {
	cabCsv := lib.GetFilesByEnv("data/cab-cap-istat.csv")
	df := lib.CsvToDataframe(cabCsv)
	return setRowLifeEmit(policy, df, *trans, time.Now())
}
func (_ GeneratorAxaTrack) Payment(policy models.Policy, trans *models.Transaction) [][]string {
	return [][]string{}
}
func (_ GeneratorAxaTrack) Delete(policy models.Policy, trans *models.Transaction) [][]string {
	return [][]string{}
}

func GenerateTrackFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	transactionUid := chi.URLParam(r, "transactionUid")
	operation := chi.URLParam(r, "operation")
	tr := transaction.GetTransactionByUid(transactionUid, "")
	if tr == nil {
		return "", nil, errors.New("Transaction not found")
	}
	policy := policy.GetPolicyByUid(tr.PolicyUid, "")
	var generator GeneratorTrack

	switch policy.Name {
	case models.LifeProduct:
		generator = GeneratorAxaTrack{}
	default:
		return "", nil, errors.New("Generator track not implement")
	}
	var track [][]string
	switch operation {
	case "emit":
		track = generator.Emit(policy, tr)
	case "delete":
		track = generator.Delete(policy, tr)
	case "Payment":
		track = generator.Payment(policy, tr)
	default:
		return "", nil, errors.New("Generator track not implement for " + operation)
	}

	if len(track) == 0 {
		return "", nil, errors.New("Error generation track")
	}
	bytes, err := lib.GetCsvByte(track, ';')
	if err != nil {
		return "", nil, err
	}
	res := struct{ Bytes []byte }{
		Bytes: bytes,
	}
	respBytes, err := json.Marshal(res)
	return string(respBytes), respBytes, err

}
