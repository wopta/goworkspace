package consens

import (
	"bytes"
	"text/template"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
)

const templatePath = "mail/network-consens.html"

func sendConsensMail(networkNode *models.NetworkNode, consens SystemConsens, nodeConsens models.NodeConsens) error {
	var (
		templateBytes []byte
		tpl           = new(bytes.Buffer)
		htmlTemplate  = template.New("consens-mail")
		to            = mail.GetNetworkNodeEmail(networkNode)
		err           error
	)

	content := markdownParser(consens, nodeConsens)

	loc, _ := time.LoadLocation("Europe/Rome")

	bodyData := BodyData{
		AvailableTime: nodeConsens.GivenAt.In(loc).Format(time.TimeOnly),
		AvailableDate: nodeConsens.GivenAt.In(loc).Format("02/01/2006"),
		HtmlContent:   content,
	}

	if templateBytes, err = lib.GetFilesByEnvV2(templatePath); err != nil {
		return err
	}
	if htmlTemplate, err = htmlTemplate.Parse(string(templateBytes)); err != nil {
		return err
	}
	if err = htmlTemplate.Execute(tpl, bodyData); err != nil {
		return err
	}

	mailRequest := mail.MailRequest{
		FromAddress: mail.AddressAnna,
		To:          []string{to.Address},
		Message:     tpl.String(),
		Title:       nodeConsens.Title,
		SubTitle:    nodeConsens.Subtitle,
		Subject:     nodeConsens.Title,
		IsHtml:      true,
	}

	mail.SendMail(mailRequest)

	return nil
}

type BodyData struct {
	AvailableTime string
	AvailableDate string
	HtmlContent   string
}

const style = "-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px"

func markdownParser(consens SystemConsens, nodeConsens models.NodeConsens) string {
	fullText := ContentToString(consens.Content, nodeConsens.Answers, false)

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.HardLineBreak
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(fullText))

	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if p, ok := node.(*ast.Paragraph); ok && entering {
			attr := p.Attribute
			if attr == nil {
				attr = &ast.Attribute{
					Attrs: make(map[string][]byte),
				}
			}

			attr.Attrs["style"] = []byte(style)
			p.Attribute = attr
		}
		return ast.GoToNext
	})

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	resp := markdown.Render(doc, renderer)

	return string(resp)
}
