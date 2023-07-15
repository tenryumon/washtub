package entities

const (
	LoginTypePassword = 1
)

const (
	LoginDeviceDesktop   = 1
	LoginDeviceMobileApp = 2
)

type Login struct {
	ID          int64  `db:"id"`
	UserID      int64  `db:"user_id"`
	LoginType   int    `db:"login_type"`
	Token       string `db:"token"`
	Status      int    `db:"status"`
	CreatedBy   int64  `db:"created_by"`
	CreatedTime string `db:"created_time"`
	UpdatedBy   int64  `db:"updated_by"`
	UpdatedTime string `db:"updated_time"`
}

func (l Login) Exist() bool {
	return l.ID != 0
}

func (l Login) IsActive() bool {
	return l.Status == 1
}

func (l Login) IsPassword() bool {
	return l.LoginType == LoginTypePassword
}
