package user

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/nyelonong/boilerplate-go/internal/entities"
)

var (
	querySelectLogin              = "SELECT id, user_id, login_type, token, status, created_by, created_time, updated_by, updated_time"
	queryGetPasswordByUserID      = querySelectLogin + " FROM logins WHERE user_id = :user_id AND status = 1 AND login_type = :login_type ORDER BY id DESC LIMIT 1"
	queryUpdateUserLastActionTime = `UPDATE users SET last_action_time = :last_action_time, updated_by = :updated_by WHERE id = :id`
)

func (us *UserObject) GetLoginInfo(ctx context.Context, userID int64) (entities.Login, error) {
	return us.getLogin(ctx, userID)
}

func (us *UserObject) getLogin(ctx context.Context, userID int64) (entities.Login, error) {
	result := entities.Login{}
	param := map[string]interface{}{
		"user_id":    userID,
		"login_type": entities.LoginTypePassword,
	}

	err := us.db.Get(ctx, &result, queryGetPasswordByUserID, param)
	return result, err
}

func (us *UserObject) CheckPassword(ctx context.Context, id int64, password string) (validPass bool, existPass bool, err error) {
	login, err := us.getLogin(ctx, id)
	if err != nil {
		return validPass, existPass, err
	}
	if !login.Exist() {
		// User not create password yet
		return validPass, existPass, nil
	}
	existPass = true

	err = bcrypt.CompareHashAndPassword([]byte(login.Token), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		// Email & Password combination is wrong
		return validPass, existPass, nil
	} else if err != nil {
		return validPass, existPass, err
	}
	validPass = true

	return validPass, existPass, nil
}

func (us *UserObject) UpdateLastAction(ctx context.Context, id int64) error {
	param := map[string]interface{}{
		"id":               id,
		"last_action_time": time.Now(),
		"updated_by":       id,
	}
	_, err := us.db.Exec(ctx, queryUpdateUserLastActionTime, param)
	return err
}
