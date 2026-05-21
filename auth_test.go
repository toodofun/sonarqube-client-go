// Copyright 2025 The Toodofun Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http:www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
