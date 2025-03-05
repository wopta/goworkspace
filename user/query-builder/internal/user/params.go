package user

var (
	paramsHierarchy = []map[string][]string{
		{"fiscalCode": []string{"fiscalCode"}},
		{"mail": []string{"mail"}},
	}

	paramsWhereClause = map[string]string{
		"fiscalCode": "(**tableAlias**.fiscalCode = @%s)",
		"mail":       "(**tableAlias**.mail = @%s)",
	}
)
