module github.com/wopta/goworkspace/rules

go 1.16

replace github.com/wopta/goworkspace/rules => ./

require (
	cloud.google.com/go v0.105.0
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1
	github.com/go-gota/gota v0.12.0
	github.com/hyperjumptech/grule-rule-engine v1.12.0 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/wopta/goworkspace/lib v0.0.0-20230215143456-09d02c79dcf4
	github.com/wopta/goworkspace/models v0.0.0-20230215143456-09d02c79dcf4
	github.com/wopta/goworkspace/quote v0.0.0-20230215143456-09d02c79dcf4
)
