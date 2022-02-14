package usage

import (
	"context"
	"encoding/json"
	"net/http"
	tgstat "tgstat"
	"tgstat/endpoints"
	"tgstat/schema"
)

type Client struct {
	API   tgstat.API
	Token string
}

func UsageStat(ctx context.Context) (*schema.StatResponse, *http.Response, error) {
	return getClient().UsageStat(ctx)
}

func (c Client) UsageStat(ctx context.Context) (*schema.StatResponse, *http.Response, error) {
	path := endpoints.UsageStat

	body := make(map[string]string)
	req, err := c.API.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.StatResponse
	result, err := c.API.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func getClient() Client {
	return Client{tgstat.GetAPI(), tgstat.Token}
}
