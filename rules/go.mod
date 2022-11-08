module github.com/wopta/goworkspace/rules

go 1.16

replace github.com/wopta/goworkspace/rules => ./

require github.com/hyperjumptech/grule-rule-engine v1.11.0

require cloud.google.com/go/storage v1.25.0 // indirect

require (
	cloud.google.com/go/firestore v1.8.0 // indirect
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1
	github.com/go-gota/gota v0.12.0
	github.com/wopta/goworkspace/lib v0.0.0-20221106140834-bf7189807e2e
	github.com/wopta/goworkspace/models v0.0.0-20221106140834-bf7189807e2e

)
