TGStat GO API Wrapper
=====================

[![LICENSE](https://img.shields.io/badge/license-MIT-orange.svg)](LICENSE)
[![Go](https://github.com/helios-ag/tgstat-go/actions/workflows/go.yaml/badge.svg)](https://github.com/helios-ag/tgstat-go/actions/workflows/go.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/helios-ag/tgstat-go)](https://goreportcard.com/report/github.com/helios-ag/tgstat-go)
[![Godocs](https://img.shields.io/badge/golang-documentation-blue.svg)](https://godoc.org/github.com/helios-ag/tgstat-go)

[TGStat API](https://api.tgstat.ru/docs/ru/start/intro.html) written in Go

[TGStat](https://tgstat.ru) is service that collects information about different channels and chats    

- [Installation](#installation)
- [Getting started](#getting-started)
    * [Step 1](#step-1)
    * [Step 2 ](#step-2)
    * [Step 3](#step-3)
    * [Step 4](#step-4)
- [Available methods](#available-methods)
    * [Channels](#channels)
        + [Get Channel Info  ](#get-channel-info)
        + [Search among channels](#search-among-channels)
        + [Get channel stat](#get-channel-stat)
        + [Get channel posts](#get-channel-posts)
        + [Get channel mentions](#get-channel-mentions)
        + [Get channel forwards](#get-channel-forwards)
        + [Get channel subscribers](#get-channel-subscribers)
        + [Get channel views](#get-channel-views)
        + [Get channel average posts reach](#get-channel-average-posts-reach)
        + [Add channel ](#add-channel)
        + [Get channel ERR rate](#get-channel-err-rate)
    * [Posts](#posts)
        + [Get post](#get-post)
        + [Post statistics](#post-statistics)
        + [Post search](#post-search)
    * [Words](#words)
        + [Mentions by period](#mentions-by-period)
        + [Mentions by channel](#mentions-by-channel)
    * [Database](#database)
        + [Categories](#categories)
        + [Countries](#countries)
        + [Languages](#languages)
    * [Usage](#usage)
        + [Statistics](#statistics)
    * [API Callback](#api-callback)

## Installation

Make sure your project is using Go Modules (it will have a `go.mod` file in its
root if it already is):

``` sh
go mod init
```

Then, reference stripe-go in a Go program with `import`:

``` go
import (
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/channels"
)
```

Run any of the normal `go` commands (`build`/`install`/`test`). The Go
toolchain will resolve and fetch the stripe-go module automatically.

Alternatively, you can also explicitly `go get` the package into a project:

```bash
go get -u github.com/helios-ag/tgstat-go
```

## Getting started

### Step 1
Obtain key by authorizing on https://api.tgstat.ru/docs/ru/start/token.html

### Step 2 
After getting token, you must set token assigning it to `tgstat.Token` value. 

### Step 3

After setting you token, you can call, for example, method from channels package: `channels.Get(context.Background(), "https://t.me/nim_ru")`

Example below: 

```go
// example.go
package main

import (
	"context"
	"fmt"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/channels"
	"os"
)

func main() {
        ...

	tgstat.Token = "yourtoken"

	channelInfo, _, err := channels.Get(context.Background(), "https://t.me/nim_ru")

	if err != nil {
		fmt.Printf("error getting data: %v\n", err)
		os.Exit(1)
	}

	fmt.Print("Channel Info")
	...
	fmt.Printf("Title: %s\n", channelInfo.Response.Title)
	...
}
```

### Step 4

Run example `go build example.go`

All examples available at [examples repository](https://github.com/helios-ag/tgstat-go-examples)

## Available methods

### Channels

#### Get Channel Info  

Docs at: https://api.tgstat.ru/docs/ru/channels/get.html

`func Get(ctx context.Context, channelId string)`

#### Search among channels

Docs at: https://api.tgstat.ru/docs/ru/channels/search.html

`func Search(ctx context.Context, request SearchRequest)`

Example [`/examples/channels-search/main.go`](examples/channels-search/main.go)

#### Get channel stat

Docs at: https://api.tgstat.ru/docs/ru/channels/search.html

`func Search(ctx context.Context, request SearchRequest)`

#### Get channel posts

Docs at: https://api.tgstat.ru/docs/ru/channels/posts.html

`func Posts(ctx context.Context, request PostsRequest)`

#### Get channel mentions

Docs at: https://api.tgstat.ru/docs/ru/channels/mentions.html

`func Mentions(ctx context.Context, request ChannelMentionsRequest)`

#### Get channel forwards

Docs at: https://api.tgstat.ru/docs/ru/channels/forwards.html

`func (c Client) Forwards(ctx context.Context, request ChannelForwardRequest)`

#### Get channel subscribers

Docs at: https://api.tgstat.ru/docs/ru/channels/subscribers.html

`func Subscribers(ctx context.Context, request ChannelSubscribersRequest)`

#### Get channel views

Docs at: https://api.tgstat.ru/docs/ru/channels/views.html

`func Views(ctx context.Context, request ChannelViewsRequest)`

#### Get channel average posts reach

Docs at: https://api.tgstat.ru/docs/ru/channels/avg-posts-reach.html

`func AvgPostsReach(ctx context.Context, request ChannelViewsRequest)`

#### Add channel 

Docs at: https://api.tgstat.ru/docs/ru/channels/add.html

`func Add(ctx context.Context, request ChannelAddRequest)`

#### Get channel ERR rate

Docs at: https://api.tgstat.ru/channels/err

`func Err(ctx context.Context, request ChannelViewsRequest)`

### Posts

#### Get post

Docs at: https://api.tgstat.ru/docs/ru/posts/get.html

`func Get(ctx context.Context, postId string)`

#### Post statistics

Docs at: https://api.tgstat.ru/docs/ru/posts/stat.html

`func PostStat(ctx context.Context, request PostStatRequest)`

#### Post search

Docs at: https://api.tgstat.ru/docs/ru/posts/search.html

`func PostSearch(ctx context.Context, request PostSearchRequest)`

and extended search

`func PostSearchExtended(ctx context.Context, request PostSearchRequest)`

### Words

#### Mentions by period

Docs at: https://api.tgstat.ru/docs/ru/words/mentions-by-period.html

`func MentionsByPeriod(ctx context.Context, request MentionPeriodRequest)`

#### Mentions by channel

Docs at: https://api.tgstat.ru/words/mentions-by-channels

`func MentionsByChannels(ctx context.Context, request MentionsByChannelRequest)`

### Database

#### Categories

Docs at: https://api.tgstat.ru/docs/ru/database/categories.html

`func CategoriesGet(ctx context.Context, lang string)`
####

#### Countries

Docs at: https://api.tgstat.ru/docs/ru/database/countries.html

`func CountriesGet(ctx context.Context, lang string)`
####

#### Languages

Docs at: https://api.tgstat.ru/docs/ru/database/languages.html

`func LanguagesGet(ctx context.Context, lang string)`
####

### Usage

#### Statistics

Example [`/examples/usage/main.go`](examples/usage/main.go)

Docs available at https://api.tgstat.ru/docs/ru/usage/stat.html

`func Stat(ctx context.Context)`


### API Callback

#### Set Callback URl

Docs available at https://api.tgstat.ru/docs/ru/callback/set-callback-url.html

`func SetCallback(ctx context.Context, callbackUrl string`

#### Get Callback Info

Docs available at https://api.tgstat.ru/docs/ru/callback/get-callback-info.html

`func GetCallbackInfo(ctx context.Context)`

#### Subscribe to channel

Docs available at https://api.tgstat.ru/docs/ru/callback/subscribe-channel.html

`func SubscribeChannel(ctx context.Context, request SubscribeChannelRequest)`

#### Subscribe to word

Docs available at https://api.tgstat.ru/docs/ru/callback/subscribe-word.html

`func SubscribeWord(ctx context.Context, request SubscribeWordRequest)`

#### Subscriptions list 

Docs available at https://api.tgstat.ru/docs/ru/callback/subscriptions-list.html

`func SubscriptionsList(ctx context.Context, subscriptionsListRequest SubscriptionsListRequest)`

#### Unsubscribe channel 

Docs available at https://api.tgstat.ru/docs/ru/callback/unsubscribe.html

`func Unsubscribe(ctx context.Context, subscriptionId string)`
