package _script

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func ListFacileBrokerLogin() {
	log.Println("ListFacileBrokerLogin")

	log.Printf("Executing script in env: %s", os.Getenv("env"))

	log.Println("Getting all nodes from facile broker")
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
		if err := doc.DataTo(&nn); err != nil {
			log.Fatalf("errored marshalling data: %v", err)
		}
		if len(mailList[chunkIdx]) >= 100 {
			chunk = make([]auth.UserIdentifier, 0)
			mailList = append(mailList, chunk)
			chunkIdx++
		}
		mailList[chunkIdx] = append(mailList[chunkIdx], auth.EmailIdentifier{Email: strings.ToLower(nn.Mail)})
	}
	log.Printf("Got %d chunks of up to 100 nodes, with a max of %d nodes", len(mailList), len(mailList)*100)

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
		Email     string
		CreatedAt string
		LastLogin string
	}

	nodes := make([]NodeInfo, 0)

	log.Println("Getting all nodes authetication info")
	for i := range mailList {
		res, err := authClient.GetUsers(ctx, mailList[i])
		if err != nil {
			log.Fatalf("errored getting auth users: %v", err)
		}

		log.Printf("Batch %d: found %d nodes", i+1, len(res.Users))
		for j := range res.Users {
			nodes = append(nodes, NodeInfo{
				Email:     res.Users[j].Email,
				CreatedAt: time.Unix(res.Users[j].UserMetadata.CreationTimestamp/1000, 0).Format(time.DateOnly),
				LastLogin: time.Unix(res.Users[j].UserMetadata.LastLogInTimestamp/1000, 0).Format(time.DateOnly),
			})
		}
		log.Printf("Batch %d: %d nodes not found", i+1, len(res.NotFound))
		for k := range res.NotFound {
			nodes = append(nodes, NodeInfo{
				Email: res.NotFound[k].(auth.EmailIdentifier).Email,
			})
		}
	}

	log.Println("Writing info to CSV output")

	filename := fmt.Sprintf("./_script/%s_facile-broker_auth-report_%s.csv", os.Getenv("env"), time.Now().Format(time.DateOnly))

	writer, err := os.Create(filename)
	if err != nil {
		log.Fatalf("error creating file writer: %s", err)
	}
	defer writer.Close()

	csvWriter := csv.NewWriter(writer)

	headers := []string{"Email", "Autenticazione creata in", "Ultimo accesso"}

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
	log.Println("CSV written. Script done")
}
