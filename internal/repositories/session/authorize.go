package session

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/nyelonong/boilerplate-go/internal/entities"
)

var keyRedisSession = "session:token_%s"

func getSessionRedisKey(code string) string {
	return fmt.Sprintf(keyRedisSession, code)
}

func (so *SessionObject) CreateSession(ctx context.Context, user entities.User, rememberMe bool) (resp entities.Session, err error) {
	sessionToken := uuid.New().String()

	expDuration := so.shortExpDuration
	if rememberMe {
		expDuration = so.longExpDuration
	}
	expireMilliSeconds := expDuration.Milliseconds() / 100

	err = so.cache.SetExpire(ctx, getSessionRedisKey(sessionToken), fmt.Sprintf("%d", user.ID), expireMilliSeconds)
	if err != nil {
		return resp, err
	}

	resp.Token = sessionToken
	resp.ExpireAt = time.Now().Add(expDuration)

	return resp, err
}

func (so *SessionObject) DeleteSession(ctx context.Context, token string) error {
	return so.cache.Del(ctx, getSessionRedisKey(token))
}

func (so *SessionObject) GetSession(ctx context.Context, token string) (userID int64, err error) {
	uid, err := so.cache.Get(ctx, getSessionRedisKey(token))
	if err != nil {
		return userID, err
	}

	userID, err = strconv.ParseInt(uid, 10, 64)
	if err != nil {
		return userID, err
	}
	return userID, nil
}

func (so *SessionObject) ExtendSession(ctx context.Context, sessionToken string, rememberMe bool) (resp entities.Session, err error) {
	expDuration := so.shortExpDuration
	if rememberMe {
		expDuration = so.longExpDuration
	}
	expireMilliSeconds := expDuration.Milliseconds() / 100

	err = so.cache.Expire(ctx, getSessionRedisKey(sessionToken), expireMilliSeconds)
	if err != nil {
		return resp, err
	}

	resp.Token = sessionToken
	resp.ExpireAt = time.Now().Add(expDuration)

	return resp, err
}
