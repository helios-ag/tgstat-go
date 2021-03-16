package schema

type WordsMentionsResponse struct {
	Status   string `json:"status"`
	Response struct {
		Items []struct {
			Period        string `json:"period"`
			MentionsCount int    `json:"mentions_count"`
			ViewsCount    int    `json:"views_count"`
		} `json:"items"`
	} `json:"response"`
}

type WordsMentionsByChannelResponse struct {
	Status   string `json:"status"`
	Response struct {
		Items []struct {
			ChannelID       int `json:"channel_id"`
			MentionsCount   int `json:"mentions_count"`
			ViewsCount      int `json:"views_count"`
			LastMentionDate int `json:"last_mention_date"`
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