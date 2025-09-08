package lib

import "strings"

func JoinNoEmptyStrings(s []string, sep string) string {
	res := strings.Builder{}
	if len(s) == 0 {
		return ""
	}
	for i := range s {
		if s[i] == "" {
			continue
		}
		if res.Len() != 0 {
			res.WriteString(sep)
		}
		res.WriteString(s[i])
	}
	return res.String()
}
