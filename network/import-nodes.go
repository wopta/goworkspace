package network

import (
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
	"strings"
)

type ImportNodesReq struct {
	Filename      string `json:"filename"`
	Bytes         string `json:"bytes"`
	MimeType      string `json:"mimeType"`
	StartPipeline *bool  `json:"startPipeline,omitempty"`
}

type nodeInfo struct {
	Warrant        string
	IsActive       bool
	HasAnnex       bool
	IsMgaProponent bool
	Type           string
}

var boolMap = map[string]bool{
	"NO": false,
	"SI": true,
}

func ImportNodesFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err           error
		req           ImportNodesReq
		startPipeline = false
		warrants      []models.Warrant
		dbNodes       []models.NetworkNode
		nodesMap      = make(map[string]nodeInfo)
		warrantsMap   = make(map[string][]string)
		validatedRows = make(map[string][][]string)
		skippedRows   = make(map[string][]string)
	)

	log.SetPrefix("ImportNodesFx ")
	defer log.SetPrefix("")

	log.Println("Handler Start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
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

	header := df.Records()[0]

	// load all nodes from Firestore
	log.Printf("Fetching all network nodes from Firestore...")
	dbNodes, err = GetAllNetworkNodes()
	if err != nil {
		log.Printf("Error fetching all network nodes from Firestore: %s", err.Error())
		return "{}", nil, err
	}
	log.Printf("Network nodes fetched from Firestore, #node: %02d", len(dbNodes))

	//load all warrant from Google Bucket
	log.Printf("Loading all warrants from Google Bucket...")
	warrants, err = getWarrants()
	if err != nil {
		log.Printf("Error loading warrants from Google Bucket: %s", err.Error())
		return "{}", nil, err
	}
	log.Printf("Warrants loaded from Google Bucket, #warrants: %02d", len(warrants))

	// build map[warrant_name] = allowed sub warrants
	warrantsMap = buildWarrantsCompatibilityMap(warrants)

	// build map[networkcode] = nodeInfo with essentials node info
	nodesMap = buildNetworkNodesMap(dbNodes)

	// validate csv rows

	for rowIndex, row := range df.Records()[1:] {
		// normalize cells content if err add to skipped rows
		row = normalizeFields(row)

		// check if all required fields have been compiled
		var optionalFields []int
		if row[2] == models.AgencyNetworkNodeType {
			optionalFields = []int{1, 23, 24}
		} else if row[2] == models.AgentNetworkNodeType {
			optionalFields = []int{1, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 23, 24}
		}
		err = validateRow(row, optionalFields)
		if err != nil {
			log.Printf("Error validating row %02d: %s", rowIndex+1, err.Error())
			skippedRows[row[2]] = append(skippedRows[row[2]], row[0])
			continue
		}

		// validated rows by node type
		validatedRows[row[2]] = append(validatedRows[row[2]], row)
	}

	if validatedRows[models.AgencyNetworkNodeType] != nil {
		for rowIndex, row := range validatedRows[models.AgencyNetworkNodeType] {
			nodeCode := row[0]
			warrantName := row[4]
			fatherNodeCode := row[5]
			isMgaProponent := boolMap[row[28]]
			hasAnnex := boolMap[row[29]]

			// get father
			fatherNode := nodesMap[fatherNodeCode]

			// check if parent is an agent or not present in nodesMap, if so skip
			if reflect.ValueOf(fatherNode).IsZero() || fatherNode.Type == models.AgentNetworkNodeType {
				skippedRows[models.AgencyNetworkNodeType] = append(skippedRows[models.AgencyNetworkNodeType], nodeCode)
				validatedRows[models.AgencyNetworkNodeType] = append(validatedRows[models.AgencyNetworkNodeType][:rowIndex], validatedRows[models.AgencyNetworkNodeType][rowIndex+1:]...)
				continue
			}

			/*
				check current agency configuration against father configuration, with following checks:
				- check has annex compatibility with father
				- check is mga proponent with father
				- check warrant compatibility with father
			*/
			if fatherNode.Type != models.AreaManagerNetworkNodeType && fatherNode.HasAnnex != hasAnnex {
				skippedRows[models.AgencyNetworkNodeType] = append(skippedRows[models.AgencyNetworkNodeType], nodeCode)
				validatedRows[models.AgencyNetworkNodeType] = append(validatedRows[models.AgencyNetworkNodeType][:rowIndex], validatedRows[models.AgencyNetworkNodeType][rowIndex+1:]...)
				continue
			}
			if fatherNode.Type != models.AreaManagerNetworkNodeType && fatherNode.IsMgaProponent != isMgaProponent {
				skippedRows[models.AgencyNetworkNodeType] = append(skippedRows[models.AgencyNetworkNodeType], nodeCode)
				validatedRows[models.AgencyNetworkNodeType] = append(validatedRows[models.AgencyNetworkNodeType][:rowIndex], validatedRows[models.AgencyNetworkNodeType][rowIndex+1:]...)
				continue
			}
			if !lib.SliceContains(warrantsMap[fatherNode.Warrant], warrantName) {
				skippedRows[models.AgencyNetworkNodeType] = append(skippedRows[models.AgencyNetworkNodeType], nodeCode)
				validatedRows[models.AgencyNetworkNodeType] = append(validatedRows[models.AgencyNetworkNodeType][:rowIndex], validatedRows[models.AgencyNetworkNodeType][rowIndex+1:]...)
				continue
			}

			// TODO: check if fields for simplo are configured correctly

			// add node to nodeMap
			nodesMap[nodeCode] = nodeInfo{
				Warrant:        warrantName,
				HasAnnex:       hasAnnex,
				IsMgaProponent: isMgaProponent,
				Type:           models.AgencyNetworkNodeType,
			}
		}
	}

	if validatedRows[models.AgentNetworkNodeType] != nil {
		for rowIndex, row := range validatedRows[models.AgentNetworkNodeType] {
			nodeCode := row[0]
			warrantName := row[4]
			fatherNodeCode := row[5]
			isMgaProponent := boolMap[row[28]]
			hasAnnex := boolMap[row[29]]

			// get father
			fatherNode := nodesMap[fatherNodeCode]

			// check if parent is an agent or not present in nodesMap, if so skip
			if reflect.ValueOf(fatherNode).IsZero() || fatherNode.Type == models.AgentNetworkNodeType {
				skippedRows[models.AgentNetworkNodeType] = append(skippedRows[models.AgentNetworkNodeType], nodeCode)
				validatedRows[models.AgentNetworkNodeType] = append(validatedRows[models.AgentNetworkNodeType][:rowIndex], validatedRows[models.AgentNetworkNodeType][rowIndex+1:]...)
				continue
			}

			/*
				check current agent configuration against father configuration, with following checks:
				- check has annex compatibility with father
				- check is mga proponent with father
				- check warrant compatibility with father
			*/
			if fatherNode.Type != models.AreaManagerNetworkNodeType && fatherNode.HasAnnex != hasAnnex {
				skippedRows[models.AgentNetworkNodeType] = append(skippedRows[models.AgentNetworkNodeType], nodeCode)
				validatedRows[models.AgentNetworkNodeType] = append(validatedRows[models.AgentNetworkNodeType][:rowIndex], validatedRows[models.AgentNetworkNodeType][rowIndex+1:]...)
				continue
			}
			if fatherNode.Type != models.AreaManagerNetworkNodeType && fatherNode.IsMgaProponent != isMgaProponent {
				skippedRows[models.AgentNetworkNodeType] = append(skippedRows[models.AgentNetworkNodeType], nodeCode)
				validatedRows[models.AgentNetworkNodeType] = append(validatedRows[models.AgentNetworkNodeType][:rowIndex], validatedRows[models.AgentNetworkNodeType][rowIndex+1:]...)
				continue
			}
			if !lib.SliceContains(warrantsMap[fatherNode.Warrant], warrantName) {
				skippedRows[models.AgentNetworkNodeType] = append(skippedRows[models.AgentNetworkNodeType], nodeCode)
				validatedRows[models.AgentNetworkNodeType] = append(validatedRows[models.AgentNetworkNodeType][:rowIndex], validatedRows[models.AgentNetworkNodeType][rowIndex+1:]...)
				continue
			}

			// TODO: check if fields for simplo are configured correctly

			// add node to nodeMap
			nodesMap[nodeCode] = nodeInfo{
				Warrant:        warrantName,
				HasAnnex:       hasAnnex,
				IsMgaProponent: isMgaProponent,
				Type:           models.AgentNetworkNodeType,
			}
		}
	}

	if req.StartPipeline != nil {
		startPipeline = *req.StartPipeline
	}

	if startPipeline && len(skippedRows) == 0 {
		// write csv to Google Bucket
		err = writeCSVToBucket(header, validatedRows, req.Filename)
		if err != nil {
			return "{}", nil, err
		}

		// TODO: start dataflow pipeline
	} else if len(skippedRows) > 0 {
		// TODO: return error to frontend
	}

	log.Printf("Skipped Rows: %v", skippedRows)

	log.Println("Handler End -------------------------------------------------")

	return "{}", nil, nil
}

func buildNetworkNodesMap(dbNodes []models.NetworkNode) map[string]nodeInfo {
	nodesMap := make(map[string]nodeInfo)
	for _, nn := range dbNodes {
		nodesMap[nn.Code] = nodeInfo{
			Warrant:        nn.Warrant,
			IsActive:       nn.IsActive,
			HasAnnex:       nn.HasAnnex,
			IsMgaProponent: nn.IsMgaProponent,
			Type:           nn.Type,
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

func validateRow(row []string, optionalFields []int) error {
	for fieldIndex, fieldValue := range row {
		if (fieldValue == "" || strings.EqualFold(fieldValue, "NaN")) && !lib.SliceContains(optionalFields, fieldIndex) {
			return errors.New("missing required field")
		}
	}
	return nil
}

func writeCSVToBucket(header []string, validatedRows map[string][][]string, filename string) error {
	// generate new csv
	outputRows := [][]string{
		header,
	}
	outputRows = append(outputRows, validatedRows[models.AgencyNetworkNodeType]...)
	outputRows = append(outputRows, validatedRows[models.AgentNetworkNodeType]...)

	err := lib.WriteCsv("../tmp/"+filename, outputRows, ';')
	if err != nil {
		log.Printf("Error writing csv: %s", err.Error())
		return err
	}
	rawDoc, err := os.ReadFile("../tmp/" + filename)
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
