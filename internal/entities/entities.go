package entities

import "time"

type Aggregate struct {
	Count int64 `db:"count"`
}

const (
	MessageFailedLogin    = "Gagal masuk menggunakan kombinasi email dan password."
	MessageUserNotActive  = "Pengguna sudah tidak aktif."
	MessagePasswordNotSet = "Harap verifikasi email terlebih dahulu."

	DBTimestampFormat        = "2006-01-02 15:04:05"
	DBDateFormat             = "2006-01-02"
	DBTimeFormat             = "15:04"
	DBDefaultTimeStampFormat = time.RFC3339

	ActionBySystem = 0
)

type Pagination struct {
	Total int64
	Page  int64
	Limit int
}
