package _script

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func ListFacileBrokerLogin() {
	iterator, err := lib.QueryWhereFirestore(lib.NetworkNodesCollection, "callbackConfig.name", "==", "facileBrokerClient")
	if err != nil {
		log.Fatalf("errored getting nodes: %v", err)
	}

	docsnaps, err := iterator.GetAll()
	if err != nil {
		log.Fatalf("errored getting documents: %v", err)
	}

	mailList := make([][]auth.UserIdentifier, 0)
	chunk := make([]auth.UserIdentifier, 0)
	mailList = append(mailList, chunk)
	chunkIdx := 0
	for _, doc := range docsnaps {
		var nn models.NetworkNode
		err := doc.DataTo(&nn)
		if err != nil {
			log.Fatalf("errored marshalloing data: %v", err)
		}
		if len(mailList[chunkIdx]) >= 100 {
			chunk = make([]auth.UserIdentifier, 0)
			mailList = append(mailList, chunk)
			chunkIdx++
		}
		mailList[chunkIdx] = append(mailList[chunkIdx], auth.EmailIdentifier{Email: nn.Mail})
	}

	ctx := context.Background()
	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: os.Getenv("GOOGLE_PROJECT_ID")})
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	type NodeInfo struct {
		Email     string `json:"email"`
		CreatedAt string `json:"createdAt"`
		LastLogin string `json:"lastLogin"`
	}

	nodes := make([]NodeInfo, 0)

	for i := range mailList {
		res, err := authClient.GetUsers(ctx, mailList[i])
		if err != nil {
			log.Fatalf("errored getting auth users: %v", err)
		}

		for j := range res.Users {
			nodes = append(nodes, NodeInfo{
				Email:     res.Users[j].Email,
				CreatedAt: time.Unix(res.Users[j].UserMetadata.CreationTimestamp/1000, 0).Format(time.RFC3339),
				LastLogin: time.Unix(res.Users[j].UserMetadata.LastLogInTimestamp/1000, 0).Format(time.RFC3339),
			})
		}
		for k := range res.NotFound {
			nodes = append(nodes, NodeInfo{
				Email: res.NotFound[k].(auth.EmailIdentifier).Email,
			})
		}
	}

	writer, err := os.Create("./facile-broker_auth-report.csv")
	if err != nil {
		log.Fatalf("error creating file writer: %s", err)
	}
	defer writer.Close()

	csvWriter := csv.NewWriter(writer)

	headers := []string{"Email", "Autenticazione create in", "Ultimo accesso"}

	if err = csvWriter.Write(headers); err != nil {
		log.Fatalf("error writing to csv writer: %s", err)
	}
	for _, n := range nodes {
		if err = csvWriter.Write([]string{n.Email, n.CreatedAt, n.LastLogin}); err != nil {
			log.Fatalf("error writing to csv writer: %s", err)
		}
	}
	csvWriter.Flush()
	if err = csvWriter.Error(); err != nil {
		log.Fatalf("error writing to csv writer: %s", err)
	}
}
