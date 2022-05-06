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
