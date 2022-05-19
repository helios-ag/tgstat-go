package main

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/callback"
	"os"
)

var qs = []*survey.Question{
	{
		Name:     "Token",
		Prompt:   &survey.Input{Message: "Enter your token"},
		Validate: survey.Required,
	},
	{
		Name:   "ChannelId",
		Prompt: &survey.Input{Message: "Enter Channel ID"},
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

	req := callback.SubscribeChannelRequest{
		SubscriptionId: nil,
		ChannelId:      "",
		EventTypes:     "",
	}

	sub, _, err := callback.SubscribeChannel(context.Background(), req)

	if err != nil {
		fmt.Printf("error getting data: %v\n", err)
		os.Exit(1)
	}

	fmt.Print("Subscription ID")
	fmt.Printf("Title: %d\n", sub.SubscriptionId)

	os.Exit(0)
}
