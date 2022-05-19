package usage

import (
	"context"
	"encoding/json"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/endpoints"
	"net/http"
)

type Client struct {
	api   tgstat.API
	token string
}

// Stat request
// see https://api.tgstat.ru/docs/ru/usage/stat.html
func Stat(ctx context.Context) (*tgstat.StatResult, *http.Response, error) {
	return getClient().Stat(ctx)
}

// Stat request
// see https://api.tgstat.ru/docs/ru/usage/stat.html
func (c Client) Stat(ctx context.Context) (*tgstat.StatResult, *http.Response, error) {
	path := endpoints.UsageStat

	body := make(map[string]string)
	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response tgstat.StatResult
	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func getClient() Client {
	return Client{
		tgstat.GetAPI(),
		tgstat.Token,
	}
}
