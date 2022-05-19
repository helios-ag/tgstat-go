package tgstat_go

type Media struct {
	MediaType string `json:"media_type"`
	Caption   string `json:"caption"`
}

type PostResponse struct {
	ID            int         `json:"id"`
	Date          int         `json:"date"`
	Views         int         `json:"views"`
	Link          string      `json:"link"`
	ChannelID     int         `json:"channel_id"`
	ForwardedFrom interface{} `json:"forwarded_from"`
	IsDeleted     int         `json:"is_deleted"`
	Text          string      `json:"text"`
	Media         `json:"media"`
}

type PostResult struct {
	Status   string       `json:"status"`
	Response PostResponse `json:"response"`
}

type Forward []struct {
	PostID    string `json:"postId"`
	PostLink  string `json:"postLink"`
	PostDate  string `json:"postDate"`
	ChannelID int    `json:"channelId"`
}

type Mention struct {
	PostID    string `json:"postId,omitempty"`
	PostLink  string `json:"postLink,omitempty"`
	PostDate  string `json:"postDate,omitempty"`
	ChannelID int    `json:"channelId,omitempty"`
}

type View struct {
	Date        string `json:"date"`
	ViewsGrowth int    `json:"viewsGrowth"`
}

type PostStatResponse struct {
	ViewsCount    int       `json:"viewsCount"`
	ForwardsCount int       `json:"forwardsCount"`
	MentionsCount int       `json:"mentionsCount"`
	Forwards      []Forward `json:"forwards"`
	Mentions      []Mention `json:"mentions"`
	Views         []View    `json:"views"`
}

type PostStatResult struct {
	Status   string           `json:"status"`
	Response PostStatResponse `json:"response"`
}

type PostSearchResultItem struct {
	ID            int64       `json:"id"`
	Date          int         `json:"date"`
	Views         int         `json:"views"`
	Link          string      `json:"link"`
	ChannelID     int         `json:"channel_id"`
	ForwardedFrom interface{} `json:"forwarded_from"`
	IsDeleted     int         `json:"is_deleted"`
	Text          string      `json:"text"`
	Snippet       string      `json:"snippet"`
	Media         struct {
		MediaType string `json:"media_type"`
		MimeType  string `json:"mime_type"`
		Size      int    `json:"size"`
	} `json:"media"`
}

type PostSearchResultResponse struct {
	Count      int                    `json:"count"`
	TotalCount int                    `json:"total_count"`
	Items      []PostSearchResultItem `json:"items"`
}

type PostSearchResult struct {
	Status   string                   `json:"status"`
	Response PostSearchResultResponse `json:"response"`
}
type PostSearchExtendedResponseItem struct {
	ID            int64       `json:"id"`
	Date          int         `json:"date"`
	Views         int         `json:"views"`
	Link          string      `json:"link"`
	ChannelID     int         `json:"channel_id"`
	ForwardedFrom interface{} `json:"forwarded_from"`
	IsDeleted     int         `json:"is_deleted"`
	Text          string      `json:"text"`
	Snippet       string      `json:"snippet"`
	Media         struct {
		MediaType string `json:"media_type"`
		MimeType  string `json:"mime_type"`
		Size      int    `json:"size"`
	} `json:"media"`
}
type PostSearchExtendedChannel struct {
	ID                int    `json:"id"`
	Link              string `json:"link"`
	Username          string `json:"username"`
	Title             string `json:"title"`
	About             string `json:"about"`
	Image100          string `json:"image100"`
	Image640          string `json:"image640"`
	ParticipantsCount int    `json:"participants_count"`
}

type PostSearchExtendedResponse struct {
	Count      int                              `json:"count"`
	TotalCount int                              `json:"total_count"`
	Items      []PostSearchExtendedResponseItem `json:"items"`
	Channels   []PostSearchExtendedChannel      `json:"channels"`
}

type PostSearchExtendedResult struct {
	Status   string                     `json:"status"`
	Response PostSearchExtendedResponse `json:"response"`
}
