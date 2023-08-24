module github.com/wopta/goworkspace/auth

go 1.19

replace github.com/wopta/goworkspace/auth => ./

require (
	github.com/GoogleCloudPlatform/functions-framework-go v1.7.4
	github.com/maja42/goval v1.3.1
	github.com/wopta/goworkspace/lib v1.0.66
	github.com/wopta/goworkspace/models v1.1.25
)

require github.com/golang-jwt/jwt/v5 v5.0.0 // indirect
