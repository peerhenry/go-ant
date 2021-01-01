package ant

import "testing"

func TestSum(t *testing.T) {
	total := 5 + 5
	if total != 10 {
		t.Errorf("Expected 10 but got %d", total)
	}
}
