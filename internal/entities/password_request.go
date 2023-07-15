package entities

import "time"

const (
	PasswordRequestDeleted = 0
	PasswordRequestActive  = 1
	PasswordRequestFinish  = 2

	PasswordRequestTypeNewUserPass    = 1
	PasswordRequestTypeForgotPass     = 2
	PasswordRequestTypeStaffInfoEmail = 3
	PasswordRequestTypeStaffInfoPhone = 4
)

var (
	passwordRequestTypeStringList = map[int]string{
		PasswordRequestTypeNewUserPass:    "NewUserPass",
		PasswordRequestTypeForgotPass:     "ForgotPass",
		PasswordRequestTypeStaffInfoEmail: "StaffInfoEmail",
		PasswordRequestTypeStaffInfoPhone: "StaffInfoPhone",
	}
)

type PasswordRequest struct {
	ID          int64     `db:"id"`
	UserID      int64     `db:"user_id"`
	Token       string    `db:"token"`
	RequestType int       `db:"request_type"`
	ExpiredTime time.Time `db:"expired_time"`
	Status      int       `db:"status"`
	CreatedBy   int64     `db:"created_by"`
	CreatedTime string    `db:"created_time"`
	UpdatedBy   int64     `db:"updated_by"`
	UpdatedTime string    `db:"updated_time"`
}

func (pr PasswordRequest) Exist() bool {
	return pr.ID != 0
}

func (pr PasswordRequest) IsActive() bool {
	return pr.Status == PasswordRequestActive
}

func (pr PasswordRequest) IsFinished() bool {
	return pr.Status == PasswordRequestFinish
}

func (pr PasswordRequest) IsExpired() bool {
	return time.Now().After(pr.ExpiredTime)
}

func IsValidPasswordRequestType(requestType int) bool {
	_, ok := passwordRequestTypeStringList[requestType]
	return ok
}

const (
	ValidPassRequestExpired  = 1
	ValidPassRequestNotFound = 2
	ValidPassRequestInactive = 3
)

type ValidPasswordRequest struct {
	Valid       bool
	InvalidType int
	PassReq     PasswordRequest
}

func (vpr ValidPasswordRequest) IsExpired() bool {
	return vpr.InvalidType == ValidPassRequestExpired
}

func (vpr ValidPasswordRequest) IsInvalid() bool {
	return vpr.InvalidType == ValidPassRequestNotFound || vpr.InvalidType == ValidPassRequestInactive
}
