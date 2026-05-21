package main

import "testing"

func Test_jsonKeyToIdent(t *testing.T) {
	cases := map[string]string{
		"1":                                "Key1",
		"squid:MethodCyclomaticComplexity": "SquidMethodCyclomaticComplexity",
		"Non Heap Init (MB)":               "NonHeapInitMB",
		"Web JVM State":                    "WebJVMState",
		"":                                 "Field",
	}
	for in, want := range cases {
		if got := jsonKeyToIdent(in); got != want {
			t.Errorf("jsonKeyToIdent(%q) = %q, want %q", in, got, want)
		}
	}
}
