package database

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

// CountriesGet request
// See https://api.tgstat.ru/docs/ru/database/countries.html
func CountriesGet(ctx context.Context, lang string) (*tgstat.CountryResult, *http.Response, error) {
	return getClient().CountriesGet(ctx, lang)
}

// CountriesGet request
// See https://api.tgstat.ru/docs/ru/database/countries.html
func (c Client) CountriesGet(ctx context.Context, lang string) (*tgstat.CountryResult, *http.Response, error) {
	path := endpoints.DatabaseCountries

	body := make(map[string]string)
	body["lang"] = lang
	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response tgstat.CountryResult
	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

// CategoriesGet request
// See https://api.tgstat.ru/docs/ru/database/categories.html
func CategoriesGet(ctx context.Context, lang string) (*tgstat.CategoryResult, *http.Response, error) {
	return getClient().CategoriesGet(ctx, lang)
}

// CategoriesGet request
// See https://api.tgstat.ru/docs/ru/database/categories.html
func (c Client) CategoriesGet(ctx context.Context, lang string) (*tgstat.CategoryResult, *http.Response, error) {
	path := endpoints.DatabaseCategories

	body := make(map[string]string)
	body["lang"] = lang

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response tgstat.CategoryResult
	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

// LanguagesGet request
// See https://api.tgstat.ru/docs/ru/database/languages.html
func LanguagesGet(ctx context.Context, lang string) (*tgstat.LanguageResult, *http.Response, error) {
	return getClient().LanguagesGet(ctx, lang)
}

// LanguagesGet request
// See https://api.tgstat.ru/docs/ru/database/languages.html
func (c Client) LanguagesGet(ctx context.Context, lang string) (*tgstat.LanguageResult, *http.Response, error) {
	path := endpoints.DatabaseLanguages

	body := make(map[string]string)
	body["lang"] = lang
	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response tgstat.LanguageResult
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
