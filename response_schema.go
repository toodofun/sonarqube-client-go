package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// responseGoType 由响应 JSON 示例推导出的 Go 结构体定义。
type responseGoType struct {
	Name   string
	Fields []responseGoField
}

type responseGoField struct {
	Name    string
	JSONKey string
	Type    string
}

type schemaBuilder struct {
	types []responseGoType
	seen  map[string]struct{}
}

const (
	goTypeInterface      = "interface{}"
	goTypeSliceInterface = "[]interface{}"
)

func buildResponseSchema(rootTypeName, exampleJSON string) ([]responseGoType, error) {
	exampleJSON = strings.TrimSpace(exampleJSON)
	if exampleJSON == "" || !looksLikeJSONObject(exampleJSON) {
		return nil, nil
	}
	var sample interface{}
	if err := json.Unmarshal([]byte(exampleJSON), &sample); err != nil {
		return nil, fmt.Errorf("parse response example json: %w", err)
	}
	root, ok := sample.(map[string]interface{})
	if !ok {
		return nil, nil
	}
	b := &schemaBuilder{seen: make(map[string]struct{})}
	b.infer(rootTypeName, root)
	pointerizeStructFields(b.types)
	return b.types, nil
}

// pointerizeStructFields 将字段中的结构体类型改为指针（含切片元素）。
func pointerizeStructFields(types []responseGoType) {
	names := make(map[string]struct{}, len(types))
	for _, t := range types {
		names[t.Name] = struct{}{}
	}
	for i := range types {
		for j := range types[i].Fields {
			types[i].Fields[j].Type = pointerizeFieldType(types[i].Fields[j].Type, names)
		}
	}
}

func pointerizeFieldType(goType string, structNames map[string]struct{}) string {
	switch goType {
	case "string", "bool", "int", "float64", goTypeInterface, goTypeSliceInterface:
		return goType
	}
	if strings.HasPrefix(goType, "[]") {
		elem := goType[2:]
		if _, ok := structNames[elem]; ok {
			return "[]*" + elem
		}
		return goType
	}
	if _, ok := structNames[goType]; ok {
		return "*" + goType
	}
	return goType
}

func looksLikeJSONObject(s string) bool {
	return len(s) > 0 && s[0] == '{'
}

func (b *schemaBuilder) infer(typeName string, v interface{}) string {
	switch val := v.(type) {
	case map[string]interface{}:
		b.defineStruct(typeName, val)
		return typeName
	case []interface{}:
		if len(val) == 0 {
			return goTypeSliceInterface
		}
		elemType := b.infer(typeName, val[0])
		return "[]" + elemType
	case string:
		return "string"
	case bool:
		return "bool"
	case float64:
		if val == float64(int64(val)) {
			return "int"
		}
		return "float64"
	case nil:
		return goTypeInterface
	default:
		return goTypeInterface
	}
}

func (b *schemaBuilder) defineStruct(typeName string, obj map[string]interface{}) {
	if _, ok := b.seen[typeName]; ok {
		return
	}
	b.seen[typeName] = struct{}{}

	usedNames := make(map[string]int)
	fields := make([]responseGoField, 0, len(obj))
	for key, value := range obj {
		fieldName := jsonKeyToIdent(key)
		if n := usedNames[fieldName]; n > 0 {
			fieldName = fmt.Sprintf("%s%d", fieldName, n+1)
		}
		usedNames[fieldName]++

		nestedName := typeName + singularizeJSONKey(key)
		goType := b.infer(nestedName, value)
		fields = append(fields, responseGoField{
			Name:    fieldName,
			JSONKey: key,
			Type:    goType,
		})
	}
	b.types = append(b.types, responseGoType{
		Name:   typeName,
		Fields: fields,
	})
}

// singularizeJSONKey 将 JSON 字段名转为嵌套类型名片段（projects -> Project）。
func singularizeJSONKey(key string) string {
	ident := jsonKeyToIdent(key)
	lower := strings.ToLower(key)
	if strings.HasSuffix(lower, "ies") && len(ident) > 3 && strings.HasSuffix(ident, "ies") {
		return ident[:len(ident)-3] + "y"
	}
	if strings.HasSuffix(ident, "s") && len(ident) > 1 {
		return ident[:len(ident)-1]
	}
	return ident
}
