module github.com/wopta/goworkspace/wise-broker

go 1.16

replace github.com/wopta/goworkspace/wise-broker => ./

require (
	firebase.google.com/go/v4 v4.8.0 // indirect
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1

	github.com/wopta/goworkspace/lib v0.0.0-20221002135750-c1075f44b3b4
)
