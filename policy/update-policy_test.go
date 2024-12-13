package policy

import (
	"os"
	"testing"
)

type testAddr struct {
	city, postalCode, cityCode string
	res                        bool
}

func TestVerifyManualAddress(t *testing.T) {
	os.Setenv("env", "local-test")

	var inputs = []testAddr{
		{"Polonghera", "12030", "CN", true},
		{"Monta'", "12046", "CN", true},
		{"Agrate Conturbia", "28010", "NO", true},
		{"Asinara Cala D'Oliva", "07046", "SS", true},
		{"Padoa", "35131", "PD", false},
		{"Senorbì", "09040", "Ca", true},
		{"San Nicolò Gerrei", "09040", "ca", true},
		{"Sant’Angelo Limosano", "86020", "cB", true},
	}

	for _, input := range inputs {
		err := verifyManualAddress(input.city, input.postalCode, input.cityCode)
		if (err != nil) && (input.res != false) {
			t.Fatalf("expected %v got %v", input.res, err)
		}
	}
}

type testStr struct {
	in, out string
}

func TestNormalizeString(t *testing.T) {
	var inputs = []testStr{
		{"perchè", "perche"},
		{"perche'", "perche"},
		{"perché", "perche"},
		{"Rocca Ciglie'", "RoccaCiglie"},
		{"Asinara Cala D'Oliva", "AsinaraCalaDOliva"},
		{"St. Martin in Thurn/S. Martin de Tor", "StMartininThurnSMartindeTor"},
		{"Mühlen/Truden", "MuhlenTruden"},
		{"St. Andrä_", "StAndra"},
		{"Weißenbach/Sarntal", "WeissenbachSarntal"},
		{"Vöran", "Voran"},
		{"Astfeld-Nordheim", "AstfeldNordheim"},
	}

	for _, input := range inputs {
		got := normalizeString(input.in)
		if got != input.out {
			t.Fatalf("expected %v got %v", input.out, got)
		}
	}
}
