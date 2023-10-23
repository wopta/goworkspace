package mail

import "os"

func SendInviteMail(inviteUid, email string, isNetworkNode bool) {
	var mailRequest MailRequest

	mailRequest.FromAddress = AddressAnna
	mailRequest.To = []string{email}
	mailRequest.Subject = "Benvenuto in Wopta!"
	mailRequest.IsHtml = true

	// move to template in function data
	lines := []string{
		"Ciao,",
		"Ecco il tuo invito al tuo account wopta.it.",
		"Accedi al link sottostante e crea la tua password.",
	}
	for _, line := range lines {
		mailRequest.Message = mailRequest.Message + `<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px">` + line + `</p>`
	}

	mailRequest.Message = mailRequest.Message + ` 
	<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px"><br></p><p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px">A presto,</p>
	<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#e50075;font-size:14px"><strong>Anna</strong> di Wopta Assicurazioni</p> `
	mailRequest.Title = "Invito a wopta.it"
	mailRequest.IsLink = true
	mailRequest.Link = os.Getenv("WOPTA_CUSTOMER_AREA_BASE_URL") + "/login/inviteregistration?inviteUid=" + inviteUid
	if isNetworkNode {
		mailRequest.Link += "&isNetworkNode=true"
	}
	mailRequest.LinkLabel = "Crea la tua password"

	SendMail(mailRequest)
}
