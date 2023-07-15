package user

import (
	"github.com/nyelonong/boilerplate-go/core/database"
	"github.com/nyelonong/boilerplate-go/core/redis"
	"github.com/nyelonong/boilerplate-go/internal/interfaces"
)

type Configuration struct {
	Database *database.DB
	Redis    *redis.Redis
}

type UserObject struct {
	db    *database.DB
	cache *redis.Redis
}

func New(config Configuration) interfaces.User {
	return &UserObject{
		db:    config.Database,
		cache: config.Redis,
	}
}
