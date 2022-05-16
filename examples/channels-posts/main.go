package main

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/channels"
	"os"
	"strconv"
	"time"
)

var qs = []*survey.Question{
	{
		Name:      "Token",
		Prompt:    &survey.Input{Message: "Enter your token"},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name:      "ChannelId",
		Prompt:    &survey.Input{Message: "Enter Channel ID"},
		Transform: survey.Title,
	},
	{
		Name:      "Limit",
		Prompt:    &survey.Input{Message: "Limit", Default: "10"},
		Transform: survey.Title,
	},
	{
		Name:      "Offset",
		Prompt:    &survey.Input{Message: "Offset", Default: "0"},
		Transform: survey.Title,
	},
	{
		Name:      "StartTime",
		Prompt:    &survey.Input{Message: "Start Time", Default: ""},
		Transform: survey.Title,
	},
	{
		Name:      "EndTime",
		Prompt:    &survey.Input{Message: "End Time", Default: ""},
		Transform: survey.Title,
	},
	{
		Name:      "HideForwards",
		Prompt:    &survey.Input{Message: "Hide Forwards", Default: "1"},
		Transform: survey.Title,
	},
	{
		Name:      "HideDeleted",
		Prompt:    &survey.Input{Message: "Hide Deleted", Default: "1"},
		Transform: survey.Title,
	},
}

func main() {
	answers := struct {
		Token        string
		ChannelId    string
		Limit        string
		Offset       string
		StartTime    string
		EndTime      string
		HideForwards string
		HideDeleted  string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tgstat.Token = answers.Token
	var limit, offset uint64
	if answers.Limit != "0" {
		limit, _ = strconv.ParseUint(answers.Limit, 10, 64)
	}

	if answers.Offset != "0" {
		offset, _ = strconv.ParseUint(answers.Offset, 10, 64)
	}

	var startTime, endTime string
	if answers.StartTime != "" {
		startTime = strconv.FormatInt(time.Now().Unix()-86400, 10)
	}

	if answers.EndTime != "" {
		endTime = strconv.FormatInt(time.Now().Unix(), 10)
	}

	var hideForwards, hideDeleted bool
	if answers.HideForwards != "" {
		hideForwards, _ = strconv.ParseBool(answers.HideForwards)
	}

	if answers.HideDeleted != "" {
		hideDeleted, _ = strconv.ParseBool(answers.HideDeleted)
	}

	req := channels.PostsRequest{
		ChannelId:    answers.ChannelId,
		Limit:        Uint(limit),
		Offset:       Uint(offset),
		StartTime:    String(startTime),
		EndTime:      String(endTime),
		HideForwards: &hideForwards,
		HideDeleted:  &hideDeleted,
	}
	info, _, err := channels.Posts(context.Background(), req)

	if err != nil {
		fmt.Printf("error getting data: %v\n", err)
		os.Exit(1)
	}

	fmt.Print("Channel posts")

	fmt.Printf("Count: %d\n", info.Response.Count)
	fmt.Printf("Total: %d\n", info.Response.TotalCount)
	for _, item := range info.Response.Items {
		fmt.Printf("ID: %d\n", item.ID)
		fmt.Printf("Date: %d\n", item.Date)
		fmt.Printf("Views: %d\n", item.Views)
		fmt.Printf("Link: %s\n", item.Link)
		fmt.Printf("ChannelID: %d\n", item.ChannelID)
		fmt.Printf("ForwardedFrom: %s\n", item.ForwardedFrom)
		fmt.Printf("Is_deleted: %b\n", item.IsDeleted)
		fmt.Printf("Text: %s\n", item.Text)
		fmt.Printf("Media_type: %s\n", item.Media.MediaType)
		fmt.Printf("mime_type: %s\n", item.Media.MimeType)
		fmt.Printf("size: %d\n", item.Media.Size)
	}
}

func Uint(i uint64) *uint64 {
	return &i
}
func String(v string) *string {
	return &v
}
