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
		Name:     "Token",
		Prompt:   &survey.Input{Message: "Enter your token"},
		Validate: survey.Required,
	},
	{
		Name:     "ChannelId",
		Prompt:   &survey.Input{Message: "Enter Channel ID"},
		Validate: survey.Required,
	},
	{
		Name:   "Limit",
		Prompt: &survey.Input{Message: "Limit", Default: "10"},
	},
	{
		Name:   "Offset",
		Prompt: &survey.Input{Message: "Offset", Default: "0"},
	},
	{
		Name:   "StartTime",
		Prompt: &survey.Input{Message: "Start Time", Default: ""},
	},
	{
		Name:   "EndTime",
		Prompt: &survey.Input{Message: "End Time", Default: ""},
	},
}

func main() {
	answers := struct {
		Token     string
		ChannelId string
		Limit     string
		Offset    string
		StartTime string
		EndTime   string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

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

	tgstat.Token = answers.Token

	req := channels.ChannelForwardRequest{
		ChannelId: answers.ChannelId,
		Limit:     Uint(limit),
		Offset:    Uint(offset),
		StartDate: String(startTime),
		EndDate:   String(endTime),
	}

	forwards, _, err := channels.Forwards(context.Background(), req)

	if err != nil {
		fmt.Printf("error getting data: %v\n", err)
		os.Exit(1)
	}

	for _, info := range forwards.Response.Items {
		fmt.Printf("ChannelId: %d\n", info.ChannelID)
		fmt.Printf("PostID: %d\n", info.PostID)
		fmt.Printf("ForwardID: %d\n", info.ForwardID)
		fmt.Printf("PostDate: %d\n", info.PostDate)
		fmt.Printf("PostLink: %s\n", info.PostLink)
	}

	mentions, _, err := channels.Mentions(context.Background(), req)

	if err != nil {
		fmt.Printf("error getting data: %v\n", err)
		os.Exit(1)
	}

	for _, info := range mentions.Response.Items {
		fmt.Printf("ChannelId: %d\n", info.ChannelID)
		fmt.Printf("PostID: %d\n", info.PostID)
		fmt.Printf("MentionID: %d\n", info.MentionID)
		fmt.Printf("MentionType: %s\n", info.MentionType)
		fmt.Printf("PostDate: %d\n", info.PostDate)
		fmt.Printf("PostLink: %s\n", info.PostLink)
	}

	os.Exit(0)
}

func Uint(i uint64) *uint64 {
	return &i
}
func String(v string) *string {
	return &v
}
