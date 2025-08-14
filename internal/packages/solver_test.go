package packages

import "testing"

func TestExamples(t *testing.T) {
	sizes := []int{250, 500, 1000, 2000, 5000}

	cases := []struct {
		items     int
		wantTotal int
	}{
		{1, 250},
		{250, 250},
		{251, 500},
		{501, 750},
		{12001, 12250},
	}

	for _, tc := range cases {
		got, err := Solve(tc.items, sizes)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if got.TotalItems != tc.wantTotal {
			t.Fatalf("items=%d got total=%d want=%d packs=%v", tc.items, got.TotalItems, tc.wantTotal, got.PacksUsed)
		}
	}
}
