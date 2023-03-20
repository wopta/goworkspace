package quote

import (
	"io/ioutil"
	"net/http"

	lib "github.com/wopta/goworkspace/lib"
)

func LifeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	jsonData, err := ioutil.ReadAll(r.Body)

	res := <-Life(jsonData)
	return res, nil, err

}
func Life(r []byte) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		lib.GetFilesByEnv("life_matrix.csv")

		ch <- ""
	}()
	return ch
}
