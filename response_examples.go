package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	responseExamplePath    = "/api/webservices/response_example"
	responseExampleWorkers = 8
)

type responseExamplePayload struct {
	Format  string `json:"format"`
	Example string `json:"example"`
}

// enrichResponseSchemas 为 hasResponseExample=true 的 action 拉取 JSON 示例并生成结构体定义。
func enrichResponseSchemas(client *http.Client, host, auth string, def *apiDefinition) {
	if client == nil {
		client = http.DefaultClient
	}
	type job struct {
		ws     *webService
		action *action
	}
	var jobs []job
	for _, ws := range def.WebServices {
		for _, action := range ws.Actions {
			if action.HasResponseExample {
				jobs = append(jobs, job{ws: ws, action: action})
			}
		}
	}
	if len(jobs) == 0 {
		return
	}

	sem := make(chan struct{}, responseExampleWorkers)
	var wg sync.WaitGroup
	for _, j := range jobs {
		wg.Add(1)
		go func(ws *webService, action *action) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			example, format, err := fetchResponseExample(client, host, auth, ws.Path, action.Key)
			if err != nil {
				log.Printf("warn: response example %s/%s: %v", ws.Path, action.Key, err)
				return
			}
			if !isJSONResponseExample(format, example) {
				return
			}
			okType := action.MethodName() + responseOKSuffix
			types, err := buildResponseSchema(okType, example)
			if err != nil {
				log.Printf("warn: response schema %s/%s: %v", ws.Path, action.Key, err)
				return
			}
			if len(types) == 0 {
				return
			}
			action.ResponseOKType = okType
			action.ResponseTypes = types
		}(j.ws, j.action)
	}
	wg.Wait()
}

func isJSONResponseExample(format, example string) bool {
	if format != "" && !strings.EqualFold(strings.TrimSpace(format), "json") {
		return false
	}
	return looksLikeJSONObject(strings.TrimSpace(example))
}

func fetchResponseExample(
	client *http.Client,
	host, auth, controller, actionKey string,
) (example, format string, err error) {
	u, err := url.Parse(host + responseExamplePath)
	if err != nil {
		return "", "", err
	}
	q := u.Query()
	q.Set("controller", controller)
	q.Set("action", actionKey)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", "", err
	}
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode == http.StatusUnauthorized {
		return "", "", errors.New("authorization failed to fetch response example")
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("response example status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	var payload responseExamplePayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", "", err
	}
	if payload.Example == "" {
		return "", "", fmt.Errorf("empty example for %s/%s", controller, actionKey)
	}
	return payload.Example, payload.Format, nil
}
