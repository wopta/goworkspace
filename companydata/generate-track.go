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
	Emitted(policy models.Policy, trans *models.Transaction) [][]string
	Deleted(policy models.Policy, trans *models.Transaction) [][]string
	Paid(policy models.Policy, trans *models.Transaction) [][]string
}
type GeneratorAxaTrack struct{}

func (_ GeneratorAxaTrack) Emitted(policy models.Policy, trans *models.Transaction) (result [][]string) {
	return [][]string{}
}
func (_ GeneratorAxaTrack) Paid(policy models.Policy, trans *models.Transaction) (result [][]string) {
	cabCsv := lib.GetFilesByEnv("data/cab-cap-istat.csv")
	df := lib.CsvToDataframe(cabCsv)
	result = append(result, getHeader())
	result = append(result, setRowLifeEmit(policy, df, *trans, time.Now())...)
	return result
}
func (_ GeneratorAxaTrack) Deleted(policy models.Policy, trans *models.Transaction) [][]string {
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
	case "emitted":
		track = generator.Emitted(policy, tr)
	case "deleted":
		track = generator.Deleted(policy, tr)
	case "paid":
		track = generator.Paid(policy, tr)
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
	res := struct {
		Bytes []byte `json:"bytes"`
	}{
		Bytes: bytes,
	}
	respBytes, err := json.Marshal(res)
	return string(respBytes), respBytes, err

}
