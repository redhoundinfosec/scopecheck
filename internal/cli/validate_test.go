package cli

import "testing"

func TestDedupeTargetsPreservesOrder(t *testing.T) {
	in := []string{"b", "a", "b", "c", "a"}
	out := dedupeTargets(in)

	want := []string{"b", "a", "c"}
	if len(out) != len(want) {
		t.Fatalf("unexpected length: got %d want %d", len(out), len(want))
	}
	for i := range want {
		if out[i] != want[i] {
			t.Fatalf("unexpected output[%d]: got %q want %q", i, out[i], want[i])
		}
	}
}

func TestValidateTargetCount(t *testing.T) {
	var targets []string
	for i := 0; i < 10000; i++ {
		targets = append(targets, "t")
	}
	if err := validateTargetCount(targets); err != nil {
		t.Fatalf("unexpected error at max limit: %v", err)
	}
	targets = append(targets, "overflow")
	if err := validateTargetCount(targets); err == nil {
		t.Fatalf("expected error when exceeding target cap")
	}
}
