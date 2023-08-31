module github.com/wopta/goworkspace/payment

go 1.19

replace github.com/wopta/goworkspace/payment => ./

require (
	cloud.google.com/go v0.110.7
	github.com/GoogleCloudPlatform/functions-framework-go v1.7.4
	github.com/google/uuid v1.3.0
	github.com/wopta/goworkspace/document v1.0.93
	github.com/wopta/goworkspace/lib v1.0.67
	github.com/wopta/goworkspace/mail v1.0.32
	github.com/wopta/goworkspace/models v1.1.36
	github.com/wopta/goworkspace/policy v1.0.11
	github.com/wopta/goworkspace/product v1.0.41
	github.com/wopta/goworkspace/transaction v1.0.12
	github.com/wopta/goworkspace/user v1.0.21
)