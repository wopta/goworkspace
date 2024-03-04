package sellable

// import (
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"

// 	"github.com/wopta/goworkspace/models"
// )

// type MockSellableProduct struct{}

// func (m MockSellableProduct) GetProduct(name, version string) *models.Product {
// 	product := models.Product{
// 		Companies: []models.Company{{
// 			GuaranteesMap: map[string]*models.Guarante{
// 				"death":                {IsSellable: false},
// 				"permanent-disability": {IsSellable: false},
// 				"temporary-disability": {IsSellable: false},
// 				"serious-ill":          {IsSellable: false},
// 			},
// 		}},
// 	}

// 	return &product
// }

// func (m MockSellableProduct) GetRulesFile() []byte {
// 	return []byte(`[{
// 		"name": "mock_step",
// 		"desc": "Mocked rule",
// 		"salience": 1,
// 		"when": "in.age >= 18",
// 		"then": [
// 			"Log('mock_step')",
// 			"out.Companies[0].GuaranteesMap['death'].IsSellable = true",
// 			"out.Companies[0].GuaranteesMap['permanent-disability'].IsSellable = true",
// 			"out.Companies[0].GuaranteesMap['temporary-disability'].IsSellable = true",
// 			"out.Companies[0].GuaranteesMap['serious-ill'].IsSellable = true",
// 			"Retract('mock_step')"
// 		]
// 	}]`)
// }

// func TestV2Life(t *testing.T) {
// 	birthDate := time.Now().UTC().AddDate(-18, 0, 0).Format(time.RFC3339)
// 	basePolicy := models.Policy{Contractor: models.Contractor{BirthDate: birthDate}}
// 	got, err := LifeMod(basePolicy, MockSellableProduct{})
// 	if err != nil {
// 		t.Fatal("got unexpected error")
// 	}
// 	expected := models.Product{Companies: []models.Company{{GuaranteesMap: map[string]*models.Guarante{
// 		"death":                {IsSellable: true},
// 		"permanent-disability": {IsSellable: true},
// 		"temporary-disability": {IsSellable: true},
// 		"serious-ill":          {IsSellable: true},
// 	}}}}

// 	if got == nil {
// 		t.Fatal("got nil product")
// 	}

// 	for key := range got.Companies[0].GuaranteesMap {
// 		assert.Equal(t, expected.Companies[0].GuaranteesMap[key], got.Companies[0].GuaranteesMap[key], "ooops")
// 	}
// }
