package role

import (
	"sync"
	"time"

	"github.com/nyelonong/boilerplate-go/core/database"
	"github.com/nyelonong/boilerplate-go/core/redis"
	"github.com/nyelonong/boilerplate-go/internal/entities"
	"github.com/nyelonong/boilerplate-go/internal/interfaces"
)

type Configuration struct {
	Database    *database.DB
	Redis       *redis.Redis
	ExpDuration time.Duration
}

type RoleObject struct {
	db          *database.DB
	cache       *redis.Redis
	expDuration time.Duration
	mtx         sync.Mutex
	byCode      map[string]int64
	byID        map[int64]RoleData
}

type RoleData struct {
	expiryTime time.Time
	data       entities.Role
}

func (dt RoleData) IsExpired() bool {
	return time.Now().After(dt.expiryTime)
}

func New(config Configuration) interfaces.Role {
	expDuration := 5 * time.Minute
	if config.ExpDuration > 0 {
		expDuration = config.ExpDuration * time.Second
	}

	return &RoleObject{
		db:          config.Database,
		cache:       config.Redis,
		expDuration: expDuration,
		byCode:      map[string]int64{},
		byID:        map[int64]RoleData{},
	}
}
