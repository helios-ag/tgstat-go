package tgstat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	//"net/url"
	"strconv"
	"tgstat/endpoints"
	//"github.com/helios-ag/tgstat-go/endpoints"
	//"github.com/helios-ag/tgstat-go/schema"
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
		return fmt.Errorf("channel ID must be set")
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

func (c *Client) ChannelSearch(ctx context.Context, request SearchRequest) (*schema.ChannelSearchResponse, *http.Response, error) {
	path := endpoints.ChannelsSearch

	if err := validateSearchRequest(request); err != nil {
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

func validateSearchRequest(request SearchRequest) error {
	if request.Country == "" {
		return fmt.Errorf("country must be set")
	}

	if request.Q != "" {
		if request.Q == "" || request.Category == "" {
			return fmt.Errorf("category or query word is empty")
		}
	}

	return nil
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
	ChannelId           string
	Limit               *uint64
	Offset              *uint64
	StartTime			*string
	EndTime             *string
	HideForwards        *int
	HideDeleted         *int
}

func (c *Client) ChannelPosts(ctx context.Context, request PostsRequest) (*schema.ChannelPostsResponse, *http.Response, error) {
	path := endpoints.ChannelsStat

	if err := validatePostsRequests(request); err != nil {
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

	if err := validatePostsRequests(request); err != nil {
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

func validatePostsRequests(request PostsRequest) error {
	if request.ChannelId == "" {
		return fmt.Errorf("channel ID must be set")
	}

	if *request.Limit > 50 {
		return fmt.Errorf("max limit is 50")
	}

	if *request.Offset > 1000 {
		return fmt.Errorf("max offset is 1000")
	}

	return nil
}

type ChannelMentionsRequest struct {
	ChannelId           string
	Limit               *uint64
	Offset              *uint64
	StartDate			*string
	EndDate             *string
}

func (c *Client) ChannelMentions(ctx context.Context, request ChannelMentionsRequest) (*schema.ChannelMentions, *http.Response, error) {
	path := endpoints.ChannelsMentions

	if err := validateChannelMentionsRequests(request); err != nil {
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

	if err := validateChannelMentionsRequests(request); err != nil {
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

func validateChannelMentionsRequests(request ChannelMentionsRequest) error {
	if request.ChannelId == "" {
		return fmt.Errorf("channel ID must be set")
	}

	if *request.Limit > 50 {
		return fmt.Errorf("max limit is 50")
	}

	if *request.Offset > 1000 {
		return fmt.Errorf("max offset is 1000")
	}

	return nil
}

type ChannelForwardRequest struct {
	ChannelId           string
	Limit               *uint64
	Offset              *uint64
	StartDate			*string
	EndDate             *string
}

func (c *Client) ChannelForwards(ctx context.Context, request ChannelForwardRequest) (*schema.ChannelForwards, *http.Response, error) {
	path := endpoints.ChannelsMentions

	if err := validateChannelForwardRequest(request); err != nil {
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
	path := endpoints.ChannelsMentions

	if err := validateChannelForwardRequest(request); err != nil {
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

func validateChannelForwardRequest(request ChannelForwardRequest) error {
	if request.ChannelId == "" {
		return fmt.Errorf("channel ID must be set")
	}

	if *request.Limit > 50 {
		return fmt.Errorf("max limit is 50")
	}

	return nil
}

type ChannelSubscribersRequest struct {
	ChannelId           string
	StartDate			*string
	EndDate             *string
	Group               *string
}

func (c *Client) ChannelSubscribers(ctx context.Context, request ChannelSubscribersRequest) (*schema.ChannelSubscribers, *http.Response, error) {
	path := endpoints.ChannelsSubscribers

	if err := validateChannelSubscriberRequest(request); err != nil {
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

func validateChannelSubscriberRequest(request ChannelSubscribersRequest) error {
	if request.ChannelId == "" {
		return fmt.Errorf("channel ID must be set")
	}

	if !validateGroup(*request.Group) {
		return fmt.Errorf("improper group value")
	}

	return nil
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
	ChannelId           string
	StartDate			*string
	EndDate             *string
	Group               *string
}

func (c *Client) ChannelViews(ctx context.Context, request ChannelViewsRequest) (*schema.ChannelViews, *http.Response, error) {
	path := endpoints.ChannelsSubscribers

	if err := validateChannelViewsRequest(request); err != nil {
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

func validateChannelViewsRequest(request ChannelViewsRequest) error {
	if request.ChannelId == "" {
		return fmt.Errorf("channel ID must be set")
	}

	if !validateGroup(*request.Group) {
		return fmt.Errorf("improper group value")
	}

	return nil
}

type ChannelAddRequest struct {
	ChannelName   string
	Country	      *string
	Language      *string
	Category      *string
}

func (c *Client) ChannelAdd(ctx context.Context, request ChannelAddRequest) (*schema.ChannelViews, *http.Response, error) {
	path := endpoints.ChannelsSubscribers

	if err := validateChannelAddRequest(request); err != nil {
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

func validateChannelAddRequest(request ChannelAddRequest) error {
	if request.ChannelName == "" {
		return fmt.Errorf("channel name must be set")
	}

	return nil
}
