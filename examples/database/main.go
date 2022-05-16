package main

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/database"
	"os"
)

var qs = []*survey.Question{
	{
		Name:     "Token",
		Prompt:   &survey.Input{Message: "Enter your token"},
		Validate: survey.Required,
	},
	{
		Name:   "Language",
		Prompt: &survey.Input{Message: "Enter Language"},
	},
}

func main() {
	answers := struct {
		Token    string
		Language string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tgstat.Token = answers.Token
	languages, _, err := database.LanguagesGet(context.Background(), answers.Language)

	if err != nil {
		fmt.Printf("error getting data: %v\n", err)
		os.Exit(1)
	}

	fmt.Print("Languages list")
	for _, info := range languages.Response {
		fmt.Printf("Name: %s\n", info.Name)
		fmt.Printf("Code: %s\n", info.Code)
	}

	categories, _, err := database.CategoriesGet(context.Background(), "ru")

	if err != nil {
		fmt.Printf("error getting data: %v\n", err)
		os.Exit(1)
	}
	fmt.Print("Categories list")
	for _, info := range categories.Response {
		fmt.Printf("Name: %s\n", info.Name)
		fmt.Printf("Code: %s\n", info.Code)
	}
	countries, _, err := database.CountriesGet(context.Background(), "ru")
	if err != nil {
		fmt.Printf("error getting data: %v\n", err)
		os.Exit(1)
	}

	fmt.Print("Countries list")
	for _, info := range countries.Response {
		fmt.Printf("Name: %s\n", info.Name)
		fmt.Printf("Code: %s\n", info.Code)
	}

	os.Exit(0)
}
