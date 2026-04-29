package solislog

import "testing"

func TestCloneExtraWithNil(t *testing.T) {
	cloned := cloneExtra(nil)

	if cloned == nil {
		t.Fatal("cloneExtra(nil) returned nil, want empty Extra")
	}

	if len(cloned) != 0 {
		t.Fatalf("cloneExtra(nil) length = %d, want 0", len(cloned))
	}
}

func TestCloneExtraCopiesValues(t *testing.T) {
	src := Extra{
		"source": "telegram",
		"id":     "123",
	}

	cloned := cloneExtra(src)

	if cloned["source"] != "telegram" {
		t.Fatalf("cloned source = %q, want %q", cloned["source"], "telegram")
	}

	if cloned["id"] != "123" {
		t.Fatalf("cloned id = %q, want %q", cloned["id"], "123")
	}
}

func TestCloneExtraDoesNotShareMap(t *testing.T) {
	src := Extra{
		"id": "123",
	}

	cloned := cloneExtra(src)

	src["id"] = "456"

	if cloned["id"] != "123" {
		t.Fatalf("cloned id changed to %q, want %q", cloned["id"], "123")
	}
}

func TestMergeExtraWithNilBase(t *testing.T) {
	merged := mergeExtra(nil, Extra{
		"id": "123",
	})

	if merged["id"] != "123" {
		t.Fatalf("merged id = %q, want %q", merged["id"], "123")
	}
}

func TestMergeExtraWithNilOverride(t *testing.T) {
	merged := mergeExtra(Extra{
		"source": "telegram",
	}, nil)

	if merged["source"] != "telegram" {
		t.Fatalf("merged source = %q, want %q", merged["source"], "telegram")
	}
}

func TestMergeExtraOverrideWins(t *testing.T) {
	base := Extra{
		"source": "telegram",
		"id":     "-1",
	}

	override := Extra{
		"id": "123",
	}

	merged := mergeExtra(base, override)

	if merged["source"] != "telegram" {
		t.Fatalf("merged source = %q, want %q", merged["source"], "telegram")
	}

	if merged["id"] != "123" {
		t.Fatalf("merged id = %q, want %q", merged["id"], "123")
	}
}

func TestMergeExtraDoesNotMutateBase(t *testing.T) {
	base := Extra{
		"id": "123",
	}

	override := Extra{
		"id": "456",
	}

	_ = mergeExtra(base, override)

	if base["id"] != "123" {
		t.Fatalf("base id changed to %q, want %q", base["id"], "123")
	}
}

func TestMergeExtraDoesNotShareMap(t *testing.T) {
	base := Extra{
		"id": "123",
	}

	merged := mergeExtra(base, nil)

	base["id"] = "456"

	if merged["id"] != "123" {
		t.Fatalf("merged id changed to %q, want %q", merged["id"], "123")
	}
}
