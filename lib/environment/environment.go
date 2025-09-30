package env

import (
	"fmt"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/joho/godotenv"
)

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

func Start() {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Sprintf("Error loading .env file: %s", err))
	}
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	log.Println(port)

	if os.Getenv("GOOGLE_PROJECT_ID") != "positive-apex-350507" || os.Getenv("GOOGLE_STORAGE_BUCKET") != "function-data" {
		log.Println("\a\x1b[49;31mYou are not on dev\x1b[39;49m")
	}
	if err := funcframework.Start(port); err != nil {
		panic(fmt.Sprintf("funcframework.Start: %v\n", err))
	}
}
