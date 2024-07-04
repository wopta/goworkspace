package renew

import (
	"log"

	"github.com/wopta/goworkspace/lib"
)

var (
	paramsHierarchy = []map[string][]string{
		{
			"codeCompany":       []string{"codeCompany"},
			"insuredFiscalCode": []string{"insuredFiscalCode"},
		},
	}
	paramsQuery = map[string]string{
		"codeCompany":       "(JSON_VALUE(p.data, '$.codeCompany') = '%s')",
		"insuredFiscalCode": "(JSON_VALUE(p.data, '$.assets[0].person.fiscalCode) = '%s')",
	}
)

type QueryBuilder interface {
	BuildQuery(map[string]string) string
}

type BigQueryQueryBuilder struct{}

func NewQueryBuilder() BigQueryQueryBuilder {
	return BigQueryQueryBuilder{}
}

func (qb *BigQueryQueryBuilder) BuildQuery(params map[string]string) string {
	paramsKeys := lib.GetMapKeys(params)
	allowedParams := make([]string, 0)

	for _, value := range paramsHierarchy {
		for k, v := range value {
			if lib.SliceContains(paramsKeys, k) {
				allowedParams = v
				break
			}
		}
	}

	log.Println(allowedParams)

	return ""
}
