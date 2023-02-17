module github.com/wopta/goworkspace/claim

go 1.16

replace github.com/wopta/goworkspace/claim => ./

require (
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1
	github.com/google/uuid v1.3.0
	github.com/wopta/goworkspace/lib v0.0.0-20230209174247-cc44270f2204
	github.com/wopta/goworkspace/mail v0.0.0-20230115160613-dfa851a70521
	github.com/wopta/goworkspace/models v0.0.0-20230207110643-e018cac2446c
)
