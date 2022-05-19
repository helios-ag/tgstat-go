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
		Prompt:   &survey.Input{Message: "Enter channel id"},
		Validate: survey.Required,
	},
	{
		Name:   "Group",
		Prompt: &survey.Input{Message: "Enter group"},
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
}

func main() {
	answers := struct {
		Token     string
		ChannelId string
		Group     string
		StartTime string
		EndTime   string
	}{}

	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var startTime, endTime string
	if answers.StartTime != "" {
		startTime = strconv.FormatInt(time.Now().Unix()-86400, 10)
	}

	if answers.EndTime != "" {
		endTime = strconv.FormatInt(time.Now().Unix(), 10)
	}

	var group *string
	if answers.Group != "" {
		group = String(answers.Group)
	}
	req := channels.ChannelViewsRequest{
		ChannelId: answers.ChannelId,
		StartDate: String(startTime),
		EndDate:   String(endTime),
		Group:     group,
	}

	tgstat.Token = answers.Token

	info, _, err := channels.Err(context.Background(), req)

	if err != nil {
		fmt.Printf("error getting data: %v\n", err)
		os.Exit(1)
	}
	fmt.Print("Err values")
	for _, info := range info.Response {
		fmt.Printf("Err value: %f\n", info.Err)
		fmt.Printf("Period: %s\n", info.Period)
	}

	views, _, err := channels.Views(context.Background(), req)

	if err != nil {
		fmt.Printf("error getting data: %v\n", err)
		os.Exit(1)
	}

	fmt.Print("Views")
	for _, info := range views.Response {
		fmt.Print("Err values")
		fmt.Printf("Err value: %f\n", info.ViewsCount)
		fmt.Printf("Period: %s\n", info.Period)
	}

	avginfo, _, err := channels.AvgPostsReach(context.Background(), req)

	if err != nil {
		fmt.Printf("error getting data: %v\n", err)
		os.Exit(1)
	}

	fmt.Print("Avg Post Reach")
	for _, info := range avginfo.Response {
		fmt.Print("Err values")
		fmt.Printf("Err value: %f\n", info.AvgPostsReach)
		fmt.Printf("Period: %s\n", info.Period)
	}

}

func String(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
