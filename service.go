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
