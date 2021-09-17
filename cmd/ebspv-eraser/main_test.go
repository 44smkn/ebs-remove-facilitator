package main

import "testing"

func TestFormatVerion(t *testing.T) {
	expects := "ebspv-eraser version 0.1.0 (2021-09-17)\n"
	if got := FormatVersion("0.1.0", "2021-09-17"); got != expects {
		t.Errorf("FormatVersion() = %q, wants %q", got, expects)
	}
}
