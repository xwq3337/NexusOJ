package chttp

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client HTTP客户端封装
type Client struct {
	cli     *resty.Client
	baseURL string
	headers map[string]string
}

// New 创建新的HTTP客户端
func New(opts ...Option) *Client {
	c := &Client{
		cli:     resty.New(),
		headers: make(map[string]string),
	}

	// 应用配置选项
	for _, opt := range opts {
		opt(c)
	}

	// 设置默认超时时间
	if c.cli.GetClient().Timeout == 0 {
		c.SetTimeout(10 * time.Second)
	}

	return c
}

// Option 客户端配置函数
type Option func(*Client)

// WithBaseURL 设置基础URL
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithTimeout 设置请求超时时间
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.SetTimeout(d)
	}
}

// WithRetry 设置重试策略
func WithRetry(count int, waitTime, maxWaitTime time.Duration) Option {
	return func(c *Client) {
		c.cli.SetRetryCount(count).
			SetRetryWaitTime(waitTime).
			SetRetryMaxWaitTime(maxWaitTime)
	}
}

// WithHeader 设置公共请求头
func WithHeader(key, value string) Option {
	return func(c *Client) {
		c.headers[key] = value
	}
}

// WithDebug 启用调试模式
func WithDebug(enable bool) Option {
	return func(c *Client) {
		c.cli.SetDebug(enable)
	}
}

// SetTimeout 设置请求超时时间
func (c *Client) SetTimeout(d time.Duration) {
	c.cli.SetTimeout(d)
}

// SetBaseURL 设置基础URL
func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

// SetHeader 设置公共请求头
func (c *Client) SetHeader(key, value string) {
	c.headers[key] = value
}

// buildRequest 构建请求对象
func (c *Client) buildRequest(ctx context.Context) *resty.Request {
	req := c.cli.R().SetContext(ctx)
	for k, v := range c.headers {
		req.SetHeader(k, v)
	}
	return req
}

// Get 发送GET请求
func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	url := c.baseURL + path
	resp, err := c.buildRequest(ctx).
		SetResult(result).
		Get(url)

	return c.handleResponse(resp, err)
}

// Post 发送POST请求
func (c *Client) Post(ctx context.Context, path string, body, result interface{}) error {
	url := c.baseURL + path
	resp, err := c.buildRequest(ctx).
		SetBody(body).
		SetResult(result).
		Post(url)

	return c.handleResponse(resp, err)
}

// Put 发送PUT请求
func (c *Client) Put(ctx context.Context, path string, body, result interface{}) error {
	url := c.baseURL + path
	resp, err := c.buildRequest(ctx).
		SetBody(body).
		SetResult(result).
		Put(url)

	return c.handleResponse(resp, err)
}

// Delete 发送DELETE请求
func (c *Client) Delete(ctx context.Context, path string, result interface{}) error {
	url := c.baseURL + path
	resp, err := c.buildRequest(ctx).
		SetResult(result).
		Delete(url)

	return c.handleResponse(resp, err)
}

// handleResponse 统一处理响应
func (c *Client) handleResponse(resp *resty.Response, err error) error {
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("http error: %s, status code: %d", resp.Status(), resp.StatusCode())
	}

	return nil
}

// Request 自定义请求方法
func (c *Client) Request(ctx context.Context) *RequestBuilder {
	return &RequestBuilder{
		client: c,
		req:    c.buildRequest(ctx),
	}
}

// RequestBuilder 请求构建器
type RequestBuilder struct {
	client *Client
	req    *resty.Request
}

// SetHeader 设置请求头
func (b *RequestBuilder) SetHeader(key, value string) *RequestBuilder {
	b.req.SetHeader(key, value)
	return b
}

// SetQueryParam 设置查询参数
func (b *RequestBuilder) SetQueryParam(key, value string) *RequestBuilder {
	b.req.SetQueryParam(key, value)
	return b
}

// SetBody 设置请求体
func (b *RequestBuilder) SetBody(body interface{}) *RequestBuilder {
	b.req.SetBody(body)
	return b
}

// SetResult 设置结果接收对象
func (b *RequestBuilder) SetResult(result interface{}) *RequestBuilder {
	b.req.SetResult(result)
	return b
}

// SetError 设置错误接收对象
func (b *RequestBuilder) SetError(error interface{}) *RequestBuilder {
	b.req.SetError(error)
	return b
}

// Execute 执行请求
func (b *RequestBuilder) Execute(method, path string) error {
	url := b.client.baseURL + path
	resp, err := b.req.Execute(method, url)
	return b.client.handleResponse(resp, err)
}
