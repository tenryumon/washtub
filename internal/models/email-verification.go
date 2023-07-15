package models

type EmailVerifReq struct {
	Email    string
	Token    string
	Password string
}

type EmailVerifResp struct {
	Errors
	IsNewUser  bool `json:"is_new_user"`
}
