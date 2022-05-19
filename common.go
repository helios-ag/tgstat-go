package tgstat_go

type ErrorResult struct {
	Status     string `json:"status"`
	Error      string `json:"error"`
	VerifyCode string `json:"verify_code"`
}

type SuccessResult struct {
	Status string `json:"status"`
}
