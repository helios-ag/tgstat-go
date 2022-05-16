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
	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.ChannelResponse
	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	switch x := response.Response.TGStatRestriction.(type) {
	case []interface{}:
		response.Response.TGStatRestriction = nil
	case map[string]interface{}:
		var restrictions schema.TGStatRestrictions
		jsonString, _ := json.Marshal(x)
		errors := json.Unmarshal(jsonString, &restrictions)
		if errors == nil {
			response.Response.TGStatRestriction = restrictions
		}
	default:
		return nil, result, fmt.Errorf("something wrong with Restrictions response")
	}

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
	Category            string
	Limit               *int
}

func (searchRequest SearchRequest) Validate() error {
	return validation.ValidateStruct(&searchRequest,
		validation.Field(&searchRequest.Country, validation.Required),
		validation.Field(&searchRequest.Q, validation.Required.When(searchRequest.Category == "").Error("Either query or category is required.")),
		validation.Field(&searchRequest.Category, validation.Required.When(searchRequest.Q == "").Error("Either query or category  is required.")),
	)
}

// Search request
// see https://api.tgstat.ru/docs/ru/channels/search.html
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
	if request.Category != "" {
		body["category"] = request.Category
	}

	if nil != request.Limit {
		body["limit"] = strconv.Itoa(*request.Limit)
	}

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

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

// Stat request
// See https://api.tgstat.ru/docs/ru/channels/stat.html
func Stat(ctx context.Context, channelId string) (*schema.ChannelStatResponse, *http.Response, error) {
	return getClient().Stat(ctx, channelId)
}

// Stat request
// See https://api.tgstat.ru/docs/ru/channels/stat.html
func (c Client) Stat(ctx context.Context, channelId string) (*schema.ChannelStatResponse, *http.Response, error) {
	path := endpoints.ChannelsStat

	if err := validateGetChannelId(channelId); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = channelId
	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

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
	HideForwards *bool
	HideDeleted  *bool
}

func (postsRequest PostsRequest) Validate() error {
	return validation.ValidateStruct(&postsRequest,
		validation.Field(&postsRequest.ChannelId, validation.Required),
		validation.Field(&postsRequest.StartTime, validation.Date("1643113399")),
		validation.Field(&postsRequest.EndTime, validation.Date("1643113399")),
	)
}

// Posts request
// see https://api.tgstat.ru/docs/ru/channels/posts.html
func Posts(ctx context.Context, request PostsRequest) (*schema.ChannelPosts, *http.Response, error) {
	return getClient().Posts(ctx, request)
}

// Posts request
// see https://api.tgstat.ru/docs/ru/channels/posts.html
func (c Client) Posts(ctx context.Context, request PostsRequest) (*schema.ChannelPosts, *http.Response, error) {
	path := endpoints.ChannelsPosts

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
		body["hideForwards"] = strconv.FormatBool(*request.HideForwards)
	}

	if nil != request.HideDeleted {
		body["hideDeleted"] = strconv.FormatBool(*request.HideDeleted)
	}

	body["extended"] = "0"

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}
	var response schema.ChannelPosts

	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

// PostsExtended request extended
// see https://api.tgstat.ru/docs/ru/channels/posts.html
func PostsExtended(ctx context.Context, request PostsRequest) (*schema.ChannelPostsWithChannelResponse, *http.Response, error) {
	return getClient().PostsExtended(ctx, request)
}

// PostsExtended request extended
// see https://api.tgstat.ru/docs/ru/channels/posts.html
func (c Client) PostsExtended(ctx context.Context, request PostsRequest) (*schema.ChannelPostsWithChannelResponse, *http.Response, error) {
	path := endpoints.ChannelsPosts

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
		body["hideForwards"] = strconv.FormatBool(*request.HideForwards)
	}

	if nil != request.HideDeleted {
		body["hideDeleted"] = strconv.FormatBool(*request.HideDeleted)
	}

	body["extended"] = "1"

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

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

// Mentions request
// see https://api.tgstat.ru/docs/ru/channels/mentions.html
func Mentions(ctx context.Context, request ChannelForwardRequest) (*schema.ChannelMentions, *http.Response, error) {
	return getClient().Mentions(ctx, request)
}

// Mentions request
// see https://api.tgstat.ru/docs/ru/channels/mentions.html
func (c Client) Mentions(ctx context.Context, request ChannelForwardRequest) (*schema.ChannelMentions, *http.Response, error) {
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

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

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

// MentionsExtended request
// see https://api.tgstat.ru/docs/ru/channels/mentions.html
func MentionsExtended(ctx context.Context, request ChannelForwardRequest) (*schema.ChannelMentionsExtended, *http.Response, error) {
	return getClient().MentionsExtended(ctx, request)
}

// MentionsExtended Mentions request
// see https://api.tgstat.ru/docs/ru/channels/mentions.html
func (c Client) MentionsExtended(ctx context.Context, request ChannelForwardRequest) (*schema.ChannelMentionsExtended, *http.Response, error) {
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

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

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
	)
}

// Forwards request
// see https://api.tgstat.ru/docs/ru/channels/forwards.html
func Forwards(ctx context.Context, request ChannelForwardRequest) (*schema.ChannelForwards, *http.Response, error) {
	return getClient().Forwards(ctx, request)
}

// Forwards request
// see https://api.tgstat.ru/docs/ru/channels/forwards.html
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

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

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

// ForwardsExtended Forwards request extended
// see https://api.tgstat.ru/docs/ru/channels/forwards.html
func ForwardsExtended(ctx context.Context, request ChannelForwardRequest) (*schema.ChannelForwardsExtended, *http.Response, error) {
	return getClient().ForwardsExtended(ctx, request)
}

// ForwardsExtended request extended
// see https://api.tgstat.ru/docs/ru/channels/forwards.html
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

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

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

// Subscribers request
// see https://api.tgstat.ru/docs/ru/channels/subscribers.html
func Subscribers(ctx context.Context, request ChannelSubscribersRequest) (*schema.ChannelSubscribers, *http.Response, error) {
	return getClient().Subscribers(ctx, request)
}

// Subscribers request
// see https://api.tgstat.ru/docs/ru/channels/subscribers.html
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

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

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

// Views request
// see https://api.tgstat.ru/docs/ru/channels/views.html
func Views(ctx context.Context, request ChannelViewsRequest) (*schema.ChannelViews, *http.Response, error) {
	return getClient().Views(ctx, request)
}

// Views request
// see https://api.tgstat.ru/docs/ru/channels/views.html
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

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

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

// Add request
// see https://api.tgstat.ru/docs/ru/channels/add.html
func Add(ctx context.Context, request ChannelAddRequest) (*schema.ChannelViews, *http.Response, error) {
	return getClient().Add(ctx, request)
}

// Add request
// see https://api.tgstat.ru/docs/ru/channels/add.html
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

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodPost, path, body)

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

// AvgPostsReach request
// See https://api.tgstat.ru/docs/ru/channels/avg-posts-reach.html
func AvgPostsReach(ctx context.Context, request ChannelViewsRequest) (*schema.ChannelAvgReach, *http.Response, error) {
	return getClient().AvgPostsReach(ctx, request)
}

// AvgPostsReach request
// See https://api.tgstat.ru/docs/ru/channels/avg-posts-reach.html
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

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

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

// Err request
// See https://api.tgstat.ru/docs/ru/channels/err.html
func Err(ctx context.Context, request ChannelViewsRequest) (*schema.ChannelErr, *http.Response, error) {
	return getClient().Err(ctx, request)
}

// Err request
// See https://api.tgstat.ru/docs/ru/channels/err.html
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

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

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
