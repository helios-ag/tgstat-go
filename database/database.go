package database

import (
	"context"
	"encoding/json"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/endpoints"
	"github.com/helios-ag/tgstat-go/schema"
	"net/http"
)

type Client struct {
	api   tgstat.API
	token string
}

func CountriesGet(ctx context.Context, lang string) (*schema.CountryResponse, *http.Response, error) {
	return getClient().CountriesGet(ctx, lang)
}

func (c Client) CountriesGet(ctx context.Context, lang string) (*schema.CountryResponse, *http.Response, error) {
	path := endpoints.DatabaseCountries

	body := make(map[string]string)
	body["lang"] = lang
	req, err := c.api.NewRestRequest(ctx, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.CountryResponse
	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func CategoriesGet(ctx context.Context, lang string) (*schema.CategoryResponse, *http.Response, error) {
	return getClient().CategoriesGet(ctx, lang)
}

func (c Client) CategoriesGet(ctx context.Context, lang string) (*schema.CategoryResponse, *http.Response, error) {
	path := endpoints.DatabaseCategories

	body := make(map[string]string)
	body["lang"] = lang

	req, err := c.api.NewRestRequest(ctx, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.CategoryResponse
	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func LanguagesGet(ctx context.Context, lang string) (*schema.LanguageResponse, *http.Response, error) {
	return getClient().LanguagesGet(ctx, lang)
}

func (c Client) LanguagesGet(ctx context.Context, lang string) (*schema.LanguageResponse, *http.Response, error) {
	path := endpoints.DatabaseLanguages

	body := make(map[string]string)
	body["lang"] = lang
	req, err := c.api.NewRestRequest(ctx, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.LanguageResponse
	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func getClient() Client {
	return Client{tgstat.GetAPI(), tgstat.Token}
}
