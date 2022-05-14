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

	channelInfo, _, err := channels.Get(context.Background(), "https://t.me/nim_ru")

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
	fmt.Printf("RedLabel: %s\n", bool2string(channelInfo.Response.TGStatRestriction.RedLabel))
	fmt.Printf("BlackLabel: %s\n", bool2string(channelInfo.Response.TGStatRestriction.BlackLabel))

	os.Exit(0)
}

func bool2string(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}
