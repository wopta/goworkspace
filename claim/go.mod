module github.com/wopta/goworkspace/claim

go 1.16

replace github.com/wopta/goworkspace/claim => ./

require (
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1
	github.com/google/uuid v1.3.0

	github.com/wopta/goworkspace/lib v0.0.0-20221129170123-a0bf01f8ecd4
	github.com/wopta/goworkspace/mail v0.0.0-20221129170123-a0bf01f8ecd4
	github.com/wopta/goworkspace/models v0.0.0-20221129170123-a0bf01f8ecd4
)
