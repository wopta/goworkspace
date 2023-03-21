module github.com/wopta/goworkspace/rules

go 1.19

replace github.com/wopta/goworkspace/rules => ./

require (
	cloud.google.com/go v0.105.0
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1
	github.com/go-gota/gota v0.12.0
	github.com/hyperjumptech/grule-rule-engine v1.12.0
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/wopta/goworkspace/lib v1.0.3
	github.com/wopta/goworkspace/models v1.0.6
	github.com/wopta/goworkspace/quote v1.0.1
)