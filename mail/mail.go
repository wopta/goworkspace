package rules

/*
mail/send POST
REQUEST EXAMPLE:
{
    "from": "website@wopta.it",
    "to": ["luca.barbieri@wopta.it"],
    "message":"messaggio lungo lungo lungo lungo che arriva dalla nostro backend dedidacao a tutti quelli che leggono e c'è anche il cc se ci fosse bisogno, qui si spra sia andato a capo da aggiungere gli allegati in TO DO or ora ma non subito   ",
    "subject": "Wopta per te Firma",
  "isHtml":true,
  "cc":"luca.barbieri.81@gmail.com",
  "isAttachment":true,
  "attachments":[{"name":"text.txt",
                 "byte":"Q2lhbyBEYXZpZGUsDQppIGZhdHRpIGUgaWwgY29kaWNlIHByb2RvdHRvIHBhcmxhbm8gZGEgc29saSwgbGUgcGVyaXppZSBhbmNoZSBzZSB0dSBpbnN1bHRpIHNlbnphIG5lYW5jaGUgc2FwZXJlIGNoaSwgc29ubyBwcm9mZXNzaW9uaXN0aSBjaGUgbGF2b3Jhbm8gcGVyIGdyb3NzZSBhemllbmRlIGUgbmVzc3VubyBkaSBlc3NpIGxhdm9yYSBwZXIgbmV4dW0gbWEgc29ubyBlc3Rlcm5pIGNoaWFtYXRpIHNvbG8gcGVyIHBlcml6aWUgcHJldmVudGl2ZSBzZSB2dW9pIHRlIGxpIHByZXNlbnRvIG1hZ2FyaSB0aSBhaXV0YW5vIGEgbWlnbGlvcmFyIGlsIHR1byBzdGFmZi4NClZvaSBub24gdmkgc2lldGUgcHJlc2VudGF0aSBhIGNvbiBuZXNzdW5hIHNvbHV6aW9uZSBzZSBub24gY29uIHVuYSBjaXRhemlvbmUgZGkgdW4gQXBpa2V5LCBtYSBzZW56YSBhcmdvbWVudGF6aW9uZSBlIHNlbnphIHVuIGRvY3VtZW50byBhIHVsdGVyaW9yZSBwcm92YSBkZWxsYSB2b3N0cm8gYXBwcm9jY2lvIHNiYWdsaWF0byBjaGUgdmkgaGEgcG9ydGF0byBpbiBxdWVzdGEgc2l0dWF6aW9uZSwgbmV4dW0gIHZpIGhhIGFyZ29tZW50YXRvIGUgaGEgcG9ydGF0byB1biBkb2N1bWVudG8gYSBzdXBwb3J0byBjb24gYWRkaXJpdHR1cmEgaWwgYm9pbGVycGxhdGUgZGVsIGNvZGljZSBwZXIgbOKAmWludGVncmF6aW9uZSBlIGlsIGdpb3JubyBkb3BvIGVyYSB0dXR0byBwcm9udG8gZSBjb25maWd1cmF0byBlIHF1ZXN0byBub24gbG8gcHVvaSBuZWdhcmUgYW5jaGUgcXVpIMOoIHR1dHRvIGZhY2lsbWVudGUgZGltb3N0cmFiaWxlLg0KVW4gY29udHJhdHRvIHNpIGZhIHBlciBldmFkZXJlIGRlbGxlIGF0dGl2aXTDoCwgbG8gc3RhZmZpbmcgIMOoIHN0YXJvIGNvbmNvcmRhdG8gaW5zaWVtZSBhIHRlIHBlciBxdWVzdG8gdGFzayBkaSBwcmV2ZW50aXZhemlvbmUgZSBhIGx1Z2xpbyBhYmJpYW1vIHBhZ2F0byBpbCBkb3BwaW8gcGVyIG5lc3N1bmEgYXR0aXZpdMOgIGluIG1lcml0byBtYSBsaSBub24gZXJhdmFtbyBwcm9udGkgbm9pIGVkIGluZmF0dGkgbm9uIGNpdG8gbWFpIGx1Z2xpbyBxdWluZGkgY29zaSB0aSBzdGFpIGRhbmRvIGxhIHphcHBhIG5laSBwaWVkaSBkYSBzb2xvIHBlcmNow6kgIMOoIGZhY2lsbWVudGUgZGltb3N0cmFiaWxlIGNoZSBkYSBhZ29zdG8gbm9uIGF2ZXRlIG5lYW5jaGUgaW5pemlhdG8uDQpEaW1taSBpbiBjaGUgZm9ybWF0byB2dW9pIGkgbG9nIG5lIGhvIGRpIHR1dHRlIGxlIG5hdHVyZSBzZSBub24gdGkgYmFzdGFubyBsZSAxMDAgZGljbyAxMDAgY2hpYW1hdGUgIGNoZSB0aSBobyBnacOgIGdpcmF0bw0KSW4gYWxsZWdhdG8gdHJvdmkgdW5hIG1haWwgY29uIHVuIHBheWxvYWQgY2hlIGdpw6AgZnVuemlvbmF2YSAoaG8gaSBsb2cgZGVsbGEgY2hpYW1hdGEgcGVyY2jDqSBs4oCZYWJiaWFtbyBmYXR0YSBpbnNpZW1lIGFkIEVucmljbyBhbmNoZSBzZSBobyBkb3Z1dG8gbWV0dGVybGEgcHViYmxpY2EgcGVyY2jDqSBub24gc2lldGUgcml1c2NpdGkgYSBzdmlsdXBwYXJlIGwgYXV0ZW50aWNhemlvbmUgdGUgbG8gcHXDsiBjb25mZXJtYXJlIGx1aSBkaXJldHRhbWVudGUpDQpJbiBxdWVsbGEgbWFpbCBkZWwgMjcvMDggIHJpY2hpZWRvIHVuYSBkZW1vIGlsIGRvcG8gYmVuIHVuIG1lc2UgY2hlICBjaSBzdGF2YXRlIGxhdm9yYW5kbyAoc2VtcHJlIHNlIGNpIGF2ZXRlIG1haSBsYXZvcmF0bykgaW4gY3VpIG5vbiBhdmV0ZSBtYWkgcmlzcG9zdG8gc2UgY2VyY2hpIGJlbmUgcG9pIHRyb3ZlcmFpIGFuY2hlIGlsIHNvbGxlY2l0byBhZCB1bHRlcmlvcmUgcmlwcm92YSBkaSBxdWFudG8gc2VpIHNmYXNhdG8gY29uIGxhIHJlYWx0w6AuDQpJbm9sdHJlIGluIGFsbGVnYXRvIGEgcXVlbGxhIG1haWwgdmkgY29uc2Vnbm8gaWwgbW9kZWwgZ2nDoCBzdmlsdXBwYXRvIHBlciBxdWVsbGEgY2hpYW1hdGEgYSByaXByb3ZhIGNoZSBhdmV2YXRlIHR1dHRlIGxlIGluZm8gcGVyIHN2aWx1cHBhcmxhLg0KUmliYWRpc2NvIHZpIHNpZXRlIGluY2FydGF0aSBvIGzigJlhdmV0ZSBwcmVzYSBzb3R0byBnYW1iYSDDqCBwYWxlc2UgYSB0dXR0aSBzb2xvIHR1IG5lZ2hpIGzigJlldmlkZW56YSBkZWkgZmF0dGkuDQpncmF6aWUNCkJ1b25hIHNlcmF0YQ0KDQoNCg=="
},{"name":"text.txt",
                 "byte":"Q2lhbyBEYXZpZGUsDQppIGZhdHRpIGUgaWwgY29kaWNlIHByb2RvdHRvIHBhcmxhbm8gZGEgc29saSwgbGUgcGVyaXppZSBhbmNoZSBzZSB0dSBpbnN1bHRpIHNlbnphIG5lYW5jaGUgc2FwZXJlIGNoaSwgc29ubyBwcm9mZXNzaW9uaXN0aSBjaGUgbGF2b3Jhbm8gcGVyIGdyb3NzZSBhemllbmRlIGUgbmVzc3VubyBkaSBlc3NpIGxhdm9yYSBwZXIgbmV4dW0gbWEgc29ubyBlc3Rlcm5pIGNoaWFtYXRpIHNvbG8gcGVyIHBlcml6aWUgcHJldmVudGl2ZSBzZSB2dW9pIHRlIGxpIHByZXNlbnRvIG1hZ2FyaSB0aSBhaXV0YW5vIGEgbWlnbGlvcmFyIGlsIHR1byBzdGFmZi4NClZvaSBub24gdmkgc2lldGUgcHJlc2VudGF0aSBhIGNvbiBuZXNzdW5hIHNvbHV6aW9uZSBzZSBub24gY29uIHVuYSBjaXRhemlvbmUgZGkgdW4gQXBpa2V5LCBtYSBzZW56YSBhcmdvbWVudGF6aW9uZSBlIHNlbnphIHVuIGRvY3VtZW50byBhIHVsdGVyaW9yZSBwcm92YSBkZWxsYSB2b3N0cm8gYXBwcm9jY2lvIHNiYWdsaWF0byBjaGUgdmkgaGEgcG9ydGF0byBpbiBxdWVzdGEgc2l0dWF6aW9uZSwgbmV4dW0gIHZpIGhhIGFyZ29tZW50YXRvIGUgaGEgcG9ydGF0byB1biBkb2N1bWVudG8gYSBzdXBwb3J0byBjb24gYWRkaXJpdHR1cmEgaWwgYm9pbGVycGxhdGUgZGVsIGNvZGljZSBwZXIgbOKAmWludGVncmF6aW9uZSBlIGlsIGdpb3JubyBkb3BvIGVyYSB0dXR0byBwcm9udG8gZSBjb25maWd1cmF0byBlIHF1ZXN0byBub24gbG8gcHVvaSBuZWdhcmUgYW5jaGUgcXVpIMOoIHR1dHRvIGZhY2lsbWVudGUgZGltb3N0cmFiaWxlLg0KVW4gY29udHJhdHRvIHNpIGZhIHBlciBldmFkZXJlIGRlbGxlIGF0dGl2aXTDoCwgbG8gc3RhZmZpbmcgIMOoIHN0YXJvIGNvbmNvcmRhdG8gaW5zaWVtZSBhIHRlIHBlciBxdWVzdG8gdGFzayBkaSBwcmV2ZW50aXZhemlvbmUgZSBhIGx1Z2xpbyBhYmJpYW1vIHBhZ2F0byBpbCBkb3BwaW8gcGVyIG5lc3N1bmEgYXR0aXZpdMOgIGluIG1lcml0byBtYSBsaSBub24gZXJhdmFtbyBwcm9udGkgbm9pIGVkIGluZmF0dGkgbm9uIGNpdG8gbWFpIGx1Z2xpbyBxdWluZGkgY29zaSB0aSBzdGFpIGRhbmRvIGxhIHphcHBhIG5laSBwaWVkaSBkYSBzb2xvIHBlcmNow6kgIMOoIGZhY2lsbWVudGUgZGltb3N0cmFiaWxlIGNoZSBkYSBhZ29zdG8gbm9uIGF2ZXRlIG5lYW5jaGUgaW5pemlhdG8uDQpEaW1taSBpbiBjaGUgZm9ybWF0byB2dW9pIGkgbG9nIG5lIGhvIGRpIHR1dHRlIGxlIG5hdHVyZSBzZSBub24gdGkgYmFzdGFubyBsZSAxMDAgZGljbyAxMDAgY2hpYW1hdGUgIGNoZSB0aSBobyBnacOgIGdpcmF0bw0KSW4gYWxsZWdhdG8gdHJvdmkgdW5hIG1haWwgY29uIHVuIHBheWxvYWQgY2hlIGdpw6AgZnVuemlvbmF2YSAoaG8gaSBsb2cgZGVsbGEgY2hpYW1hdGEgcGVyY2jDqSBs4oCZYWJiaWFtbyBmYXR0YSBpbnNpZW1lIGFkIEVucmljbyBhbmNoZSBzZSBobyBkb3Z1dG8gbWV0dGVybGEgcHViYmxpY2EgcGVyY2jDqSBub24gc2lldGUgcml1c2NpdGkgYSBzdmlsdXBwYXJlIGwgYXV0ZW50aWNhemlvbmUgdGUgbG8gcHXDsiBjb25mZXJtYXJlIGx1aSBkaXJldHRhbWVudGUpDQpJbiBxdWVsbGEgbWFpbCBkZWwgMjcvMDggIHJpY2hpZWRvIHVuYSBkZW1vIGlsIGRvcG8gYmVuIHVuIG1lc2UgY2hlICBjaSBzdGF2YXRlIGxhdm9yYW5kbyAoc2VtcHJlIHNlIGNpIGF2ZXRlIG1haSBsYXZvcmF0bykgaW4gY3VpIG5vbiBhdmV0ZSBtYWkgcmlzcG9zdG8gc2UgY2VyY2hpIGJlbmUgcG9pIHRyb3ZlcmFpIGFuY2hlIGlsIHNvbGxlY2l0byBhZCB1bHRlcmlvcmUgcmlwcm92YSBkaSBxdWFudG8gc2VpIHNmYXNhdG8gY29uIGxhIHJlYWx0w6AuDQpJbm9sdHJlIGluIGFsbGVnYXRvIGEgcXVlbGxhIG1haWwgdmkgY29uc2Vnbm8gaWwgbW9kZWwgZ2nDoCBzdmlsdXBwYXRvIHBlciBxdWVsbGEgY2hpYW1hdGEgYSByaXByb3ZhIGNoZSBhdmV2YXRlIHR1dHRlIGxlIGluZm8gcGVyIHN2aWx1cHBhcmxhLg0KUmliYWRpc2NvIHZpIHNpZXRlIGluY2FydGF0aSBvIGzigJlhdmV0ZSBwcmVzYSBzb3R0byBnYW1iYSDDqCBwYWxlc2UgYSB0dXR0aSBzb2xvIHR1IG5lZ2hpIGzigJlldmlkZW56YSBkZWkgZmF0dGkuDQpncmF6aWUNCkJ1b25hIHNlcmF0YQ0KDQoNCg=="
}]

}
*/
import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
	"text/template"

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
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	log.Println(string(req))
	var send Request
	// Unmarshal or Decode the JSON to the interface.
	//json.NewDecoder(req).Decode(&send)
	defer r.Body.Close()

	json.Unmarshal([]byte(req), &send)
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
		Send(w, &send)
	case base + "/read":

	default:
		fmt.Fprintf(w, " select path method send or read")
	}

	//lib.Files("")

}

type Data struct {
	Title   string
	Content string
}
type Attachment struct {
	Name        string `json:"name"`
	Byte        string `json:"byte"`
	ContentType string `json:"contentType,omitempty"`
}
type Request struct {
	From         string       `json:"from"`
	To           []string     `json:"to"`
	Message      string       `json:"message"`
	Subject      string       `json:"subject"`
	IsHtml       bool         `json:"isHtml,omitempty"`
	IsAttachment bool         `json:"isAttachment,omitempty"`
	Attachments  []Attachment `json:"attachments,omitempty"`
	Cc           string       `json:"cc,omitempty"`
	TemplateName string       `json:"templateName,omitempty"`
}

func Send(resp http.ResponseWriter, obj *Request) {
	var (
		host       = os.Getenv("EMAIL_HOST")
		username   = os.Getenv("EMAIL_USERNAME")
		password   = os.Getenv("EMAIL_PASSWORD")
		portNumber = os.Getenv("EMAIL_PORT")
		file       []byte
	)

	const (
		boundary = "my-boundary-779"
	)
	log.Println(username)
	log.Println(password)
	log.Println(host)
	log.Println(portNumber)
	switch os.Getenv("env") {
	case "local":
		file = lib.ErrorByte(ioutil.ReadFile("function-data/mail/mail_template.html"))

	case "dev":
		file = lib.GetFromStorage("function-data", "mail/mail_template.html", "")

	case "prod":
		file = lib.GetFromStorage("core-350507-function-data", "mail/mail_template.html", "")

	default:

	}
	tmplt := template.New("action")
	var tpl bytes.Buffer

	tmplt, err := tmplt.Parse(string(file))
	lib.CheckError(err)
	data := Data{Title: obj.Subject, Content: obj.Message}
	tmplt.Execute(&tpl, data)
	for _, _to := range obj.To {
		//password := "We20-tE22?"
		from := mail.Address{Name: "Wopta assicurazioni", Address: "website@wopta.it"}
		to := mail.Address{Name: _to, Address: _to}
		subj := obj.Subject
		body := obj.Message
		// Setup headers
		headers := make(map[string]string)
		headers["From"] = from.String()
		headers["To"] = _to
		headers["Subject"] = subj
		if len(obj.Cc) > 2 {
			headers["Cc"] = obj.Cc
		}

		// Setup message
		message := ""
		for k, v := range headers {
			message += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		message += fmt.Sprintf("MIME-Version: 1.0\r\n")
		if obj.IsAttachment {
			message += fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\n", boundary)
			message += fmt.Sprintf("\r\n--%s\r\n", boundary)
		}
		if obj.IsHtml {
			message += "Content-Type: text/html; charset=\"UTF-8\";\r\n"
			message += "\r\n" + tpl.String()
			if obj.IsAttachment {
				message += fmt.Sprintf("\r\n--%s\r\n", boundary)
			}

		} else {
			message += "Content-Type:text/plain; charset=\"UTF-8\";\r\n"
			message += "\r\n" + body
			if obj.IsAttachment {
				message += fmt.Sprintf("\r\n--%s\r\n", boundary)
			}
		}

		if obj.IsAttachment {
			for k, v := range obj.Attachments {

				var close bool
				if k == len(obj.Attachments)-1 {
					close = true
				}
				message = addAttachment(message, v.Name, v.ContentType, []byte(v.Byte), close)
			}

		}
		//message += "\r\n" + body
		log.Println("MESSAGE:----------------------")
		log.Println(message)
		// Connect to the SMTP Server
		servername := "smtp.office365.com:587"
		host, _, err := net.SplitHostPort(servername)
		lib.CheckError(err)
		//auth := smtp.PlainAuth("", "website@wopta.it", "We20-tE22?", host)

		// TLS config
		tlsconfig := &tls.Config{
			//InsecureSkipVerify: true,
			ServerName: host,
		}

		// Here is the key, you need to call tls.Dial instead of smtp.Dial
		// for smtp servers running on 465 that require an ssl connection
		// from the very beginning (no starttls)40.99.214.146
		log.Println("end MESSAGE:----------------------")
		conn, err := net.Dial("tcp", "smtp.office365.com:587")
		log.Println("end DIAL:----------------------")
		lib.CheckError(err)
		c, err := smtp.NewClient(conn, host)
		lib.CheckError(err)
		c.StartTLS(tlsconfig)
		lib.CheckError(err)
		log.Println("end Tls:----------------------")
		// Auth
		err = c.Auth(LoginAuth(username, password))
		lib.CheckError(err)
		// To && From
		log.Println("start mail:----------------------")
		err = c.Mail(from.Address)
		log.Println("end Mail:----------------------")
		lib.CheckError(err)
		err = c.Rcpt(to.Address)
		log.Println("end Rcpt:----------------------")
		lib.CheckError(err)
		// Data
		w, err := c.Data()

		lib.CheckError(err)
		log.Println("start write massage:----------------------")
		_, err = w.Write([]byte(message))
		log.Println("end write massage:----------------------")
		lib.CheckError(err)
		err = w.Close()
		lib.CheckError(err)
		c.Quit()
		fmt.Fprintf(resp, " mail send")
	}

}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown fromServer")
		}
	}
	return nil, nil
}
func addAttachment(message string, name string, contentType string, data []byte, close bool) string {

	const (
		boundary = "my-boundary-779"
	)
	var ct string
	if contentType == "" {
		sct := strings.Split(name, ".")
		ct = getContentType(sct[1])
	} else {
		ct = contentType
	}

	//b := base64.URLEncoding.EncodeToString(data)

	//message += fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary)
	message += fmt.Sprintf("\r\n")
	message += fmt.Sprintf("Content-Type: " + ct + "; charset=\"utf-8\"\r\n")
	message += fmt.Sprintf("Content-Transfer-Encoding: base64\r\n")
	message += fmt.Sprintf("Content-Disposition: attachment; filename=" + name + "\r\n")
	message += fmt.Sprintf("Content-ID: <" + name + ">\r\n")
	message += fmt.Sprintf(string(data))
	message += fmt.Sprintf("\r\n--%s", boundary)
	if close {
		message += fmt.Sprintf("--")
	}

	return message
}

func getContentType(ext string) string {
	m := make(map[string]string)
	m["doc"] = "applicazione/msword"
	m["docx"] = "applicazione/msword"
	m["pdf"] = "applicazione/pdf"
	m["GIF"] = "immagine/gif"
	m["jpeg"] = "immagine/jpeg"
	m["jpg"] = "immagine/jpeg"
	m["jpe"] = "immagine/jpeg"
	m["PNG"] = "immagine/png"
	m["png"] = "immagine/png"
	m["tiff"] = "immagine/tiff"
	m["tif"] = "immagine/tiff"
	m["xls"] = "application/vnd.ms-excel"
	m["xlsx"] = "application/vnd.ms-excel"
	m["pptx"] = "application/vnd.ms-powerpoint"
	m["ppt"] = "application/vnd.ms-powerpoint"
	m["txt"] = "text/plain"
	m["zip"] = "applicazione/zip"
	m["gzip"] = "applicazione/x-gzip"
	return m[ext]
}
