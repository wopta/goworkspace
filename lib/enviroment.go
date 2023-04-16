package lib

import (
	"log"
	"strings"
)

func GetDatasetByContractorName(name string, dataset string) string {
	var result string
	if strings.Contains(name, "Woptatest") {
		result = "uat_" + dataset
	} else {
		result = dataset
	}
	log.Println("GetDatasetByContractorName: name:", name)
	log.Println("GetDatasetByContractorName result: ", result)
	return result
}
func GetDatasetByEnv(origin string, dataset string) string {
	var result string
	if strings.Contains(origin, "uat") || strings.Contains(origin, "dev") {
		result = "uat_" + dataset
	} else {
		result = dataset
	}
	log.Println("GetDatasetByEnv: name:", origin)
	log.Println("GetDatasetByEnv result: ", result)
	return result
}
