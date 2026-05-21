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
	"regexp"
	"strings"
	"text/template"
)

var snakeToCamelRE = regexp.MustCompile("_([a-z])")

func snakeToCamel(str string) string {
	return snakeToCamelRE.ReplaceAllStringFunc(str, func(match string) string {
		return strings.ToUpper(strings.TrimPrefix(match, "_"))
	})
}

func makeExported(str string) string {
	first := string(str[0])
	return strings.Replace(str, first, strings.ToUpper(first), 1)
}

func makeUnexported(str string) string {
	first := string(str[0])
	return strings.Replace(str, first, strings.ToLower(first), 1)
}

func sanitizeItentifier(str string) string {
	return strings.ReplaceAll(str, "-", "_")
}

func formatFieldName(str string) string {
	return strings.ReplaceAll(str, ".", "_")
}

func tick() string {
	return "`"
}

func formatSince(since version) string {
	return since.String()
}

func replaceTags(str string) string {
	repl := strings.NewReplacer(
		"\n", "\n// ",
		"<br> ", "\n// ",
		"<br>", "\n// ",
		"<br/>", "\n// ",
		"<br />", "\n// ",
		"<ul>", "\n// ",
		"<li>", " * ",
		"</li>", "\n// ",
		"</ul>", "",
	)
	return repl.Replace(str)
}

var templateHelpers = template.FuncMap{
	"formatDescription": replaceTags,
	"tick":              tick,
	"formatSince":       formatSince,
}
