package reserved

import "testing"

type testBMI struct {
	weight   int
	height   int
	expected bool
}

func TestCheckOutOfRangeBMI(t *testing.T) {
	var tests = []testBMI{
		{
			weight:   80,
			height:   180,
			expected: false, // bmi = 24,69135
		},
		{
			weight:   180,
			height:   80,
			expected: true, // bmi = 281.25
		},
		{
			weight:   40,
			height:   180,
			expected: true, // bmi = 12,34567
		},
		{
			weight:   50,
			height:   177,
			expected: true, // bmi = 15,95965
		},
		{
			weight:   51,
			height:   177,
			expected: false, // bmi = 16,27884
		},
		{
			weight:   112,
			height:   194,
			expected: false, // bmi = 29,75874
		},
		{
			weight:   112,
			height:   193,
			expected: true, // bmi = 30,06792
		},
	}

	for i, test := range tests {
		bmi, isOutOfRange := checkOutOfRangeBMI(test.weight, test.height)
		if isOutOfRange != test.expected {
			t.Errorf("BMI check failed for test n° %d: result is %t because bmi is %f", i+1, isOutOfRange, bmi)
		}
	}
}
