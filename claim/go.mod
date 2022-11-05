module github.com/wopta/goworkspace/claim

go 1.16

replace github.com/wopta/goworkspace/claim => ./

require (
	cloud.google.com/go/firestore v1.8.0 // indirect
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1
	github.com/hyperjumptech/grule-rule-engine v1.11.0 // indirect
	github.com/wopta/goworkspace/lib v0.0.0-20221105140759-86481980b823
	github.com/wopta/goworkspace/models v0.0.0-20221105140759-86481980b823
)
