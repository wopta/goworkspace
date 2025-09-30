package main

import (
	_ "gitlab.dev.wopta.it/goworkspace/auth"
	_ "gitlab.dev.wopta.it/goworkspace/broker"
	_ "gitlab.dev.wopta.it/goworkspace/callback"
	_ "gitlab.dev.wopta.it/goworkspace/claim"
	_ "gitlab.dev.wopta.it/goworkspace/companydata"
	_ "gitlab.dev.wopta.it/goworkspace/document"
	_ "gitlab.dev.wopta.it/goworkspace/enrich"
	_ "gitlab.dev.wopta.it/goworkspace/form"
	env "gitlab.dev.wopta.it/goworkspace/lib/environment"
	_ "gitlab.dev.wopta.it/goworkspace/mail"
	_ "gitlab.dev.wopta.it/goworkspace/mga"
	_ "gitlab.dev.wopta.it/goworkspace/partnership"
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
	defer env.Start()
}
