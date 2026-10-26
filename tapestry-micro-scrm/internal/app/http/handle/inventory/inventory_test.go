package inventory

import "testing"

func TestFirstStoreCode(t *testing.T) {
	if got := firstStoreCode(nil); got != "" {
		t.Fatalf("empty store list must not panic or return a code, got %q", got)
	}
	if got := firstStoreCode([]string{"S001"}); got != "S001" {
		t.Fatalf("expected S001, got %q", got)
	}
}
