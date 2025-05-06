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

func IsProduction() bool {
	return os.Getenv("env") == Production
}
