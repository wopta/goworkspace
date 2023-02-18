module github.com/wopta/goworkspace/callback

go 1.16

replace github.com/wopta/goworkspace/callback => ./

require (
	firebase.google.com/go/v4 v4.10.0 // indirect
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1
	github.com/wopta/goworkspace/document v0.0.0-20230115160613-dfa851a70521
	github.com/wopta/goworkspace/lib v0.0.0-20230218133800-746a7c429fe6
	github.com/wopta/goworkspace/mail v0.0.0-20230115160613-dfa851a70521
	github.com/wopta/goworkspace/models v0.0.0-20230218133800-746a7c429fe6
)
