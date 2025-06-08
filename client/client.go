package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
)

type ClientInterface interface {
	Do(req *http.Request) (resp *http.Response, err error)
	DoWithParse(req *http.Request, value interface{}) (*Response, error)
	URL(uri string) string
	Get(ctx context.Context, uri string, headers map[string]string, value interface{}) (resp *Response, err error)
	Post(ctx context.Context, uri string, headers map[string]string, body io.Reader, value interface{}) (resp *Response, err error)
	Put(ctx context.Context, uri string, headers map[string]string, body io.Reader, value interface{}) (resp *Response, err error)
	Delete(ctx context.Context, uri string, headers map[string]string, body io.Reader, value interface{}) (resp *Response, err error)
	Patch(ctx context.Context, uri string, headers map[string]string, body io.Reader, value interface{}) (resp *Response, err error)
}

type Client struct {
	httpClient *http.Client
	Config     *Config
}

func NewClient(config *Config) *Client {
	client := &http.Client{}
	if config.IgnoreCertificateErrors {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: transport}
	}

	if config.Timeout > 0 {
		client.Timeout = config.Timeout
	}

	return &Client{
		httpClient: client,
		Config:     config,
	}
}

func (c *Client) Do(req *http.Request) (resp *http.Response, err error) {
	ctx := req.Context()
	resp, err = c.httpClient.Do(req.WithContext(ctx))
	return
}

func (c *Client) DoWithParse(req *http.Request, value interface{}) (*Response, error) {
	underlying, err := c.Do(req)
	parsedResp := Response{
		HTTPResponse: underlying,
	}
	if err != nil {
		if underlying != nil {
			return &parsedResp, err
		}
		return nil, err
	}

	err = parsedResp.Unmarshal(value)
	return &parsedResp, err
}

func (c *Client) Get(ctx context.Context, uri string, headers map[string]string, value interface{}) (resp *Response, err error) {
	return c.do(ctx, uri, http.MethodGet, headers, nil, value)
}

func (c *Client) Post(ctx context.Context, uri string, headers map[string]string, body io.Reader, value interface{}) (resp *Response, err error) {
	return c.do(ctx, uri, http.MethodPost, headers, body, value)
}

func (c *Client) Put(ctx context.Context, uri string, headers map[string]string, body io.Reader, value interface{}) (resp *Response, err error) {
	return c.do(ctx, uri, http.MethodPut, headers, body, value)
}

func (c *Client) Delete(ctx context.Context, uri string, headers map[string]string, body io.Reader, value interface{}) (resp *Response, err error) {
	return c.do(ctx, uri, http.MethodDelete, headers, body, value)
}

func (c *Client) Patch(ctx context.Context, uri string, headers map[string]string, body io.Reader, value interface{}) (resp *Response, err error) {
	return c.do(ctx, uri, http.MethodPatch, headers, body, value)
}

func (c *Client) do(ctx context.Context, uri, method string, headers map[string]string, body io.Reader, value interface{}) (resp *Response, err error) {
	req, err := http.NewRequest(method, c.URL(uri), body)
	if err != nil {
		return nil, err
	}
	header := http.Header{}
	for key, value := range headers {
		header.Add(key, value)
	}
	header.Set("Content-Type", "application/json")
	req.Header = header

	return c.DoWithParse(req.WithContext(ctx), value)
}

func (c *Client) URL(uri string) string {
	if len(uri) > 0 && uri[0] == '/' {
		uri = uri[1:]
	}
	return c.AddScheme(fmt.Sprintf("%s/%s", c.Config.EndpointURL, uri))
}

func (c *Client) AddScheme(url string) string {
	if c.Config.UseSecure {
		return fmt.Sprintf("https://%s", url)
	}
	return fmt.Sprintf("http://%s", url)

}
