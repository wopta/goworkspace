package main

import (
	"encoding/base64"

	_ "gitlab.dev.wopta.it/goworkspace/auth"
	_ "gitlab.dev.wopta.it/goworkspace/broker"
	_ "gitlab.dev.wopta.it/goworkspace/callback"
	_ "gitlab.dev.wopta.it/goworkspace/claim"
	_ "gitlab.dev.wopta.it/goworkspace/companydata"
	_ "gitlab.dev.wopta.it/goworkspace/document"
	_ "gitlab.dev.wopta.it/goworkspace/enrich"
	_ "gitlab.dev.wopta.it/goworkspace/form"
	"gitlab.dev.wopta.it/goworkspace/lib"
	env "gitlab.dev.wopta.it/goworkspace/lib/environment"
	_ "gitlab.dev.wopta.it/goworkspace/mail"
	_ "gitlab.dev.wopta.it/goworkspace/mga"
	"gitlab.dev.wopta.it/goworkspace/models/catnat"
	_ "gitlab.dev.wopta.it/goworkspace/partnership"
	"gitlab.dev.wopta.it/goworkspace/policy"
	_ "gitlab.dev.wopta.it/goworkspace/policy"
	_ "gitlab.dev.wopta.it/goworkspace/question"
	_ "gitlab.dev.wopta.it/goworkspace/quote"
	_ "gitlab.dev.wopta.it/goworkspace/renew"
	_ "gitlab.dev.wopta.it/goworkspace/reserved"
	_ "gitlab.dev.wopta.it/goworkspace/rules"
	_ "gitlab.dev.wopta.it/goworkspace/sellable"
	_ "gitlab.dev.wopta.it/goworkspace/test"
	_ "gitlab.dev.wopta.it/goworkspace/user"
)

func main() {
	env.Start(false)
	policy, e := policy.GetPolicy("rn5jC1QATMtzznGizFsD")
	if e != nil {
		panic(e)
	}
	client := catnat.NewNetClient()
	attachmentName := "Contratto NetInsurance"
	var document string
	for i := range *policy.Attachments {
		if (*policy.Attachments)[i].Name == attachmentName {
			bytes, _ := lib.ReadFileFromGoogleStorageEitherGsOrNot((*policy.Attachments)[i].Link)
			document = base64.StdEncoding.EncodeToString(bytes)
			break
		}
	}
	panic(client.UploadDocument(policy, document))
}
