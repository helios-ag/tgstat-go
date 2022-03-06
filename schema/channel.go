package schema

type Channel struct {
	Id                int               `json:"id"`
	Link              string            `json:"link"`
	Username          string            `json:"username"`
	Title             string            `json:"title"`
	About             string            `json:"about"`
	Image100          string            `json:"image100"`
	Image640          string            `json:"image640"`
	ParticipantsCount int               `json:"participants_count"`
	TGStatRestriction TGStatRestriction `json:"tgstat_restrictions"`
}

type TGStatRestriction struct {
	RedLabel   bool `json:"red_label"`
	BlackLabel bool `json:"black_label"`
}

type ChannelResponse struct {
	Status   string  `json:"status,string"`
	Response Channel `json:"response"`
}

type ChannelSearchResponse struct {
	Status   string `json:"status,string"`
	Response struct {
		Count int `json:"count"`
		Items []struct {
			Id                int    `json:"id"`
			Link              string `json:"link"`
			Username          string `json:"username"`
			Title             string `json:"title"`
			About             string `json:"about"`
			Image100          string `json:"image100"`
			Image640          string `json:"image640"`
			ParticipantsCount int    `json:"participants_count"`
		} `json:"items"`
	} `json:"response"`
}

type ChannelStatResponse struct {
	Status   string `json:"status,string"`
	Response struct {
		Id                int    `json:"id"`
		Title             string `json:"title"`
		Username          string `json:"username"`
		ParticipantsCount int    `json:"participants_count"`
		AvgPostReach      int    `json:"avg_post_reach"`
		ErrPercent        int    `json:"err_percent"`
		DailyReach        int    `json:"daily_reach"`
		CiIndex           int    `json:"ci_index"`
	} `json:"response"`
}

type ChannelPostsWithChannelResponse struct {
	Status   string `json:"status"`
	Response struct {
		Count      int `json:"count"`
		TotalCount int `json:"total_count"`
		Channel    struct {
			ID                int    `json:"id"`
			Link              string `json:"link"`
			Username          string `json:"username"`
			Title             string `json:"title"`
			About             string `json:"about"`
			Image100          string `json:"image100"`
			Image640          string `json:"image640"`
			ParticipantsCount int    `json:"participants_count"`
		} `json:"channel"`
		Items []struct {
			ID            int64       `json:"id"`
			Date          int         `json:"date"`
			Views         int         `json:"views"`
			Link          string      `json:"link"`
			ChannelID     int         `json:"channel_id"`
			ForwardedFrom interface{} `json:"forwarded_from"`
			IsDeleted     int         `json:"is_deleted"`
			Text          string      `json:"text"`
			Media         struct {
				MediaType string `json:"media_type"`
				MimeType  string `json:"mime_type"`
				Size      int    `json:"size"`
			} `json:"media"`
		} `json:"items"`
	} `json:"response"`
}

type ChannelPostsResponse struct {
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
			Media         struct {
				MediaType string `json:"media_type"`
				MimeType  string `json:"mime_type"`
				Size      int    `json:"size"`
			} `json:"media"`
		} `json:"items"`
	} `json:"response"`
}

type ChannelMentions struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type ChannelMentionsExtended struct {
	Status   string `json:"status"`
	Response struct {
		Items []struct {
			MentionID   int    `json:"mentionId"`
			MentionType string `json:"mentionType"`
			PostID      int64  `json:"postId"`
			PostLink    string `json:"postLink"`
			PostDate    int    `json:"postDate"`
			ChannelID   int    `json:"channelId"`
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

type ChannelForwardsExtended struct {
	Status   string `json:"status"`
	Response struct {
		Items []struct {
			ForwardID int    `json:"forwardId"`
			PostID    int64  `json:"postId"`
			PostLink  string `json:"postLink"`
			PostDate  int    `json:"postDate"`
			ChannelID int    `json:"channelId"`
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

type ChannelForwards struct {
	Status   string `json:"status"`
	Response struct {
		Items []struct {
			ForwardID int    `json:"forwardId"`
			PostID    int64  `json:"postId"`
			PostLink  string `json:"postLink"`
			PostDate  int    `json:"postDate"`
			ChannelID int    `json:"channelId"`
		} `json:"items"`
	} `json:"response"`
}

type ChannelSubscribers struct {
	Status   string `json:"status"`
	Response struct {
		Items []struct {
			Period            string `json:"period"`
			ParticipantsCount int    `json:"participants_count"`
		} `json:"items"`
	} `json:"response"`
}

type ChannelViews struct {
	Status   string `json:"status"`
	Response struct {
		Items []struct {
			Period     string `json:"period"`
			ViewsCount int    `json:"views_count"`
		} `json:"items"`
	} `json:"response"`
}

type ChannelAvgReach struct {
	Status   string `json:"status"`
	Response []struct {
		Period        string `json:"period"`
		AvgPostsReach int    `json:"avg_posts_reach"`
	} `json:"response"`
}

type ChannelErr struct {
	Status   string `json:"status"`
	Response []struct {
		Period string `json:"period"`
		Err    int    `json:"err"`
	} `json:"response"`
}

type ChannelAddPendingResponse struct {
	Status string `json:"status"`
}

type ChannelAddSuccessResponse struct {
	Status   string `json:"status"`
	Response []struct {
		ChannelId string `json:"channelId"`
	} `json:"response"`
}
