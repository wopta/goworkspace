module github.com/wopta/goworkspace/broker

go 1.16

replace github.com/wopta/goworkspace/broker => ./

require (
	cloud.google.com/go v0.105.0
	cloud.google.com/go/firestore v1.9.0
	firebase.google.com/go/v4 v4.10.0 // indirect
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1
	github.com/wopta/goworkspace/document v0.0.0-20230307163124-33286ac35286
	github.com/wopta/goworkspace/lib v0.0.0-20230307163124-33286ac35286
	github.com/wopta/goworkspace/mail v0.0.0-20230307163124-33286ac35286
	github.com/wopta/goworkspace/models v0.0.0-20230307163124-33286ac35286
	github.com/wopta/goworkspace/payment v0.0.0-20230307163124-33286ac35286
)
