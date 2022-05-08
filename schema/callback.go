package schema

type SetCallbackResponse struct {
	Status string `json:"status,string"`
}

type GetCallbackResponse struct {
	Status   string           `json:"status,string"`
	Response CallbackResponse `json:"response"`
}

type CallbackResponse struct {
	Url                string `json:"url"`
	PendingUpdateCount int    `json:"pending_update_count"`
	LastErrorDate      int    `json:"last_error_date"`
	LastErrorMessage   string `json:"last_error_message"`
}

type SubscribeResponse struct {
	Status   string `json:"status"`
	Response struct {
		SubscriptionId int `json:"subscription_id"`
	} `json:"response"`
}

type SubscriptionList struct {
	Status   string                   `json:"status"`
	Response SubscriptionListResponse `json:"response"`
}

type SubscriptionListResponse struct {
	TotalCount    int            `json:"total_count"`
	Subscriptions []Subscription `json:"subscriptions"`
}

type Subscription struct {
	SubscriptionId int      `json:"subscription_id"`
	EventTypes     []string `json:"event_types"`
	Type           string   `json:"type"`
	Channel        struct {
		Id                int    `json:"id"`
		Link              string `json:"link"`
		Username          string `json:"username"`
		Title             string `json:"title"`
		About             string `json:"about"`
		Image100          string `json:"image100"`
		Image640          string `json:"image640"`
		ParticipantsCount int    `json:"participants_count"`
	} `json:"channel,omitempty"`
	CreatedAt int `json:"created_at"`
	Keyword   struct {
		Q              string `json:"q"`
		StrongSearch   bool   `json:"strong_search"`
		MinusWords     string `json:"minus_words"`
		ExtendedSyntax bool   `json:"extended_syntax"`
		PeerTypes      string `json:"peer_types"`
	} `json:"keyword,omitempty"`
}
