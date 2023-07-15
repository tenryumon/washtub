package entities

import (
	"time"
)

type Session struct {
	Token    string
	ExpireAt time.Time
}
