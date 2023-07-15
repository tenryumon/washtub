package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/nyelonong/boilerplate-go/core/database"
	"github.com/nyelonong/boilerplate-go/core/format"
	"github.com/nyelonong/boilerplate-go/internal/entities"
)

var (
	queryInsertUser = `INSERT INTO users (name, email, phone, photo, status, user_type, last_action_time, timezone, created_by, updated_by) 
						VALUES (:name, :email, :phone, :photo, 0, :user_type, :last_action_time, :timezone, :created_by, :updated_by);`
	queryInsertPasswordRequest = `INSERT INTO password_requests (user_id, token, request_type, status, expired_time, created_by, updated_by) 
						VALUES (:user_id, :token, :request_type, 1, :expired_time, :created_by, :updated_by);`
)

func (us *UserObject) CreateNewUser(ctx context.Context, user entities.User) (entities.User, entities.PasswordRequest, error) {
	pr := entities.PasswordRequest{
		Token:       uuid.New().String(),
		RequestType: entities.PasswordRequestTypeNewUserPass,
		ExpiredTime: time.Now().Add(7 * 24 * time.Hour),
		CreatedBy:   user.CreatedBy,
		UpdatedBy:   user.UpdatedBy,
	}

	tx, err := us.db.Begin()
	if err != nil {
		return user, pr, fmt.Errorf("Failed to create db Tx because %s", err)
	}

	defer tx.Rollback()

	user, err = us.insertNewUser(ctx, tx, user)
	if err != nil {
		return user, pr, fmt.Errorf("Failed to insertNewUser because %s", err)
	}

	pr.UserID = user.ID
	err = us.insertNewPasswordRequest(ctx, tx, pr)
	if err != nil {
		err = fmt.Errorf("Failed to insertNewPasswordRequest because %s", err)
	}

	tx.Commit()

	return user, pr, err
}

func (us *UserObject) insertNewUser(ctx context.Context, tx *database.Tx, user entities.User) (entities.User, error) {
	user.Phone = format.Phone(user.Phone)
	user.LastActionTime = entities.ZeroUserLastActionTime
	result, err := tx.Exec(ctx, queryInsertUser, user)
	if err != nil {
		return user, err
	}

	user.ID, err = result.LastInsertId()
	return user, err
}

func (us *UserObject) insertNewPasswordRequest(ctx context.Context, tx *database.Tx, req entities.PasswordRequest) error {
	_, err := tx.Exec(ctx, queryInsertPasswordRequest, req)
	return err
}
