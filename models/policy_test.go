package models_test

import (
	"github.com/wopta/goworkspace/models"
	"testing"
	"time"
)

type dateInput struct {
	day   int
	month int
	year  int
}

type calculateContractorAgeInput struct {
	birthDate dateInput
	startDate dateInput
}

func getPolicy(birthDate, startDate dateInput) models.Policy {
	formattedStartDate := time.Date(startDate.year, time.Month(startDate.month), startDate.day, 0, 0, 0, 0, time.Local)
	formattedBirthDate := time.Date(birthDate.year, time.Month(birthDate.month), birthDate.day, 0, 0, 0, 0, time.Local).Format(time.RFC3339)
	return models.Policy{
		StartDate: formattedStartDate,
		Contractor: models.Contractor{
			BirthDate: formattedBirthDate,
		},
	}
}

func calculateContractorAge(t *testing.T, in calculateContractorAgeInput, expectedAge int) {
	policy := getPolicy(in.birthDate, in.startDate)
	contractorAge, _ := policy.CalculateContractorAge()
	if contractorAge != expectedAge {
		t.Fatalf("input: %v contractorAge %02d - expectedAge: %02d", in, contractorAge, expectedAge)
	}
}

func TestCalculateAgeStartDateLeapBirthDateNonLeap(t *testing.T) {
	// START DATE LEAP/BIRTH DATE NON LEAP
	inputs := []calculateContractorAgeInput{
		{birthDate: dateInput{
			day:   21,
			month: 3,
			year:  1969,
		}, startDate: dateInput{
			day:   20,
			month: 3,
			year:  2024,
		}},
		{birthDate: dateInput{
			day:   21,
			month: 3,
			year:  1969,
		}, startDate: dateInput{
			day:   21,
			month: 3,
			year:  2024,
		}},
		{birthDate: dateInput{
			day:   21,
			month: 3,
			year:  1969,
		}, startDate: dateInput{
			day:   5,
			month: 1,
			year:  2024,
		}},
		{birthDate: dateInput{
			day:   21,
			month: 1,
			year:  1969,
		}, startDate: dateInput{
			day:   5,
			month: 3,
			year:  2024,
		}},
	}

	/*inputs := [][]int{
		// START DATE LEAP/BIRTH DATE LEAP
		{21, 3, 1980, 20, 3, 2024},
		{21, 3, 1980, 21, 3, 2024},
		{21, 3, 1980, 05, 1, 2024},
		{21, 1, 1980, 05, 1, 2024},
		// START DATE NON LEAP/BIRTH DATE NON LEAP
		{07, 10, 1994, 20, 3, 2023},
		{27, 9, 1998, 15, 6, 2023},
		{12, 3, 1987, 10, 4, 2023},
	}*/

	output := []int{54, 55, 54, 55} //43, 44, 43, 43, 28, 24, 36}

	for index, in := range inputs {
		calculateContractorAge(t, in, output[index])
	}
}

func TestCalculateAgeStartDateLeapBirthDateLeap(t *testing.T) {
	// START DATE LEAP/BIRTH DATE NON LEAP
	inputs := []calculateContractorAgeInput{
		{birthDate: dateInput{
			day:   21,
			month: 3,
			year:  1980,
		}, startDate: dateInput{
			day:   20,
			month: 3,
			year:  2024,
		}},
		{birthDate: dateInput{
			day:   21,
			month: 3,
			year:  1980,
		}, startDate: dateInput{
			day:   21,
			month: 3,
			year:  2024,
		}},
		{birthDate: dateInput{
			day:   21,
			month: 3,
			year:  1980,
		}, startDate: dateInput{
			day:   5,
			month: 1,
			year:  2024,
		}},
		{birthDate: dateInput{
			day:   21,
			month: 1,
			year:  1980,
		}, startDate: dateInput{
			day:   5,
			month: 3,
			year:  2024,
		}},
	}

	output := []int{43, 44, 43, 44}

	for index, in := range inputs {
		calculateContractorAge(t, in, output[index])
	}
}

func TestCalculateAgeStartDateNonLeapBirthDateNonLeap(t *testing.T) {
	// START DATE NON LEAP/BIRTH DATE NON LEAP
	inputs := []calculateContractorAgeInput{
		{birthDate: dateInput{
			day:   7,
			month: 10,
			year:  1994,
		}, startDate: dateInput{
			day:   20,
			month: 3,
			year:  2023,
		}},
		{birthDate: dateInput{
			day:   27,
			month: 9,
			year:  1998,
		}, startDate: dateInput{
			day:   15,
			month: 6,
			year:  2023,
		}},
		{birthDate: dateInput{
			day:   12,
			month: 3,
			year:  1987,
		}, startDate: dateInput{
			day:   10,
			month: 4,
			year:  2023,
		}},
	}

	output := []int{28, 24, 36}

	for index, in := range inputs {
		calculateContractorAge(t, in, output[index])
	}
}
