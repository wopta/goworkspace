package mail

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
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

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
	log.Println(r.RequestURI)
	lib.EnableCors(&w, r)
	route := lib.RouteData{
		Routes: []lib.Route{{
			Route:   "/v1/send",
			Handler: Send,
			Method:  "POST",
		},
			{
				Route:   "/v1/score",
				Handler: Score,
				Method:  "POST",
			},

			{
				Route:   "/v1/validate",
				Handler: Validate,
				Method:  "POST",
			},
		},
	}
	route.Router(w, r)

	//lib.Files("")

}

type Data struct {
	Title     string
	SubTitle  string
	Content   string
	Link      string
	LinkLabel string
	IsLink    bool
	IsApp     bool
}
type Attachment struct {
	Name        string `firestore:"name,omitempty" json:"name,omitempty"`
	Link        string `firestore:"link,omitempty" json:"link,omitempty"`
	Byte        string `firestore:"byte,omitempty" json:"byte,omitempty"`
	FileName    string `firestore:"fileName,omitempty" json:"fileName,omitempty"`
	MimeType    string `firestore:"mimeType,omitempty" json:"mimeType,omitempty"`
	Url         string `firestore:"url,omitempty" json:"url,omitempty"`
	ContentType string `firestore:"contentType,omitempty" json:"contentType,omitempty"`
}
type MailRequest struct {
	From         string        `json:"from"`
	To           []string      `json:"to"`
	Message      string        `json:"message"`
	Subject      string        `json:"subject"`
	IsHtml       bool          `json:"isHtml,omitempty"`
	IsAttachment bool          `json:"isAttachment,omitempty"`
	Attachments  *[]Attachment `json:"attachments,omitempty"`
	Cc           string        `json:"cc,omitempty"`
	TemplateName string        `json:"templateName,omitempty"`
	Title        string        `json:"title,omitempty"`
	SubTitle     string        `json:"subTitle,omitempty"`
	Content      string        `json:"content,omitempty"`
	Link         string        `json:"link,omitempty"`
	LinkLabel    string        `json:"linkLabel,omitempty"`
	IsLink       bool          `json:"isLink,omitempty"`
	IsApp        bool          `json:"isApp,omitempty"`
}
type MailValidate struct {
	Mail      string `firestore:"mail,omitempty" json:"mail,omitempty"`
	IsValid   bool   `firestore:"isValid" json:"isValid"`
	FidoScore int64  `firestore:"fidoScore" json:"fidoScore"`
}

func Send(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	log.Println(string(req))
	var obj MailRequest
	// Unmarshal or Decode the JSON to the interface.
	//json.NewDecoder(req).Decode(&send)
	defer r.Body.Close()

	json.Unmarshal([]byte(req), &obj)
	SendMail(obj)

	return `{"message":"Success send "}`, nil, nil
}
func Score(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var result map[string]string
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	log.Println(string(req))

	// Unmarshal or Decode the JSON to the interface.
	//json.NewDecoder(req).Decode(&send)
	defer r.Body.Close()

	json.Unmarshal([]byte(req), &result)
	ScoreFido(result["email"])

	return `{"message":"Success send "}`, nil, nil
}
func Validate(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var result map[string]string

	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	log.Println(string(req))
	json.Unmarshal([]byte(req), &result)
	defer r.Body.Close()
	resObj := MailValidate{
		Mail:    result["email"],
		IsValid: false,
	}

	resfire := lib.WhereFirestore("mail", "mail", "==", result["email"])
	objmail, _ := ToListData(resfire)
	if len(objmail) > 0 {
		if objmail[0].IsValid {
			res, _ := json.Marshal(resObj)
			return string(res), res, nil
		}

	} else {
		fido := <-ScoreFido(result["email"])
		log.Println(fido.Email.Score)

		if fido.Email.Score >= 480 {
			log.Println("valid")
			resObj.IsValid = true
			res, _ := json.Marshal(resObj)
			return string(res), res, nil
		} else {
			log.Println("invalid")
			lib.PutFirestore("mail", resObj)
			VerifyEmail(result["email"])
		}

	}

	res, e := json.Marshal(resObj)
	lib.CheckError(e)
	log.Println(string(res))
	return string(res), res, nil
}
