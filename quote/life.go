package quote

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func LifeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	var data models.Policy
	defer r.Body.Close()
	e := json.Unmarshal([]byte(req), &data)
	res := <-Life(data)
	s, e := json.Marshal(res)
	return string(s), nil, e

}
func Life(data models.Policy) <-chan models.Policy {
	ch := make(chan models.Policy)
	go func() {
		defer close(ch)
		birthDate, e := time.Parse("DD-MM-YYYY", data.Contractor.BirthDate)
		lib.CheckError(e)
		year := time.Now().Year() - birthDate.Year()

		b := lib.GetFilesByEnv("quote/life_matrix.csv")
		df := lib.CsvToDataframe(b)
		for _, row := range df.Records() {
			if row[0] == string(year) {

			}
		}

		ch <- data
	}()
	return ch
}
