package questions

import (
	"testing"
)

func TestQuestion2_IsFibNumber(t *testing.T) {
	testCases := []struct {
		val  int
		want bool
	}{
		{1, true},
		{2, true},
		{3, true},
		{4, false},
		{5, true},
		{6, false},
		{8, true},
	}
	for _, tc := range testCases {
		got := IsFibNumber(tc.val)
		if tc.want != got {
			t.Errorf("IsFibNumber(%d): want: %v, got: %v", tc.val, tc.want, got)
		}
	}
}
