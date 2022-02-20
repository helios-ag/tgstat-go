package schema

type ErrorResponse struct {
	Status string `json:"status,string,omitempty"`
	Error  string `json:"error,string,omitempty"`
}

type Country = struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CountryResponse struct {
	Status   string    `json:"status"`
	Response []Country `json:"response"`
}

type Category = struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CategoryResponse struct {
	Status   string     `json:"status"`
	Response []Category `json:"response"`
}

type Language = struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type LanguageResponse struct {
	Status   string     `json:"status"`
	Response []Language `json:"response"`
}
