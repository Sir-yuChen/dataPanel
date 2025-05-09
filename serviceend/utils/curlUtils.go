package utils

import (
	"bytes"
	"context"
	"dataPanel/serviceend/global"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Curl struct {
	client  *http.Client      // http client
	baseURL string            // base url
	headers map[string]string // headers
}

func NewCurl(baseURL string, timeout time.Duration) *Curl {
	return &Curl{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: strings.TrimSuffix(baseURL, "/"),
		headers: make(map[string]string),
	}
}
func (c *Curl) SetHeader(key, value string) {
	c.headers[key] = value
}
func (c *Curl) buildRequest(ctx context.Context, method, urlPath string, queryParams map[string]string, body io.Reader) (*http.Request, error) {
	fullURL := c.baseURL + urlPath
	if queryParams != nil {
		query := url.Values{}
		for key, value := range queryParams {
			query.Add(key, value)
		}
		fullURL += "?" + query.Encode()
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, err
	}
	// 设置请求头
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	return req, nil
}
func (c *Curl) doRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Curl) Get(ctx context.Context, urlPath string, queryParams map[string]string) (*http.Response, error) {
	req, err := c.buildRequest(ctx, http.MethodGet, urlPath, queryParams, nil)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req)
}

func (c *Curl) Post(ctx context.Context, urlPath string, body []byte) (*http.Response, error) {
	req, err := c.buildRequest(ctx, http.MethodPost, urlPath, nil, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	return c.doRequest(req)
}

func (c *Curl) Put(ctx context.Context, urlPath string, body []byte) (*http.Response, error) {
	req, err := c.buildRequest(ctx, http.MethodPut, urlPath, nil, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	return c.doRequest(req)
}

func (c *Curl) Delete(ctx context.Context, urlPath string) (*http.Response, error) {
	req, err := c.buildRequest(ctx, http.MethodDelete, urlPath, nil, nil)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req)
}
func (c *Curl) ReadResponseBody(resp *http.Response) ([]byte, error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			global.GvaLog.Error("close response body failed", zap.Error(err))
		}
	}(resp.Body)
	return io.ReadAll(resp.Body)
}
