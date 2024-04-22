package lib_test

import (
	"github.com/wopta/goworkspace/lib"
	"testing"
	"time"
)

type testAddMonths struct {
	input  time.Time
	output time.Time
}

func TestAddMonth(t *testing.T) {
	var inputs = []testAddMonths{
		{
			input:  time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC),
			output: time.Date(2024, time.January, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			input:  time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			output: time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			input:  time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			output: time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			input:  time.Date(2024, time.March, 31, 0, 0, 0, 0, time.UTC),
			output: time.Date(2024, time.April, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			input:  time.Date(2024, time.January, 31, 0, 0, 0, 0, time.UTC),
			output: time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			input:  time.Date(2023, time.January, 31, 0, 0, 0, 0, time.UTC),
			output: time.Date(2023, time.February, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			input:  time.Date(2023, time.February, 28, 0, 0, 0, 0, time.UTC),
			output: time.Date(2023, time.March, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			input:  time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC),
			output: time.Date(2024, time.March, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			input:  time.Date(2024, time.February, 28, 0, 0, 0, 0, time.UTC),
			output: time.Date(2024, time.March, 28, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, input := range inputs {
		calculatedDate := lib.AddMonths(input.input, 1)
		if calculatedDate != input.output {
			t.Fatalf("expected %v got %v", input.output, calculatedDate)
		}
	}
}

func TestAddTwoMonths(t *testing.T) {
	startDate := time.Date(2024, time.March, 31, 0, 0, 0, 0, time.UTC)
	expectedDate := time.Date(2024, time.May, 31, 0, 0, 0, 0, time.UTC)
	calculatedDate := lib.AddMonths(startDate, 2)

	if calculatedDate != expectedDate {
		t.Fatalf("expected %v got %v", expectedDate, calculatedDate)
	}

	startDate = time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)
	expectedDate = time.Date(2024, time.April, 29, 0, 0, 0, 0, time.UTC)
	calculatedDate = lib.AddMonths(startDate, 2)

	if calculatedDate != expectedDate {
		t.Fatalf("expected %v got %v", expectedDate, calculatedDate)
	}
}

func TestAddOneYear(t *testing.T) {
	startDate := time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)
	expectedDate := time.Date(2025, time.February, 28, 0, 0, 0, 0, time.UTC)
	calculatedDate := lib.AddMonths(startDate, 12)

	if calculatedDate != expectedDate {
		t.Fatalf("expected %v got %v", expectedDate, calculatedDate)
	}
}

func TestAddOneYearOneMonth(t *testing.T) {
	startDate := time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)
	expectedDate := time.Date(2025, time.March, 29, 0, 0, 0, 0, time.UTC)
	calculatedDate := lib.AddMonths(startDate, 13)

	if calculatedDate != expectedDate {
		t.Fatalf("expected %v got %v", expectedDate, calculatedDate)
	}
}
