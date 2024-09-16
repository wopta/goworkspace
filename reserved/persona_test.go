package reserved

import "testing"

type testBMI struct {
	description string
	weight      int
	height      int
	expected    bool
}

func TestCheckOutOfRangeBMI(t *testing.T) {
	var tests = []testBMI{
		{
			description: "A value within the limits",
			weight:      80,
			height:      180,
			expected:    false, // bmi = 24,69135
		},
		{
			description: "Far above bmiUpperLimit",
			weight:      180,
			height:      80,
			expected:    true, // bmi = 281.25
		},
		{
			description: "Far below bmiLowerLimit",
			weight:      40,
			height:      180,
			expected:    true, // bmi = 12,34567
		},
		{
			description: "Slightly below bmiLowerLimit",
			weight:      50,
			height:      177,
			expected:    true, // bmi = 15,95965
		},
		{
			description: "Slightly above bmiLowerLimit",
			weight:      51,
			height:      177,
			expected:    false, // bmi = 16,27884
		},
		{
			description: "Slightly below bmiUpperLimit",
			weight:      112,
			height:      168,
			expected:    false, // bmi = 39,68253
		},
		{
			description: "Slightly above bmiUpperLimit",
			weight:      112,
			height:      167,
			expected:    true, // bmi = 40,15920
		},
	}

	for i, test := range tests {
		_, got := checkOutOfRangeBMI(test.weight, test.height)
		if test.expected != got {
			t.Errorf("error in test %s %d: expected %t, got %t", test.description, i+1, test.expected, got)
		}
	}
}
