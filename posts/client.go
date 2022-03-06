package posts

import (
	"context"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/endpoints"
	"github.com/helios-ag/tgstat-go/schema"
	"net/http"
	"strconv"
)

type Client struct {
	api   tgstat.API
	token string
}

// Get request
// see https://api.tgstat.ru/docs/ru/posts/get.html
func Get(ctx context.Context, postId string) (*schema.PostResponse, *http.Response, error) {
	return getClient().Get(ctx, postId)
}

// Get request
// see https://api.tgstat.ru/docs/ru/posts/get.html
func (c Client) Get(ctx context.Context, postId string) (*schema.PostResponse, *http.Response, error) {
	path := endpoints.PostsGet

	if postId == "" {
		return nil, nil, fmt.Errorf("postId can not be empty")
	}

	body := make(map[string]string)
	body["postId"] = postId
	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.PostResponse
	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

type PostStatRequest struct {
	PostId string
	Group  *string
}

func (postStatRequest PostStatRequest) Validate() error {
	return validation.ValidateStruct(&postStatRequest,
		validation.Field(&postStatRequest.PostId, validation.Required),
		validation.Field(&postStatRequest.Group, validation.In("hour", "day")),
	)
}

// PostStat request
// see https://api.tgstat.ru/docs/ru/posts/get.html
func PostStat(ctx context.Context, request PostStatRequest) (*schema.PostStatResponse, *http.Response, error) {
	return getClient().PostStat(ctx, request)
}

// PostStat request
// see https://api.tgstat.ru/docs/ru/posts/stat.html
func (c Client) PostStat(ctx context.Context, request PostStatRequest) (*schema.PostStatResponse, *http.Response, error) {
	path := endpoints.PostsStat

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["postId"] = request.PostId
	if nil != request.Group {
		body["group"] = *request.Group
	}

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.PostStatResponse
	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

type PostSearchRequest struct {
	Q              string
	Limit          *int
	Offset         *int
	PeerType       *string
	StartDate      *string
	EndDate        *string
	HideForwards   *bool
	HideDeleted    *bool
	StrongSearch   *bool
	MinusWords     *string
	ExtendedSyntax *bool
}

func (postSearchRequest PostSearchRequest) Validate() error {
	return validation.ValidateStruct(&postSearchRequest,
		validation.Field(&postSearchRequest.Q, validation.Required),
		validation.Field(&postSearchRequest.Limit, validation.Max(50)),
		validation.Field(&postSearchRequest.Offset, validation.Max(50)),
	)
}

// PostSearch request
// see https://api.tgstat.ru/docs/ru/posts/search.html
func PostSearch(ctx context.Context, request PostSearchRequest) (*schema.PostSearchResponse, *http.Response, error) {
	return getClient().PostSearch(ctx, request)
}

// PostSearch request
// see https://api.tgstat.ru/docs/ru/posts/search.html
func (c Client) PostSearch(ctx context.Context, request PostSearchRequest) (*schema.PostSearchResponse, *http.Response, error) {
	path := endpoints.PostsSearch

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := makeRequestBody(request)

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.PostSearchResponse
	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

// PostSearchExtended request
// see https://api.tgstat.ru/docs/ru/posts/search.html
func PostSearchExtended(ctx context.Context, request PostSearchRequest) (*schema.PostSearchExtendedResponse, *http.Response, error) {
	return getClient().PostSearchExtended(ctx, request)
}

// PostSearchExtended request
// see https://api.tgstat.ru/docs/ru/posts/search.html
func (c Client) PostSearchExtended(ctx context.Context, request PostSearchRequest) (*schema.PostSearchExtendedResponse, *http.Response, error) {
	path := endpoints.PostsSearch

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := makeRequestBody(request)

	body["extended"] = "1"

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.PostSearchExtendedResponse
	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func makeRequestBody(request PostSearchRequest) map[string]string {
	body := make(map[string]string)
	body["q"] = request.Q
	if nil != request.Limit {
		body["limit"] = strconv.Itoa(*request.Limit)
	}

	if nil != request.Offset {
		body["offset"] = strconv.Itoa(*request.Offset)
	}

	if nil != request.PeerType {
		body["peerType"] = *request.PeerType
	}

	if nil != request.StartDate {
		body["startDate"] = *request.StartDate
	}

	if nil != request.EndDate {
		body["EndDate"] = *request.EndDate
	}

	body["hideForwards"] = func() string {
		if nil != request.HideForwards && *request.HideForwards {
			return "1"
		} else {
			return "0"
		}
	}()

	body["hideDeleted"] = func() string {
		if nil != request.HideDeleted && *request.HideDeleted {
			return "1"
		} else {
			return "0"
		}
	}()
	if nil != request.MinusWords {
		body["minusWords"] = *request.MinusWords
	}

	body["extendedSyntax"] = func() string {
		if nil != request.ExtendedSyntax && *request.ExtendedSyntax {
			return "1"
		} else {
			return "0"
		}
	}()

	return body
}

func getClient() Client {
	return Client{tgstat.GetAPI(), tgstat.Token}
}
