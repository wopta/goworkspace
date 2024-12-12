package policy

import (
	"fmt"
	"os"
	"testing"

	"github.com/wopta/goworkspace/lib"
)

type test struct {
	city, postalCode, cityCode string
	res                        bool
}

func TestVerifyManualAddress(t *testing.T) {
	var inputs = []test{
		{"Polonghera", "12030", "CN", true},
		{"Monta'", "12046", "CN", true},
		{"Agrate Conturbia", "28010", "NO", true},
		{"Asinara Cala D'Oliva", "07046", "SS", true},
		{"Padoa", "35131", "PD", false},
	}

	fileName := "enrich/postal-codes.csv"
	res, err := os.ReadFile("../../function-data/dev/" + fileName)
	if err != nil {
		fmt.Printf("error reading file %s: %v", fileName, err)
	}
	df, err := lib.CsvToDataframeV2(res, ';', true)
	if err != nil {
		fmt.Printf("==> error reading df: %v", err)
	}

	for _, input := range inputs {
		err := verifyManualAddress(input.city, input.postalCode, input.cityCode, df)
		if (err != nil) && (input.res != false) {
			t.Fatalf("expected %v got %v", input.res, err)
		}
	}
}
