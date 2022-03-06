package main

import (
	"context"
	"fmt"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/channels"
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

	channelInfo, _, err := channels.Get(context.Background(), "https://t.me/readovkanews")

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

	os.Exit(0)
}
