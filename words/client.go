package words

import (
	"context"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/endpoints"
	"github.com/helios-ag/tgstat-go/schema"
	"net/http"
)

type Client struct {
	api   tgstat.API
	token string
}

type MentionPeriodRequest struct {
	Q              string
	PeerType       *string
	StartDate      *string
	EndDate        *string
	HideForwards   *bool
	StrongSearch   *bool
	MinusWords     *string
	Group          *string
	ExtendedSyntax *bool
}

func (mentionPeriodRequest MentionPeriodRequest) Validate() error {
	return validation.ValidateStruct(&mentionPeriodRequest,
		validation.Field(&mentionPeriodRequest.Q, validation.Required),
		validation.Field(&mentionPeriodRequest.Group, validation.In("channel", "chat", "all")),
	)
}

// MentionsByPeriod request
// See https://api.tgstat.ru/docs/ru/words/mentions-by-period.html
func MentionsByPeriod(ctx context.Context, request MentionPeriodRequest) (*schema.WordsMentionsResponse, *http.Response, error) {
	return getClient().MentionsByPeriod(ctx, request)
}

// MentionsByPeriod request
// See https://api.tgstat.ru/docs/ru/words/mentions-by-period.html
func (c Client) MentionsByPeriod(ctx context.Context, request MentionPeriodRequest) (*schema.WordsMentionsResponse, *http.Response, error) {
	path := endpoints.WordsMentionsByPeriod

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["q"] = request.Q
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
	body["strongSearch"] = func() string {
		if nil != request.StrongSearch && *request.StrongSearch {
			return "1"
		} else {
			return "0"
		}
	}()

	if nil != request.MinusWords {
		body["minusWords"] = *request.MinusWords
	}

	if nil != request.Group {
		body["group"] = *request.Group
	}

	body["extendedSyntax"] = func() string {
		if nil != request.ExtendedSyntax && *request.ExtendedSyntax {
			return "1"
		} else {
			return "0"
		}
	}()
	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.WordsMentionsResponse
	result, err := c.api.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

type MentionsByChannelRequest struct {
	Q              string
	PeerType       *string
	StartDate      *string
	EndDate        *string
	HideForwards   *bool
	StrongSearch   *bool
	MinusWords     *string
	ExtendedSyntax *bool
}

func (mentionsByChannelRequest MentionsByChannelRequest) Validate() error {
	return validation.ValidateStruct(&mentionsByChannelRequest,
		validation.Field(&mentionsByChannelRequest.Q, validation.Required),
		validation.Field(&mentionsByChannelRequest.PeerType, validation.In("channel", "chat", "all")),
	)
}

// MentionsByChannels request
// See https://api.tgstat.ru/docs/ru/words/mentions-by-channels.html
func MentionsByChannels(ctx context.Context, request MentionsByChannelRequest) (*schema.WordsMentionsResponse, *http.Response, error) {
	return getClient().MentionsByChannels(ctx, request)
}

// MentionsByChannels request
// See https://api.tgstat.ru/docs/ru/words/mentions-by-channels.html
func (c Client) MentionsByChannels(ctx context.Context, request MentionsByChannelRequest) (*schema.WordsMentionsResponse, *http.Response, error) {
	path := endpoints.WordsMentionsByChannels

	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["q"] = request.Q
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
	body["strongSearch"] = func() string {
		if nil != request.StrongSearch && *request.StrongSearch {
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

	req, err := c.api.NewRestRequest(ctx, c.token, http.MethodGet, path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.WordsMentionsResponse
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
