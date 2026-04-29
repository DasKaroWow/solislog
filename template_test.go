package solislog

import "testing"

func TestParseTemplateParsesTextAndFields(t *testing.T) {
	parts := parseTemplate("{time} | {level} | {message}\n")

	want := []templatePart{
		{mode: fieldMode, value: "time"},
		{mode: textMode, value: " | "},
		{mode: fieldMode, value: "level"},
		{mode: textMode, value: " | "},
		{mode: fieldMode, value: "message"},
		{mode: textMode, value: "\n"},
	}

	if len(parts) != len(want) {
		t.Fatalf("len(parts) = %d, want %d", len(parts), len(want))
	}

	for i := range want {
		if parts[i] != want[i] {
			t.Fatalf("parts[%d] = %+v, want %+v", i, parts[i], want[i])
		}
	}
}

func TestParseTemplateParsesExtraField(t *testing.T) {
	parts := parseTemplate("{extra[source]}")

	want := []templatePart{
		{mode: extraMode, value: "source"},
	}

	if len(parts) != len(want) {
		t.Fatalf("len(parts) = %d, want %d", len(parts), len(want))
	}

	if parts[0] != want[0] {
		t.Fatalf("parts[0] = %+v, want %+v", parts[0], want[0])
	}
}

func TestParseTemplateParsesFullExtraField(t *testing.T) {
	parts := parseTemplate("{extra}")

	want := []templatePart{
		{mode: fieldMode, value: "extra"},
	}

	if len(parts) != len(want) {
		t.Fatalf("len(parts) = %d, want %d", len(parts), len(want))
	}

	if parts[0] != want[0] {
		t.Fatalf("parts[0] = %+v, want %+v", parts[0], want[0])
	}
}

func TestParseTemplatePanicsOnUnknownField(t *testing.T) {
	assertPanics(t, func() {
		parseTemplate("{unknown}")
	})
}

func TestParseTemplatePanicsOnEmptyPlaceholder(t *testing.T) {
	assertPanics(t, func() {
		parseTemplate("{}")
	})
}

func TestParseTemplatePanicsOnEmptyExtraKey(t *testing.T) {
	assertPanics(t, func() {
		parseTemplate("{extra[]}")
	})
}

func TestParseTemplatePanicsOnUnclosedPlaceholder(t *testing.T) {
	assertPanics(t, func() {
		parseTemplate("{time")
	})
}

func TestParseTemplatePanicsOnUnexpectedClosingBrace(t *testing.T) {
	assertPanics(t, func() {
		parseTemplate("time}")
	})
}

func assertPanics(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()

	fn()
}
