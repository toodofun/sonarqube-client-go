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
	"embed"
	"fmt"
	"path/filepath"
	"runtime"
	"text/template"
)

//go:embed tpl/client.tpl tpl/service.tpl
var embeddedTemplates embed.FS

func toolDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Dir(file)
}

func clientTemplatePath() string {
	if templateDir != nil && *templateDir != "" {
		return filepath.Join(toolDir(), *templateDir, "client.tpl")
	}
	return ""
}

func serviceTemplatePath() string {
	if templateDir != nil && *templateDir != "" {
		return filepath.Join(toolDir(), *templateDir, "service.tpl")
	}
	return ""
}

func parseClientTemplate() (*template.Template, error) {
	if path := clientTemplatePath(); path != "" {
		return template.New(clientTemplateName).Funcs(templateHelpers).ParseFiles(path)
	}
	return template.New(clientTemplateName).Funcs(templateHelpers).ParseFS(embeddedTemplates, "tpl/client.tpl")
}

func parseServiceTemplate() (*template.Template, error) {
	if path := serviceTemplatePath(); path != "" {
		return template.New(serviceTemplateName).Funcs(templateHelpers).ParseFiles(path)
	}
	return template.New(serviceTemplateName).Funcs(templateHelpers).ParseFS(embeddedTemplates, "tpl/service.tpl")
}

func mustParseClientTemplate() *template.Template {
	tpl, err := parseClientTemplate()
	if err != nil {
		panic(fmt.Errorf("parse client template: %w", err))
	}
	return tpl
}

func mustParseServiceTemplate() *template.Template {
	tpl, err := parseServiceTemplate()
	if err != nil {
		panic(fmt.Errorf("parse service template: %w", err))
	}
	return tpl
}
