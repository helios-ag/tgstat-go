package callback

import (
	"context"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/endpoints"
	"github.com/helios-ag/tgstat-go/schema"
	"net/http"
)

type Client struct {
	api   tgstat.API
	token string
}

type SetCallbackRequest struct {
	CallbackUrl string
}

func (setCallbackRequest SetCallbackRequest) Validate() error {
	return validation.ValidateStruct(&setCallbackRequest,
		validation.Field(&setCallbackRequest.CallbackUrl, validation.Required, is.URL),
	)
}

// SetCallback request
// https://api.tgstat.ru/docs/ru/callback/set-callback-url.html
func SetCallback(ctx context.Context, request SetCallbackRequest) (*schema.SetCallbackResponse, *http.Response, error) {
	return getClient().SetCallback(ctx, request)
}
func (c Client) SetCallback(ctx context.Context, request SetCallbackRequest) (*schema.SetCallbackResponse, *http.Response, error) {
	path := endpoints.SetCallbackURL

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["callback_url"] = request.CallbackUrl

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodPost, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.SetCallbackResponse

	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

// GetCallbackInfo GetCallback request
// https://api.tgstat.ru/docs/ru/callback/get-callback-url.html
func GetCallbackInfo(ctx context.Context) (*schema.GetCallbackResponse, *http.Response, error) {
	return getClient().GetCallbackInfo(ctx)
}
func (c Client) GetCallbackInfo(ctx context.Context) (*schema.GetCallbackResponse, *http.Response, error) {
	path := endpoints.GetCallbackURL
	body := make(map[string]string)
	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.GetCallbackResponse
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
		validation.Field(&subscribeChannelRequest.EventTypes, validation.Required),
	)
}

// SubscribeChannel SetCallback request
// https://api.tgstat.ru/docs/ru/callback/subscribe-channel.html
func SubscribeChannel(ctx context.Context, request SubscribeChannelRequest) (*schema.SubscribeResponse, *http.Response, error) {
	return getClient().SubscribeChannel(ctx, request)
}
func (c Client) SubscribeChannel(ctx context.Context, request SubscribeChannelRequest) (*schema.SubscribeResponse, *http.Response, error) {
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

	var response schema.SubscribeResponse

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
