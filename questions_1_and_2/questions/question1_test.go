package questions

import (
	"testing"
)

func TestQuestion1_Parse(t *testing.T) {
	testCases := []struct {
		id          string
		wantError   error
		wantDetails *SAIDDetails
	}{
		{"910310501708", ErrorSize, nil},
		{"P103105017084", ErrorYear, nil},
		{"910A105017084", ErrorMonth, nil},
		{"91031R5017084", ErrorDay, nil},
		{"9103105X17084", ErrorGender, nil},
		{"9103105017284", ErrorCitizenship, nil},
		{"910310501708N", ErrorCheckDigit, nil},
		{"9103105017084", nil, &SAIDDetails{
			DateOfBirth: struct {
				Year  int
				Month int
				Day   int
			}{
				Day:   10,
				Month: 3,
				Year:  91,
			},
			Gender:      "MALE",
			CitizenShip: true,
		}},
	}

	for _, tc := range testCases {
		details, err := Parse(tc.id)
		if tc.wantError != err {
			t.Errorf("Parse(%q): want: %q, got: %q", tc.id, tc.wantError, err)
		}
		if err == nil {
			if details == nil {
				t.Errorf("Parse(%q): no details returned", tc.id)
				continue
			}
			switch {
			case details.DateOfBirth.Day != tc.wantDetails.DateOfBirth.Day:
				t.Errorf("Parse(%q): wanted: %v got: %v", tc.id, tc.wantDetails, details)
			case details.DateOfBirth.Month != tc.wantDetails.DateOfBirth.Month:
				t.Errorf("Parse(%q): wanted: %v got: %v", tc.id, tc.wantDetails, details)
			case details.DateOfBirth.Year != tc.wantDetails.DateOfBirth.Year:
				t.Errorf("Parse(%q): wanted: %v got: %v", tc.id, tc.wantDetails, details)
			case details.Gender != tc.wantDetails.Gender:
				t.Errorf("Parse(%q): wanted: %v got: %v", tc.id, tc.wantDetails, details)
			case details.CitizenShip != tc.wantDetails.CitizenShip:
				t.Errorf("Parse(%q): wanted: %v got: %v", tc.id, tc.wantDetails, details)
			}
		}
	}
}

func TestQuestionq_Luhn(t *testing.T) {
	testCases := []struct {
		id   string
		want bool
	}{
		{"9103105017084", true},
		{"9103105017085", false},
	}

	for _, tc := range testCases {
		got := luhnCheck(tc.id)
		if tc.want != got {
			t.Errorf("Parse(%q): want: %v, got: %v", tc.id, tc.want, got)
		}
	}
}
