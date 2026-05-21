package main

import (
	"encoding/base64"
	"strings"
	"testing"
)

func Test_resolveAuthorization(t *testing.T) {
	token := "squ_test"
	want := "Basic " + base64.StdEncoding.EncodeToString([]byte(token+":"))

	if got := resolveAuthorization("", token); got != want {
		t.Fatalf("token auth = %q, want %q", got, want)
	}
	if got := resolveAuthorization("Bearer x", token); got != "Bearer x" {
		t.Fatalf("explicit auth = %q", got)
	}
	if got := resolveAuthorization("", ""); got != "" {
		t.Fatalf("empty = %q", got)
	}
	if !strings.HasPrefix(resolveAuthorization("  Basic abc  ", ""), "Basic abc") {
		t.Fatal("trim auth")
	}
}
