package tgstat_go

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/helios-ag/tgstat-go/schema"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	APIURI string = "https://api.tgstat.ru"
)

// APIs are the currently supported endpoints.
type APIs struct {
	Api API
	mu  sync.RWMutex
}

var Token string

var apis APIs

type API interface {
	NewRestRequest(ctx context.Context, method, urlPath string, data map[string]string) (*http.Request, error)
	Do(r *http.Request, v interface{}) (*http.Response, error)
}

// ClientConfig is used to set client configuration
type ClientConfig struct {
	Token    string
	Endpoint string
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

// WithEndpoint configures a Client to use the specified API endpoint.
func WithEndpoint(endpoint string) {
	cfg.Endpoint = strings.TrimRight(endpoint, "/")
}

func (c *Client) NewRestRequest(ctx context.Context, method, urlPath string, data map[string]string) (*http.Request, error) {
	return newRestRequest(c, ctx, method, urlPath, data)
}

var newRestRequest = func(c *Client, ctx context.Context, method, urlPath string, data map[string]string) (*http.Request, error) {
	uri := APIURI + urlPath

	if c.Config.Endpoint != "" {
		uri = c.Config.Endpoint + urlPath
	}

	body := url.Values{}

	for key, value := range data {
		body.Add(key, value)
	}

	body.Add("token", c.Config.Token)

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

var reader = func(r io.Reader) ([]byte, error) {
	return ioutil.ReadAll(r)
}

// Do perform an HTTP request against the API.
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

	var respBody schema.ErrorResponse
	if err := json.Unmarshal(body, &respBody); err != nil {
		return nil
	}
	if respBody.Error == "" {
		return nil
	}
	return fmt.Errorf(respBody.Error)
}

func (c *ClientConfig) validate() error {
	if c.Token == "" {
		return errors.New("token can't be empty")
	}

	if _, err := url.Parse(c.Endpoint); err != nil {
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

// newAPI creates a new client.
func newAPI(cfg *ClientConfig, options ...ClientOption) *Client {
	client := &Client{
		Config:     cfg,
		httpClient: &http.Client{},
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func GetAPI(options ...ClientOption) API {
	var api API

	apis.mu.RLock()
	api = apis.Api
	apis.mu.RUnlock()

	if api != nil {
		return api
	}

	return newAPI(&cfg, options...)
}

var cfg ClientConfig

func SetConfig(config ClientConfig) {
	cfg = config
}
