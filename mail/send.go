package mail

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
	"text/template"

	"github.com/wopta/goworkspace/lib"
)

const (
	outerBoundary = "outer"
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
			return nil, errors.New("unknown fromServer")
		}
	}
	return nil, nil
}

func addAttachment(message, filename, contentType, data string) string {
	var ct string
	if contentType == "" {
		sct := strings.Split(filename, ".")
		ct = getContentType(sct[1])
	} else {
		ct = contentType
	}

	message += fmt.Sprintf("\r\n--%s\r\n", outerBoundary)
	message += fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", ct, filename)
	message += fmt.Sprintf("Content-Description: %s\r\n", filename)
	message += fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", filename)
	message += fmt.Sprintf("Content-Transfer-Encoding: base64\r\n")
	message += fmt.Sprintf("\r\n%s\r\n", string(data))

	return message
}

func getContentType(ext string) string {
	m := make(map[string]string)
	m["doc"] = "application/msword"
	m["docx"] = "application/msword"
	m["pdf"] = "application/pdf"
	m["GIF"] = "image/gif"
	m["jpeg"] = "image/jpeg"
	m["jpg"] = "image/jpeg"
	m["jpe"] = "image/jpeg"
	m["PNG"] = "image/png"
	m["png"] = "image/png"
	m["tiff"] = "image/tiff"
	m["tif"] = "image/tiff"
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
	log.Println("[SendMail] start --------------------------------------------")
	var (
		username = os.Getenv("EMAIL_USERNAME")
		password = os.Getenv("EMAIL_PASSWORD")
		from     = AddressAnna
		file     []byte
		tpl      bytes.Buffer
	)

	switch os.Getenv("env") {
	case "local":
		file = lib.ErrorByte(os.ReadFile("../function-data/dev/mail/mail_template.html"))
	case "dev":
		file = lib.GetFromStorage("function-data", "mail/mail_template.html", "")
	case "prod":
		file = lib.GetFromStorage("core-350507-function-data", "mail/mail_template.html", "")
	}

	tmplt := template.New("action")
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

	emptyAddress := mail.Address{}
	if obj.FromAddress.String() != emptyAddress.String() {
		from = obj.FromAddress
	} else if obj.From != "" {
		from.Address = obj.From
		if obj.FromName != "" {
			from.Name = obj.FromName
		} else {
			from.Name = obj.From
		}
	}

	for _, _to := range obj.To {
		to := mail.Address{Name: _to, Address: _to}
		subj := obj.Subject
		body := obj.Message

		// Setup headers
		headers := make(map[string]string)
		headers["From"] = from.String()
		headers["To"] = _to
		headers["Subject"] = subj
		headers["Cc"] = obj.Cc

		// Setup message
		message := ""
		for k, v := range headers {
			message += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		message += "MIME-Version: 1.0\r\n"

		if obj.IsAttachment {
			message += fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\n", outerBoundary)
			message += fmt.Sprintf("\r\n--%s\r\n", outerBoundary)
		}

		if obj.IsHtml {
			message += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
			message += fmt.Sprintf("\r\n%s", tpl.String())

			if obj.IsAttachment {
				for _, v := range *obj.Attachments {
					message = addAttachment(message, v.Name, v.ContentType, v.Byte)
				}
				message += fmt.Sprintf("\r\n--%s--\r\n", outerBoundary)
			}
		} else {
			message += "Content-Type:text/plain; charset=\"UTF-8\"\r\n"
			message += fmt.Sprintf("\r\n%s", body)

			if obj.IsAttachment {
				for _, v := range *obj.Attachments {
					message = addAttachment(message, v.Name, v.ContentType, v.Byte)
				}
				message += fmt.Sprintf("\r\n\n--%s--\r\n", outerBoundary)
			}
		}

		log.Println("[SendMail] sending message...")

		// Connect to the SMTP Server
		servername := "smtp.office365.com:587"
		host, _, err := net.SplitHostPort(servername)
		lib.CheckError(err)

		// TLS config
		tlsconfig := &tls.Config{
			ServerName: host,
		}

		// Here is the key, you need to call tls.Dial instead of smtp.Dial
		// for smtp servers running on 465 that require an ssl connection
		// from the very beginning (no starttls)40.99.214.146
		conn, err := net.Dial("tcp", "smtp.office365.com:587")
		lib.CheckError(err)

		c, err := smtp.NewClient(conn, host)
		lib.CheckError(err)

		c.StartTLS(tlsconfig)
		lib.CheckError(err)

		// Auth
		err = c.Auth(LoginAuth(username, password))
		lib.CheckError(err)

		// To, From and Cc
		log.Printf("[SendMail] setting address from: %s", from.Address)
		err = c.Mail(from.Address)
		lib.CheckError(err)

		log.Printf("[SendMail] setting address to: %s", to.Address)
		err = c.Rcpt(to.Address)
		lib.CheckError(err)

		if obj.Cc != "" {
			// TODO: in the future we might need to handle multiple Ccs
			log.Printf("[SendMail] setting cc to: %s", obj.Cc)
			err = c.Rcpt(obj.Cc)
			lib.CheckError(err)
		}

		// Data
		w, err := c.Data()
		lib.CheckError(err)

		_, err = w.Write([]byte(message))
		lib.CheckError(err)

		err = w.Close()
		lib.CheckError(err)

		c.Quit()

		log.Println("[SendMail] message sent")
	}

	log.Println("[SendMail] end ----------------------------------------------")
}
