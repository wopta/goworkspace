module github.com/wopta/goworkspace/product

go 1.16

replace github.com/wopta/goworkspace/product => ./

require (
	firebase.google.com/go/v4 v4.10.0 // indirect
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1
	github.com/wopta/goworkspace/lib v0.0.0-20230304110807-c25080172c0e
	github.com/wopta/goworkspace/models v0.0.0-20230304110807-c25080172c0e
)
