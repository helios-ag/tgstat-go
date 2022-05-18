package schema

type WordsMentionsResponseItem struct {
	Period        string `json:"period"`
	MentionsCount int    `json:"mentions_count"`
	ViewsCount    int    `json:"views_count"`
}

type WordsMentionsResponse struct {
	Items []WordsMentionsResponseItem `json:"items"`
}

type WordsMentions struct {
	Status   string                `json:"status"`
	Response WordsMentionsResponse `json:"response"`
}

type WordsMentionsByChannelItem struct {
	ChannelID       int `json:"channel_id"`
	MentionsCount   int `json:"mentions_count"`
	ViewsCount      int `json:"views_count"`
	LastMentionDate int `json:"last_mention_date"`
}

type WordsMentionsByChannelChannel struct {
	ID                int    `json:"id"`
	Link              string `json:"link"`
	Username          string `json:"username"`
	Title             string `json:"title"`
	About             string `json:"about"`
	Image100          string `json:"image100"`
	Image640          string `json:"image640"`
	ParticipantsCount int    `json:"participants_count"`
}

type WordsMentionsByChannelResponse struct {
	Items    []WordsMentionsByChannelItem    `json:"items"`
	Channels []WordsMentionsByChannelChannel `json:"channels"`
}

type WordsMentionsByChannel struct {
	Status   string                         `json:"status"`
	Response WordsMentionsByChannelResponse `json:"response"`
}
