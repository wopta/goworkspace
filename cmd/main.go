package main

import (
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/joho/godotenv"
	_ "github.com/wopta/goworkspace/auth"
	_ "github.com/wopta/goworkspace/broker"
	_ "github.com/wopta/goworkspace/callback"
	_ "github.com/wopta/goworkspace/claim"
	_ "github.com/wopta/goworkspace/companydata"
	_ "github.com/wopta/goworkspace/document"
	_ "github.com/wopta/goworkspace/enrich"
	_ "github.com/wopta/goworkspace/form"
	_ "github.com/wopta/goworkspace/mail"
	_ "github.com/wopta/goworkspace/mga"
	_ "github.com/wopta/goworkspace/partnership"
	_ "github.com/wopta/goworkspace/policy"
	_ "github.com/wopta/goworkspace/question"
	_ "github.com/wopta/goworkspace/quote"
	_ "github.com/wopta/goworkspace/renew"
	_ "github.com/wopta/goworkspace/reserved"
	_ "github.com/wopta/goworkspace/rules"
	_ "github.com/wopta/goworkspace/sellable"
	_ "github.com/wopta/goworkspace/test"
	_ "github.com/wopta/goworkspace/user"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	log.Println(port)

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
