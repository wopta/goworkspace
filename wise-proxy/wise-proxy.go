package rules

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	// Blank-import the function package so the init() runs
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
)

func init() {
	log.Println("INIT WiseProxy")

	functions.HTTP("WiseProxy", WiseProxy)
}

func WiseProxy(w http.ResponseWriter, r *http.Request) {
	log.Println("WiseProxy")
	log.Println(r.RequestURI)
	lib.EnableCors(&w, r)
	client := http.Client{Timeout: time.Duration(1) * time.Second}
	if strings.Contains(r.RequestURI, "/WebApiProduct") {
		var urlstring = "http://test-wopta.northeurope.cloudapp.azure.com/" + r.RequestURI
		jsonData, _ := ioutil.ReadAll(r.Body)
		login := strings.NewReader("{\"username\": \"" + os.Getenv("wiseUser") + "\",\"password\":\"" + GetMD5Hash(os.Getenv("wisePws")) + "\",\"cdLingua\": \"it\"}")
		tokenReq, _ := http.Post("http://test-wopta.northeurope.cloudapp.azure.com/WebApiProduct/Api/loginWise", "", login)
		var result map[string]interface{}
		err := json.NewDecoder(tokenReq.Body).Decode(&result)
		lib.CheckError(err)

		log.Println(result)

		req, _ := http.NewRequest(http.MethodPost, urlstring, bytes.NewBuffer(jsonData))
		req.Header.Set("Autentication", string(result["token"].(string)))
		res, err := client.Do(req)
		lib.CheckError(err)
		if res != nil {
			body, err := ioutil.ReadAll(res.Body)
			lib.CheckError(err)
			res.Body.Close()

			fmt.Fprintf(w, string(body))
		}
	} else {
		fmt.Fprintf(w, "")
	}
	//lib.Files("")

}
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
