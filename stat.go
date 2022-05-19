package schema

type StatResponse struct {
	ServiceKey    string `json:"serviceKey"`
	Title         string `json:"title"`
	SpentChannels string `json:"spentChannels,omitempty"`
	SpentRequests string `json:"spentRequests"`
	ExpiredAt     int64  `json:"expiredAt"`
	SpentWords    string `json:"spentWords,omitempty"`
	SpentObjects  string `json:"spentObjects,omitempty"`
}

type StatResult struct {
	Status   string         `json:"status"`
	Response []StatResponse `json:"response"`
}
