package main

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/channels"
	"os"
	"strconv"
)

var qs = []*survey.Question{
	{
		Name:      "Token",
		Prompt:    &survey.Input{Message: "Enter your token"},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name:      "Term",
		Prompt:    &survey.Input{Message: "Enter search term"},
		Transform: survey.Title,
	},
	{
		Name:      "Country",
		Prompt:    &survey.Input{Message: "Enter country (default Russia)", Default: "ru"},
		Transform: survey.Title,
	},
	{
		Name:   "SearchInDescription",
		Prompt: &survey.Confirm{Message: "Search in description?"},
	},
	{
		Name:      "Language",
		Prompt:    &survey.Input{Message: "Enter language here", Default: ""},
		Transform: survey.Title,
	},
	{
		Name:      "Category",
		Prompt:    &survey.Input{Message: "Enter category here", Default: ""},
		Transform: survey.Title,
	},
	{
		Name:      "Limit",
		Prompt:    &survey.Input{Message: "Enter limit here", Default: "10"},
		Transform: survey.Title,
	},
}

func main() {
	answers := struct {
		Token               string
		Term                string
		Country             string
		Category            string
		Language            string
		SearchInDescription bool
		Limit               string
	}{}

	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var lang *string
	lang = nil

	if answers.Language != "" {
		lang = String(answers.Language)
	}
	req := channels.SearchRequest{
		Q:                   answers.Term,
		SearchByDescription: bool2int(answers.SearchInDescription),
		Country:             answers.Country,
		Language:            lang,
		Category:            answers.Category,
		Limit:               string2int(answers.Limit),
	}

	tgstat.Token = answers.Token

	info, _, err := channels.Search(context.Background(), req)

	if err != nil {
		fmt.Printf("error getting data: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Count: %d\n", info.Response.Count)
	for _, channelInfo := range info.Response.Items {
		fmt.Print("Channel Info")
		fmt.Printf("Title: %s\n", channelInfo.Title)
		fmt.Printf("Id: %d\n", channelInfo.Id)
		fmt.Printf("Username: %s\n", channelInfo.Username)
		fmt.Printf("Title: %s\n", channelInfo.Title)
		fmt.Printf("About: %s\n", channelInfo.About)
		fmt.Printf("Image100: %s\n", channelInfo.Image100)
		fmt.Printf("Image640: %s\n", channelInfo.Image640)
		fmt.Printf("ParticipantsCount: %d\n", channelInfo.ParticipantsCount)
	}
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}

func string2int(v string) *int {
	res, _ := strconv.Atoi(v)
	return &res
}

func String(v string) *string {
	return &v
}
