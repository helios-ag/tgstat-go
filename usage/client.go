package usage

import (
	"context"
	"encoding/json"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/endpoints"
	"github.com/helios-ag/tgstat-go/schema"
	"net/http"
)

type Client struct {
	API   tgstat.API
	Token string
}

func Stat(ctx context.Context) (*schema.StatResponse, *http.Response, error) {
	return getClient().Stat(ctx)
}

func (c Client) Stat(ctx context.Context) (*schema.StatResponse, *http.Response, error) {
	path := endpoints.UsageStat

	body := make(map[string]string)
	req, err := c.API.NewRestRequest(ctx, http.MethodGet, path, body)

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
