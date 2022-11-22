module github.com/wopta/goworkspace/broker

go 1.16

replace github.com/wopta/goworkspace/broker => ./

require (
	firebase.google.com/go/v4 v4.8.0 // indirect
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1
	github.com/hyperjumptech/grule-rule-engine v1.11.0 // indirect
	github.com/wopta/goworkspace/document v0.0.0-20221119181257-41328d67d709 // indirect
	github.com/wopta/goworkspace/lib v0.0.0-20221122094244-96e541ee7842
	github.com/wopta/goworkspace/mail v0.0.0-20221119181257-41328d67d709 // indirect
	github.com/wopta/goworkspace/models v0.0.0-20221122094244-96e541ee7842 // indirect

)
