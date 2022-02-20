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
	API   tgstat.API
	Token string
}

func PostGet(ctx context.Context, postId string) (*schema.PostResponse, *http.Response, error) {
	return getClient().PostGet(ctx, postId)
}

func (c Client) PostGet(ctx context.Context, postId string) (*schema.PostResponse, *http.Response, error) {
	path := endpoints.PostsGet

	if postId == "" {
		return nil, nil, fmt.Errorf("postId can not be empty")
	}

	body := make(map[string]string)
	body["postId"] = postId
	req, err := c.API.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.PostResponse
	result, err := c.API.Do(req, &response)
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

func PostStat(ctx context.Context, request PostStatRequest) (*schema.PostStatResponse, *http.Response, error) {
	return getClient().PostStat(ctx, request)
}

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

	req, err := c.API.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.PostStatResponse
	result, err := c.API.Do(req, &response)
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

func PostSearch(ctx context.Context, request PostSearchRequest) (*schema.PostSearchResponse, *http.Response, error) {
	return getClient().PostSearch(ctx, request)
}

func (c Client) PostSearch(ctx context.Context, request PostSearchRequest) (*schema.PostSearchResponse, *http.Response, error) {
	path := endpoints.PostsSearch

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := makeRequestBody(request)

	req, err := c.API.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.PostSearchResponse
	result, err := c.API.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func PostSearchExtended(ctx context.Context, request PostSearchRequest) (*schema.PostSearchExtendedResponse, *http.Response, error) {
	return getClient().PostSearchExtended(ctx, request)
}

func (c Client) PostSearchExtended(ctx context.Context, request PostSearchRequest) (*schema.PostSearchExtendedResponse, *http.Response, error) {
	path := endpoints.PostsSearch

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := makeRequestBody(request)

	body["extended"] = "1"

	req, err := c.API.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.PostSearchExtendedResponse
	result, err := c.API.Do(req, &response)
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
		if nil != request.HideForwards && *request.HideForwards == true {
			return "1"
		} else {
			return "0"
		}
	}()

	body["hideDeleted"] = func() string {
		if nil != request.HideDeleted && *request.HideDeleted == true {
			return "1"
		} else {
			return "0"
		}
	}()
	if nil != request.MinusWords {
		body["minusWords"] = *request.MinusWords
	}

	body["extendedSyntax"] = func() string {
		if nil != request.ExtendedSyntax && *request.ExtendedSyntax == true {
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