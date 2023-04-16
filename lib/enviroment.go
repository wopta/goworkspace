package lib

import "strings"

func GetDatasetByContractorName(name string, dataset string) string {
	var result string
	if strings.Contains(name, "Woptatest") {
		result = "uat-" + dataset
	} else {
		result = dataset
	}
	return result
}
