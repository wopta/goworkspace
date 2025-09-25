package env

import "os"

type Environment = string

const (
	Local       Environment = "local"
	LocalTest   Environment = "local-test"
	Development Environment = "dev"
	Uat         Environment = "uat"
	Production  Environment = "prod"
)

func IsLocal() bool {
	return os.Getenv("env") == Local
}

func IsLocalTest() bool {
	return os.Getenv("env") == LocalTest
}

func IsDevelopment() bool {
	return os.Getenv("env") == Development
}

func IsProduction() bool {
	return os.Getenv("env") == Production
}

func GetExecutionId() string {
	return os.Getenv("Execution-Id")
}
