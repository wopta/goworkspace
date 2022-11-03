module github.com/wopta/goworkspace/broker

go 1.16

replace github.com/wopta/goworkspace/broker => ./

require (
	firebase.google.com/go/v4 v4.8.0 // indirect
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/hyperjumptech/grule-rule-engine v1.11.0 // indirect
	github.com/wopta/goworkspace/lib v0.0.0-20221002135750-c1075f44b3b4
	go.uber.org/zap v1.21.0 // indirect
)
