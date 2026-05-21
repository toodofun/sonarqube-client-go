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
	"bytes"
	"fmt"
	"go/format"
	"io"
	"log"
)

const (
	serviceTemplateName = "service.tpl"
)

func renderService(in io.Writer, data *webService) error {
	buff := bytes.NewBuffer([]byte{})

	serviceTemplate := mustParseServiceTemplate()
	if err := serviceTemplate.Execute(buff, data); err != nil {
		return fmt.Errorf("failed to render service %s：%w", data.ServiceName(), err)
	}

	src := buff.Bytes()

	formatted, err := format.Source(src)
	if err != nil {
		log.Printf("failed to format source of %s, %s", data.ServiceName(), err.Error())
		formatted = src
	}

	_, err = in.Write(formatted)

	return err
}
