package main

import (
	"context"
	"fmt"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/usage"
	"os"
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

	tgstat.Token = token
	stat, _, err := usage.Stat(context.Background())

	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting data: %v\n", err)
		os.Exit(1)
	}

	for _, info := range stat.Response {
		fmt.Sprintf("ServiceKey :%s", info.ServiceKey)
		fmt.Sprintf("Title :%s", info.Title)
		fmt.Sprintf("SpentChannels :%s", info.SpentChannels)
		fmt.Sprintf("SpentRequests :%s", info.SpentRequests)
		fmt.Sprintf("ExpiredAt :%s", info.ExpiredAt)
		fmt.Sprintf("SpentWords :%s", info.SpentWords)
	}
}
