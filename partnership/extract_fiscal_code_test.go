package partnership

import (
	"testing"
	"time"
)

func TestExtractBirthdate(t *testing.T) {
	tests := []struct {
		fiscalCode string
		birthdate  string
	}{
		{"RSSMRA80A01H501L", "01/01/1980"},
		{"RSSMRA80A41H501L", "01/01/1980"},
		{"MRCMNL80A01H501B", "01/01/1980"},
		{"LRACNO70T24A089L", "24/12/1970"},
	}

	for _, tt := range tests {
		t.Run(tt.fiscalCode, func(t *testing.T) {
			expected, _ := time.Parse("02/01/2006", tt.birthdate)
			actual := extractBirthdateFromItalianFiscalCode(tt.fiscalCode)
			if actual != expected {
				t.Errorf("expected %v, but got %v", expected, actual)
			}
		})
	}
}