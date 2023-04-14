package mail

import (
	"log"
	"os"

	"cloud.google.com/go/firestore"
	lib "github.com/wopta/goworkspace/lib"
	"google.golang.org/api/iterator"
)

func VerifyEmail(data string) <-chan MailValidate {
	r := make(chan MailValidate)
	go func() {
		defer close(r)
		m := MailValidate{Mail: data, IsValid: false}
		res := lib.WhereFirestore("mail", "email", "==", data)
		objmail, _ := ToListData(res)
		if len(objmail) > 0 {
			if objmail[0].IsValid {
				r <- objmail[0]
			}
		} else {
			lib.PutFirestore("mail", m)
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
			r <- m
		}

	}()
	return r
}
func ToListData(query *firestore.DocumentIterator) ([]MailValidate, []string) {
	log.Println("MailValidate ToListData")
	var result []MailValidate
	var uid []string
	for {
		d, err := query.Next()

		if err != nil {
			log.Println("error")
			if err == iterator.Done {
				log.Println("MailValidate ToListData iterator.Done")
				break
			}

		} else {
			log.Println("else")
			var value MailValidate
			e := d.DataTo(&value)
			lib.CheckError(e)
			result = append(result, value)
			uid = append(uid, d.Ref.ID)
			log.Println(len(result))
		}

	}
	return result, uid
}
