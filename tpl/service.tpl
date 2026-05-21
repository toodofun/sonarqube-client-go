package {{.PackageName}}

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// {{.ServiceName}} {{.Description | formatDescription }}
{{- if .Since }}
// Since : {{.Since}}
{{- end}}
{{- if .Deprecated }}
// Deprecated
{{- end}}
{{- if .Internal }}
// Internal
{{- end}}
type {{.ServiceName}} struct {
	client *Client
	url string
}


{{- if .Deprecated }}
// Deprecated
{{- end}}
func new{{.ServiceName}} (client *Client) *{{.ServiceName}}{
	s := &{{.ServiceName}}{
		client: client,
		url: "{{.Path}}",
	}
	return s
}

{{- range $index, $element := .Actions}}
{{ template "action" $element}}
{{- end}}

{{- define "action"}}
// {{ .MethodName }} {{.Description | formatDescription }}
{{- if .Since}}
// Since {{.Since}}
{{- end}}
{{- if len .Changelog }}
//
// Changelog:
	{{- range .Changelog }}
// {{.String}}
	{{- end}}
{{- end}}
{{- if .Deprecated }}
//
// Deprecated since {{.DeprecatedSince}}
{{- end}}
{{- if .ResponseOKType}}
func (s *{{.ServiceName}}) {{.MethodName}} (ctx context.Context{{- if .Params}}, request *{{.RequestType}}{{- end}}) (*{{.ResponseOKType}}, *http.Response, error) {
	resp, err := s.client.invoke(ctx, {{.Post}}, fmt.Sprintf("%s/%s", s.url, "{{.Key}}"), {{- if .Params}} request {{- else}} nil {{- end}})
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to call {{.ServiceName}}.{{.MethodName}}")
	}
	out := &{{.ResponseOKType}}{}
	if err = decodeBody(resp, out); err != nil {
		return nil, resp, errors.Wrap(err, "failed to decode {{.ServiceName}}.{{.MethodName}} response")
	}
	return out, resp, nil
}
{{- else}}
func (s *{{.ServiceName}}) {{.MethodName}} (ctx context.Context{{- if .Params}}, request *{{.RequestType}}{{- end}}) (*http.Response, error) {
	resp, err := s.client.invoke(ctx, {{.Post}}, fmt.Sprintf("%s/%s", s.url, "{{.Key}}"), {{- if .Params}} request {{- else}} nil {{- end}})
	if err != nil {
		return nil, errors.Wrap(err, "failed to call {{.ServiceName}}.{{.MethodName}}")
	}
	return resp, nil
}
{{- end}}
{{- if .Params }}
{{ template "request" .}}
{{- end}}

{{ template "response" .}}

{{- end}}

{{- define "request"}}
{{$post := .Post}}
type {{.RequestType}} struct {
{{- /* see https://github.com/golang/go/issues/18221#issuecomment-394255883 */}}
{{- range .Params}}
	// {{.Description | formatDescription }}
	{{- if .Since | formatSince }}
	// Since {{ .Since | formatSince }}
	{{- end}}
	{{- if .Required}}
	// Required
	{{- end }}
	{{- if .Internal}}
	// Internal
	{{- end }}
	{{- if .DefaultValue }}
	// Default: {{.DefaultValue}}
	{{- end}}
	{{- if .ExampleValue}}
	// Example: {{.ExampleValue}}
	{{- end }}
	{{- if .PossibleValues }}
	// Possible values: {{- range .PossibleValues}} "{{.}}", {{- end}}
	{{- end}}
	{{- if .Deprecated}}
	// Deprecated since {{.DeprecatedSince.String}}
	{{- end }}
	{{.ParamName}} *string {{- if $post }} {{tick}}json:"{{.Key}}{{ if not .Required}},omitempty{{ end }}"{{tick}} {{- else}} {{tick}}url:"{{.Key}}{{ if not .Required}},omitempty{{ end }}"{{tick}} {{- end}}
{{- end}}
}
{{- end}}

{{- define "response"}}
{{- range .ResponseTypes}}
// {{.Name}} response body (from SonarQube response example)
type {{.Name}} struct {
{{- range .Fields}}
	{{.Name}} {{.Type}} {{tick}}json:"{{.JSONKey}}"{{tick}}
{{- end}}
}
{{- end}}
{{- end}}
