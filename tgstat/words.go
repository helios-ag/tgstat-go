package tgstat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"tgstat/endpoints"
	"tgstat/schema"
)

type MentionPeriodRequest struct {
	Q string
	PeerType *string
	StartDate *string
	EndDate *string
	HideForwards *bool
	StrongSearch *bool
	MinusWords *string
	Group *string
	ExtendedSyntax *bool
}

func (c *Client) ChannelMentionsByPeriod(ctx context.Context, request MentionPeriodRequest) (*schema.WordsMentionsResponse, *http.Response, error) {
	path := endpoints.WordsMentionsByPeriod

	if err := validateMentionsByPeriod(request); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["q"] = request.Q
	body["peerType"] = *request.PeerType
	body["startDate"] = *request.StartDate
	body["EndDate"] = *request.EndDate
	body["hideForwards"] = func() string { if *request.HideForwards == true { return "1" } else { return "0" } }()
	body["strongSearch"] = func() string { if *request.StrongSearch == true { return "1" } else { return "0" } }()
	body["minusWords"] = *request.MinusWords
	body["group"] = *request.Group
	body["extendedSyntax"] = func() string { if *request.ExtendedSyntax == true { return "1" } else { return "0" } }()
	req, err := c.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.WordsMentionsResponse
	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func validateMentionsByPeriod(request MentionPeriodRequest) error {
	if request.Q == "" {
		return fmt.Errorf("q must be set")
	}
	return nil
}

type MentionsByChannelRequest struct {
	Q string
	PeerType *string
	StartDate *string
	EndDate *string
	HideForwards *bool
	StrongSearch *bool
	MinusWords *string
	ExtendedSyntax *bool
}

func (c *Client) WordsMentionsByChannels(ctx context.Context, request MentionsByChannelRequest) (*schema.WordsMentionsResponse, *http.Response, error) {
	path := endpoints.WordsMentionsByChannels

	if err := validateMentionsByChannelsRequest(request); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["q"] = request.Q
	body["peerType"] = *request.PeerType
	body["startDate"] = *request.StartDate
	body["EndDate"] = *request.EndDate
	body["hideForwards"] = func() string { if *request.HideForwards == true { return "1" } else { return "0" } }()
	body["strongSearch"] = func() string { if *request.StrongSearch == true { return "1" } else { return "0" } }()
	body["minusWords"] = *request.MinusWords
	body["extendedSyntax"] = func() string { if *request.ExtendedSyntax == true { return "1" } else { return "0" } }()

	req, err := c.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.WordsMentionsResponse
	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func validateMentionsByChannelsRequest(request MentionsByChannelRequest) error {

	if request.Q == "" {
		return fmt.Errorf("q must be set")
	}

	return nil
}
