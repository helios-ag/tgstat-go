package tgstat

import (
	"context"
	"encoding/json"
	"net/http"
	"tgstat/endpoints"
	"tgstat/schema"
)

func (c *Client) UsageStat(ctx context.Context) (*schema.StatResponse, *http.Response, error) {
	path := endpoints.UsageStat

	body := make(map[string]string)
	req, err := c.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.StatResponse
	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}
