package main

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/usage"
	"os"
	"time"
)

var qs = []*survey.Question{
	{
		Name:     "Token",
		Prompt:   &survey.Input{Message: "Enter your token"},
		Validate: survey.Required,
	},
}

func main() {
	answers := struct {
		Token string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tgstat.Token = answers.Token
	req, _, err := usage.Stat(context.Background())

	if err != nil {
		fmt.Printf("error getting data: %v\n", err)
		os.Exit(1)
	}

	for _, stat := range req.Response {
		fmt.Print("Stat Info\n")
		fmt.Printf("Title: %s\n", stat.Title)
		fmt.Printf("Service key: %s\n", stat.ServiceKey)
		if stat.SpentChannels != "" {
			fmt.Printf("SpentChannels: %s\n", stat.SpentChannels)
		}
		fmt.Printf("spentRequests: %s\n", stat.SpentRequests)
		if stat.SpentWords != "" {
			fmt.Printf("spentWords: %s\n", stat.SpentWords)
		}
		if stat.SpentObjects != "" {
			fmt.Printf("spentObjects: %s\n", stat.SpentObjects)
		}
		fmt.Printf("Expired at: %s\n", time.Unix(stat.ExpiredAt, 0))
	}
}
