package channels

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
// see https://api.tgstat.ru/docs/ru/channels/get.html
func Get(ctx context.Context, channelId string) (*schema.ChannelResponse, *http.Response, error) {
	return getClient().Get(ctx, channelId)
}

func (c Client) Get(ctx context.Context, channelId string) (*schema.ChannelResponse, *http.Response, error) {
	path := endpoints.ChannelsGet

	if err := validateGetChannelId(channelId); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = channelId
	req, err := c.api.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelResponse
	result, err := c.api.Do(req, &response)
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
	Language            *string
	Category            *string
	Limit               *int
}

func (searchRequest SearchRequest) Validate() error {
	return validation.ValidateStruct(&searchRequest,
		validation.Field(&searchRequest.Country, validation.Required),
		validation.Field(&searchRequest.Q, validation.Required.When(*searchRequest.Category == "").Error("Either query or category is required.")),
		validation.Field(&searchRequest.Category, validation.Required.When(searchRequest.Q == "").Error("Either query or category  is required.")),
	)
}

// Search
func Search(ctx context.Context, request SearchRequest) (*schema.ChannelSearchResponse, *http.Response, error) {
	return getClient().Search(ctx, request)
}
func (c Client) Search(ctx context.Context, request SearchRequest) (*schema.ChannelSearchResponse, *http.Response, error) {
	path := endpoints.ChannelsSearch

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["q"] = request.Q
	body["search_by_description"] = strconv.Itoa(request.SearchByDescription)
	body["country"] = request.Country
	if nil != request.Language {
		body["language"] = *request.Language
	}
	if nil != request.Category {
		body["category"] = *request.Category
	}

	if nil != request.Limit {
		body["limit"] = strconv.Itoa(*request.Limit)
	}

	req, err := c.api.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelSearchResponse
	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func Stat(ctx context.Context, channelId string) (*schema.ChannelStatResponse, *http.Response, error) {
	return getClient().Stat(ctx, channelId)
}

func (c Client) Stat(ctx context.Context, channelId string) (*schema.ChannelStatResponse, *http.Response, error) {
	path := endpoints.ChannelsStat

	if err := validateGetChannelId(channelId); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = channelId
	req, err := c.api.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelStatResponse
	result, err := c.api.Do(req, &response)
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

func Posts(ctx context.Context, request PostsRequest) (*schema.ChannelPostsResponse, *http.Response, error) {
	return getClient().Posts(ctx, request)
}

func (c Client) Posts(ctx context.Context, request PostsRequest) (*schema.ChannelPostsResponse, *http.Response, error) {
	path := endpoints.ChannelsStat

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = request.ChannelId
	if nil != request.Limit {
		body["limit"] = strconv.FormatUint(*request.Limit, 10)
	}
	if nil != request.Offset {
		body["offset"] = strconv.FormatUint(*request.Offset, 10)
	}

	if nil != request.StartTime {
		body["startTime"] = *request.StartTime
	}

	if nil != request.EndTime {
		body["endTime"] = *request.EndTime
	}

	if nil != request.HideForwards {
		body["hideForwards"] = strconv.Itoa(*request.HideForwards)
	}

	if nil != request.HideDeleted {
		body["hideDeleted"] = strconv.Itoa(*request.HideDeleted)
	}

	body["extended"] = "0"

	req, err := c.api.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}
	var response schema.ChannelPostsResponse

	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func PostsExtended(ctx context.Context, request PostsRequest) (*schema.ChannelPostsWithChannelResponse, *http.Response, error) {
	return getClient().PostsExtended(ctx, request)
}

func (c Client) PostsExtended(ctx context.Context, request PostsRequest) (*schema.ChannelPostsWithChannelResponse, *http.Response, error) {
	path := endpoints.ChannelsStat

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = request.ChannelId

	if nil != request.Limit {
		body["limit"] = strconv.FormatUint(*request.Limit, 10)
	}
	if nil != request.Offset {
		body["offset"] = strconv.FormatUint(*request.Offset, 10)
	}

	if nil != request.StartTime {
		body["startTime"] = *request.StartTime
	}

	if nil != request.EndTime {
		body["endTime"] = *request.EndTime
	}

	if nil != request.HideForwards {
		body["hideForwards"] = strconv.Itoa(*request.HideForwards)
	}

	if nil != request.HideDeleted {
		body["hideDeleted"] = strconv.Itoa(*request.HideDeleted)
	}

	body["extended"] = "1"

	req, err := c.api.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelPostsWithChannelResponse

	result, err := c.api.Do(req, &response)
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
func Mentions(ctx context.Context, request ChannelMentionsRequest) (*schema.ChannelMentions, *http.Response, error) {
	return getClient().Mentions(ctx, request)
}

func (c Client) Mentions(ctx context.Context, request ChannelMentionsRequest) (*schema.ChannelMentions, *http.Response, error) {
	path := endpoints.ChannelsMentions

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = request.ChannelId

	if nil != request.Limit {
		body["limit"] = strconv.FormatUint(*request.Limit, 10)
	}

	if nil != request.Offset {
		body["offset"] = strconv.FormatUint(*request.Offset, 10)
	}

	if nil != request.StartDate {
		body["startDate"] = *request.StartDate
	}

	if nil != request.EndDate {
		body["endDate"] = *request.EndDate
	}

	body["extended"] = "0"

	req, err := c.api.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelMentions

	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func MentionsExtended(ctx context.Context, request ChannelMentionsRequest) (*schema.ChannelMentionsExtended, *http.Response, error) {
	return getClient().MentionsExtended(ctx, request)
}

func (c Client) MentionsExtended(ctx context.Context, request ChannelMentionsRequest) (*schema.ChannelMentionsExtended, *http.Response, error) {
	path := endpoints.ChannelsMentions

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = request.ChannelId

	if nil != request.Limit {
		body["limit"] = strconv.FormatUint(*request.Limit, 10)
	}

	if nil != request.Offset {
		body["offset"] = strconv.FormatUint(*request.Offset, 10)
	}

	if nil != request.StartDate {
		body["startDate"] = *request.StartDate
	}

	if nil != request.EndDate {
		body["endDate"] = *request.EndDate
	}

	body["extended"] = "1"

	req, err := c.api.NewRestRequest(ctx, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelMentionsExtended

	result, err := c.api.Do(req, &response)
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
	)
}

func Forwards(ctx context.Context, request ChannelForwardRequest) (*schema.ChannelForwards, *http.Response, error) {
	return getClient().Forwards(ctx, request)
}

func (c Client) Forwards(ctx context.Context, request ChannelForwardRequest) (*schema.ChannelForwards, *http.Response, error) {
	path := endpoints.ChannelsForwards

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = request.ChannelId

	if nil != request.Limit {
		body["limit"] = strconv.FormatUint(*request.Limit, 10)
	}

	if nil != request.Offset {
		body["offset"] = strconv.FormatUint(*request.Offset, 10)
	}

	if nil != request.StartDate {
		body["startDate"] = *request.StartDate
	}

	if nil != request.EndDate {
		body["endDate"] = *request.EndDate
	}

	body["extended"] = "0"

	req, err := c.api.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelForwards

	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func ForwardsExtended(ctx context.Context, request ChannelForwardRequest) (*schema.ChannelForwardsExtended, *http.Response, error) {
	return getClient().ForwardsExtended(ctx, request)
}

func (c Client) ForwardsExtended(ctx context.Context, request ChannelForwardRequest) (*schema.ChannelForwardsExtended, *http.Response, error) {
	path := endpoints.ChannelsForwards

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = request.ChannelId

	if nil != request.Limit {
		body["limit"] = strconv.FormatUint(*request.Limit, 10)
	}

	if nil != request.Offset {
		body["offset"] = strconv.FormatUint(*request.Offset, 10)
	}

	if nil != request.StartDate {
		body["startDate"] = *request.StartDate
	}

	if nil != request.EndDate {
		body["endDate"] = *request.EndDate
	}

	body["extended"] = "1"

	req, err := c.api.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelForwardsExtended

	result, err := c.api.Do(req, &response)
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
		validation.Field(&channelSubscribersRequest.Group, validation.In("hour", "day", "week", "month")),
	)
}

func Subscribers(ctx context.Context, request ChannelSubscribersRequest) (*schema.ChannelSubscribers, *http.Response, error) {
	return getClient().Subscribers(ctx, request)
}

func (c Client) Subscribers(ctx context.Context, request ChannelSubscribersRequest) (*schema.ChannelSubscribers, *http.Response, error) {
	path := endpoints.ChannelsSubscribers

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = request.ChannelId
	if nil != request.StartDate {
		body["startDate"] = *request.StartDate
	}

	if nil != request.EndDate {
		body["endDate"] = *request.EndDate
	}

	if nil != request.Group {
		body["group"] = *request.Group
	}

	req, err := c.api.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelSubscribers

	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
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
		validation.Field(&channelViewsRequest.Group, validation.In("day", "week", "month")),
	)
}

func Views(ctx context.Context, request ChannelViewsRequest) (*schema.ChannelViews, *http.Response, error) {
	return getClient().Views(ctx, request)
}

func (c Client) Views(ctx context.Context, request ChannelViewsRequest) (*schema.ChannelViews, *http.Response, error) {
	path := endpoints.ChannelsViews

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = request.ChannelId
	if nil != request.StartDate {
		body["startDate"] = *request.StartDate
	}

	if nil != request.EndDate {
		body["endDate"] = *request.EndDate
	}

	if nil != request.Group {
		body["group"] = *request.Group
	}

	req, err := c.api.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelViews

	result, err := c.api.Do(req, &response)
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

func Add(ctx context.Context, request ChannelAddRequest) (*schema.ChannelViews, *http.Response, error) {
	return getClient().Add(ctx, request)
}

func (c Client) Add(ctx context.Context, request ChannelAddRequest) (*schema.ChannelViews, *http.Response, error) {
	path := endpoints.ChannelsAdd

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelName"] = request.ChannelName

	if nil != request.Country {
		body["country"] = *request.Country
	}

	if nil != request.Language {
		body["language"] = *request.Language
	}

	if nil != request.Category {
		body["category"] = *request.Category
	}

	req, err := c.api.NewRestRequest(ctx, "POST", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelViews

	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func AvgPostsReach(ctx context.Context, request ChannelViewsRequest) (*schema.ChannelAvgReach, *http.Response, error) {
	return getClient().AvgPostsReach(ctx, request)
}

func (c Client) AvgPostsReach(ctx context.Context, request ChannelViewsRequest) (*schema.ChannelAvgReach, *http.Response, error) {
	path := endpoints.ChannelAVGPostsReach

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = request.ChannelId
	if nil != request.StartDate {
		body["startDate"] = *request.StartDate
	}

	if nil != request.EndDate {
		body["endDate"] = *request.EndDate
	}

	if nil != request.Group {
		body["group"] = *request.Group
	}

	req, err := c.api.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelAvgReach

	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func Err(ctx context.Context, request ChannelViewsRequest) (*schema.ChannelErr, *http.Response, error) {
	return getClient().Err(ctx, request)
}

func (c Client) Err(ctx context.Context, request ChannelViewsRequest) (*schema.ChannelErr, *http.Response, error) {
	path := endpoints.ChannelErr

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = request.ChannelId
	if nil != request.StartDate {
		body["startDate"] = *request.StartDate
	}

	if nil != request.EndDate {
		body["endDate"] = *request.EndDate
	}

	if nil != request.Group {
		body["group"] = *request.Group
	}

	req, err := c.api.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelErr

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
