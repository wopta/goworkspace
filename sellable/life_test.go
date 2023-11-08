package sellable

import (
	"os"
	"testing"
	"time"

	"github.com/wopta/goworkspace/models"
)

func getPolicyByContractorAge(age int) models.Policy {
	return models.Policy{
		Name:           models.LifeProduct,
		ProductVersion: models.ProductV2,
		Contractor: models.User{
			BirthDate: time.Now().UTC().AddDate(-age, 0, 0).Format(time.RFC3339),
		},
	}
}

func TestLife(t *testing.T) {
	var (
		output *models.Product
		err    error
	)
	channel := models.ECommerceChannel
	os.Setenv("env", "test")

	inputs := []int{17, 18, 53, 54, 55, 58, 59, 60, 63, 64, 65, 68, 69, 70, 71}
	outputs := []map[string]*models.Guarante{
		{},
		{
			"death":                {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 20}}},
			"permanent-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 20}}},
			"temporary-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
			"serious-ill":          {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
		},
		{
			"death":                {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 20}}},
			"permanent-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 20}}},
			"temporary-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
			"serious-ill":          {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
		},
		{
			"death":                {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 20}}},
			"permanent-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 20}}},
			"temporary-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
			"serious-ill":          {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
		},
		{
			"death":                {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 15}}},
			"permanent-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 15}}},
			"temporary-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
			"serious-ill":          {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 5}}},
		},
		{
			"death":                {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 15}}},
			"permanent-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 15}}},
			"temporary-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
			"serious-ill":          {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 5}}},
		},
		{
			"death":                {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 15}}},
			"permanent-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 15}}},
			"temporary-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
			"serious-ill":          {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 5}}},
		},
		{
			"death":                {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
			"permanent-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
			"temporary-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
			"serious-ill":          {IsSellable: false, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
		},
		{
			"death":                {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
			"permanent-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
			"temporary-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
			"serious-ill":          {IsSellable: false, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
		},
		{
			"death":                {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
			"permanent-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
			"temporary-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
			"serious-ill":          {IsSellable: false, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
		},
		{
			"death":                {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 5}}},
			"permanent-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 5}}},
			"temporary-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 5}}},
			"serious-ill":          {IsSellable: false, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
		},
		{
			"death":                {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 5}}},
			"permanent-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 5}}},
			"temporary-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 5}}},
			"serious-ill":          {IsSellable: false, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
		},
		{
			"death":                {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 5}}},
			"permanent-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 5}}},
			"temporary-disability": {IsSellable: true, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 5}}},
			"serious-ill":          {IsSellable: false, Config: &models.GuaranteValue{DurationValues: &models.DurationFieldValue{Max: 10}}},
		},
		{},
		{},
	}

	for index, age := range inputs {
		policy := getPolicyByContractorAge(age)
		expected := outputs[index]
		output, err = Life(&policy, channel, nil, nil)
		if err != nil {
			t.Fatalf("error on sellable age %d", age)
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
