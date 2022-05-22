package tgstat_go

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
	"sync"
)

const (
	APIURL string = "https://api.tgstat.ru"
)

// APIs are the currently supported endpoints.
type APIs struct {
	Api API
	mu  sync.RWMutex
}

var Token string

var apis APIs

type API interface {
	NewRestRequest(ctx context.Context, token, method, urlPath string, data map[string]string) (*http.Request, error)
	Do(r *http.Request, v interface{}) (*http.Response, error)
}

// Client is a client to TG Stat API
type Client struct {
	Url        string
	httpClient *http.Client
}

var TGStatClient Client

// ClientOption is used to configure a Client.
type ClientOption func(*Client)

// WithEndpoint configures a Client to use the specified API endpoint.
func WithEndpoint(endpoint string) {
	TGStatClient.Url = strings.TrimRight(endpoint, "/")
}

func (c *Client) NewRestRequest(ctx context.Context, token, method, urlPath string, data map[string]string) (*http.Request, error) {
	return NewRestRequest(c, ctx, token, method, urlPath, data)
}

var NewRestRequest = func(c *Client, ctx context.Context, token, method, urlPath string, data map[string]string) (*http.Request, error) {
	uri := APIURL + urlPath

	if c == nil {
		return nil, errors.New("client not configured")
	}

	if c.Url != "" {
		uri = c.Url + urlPath
	}

	if token == "" {
		return nil, errors.New("token not found")
	}

	//var body string
	body := url.Values{}

	for key, value := range data {
		body.Add(key, value)
	}

	body.Add("token", token)
	reqBodyData, _ := json.Marshal(body)
	// On `GET`, move the payload into the URL
	if method == http.MethodGet {
		uri += "?" + body.Encode()
		reqBodyData = nil
	}

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

	err = errorFromResponse(resp, body)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode >= 400 && resp.StatusCode <= 599 {
		err = fmt.Errorf("tgstat server responded with status code %d", resp.StatusCode)
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

	var respBody ErrorResult
	if err := json.Unmarshal(body, &respBody); err != nil {
		return nil
	}
	if respBody.Error == "" || respBody.VerifyCode != "" {
		return nil
	}
	return fmt.Errorf(respBody.Error)
}

// NewClient creates a new client.
func newClient(uri string, options ...ClientOption) (*Client, error) {

	if uri == "" {
		return nil, errors.New("URL is empty")
	}

	if _, err := url.ParseRequestURI(uri); err != nil {
		return nil, fmt.Errorf("unable to parse URL: %v", err)
	}

	client := &Client{
		Url:        uri,
		httpClient: &http.Client{},
	}

	for _, option := range options {
		option(client)
	}

	return client, nil
}

// newAPI creates a new client.
func newAPI(url string, options ...ClientOption) *Client {
	client := &Client{
		Url:        url,
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

	if TGStatClient.Url == "" {
		TGStatClient.Url = APIURL
	}

	return newAPI(TGStatClient.Url, options...)
}

func String(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func Bool(b bool) *bool {
	return &b
}

func Int(v int) *int {
	return &v
}
