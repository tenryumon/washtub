package user

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/nyelonong/boilerplate-go/internal/entities"
)

var (
	keyLoginOtpRedisKey = "user_login_otp:id_%d"
	keyLoginAttemptKey  = "user_login_attempt:id_%d"
	maxLoginAttempt     = 5
)

func getLoginOtpRedisKey(id int64) string {
	return fmt.Sprintf(keyLoginOtpRedisKey, id)
}

func getLoginAttemptRedisKey(id int64) string {
	return fmt.Sprintf(keyLoginAttemptKey, id)
}

func (us *UserObject) GenerateLoginOTP(ctx context.Context, id int64) (string, error) {
	otp, err := us.getLoginOTP(ctx, id)
	if err == nil && otp != "" {
		// Previous OTP is still valid
		return otp, nil
	}

	otp = fmt.Sprint(time.Now().Nanosecond())[:6]
	expDuration := 10 * 60 * 60 // 10 minutes
	err = us.cache.SetExpire(ctx, getLoginOtpRedisKey(id), otp, int64(expDuration))
	return otp, err
}

func (us *UserObject) CheckLoginOTP(ctx context.Context, id int64, otp string) error {
	// Always Increase login attempt if not success
	success := false
	defer func() {
		if !success {
			us.increaseLoginAttempt(ctx, id)
		}
	}()

	token, err := us.getLoginOTP(ctx, id)
	if err != nil {
		return err
	}
	if token == "" {
		return entities.FailureLoginOtpExpired
	}

	// If failed get attempt, it's ok to bypass the check
	attempt, _ := us.getLoginAttempt(ctx, id)
	if attempt > maxLoginAttempt {
		return entities.FailureLoginOtpLimitExceed
	}
	if token != otp {
		return entities.FailureLoginOtpWrong
	}

	success = true
	return nil
}

func (us *UserObject) RemoveLoginData(ctx context.Context, id int64) error {
	us.removeLoginOTP(ctx, id)
	us.removeLoginAttempt(ctx, id)
	return nil
}

func (us *UserObject) getLoginAttempt(ctx context.Context, id int64) (int, error) {
	result, err := us.cache.Get(ctx, getLoginAttemptRedisKey(id))
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(result)
}

func (us *UserObject) increaseLoginAttempt(ctx context.Context, id int64) error {
	_, err := us.cache.Incr(ctx, getLoginAttemptRedisKey(id))
	return err
}

func (us *UserObject) removeLoginAttempt(ctx context.Context, id int64) error {
	return us.cache.Del(ctx, getLoginAttemptRedisKey(id))
}

func (us *UserObject) getLoginOTP(ctx context.Context, id int64) (string, error) {
	return us.cache.Get(ctx, getLoginOtpRedisKey(id))
}

func (us *UserObject) removeLoginOTP(ctx context.Context, id int64) error {
	return us.cache.Del(ctx, getLoginOtpRedisKey(id))
}
