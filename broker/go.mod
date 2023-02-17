module github.com/wopta/goworkspace/broker

go 1.16

replace github.com/wopta/goworkspace/broker => ./

require (
	cloud.google.com/go v0.105.0
	cloud.google.com/go/firestore v1.9.0
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1
	github.com/wopta/goworkspace/document v0.0.0-20230210120851-c4c9af79f605
	github.com/wopta/goworkspace/lib v0.0.0-20230210120851-c4c9af79f605
	github.com/wopta/goworkspace/mail v0.0.0-20230210120851-c4c9af79f605
	github.com/wopta/goworkspace/models v0.0.0-20230210120851-c4c9af79f605
	github.com/wopta/goworkspace/payment v0.0.0-20230210120851-c4c9af79f605
)
