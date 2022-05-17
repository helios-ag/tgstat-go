package main

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/words"
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
		Name:   "Q",
		Prompt: &survey.Input{Message: "Q"},
	},
	{
		Name:   "PeerType",
		Prompt: &survey.Select{Message: "Choose grouping", Options: []string{"channel", "chat", "all"}},
	},
	{
		Name:   "StartDate",
		Prompt: &survey.Input{Message: "Start Time", Default: ""},
	},
	{
		Name:   "EndDate",
		Prompt: &survey.Input{Message: "End Time", Default: ""},
	},
	{
		Name:   "HideForwards",
		Prompt: &survey.Confirm{Message: "Hide Forwards", Default: false},
	},
	{
		Name:   "StrongSearch",
		Prompt: &survey.Confirm{Message: "Strong Search", Default: false},
	},
	{
		Name:   "Group",
		Prompt: &survey.Select{Message: "Choose grouping", Options: []string{"day", "week", "month"}},
	},
	{
		Name:   "ExtendedSyntax",
		Prompt: &survey.Confirm{Message: "Enable extended syntax", Default: false},
	},
	{
		Name:   "MinusWords",
		Prompt: &survey.Input{Message: "Minus Words"},
	},
}

func main() {
	answers := struct {
		Token          string
		Q              string
		PeerType       string
		StartDate      string
		EndDate        string
		HideForwards   bool
		StrongSearch   bool
		MinusWords     string
		Group          string
		ExtendedSyntax bool
	}{}

	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var startTime, endTime string
	if answers.StartDate != "" {
		startTime = strconv.FormatInt(time.Now().Unix()-86400, 10)
	}

	if answers.EndDate != "" {
		endTime = strconv.FormatInt(time.Now().Unix(), 10)
	}

	var group *string
	if answers.Group != "" {
		group = String(answers.Group)
	}
	req := words.MentionPeriodRequest{
		Q:              answers.Q,
		PeerType:       String(answers.PeerType),
		StartDate:      String(startTime),
		EndDate:        String(endTime),
		HideForwards:   Bool(answers.HideForwards),
		StrongSearch:   Bool(answers.StrongSearch),
		MinusWords:     String(answers.MinusWords),
		Group:          group,
		ExtendedSyntax: Bool(answers.ExtendedSyntax),
	}

	tgstat.Token = answers.Token

	info, _, err := words.MentionsByPeriod(context.Background(), req)

	if err != nil {
		fmt.Printf("error getting data: %v\n", err)
		os.Exit(1)
	}

	for _, item := range info.Response.Items {
		fmt.Printf("ViewsCount %d\n", item.ViewsCount)
		fmt.Printf("Period %s\n", item.Period)
		fmt.Printf("MentionsCount %d\n", item.MentionsCount)
	}

	chanReq := words.MentionsByChannelRequest{
		Q:              answers.Q,
		PeerType:       String(answers.PeerType),
		StartDate:      String(startTime),
		EndDate:        String(endTime),
		HideForwards:   Bool(answers.HideForwards),
		StrongSearch:   Bool(answers.StrongSearch),
		MinusWords:     String(answers.MinusWords),
		ExtendedSyntax: Bool(answers.ExtendedSyntax),
	}

	tgstat.Token = answers.Token

	mentions, _, err := words.MentionsByChannels(context.Background(), chanReq)

	if err != nil {
		fmt.Printf("error getting data: %v\n", err)
		os.Exit(1)
	}

	for _, mentions := range mentions.Response.Items {
		fmt.Printf("ViewsCount %d\n", mentions.ViewsCount)
		fmt.Printf("ChannelID %d\n", mentions.ChannelID)
		fmt.Printf("MentionsCount %d\n", mentions.MentionsCount)
		fmt.Printf("LastMentionDate %d\n", mentions.LastMentionDate)
	}

	for _, channelInfo := range mentions.Response.Channels {
		fmt.Print("Channel Info")
		fmt.Printf("Title: %s\n", channelInfo.Title)
		fmt.Printf("Id: %d\n", channelInfo.ID)
		fmt.Printf("Username: %s\n", channelInfo.Username)
		fmt.Printf("Title: %s\n", channelInfo.Title)
		fmt.Printf("About: %s\n", channelInfo.About)
		fmt.Printf("Image100: %s\n", channelInfo.Image100)
		fmt.Printf("Image640: %s\n", channelInfo.Image640)
		fmt.Printf("ParticipantsCount: %d\n", channelInfo.ParticipantsCount)
	}

}

func String(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func Bool(b bool) *bool {
	return &b
}
