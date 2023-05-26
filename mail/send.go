package mail

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
	"text/template"

	lib "github.com/wopta/goworkspace/lib"
)

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
func addAttachment(message string, name string, contentType string, data string, close bool) string {

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

	//b := base64.URLEncoding.EncodeToString(data) iso-8859-1

	//message += fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary)
	message += fmt.Sprintf("\r\n")
	message += fmt.Sprintf("Content-Type: " + ct + ";\r\n")
	message += fmt.Sprintf("Content-Transfer-Encoding: base64\r\n")
	message += fmt.Sprintf("Content-Disposition: attachment; filename=\"" + name + "\"\r\n")

	message += fmt.Sprintf("\r\n" + string(data) + "\r\n")
	message += fmt.Sprintf("\r\n--%s", boundary)
	if close {
		message += fmt.Sprintf("--")
	}

	return message
}

func getContentType(ext string) string {
	m := make(map[string]string)
	m["doc"] = "application/msword"
	m["docx"] = "application/msword"
	m["pdf"] = "application/pdf"
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
	m["zip"] = "application/zip"
	m["gzip"] = "application/x-gzip"
	return m[ext]
}
func SendMail(obj MailRequest) {
	var (
		username = os.Getenv("EMAIL_USERNAME")
		password = os.Getenv("EMAIL_PASSWORD")

		file []byte
	)

	const (
		boundary = "my-boundary-779"
	)

	switch os.Getenv("env") {
	case "local":
		file = lib.ErrorByte(ioutil.ReadFile("../function-data/dev/mail/mail_template.html"))

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
	data := Data{
		Title:     obj.Title,
		SubTitle:  obj.SubTitle,
		IsLink:    obj.IsLink,
		Link:      obj.Link,
		LinkLabel: obj.LinkLabel,
		IsApp:     obj.IsApp,
		Content:   obj.Message,
	}
	tmplt.Execute(&tpl, data)
	log.Println()
	for _, _to := range obj.To {

		from := mail.Address{Name: "Anna di Wopta Assicurazioni", Address: obj.From}
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
			for k, v := range *obj.Attachments {

				var close bool
				if k == len(*obj.Attachments)-1 {
					close = true
				}
				message = addAttachment(message, v.Name, v.ContentType, v.Byte, close)
			}

		}
		//message += "\r\n" + body
		log.Println("MESSAGE:----------------------")
		//log.Println(message)
		// Connect to the SMTP Server
		servername := "smtp.office365.com:587"
		host, _, err := net.SplitHostPort(servername)
		lib.CheckError(err)

		// TLS config
		tlsconfig := &tls.Config{
			//InsecureSkipVerify: true,
			ServerName: host,
		}

		// Here is the key, you need to call tls.Dial instead of smtp.Dial
		// for smtp servers running on 465 that require an ssl connection
		// from the very beginning (no starttls)40.99.214.146
		//log.Println("end MESSAGE:----------------------")
		conn, err := net.Dial("tcp", "smtp.office365.com:587")
		//log.Println("end DIAL:----------------------")
		lib.CheckError(err)
		c, err := smtp.NewClient(conn, host)
		lib.CheckError(err)
		c.StartTLS(tlsconfig)
		lib.CheckError(err)
		//log.Println("end Tls:----------------------")
		// Auth
		err = c.Auth(LoginAuth(username, password))
		lib.CheckError(err)
		// To && From
		//log.Println("start mail:----------------------")
		err = c.Mail(from.Address)
		//log.Println("end Mail:----------------------")
		lib.CheckError(err)
		err = c.Rcpt(to.Address)
		//log.Println("end Rcpt:----------------------")
		lib.CheckError(err)
		// Data
		w, err := c.Data()
		lib.CheckError(err)
		//log.Println("start write massage:----------------------")
		_, err = w.Write([]byte(message))
		//log.Println(message)
		log.Println("end write massage:----------------------")
		//log.Println(message)
		lib.CheckError(err)
		err = w.Close()
		lib.CheckError(err)
		c.Quit()

	}
}
