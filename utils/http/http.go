// Package http
// @title
// @description
// @author njy
// @since 2022/12/1 13:27
package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

var client *http.Client

const (
	contentType = "application/json;charset=utf-8"
)

func init() {
	client = &http.Client{}
}

// Get 函数用于发送一个 Get 请求
func Get(url string, params map[string]string) ([]byte, error) {
	return withCloser(http.Get(generateUrl(url, params)))
}

// Post 函数用于发送一个 Post 请求
func Post(url, body string) ([]byte, error) {
	return withCloser(http.Post(url, contentType, bytes.NewReader([]byte(body))))
}

// Download 函数用于下载文件
func Download(url string) (string, error) {
	filename := path.Base(url)
	resp, err := http.Get(url)
	if err != nil {
		return filename, err
	}
	defer resp.Body.Close()
	out, err := os.Create(filename)
	defer out.Close()
	if err != nil {
		return filename, err
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return filename, err
	}
	return filename, nil
}

func generateUrl(url string, params map[string]string) string {
	if len(params) > 0 {
		values := make([]string, 0, len(params))
		for k, v := range params {
			values = append(values, fmt.Sprintf("%s=%s", k, v))
		}
		url += "?" + strings.Join(values, "&")
	}
	return url
}

// GetWithHeaders 函数用于发送一个 Get 请求并带有请求头
func GetWithHeaders(url string, headers, params map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", generateUrl(url, params), nil)
	return requestWithCloser(req, err, headers)
}

// PostWithHeaders 函数用于发送一个 POST 请求并带有请求头
func PostWithHeaders(url, body string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(body)))
	return requestWithCloser(req, err, headers)
}

func requestWithCloser(req *http.Request, err error, headers map[string]string) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	// 设置请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return withCloser(client.Do(req))
}

func withCloser(resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf := bytes.NewBuffer(nil)
	// 使用io.Copy将响应内容拷贝到bytes.Buffer中
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
