package quote

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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
		var urlstring = os.Getenv("MUNICHREBASEURL") + "/api/quote/rate/"
		client := lib.ClientCredentials(os.Getenv("MUNICHRECLIENTID"),
			os.Getenv("MUNICHRECLIENTSECRET"), os.Getenv("MUNICHRESCOPE"), os.Getenv("MUNICHRETOKENENDPOINT"))
		req, _ := http.NewRequest(http.MethodPost, urlstring, bytes.NewBuffer(r))
		req.Header.Set("Ocp-Apim-Subscription-Key", os.Getenv("MUNICHRESUBSCRIPTIONKEY"))
		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)
		lib.CheckError(err)

		if res != nil {
			body, err := ioutil.ReadAll(res.Body)
			lib.CheckError(err)
			res.Body.Close()
			log.Println("quote res")
			ch <- string(body)
			if res.StatusCode == 500 {
				log.Println("StatusCode == 500")

			}

		}

	}()
	return ch
}
