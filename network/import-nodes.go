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

func ImportNodesFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err           error
		req           ImportNodesReq
		startPipeline = false
		warrants      []models.Warrant
		dbNodes       []models.NetworkNode
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

	// load dataframe
	df := lib.CsvToDataframe(data)
	log.Printf("#rows: %02d #cols: %02d", df.Nrow(), df.Ncol())

	// TODO: load all nodes from Firestore

	// TODO: build map[networkcode] = node

	// validate csv rows

	for rowIndex, row := range df.Records()[1:] {
		// TODO: normalize cells content if err add to skipped rows

		// check if all required fields have been compiled

		optionalFields := []int{}
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

		// TODO: substitute values in row (e.g. father code with father uid) if err add to skipped rows

		// TODO: get uid for current node from Firestore

		// TODO: reshape row by adding at the beginning the uid generated at previous step

		validatedRows = append(validatedRows, row)
	}

	// TODO: generate new csv containing validatedRows

	if req.StartPipeline != nil {
		startPipeline = *req.StartPipeline
	}

	if startPipeline {
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
