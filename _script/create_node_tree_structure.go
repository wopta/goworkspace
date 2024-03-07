package _script

import (
	"cloud.google.com/go/bigquery"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"log"
	"strings"
	"time"
)

type nodeInfo struct {
	Uid        string
	Code       string
	Level      int
	BreadCrumb string
	Parents    []nodeInfo
}

type dbNodeInfo struct {
	ChildUid     string                `bigquery:"childUid"`
	Level        int                   `bigquery:"level"`
	BreadCrumb   string                `bigquery:"breadCrumb"`
	ParentUid    string                `bigquery:"parentUid"`
	CreationDate bigquery.NullDateTime `bigquery:"creationDate"`
}

func CreateNodeTreeStructure() {
	log.Println("function start ----------------------------------------------")

	//nodesMap := make(map[string][]Child)

	nodesList, err := getAllNodes()
	if err != nil {
		panic(err)
	}

	nodesList = lib.SliceFilter(nodesList, func(node models.NetworkNode) bool {
		return node.Type != models.PartnershipNetworkNodeType
	})

	dbNodes := lib.SliceFilter(nodesList, func(node models.NetworkNode) bool {
		return node.ParentUid == ""
	})

	firstLevelNodes := lib.SliceMap(dbNodes, func(node models.NetworkNode) models.NetworkNode {
		return node
	})

	toBeVisitedNodes := make([]nodeInfo, 0)
	visitedNodes := make([]nodeInfo, 0)

	for _, node := range firstLevelNodes {
		toBeVisitedNodes = append(toBeVisitedNodes, nodeInfo{
			Uid:        node.Uid,
			Code:       node.Code,
			Level:      1,
			BreadCrumb: node.Code,
			Parents:    []nodeInfo{},
		})
	}

	for len(toBeVisitedNodes) > 0 {
		index := len(toBeVisitedNodes) - 1
		currentNode := toBeVisitedNodes[index]
		children := lib.SliceFilter(nodesList, func(node models.NetworkNode) bool {
			return node.ParentUid == currentNode.Uid
		})
		toBeVisitedNodes = toBeVisitedNodes[:index]
		parents := append(currentNode.Parents, currentNode)
		for _, child := range children {
			toBeVisitedNodes = append(toBeVisitedNodes, nodeInfo{
				Uid:        child.Uid,
				Code:       child.Code,
				Level:      currentNode.Level + 1,
				BreadCrumb: currentNode.BreadCrumb + " > " + child.Code,
				Parents:    parents,
			})
		}
		visitedNodes = append(visitedNodes, currentNode)

	}

	for _, nn := range visitedNodes {
		if len(nn.Parents) > 0 {
			dbNode := dbNodeInfo{
				ChildUid:     nn.Uid,
				Level:        nn.Level,
				BreadCrumb:   "",
				ParentUid:    "",
				CreationDate: lib.GetBigQueryNullDateTime(time.Now().UTC()),
			}
			splittedBreadCrumb := strings.Split(nn.BreadCrumb, " > ")
			for _, p := range nn.Parents {
				log.Printf("child: %s\tparent: %s\t childLevel: %02d\tparentLevel: %02d\t relativeBreadCrumb: %s\tabsoluteBreadCrumb: %s\t\n", nn.Code, p.Code,
					nn.Level, p.Level, strings.Join(splittedBreadCrumb[p.Level:], " > "), nn.BreadCrumb)
				dbNode.ParentUid = p.Uid
				dbNode.BreadCrumb = strings.Join(splittedBreadCrumb[p.Level:], " > ")
				err = lib.InsertRowsBigQuery("wopta", "node-tree-structure", dbNode)
				if err != nil {
					log.Printf("insert error: %s", err.Error())
					return
				}
			}
		} /*else {
			log.Printf("child: %s\tparent: %s\tbreadCrumb: %s\tlevel: %02d\n", nn.Code, "",
				nn.BreadCrumb, nn.Level)
			lib.InsertRowsBigQuery("wopta", "node-tree-structure", nn)
		}*/
	}

	log.Println("function end ------------------------------------------------")
}

func getAllNodes() ([]models.NetworkNode, error) {
	log.Println("fetching all nodes from Firestore...")
	allNodes, err := network.GetAllNetworkNodes()
	if err != nil {
		log.Printf("error fetching all nodes from Firestore: %s", err.Error())
		return nil, err
	}
	filteredNodes := lib.SliceFilter(allNodes, func(node models.NetworkNode) bool {
		return node.Type != models.PartnershipNetworkNodeType
	})
	log.Printf("found %02d nodes", len(filteredNodes))
	return filteredNodes, nil
}
