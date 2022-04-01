package main

import (
	"context"
	"fmt"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/database"
	"os"
)

func getToken() (key string, err error) {
	key = os.Getenv("TOKEN")
	if key == "" {
		return "", fmt.Errorf("token not found")
	}

	return key, nil
}

func main() {
	token, err := getToken()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Token is %s\n", token)

	tgstat.Token = token
	languages, _, err := database.LanguagesGet(context.Background(), "ru")

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
