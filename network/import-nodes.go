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
	TotalInputNodes int              `json:"totalInputNodes"`
	TotalErrorNodes int              `json:"totalErrorNodes"`
	TotalValidNodes int              `json:"totalValidNodes"`
	ErrorNodes      map[string][]int `json:"errorNodes"`
	ValidNodes      []string         `json:"validNodes"`
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

const (
	expectedColumns        = 33
	fiscalCodeRegexPattern = "^(?:[A-Z][AEIOU][AEIOUX]|[AEIOU]X{2}|[B-DF-HJ-NP-TV-Z]{2}[A-Z]){2}(?:[\\dLMNP-V]{2}(?:[" +
		"A-EHLMPR-T](?:[04LQ][1-9MNP-V]|[15MR][\\dLMNP-V]|[26NS][0-8LMNP-U])|[DHPS][37PT][0L]|[ACELMRT][37PT][01LM]|[" +
		"AC-EHLMPR-T][26NS][9V])|(?:[02468LNQSU][048LQU]|[13579MPRTV][26NS])B[26NS][9V])(?:[A-MZ][1-9MNP-V][\\dLMNP-V]" +
		"{2}|[A-M][0L](?:[1-9MNP-V][\\dLMNP-V]|[0L][1-9MNP-V]))[A-Z]$"
	codeCol                  int = 0
	externalNetworkCodeCol   int = 1
	typeCol                  int = 2
	mailCol                  int = 3
	warrantCol               int = 4
	parentUidCol             int = 5
	agencyNameCol            int = 6
	agencyRuiCodeCol         int = 7
	agencyRuiSectionCol      int = 8
	agencyRuiRegistrationCol int = 9
	agencyVatCodeCol         int = 10
	agencyPecCol             int = 11
	agencyWebsiteCol         int = 12
	agencyPhoneCol           int = 13
	agencyStreetNameCol      int = 14
	agencyStreetNumberCol    int = 15
	agencyLocalityCol        int = 16
	agencyCityCol            int = 17
	agencyPostalCodeCol      int = 18
	agencyCityCodeCol        int = 19
	agentNameCol             int = 20
	agentSurnameCol          int = 21
	agentFiscalCodeCol       int = 22
	agentVatCodeCol          int = 23
	agentPhoneCol            int = 24
	agentRuiCodeCol          int = 25
	agentRuiSectionCol       int = 26
	agentRuiRegistrationCol  int = 27
	isMgaProponentCol        int = 28
	hasAnnexCol              int = 29
	designationCol           int = 30
	worksForUidCol           int = 31
	isActiveCol              int = 32
)

func ImportNodesFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err           error
		req           ImportNodesReq
		resp          ImportNodesResp
		startPipeline = false
		warrants      []models.Warrant
		dbNodes       []models.NetworkNode
		emailsList    []string
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
	if df.Ncol() != expectedColumns {
		log.Printf("#columns isn't correct, expected %02d got %02d", expectedColumns, df.Ncol())
		return "{}", nil, fmt.Errorf("invalid file content")
	}

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
	nodesMap, emailsList = buildNetworkNodesMap(dbNodes)

	// init resp object

	resp = ImportNodesResp{
		TotalInputNodes: len(df.Records()[1:]),
		TotalErrorNodes: 0,
		TotalValidNodes: 0,
		ErrorNodes:      make(map[string][]int),
		ValidNodes:      make([]string, 0),
	}

	// validate csv rows

	for _, row := range df.Records()[1:] {
		// normalize cells content
		row = normalizeFields(row)

		// check if all required fields have been compiled
		columnsErr := validateRow(row)
		if len(columnsErr) > 0 {
			log.Printf("Error processing node %s at indexes: %v", row[codeCol], columnsErr)
			resp.ErrorNodes[row[codeCol]] = columnsErr
			resp.TotalErrorNodes++
			continue
		}

		// validated rows by node type
		validatedRows[row[typeCol]] = append(validatedRows[row[typeCol]], row)
	}

	if validatedRows[models.AgencyNetworkNodeType] != nil {
		validNodes, validRows := nodeConfigurationValidation(resp.ErrorNodes, models.AgencyNetworkNodeType,
			validatedRows[models.AgencyNetworkNodeType], nodesMap, warrantsMap, emailsList)
		outputRows = append(outputRows, validRows...)
		resp.ValidNodes = append(resp.ValidNodes, validNodes...)
	}

	if validatedRows[models.AgentNetworkNodeType] != nil {
		validNodes, validRows := nodeConfigurationValidation(resp.ErrorNodes, models.AgentNetworkNodeType,
			validatedRows[models.AgentNetworkNodeType], nodesMap, warrantsMap, emailsList)
		outputRows = append(outputRows, validRows...)
		resp.ValidNodes = append(resp.ValidNodes, validNodes...)
	}

	resp.TotalValidNodes = len(resp.ValidNodes)
	resp.TotalErrorNodes = len(resp.ErrorNodes)

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

		log.Println("Getting invoker address...")
		authToken, err := models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
		if err != nil {
			log.Printf("Error getting invoker authToken: %s", err.Error())
			return "{}", nil, err
		}
		invokerAddress := authToken.Email
		log.Printf("Invoker address: %s", invokerAddress)

		pubSubClient, err := pubsub.NewClient(context.Background(), os.Getenv("GOOGLE_PROJECT_ID"))
		if err != nil {
			return "{}", nil, err
		}
		topic := pubSubClient.Topic("dataflow")
		topic.Publish(context.Background(), &pubsub.Message{
			Attributes: map[string]string{
				"filename":       filename,
				"invokerAddress": invokerAddress,
				"module":         "in_network_node",
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

func buildNetworkNodesMap(dbNodes []models.NetworkNode) (map[string]nodeInfo, []string) {
	nodesMap := make(map[string]nodeInfo)
	emailsList := make([]string, 0)
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

		if !lib.SliceContains(emailsList, nn.Mail) {
			emailsList = append(emailsList, nn.Mail)
		}
	}
	return nodesMap, emailsList
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
	trimFields := []int{codeCol, externalNetworkCodeCol, parentUidCol, agencyRuiRegistrationCol, agencyVatCodeCol,
		agencyPhoneCol, agentVatCodeCol, agentPhoneCol, worksForUidCol}
	toUpperFields := []int{mailCol, agencyNameCol, agencyRuiCodeCol, agencyRuiSectionCol, agencyPecCol,
		agencyWebsiteCol, agencyStreetNameCol, agencyStreetNumberCol, agencyLocalityCol, agencyCityCol,
		agencyPostalCodeCol, agencyCityCodeCol, agentNameCol, agentSurnameCol, agentFiscalCodeCol, agentRuiCodeCol,
		agentRuiSectionCol, isMgaProponentCol, hasAnnexCol, isActiveCol}
	toLowerFields := []int{typeCol, warrantCol}

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

func validateRow(row []string) []int {
	columnsError := make([]int, 0)
	nodeCode := row[codeCol]

	if !lib.SliceContains(nodeTypeList, row[typeCol]) {
		log.Printf("Error processing node %s: invalid node type %s", nodeCode, row[typeCol])
		columnsError = append(columnsError, typeCol)
	}

	var requiredFields []int
	if row[typeCol] == models.AgencyNetworkNodeType {
		requiredFields = []int{codeCol, typeCol, mailCol, warrantCol, parentUidCol, agencyNameCol, agencyRuiCodeCol,
			agencyRuiSectionCol, agencyRuiRegistrationCol, agencyVatCodeCol, agencyPecCol, agencyWebsiteCol,
			agencyPhoneCol, agencyStreetNameCol, agencyStreetNumberCol, agencyLocalityCol, agencyCityCol,
			agencyPostalCodeCol, agencyCityCodeCol, agentNameCol, agentSurnameCol, agentFiscalCodeCol,
			agentRuiCodeCol, agentRuiSectionCol, agentRuiRegistrationCol, isMgaProponentCol, isActiveCol}
		isMgaProponent := boolMap[row[isMgaProponentCol]]
		hasAnnex := boolMap[row[hasAnnexCol]]
		if isMgaProponent {
			requiredFields = append(requiredFields, hasAnnexCol)
			requiredFields = append(requiredFields, designationCol)
		} else if !isMgaProponent && hasAnnex {
			requiredFields = append(requiredFields, hasAnnexCol)
		}
	} else if row[typeCol] == models.AgentNetworkNodeType {
		requiredFields = []int{codeCol, typeCol, mailCol, warrantCol, parentUidCol, agentNameCol, agentSurnameCol,
			agentFiscalCodeCol, agentRuiCodeCol, agentRuiSectionCol, agentRuiRegistrationCol, isMgaProponentCol,
			hasAnnexCol, isActiveCol}
		isMgaProponent := boolMap[row[isMgaProponentCol]]
		hasAnnex := boolMap[row[hasAnnexCol]]
		if isMgaProponent || hasAnnex {
			requiredFields = append(requiredFields, designationCol)
			requiredFields = append(requiredFields, worksForUidCol)
		}
	}

	// check fiscalCode format
	regExp, _ := regexp.Compile(fiscalCodeRegexPattern)
	if lib.SliceContains(requiredFields, agentFiscalCodeCol) && !regExp.MatchString(row[agentFiscalCodeCol]) {
		log.Printf("Error processing node %s: invalid fiscalCode %s", nodeCode, row[agentFiscalCodeCol])
		columnsError = append(columnsError, agentFiscalCodeCol)
	}

	for fieldIndex, fieldValue := range row {
		if (fieldValue == "" || strings.EqualFold(fieldValue, "NaN")) && lib.SliceContains(requiredFields, fieldIndex) {
			log.Printf("Error processing node %s: missing required field at index %02d", nodeCode, fieldIndex)
			columnsError = append(columnsError, fieldIndex)
		}
	}

	if lib.SliceContains(requiredFields, designationCol) && !lib.SliceContains(designationsList, row[designationCol]) {
		log.Printf("Error processing node %s: invalid designation %s", nodeCode, row[designationCol])
		columnsError = append(columnsError, designationCol)
	}

	var dateFieldsIndexes = []int{agencyRuiRegistrationCol, agentRuiRegistrationCol}
	for _, index := range dateFieldsIndexes {
		if row[index] == "" && lib.SliceContains(requiredFields, index) {
			log.Printf("Error processing node %s: missing required field at index %02d", nodeCode, index)
			columnsError = append(columnsError, index)
		}
		_, err := time.Parse("02012006", fmt.Sprintf("%08s", row[index]))
		if err != nil && lib.SliceContains(requiredFields, index) {
			log.Printf("Error processing node %s: malformed date at index %02d", nodeCode, index)
			columnsError = append(columnsError, index)
		}
	}

	return columnsError
}

func nodeConfigurationValidation(errorNodes map[string][]int, nodeType string, rows [][]string, nodesMap map[string]nodeInfo, warrantsMap map[string][]string, emailsList []string) ([]string, [][]string) {
	var (
		validNodes = make([]string, 0)
		outputRows = make([][]string, 0)
	)

	for _, row := range rows {
		nodeCode := row[codeCol]
		email := row[mailCol]
		warrantName := row[warrantCol]
		parentNodeCode := row[parentUidCol]
		isMgaProponent := boolMap[row[isMgaProponentCol]]
		hasAnnex := boolMap[row[hasAnnexCol]]
		designation := row[designationCol]
		worksForUid := row[worksForUidCol]
		columnsError := make([]int, 0)

		// check if node is not already present
		if !reflect.ValueOf(nodesMap[nodeCode]).IsZero() {
			log.Printf("Error processing node %s: duplicated node code", nodeCode)
			columnsError = append(columnsError, parentUidCol)
		}

		// get father
		parentNode := nodesMap[parentNodeCode]

		// check if parent is present in nodesMap, if not skip
		if reflect.ValueOf(parentNode).IsZero() {
			log.Printf("Error processing node %s: parent node not found", nodeCode)
			columnsError = append(columnsError, parentUidCol)
		}

		// check if parent is an agent in nodesMap, if not skip
		if parentNode.Type == models.AgentNetworkNodeType {
			log.Printf("Error processing node %s: node can't have parent node of type agent", nodeCode)
			columnsError = append(columnsError, parentUidCol)
		}

		// check if node email is unique
		err := checkDuplicatedMails(emailsList, email)
		if err != nil {
			log.Printf("Error processing node %s: email is not unique", nodeCode)
			columnsError = append(columnsError, mailCol)
		}

		/*
			check current agency configuration against father configuration, with following checks:
			- check is mga proponent with father
			- check warrant compatibility with father
		*/
		if parentNode.Type != models.AreaManagerNetworkNodeType && parentNode.IsMgaProponent != isMgaProponent {
			log.Printf("Error processing node %s: isMgaProponent configuration not matching parent configuration", nodeCode)
			columnsError = append(columnsError, parentUidCol, isMgaProponentCol)
		}

		if !lib.SliceContains(warrantsMap[parentNode.Warrant], warrantName) {
			log.Printf("Error processing node %s: warrant configuration not matching parent configuration", nodeCode)
			columnsError = append(columnsError, warrantCol)
		}

		if nodeType == models.AgencyNetworkNodeType {
			// check if fields for simplo are configured correctly
			if worksForUid != "" {
				log.Printf("Error processing node %s: not empty worksForUid", nodeCode)
				columnsError = append(columnsError, worksForUidCol)
			}

			if isMgaProponent && (!hasAnnex || designation == "") {
				log.Printf("Error processing node %s: invalid node configuration for isMgaProponent = true", nodeCode)
				columnsError = append(columnsError, isMgaProponentCol, hasAnnexCol, designationCol)
			} else if !isMgaProponent && hasAnnex && designation == "" {
				log.Printf("Error processing node %s: invalid node configuration for isMgaProponent = false", nodeCode)
				columnsError = append(columnsError, isMgaProponentCol, hasAnnexCol, designationCol)
			}
		} else if nodeType == models.AgentNetworkNodeType {
			// check if fields for simplo are configured correctly
			if isMgaProponent && (!hasAnnex || designation == "" || worksForUid == "" || (worksForUid != models.WorksForMgaUid && nodesMap[worksForUid].RuiSection != "E")) {
				log.Printf("Error processing node %s: invalid node configuration for isMgaProponent = true", nodeCode)
				columnsError = append(columnsError, isMgaProponentCol, hasAnnexCol, designationCol, worksForUidCol)
			} else if !isMgaProponent && ((hasAnnex && designation == "" && lib.SliceContains([]string{"A", "B"}, nodesMap[worksForUid].RuiSection)) || (!hasAnnex && designation != "" && worksForUid != "")) {
				log.Printf("Error processing node %s: invalid node configuration for isMgaProponent = false", nodeCode)
				columnsError = append(columnsError, isMgaProponentCol, hasAnnexCol, designationCol, worksForUidCol)
			}
		}

		if len(columnsError) != 0 {
			if errorNodes[nodeCode] != nil {
				errorNodes[nodeCode] = append(errorNodes[nodeCode], columnsError...)
			} else {
				errorNodes[nodeCode] = columnsError
			}
			continue
		}

		validNodes = append(validNodes, nodeCode)
		outputRows = append(outputRows, row)
		emailsList = append(emailsList, email)
		// add node to nodeMap
		nodesMap[nodeCode] = nodeInfo{
			Warrant:        warrantName,
			HasAnnex:       hasAnnex,
			IsMgaProponent: isMgaProponent,
			Type:           nodeType,
			RuiSection:     row[agencyRuiSectionCol],
		}
	}

	return validNodes, outputRows
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

func checkDuplicatedMails(emailsList []string, inputEmail string) error {
	for _, email := range emailsList {
		if strings.EqualFold(email, inputEmail) {
			return errors.New("duplicated email")
		}
	}
	return nil
}
