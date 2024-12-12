package policy

import (
	"os"
	"testing"
)

type test struct {
	city, postalCode, cityCode string
	res                        bool
}

func TestVerifyManualAddress(t *testing.T) {
	os.Setenv("env", "local-test")

	var inputs = []test{
		{"Polonghera", "12030", "CN", true},
		{"Monta'", "12046", "CN", true},
		{"Agrate Conturbia", "28010", "NO", true},
		{"Asinara Cala D'Oliva", "07046", "SS", true},
		{"Padoa", "35131", "PD", false},
	}

	for _, input := range inputs {
		err := verifyManualAddress(input.city, input.postalCode, input.cityCode)
		if (err != nil) && (input.res != false) {
			t.Fatalf("expected %v got %v", input.res, err)
		}
	}
}
