package main

import (
	"context"
	"fmt"
	"github.com/helios-ag/tgstat-go/usage"
	"os"
	"time"
)

func getToken() (key string, err error) {
	key = os.Getenv("TOKEN")
	if "" == key {
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
	stat, _, err := usage.Stat(context.Background())

	if err != nil {
		fmt.Printf("error getting data: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("No error\n")
	for _, info := range stat.Response {
		fmt.Printf("ServiceKey: %s\n", info.ServiceKey)
		fmt.Printf("Title: %s\n", info.Title)
		fmt.Printf("SpentChannels: %s\n", info.SpentChannels)
		fmt.Printf("SpentRequests: %s\n", info.SpentRequests)
		fmt.Printf("ExpiredAt: %s\n", time.Unix(info.ExpiredAt, 0))
		fmt.Printf("SpentWords: %s\n", info.SpentWords)
	}
	os.Exit(0)
}
