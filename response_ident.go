package main

import (
	"regexp"
	"strings"
	"unicode"
)

var jsonKeyTokenRE = regexp.MustCompile(`[^a-zA-Z0-9]+`)

var goKeywords = map[string]struct{}{
	"break": {}, "default": {}, "func": {}, "interface": {}, "select": {},
	"case": {}, "defer": {}, "go": {}, "map": {}, "struct": {},
	"chan": {}, "else": {}, "goto": {}, "package": {}, "switch": {},
	"const": {}, "fallthrough": {}, "if": {}, "range": {}, "type": {},
	"continue": {}, "for": {}, "import": {}, "return": {}, "var": {},
}

// jsonKeyToIdent 将 JSON 字段名转为合法 Go 标识符（导出）。
func jsonKeyToIdent(key string) string {
	tokens := jsonKeyTokenRE.Split(key, -1)
	var b strings.Builder
	for _, tok := range tokens {
		if tok == "" {
			continue
		}
		if len(tok) == 1 {
			b.WriteString(strings.ToUpper(tok))
			continue
		}
		b.WriteString(strings.ToUpper(tok[:1]) + tok[1:])
	}
	ident := b.String()
	if ident == "" {
		return "Field"
	}
	r := rune(ident[0])
	if unicode.IsDigit(r) {
		ident = "Key" + ident
	}
	if _, ok := goKeywords[strings.ToLower(ident)]; ok {
		ident += "Value"
	}
	return ident
}
