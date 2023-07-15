package models

import (
	"time"
)

const (
	UserTypeAdmin = 1
	UserTypeStaff = 2
)

type DoLoginReq struct {
	Email      string
	Password   string
	UserType   int
	RememberMe bool
}

type DoLoginResp struct {
	Errors `json:"-"`

	Token    string    `json:"-"`
	ExpireAt time.Time `json:"-"`
	Name     string    `json:"name"`
	UserID   int64     `json:"user_id"`
}

type CheckLoginReq struct {
	Token       string
	LoginDevice int
}

type CheckLoginResp struct {
	UserID   int64
	OrgID    int64
	Timezone string
}

type DoLoginFromAppResp struct {
	Errors `json:"-"`

	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type DoRefreshAccessTokenReq struct {
	RefreshToken string
}
