// Package httpclient 提供通用 HTTP 请求工具（无全局状态）。
package httpclient

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

// Request 发送 HTTP 请求（行为与 utils.HttpRequest 一致）。
func Request(
	urlStr string,
	method string,
	headers map[string]string,
	params map[string]string,
	data any) (*http.Response, error) {

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	query := u.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	u.RawQuery = query.Encode()

	buf := new(bytes.Buffer)
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return http.DefaultClient.Do(req)
}
