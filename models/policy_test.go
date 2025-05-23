package models_test

import (
	"testing"
	"time"

	"gitlab.dev.wopta.it/goworkspace/models"
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
	inputs := []calculateContractorAgeInput{
		{
			birthDate: dateInput{day: 21, month: 3, year: 1969},
			startDate: dateInput{day: 20, month: 3, year: 2024},
		},
		{
			birthDate: dateInput{day: 21, month: 3, year: 1969},
			startDate: dateInput{day: 21, month: 3, year: 2024},
		},
		{
			birthDate: dateInput{day: 21, month: 3, year: 1969},
			startDate: dateInput{day: 5, month: 1, year: 2024},
		},
		{
			birthDate: dateInput{day: 21, month: 1, year: 1969},
			startDate: dateInput{day: 5, month: 3, year: 2024},
		},
	}

	output := []int{54, 55, 54, 55}

	for index, in := range inputs {
		calculateContractorAge(t, in, output[index])
	}
}

func TestCalculateAgeStartDateLeapBirthDateLeap(t *testing.T) {
	inputs := []calculateContractorAgeInput{
		{
			birthDate: dateInput{day: 21, month: 3, year: 1980},
			startDate: dateInput{day: 20, month: 3, year: 2024},
		},
		{
			birthDate: dateInput{day: 21, month: 3, year: 1980},
			startDate: dateInput{day: 21, month: 3, year: 2024},
		},
		{
			birthDate: dateInput{day: 21, month: 3, year: 1980},
			startDate: dateInput{day: 5, month: 1, year: 2024},
		},
		{
			birthDate: dateInput{day: 21, month: 1, year: 1980},
			startDate: dateInput{day: 5, month: 3, year: 2024},
		},
	}

	output := []int{43, 44, 43, 44}

	for index, in := range inputs {
		calculateContractorAge(t, in, output[index])
	}
}

func TestCalculateAgeStartDateNonLeapBirthDateLeap(t *testing.T) {
	inputs := []calculateContractorAgeInput{
		{
			birthDate: dateInput{day: 21, month: 3, year: 1980},
			startDate: dateInput{day: 20, month: 3, year: 2023},
		},
		{
			birthDate: dateInput{day: 21, month: 3, year: 1980},
			startDate: dateInput{day: 21, month: 3, year: 2023},
		},
		{
			birthDate: dateInput{day: 21, month: 3, year: 1980},
			startDate: dateInput{day: 5, month: 1, year: 2023},
		},
		{
			birthDate: dateInput{day: 21, month: 1, year: 1980},
			startDate: dateInput{day: 5, month: 3, year: 2023},
		},
	}

	output := []int{42, 43, 42, 43}

	for index, in := range inputs {
		calculateContractorAge(t, in, output[index])
	}
}

func TestCalculateAgeStartDateNonLeapBirthDateNonLeap(t *testing.T) {
	inputs := []calculateContractorAgeInput{
		{
			birthDate: dateInput{day: 7, month: 10, year: 1994},
			startDate: dateInput{day: 20, month: 3, year: 2023},
		},
		{
			birthDate: dateInput{day: 27, month: 9, year: 1998},
			startDate: dateInput{day: 15, month: 6, year: 2023},
		},
		{
			birthDate: dateInput{day: 12, month: 3, year: 1987},
			startDate: dateInput{day: 10, month: 4, year: 2023},
		},
		{
			birthDate: dateInput{day: 14, month: 2, year: 2022},
			startDate: dateInput{day: 10, month: 4, year: 2023},
		},
	}

	output := []int{28, 24, 36, 1}

	for index, in := range inputs {
		calculateContractorAge(t, in, output[index])
	}
}
