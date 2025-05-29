package mail

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
	"text/template"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"cloud.google.com/go/bigquery"
	"gitlab.dev.wopta.it/goworkspace/lib"
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

func SendFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var obj MailRequest

	log.AddPrefix("SendFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	req := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	json.Unmarshal(req, &obj)
	SendMail(obj)

	log.Println("Handler end -------------------------------------------------")

	return `{"message":"Success send "}`, nil, nil
}

func addAttachment(message, filename, contentType, data string) string {
	var ct string
	if contentType == "" {
		sct := strings.Split(filename, ".")
		ct = lib.GetContentType(sct[1])
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

func sendmail(obj MailRequest) error {
	log.AddPrefix("SendMail")
	defer log.PopPrefix()

	log.Println("start --------------------------------------------")
	var (
		username = os.Getenv("EMAIL_USERNAME")
		password = os.Getenv("EMAIL_PASSWORD")
		from     = AddressAnna
		file     []byte
		tpl      bytes.Buffer
		err      error
	)

	file, err = lib.GetFilesByEnvV2("mail/mail_template.html")
	if err != nil {
		return err
	}

	tmplt := template.New("action")
	tmplt, err = tmplt.Parse(string(file))
	if err != nil {
		return err
	}

	data := Data{
		Title:     obj.Title,
		SubTitle:  obj.SubTitle,
		IsLink:    obj.IsLink,
		Link:      obj.Link,
		LinkLabel: obj.LinkLabel,
		IsApp:     obj.IsApp,
		Content:   obj.Message,
	}
	err = tmplt.Execute(&tpl, data)
	if err != nil {
		return err
	}

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
		headers["Bcc"] = obj.Bcc

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
					message = addAttachment(message, v.FileName, v.ContentType, v.Byte)
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

		log.Println("sending message...")

		// Connect to the SMTP Server
		servername := "smtp.office365.com:587"
		host, _, err := net.SplitHostPort(servername)
		if err != nil {
			return err
		}

		// TLS config
		tlsconfig := &tls.Config{
			ServerName: host,
		}

		// Here is the key, you need to call tls.Dial instead of smtp.Dial
		// for smtp servers running on 465 that require an ssl connection
		// from the very beginning (no starttls)40.99.214.146
		conn, err := net.Dial("tcp", "smtp.office365.com:587")
		if err != nil {
			return err
		}

		c, err := smtp.NewClient(conn, host)
		if err != nil {
			return err
		}

		err = c.StartTLS(tlsconfig)
		if err != nil {
			return err
		}

		// Auth
		err = c.Auth(LoginAuth(username, password))
		if err != nil {
			return err
		}

		// To, From and Cc
		log.Printf("setting address from: %s", from.Address)
		err = c.Mail(from.Address)
		if err != nil {
			return err
		}

		log.Printf("setting address to: %s", to.Address)
		err = c.Rcpt(to.Address)
		if err != nil {
			return err
		}

		if obj.Cc != "" {
			// TODO: in the future we might need to handle multiple Ccs
			log.Printf("setting cc to: %s", obj.Cc)
			err = c.Rcpt(obj.Cc)
			if err != nil {
				return err
			}
		}

		if obj.Bcc != "" {
			// TODO: in the future we might need to handle multiple Bccs
			log.Printf("setting bcc to: %s", obj.Bcc)
			err = c.Rcpt(obj.Bcc)
			if err != nil {
				return err
			}
		}

		// Data
		w, err := c.Data()
		if err != nil {
			return err
		}

		_, err = w.Write([]byte(message))
		if err != nil {
			return err
		}

		err = w.Close()
		if err != nil {
			return err
		}

		err = c.Quit()
		if err != nil {
			return err
		}

		log.Println("message sent")
	}

	log.Println("end ----------------------------------------------")

	return nil
}

func SendMail(obj MailRequest) {
	var (
		reportError = ""
		reportRecip = strings.Join(obj.To, ",")
	)

	err := sendmail(obj)

	if err != nil {
		log.ErrorF("error sending mail: %s", err.Error())
		reportError = err.Error()
	}
	mailErr := writeMailReport(obj.Policy, obj.FromName, reportRecip, lib.GetBigQueryNullDateTime(time.Now().UTC()), reportError)
	if mailErr != nil {
		log.ErrorF("error writing report: %s", mailErr.Error())
	}
}

func writeMailReport(policyUid string, senderName string, recipientAddress string, date bigquery.NullDateTime, message string) error {

	report := MailReport{policyUid, senderName, recipientAddress, date, message}
	err := lib.InsertRowsBigQuery(lib.WoptaDataset, lib.MailReportCollection, report)
	if err != nil {
		return err
	}
	return nil
}
