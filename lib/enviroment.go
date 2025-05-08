package lib

import (
	"os"
	"strconv"
	"strings"

	"github.com/wopta/goworkspace/lib/log"
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

// DEPRECATED
func GetDatasetByEnv(origin string, dataset string) string {
	var result string
	if strings.Contains(origin, "uat") {
		result = "uat_" + dataset
	} else {
		result = dataset
	}
	log.Println("GetDatasetByEnv: name:", origin)
	log.Println("GetDatasetByEnv result: ", result)
	return result
}

func GetBoolEnv(key string) bool {
	flag, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		log.ErrorF("error loading %s environment variable", key)
		return false
	}
	return flag
}
