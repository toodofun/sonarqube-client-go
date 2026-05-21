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
	"fmt"
	"strings"
)

func assignActionTypeNames(def *apiDefinition) {
	usedRequest := make(map[string]struct{})
	usedOK := make(map[string]struct{})

	for _, ws := range def.WebServices {
		for _, action := range ws.Actions {
			if len(action.Params) > 0 {
				preferred := action.MethodName() + requestSuffix
				fallback := ws.Getter() + action.MethodName() + requestSuffix
				action.RequestType = uniqueTypeName(usedRequest, preferred, fallback)
			}
			if action.ResponseOKType == "" {
				continue
			}
			oldPrefix := action.ResponseOKType
			preferred := action.MethodName() + responseOKSuffix
			fallback := ws.Getter() + action.MethodName() + responseOKSuffix
			newPrefix := uniqueTypeName(usedOK, preferred, fallback)
			renameResponseTypes(action.ResponseTypes, oldPrefix, newPrefix)
			action.ResponseOKType = newPrefix
		}
	}
}

func uniqueTypeName(used map[string]struct{}, preferred, fallback string) string {
	if _, exists := used[preferred]; !exists {
		used[preferred] = struct{}{}
		return preferred
	}
	if _, exists := used[fallback]; !exists {
		used[fallback] = struct{}{}
		return fallback
	}
	for i := 2; ; i++ {
		name := fmt.Sprintf("%s%d", fallback, i)
		if _, exists := used[name]; !exists {
			used[name] = struct{}{}
			return name
		}
	}
}

func renameResponseTypes(types []responseGoType, oldPrefix, newPrefix string) {
	for i := range types {
		types[i].Name = strings.Replace(types[i].Name, oldPrefix, newPrefix, 1)
		for j := range types[i].Fields {
			types[i].Fields[j].Type = strings.ReplaceAll(types[i].Fields[j].Type, oldPrefix, newPrefix)
		}
	}
}
