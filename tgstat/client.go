package tgstat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"tgstat/schema"
)

const (
	APIURI string = "https://api.tgstat.ru"
)

// ClientConfig is used to set client configuration
type ClientConfig struct {
	token    string
	extended string
	endpoint string
}

// Client is a client to SB API
type Client struct {
	Config     *ClientConfig
	httpClient *http.Client
}

// Body struct
type Body struct {
	Token *string `json:"token"`
}

// ClientOption is used to configure a Client.
type ClientOption func(*Client)

// WithToken configures a Client to use the specified token for authentication.
func WithToken(token string) ClientOption {
	return func(client *Client) {
		client.Config.token = token
	}
}

// WithExtendedResponse configures a Client to receive extended response
func WithExtendedResponse() ClientOption {
	return func(client *Client) {
		client.Config.extended = "1"
	}
}

// WithEndpoint configures a Client to use the specified API endpoint.
func WithEndpoint(endpoint string) ClientOption {
	return func(client *Client) {
		client.Config.endpoint = strings.TrimRight(endpoint, "/")
	}
}

func (c *Client) NewRestRequest(ctx context.Context, method, urlPath string, data map[string]string) (*http.Request, error) {
	return newRestRequest(c, ctx, method, urlPath, data)
}

var newRestRequest = func(c *Client, ctx context.Context, method, urlPath string, data map[string]string) (*http.Request, error) {
	uri := APIURI + urlPath

	if c.Config.endpoint != "" {
		uri = c.Config.endpoint + urlPath
	}

	body := url.Values{}

	for key, value := range data {
		body.Add(key, value)
	}

	body.Add("token", c.Config.token)

	reqData := body.Encode()
	req, err := http.NewRequest(method, uri, strings.NewReader(reqData))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req = req.WithContext(ctx)
	return req, nil
}

// NewRequest creates an HTTP request against the API (mobile payments). The returned request
// is assigned with ctx and has all necessary headers set (auth, user agent, etc.).
// NewRestRequest creates an HTTP request against the API with 'rest' in path. The returned request
// is assigned with ctx and has all necessary headers set (auth, user agent, etc.).
func (c *Client) NewRequest(ctx context.Context, method, urlPath string, data interface{}) (*http.Request, error) {
	return newRequest(c, ctx, method, urlPath, data)
}

var newRequest = func(c *Client, ctx context.Context, method, urlPath string, data interface{}) (*http.Request, error) {
	if strings.Contains(urlPath, "rest") {
		return nil, fmt.Errorf("path contains rest request, use NewRestRequest instead")
	}

	uri := APIURI + urlPath

	if c.Config.endpoint != "" {
		uri = c.Config.endpoint + urlPath
	}

	reqBodyData, _ := json.Marshal(data)

	req, err := http.NewRequest(method, uri, bytes.NewReader(reqBodyData))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")

	req = req.WithContext(ctx)

	return req, nil
}

var reader = func(r io.Reader) ([]byte, error) {
	return ioutil.ReadAll(r)
}

// Do performs an HTTP request against the API.
func (c *Client) Do(r *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	body, err := reader(resp.Body)
	if err != nil {
		resp.Body.Close()
		return resp, err
	}
	resp.Body.Close()
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))

	if resp.StatusCode >= 400 && resp.StatusCode <= 599 {
		err = errorFromResponse(resp, body)
		if err == nil {
			err = fmt.Errorf("tgstat server responded with status code %d", resp.StatusCode)
		}
		return resp, err
	}
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, bytes.NewReader(body))
		} else {
			err = json.Unmarshal(body, v)
		}
	}

	return resp, err
}

func errorFromResponse(resp *http.Response, body []byte) error {
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		return nil
	}

	var respBody schema.Response
	if err := json.Unmarshal(body, &respBody); err != nil {
		return nil
	}
	if respBody.Error == "" {
		return nil
	}
	return fmt.Errorf(respBody.Error)
}

func (c *ClientConfig) validate() error {
	if c.token != "" {
		return errors.New("token can't be empty")
	}

	if _, err := url.Parse(c.endpoint); err != nil {
		return fmt.Errorf("unable to parse URL: %v", err)
	}

	return nil
}

// NewClient creates a new client.
func NewClient(cfg *ClientConfig, options ...ClientOption) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("passed in config cannot be nil")
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("unable to validate given config: %v", err)
	}

	client := &Client{
		Config:     cfg,
		httpClient: &http.Client{},
	}

	for _, option := range options {
		option(client)
	}

	return client, nil
}
