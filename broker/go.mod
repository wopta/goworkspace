module github.com/wopta/goworkspace/broker

go 1.16

replace github.com/wopta/goworkspace/broker => ./

require (
	firebase.google.com/go/v4 v4.8.0 // indirect
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1
	github.com/wopta/goworkspace/lib v0.0.0-20221109230317-146ed1c0976a
	github.com/wopta/goworkspace/models v0.0.0-20221109230317-146ed1c0976a // indirect

)
