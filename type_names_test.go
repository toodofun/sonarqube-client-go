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

import "testing"

func Test_assignActionTypeNames(t *testing.T) {
	def := &apiDefinition{
		WebServices: []*webService{
			{
				Path: "api/alm_integrations",
				Actions: []*action{
					{
						Key:            "list_azure_projects",
						Params:         []*param{{Key: "almSetting"}},
						ResponseOKType: "AlmIntegrationsServiceListAzureProjectsOK",
						ResponseTypes: []responseGoType{
							{Name: "AlmIntegrationsServiceListAzureProjectsOK"},
						},
					},
				},
			},
			{
				Path: "api/webhooks",
				Actions: []*action{
					{Key: "list", Params: []*param{{Key: "project"}}},
				},
			},
			{
				Path: "api/project_branches",
				Actions: []*action{
					{Key: "list", Params: []*param{{Key: "project"}}},
				},
			},
		},
	}
	for _, ws := range def.WebServices {
		ws.Getter()
		for _, a := range ws.Actions {
			a.MethodName()
		}
	}
	assignActionTypeNames(def)

	a0 := def.WebServices[0].Actions[0]
	if a0.RequestType != "ListAzureProjectsRequest" {
		t.Fatalf("request type = %s", a0.RequestType)
	}
	if a0.ResponseOKType != "ListAzureProjectsOK" {
		t.Fatalf("ok type = %s", a0.ResponseOKType)
	}

	a1 := def.WebServices[1].Actions[0]
	a2 := def.WebServices[2].Actions[0]
	if a1.RequestType != "ListRequest" {
		t.Fatalf("webhooks request = %s", a1.RequestType)
	}
	if a2.RequestType != "ProjectBranchesListRequest" {
		t.Fatalf("branches request = %s", a2.RequestType)
	}
}
