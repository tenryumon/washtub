package entities

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

const (
	UserTypeAdmin = 1
	UserTypeStaff = 2

	UserInactive = -1
	UserPending  = 0
	UserActive   = 1
)

var (
	FailureLoginOtpExpired     = errors.New("LOGIN_OTP_EXPIRED")
	FailureLoginOtpWrong       = errors.New("LOGIN_OTP_WRONG")
	FailureLoginOtpLimitExceed = errors.New("LOGIN_OTP_LIMIT_EXCEED")
	ZeroUserLastActionTime     = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	TimeZoneCodeDefault        = "Asia/Jakarta"

	userStatusStringList = map[int]string{
		UserInactive: "Inactive",
		UserPending:  "Pending",
		UserActive:   "Active",
	}

	userTypeStringList = map[int]string{
		UserTypeAdmin: "Admin",
		UserTypeStaff: "Staff",
	}
)

type User struct {
	ID             int64     `db:"id"`
	Name           string    `db:"name"`
	Email          string    `db:"email"`
	Phone          string    `db:"phone"`
	Photo          string    `db:"photo"`
	Status         int       `db:"status"`
	UserType       int       `db:"user_type"`
	LastActionTime time.Time `db:"last_action_time"`
	CreatedBy      int64     `db:"created_by"`
	CreatedTime    string    `db:"created_time"`
	UpdatedBy      int64     `db:"updated_by"`
	UpdatedTime    string    `db:"updated_time"`
}

func (u User) Exist() bool {
	return u.ID != 0
}

func (u User) IsActive() bool {
	return u.Status == UserActive
}
func (u User) IsInactive() bool {
	return u.Status == UserInactive
}

func (u User) IsAdmin() bool {
	return u.UserType == UserTypeAdmin
}

func (u User) IsStaff() bool {
	return u.UserType == UserTypeStaff
}

func (u User) GetStatusString() string {
	if value, ok := userStatusStringList[u.Status]; ok {
		return value
	}
	return "Unknown"
}

func IsValidUserStatus(status int) bool {
	_, ok := userStatusStringList[status]
	return ok
}

type NewPasswordToken struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (np NewPasswordToken) GetBase64() (string, error) {
	tokenStr, err := json.Marshal(np)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(tokenStr), nil
}
