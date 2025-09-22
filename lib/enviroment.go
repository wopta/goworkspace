package lib

import (
	"os"
	"strconv"
	"strings"

	"gitlab.dev.wopta.it/goworkspace/lib/log"
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

func GetBoolEnv(key string) bool {
	flag, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		log.ErrorF("error loading %s environment variable", key)
		return false
	}
	return flag
}
func IsEnv(key string) bool {
	env := os.Getenv("env")
	return env == key
}
