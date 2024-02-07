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
		dbNodes       []models.NetworkNode
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

	// load dataframe
	df := lib.CsvToDataframe(data)
	log.Printf("#rows: %02d #cols: %02d", df.Nrow(), df.Ncol())
	for _, row := range df.Records()[1:] {
		if row[2] == models.AgencyNetworkNodeType {
			log.Printf(models.AgencyNetworkNodeType)
		} else if row[2] == models.AgentNetworkNodeType {
			log.Printf(models.AgentNetworkNodeType)
		}
	}

	// load all nodes from Firestore

	if req.StartPipeline != nil {
		startPipeline = *req.StartPipeline
	}

	if startPipeline {
		log.Printf("Saving import file to Google Bucket...")
		filePath := fmt.Sprintf("dataflow/in_network_node/%s", "prova.csv")
		_, err = lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filePath, data)
		if err != nil {
			log.Printf("Error saving import file to Google Bucket: %s", err.Error())
			return "", nil, err
		}
		log.Printf("Import file saved into Google Bucket")

		// TODO: start dataflow pipeline
	}

	log.Println("Handler End -------------------------------------------------")

	return "", nil, nil
}
