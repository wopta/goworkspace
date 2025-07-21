package _script

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"google.golang.org/api/iterator"
)

func UpdatePartnershipNodeJwt() {
	nodes := make([]models.NetworkNode, 0)
	iter := lib.WhereFirestore(lib.NetworkNodesCollection, "type", "==", "partnership")
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err.Error())
		}
		var networkNode models.NetworkNode
		if err := doc.DataTo(&networkNode); err != nil {
			log.Println(err.Error())
		}
		nodes = append(nodes, networkNode)
	}

	for _, node := range nodes {
		if node.Partnership != nil {
			node.JwtConfig = node.Partnership.JwtConfig
			if node.IsJwtProtected() {
				node.JwtConfig.KeyName = fmt.Sprintf("%s_SIGNING_KEY", strings.ToUpper(node.Partnership.Name))
			}
			node.UpdatedDate = time.Now().UTC()
			err := node.SaveFirestore()
			if err != nil {
				log.Println(err.Error())
			}
			err = node.SaveBigQuery()
			if err != nil {
				log.Println(node.Uid + " error: " + err.Error())
			}
		}
	}

	log.Println("script done")
}
