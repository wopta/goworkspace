package main

import (
	"log"
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
		log.Fatalf("Error loading .env file: %s", err)
	}
	//	pl, _ := policy.GetPolicy("Ew13HXDECqGC1FIQ3GbI")
	//	net := network.GetNetworkNodeByUid(pl.ProducerUid)
	//	var warrant *models.Warrant
	//	if net != nil {
	//		warrant = net.GetWarrant()
	//	}
	//	pr := product.GetProductV2(pl.Name, pl.ProductVersion, pl.Channel, net, warrant)
	//	doc := <-document.ContractObj(pl, net, pr)
	//	os.WriteFile("prova1.pdf", doc.Bytes, os.ModePerm)
	//	return
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	log.Println(port)

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
