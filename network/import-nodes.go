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
)

type ImportNodesReq struct {
	Filename      string `json:"filename"`
	Bytes         string `json:"bytes"`
	MimeType      string `json:"mimeType"`
	StartPipeline *bool  `json:"startPipeline,omitempty"`
}

type nodeInfo struct {
	Uid            string
	Warrant        string
	IsActive       bool
	HasAnnex       bool
	IsMgaProponent bool
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
		skippedRows   []int
		validatedRows [][]string
	)

	log.SetPrefix("ImportNodesFx ")
	defer log.SetPrefix("")

	log.Println("Handler Start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	log.Printf("Request body: %s", string(body))
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("Error unmarshiling request body: %s", err.Error())
		return "", nil, err
	}

	// check file mimetype
	if req.MimeType != "text/csv" {
		log.Printf("File format %s not supported", req.MimeType)
		return "", nil, errors.New("file format not supported")
	}

	// convert csv to bytes
	data, err := base64.StdEncoding.DecodeString(req.Bytes)
	if err != nil {
		log.Printf("Error decoding file: %s", err.Error())
		return "", nil, err
	}

	// load dataframe
	df := lib.CsvToDataframe(data)
	log.Printf("#rows: %02d #cols: %02d", df.Nrow(), df.Ncol())

	// load all nodes from Firestore
	log.Printf("Fetching all network nodes from Firestore...")
	dbNodes, err = GetAllNetworkNodes()
	if err != nil {
		log.Printf("Error fetching all network nodes from Firestore: %s", err.Error())
		return "", nil, err
	}
	log.Printf("Network nodes fetched from Firestore, #node: %02d", len(dbNodes))

	//load all warrant from Google Bucket
	log.Printf("Loading all warrants from Google Bucket...")
	warrants, err = getWarrants()
	if err != nil {
		log.Printf("Error loading warrants from Google Bucket: %s", err.Error())
		return "", nil, err
	}
	log.Printf("Warrants loaded from Google Bucket, #warrants: %02d", len(warrants))

	// build map[networkcode] = nodeInfo with essentials node info
	for _, nn := range dbNodes {
		nodesMap[nn.Code] = nodeInfo{
			Uid:            nn.Uid,
			Warrant:        nn.Warrant,
			IsActive:       nn.IsActive,
			HasAnnex:       nn.HasAnnex,
			IsMgaProponent: nn.IsMgaProponent,
		}
	}

	// build map[warrant_name] = allowed sub warrants
	// TODO: improve code quality
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

	// validate csv rows

	for rowIndex, row := range df.Records()[1:] {
		// TODO: normalize cells content if err add to skipped rows

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
			skippedRows = append(skippedRows, rowIndex+1)
			continue
		}

		// TODO: check has annex compatibility with father

		// TODO: check is mga proponent with father

		// TODO: check warrant compatibility with father

		// TODO: add node to nodeMap

		validatedRows = append(validatedRows, row)
	}

	log.Printf("#validated rows: %02d", len(validatedRows))

	if req.StartPipeline != nil {
		startPipeline = *req.StartPipeline
	}

	if startPipeline && len(skippedRows) == 0 {
		// TODO: generate new csv

		// TODO: upload newly generated csv to Google Bucket
		log.Printf("Saving import file to Google Bucket...")
		filePath := fmt.Sprintf("dataflow/in_network_node/%s", req.Filename)
		_, err = lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filePath, data)
		if err != nil {
			log.Printf("Error saving import file to Google Bucket: %s", err.Error())
			return "", nil, err
		}
		log.Printf("Import file saved into Google Bucket")

		// TODO: start dataflow pipeline
	} else if len(skippedRows) > 0 {
		// TODO: return error to frontend
	}

	log.Printf("Skipped Rows: %v", skippedRows)

	log.Println("Handler End -------------------------------------------------")

	return "", nil, nil
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

func validateRow(row []string, optionalFields []int) error {
	for rowIndex, rowValue := range row {
		if (rowValue == "" || rowValue == "NaN") && !lib.SliceContains(optionalFields, rowIndex) {
			return errors.New("missing required field")
		}
	}
	return nil
}
