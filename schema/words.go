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

