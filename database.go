package tgstat_go

type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CountryResult struct {
	Status   string    `json:"status"`
	Response []Country `json:"response"`
}

type Category struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CategoryResult struct {
	Status   string     `json:"status"`
	Response []Category `json:"response"`
}

type Language struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type LanguageResult struct {
	Status   string     `json:"status"`
	Response []Language `json:"response"`
}
