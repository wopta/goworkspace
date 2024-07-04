package renew_test

import (
	"strings"
	"testing"

	"github.com/wopta/goworkspace/renew"
)

func TestQueryBuilder(t *testing.T) {
	qb := renew.NewQueryBuilder()

	t.Run("codeCompany overcome everything", func(t *testing.T) {
		params := map[string]string{
			"codeCompany":       "100100",
			"insuredFiscalCode": "LLLRRR85E05R94Z330F",
		}

		want := `(JSON_VALUE(p.data, '$.codeCompany') = "100100")`

		got := qb.BuildQuery(params)

		if !strings.EqualFold(got, want) {
			t.Fatalf("expected: %s, got: %s", want, got)
		}
	})

}
