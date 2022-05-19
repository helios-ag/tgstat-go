package main

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/channels"
	"os"
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
}

func main() {
	answers := struct {
		Token     string
		ChannelId string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tgstat.Token = answers.Token

	channelInfo, _, err := channels.Get(context.Background(), answers.ChannelId)

	if err != nil {
		fmt.Printf("error getting data: %v\n", err)
		os.Exit(1)
	}

	fmt.Print("Channel Info")
	fmt.Printf("Title: %s\n", channelInfo.Response.Title)
	fmt.Printf("Id: %d\n", channelInfo.Response.Id)
	fmt.Printf("Username: %s\n", channelInfo.Response.Username)
	fmt.Printf("Title: %s\n", channelInfo.Response.Title)
	fmt.Printf("About: %s\n", channelInfo.Response.About)
	fmt.Printf("Category: %s\n", channelInfo.Response.Category)
	fmt.Printf("Country: %s\n", channelInfo.Response.Country)
	fmt.Printf("Language: %s\n", channelInfo.Response.Language)
	fmt.Printf("Image100: %s\n", channelInfo.Response.Image100)
	fmt.Printf("Image640: %s\n", channelInfo.Response.Image640)
	fmt.Printf("ParticipantsCount: %d\n", channelInfo.Response.ParticipantsCount)
	if channelInfo.Response.TGStatRestriction != nil {
		data := channelInfo.Response.TGStatRestriction.(tgstat.TGStatRestrictions)
		fmt.Printf("RedLabel: %s\n", bool2string(data.RedLabel))
		fmt.Printf("BlackLabel: %s\n", bool2string(data.BlackLabel))
	}

	os.Exit(0)
}

func bool2string(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}
