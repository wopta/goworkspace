package mail

import env "gitlab.dev.wopta.it/goworkspace/lib/environment"

func SendBaseEmail(body string, subject string, email string) {
	var mailRequest MailRequest
	mailRequest.IsHtml = true
	mailRequest.FromAddress = AddressAnna
	mailRequest.To = []string{email}
	if env.IsDevelopment() || env.IsLocal() {
		mailRequest.Subject = "dev: "
	}
	mailRequest.Subject += subject

	mailRequest.Message = body
	mailRequest.Message = mailRequest.Message + ` 
		<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px"><br></p><p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px">A presto,</p>
		<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#e50075;font-size:14px"><strong>Anna</strong> di Wopta Assicurazioni</p> `
	SendMail(mailRequest)
}
