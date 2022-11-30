module github.com/wopta/goworkspace/claim

go 1.16

replace github.com/wopta/goworkspace/claim => ./

require (
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1
	github.com/google/uuid v1.3.0
	github.com/wopta/goworkspace/lib v0.0.0-20221130203648-5e90af04e9ef
	github.com/wopta/goworkspace/mail v0.0.0-20221130201533-5e8049a2b1db
	github.com/wopta/goworkspace/models v0.0.0-20221130203648-5e90af04e9ef
)
