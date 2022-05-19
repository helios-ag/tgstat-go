package schema

type SetCallbackSuccessResult struct {
	Status string `json:"status"`
}

type SetCallbackVerificationResult struct {
	Status     string `json:"status"`
	Error      string `json:"error"`
	VerifyCode string `json:"verify_code"`
}

type GetCallbackResponse struct {
	Status   string           `json:"status"`
	Response CallbackResponse `json:"response"`
}

type CallbackResponse struct {
	Url                string `json:"url"`
	PendingUpdateCount int    `json:"pending_update_count"`
	LastErrorDate      int    `json:"last_error_date"`
	LastErrorMessage   string `json:"last_error_message"`
}

type Subscribe struct {
	Status   string            `json:"status"`
	Response SubscribeResponse `json:"response"`
}

type SubscribeResponse struct {
	SubscriptionId int `json:"subscription_id"`
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
	Channel        Channel  `json:"channel,omitempty"`
	CreatedAt      int      `json:"created_at"`
	Keyword        Keyword  `json:"keyword,omitempty"`
}

type Keyword struct {
	Q              string `json:"q"`
	StrongSearch   bool   `json:"strong_search"`
	MinusWords     string `json:"minus_words"`
	ExtendedSyntax bool   `json:"extended_syntax"`
	PeerTypes      string `json:"peer_types"`
}
