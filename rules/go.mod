module github.com/wopta/goworkspace/rules

go 1.16

replace github.com/wopta/goworkspace/rules => ./

require github.com/hyperjumptech/grule-rule-engine v1.11.0

require cloud.google.com/go/storage v1.25.0 // indirect

require (
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1
	github.com/go-gota/gota v0.12.0
	github.com/wopta/goworkspace/lib v0.0.0-20221009154939-92933ae1b6d4
	github.com/wopta/goworkspace/models v0.0.0-20220909121553-d232bcdeb3e0
	golang.org/x/net v0.0.0-20220809184613-07c6da5e1ced // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
