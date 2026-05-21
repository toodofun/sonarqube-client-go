package main

import (
	"strings"
	"testing"
)

func Test_buildResponseSchema_listAzureProjects(t *testing.T) {
	example := `{
  "projects": [
    {
      "name": "project1",
      "description": "description1"
    }
  ]
}`
	types, err := buildResponseSchema("AlmIntegrationsServiceListAzureProjectsOK", example)
	if err != nil {
		t.Fatal(err)
	}
	if len(types) != 2 {
		t.Fatalf("types count = %d, want 2", len(types))
	}
	var root *responseGoType
	for i := range types {
		if types[i].Name == "AlmIntegrationsServiceListAzureProjectsOK" {
			root = &types[i]
			break
		}
	}
	if root == nil {
		t.Fatal("missing root type")
	}
	var projectsField *responseGoField
	for i := range root.Fields {
		if root.Fields[i].JSONKey == "projects" {
			projectsField = &root.Fields[i]
			if projectsField.Name != "Projects" {
				t.Fatalf("field name = %s, want Projects", projectsField.Name)
			}
			break
		}
	}
	if projectsField == nil {
		t.Fatal("missing projects field")
	}
	wantSlice := "[]*AlmIntegrationsServiceListAzureProjectsOKProject"
	if projectsField.Type != wantSlice {
		t.Fatalf("projects type = %s, want %s", projectsField.Type, wantSlice)
	}
	var item *responseGoType
	for i := range types {
		if strings.HasSuffix(types[i].Name, "Project") {
			item = &types[i]
			break
		}
	}
	if item == nil {
		t.Fatal("missing project item type")
	}
}

func Test_pointerizeStructFields(t *testing.T) {
	types := []responseGoType{
		{
			Name: "ValidateOK",
			Fields: []responseGoField{
				{Name: "Errors", JSONKey: "errors", Type: "[]ValidateOKError"},
			},
		},
		{Name: "ValidateOKError"},
	}
	pointerizeStructFields(types)
	if types[0].Fields[0].Type != "[]*ValidateOKError" {
		t.Fatalf("slice type = %s", types[0].Fields[0].Type)
	}
	types[1].Fields = []responseGoField{{Name: "Nested", Type: "ValidateOKError"}}
	pointerizeStructFields(types)
	if types[1].Fields[0].Type != "*ValidateOKError" {
		t.Fatalf("struct type = %s", types[1].Fields[0].Type)
	}
}
