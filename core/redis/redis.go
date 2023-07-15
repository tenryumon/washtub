package redis

import (
	"context"
	"time"

	goredis "github.com/go-redis/redis/v8"
)

type Config struct {
	Address      string
	Password     string
	WriteTimeout int
	ReadTimeout  int
	DialTimeout  int
}

type Redis struct {
	client *goredis.Client
}

func New(conf Config) (*Redis, error) {
	option := goredis.Options{
		Addr:         conf.Address,
		Password:     conf.Password,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		DialTimeout:  5 * time.Second,
	}

	if conf.WriteTimeout > 0 {
		option.WriteTimeout = time.Duration(conf.WriteTimeout) * time.Second
	}
	if conf.ReadTimeout > 0 {
		option.ReadTimeout = time.Duration(conf.ReadTimeout) * time.Second
	}
	if conf.DialTimeout > 0 {
		option.DialTimeout = time.Duration(conf.DialTimeout) * time.Second
	}

	rdb := &Redis{client: goredis.NewClient(&option)}
	return rdb, rdb.Ping()
}

func (r *Redis) Ping() error {
	_, err := r.client.Ping(context.Background()).Result()
	return err
}

func (r *Redis) Do(ctx context.Context, args ...interface{}) (interface{}, error) {
	result, err := r.client.Do(ctx, args...).Result()
	if err != goredis.Nil {
		return result, err
	}
	return result, nil
}

func (r *Redis) Set(ctx context.Context, key, value string) error {
	return r.SetExpire(ctx, key, value, 15*60)
}

func (r *Redis) SetExpire(ctx context.Context, key, value string, second int64) error {
	if second <= 0 {
		second = 0
	}
	_, err := r.client.Set(ctx, key, value, time.Duration(second)*time.Second).Result()
	return err
}

func (r *Redis) Expire(ctx context.Context, key string, second int64) error {
	if second <= 0 {
		second = 0
	}
	_, err := r.client.Expire(ctx, key, time.Duration(second)*time.Second).Result()
	return err
}

func (r *Redis) Incr(ctx context.Context, key string) (int64, error) {
	result, err := r.client.Incr(ctx, key).Result()
	if err != goredis.Nil {
		return result, err
	}
	return result, nil
}

func (r *Redis) Decr(ctx context.Context, key string) (int64, error) {
	result, err := r.client.Decr(ctx, key).Result()
	if err != goredis.Nil {
		return result, err
	}
	return result, nil
}

func (r *Redis) Del(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	if err != goredis.Nil {
		return err
	}
	return nil
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != goredis.Nil {
		return result, err
	}
	return result, nil
}

func (r *Redis) SAdd(ctx context.Context, key string, members ...interface{}) error {
	_, err := r.client.SAdd(ctx, key, members...).Result()
	if err != goredis.Nil {
		return err
	}
	return nil
}

func (r *Redis) SRem(ctx context.Context, key string, members ...interface{}) error {
	_, err := r.client.SRem(ctx, key, members...).Result()
	if err != goredis.Nil {
		return err
	}
	return nil
}

func (r *Redis) SMembers(ctx context.Context, key string) ([]string, error) {
	result, err := r.client.SMembers(ctx, key).Result()
	if err != goredis.Nil {
		return result, err
	}
	return result, nil
}

func (r *Redis) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	result, err := r.client.SIsMember(ctx, key, member).Result()
	if err != goredis.Nil {
		return result, err
	}
	return result, nil
}

func (r *Redis) SCard(ctx context.Context, key string) (int64, error) {
	result, err := r.client.SCard(ctx, key).Result()
	if err != goredis.Nil {
		return result, err
	}
	return result, nil
}

func (r *Redis) LPush(ctx context.Context, key string, values ...interface{}) error {
	_, err := r.client.LPush(ctx, key, values...).Result()
	if err != goredis.Nil {
		return err
	}
	return nil
}

func (r *Redis) RPush(ctx context.Context, key string, values ...interface{}) error {
	_, err := r.client.RPush(ctx, key, values...).Result()
	if err != goredis.Nil {
		return err
	}
	return nil
}

func (r *Redis) LPop(ctx context.Context, key string) error {
	_, err := r.client.LPop(ctx, key).Result()
	if err != goredis.Nil {
		return err
	}
	return nil
}

func (r *Redis) RPop(ctx context.Context, key string) error {
	_, err := r.client.RPop(ctx, key).Result()
	if err != goredis.Nil {
		return err
	}
	return nil
}

func (r *Redis) LRem(ctx context.Context, key string, count int64, value interface{}) error {
	_, err := r.client.LRem(ctx, key, count, value).Result()
	if err != goredis.Nil {
		return err
	}
	return nil
}

func (r *Redis) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	result, err := r.client.LRange(ctx, key, start, stop).Result()
	if err != goredis.Nil {
		return result, err
	}
	return result, nil
}

func (r *Redis) LLen(ctx context.Context, key string) (int64, error) {
	result, err := r.client.LLen(ctx, key).Result()
	if err != goredis.Nil {
		return result, err
	}
	return result, nil
}
