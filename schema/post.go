package schema

type PostResponse struct {
	Status   string `json:"status"`
	Response struct {
		ID            int         `json:"id"`
		Date          int         `json:"date"`
		Views         int         `json:"views"`
		Link          string      `json:"link"`
		ChannelID     int         `json:"channel_id"`
		ForwardedFrom interface{} `json:"forwarded_from"`
		IsDeleted     int         `json:"is_deleted"`
		Text          string      `json:"text"`
		Media         struct {
			MediaType string `json:"media_type"`
			Caption   string `json:"caption"`
		} `json:"media"`
	} `json:"response"`
}

type PostStatResponse struct {
	Status   string `json:"status"`
	Response struct {
		ViewsCount    int `json:"viewsCount"`
		ForwardsCount int `json:"forwardsCount"`
		MentionsCount int `json:"mentionsCount"`
		Forwards      []struct {
			PostID    string `json:"postId"`
			PostLink  string `json:"postLink"`
			PostDate  string `json:"postDate"`
			ChannelID int    `json:"channelId"`
		} `json:"forwards"`
		Mentions []struct {
			PostID    string `json:"postId,omitempty"`
			PostLink  string `json:"postLink,omitempty"`
			PostDate  string `json:"postDate,omitempty"`
			ChannelID int    `json:"channelId,omitempty"`
		} `json:"mentions"`
		Views []struct {
			Date        string `json:"date"`
			ViewsGrowth int    `json:"viewsGrowth"`
		} `json:"views"`
	} `json:"response"`
}

type PostSearchResponse struct {
	Status   string `json:"status"`
	Response struct {
		Count      int `json:"count"`
		TotalCount int `json:"total_count"`
		Items      []struct {
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
		} `json:"items"`
	} `json:"response"`
}

type PostSearchExtendedResponse struct {
	Status   string `json:"status"`
	Response struct {
		Count      int `json:"count"`
		TotalCount int `json:"total_count"`
		Items      []struct {
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
		} `json:"items"`
		Channels []struct {
			ID                int    `json:"id"`
			Link              string `json:"link"`
			Username          string `json:"username"`
			Title             string `json:"title"`
			About             string `json:"about"`
			Image100          string `json:"image100"`
			Image640          string `json:"image640"`
			ParticipantsCount int    `json:"participants_count"`
		} `json:"channels"`
	} `json:"response"`
}