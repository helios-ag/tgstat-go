package tgstat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"tgstat/endpoints"
	"tgstat/schema"
)

func (c *Client) PostGet(ctx context.Context, postId string) (*schema.PostResponse, *http.Response, error) {
	path := endpoints.PostsGet

	if postId == "" {
		return nil, nil, fmt.Errorf("postId can not be empty")
	}

	body := make(map[string]string)
	body["postId"] = postId
	req, err := c.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.PostResponse
	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

type PostStatRequest struct {
	PostId string
	Group *string
}

func (c *Client) PostStat(ctx context.Context, request PostStatRequest) (*schema.PostStatResponse, *http.Response, error) {
	path := endpoints.PostsStat

	if request.PostId == "" {
		return nil, nil, fmt.Errorf("postId can not be empty")
	}
	if !validateGroupStat(*request.Group) {
		return nil, nil, fmt.Errorf("improper group value")
	}

	body := make(map[string]string)
	body["postId"] = request.PostId
	body["group"] = *request.Group
	req, err := c.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.PostStatResponse
	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func validateGroupStat(group string) bool {
	switch group {
	case
		"day",
		"hour":
		return true
	}

	return false
}

type PostSearchRequest struct {
	Q string
	Limit *int
	Offset *int
	PeerType *string
	StartDate *string
	EndDate *string
	HideForwards *bool
	HideDeleted *bool
	StrongSearch *bool
	MinusWords *string
	ExtendedSyntax *bool
}

func (c *Client) PostSearch(ctx context.Context, request PostSearchRequest) (*schema.PostSearchResponse, *http.Response, error) {
	path := endpoints.PostsSearch

	if err := validateChannelSearchRequest(request); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["q"] = request.Q
	body["limit"] = strconv.Itoa(*request.Limit)
	body["offset"] = strconv.Itoa(*request.Offset)
	body["peerType"] = *request.PeerType
	body["startDate"] = *request.StartDate
	body["EndDate"] = *request.EndDate
	body["hideForwards"] = func() string { if *request.HideForwards == true { return "1" } else { return "0" } }()
	body["hideDeleted"] = func() string { if *request.HideDeleted == true { return "1" } else { return "0" } }()
	body["minusWords"] = *request.MinusWords
	body["extendedSyntax"] = func() string { if *request.ExtendedSyntax == true { return "1" } else { return "0" } }()

	req, err := c.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.PostSearchResponse
	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func validateChannelSearchRequest(request PostSearchRequest) error {
	if request.Q == "" {
		return fmt.Errorf("q must be set")
	}

	if *request.Limit > 50 {
		return fmt.Errorf("max limit is 50")
	}

	if *request.Offset > 1000 {
		return fmt.Errorf("max offset is 50")
	}

	return nil
}

func (c *Client) PostSearchExtended(ctx context.Context, request PostSearchRequest) (*schema.PostSearchExtendedResponse, *http.Response, error) {
	path := endpoints.PostsSearch

	if err := validateChannelSearchRequest(request); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["q"] = request.Q
	body["limit"] = strconv.Itoa(*request.Limit)
	body["offset"] = strconv.Itoa(*request.Offset)
	body["peerType"] = *request.PeerType
	body["startDate"] = *request.StartDate
	body["EndDate"] = *request.EndDate
	body["hideForwards"] = func() string { if *request.HideForwards == true { return "1" } else { return "0" } }()
	body["hideDeleted"] = func() string { if *request.HideDeleted == true { return "1" } else { return "0" } }()
	body["minusWords"] = *request.MinusWords
	body["extendedSyntax"] = func() string { if *request.ExtendedSyntax == true { return "1" } else { return "0" } }()
	body["extended"] = "1"

	req, err := c.NewRestRequest(ctx, "GET", path, body)

	if err != nil {
		return nil, nil, err
	}

	var response schema.PostSearchExtendedResponse
	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}