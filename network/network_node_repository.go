package network

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetNodeByUidErr(uid string) (*models.NetworkNode, error) {
	var node *models.NetworkNode
	docSnapshot, err := lib.GetFirestoreErr(lib.NetworkNodesCollection, uid)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("could not fetch node: %s", err.Error())
	}
	err = docSnapshot.DataTo(&node)
	if err != nil {
		return nil, fmt.Errorf("could not parse node: %s", err.Error())
	}
	return node, err
}

func initNode(node *models.NetworkNode) error {
	if len(node.Uid) == 0 {
		node.Uid = lib.NewDoc(lib.NetworkNodesCollection)
	}
	now := time.Now().UTC()
	node.CreationDate, node.UpdatedDate = now, now
	node.NetworkUid = node.NetworkCode
	node.Role = node.Type
	node.IsActive = true

	if node.Type != models.AreaManagerNetworkNodeType {
		if node.ExternalNetworkCode == "" {
			node.ExternalNetworkCode = node.Code
		}

		warrant := node.GetWarrant()
		if warrant == nil {
			return fmt.Errorf("could not find warrant for node with value '%s'", node.Warrant)
		}

		if node.Products == nil {
			node.Products = make([]models.Product, 0)
			for _, product := range warrant.Products {
				companies := make([]models.Company, 0)
				for _, company := range product.Companies {
					companies = append(companies, models.Company{
						Name:         company.Name,
						ProducerCode: node.Code,
					})
				}
				node.Products = append(node.Products, models.Product{
					Name:      product.Name,
					Companies: companies,
				})
			}
		} else {
			for prodIndex, product := range node.Products {
				for companyIndex, company := range product.Companies {
					if company.ProducerCode == "" {
						node.Products[prodIndex].Companies[companyIndex].ProducerCode = node.Code
					}
				}
			}
		}
	}

	if node.IsMgaProponent {
		node.HasAnnex = true
	}

	return addWorksForUid(node, node)
}

func CreateNode(node models.NetworkNode) (*models.NetworkNode, error) {
	if err := initNode(&node); err != nil {
		return nil, err
	}
	return &node, node.SaveFirestore()
}

func UpdateNode(node models.NetworkNode) error {
	var (
		err          error
		originalNode *models.NetworkNode
	)
	log.AddPrefix("UpdateNode")
	defer log.PopPrefix()
	log.Println("function start ----------------------------------")

	log.Printf("fetching network node %s from Firestore...", node.Uid)

	originalNode, err = GetNodeByUidErr(node.Uid)
	if err != nil {
		log.Printf("error fetching network node from firestore: %s", err.Error())
		return err
	}
	if originalNode == nil {
		return fmt.Errorf("error no node found: %s", node.Uid)

	}

	if originalNode.Code != node.Code {
		err = TestNetworkNodeUniqueness(node.Code)
		if err != nil {
			log.Printf("error testing network node uniqueness: %s", err.Error())
			return err
		}
	}

	if originalNode.AuthId != "" && originalNode.Mail != node.Mail {
		_, err := lib.UpdateUserEmail(node.Uid, node.Mail)
		if err != nil {
			log.Printf("error updating network node mail on Firebase Auth: %s", err.Error())
			return err
		}
	}
	originalNode.ExternalNetworkCode = node.ExternalNetworkCode
	originalNode.Mail = node.Mail
	if originalNode.Warrant != node.Warrant {
		originalNode.Warrant = node.Warrant
		// Update node products - this will clear all existing overrides
		warrant := originalNode.GetWarrant()
		if warrant == nil {
			return fmt.Errorf("could not find warrant for node with value '%s'", node.Warrant)
		}
		originalNode.Products = make([]models.Product, 0)
		for _, product := range warrant.Products {
			companies := make([]models.Company, 0)
			for _, company := range product.Companies {
				companies = append(companies, models.Company{
					Name:         company.Name,
					ProducerCode: node.Code,
				})
			}
			originalNode.Products = append(originalNode.Products, models.Product{
				Name:      product.Name,
				Companies: companies,
			})
		}
	} else {
		originalNode.Products = node.Products
	}

	if originalNode.AuthId != "" && originalNode.IsActive != node.IsActive {
		err = lib.HandleUserAuthenticationStatus(originalNode.Uid, !node.IsActive)
		if err != nil {
			// TODO: in case of error we might want to restore the old email in auth
			log.ErrorF("error updating network node auth status on Firebase Auth: %s", err.Error())
			return err
		}
	}
	originalNode.IsActive = node.IsActive
	originalNode.Designation = node.Designation
	originalNode.HasAnnex = node.HasAnnex
	originalNode.UpdatedDate = time.Now().UTC()

	originalNode.IsMgaProponent = node.IsMgaProponent
	originalNode.HasAnnex = node.HasAnnex
	if originalNode.IsMgaProponent {
		originalNode.HasAnnex = true
	}

	if originalNode.ParentUid != node.ParentUid {
		err = updateNodeTreeRelations(node)
		if err != nil {
			return err
		}
	}
	originalNode.ParentUid = node.ParentUid

	err = addWorksForUid(originalNode, &node)
	if err != nil {
		log.ErrorF("error updating WorksForUid '%s' in network node '%s': %s", originalNode.WorksForUid, originalNode.Uid, err.Error())
		return err
	}

	switch node.Type {
	case models.AgentNetworkNodeType:
		originalNode.Agent = node.Agent
	case models.AgencyNetworkNodeType:
		originalNode.Agency = node.Agency
	case models.BrokerNetworkNodeType:
		originalNode.Broker = node.Broker
	case models.AreaManagerNetworkNodeType:
		originalNode.AreaManager = node.AreaManager
	}

	if originalNode.AuthId == "" {
		originalNode.Code = node.Code
		originalNode.Type = node.Type
		originalNode.Role = node.Type
	}

	log.Printf("writing network node %s in Firestore...", originalNode.Uid)

	err = originalNode.SaveFirestore()
	if err != nil {
		log.ErrorF("error updating network node %s in Firestore", originalNode.Uid)
		return err
	}

	log.Printf("writing network node %s in BigQuery...", originalNode.Uid)

	return originalNode.SaveBigQuery("")
}

func updateNodeTreeRelations(node models.NetworkNode) error {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	if err != nil {
		log.ErrorF("error getting BigQuery client: %s", err.Error())
		return err
	}
	defer client.Close()
	table := fmt.Sprintf("%s.%s", models.WoptaDataset, models.NetworkTreeStructureTable)

	query := `DELETE FROM ` + "`" + table + "`" + `
        WHERE nodeUid IN (SELECT nodeUid FROM ` + "`" + table + "`" + ` WHERE rootUid = @rootUid)
        AND rootUid NOT IN (SELECT nodeUid FROM ` + "`" + table + "`" + ` WHERE rootUid = @rootUid);`

	params := map[string]interface{}{
		"rootUid": node.Uid,
	}

	err = lib.ExecuteQueryBigQuery(query, params)
	if err != nil {
		return err
	}

	query = `INSERT INTO ` + "`" + table + "`" + ` (rootUid, parentUid, nodeUid, relativeLevel, creationDate) 
		SELECT supertree.rootUid, CASE WHEN subtree.nodeUid = @rootUid THEN  @nodeUid ELSE subtree.parentUid END, 
		subtree.nodeUid, supertree.relativeLevel + subtree.relativeLevel + 1, CURRENT_DATETIME()
        FROM ` + "`" + table + "`" + ` AS supertree 
        CROSS JOIN ` + "`" + table + "`" + ` AS subtree
        WHERE subtree.rootUid = @rootUid
        AND supertree.nodeUid = @nodeUid;`

	params = map[string]interface{}{
		"rootUid": node.Uid,
		"nodeUid": node.ParentUid,
	}

	err = lib.ExecuteQueryBigQuery(query, params)
	if err != nil {
		return err
	}
	return nil
}

func GetNetworkNodeByUid(nodeUid string) *models.NetworkNode {
	log.AddPrefix("GetNetworkNodeByUid")
	defer log.PopPrefix()
	if nodeUid == "" {
		log.Println("nodeUid empty")
		return nil
	}

	networkNode, err := GetNodeByUidErr(nodeUid)
	if err != nil {
		log.ErrorF("error getting producer %s from Firestore", nodeUid)
	}
	if networkNode == nil {
		log.ErrorF("error no node found: %s", nodeUid)
	}

	return networkNode
}

func GetAllNetworkNodes() ([]models.NetworkNode, error) {
	var nodes []models.NetworkNode
	log.AddPrefix("GetAllNetworkNodes")
	defer log.PopPrefix()
	docIterator := lib.OrderFirestore(lib.NetworkNodesCollection, "code", firestore.Asc)

	snapshots, err := docIterator.GetAll()
	if err != nil {
		log.ErrorF("error getting nodes from Firestore: %s", err.Error())
		return nodes, err
	}

	for _, snapshot := range snapshots {
		var node models.NetworkNode
		err = snapshot.DataTo(&node)
		if err != nil {
			log.ErrorF("error parsing node %s: %s", snapshot.Ref.ID, err.Error())
		} else {
			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}

func DeleteNetworkNodeByUid(origin, nodeUid string) error {
	log.AddPrefix("DeleteNetworkNodeByUid")
	defer log.PopPrefix()
	if nodeUid == "" {
		log.ErrorF("no nodeUid specified")
		return fmt.Errorf("no nodeUid specified")
	}

	fireNetwork := lib.NetworkNodesCollection
	_, err := lib.DeleteFirestoreErr(fireNetwork, nodeUid)
	return err
}

func UpdateNetworkNodePortfolio(origin string, policy *models.Policy, networkNode *models.NetworkNode) error {
	log.AddPrefix("UpdateNetworkNodePortfolio")
	defer log.PopPrefix()
	if networkNode == nil {
		log.Printf("no networkNode specified")
		return nil
	}

	log.Printf("adding policy %s to networkNode %s portfolio", policy.Uid, networkNode.Uid)

	networkNode.Policies = append(networkNode.Policies, policy.Uid)

	if !lib.SliceContains(networkNode.Users, policy.Contractor.Uid) {
		log.Printf("adding user %s to networkNode %s users list", policy.Contractor.Uid, networkNode.Uid)
		networkNode.Users = append(networkNode.Users, policy.Contractor.Uid)
	}

	networkNode.UpdatedDate = time.Now().UTC()

	log.Printf("saving networkNode %s to Firestore...", networkNode.Uid)
	fireNetwork := lib.NetworkNodesCollection
	err := lib.SetFirestoreErr(fireNetwork, networkNode.Uid, networkNode)
	if err != nil {
		log.Printf("error saving networkNode %s to Firestore: %s", networkNode.Uid, err.Error())
		return err
	}

	log.Printf("saving networkNode %s to BigQuery...", networkNode.Uid)
	err = networkNode.SaveBigQuery(origin)

	return err
}

func GetNodeByUidBigQuery(uid string) (models.NetworkNode, error) {
	query := "select * from `%s.%s` where uid = @uid limit 1"
	query = fmt.Sprintf(query, models.WoptaDataset, lib.NetworkNodesCollection)
	params := map[string]interface{}{"uid": uid}
	nodes, err := lib.QueryParametrizedRowsBigQuery[models.NetworkNode](query, params)

	if len(nodes) == 0 {
		return models.NetworkNode{}, fmt.Errorf("could not find node with uid %s", uid)
	}
	return nodes[0], err
}

func GetAllSubNodesFromNodeBigQuery(uid string) ([]models.NetworkNode, error) {
	query := `WITH
	RECURSIVE network AS (
	SELECT
	  *
	FROM
	  ` + "`%s.%s`" + `
	WHERE
	  uid = @uid
	UNION ALL
	SELECT
	  child.*
	FROM
	  ` + "`%s.%s`" + ` child
	JOIN
	  network n
	ON
	  n.uid = child.parentUid )
  SELECT
	*
  FROM
	network n
  WHERE
	uid <> @uid`
	query = fmt.Sprintf(query, models.WoptaDataset, lib.NetworkNodesCollection, models.WoptaDataset, models.NetworkNodesCollection)
	params := map[string]interface{}{"uid": uid}
	nodes, err := lib.QueryParametrizedRowsBigQuery[models.NetworkNode](query, params)

	if len(nodes) == 0 {
		return []models.NetworkNode{}, fmt.Errorf("could not find node with uid %s", uid)
	}
	return nodes, err
}

func GetAllParentNodesFromNodeBigQuery(uid string) ([]models.NetworkNode, error) {
	query := `WITH
	RECURSIVE network AS (
	SELECT
	  *
	FROM
	  ` + "`%s.%s`" + `
	WHERE
	  uid = @uid
	UNION ALL
	SELECT
	  child.*
	FROM
	  ` + "`%s.%s`" + ` child
	JOIN
	  network n
	ON
	  n.parentUid = child.uid )
  SELECT
	*
  FROM
	network n
  WHERE
	uid <> @uid`
	query = fmt.Sprintf(query, models.WoptaDataset, lib.NetworkNodesCollection, models.WoptaDataset, models.NetworkNodesCollection)
	params := map[string]interface{}{"uid": uid}
	nodes, err := lib.QueryParametrizedRowsBigQuery[models.NetworkNode](query, params)

	if len(nodes) == 0 {
		return []models.NetworkNode{}, fmt.Errorf("could not find node with uid %s", uid)
	}
	return nodes, err
}

func addWorksForUid(originalNode, inputNode *models.NetworkNode) error {
	if inputNode.Type == models.AgentNetworkNodeType && inputNode.WorksForUid != "" && inputNode.WorksForUid != models.WorksForMgaUid {
		worksForNode := GetNetworkNodeByUid(inputNode.WorksForUid)
		// TODO: Check also for broker?
		if worksForNode.Type != models.AgencyNetworkNodeType {
			return fmt.Errorf("worksForUid must reference an agency, got %s", worksForNode.Type)
		}
		if inputNode.IsMgaProponent && !lib.SliceContains(models.GetProponentRuiSections(), worksForNode.Agency.RuiSection) {
			return fmt.Errorf(
				"worksForUid must reference an agency of RuiSection %v when IsMgaProponent is %t",
				models.GetProponentRuiSections(),
				inputNode.IsMgaProponent,
			)
		}
		if !inputNode.IsMgaProponent && !lib.SliceContains(models.GetIssuerRuiSections(), worksForNode.Agency.RuiSection) {
			return fmt.Errorf(
				"worksForUid must reference an agency of RuiSection %v when IsMgaProponent is %t",
				models.GetIssuerRuiSections(),
				inputNode.IsMgaProponent,
			)
		}
	}

	originalNode.WorksForUid = inputNode.WorksForUid
	return nil
}

var ErrNodeNotFound = errors.New("node not found")

func GetNetworkNodeByCode(code string) (*models.NetworkNode, error) {
	var node *models.NetworkNode
	log.AddPrefix("GetNetworkNodeByCode")
	defer log.PopPrefix()
	if code == "" {
		log.Println("code empty")
		return nil, fmt.Errorf("empty code")
	}

	iter := lib.WhereFirestore(lib.NetworkNodesCollection, "code", "==", code)
	nodeDocSnapshot, err := iter.Next()

	if errors.Is(err, iterator.Done) && nodeDocSnapshot == nil {
		log.Println("node not found")
		return nil, ErrNodeNotFound
	}

	if !errors.Is(err, iterator.Done) && err != nil {
		log.Printf("error getting node: %s", err.Error())
		return nil, err
	}

	err = nodeDocSnapshot.DataTo(&node)
	if node == nil || err != nil {
		log.Printf("could not parse node: %s", err)
		return nil, err
	}

	return node, nil
}

func TestNetworkNodeUniqueness(nodeCode string) error {
	_, err := GetNetworkNodeByCode(nodeCode)
	if err == nil {
		return errors.New("node already exists")
	}
	if !errors.Is(err, ErrNodeNotFound) {
		return err
	}
	return nil
}
