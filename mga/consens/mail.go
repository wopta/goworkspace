package consens

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
)

const templatePath = "mail/network-consens.html"

func sendConsensMail(networkNode *models.NetworkNode, consens models.NodeConsens) error {
	var (
		templateBytes []byte
		tpl           = new(bytes.Buffer)
		htmlTemplate  = template.New("consens-mail")
		to            = mail.GetNetworkNodeEmail(networkNode)
		bodyData      = getBodyData(networkNode, consens)
		err           error
	)

	if templateBytes, err = lib.GetFilesByEnvV2(templatePath); err != nil {
		return err
	}
	if htmlTemplate, err = htmlTemplate.Parse(string(templateBytes)); err != nil {
		return err
	}
	if err = htmlTemplate.Execute(tpl, bodyData); err != nil {
		return err
	}

	title := fmt.Sprintf("Consenso: %s", consens.Title)

	mailRequest := mail.MailRequest{
		FromAddress: mail.AddressAnna,
		To:          []string{to.Address},
		Message:     tpl.String(),
		Title:       title,
		Subject:     title,
		IsHtml:      true,
	}

	mail.SendMail(mailRequest)

	return nil
}

type BodyData struct {
	Name         string
	Date         string
	Time         string
	ConsensTitle string
	ConsensValue string
}

func getBodyData(networkNode *models.NetworkNode, consens models.NodeConsens) BodyData {
	value := strings.ReplaceAll(consens.Value, "_", " ")
	value = strings.ReplaceAll(value, "-", " ")
	return BodyData{
		Name:         networkNode.GetName(),
		Date:         consens.GivenAt.Format("02/01/2006"),
		Time:         consens.GivenAt.Format(time.TimeOnly),
		ConsensTitle: consens.Title,
		ConsensValue: lib.ToUpper(value),
	}
}
