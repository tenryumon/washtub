package models

type NewPasswordReq struct {
	Token    string
	Password string
}

type NewPasswordResp struct {
	Errors `json:"-"`

	IsNewUser  bool `json:"is_new_user"`
}
type LoginForgotPasswordResp struct {
	Errors
}
