package network

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type ImportNodesReq struct {
	Filename      string `json:"filename"`
	Bytes         string `json:"bytes"`
	MimeType      string `json:"mimeType"`
	StartPipeline *bool  `json:"startPipeline,omitempty"`
}

type ErrorCategories struct {
	DuplicatedNodes           []string `json:"duplicatedNodes"`
	InvalidConfigurationNodes []string `json:"invalidConfigurationNodes"`
}

type ImportNodesResp struct {
	TotalInputNodes int             `json:"totalInputNodes"`
	TotalErrorNodes int             `json:"totalErrorNodes"`
	TotalValidNodes int             `json:"totalValidNodes"`
	ErrorNodes      ErrorCategories `json:"errorNodes"`
	ValidNodes      []string        `json:"validNodes"`
}

type nodeInfo struct {
	Warrant        string
	IsActive       bool
	HasAnnex       bool
	IsMgaProponent bool
	Type           string
	RuiSection     string
}

var (
	boolMap = map[string]bool{
		"NO": false,
		"SI": true,
	}
	nodeTypeList = []string{
		models.AgentNetworkNodeType,
		models.AgencyNetworkNodeType,
	}
	designationsList = []string{
		"Addetto Attività intermediazione al di fuori dei locali",
		"Addetto Attività intermediazione all'interno dei locali",
		"Responsabile dell'attività di distribuzione",
		"Responsabile dell'attività di intermediazione",
	}
)

func ImportNodesFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err           error
		req           ImportNodesReq
		resp          ImportNodesResp
		startPipeline = false
		warrants      []models.Warrant
		dbNodes       []models.NetworkNode
		nodesMap      = make(map[string]nodeInfo)
		warrantsMap   = make(map[string][]string)
		validatedRows = make(map[string][][]string)
	)

	log.SetPrefix("ImportNodesFx ")
	defer log.SetPrefix("")

	log.Println("Handler Start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	log.Printf("Request body: %s", string(body))
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("Error unmarshiling request body: %s", err.Error())
		return "{}", nil, err
	}

	// check file mimetype
	if req.MimeType != "text/csv" {
		log.Printf("File format %s not supported", req.MimeType)
		return "{}", nil, errors.New("file format not supported")
	}

	// convert csv to bytes
	data, err := base64.StdEncoding.DecodeString(req.Bytes)
	if err != nil {
		log.Printf("Error decoding file: %s", err.Error())
		return "{}", nil, err
	}

	// load dataframe
	df := lib.CsvToDataframe(data)
	log.Printf("#rows: %02d #cols: %02d", df.Nrow(), df.Ncol())

	outputRows := [][]string{df.Records()[0]}

	// load all nodes from Firestore
	log.Printf("Fetching all network nodes from Firestore...")
	dbNodes, err = GetAllNetworkNodes()
	if err != nil {
		log.Printf("Error fetching all network nodes from Firestore: %s", err.Error())
		return "{}", nil, err
	}
	log.Printf("#Network nodes fetched from Firestore: %02d", len(dbNodes))

	//load all warrant from Google Bucket
	log.Printf("Loading all warrants from Google Bucket...")
	warrants, err = getWarrants()
	if err != nil {
		log.Printf("Error loading warrants from Google Bucket: %s", err.Error())
		return "{}", nil, err
	}
	log.Printf("#Warrants loaded from Google Bucket: %02d", len(warrants))

	// build map[warrant_name] = allowed sub warrants
	warrantsMap = buildWarrantsCompatibilityMap(warrants)

	// build map[networkcode] = nodeInfo with essentials node info
	nodesMap = buildNetworkNodesMap(dbNodes)

	// init resp object

	resp = ImportNodesResp{
		TotalInputNodes: len(df.Records()[1:]),
		TotalErrorNodes: 0,
		TotalValidNodes: 0,
		ErrorNodes: ErrorCategories{
			DuplicatedNodes:           make([]string, 0),
			InvalidConfigurationNodes: make([]string, 0),
		},
		ValidNodes: make([]string, 0),
	}

	// validate csv rows

	for _, row := range df.Records()[1:] {
		// normalize cells content if err add to skipped rows
		row = normalizeFields(row)

		// check if all required fields have been compiled
		err = validateRow(row)
		if err != nil {
			log.Printf("Error processing node %s: %s", row[0], err.Error())
			resp.ErrorNodes.InvalidConfigurationNodes = append(resp.ErrorNodes.InvalidConfigurationNodes, row[0])
			resp.TotalErrorNodes++
			continue
		}

		// validated rows by node type
		validatedRows[row[2]] = append(validatedRows[row[2]], row)
	}

	if validatedRows[models.AgencyNetworkNodeType] != nil {
		duplicatedNodes, invalidNodes, validNodes, validRows := nodeConfigurationValidation(models.AgencyNetworkNodeType, validatedRows[models.AgencyNetworkNodeType], nodesMap, warrantsMap)
		outputRows = append(outputRows, validRows...)
		resp.ErrorNodes.DuplicatedNodes = append(resp.ErrorNodes.DuplicatedNodes, duplicatedNodes...)
		resp.ErrorNodes.InvalidConfigurationNodes = append(resp.ErrorNodes.InvalidConfigurationNodes, invalidNodes...)
		resp.ValidNodes = append(resp.ValidNodes, validNodes...)
	}

	if validatedRows[models.AgentNetworkNodeType] != nil {
		duplicatedNodes, invalidNodes, validNodes, validRows := nodeConfigurationValidation(models.AgentNetworkNodeType, validatedRows[models.AgentNetworkNodeType], nodesMap, warrantsMap)
		outputRows = append(outputRows, validRows...)
		resp.ErrorNodes.DuplicatedNodes = append(resp.ErrorNodes.DuplicatedNodes, duplicatedNodes...)
		resp.ErrorNodes.InvalidConfigurationNodes = append(resp.ErrorNodes.InvalidConfigurationNodes, invalidNodes...)
		resp.ValidNodes = append(resp.ValidNodes, validNodes...)
	}

	resp.TotalValidNodes = len(resp.ValidNodes)
	resp.TotalErrorNodes = len(resp.ErrorNodes.DuplicatedNodes) + len(resp.ErrorNodes.InvalidConfigurationNodes)

	if req.StartPipeline != nil {
		startPipeline = *req.StartPipeline
	}

	if startPipeline && resp.TotalInputNodes == resp.TotalValidNodes {
		var filename string
		splittedFilename := strings.Split(req.Filename, ".")
		if len(splittedFilename) > 2 {
			filename = strings.Join(splittedFilename[:len(splittedFilename)-1], ".")
		} else {
			filename = splittedFilename[0]
		}
		filename += fmt.Sprintf("_%d.%s", time.Now().UTC().Unix(), splittedFilename[len(splittedFilename)-1])
		// write csv to Google Bucket
		err = writeCSVToBucket(outputRows, filename)
		if err != nil {
			return "{}", nil, err
		}

		pubSubClient, err := pubsub.NewClient(context.Background(), os.Getenv("GOOGLE_PROJECT_ID"))
		if err != nil {
			return "{}", nil, err
		}
		topic := pubSubClient.Topic("dataflow")
		topic.Publish(context.Background(), &pubsub.Message{
			Attributes: map[string]string{
				"filename": filename,
			},
		})
		defer topic.Stop()
	}

	log.Printf("#Input Nodes: %d", resp.TotalInputNodes)
	log.Printf("#Invalid Configuration Nodes: %d", resp.TotalErrorNodes)
	log.Printf("#Valid Nodes: %d", resp.TotalValidNodes)

	rawResp, err := json.Marshal(resp)

	log.Println("Handler End -------------------------------------------------")

	return string(rawResp), resp, err
}

func buildNetworkNodesMap(dbNodes []models.NetworkNode) map[string]nodeInfo {
	nodesMap := make(map[string]nodeInfo)
	for _, nn := range dbNodes {
		var ruiSection string
		if nn.Type == models.AgentNetworkNodeType {
			ruiSection = nn.Agent.RuiSection
		} else if nn.Type == models.AgencyNetworkNodeType {
			ruiSection = nn.Agency.RuiSection
		}

		nodesMap[nn.Code] = nodeInfo{
			Warrant:        nn.Warrant,
			IsActive:       nn.IsActive,
			HasAnnex:       nn.HasAnnex,
			IsMgaProponent: nn.IsMgaProponent,
			Type:           nn.Type,
			RuiSection:     ruiSection,
		}

	}
	return nodesMap
}

func getWarrants() ([]models.Warrant, error) {
	var (
		err      error
		warrants []models.Warrant
	)

	warrantsBytes := lib.GetFolderContentByEnv(models.WarrantsFolder)

	for _, warrantBytes := range warrantsBytes {
		var warrant models.Warrant
		err = json.Unmarshal(warrantBytes, &warrant)
		if err != nil {
			log.Printf("[GetWarrants] error unmarshaling warrant: %s", err.Error())
			return warrants, err
		}

		warrants = append(warrants, warrant)
	}
	return warrants, nil
}

func buildWarrantsCompatibilityMap(warrants []models.Warrant) map[string][]string {
	warrantsMap := make(map[string][]string)
	for _, outerWarrant := range warrants {
		warrantsMap[outerWarrant.Name] = make([]string, 0)
		for _, innerWarrant := range warrants {
			compatibleProducts := 0
			for _, innerProduct := range innerWarrant.Products {
				for _, outerProduct := range outerWarrant.Products {
					if innerProduct.Name == outerProduct.Name {
						compatibleProducts++
						break
					}
				}
			}
			if compatibleProducts == len(innerWarrant.Products) {
				warrantsMap[outerWarrant.Name] = append(warrantsMap[outerWarrant.Name], innerWarrant.Name)
			}
		}
	}
	return warrantsMap
}

func normalizeFields(row []string) []string {
	trimFields := []int{0, 1, 5, 9, 10, 13, 23, 24, 31}
	toUpperFields := []int{3, 6, 7, 8, 11, 12, 14, 15, 16, 17, 18, 19, 20, 21, 22, 25, 26, 28, 29, 32}
	toLowerFields := []int{2, 4}

	row = lib.SliceMap(row, func(field string) string {
		if strings.EqualFold(field, "NaN") {
			return ""
		}
		return field
	})

	for _, fieldIndex := range trimFields {
		row[fieldIndex] = lib.TrimSpace(row[fieldIndex])
	}

	for _, fieldIndex := range toLowerFields {
		row[fieldIndex] = lib.ToLower(row[fieldIndex])
	}

	for _, fieldIndex := range toUpperFields {
		row[fieldIndex] = lib.ToUpper(row[fieldIndex])
	}

	return row
}

func validateRow(row []string) error {
	if !lib.SliceContains(nodeTypeList, row[2]) {
		return errors.New("invalid node type")
	}

	var requiredFields []int
	if row[2] == models.AgencyNetworkNodeType {
		requiredFields = []int{0, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 25, 26, 27, 28, 32}
		isMgaProponent := boolMap[row[28]]
		hasAnnex := boolMap[row[29]]
		if isMgaProponent {
			requiredFields = append(requiredFields, 29)
			requiredFields = append(requiredFields, 30)
		} else if !isMgaProponent && hasAnnex {
			requiredFields = append(requiredFields, 29)
		}
	} else if row[2] == models.AgentNetworkNodeType {
		requiredFields = []int{0, 2, 3, 4, 5, 20, 21, 22, 25, 26, 27, 28, 29, 32}
		isMgaProponent := boolMap[row[28]]
		hasAnnex := boolMap[row[29]]
		if isMgaProponent || hasAnnex {
			requiredFields = append(requiredFields, 30)
			requiredFields = append(requiredFields, 31)
		}
	}

	// check fiscalCode format
	fiscalCodePattern := "^(?:[A-Z][AEIOU][AEIOUX]|[AEIOU]X{2}|[B-DF-HJ-NP-TV-Z]{2}[A-Z]){2}(?:[\\dLMNP-V]{2}(?:[" +
		"A-EHLMPR-T](?:[04LQ][1-9MNP-V]|[15MR][\\dLMNP-V]|[26NS][0-8LMNP-U])|[DHPS][37PT][0L]|[ACELMRT][37PT][01LM]|[" +
		"AC-EHLMPR-T][26NS][9V])|(?:[02468LNQSU][048LQU]|[13579MPRTV][26NS])B[26NS][9V])(?:[A-MZ][1-9MNP-V][\\dLMNP-V]" +
		"{2}|[A-M][0L](?:[1-9MNP-V][\\dLMNP-V]|[0L][1-9MNP-V]))[A-Z]$"
	regExp, _ := regexp.Compile(fiscalCodePattern)
	if lib.SliceContains(requiredFields, 22) && !regExp.MatchString(row[22]) {
		return errors.New("invalid fiscal code")
	}

	for fieldIndex, fieldValue := range row {
		if (fieldValue == "" || strings.EqualFold(fieldValue, "NaN")) && lib.SliceContains(requiredFields, fieldIndex) {
			return fmt.Errorf("missing required field at index: %02d", fieldIndex)
		}
	}

	if lib.SliceContains(requiredFields, 30) && !lib.SliceContains(designationsList, row[30]) {
		return errors.New("invalid designation")
	}

	var dateFieldsIndexes = []int{9, 27}
	for _, index := range dateFieldsIndexes {
		if row[index] == "" && lib.SliceContains(requiredFields, index) {
			return fmt.Errorf("missing required field at index: %02d", index)
		}
		_, err := time.Parse("02012006", fmt.Sprintf("%08s", row[index]))
		if err != nil && lib.SliceContains(requiredFields, index) {
			return fmt.Errorf("malformed date at index: %02d", index)
		}
	}

	return nil
}

func nodeConfigurationValidation(nodeType string, rows [][]string, nodesMap map[string]nodeInfo, warrantsMap map[string][]string) ([]string, []string, []string, [][]string) {
	var (
		duplicatedNodes           = make([]string, 0)
		invalidConfigurationNodes = make([]string, 0)
		validNodes                = make([]string, 0)
		outputRows                = make([][]string, 0)
	)

	for _, row := range rows {
		nodeCode := row[0]
		warrantName := row[4]
		parentNodeCode := row[5]
		isMgaProponent := boolMap[row[28]]
		hasAnnex := boolMap[row[29]]
		designation := row[31]
		worksForUid := row[31]

		// check if node is not already present
		if !reflect.ValueOf(nodesMap[nodeCode]).IsZero() {
			log.Printf("Error processing node %s: duplicated node code", nodeCode)
			duplicatedNodes = append(duplicatedNodes, nodeCode)
			continue
		}

		// get father
		parentNode := nodesMap[parentNodeCode]

		// check if parent is present in nodesMap, if not skip
		if reflect.ValueOf(parentNode).IsZero() {
			log.Printf("Error processing node %s: parent node not found", nodeCode)
			invalidConfigurationNodes = append(invalidConfigurationNodes, nodeCode)
			continue
		}

		// check if parent is an agent in nodesMap, if not skip
		if parentNode.Type == models.AgentNetworkNodeType {
			log.Printf("Error processing node %s: agency can't have parent node of type agent", nodeCode)
			invalidConfigurationNodes = append(invalidConfigurationNodes, nodeCode)
			continue
		}

		/*
			check current agency configuration against father configuration, with following checks:
			- check is mga proponent with father
			- check warrant compatibility with father
		*/
		if parentNode.Type != models.AreaManagerNetworkNodeType && parentNode.IsMgaProponent != isMgaProponent {
			log.Printf("Error processing node %s: isMgaProponent configuration not matching parent configuration", nodeCode)
			invalidConfigurationNodes = append(invalidConfigurationNodes, nodeCode)
			continue
		}

		if !lib.SliceContains(warrantsMap[parentNode.Warrant], warrantName) {
			log.Printf("Error processing node %s: warrant configuration not matching parent configuration", nodeCode)
			invalidConfigurationNodes = append(invalidConfigurationNodes, nodeCode)
			continue
		}

		if nodeType == models.AgencyNetworkNodeType {
			// check if fields for simplo are configured correctly
			if worksForUid != "" {
				log.Printf("Error processing node %s: not empty worksForUid", nodeCode)
				invalidConfigurationNodes = append(invalidConfigurationNodes, nodeCode)
				continue
			}

			if isMgaProponent && (!hasAnnex || designation == "") {
				log.Printf("Error processing node %s: invalid node configuration for isMgaProponent = true", nodeCode)
				invalidConfigurationNodes = append(invalidConfigurationNodes, nodeCode)
				continue
			} else if !isMgaProponent && hasAnnex && designation == "" {
				log.Printf("Error processing node %s: invalid node configuration for isMgaProponent = false", nodeCode)
				invalidConfigurationNodes = append(invalidConfigurationNodes, nodeCode)
				continue
			}
		} else if nodeType == models.AgentNetworkNodeType {
			// check if fields for simplo are configured correctly
			if isMgaProponent && (!hasAnnex || designation == "" || worksForUid == "" || (worksForUid != "__wopta__" && nodesMap[worksForUid].RuiSection != "E")) {
				log.Printf("Error processing node %s: invalid node configuration for isMgaProponent = true", nodeCode)
				invalidConfigurationNodes = append(invalidConfigurationNodes, nodeCode)
				continue
			} else if !isMgaProponent && ((hasAnnex && designation == "" && lib.SliceContains([]string{"A", "B"}, nodesMap[worksForUid].RuiSection)) || (!hasAnnex && designation != "" && worksForUid != "")) {
				log.Printf("Error processing node %s: invalid node configuration for isMgaProponent = false", nodeCode)
				invalidConfigurationNodes = append(invalidConfigurationNodes, nodeCode)
				continue
			}
		}

		validNodes = append(validNodes, nodeCode)
		outputRows = append(outputRows, row)
		// add node to nodeMap
		nodesMap[nodeCode] = nodeInfo{
			Warrant:        warrantName,
			HasAnnex:       hasAnnex,
			IsMgaProponent: isMgaProponent,
			Type:           models.AgencyNetworkNodeType,
			RuiSection:     row[8],
		}
	}

	return duplicatedNodes, invalidConfigurationNodes, validNodes, outputRows
}

func writeCSVToBucket(outputRows [][]string, filename string) error {
	tmpFilePath := fmt.Sprintf("../tmp/%s", filename)
	// generate new csv
	err := lib.WriteCsv(tmpFilePath, outputRows, ';')
	if err != nil {
		log.Printf("Error writing csv: %s", err.Error())
		return err
	}
	rawDoc, err := os.ReadFile(tmpFilePath)
	if err != nil {
		log.Printf("Error reading generated csv: %s", err.Error())
		return err
	}

	// upload newly generated csv to Google Bucket
	log.Printf("Saving import file to Google Bucket...")
	filePath := fmt.Sprintf("dataflow/in_network_node/%s", filename)
	_, err = lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filePath, rawDoc)
	if err != nil {
		log.Printf("Error saving import file to Google Bucket: %s", err.Error())
		return err
	}
	log.Printf("Import file saved into Google Bucket")
	return nil
}
