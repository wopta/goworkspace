package mail

import (
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/lib/log"
	"io"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	lib "github.com/wopta/goworkspace/lib"
	"google.golang.org/api/iterator"
)

func ValidateFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var result map[string]string

	log.AddPrefix("[ValidateFx] ")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	req := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	json.Unmarshal(req, &result)
	resObj := MailValidate{
		Mail:      result["email"],
		FidoScore: 0,
		IsValid:   false,
		IsValidS:  false,
	}

	resfire := lib.WhereFirestore("mail", "mail", "==", result["email"])
	objmail, _ := ToListData(resfire)

	if len(objmail) > 0 {
		if objmail[0].IsValid {
			log.Println("isValid: ", objmail[0])
			objmail[0].IsValidS = true
			res, _ := json.Marshal(objmail[0])
			return string(res), res, nil
		}
	} else {
		fido := <-ScoreFido(result["email"])
		log.Println(fido.Email.Score)
		resObj.FidoScore = fido.Email.Score

		if fido.Email.Score >= 480 {
			log.Println("valid")
			resObj.IsValid = true
			res1, _ := json.Marshal(resObj)
			lib.PutFirestore("mail", resObj)
			return string(res1), res1, nil
		} else {
			log.Println("invalid")
			resObj.IsValid = false
			lib.PutFirestore("mail", resObj)
			VerifyEmail(result["email"])
		}
	}

	res, e := json.Marshal(resObj)
	lib.CheckError(e)

	log.Println("Handler end -------------------------------------------------")

	return string(res), res, nil
}

func VerifyEmail(data string) {
	r := make(chan MailValidate)
	go func() {
		defer close(r)

		log.Println("saved mail")
		var obj MailRequest
		obj.From = "anna@wopta.it"
		obj.To = []string{data}
		obj.Message = `<p>Ciao </p>
			<p>Verifica la tua mail clicando nel bottone sottostante al termine di questo processo potrai continuare l'acquisto </p>`
		obj.Subject = "Wopta Verifica mail"
		obj.IsHtml = true
		obj.IsLink = true
		obj.Link = os.Getenv("WOPTA_BASEURL") + `callback/v1/emailVerify?email=` + data + `&token=` + os.Getenv("WOPTA_TOKEN_API")
		obj.LinkLabel = "Verifica la Mail"
		obj.IsApp = false
		SendMail(obj)

	}()

}
func ToListData(query *firestore.DocumentIterator) ([]MailValidate, []string) {
	log.Println("MailValidate ToListData")
	var result []MailValidate
	var uid []string
	for {
		d, err := query.Next()

		if err != nil {
			log.ErrorF("error")
			if err == iterator.Done {
				log.Println("MailValidate ToListData iterator.Done")
				break
			}

		} else {
			log.Println("else")
			var value MailValidate
			fmt.Printf("d: %v\n", d.Data())
			e := d.DataTo(&value)
			lib.CheckError(e)
			result = append(result, value)
			uid = append(uid, d.Ref.ID)
			log.Println(len(result))
		}

	}
	return result, uid
}
