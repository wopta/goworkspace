module github.com/wopta/goworkspace/appcheck-proxy

go 1.16

replace github.com/wopta/goworkspace/appcheck-proxy => ./

require (
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1
	github.com/hyperjumptech/grule-rule-engine v1.11.0 // indirect
	github.com/wopta/goworkspace/lib v0.0.0-20220926222829-8d5718324ab5
)
