package wiseProxy

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
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
	log.Println(r.Method)
	lib.EnableCors(&w, r)
	jsonData, _ := ioutil.ReadAll(r.Body)
	client := http.Client{Timeout: time.Duration(5) * time.Second}
	//value := r.RequestURI
	log.Println("len(r.RequestURI): ", len(r.RequestURI))
	substring := r.RequestURI
	log.Println("substring: " + substring)
	var token string
	var urlstring = os.Getenv("wiseBaseUrl") + substring
	var req *http.Request
	log.Println("urlstring: " + urlstring)
	if strings.Contains(r.RequestURI, "/WebApiProduct") {
		token = GetToken(false)
		req, _ = http.NewRequest(r.Method, urlstring, bytes.NewBuffer(jsonData))
		//lib.CheckError(err)
		req.Header.Set("Authorization", "Bearer "+token)
	} else {
		token = GetToken(true)
		req, _ = http.NewRequest(r.Method, urlstring, bytes.NewBuffer(jsonData))
		//lib.CheckError(err)
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")
	log.Println("call: request")
	res, err := client.Do(req)
	lib.CheckError(err)
	defer res.Body.Close()
	if res != nil {
		//body, err := ioutil.ReadAll(res.Body)
		//log.Println("body: " + string(body))
		//lib.CheckError(err)
		io.Copy(w, res.Body)
		//fmt.Fprintf(w, strings.ReplaceAll(string(body), "(MISSING)", ""))
		res.Body.Close()
	}
	//lib.Files("")

}
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
func GetToken(isFewfine bool) string {
	var login *strings.Reader
	var url string
	if isFewfine {
		login = strings.NewReader("{\"username\": \"" + os.Getenv("wiseUser") + "\",\"password\":\"" + os.Getenv("wisePwd") + "\"}")
		url = os.Getenv("wiseBaseUrl")
	} else {
		login = strings.NewReader("{\"username\": \"" + os.Getenv("wiseUser") + "\",\"password\":\"" + GetMD5Hash(os.Getenv("wisePwd")) + "\",\"cdLingua\": \"it\"}")
		url = os.Getenv("wiseBaseUrl") + "WebApiProduct/Api/loginWise"
	}
	log.Println(login)
	log.Println("url: " + url)
	tokenReq, err := http.Post(url, "application/json", login)
	lib.CheckError(err)
	defer tokenReq.Body.Close()
	//var result *WiseLoginResponse
	result := &WiseLoginResponse{}
	log.Println("decode json login")
	body, err := ioutil.ReadAll(tokenReq.Body)
	err = json.Unmarshal(body, &result)
	lib.CheckError(err)
	log.Println(result)
	return result.DatiAuth.Token
}

type WiseLoginResponse struct {
	DatiAuth struct {
		ArrivoRichiesta time.Time `json:"arrivoRichiesta,omitempty"`
		ScadenzaToken   time.Time `json:"scadenzaToken,omitempty"`
		TipoToken       string    `json:"tipoToken,omitempty"`
		Token           string    `json:"token,omitempty"`
	} `json:"datiAuth,omitempty"`
}
