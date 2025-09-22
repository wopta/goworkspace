package main

import (
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/joho/godotenv"
	_ "gitlab.dev.wopta.it/goworkspace/auth"
	_ "gitlab.dev.wopta.it/goworkspace/broker"
	_ "gitlab.dev.wopta.it/goworkspace/callback"
	_ "gitlab.dev.wopta.it/goworkspace/claim"
	_ "gitlab.dev.wopta.it/goworkspace/companydata"
	_ "gitlab.dev.wopta.it/goworkspace/document"
	_ "gitlab.dev.wopta.it/goworkspace/enrich"
	_ "gitlab.dev.wopta.it/goworkspace/form"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	_ "gitlab.dev.wopta.it/goworkspace/mail"
	_ "gitlab.dev.wopta.it/goworkspace/mga"
	_ "gitlab.dev.wopta.it/goworkspace/partnership"
	_ "gitlab.dev.wopta.it/goworkspace/policy"
	_ "gitlab.dev.wopta.it/goworkspace/question"
	_ "gitlab.dev.wopta.it/goworkspace/quote"
	_ "gitlab.dev.wopta.it/goworkspace/renew"
	_ "gitlab.dev.wopta.it/goworkspace/reserved"
	_ "gitlab.dev.wopta.it/goworkspace/rules"
	_ "gitlab.dev.wopta.it/goworkspace/sellable"
	_ "gitlab.dev.wopta.it/goworkspace/test"
	_ "gitlab.dev.wopta.it/goworkspace/user"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.ErrorF("Error loading .env file: %s", err)
	}
	//	file, _ := os.ReadFile("policies.txt")
	//	fileStr := string(file)
	//	lines := strings.Split(fileStr, "\n")
	//	datasetId := models.WoptaDataset
	//	for i := range lines[:len(lines)-1] {
	//		continue
	//		lib.DeleteRowBigQuery(datasetId, lib.PolicyCollection, fmt.Sprintf("codeCompany==%s", lines[i]))
	//	}

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	log.Println(port)

	if err := funcframework.Start(port); err != nil {
		log.ErrorF("funcframework.Start: %v\n", err)
	}
}
