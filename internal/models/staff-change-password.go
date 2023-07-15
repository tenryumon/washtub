package models

type UpdateUserPasswordReq struct {
	OrgID    int64
	ActionBy int64

	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type UpdateUserPasswordResp struct {
	Errors `json:"-"`
}
