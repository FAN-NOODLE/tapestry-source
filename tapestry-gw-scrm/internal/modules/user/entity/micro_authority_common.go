package entity

type AuthorityCheckSmsCodeRequest struct {
	Code   string `json:"code"`
	Event  string `json:"event"`
	Mobile string `json:"mobile"`
}
