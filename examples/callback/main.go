package main

import (
	"context"
	"fmt"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/callback"
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
