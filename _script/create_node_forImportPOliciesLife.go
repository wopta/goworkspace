package _script

import (
	"os"
	"slices"
	"strings"

	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/network"
)

func CheckAndCreateNodeForImportLife(path string, runDry bool) {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	fileStr := string(fileBytes)
	lines := strings.Split(fileStr, "\n")
	notFound := []string{}
	found := []string{}
	for _, line := range lines[6 : len(lines)-1] {
		fields := strings.Split(line, ";")
		node, _ := network.GetNetworkNodeByCode(fields[13])
		if node == nil {
			notFound = append(notFound, fields[13])
			continue
		}
		found = append(found, fields[13])
	}
	found = slices.Compact(found)
	notFound = slices.Compact(notFound)
	log.PrintStruct("Found: ", found)
	log.PrintStruct("NotFound: ", notFound)
	if runDry {
		return
	}
	for i := range notFound {
		createNodeForImportLife(notFound[i])
	}

}
