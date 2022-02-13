package tgstat

import (
	"context"
	"encoding/json"
	"fmt"
	//"github.com/helios-ag/tgstat-go/endpoints"
	//"github.com/helios-ag/tgstat-go/schema"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
	"strconv"
	"tgstat/endpoints"
	"tgstat/schema"
)

func (c *Client) ChannelGet(ctx context.Context, channelId string) (*schema.ChannelResponse, *http.Response, error) {
	path := endpoints.ChannelsGet

	if err := validateGetChannelId(channelId); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = channelId
	req, err := c.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelResponse
	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func validateGetChannelId(channelId string) error {
	if channelId == "" {
		return fmt.Errorf("ChannelID must be set")
	}
	return nil
}

type SearchRequest struct {
	Q                   string
	SearchByDescription int
	Country             string
	Language            string
	Category            string
	Limit               int
}

func (searchRequest SearchRequest) Validate() error {
	return validation.ValidateStruct(&searchRequest,
		validation.Field(&searchRequest.Country, validation.Required),
		validation.Field(&searchRequest.Q, validation.Required.When(searchRequest.Category == "").Error("Either query or category is required.")),
		validation.Field(&searchRequest.Category, validation.Required.When(searchRequest.Q == "").Error("Either query or category  is required.")),
	)
}

func (c *Client) ChannelSearch(ctx context.Context, request SearchRequest) (*schema.ChannelSearchResponse, *http.Response, error) {
	path := endpoints.ChannelsSearch

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["q"] = request.Q
	body["search_by_description"] = strconv.Itoa(request.SearchByDescription)
	body["country"] = request.Country
	body["language"] = request.Language
	body["category"] = request.Category
	body["limit"] = strconv.Itoa(request.Limit)

	req, err := c.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelSearchResponse
	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func (c *Client) ChannelStat(ctx context.Context, channelId string) (*schema.ChannelStatResponse, *http.Response, error) {
	path := endpoints.ChannelsStat

	if err := validateGetChannelId(channelId); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = channelId
	req, err := c.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelStatResponse
	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

type PostsRequest struct {
	ChannelId    string
	Limit        *uint64
	Offset       *uint64
	StartTime    *string
	EndTime      *string
	HideForwards *int
	HideDeleted  *int
}

func (postsRequest PostsRequest) Validate() error {
	return validation.ValidateStruct(&postsRequest,
		validation.Field(&postsRequest.ChannelId, validation.Required),
		validation.Field(&postsRequest.StartTime, validation.Date("1643113399")),
		validation.Field(&postsRequest.EndTime, validation.Date("1643113399")),
		validation.Field(&postsRequest.Offset, validation.Min(0), validation.Min(1000)),
		validation.Field(&postsRequest.Limit, validation.Min(0), validation.Min(50)),
		validation.Field(&postsRequest.HideForwards, validation.In(0, 1)),
		validation.Field(&postsRequest.HideDeleted, validation.In(0, 1)),
	)
}

func (c *Client) ChannelPosts(ctx context.Context, request PostsRequest) (*schema.ChannelPostsResponse, *http.Response, error) {
	path := endpoints.ChannelsStat

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = request.ChannelId
	body["limit"] = strconv.FormatUint(*request.Limit, 10)
	body["offset"] = strconv.FormatUint(*request.Offset, 10)
	body["startTime"] = *request.StartTime
	body["endTime"] = *request.EndTime
	body["hideForwards"] = strconv.Itoa(*request.HideForwards)
	body["hideDeleted"] = strconv.Itoa(*request.HideDeleted)
	body["extended"] = "0"

	req, err := c.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}
	var response schema.ChannelPostsResponse

	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func (c *Client) ChannelPostsExtended(ctx context.Context, request PostsRequest) (*schema.ChannelPostsWithChannelResponse, *http.Response, error) {
	path := endpoints.ChannelsStat

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = request.ChannelId
	body["limit"] = strconv.FormatUint(*request.Limit, 10)
	body["offset"] = strconv.FormatUint(*request.Offset, 10)
	body["startTime"] = *request.StartTime
	body["endTime"] = *request.EndTime
	body["hideForwards"] = strconv.Itoa(*request.HideForwards)
	body["hideDeleted"] = strconv.Itoa(*request.HideDeleted)
	body["extended"] = "1"

	req, err := c.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelPostsWithChannelResponse

	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

type ChannelMentionsRequest struct {
	ChannelId string
	Limit     *uint64
	Offset    *uint64
	StartDate *string
	EndDate   *string
}

func (channelMentionsRequest ChannelMentionsRequest) Validate() error {
	return validation.ValidateStruct(&channelMentionsRequest,
		validation.Field(&channelMentionsRequest.ChannelId, validation.Required),
		validation.Field(&channelMentionsRequest.Limit, validation.Min(0), validation.Max(50)),
		validation.Field(&channelMentionsRequest.Offset, validation.Min(0), validation.Max(1000)),
		validation.Field(&channelMentionsRequest.StartDate, validation.Date("1643113399")),
		validation.Field(&channelMentionsRequest.EndDate, validation.Date("1643113399")),
	)
}

func (c *Client) ChannelMentions(ctx context.Context, request ChannelMentionsRequest) (*schema.ChannelMentions, *http.Response, error) {
	path := endpoints.ChannelsMentions

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = request.ChannelId
	body["limit"] = strconv.FormatUint(*request.Limit, 10)
	body["offset"] = strconv.FormatUint(*request.Offset, 10)
	body["startDate"] = *request.StartDate
	body["endDate"] = *request.EndDate
	body["extended"] = "0"

	req, err := c.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelMentions

	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func (c *Client) ChannelMentionsExtended(ctx context.Context, request ChannelMentionsRequest) (*schema.ChannelMentionsExtended, *http.Response, error) {
	path := endpoints.ChannelsMentions

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = request.ChannelId
	body["limit"] = strconv.FormatUint(*request.Limit, 10)
	body["offset"] = strconv.FormatUint(*request.Offset, 10)
	body["startDate"] = *request.StartDate
	body["endDate"] = *request.EndDate
	body["extended"] = "1"

	req, err := c.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelMentionsExtended

	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

type ChannelForwardRequest struct {
	ChannelId string
	Limit     *uint64
	Offset    *uint64
	StartDate *string
	EndDate   *string
}

func (channelForwardRequest ChannelForwardRequest) Validate() error {
	return validation.ValidateStruct(&channelForwardRequest,
		validation.Field(&channelForwardRequest.ChannelId, validation.Required),
		validation.Field(&channelForwardRequest.Limit, validation.Min(0), validation.Max(50)),
		validation.Field(&channelForwardRequest.Offset, validation.Min(0), validation.Max(1000)),
		validation.Field(&channelForwardRequest.StartDate, validation.Date("1643113399")),
		validation.Field(&channelForwardRequest.EndDate, validation.Date("1643113399")),
	)
}

func (c *Client) ChannelForwards(ctx context.Context, request ChannelForwardRequest) (*schema.ChannelForwards, *http.Response, error) {
	path := endpoints.ChannelsForwards

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = request.ChannelId
	body["limit"] = strconv.FormatUint(*request.Limit, 10)
	body["offset"] = strconv.FormatUint(*request.Offset, 10)
	body["startDate"] = *request.StartDate
	body["endDate"] = *request.EndDate
	body["extended"] = "0"

	req, err := c.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelForwards

	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func (c *Client) ChannelForwardsExtended(ctx context.Context, request ChannelForwardRequest) (*schema.ChannelForwardsExtended, *http.Response, error) {
	path := endpoints.ChannelsForwards

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = request.ChannelId
	body["limit"] = strconv.FormatUint(*request.Limit, 10)
	body["offset"] = strconv.FormatUint(*request.Offset, 10)
	body["startDate"] = *request.StartDate
	body["endDate"] = *request.EndDate
	body["extended"] = "1"

	req, err := c.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelForwardsExtended

	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

type ChannelSubscribersRequest struct {
	ChannelId string
	StartDate *string
	EndDate   *string
	Group     *string
}

func (channelSubscribersRequest ChannelSubscribersRequest) Validate() error {
	return validation.ValidateStruct(&channelSubscribersRequest,
		validation.Field(&channelSubscribersRequest.ChannelId, validation.Required),
		validation.Field(&channelSubscribersRequest.StartDate, validation.Date("1643113399")),
		validation.Field(&channelSubscribersRequest.EndDate, validation.Date("1643113399")),
		validation.Field(&channelSubscribersRequest.Group, validation.In("hour", "day", "week", "month")),
	)
}

func (c *Client) ChannelSubscribers(ctx context.Context, request ChannelSubscribersRequest) (*schema.ChannelSubscribers, *http.Response, error) {
	path := endpoints.ChannelsSubscribers

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = request.ChannelId
	body["startDate"] = *request.StartDate
	body["endDate"] = *request.EndDate
	body["group"] = *request.Group

	req, err := c.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelSubscribers

	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func validateGroup(group string) bool {
	switch group {
	case
		"day",
		"week",
		"month":
		return true
	}

	return false
}

type ChannelViewsRequest struct {
	ChannelId string
	StartDate *string
	EndDate   *string
	Group     *string
}

func (channelViewsRequest ChannelViewsRequest) Validate() error {
	return validation.ValidateStruct(&channelViewsRequest,
		validation.Field(&channelViewsRequest.ChannelId, validation.Required),
		validation.Field(&channelViewsRequest.StartDate, validation.Date("1643113399")),
		validation.Field(&channelViewsRequest.EndDate, validation.Date("1643113399")),
		validation.Field(&channelViewsRequest.Group, validation.In("day", "week", "month")),
	)
}

func (c *Client) ChannelViews(ctx context.Context, request ChannelViewsRequest) (*schema.ChannelViews, *http.Response, error) {
	path := endpoints.ChannelsViews

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = request.ChannelId
	body["startDate"] = *request.StartDate
	body["endDate"] = *request.EndDate
	body["group"] = *request.Group

	req, err := c.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelViews

	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

type ChannelAddRequest struct {
	ChannelName string
	Country     *string
	Language    *string
	Category    *string
}

func (channelAddRequest ChannelAddRequest) Validate() error {
	return validation.ValidateStruct(&channelAddRequest,
		validation.Field(&channelAddRequest.ChannelName, validation.Required),
	)
}

func (c *Client) ChannelAdd(ctx context.Context, request ChannelAddRequest) (*schema.ChannelViews, *http.Response, error) {
	path := endpoints.ChannelsAdd

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelName"] = request.ChannelName
	body["country"] = *request.Country
	body["language"] = *request.Language
	body["category"] = *request.Category

	req, err := c.NewRestRequest(ctx, "POST", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelViews

	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func (c *Client) ChannelAvgPostsReach(ctx context.Context, request ChannelViewsRequest) (*schema.ChannelAvgReach, *http.Response, error) {
	path := endpoints.ChannelAVGPostsReach

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = request.ChannelId
	body["startDate"] = *request.StartDate
	body["endDate"] = *request.EndDate
	body["group"] = *request.Group

	req, err := c.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelAvgReach

	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func (c *Client) ChannelErr(ctx context.Context, request ChannelViewsRequest) (*schema.ChannelErr, *http.Response, error) {
	path := endpoints.ChannelErr

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = request.ChannelId
	body["startDate"] = *request.StartDate
	body["endDate"] = *request.EndDate
	body["group"] = *request.Group

	req, err := c.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelErr

	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}
