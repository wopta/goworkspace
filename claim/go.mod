module github.com/wopta/goworkspace/claim

go 1.16

replace github.com/wopta/goworkspace/claim => ./

require (
	firebase.google.com/go/v4 v4.10.0 // indirect
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1
	github.com/google/uuid v1.3.0
	github.com/wopta/goworkspace/lib v0.0.0-20230304150751-48959cb12ccd
	github.com/wopta/goworkspace/mail v0.0.0-20230115160613-dfa851a70521
	github.com/wopta/goworkspace/models v0.0.0-20230304150751-48959cb12ccd
)
