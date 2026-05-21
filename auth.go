package main

import (
	"encoding/base64"
	"strings"
)

// resolveAuthorization 合并 -auth 与 -token（SonarQube 用户令牌：login=token，password 为空）。
func resolveAuthorization(authHeader, token string) string {
	if strings.TrimSpace(authHeader) != "" {
		return strings.TrimSpace(authHeader)
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return ""
	}
	raw := token + ":"
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(raw))
}
