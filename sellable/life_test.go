package sellable

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func getPolicyByContractorAge(age int) models.Policy {
	return models.Policy{
		Name:           models.LifeProduct,
		ProductVersion: models.ProductV2,
		Contractor: models.Contractor{
			BirthDate: time.Now().UTC().AddDate(-age, 0, 0).Format(time.RFC3339),
		},
	}
}

func TestLife(t *testing.T) {
	var (
		output  *models.Product
		err     error
		channel = models.ECommerceChannel
		inputs  []int
		outputs []map[string]*models.Guarante
	)

	os.Setenv("env", "local-test")

	inputFile := lib.GetFilesByEnv("data/test/sellable/input.json")
	err = json.Unmarshal(inputFile, &inputs)
	lib.CheckError(err)
	outputFile := lib.GetFilesByEnv("data/test/sellable/output.json")
	err = json.Unmarshal(outputFile, &outputs)
	lib.CheckError(err)

	for index, age := range inputs {
		policy := getPolicyByContractorAge(age)
		expected := outputs[index]
		output, err = Life(&policy, channel, nil, nil)
		if err != nil {
			t.Fatalf("error on sellable age %d: %s", age, err.Error())
		}
		if len(expected) == 0 {
			if len(output.Companies[0].GuaranteesMap) > 0 {
				t.Fatalf("age %d - expected %v - got %v", age, expected, output.Companies[0].GuaranteesMap)
			}
		} else {
			if len(output.Companies[0].GuaranteesMap) == 0 {
				t.Fatalf("age %d - expected %v - got %v", age, expected, output.Companies[0].GuaranteesMap)
			}

			for _, guarantee := range output.Companies[0].GuaranteesMap {
				if guarantee.IsSellable != expected[guarantee.Slug].IsSellable {
					t.Fatalf("age %d - guarantee %s - expected %v - got %v", age, guarantee.Slug, expected[guarantee.Slug].IsSellable, guarantee.IsSellable)
				}
				if guarantee.Config.DurationValues.Max != expected[guarantee.Slug].Config.DurationValues.Max {
					t.Fatalf("age %d - guarantee %s - expected %v - got %v", age, guarantee.Slug, expected[guarantee.Slug].Config.DurationValues.Max, guarantee.Config.DurationValues.Max)
				}
			}
		}
	}
}
