package callback

import (
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
)

func EmailVerify(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("EmailVerify")
	log.Println("GET params were:", r.URL.Query())
	email := r.URL.Query().Get("email")
	token := r.URL.Query().Get("token")
	log.Println(token)
	res := lib.WhereFirestore("mail", "mail", "==", email)

	objmail, uid := mail.ToListData(res)
	log.Println(objmail)
	if len(objmail) > 0 {
		objmail[0].IsValid = true
		lib.SetFirestore("mail", uid[0], objmail[0])
	}

	return getResponse("<p>Grazie la tua mail è stata validata poi continuare l'acquisto</p>", "Validazione Mail", email), nil, nil
}

func getResponse(content string, title string, sub string) string {
	return `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
	<html xmlns="http://www.w3.org/1999/xhtml" xmlns:o="urn:schemas-microsoft-com:office:office" style="font-family:arial, 'helvetica neue', helvetica, sans-serif">
	 <head>
	  <meta charset="UTF-8">
	  <meta content="width=device-width, initial-scale=1" name="viewport">
	  <meta name="x-apple-disable-message-reformatting">
	  <meta http-equiv="X-UA-Compatible" content="IE=edge">
	  <meta content="telephone=no" name="format-detection">
	  <title>Wopta email validazione</title><!--[if (mso 16)]>
		<style type="text/css">
		a {text-decoration: none;}
		</style>
		<![endif]--><!--[if gte mso 9]><style>sup { font-size: 100% !important; }</style><![endif]--><!--[if gte mso 9]>
	<xml>
		<o:OfficeDocumentSettings>
		<o:AllowPNG></o:AllowPNG>
		<o:PixelsPerInch>96</o:PixelsPerInch>
		</o:OfficeDocumentSettings>
	</xml>
	<![endif]--><!--[if !mso]><!-- -->
	  <link href="https://fonts.googleapis.com/css2?family=Imprima&display=swap" rel="stylesheet"><!--<![endif]-->
	  <style type="text/css">
	#outlook a {
		padding:0;
	}
	.es-button {
		mso-style-priority:100!important;
		text-decoration:none!important;
	}
	a[x-apple-data-detectors] {
		color:inherit!important;
		text-decoration:none!important;
		font-size:inherit!important;
		font-family:inherit!important;
		font-weight:inherit!important;
		line-height:inherit!important;
	}
	.es-desk-hidden {
		display:none;
		float:left;
		overflow:hidden;
		width:0;
		max-height:0;
		line-height:0;
		mso-hide:all;
	}
	.es-button-border:hover a.es-button, .es-button-border:hover button.es-button {
		background:#ff1990!important;
	}
	.es-button-border:hover {
		border-color:#42d159 #42d159 #42d159 #42d159!important;
		background:#ff1990!important;
	}
	[data-ogsb] .es-button.es-button-1 {
		padding:10px 5px!important;
	}
	td .es-button-border:hover a.es-button-2 {
		background:#e50075!important;
	}
	@media only screen and (max-width:600px) {p, ul li, ol li, a { line-height:120%!important } h1, h2, h3, h1 a, h2 a, h3 a { line-height:120% } h1 { font-size:14px!important; text-align:left } h2 { font-size:14px!important; text-align:left } h3 { font-size:14px!important; text-align:left } .es-header-body h1 a, .es-content-body h1 a, .es-footer-body h1 a { font-size:14px!important; text-align:left } .es-header-body h2 a, .es-content-body h2 a, .es-footer-body h2 a { font-size:14px!important; text-align:left } .es-header-body h3 a, .es-content-body h3 a, .es-footer-body h3 a { font-size:14px!important; text-align:left } .es-menu td a { font-size:12px!important } .es-header-body p, .es-header-body ul li, .es-header-body ol li, .es-header-body a { font-size:12px!important } .es-content-body p, .es-content-body ul li, .es-content-body ol li, .es-content-body a { font-size:12px!important } .es-footer-body p, .es-footer-body ul li, .es-footer-body ol li, .es-footer-body a { font-size:12px!important } .es-infoblock p, .es-infoblock ul li, .es-infoblock ol li, .es-infoblock a { font-size:12px!important } *[class="gmail-fix"] { display:none!important } .es-m-txt-c, .es-m-txt-c h1, .es-m-txt-c h2, .es-m-txt-c h3 { text-align:center!important } .es-m-txt-r, .es-m-txt-r h1, .es-m-txt-r h2, .es-m-txt-r h3 { text-align:right!important } .es-m-txt-l, .es-m-txt-l h1, .es-m-txt-l h2, .es-m-txt-l h3 { text-align:left!important } .es-m-txt-r img, .es-m-txt-c img, .es-m-txt-l img { display:inline!important } .es-button-border { display:block!important } a.es-button, button.es-button { font-size:14px!important; display:block!important; border-right-width:0px!important; border-left-width:0px!important; border-top-width:15px!important; border-bottom-width:15px!important; padding-left:0px!important; padding-right:0px!important } .es-adaptive table, .es-left, .es-right { width:100%!important } .es-content table, .es-header table, .es-footer table, .es-content, .es-footer, .es-header { width:100%!important; max-width:600px!important } .es-adapt-td { display:block!important; width:100%!important } .adapt-img { width:100%!important; height:auto!important } .es-m-p0 { padding:0px!important } .es-m-p0r { padding-right:0px!important } .es-m-p0l { padding-left:0px!important } .es-m-p0t { padding-top:0px!important } .es-m-p0b { padding-bottom:0!important } .es-m-p20b { padding-bottom:20px!important } .es-mobile-hidden, .es-hidden { display:none!important } tr.es-desk-hidden, td.es-desk-hidden, table.es-desk-hidden { width:auto!important; overflow:visible!important; float:none!important; max-height:inherit!important; line-height:inherit!important } tr.es-desk-hidden { display:table-row!important } table.es-desk-hidden { display:table!important } td.es-desk-menu-hidden { display:table-cell!important } .es-menu td { width:1%!important } table.es-table-not-adapt, .esd-block-html table { width:auto!important } table.es-social { display:inline-block!important } table.es-social td { display:inline-block!important } .es-desk-hidden { display:table-row!important; width:auto!important; overflow:visible!important; max-height:inherit!important } }
	</style>
	 </head>
	 <body style="width:100%;font-family:arial, 'helvetica neue', helvetica, sans-serif;-webkit-text-size-adjust:100%;-ms-text-size-adjust:100%;padding:0;Margin:0">
	  <div class="es-wrapper-color" style="background-color:#FFFFFF"><!--[if gte mso 9]>
				<v:background xmlns:v="urn:schemas-microsoft-com:vml" fill="t">
					<v:fill type="tile" color="#ffffff"></v:fill>
				</v:background>
			<![endif]-->
	   <table class="es-wrapper" width="100%" cellspacing="0" cellpadding="0" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;padding:0;Margin:0;width:100%;height:100%;background-repeat:repeat;background-position:center top;background-color:#FFFFFF">
		 <tr>
		  <td valign="top" style="padding:0;Margin:0">
		   <table cellpadding="0" cellspacing="0" class="es-header" align="center" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;table-layout:fixed !important;width:100%;background-color:transparent;background-repeat:repeat;background-position:center top">
			 <tr>
			  <td align="center" style="padding:0;Margin:0">
			   <table bgcolor="#ffffff" class="es-header-body" align="center" cellpadding="0" cellspacing="0" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;background-color:#ffffff;width:600px">
				 <tr>
				  <td align="left" style="Margin:0;padding-top:20px;padding-bottom:20px;padding-left:40px;padding-right:40px">
				   <table cellpadding="0" cellspacing="0" width="100%" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px">
					 <tr>
					  <td align="center" valign="top" style="padding:0;Margin:0;width:520px">
					   <table cellpadding="0" cellspacing="0" width="100%" role="presentation" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px">
						 <tr>
						
						 </tr>
					   </table></td>
					 </tr>
				   </table></td>
				 </tr>
			   </table></td>
			 </tr>
		   </table>
		   <table cellpadding="0" cellspacing="0" class="es-content" align="center" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;table-layout:fixed !important;width:100%">
			 <tr>
			  <td align="center" style="padding:0;Margin:0">
			   <table bgcolor="#f6f6f6" class="es-content-body" align="center" cellpadding="0" cellspacing="0" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;background-color:#f6f6f6;border-radius:20px 20px 0px 0px;width:600px">
				 <tr>
				  <td align="left" style="padding:0;Margin:0;padding-top:20px;padding-left:40px;padding-right:40px">
				   <table cellpadding="0" cellspacing="0" width="100%" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px">
					 <tr>
					  <td align="center" valign="top" style="padding:0;Margin:0;width:520px">
					   <table cellpadding="0" cellspacing="0" width="100%" role="presentation" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px">
						 <tr>
						  <td align="left" bgcolor="#ffffff" style="padding:10px;Margin:0"><h3 style="Margin:0;line-height:19px;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;font-size:16px;font-style:normal;font-weight:bold;color:#e50075"><strong>` + title + `</strong></h3><h3 style="Margin:0;line-height:19px;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;font-size:16px;font-style:normal;font-weight:bold;color:#e50075">` + sub + `</h3></td>
						 </tr>
					   </table></td>
					 </tr>
				   </table></td>
				 </tr>
				 <tr>
				  <td align="left" style="padding:0;Margin:0;padding-top:20px;padding-left:40px;padding-right:40px">
				   <table cellpadding="0" cellspacing="0" width="100%" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px">
					 <tr>
					  <td align="center" valign="top" style="padding:0;Margin:0;width:520px">
					   <table cellpadding="0" cellspacing="0" width="100%" bgcolor="#fafafa" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:separate;border-spacing:0px;background-color:#fafafa;border-radius:10px" role="presentation">
						 <tr>
						  <td align="left" bgcolor="#ffffff" style="padding:10px;Margin:0">
						  ` + content + `
						</td>
						 </tr>
					   </table></td>
					 </tr>
				   </table></td>
				 </tr>
			   </table></td>
			 </tr>
		   </table>
		   <table cellpadding="0" cellspacing="0" class="es-content" align="center" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;table-layout:fixed !important;width:100%">
			 <tr>
			  <td align="center" style="padding:0;Margin:0">
			   <table bgcolor="#f6f6f6" class="es-content-body" align="center" cellpadding="0" cellspacing="0" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;background-color:#F6F6F6;width:600px">
				 <tr>
				  <td align="left" style="Margin:0;padding-bottom:10px;padding-top:20px;padding-left:40px;padding-right:40px">
				   <table cellpadding="0" cellspacing="0" width="100%" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px">
					 <tr>
					  <td align="center" valign="top" style="padding:0;Margin:0;width:520px">
					   <table cellpadding="0" cellspacing="0" width="100%" role="presentation" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px">
						 <tr>
						    
	</tr>
					   </table></td>
					 </tr>
				   </table></td>
				 </tr>
				 <tr>
				  <td align="left" style="padding:0;Margin:0;padding-left:40px;padding-right:40px">
				   <table cellpadding="0" cellspacing="0" width="100%" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px">
					 <tr>
					  <td align="center" valign="top" style="padding:0;Margin:0;width:520px">
					   <table cellpadding="0" cellspacing="0" width="100%" role="presentation" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px">
						 <tr>
						  <td align="center" style="padding:0;Margin:0;padding-top:10px;padding-bottom:20px;font-size:0">
						   <table border="0" width="100%" height="100%" cellpadding="0" cellspacing="0" role="presentation" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px">
							 <tr>
							  <td style="padding:0;Margin:0;border-bottom:1px solid #000000;background:unset;height:1px;width:100%;margin:0px"></td>
							 </tr>
						   </table></td>
						 </tr>
					   </table></td>
					 </tr>
				   </table></td>
				 </tr>
			   </table></td>
			 </tr>
		   </table>
		   <table cellpadding="0" cellspacing="0" class="es-footer" align="center" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;table-layout:fixed !important;width:100%;background-color:transparent;background-repeat:repeat;background-position:center top">
			 <tr>
			  <td align="center" style="padding:0;Margin:0">
			   <table bgcolor="#bcb8b1" class="es-footer-body" align="center" cellpadding="0" cellspacing="0" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;background-color:#FFFFFF;width:600px">
				 <tr>
				  <td align="left" bgcolor="#f6f6f6" style="padding:10px;Margin:0;background-color:#f6f6f6">
				   <table cellpadding="0" cellspacing="0" width="100%" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px">
					 <tr>
					  <td align="left" style="padding:0;Margin:0;width:580px">
					   <table cellpadding="0" cellspacing="0" width="100%" role="presentation" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px">
						 <tr>
						  <td align="center" class="es-m-txt-c" style="padding:0;Margin:0;font-size:0px"><a target="_blank" href="https://wopta.it" style="-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;text-decoration:underline;color:#2D3142;font-size:9px"><img src="https://mkufhg.stripocdn.email/content/guids/CABINET_2d4d939d8981dffb82d9eb4d43948346a115534e6e2e3c8ba27eaeab5e33b411/images/artw_logo_rgb_7i8.png" alt="Logo" style="display:block;border:0;outline:none;text-decoration:none;-ms-interpolation-mode:bicubic" title="Logo" height="60"></a></td>
						 </tr>
						 <tr>
						  <td align="center" class="es-m-txt-c" style="padding:10px;Margin:0;font-size:0">
						   <table cellpadding="0" cellspacing="0" class="es-table-not-adapt es-social" role="presentation" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px">
							 <tr>
							  <td align="center" valign="top" style="padding:0;Margin:0;padding-right:30px"><a target="_blank" href="https://www.facebook.com/woptaassicurazioni/" style="-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;text-decoration:underline;color:#2D3142;font-size:9px"><img src="https://mkufhg.stripocdn.email/content/assets/img/social-icons/rounded-colored/facebook-rounded-colored.png" alt="Fb" title="Facebook" height="24" style="display:block;border:0;outline:none;text-decoration:none;-ms-interpolation-mode:bicubic"></a></td>
							  <td align="center" valign="top" style="padding:0;Margin:0;padding-right:30px"><a target="_blank" href="https://www.linkedin.com/company/wopta-assicurazioni" style="-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;text-decoration:underline;color:#2D3142;font-size:9px"><img src="https://mkufhg.stripocdn.email/content/assets/img/social-icons/rounded-colored/linkedin-rounded-colored.png" alt="In" title="Linkedin" height="24" style="display:block;border:0;outline:none;text-decoration:none;-ms-interpolation-mode:bicubic"></a></td>
							  <td align="center" valign="top" style="padding:0;Margin:0;padding-right:30px"><a target="_blank" href="https://www.instagram.com/wopta_assicurazioni/" style="-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;text-decoration:underline;color:#2D3142;font-size:9px"><img src="https://mkufhg.stripocdn.email/content/assets/img/social-icons/rounded-colored/instagram-rounded-colored.png" alt="Ig" title="Instagram" height="24" style="display:block;border:0;outline:none;text-decoration:none;-ms-interpolation-mode:bicubic"></a></td>
							  <td align="center" valign="top" style="padding:0;Margin:0"><a target="_blank" href="mailto:info@wopta.it?subject=Richiesta%20info" style="-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;text-decoration:underline;color:#2D3142;font-size:9px"><img src="https://mkufhg.stripocdn.email/content/assets/img/other-icons/rounded-colored/mail-rounded-colored.png" alt="Email" title="Email" height="24" style="display:block;border:0;outline:none;text-decoration:none;-ms-interpolation-mode:bicubic"></a></td>
							 </tr>
						   </table></td>
						 </tr>
						 <tr>
						  <td align="center" style="padding:10px;Margin:0"><p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:11px;color:#585a5d;font-size:9px"><a target="_blank" style="-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;text-decoration:underline;color:#2D3142;font-size:9px" href=""></a></p><p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:11px;color:#585a5d;font-size:9px">© Wopta Assicurazioni s.r.l. | P. IVA 12072020964 | Galleria del Corso, 1 – 20122 Milano (MI).</p><p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:11px;color:#585a5d;font-size:9px">Wopta Assicurazioni s.r.l. è un intermediario assicurativo soggetto alla vigilanza dell’IVASS ed iscritto alla Sezione A del Registro Unico</p><p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:11px;color:#585a5d;font-size:9px">degli Intermediari Assicurativi con numero A000701923. Consulta gli estremi dell’iscrizione al sito <a href="https://servizi.ivass.it/RuirPubblica/" style="-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;text-decoration:underline;color:#2D3142;font-size:9px;font-family:arial, 'helvetica neue', helvetica, sans-serif">servizi.ivass.it</a></p><p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:11px;color:#585a5d;font-size:9px"><a target="_blank" style="-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;text-decoration:underline;color:#2D3142;font-size:9px" href=""></a></p></td>
						 </tr>
					   </table></td>
					 </tr>
				   </table></td>
				 </tr>
			   </table></td>
			 </tr>
		   </table></td>
		 </tr>
	   </table>
	  </div>
	 </body>
	</html>`
}
