package rules

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
)

func init() {
	log.Println("INIT Mail")

	functions.HTTP("Mail", Mail)
}

func Mail(w http.ResponseWriter, r *http.Request) {

	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	log.Println("mail")
	request := lib.ErrorByte(ioutil.ReadAll(r.Body))
	var send Request
	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(request), &send)
	base := "/mail"
	if strings.Contains(r.RequestURI, "/mail") {
		base = "/mail"
	} else {
		base = ""
	}
	log.Println(r.RequestURI)
	lib.EnableCors(&w, r)
	switch os := r.RequestURI; os {
	case base + "/send":
		Send(w, send)
	case base + "/pmi":

	default:
		fmt.Fprintf(w, " select path method send or read")
	}

	//lib.Files("")

}

type Request struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Message string   `json:"message"`
	Subject string   `json:"subject"`
}

func Send(resp http.ResponseWriter, obj Request) {
	for _, dest := range obj.To {
		password := "We20-tE22?"
		from := mail.Address{Name: "Wopta assicurazioni", Address: "website@wopta.it"}
		to := mail.Address{Name: dest, Address: dest}
		subj := obj.Subject
		body := obj.Message

		// Setup headers
		headers := make(map[string]string)
		headers["From"] = from.String()
		headers["To"] = to.String()
		headers["Subject"] = subj

		// Setup message
		message := ""
		for k, v := range headers {
			message += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		message += "\r\n" + body
		// Connect to the SMTP Server
		servername := "smtp.office365.com:587"
		host, _, _ := net.SplitHostPort(servername)
		auth := smtp.PlainAuth("", "username@example.tld", password, host)

		// TLS config
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}

		// Here is the key, you need to call tls.Dial instead of smtp.Dial
		// for smtp servers running on 465 that require an ssl connection
		// from the very beginning (no starttls)
		conn, err := tls.Dial("tcp", servername, tlsconfig)
		lib.CheckError(err)
		c, err := smtp.NewClient(conn, host)
		lib.CheckError(err)
		// Auth
		err = c.Auth(auth)
		lib.CheckError(err)
		// To && From
		err = c.Mail(from.Address)
		lib.CheckError(err)
		err = c.Rcpt(to.Address)
		lib.CheckError(err)
		// Data
		w, err := c.Data()
		lib.CheckError(err)
		_, err = w.Write([]byte(message))
		lib.CheckError(err)
		err = w.Close()
		lib.CheckError(err)
		c.Quit()
	}

}
