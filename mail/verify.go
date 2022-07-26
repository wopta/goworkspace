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
			log.Println("saved")
			var obj MailRequest
			obj.From = "noreply@wopta.it"
			obj.To = []string{data}
			obj.Message = `<p>ciao </p><p>verifica la tua mail clicando il link quii sotto </p></p><p>verifica la tua mail clicando il link qui sotto </p>
			<p><a class="button" href='` + os.Getenv("WOPTA_BASEURL") + `callback/v1/emailVerify?email=` + data + `&token=` + os.Getenv("WOPTA_TOKEN_API") + `'>Veriifica la tua mail</a> </p>`
			obj.Subject = "Wopta Verifica mail"
			obj.IsHtml = true
			SendMail(obj)
			r <- m
		}

	}()
	return r
}
func ToListData(query *firestore.DocumentIterator) ([]MailValidate, []string) {
	var result []MailValidate
	var uid []string
	for {
		d, err := query.Next()

		log.Println(d.Ref.ID)
		if err != nil {
			log.Println("error")
		}
		if err != nil {
			if err == iterator.Done {
				log.Println("iterator.Done")
				break
			}

		}

		var value MailValidate

		e := d.DataTo(&value)

		log.Println("todata")
		lib.CheckError(e)
		result = append(result, value)
		uid = append(uid, d.Ref.ID)

		log.Println(len(result))
	}
	return result, uid
}
