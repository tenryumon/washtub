package session

import (
	"time"

	"github.com/nyelonong/boilerplate-go/core/database"
	"github.com/nyelonong/boilerplate-go/core/redis"
	"github.com/nyelonong/boilerplate-go/internal/interfaces"
)

type Configuration struct {
	Database *database.DB
	Redis    *redis.Redis

	ShortExpDuration time.Duration
	LongExpDuration  time.Duration
}

type SessionObject struct {
	db    *database.DB
	cache *redis.Redis

	shortExpDuration time.Duration
	longExpDuration  time.Duration
}

func New(config Configuration) interfaces.Session {
	shortExpDuration := 1 * 24 * time.Hour
	if config.ShortExpDuration > 0 {
		shortExpDuration = config.ShortExpDuration
	}

	longExpDuration := 7 * 24 * time.Hour
	if config.LongExpDuration > 0 {
		longExpDuration = config.LongExpDuration
	}

	return &SessionObject{
		db:    config.Database,
		cache: config.Redis,

		shortExpDuration: shortExpDuration,
		longExpDuration:  longExpDuration,
	}
}
