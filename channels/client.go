package channels

import (
	"context"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/endpoints"
	"net/http"
	"strconv"
)

type Client struct {
	api   tgstat.API
	token string
}

// Get request
// see https://api.tgstat.ru/docs/ru/channels/get.html
func Get(ctx context.Context, channelId string) (*tgstat.ChannelResponseResult, *http.Response, error) {
	return getClient().Get(ctx, channelId)
}

func (c Client) Get(ctx context.Context, channelId string) (*tgstat.ChannelResponseResult, *http.Response, error) {
	path := endpoints.ChannelsGet

	if err := validateChannelId(channelId); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = channelId
	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response tgstat.ChannelResponseResult
	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	switch x := response.Response.TGStatRestriction.(type) {
	case []interface{}:
		response.Response.TGStatRestriction = nil
	case map[string]interface{}:
		var restrictions tgstat.TGStatRestrictions
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

func validateChannelId(channelId string) error {
	if channelId == "" {
		return fmt.Errorf("ChannelId: cannot be blank")
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
func Search(ctx context.Context, request SearchRequest) (*tgstat.ChannelSearchResult, *http.Response, error) {
	return getClient().Search(ctx, request)
}
func (c Client) Search(ctx context.Context, request SearchRequest) (*tgstat.ChannelSearchResult, *http.Response, error) {
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

	var response tgstat.ChannelSearchResult
	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

// Stat request
// See https://api.tgstat.ru/docs/ru/channels/stat.html
func Stat(ctx context.Context, channelId string) (*tgstat.ChannelStatResult, *http.Response, error) {
	return getClient().Stat(ctx, channelId)
}

// Stat request
// See https://api.tgstat.ru/docs/ru/channels/stat.html
func (c Client) Stat(ctx context.Context, channelId string) (*tgstat.ChannelStatResult, *http.Response, error) {
	path := endpoints.ChannelsStat

	if err := validateChannelId(channelId); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["channelId"] = channelId
	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response tgstat.ChannelStatResult
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

// Posts request
// see https://api.tgstat.ru/docs/ru/channels/posts.html
func Posts(ctx context.Context, request PostsRequest) (*tgstat.ChannelPostsResult, *http.Response, error) {
	return getClient().Posts(ctx, request)
}

// Posts request
// see https://api.tgstat.ru/docs/ru/channels/posts.html
func (c Client) Posts(ctx context.Context, request PostsRequest) (*tgstat.ChannelPostsResult, *http.Response, error) {
	path := endpoints.ChannelsPosts

	if err := validateChannelId(request.ChannelId); err != nil {
		return nil, nil, err
	}

	body := posts(request, false)

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}
	var response tgstat.ChannelPostsResult

	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

// PostsExtended request extended
// see https://api.tgstat.ru/docs/ru/channels/posts.html
func PostsExtended(ctx context.Context, request PostsRequest) (*tgstat.ChannelPostsWithChannelResult, *http.Response, error) {
	return getClient().PostsExtended(ctx, request)
}

// PostsExtended request extended
// see https://api.tgstat.ru/docs/ru/channels/posts.html
func (c Client) PostsExtended(ctx context.Context, request PostsRequest) (*tgstat.ChannelPostsWithChannelResult, *http.Response, error) {
	path := endpoints.ChannelsPosts

	if err := validateChannelId(request.ChannelId); err != nil {
		return nil, nil, err
	}

	body := posts(request, true)

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response tgstat.ChannelPostsWithChannelResult

	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func posts(request PostsRequest, extended bool) map[string]string {

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

	body["extended"] = strconv.FormatBool(extended)

	return body
}

// Mentions request
// see https://api.tgstat.ru/docs/ru/channels/mentions.html
func Mentions(ctx context.Context, request ChannelForwardRequest) (*tgstat.ChannelMentionsResult, *http.Response, error) {
	return getClient().Mentions(ctx, request)
}

// Mentions request
// see https://api.tgstat.ru/docs/ru/channels/mentions.html
func (c Client) Mentions(ctx context.Context, request ChannelForwardRequest) (*tgstat.ChannelMentionsResult, *http.Response, error) {
	path := endpoints.ChannelsMentions

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := mentions(request, false)

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response tgstat.ChannelMentionsResult

	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

// MentionsExtended request
// see https://api.tgstat.ru/docs/ru/channels/mentions.html
func MentionsExtended(ctx context.Context, request ChannelForwardRequest) (*tgstat.ChannelMentionsExtended, *http.Response, error) {
	return getClient().MentionsExtended(ctx, request)
}

// MentionsExtended Mentions request
// see https://api.tgstat.ru/docs/ru/channels/mentions.html
func (c Client) MentionsExtended(ctx context.Context, request ChannelForwardRequest) (*tgstat.ChannelMentionsExtended, *http.Response, error) {
	path := endpoints.ChannelsMentions

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := mentions(request, true)

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response tgstat.ChannelMentionsExtended

	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func mentions(request ChannelForwardRequest, extended bool) map[string]string {
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

	body["extended"] = strconv.FormatBool(extended)

	return body
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
func Forwards(ctx context.Context, request ChannelForwardRequest) (*tgstat.ChannelForwards, *http.Response, error) {
	return getClient().Forwards(ctx, request)
}

// Forwards request
// see https://api.tgstat.ru/docs/ru/channels/forwards.html
func (c Client) Forwards(ctx context.Context, request ChannelForwardRequest) (*tgstat.ChannelForwards, *http.Response, error) {
	path := endpoints.ChannelsForwards

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := forwards(request, false)

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response tgstat.ChannelForwards

	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

// ForwardsExtended Forwards request extended
// see https://api.tgstat.ru/docs/ru/channels/forwards.html
func ForwardsExtended(ctx context.Context, request ChannelForwardRequest) (*tgstat.ChannelForwardsExtended, *http.Response, error) {
	return getClient().ForwardsExtended(ctx, request)
}

// ForwardsExtended request extended
// see https://api.tgstat.ru/docs/ru/channels/forwards.html
func (c Client) ForwardsExtended(ctx context.Context, request ChannelForwardRequest) (*tgstat.ChannelForwardsExtended, *http.Response, error) {
	path := endpoints.ChannelsForwards

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := forwards(request, true)

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response tgstat.ChannelForwardsExtended

	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func forwards(request ChannelForwardRequest, extended bool) map[string]string {
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

	if extended {
		body["extended"] = "1"
	}

	return body
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
func Subscribers(ctx context.Context, request ChannelSubscribersRequest) (*tgstat.ChannelSubscribers, *http.Response, error) {
	return getClient().Subscribers(ctx, request)
}

// Subscribers request
// see https://api.tgstat.ru/docs/ru/channels/subscribers.html
func (c Client) Subscribers(ctx context.Context, request ChannelSubscribersRequest) (*tgstat.ChannelSubscribers, *http.Response, error) {
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

	var response tgstat.ChannelSubscribers

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
func Views(ctx context.Context, request ChannelViewsRequest) (*tgstat.ChannelViews, *http.Response, error) {
	return getClient().Views(ctx, request)
}

// Views request
// see https://api.tgstat.ru/docs/ru/channels/views.html
func (c Client) Views(ctx context.Context, request ChannelViewsRequest) (*tgstat.ChannelViews, *http.Response, error) {
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

	var response tgstat.ChannelViews

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
func Add(ctx context.Context, request ChannelAddRequest) (*tgstat.ChannelAddSuccess, *http.Response, error) {
	return getClient().Add(ctx, request)
}

// Add request
// see https://api.tgstat.ru/docs/ru/channels/add.html
func (c Client) Add(ctx context.Context, request ChannelAddRequest) (*tgstat.ChannelAddSuccess, *http.Response, error) {
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

	var response tgstat.ChannelAddSuccess

	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

// AvgPostsReach request
// See https://api.tgstat.ru/docs/ru/channels/avg-posts-reach.html
func AvgPostsReach(ctx context.Context, request ChannelViewsRequest) (*tgstat.ChannelAvgReach, *http.Response, error) {
	return getClient().AvgPostsReach(ctx, request)
}

// AvgPostsReach request
// See https://api.tgstat.ru/docs/ru/channels/avg-posts-reach.html
func (c Client) AvgPostsReach(ctx context.Context, request ChannelViewsRequest) (*tgstat.ChannelAvgReach, *http.Response, error) {
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

	var response tgstat.ChannelAvgReach

	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

// Err request
// See https://api.tgstat.ru/docs/ru/channels/err.html
func Err(ctx context.Context, request ChannelViewsRequest) (*tgstat.ChannelErr, *http.Response, error) {
	return getClient().Err(ctx, request)
}

// Err request
// See https://api.tgstat.ru/docs/ru/channels/err.html
func (c Client) Err(ctx context.Context, request ChannelViewsRequest) (*tgstat.ChannelErr, *http.Response, error) {
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

	var response tgstat.ChannelErr

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
