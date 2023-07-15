package upload

import (
	"github.com/nyelonong/boilerplate-go/core/database"
	"github.com/nyelonong/boilerplate-go/core/redis"
	"github.com/nyelonong/boilerplate-go/core/storage"
	"github.com/nyelonong/boilerplate-go/internal/interfaces"
)

type Configuration struct {
	Database *database.DB
	Redis    *redis.Redis
	Storage  *storage.Storage
}

type UploadObject struct {
	db      *database.DB
	cache   *redis.Redis
	storage *storage.Storage
}

func New(config Configuration) interfaces.Upload {
	return &UploadObject{
		db:      config.Database,
		cache:   config.Redis,
		storage: config.Storage,
	}
}
