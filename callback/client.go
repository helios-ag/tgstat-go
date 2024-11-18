package callback

import (
	"context"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/endpoints"
	"net/http"
	"net/url"
)

type Client struct {
	api   tgstat.API
	token string
}

// SetCallback request
// https://api.tgstat.ru/docs/ru/callback/set-callback-url.html
func SetCallback(ctx context.Context, callbackUrl string) (*tgstat.SetCallbackVerificationResult, *http.Response, error) {
	return getClient().SetCallback(ctx, callbackUrl)
}
func (c Client) SetCallback(ctx context.Context, callbackUrl string) (*tgstat.SetCallbackVerificationResult, *http.Response, error) {
	path := endpoints.SetCallbackURL

	if err := validateCallbackUrl(callbackUrl); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["callback_url"] = callbackUrl

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodPost, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response tgstat.SetCallbackVerificationResult

	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func validateCallbackUrl(callbackUrl string) error {
	if callbackUrl == "" {
		return fmt.Errorf("CallbackUrl must be set")
	}

	if _, err := url.ParseRequestURI(callbackUrl); err != nil {
		return fmt.Errorf("unable to parse URL: %v", err)
	}

	return nil
}

// GetCallbackInfo request
// https://api.tgstat.ru/docs/ru/callback/get-callback-url.html
func GetCallbackInfo(ctx context.Context) (*tgstat.GetCallbackResponse, *http.Response, error) {
	return getClient().GetCallbackInfo(ctx)
}
func (c Client) GetCallbackInfo(ctx context.Context) (*tgstat.GetCallbackResponse, *http.Response, error) {
	path := endpoints.GetCallbackURL
	body := make(map[string]string)
	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response tgstat.GetCallbackResponse
	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

type SubscribeChannelRequest struct {
	SubscriptionId *string
	ChannelId      string
	EventTypes     string
}

func (subscribeChannelRequest SubscribeChannelRequest) Validate() error {
	return validation.ValidateStruct(&subscribeChannelRequest,
		validation.Field(&subscribeChannelRequest.ChannelId, validation.Required),
		validation.Field(&subscribeChannelRequest.EventTypes, validation.Required, validation.In("new_post", "edit_post", "remove_post")),
	)
}

// SubscribeChannel request
// https://api.tgstat.ru/docs/ru/callback/subscribe-channel.html
func SubscribeChannel(ctx context.Context, request SubscribeChannelRequest) (*tgstat.SubscribeResponse, *http.Response, error) {
	return getClient().SubscribeChannel(ctx, request)
}
func (c Client) SubscribeChannel(ctx context.Context, request SubscribeChannelRequest) (*tgstat.SubscribeResponse, *http.Response, error) {
	path := endpoints.SubscribeChannel

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	if nil != request.SubscriptionId {
		body["subscription_id"] = *request.SubscriptionId
	}

	body["channel_id"] = request.ChannelId
	body["event_types"] = request.EventTypes

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodPost, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response tgstat.SubscribeResponse

	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

type SubscribeWordRequest struct {
	SubscriptionId *string
	Q              string
	EventTypes     string
	StrongSearch   *bool
	MinusWords     *string
	ExtendedSyntax *bool
	PeerTypes      *string
}

func (subscribeWordRequest SubscribeWordRequest) Validate() error {
	return validation.ValidateStruct(&subscribeWordRequest,
		validation.Field(&subscribeWordRequest.Q, validation.Required),
		validation.Field(&subscribeWordRequest.EventTypes, validation.Required, validation.In("new_post")),
		validation.Field(&subscribeWordRequest.PeerTypes, validation.In("channel", "chat", "all")),
	)
}

// SubscribeWord request
// https://api.tgstat.ru/docs/ru/callback/subscribe-word.html
func SubscribeWord(ctx context.Context, request SubscribeWordRequest) (*tgstat.Subscribe, *http.Response, error) {
	return getClient().SubscribeWord(ctx, request)
}

func (c Client) SubscribeWord(ctx context.Context, request SubscribeWordRequest) (*tgstat.Subscribe, *http.Response, error) {
	path := endpoints.SubscribeWord

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	if nil != request.SubscriptionId {
		body["subscription_id"] = *request.SubscriptionId
	}

	body["q"] = request.Q
	body["event_types"] = request.EventTypes

	if nil != request.StrongSearch {
		body["strong_search"] = boolValue(request.StrongSearch)
	}

	if nil != request.MinusWords {
		body["minus_words"] = *request.MinusWords
	}

	if nil != request.ExtendedSyntax {
		body["extended_syntax"] = boolValue(request.ExtendedSyntax)
	}

	if nil != request.PeerTypes {
		body["peer_types"] = *request.PeerTypes
	}

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodPost, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response tgstat.Subscribe

	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

type SubscriptionsListRequest struct {
	SubscriptionId   *string
	SubscriptionType *string
}

func (subscriptionsListRequest SubscriptionsListRequest) Validate() error {
	return validation.ValidateStruct(&subscriptionsListRequest,
		validation.Field(&subscriptionsListRequest.SubscriptionType, validation.In("channel", "keyword")),
	)
}

// SubscriptionsList request
// https://api.tgstat.ru/docs/ru/callback/get-callback-url.html
func SubscriptionsList(ctx context.Context, subscriptionsListRequest SubscriptionsListRequest) (*tgstat.SubscriptionList, *http.Response, error) {
	return getClient().SubscriptionsList(ctx, subscriptionsListRequest)
}
func (c Client) SubscriptionsList(ctx context.Context, subscriptionsListRequest SubscriptionsListRequest) (*tgstat.SubscriptionList, *http.Response, error) {
	path := endpoints.SubscriptionsList
	body := make(map[string]string)
	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	if nil != subscriptionsListRequest.SubscriptionId {
		body["subscription_id"] = *subscriptionsListRequest.SubscriptionId
	}

	if nil != subscriptionsListRequest.SubscriptionType {
		body["subscription_type"] = *subscriptionsListRequest.SubscriptionType
	}

	var response tgstat.SubscriptionList
	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

// Unsubscribe request
// https://api.tgstat.ru/docs/ru/callback/unsubscribe.html
func Unsubscribe(ctx context.Context, subscriptionId string) (*tgstat.SuccessResult, *http.Response, error) {
	return getClient().Unsubscribe(ctx, subscriptionId)
}

func (c Client) Unsubscribe(ctx context.Context, subscriptionId string) (*tgstat.SuccessResult, *http.Response, error) {
	path := endpoints.Unsubscribe

	if subscriptionId == "" {
		return nil, nil, fmt.Errorf("subscriptionId must be set")
	}

	body := make(map[string]string)
	
	
	body["subscription_id"] = subscriptionId
	
	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodPost, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response tgstat.SuccessResult

	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func boolValue(v *bool) string {
	if *v {
		return "0"
	}
	return "1"
}

func getClient() Client {
	return Client{tgstat.GetAPI(), tgstat.Token}
}
