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

	info, _, err := channels.Stat(context.Background(), answers.ChannelId)

	if err != nil {
		fmt.Printf("error getting data: %v\n", err)
		os.Exit(1)
	}

	fmt.Print("Channel Info")
	fmt.Printf("Title: %s\n", info.Response.Title)
	fmt.Printf("Id: %d\n", info.Response.Id)
	fmt.Printf("Username: %s\n", info.Response.Username)
	fmt.Printf("ParticipantsCount: %d\n", info.Response.ParticipantsCount)
	fmt.Printf("Average post reach: %d\n", info.Response.AvgPostReach)
	fmt.Printf("Error percent: %f\n", info.Response.ErrPercent)
	fmt.Printf("Daily Reach: %d\n", info.Response.DailyReach)
	fmt.Printf("CI Index: %f\n", info.Response.CiIndex)

}
