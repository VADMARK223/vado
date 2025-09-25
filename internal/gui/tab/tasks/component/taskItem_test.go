package component

import "testing"

func TestGetDescText(t *testing.T) {
	got := GetDescText("vadmark")
	want := " (vadmark)"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
