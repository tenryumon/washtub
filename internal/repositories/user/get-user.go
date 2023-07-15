package user

import (
	"context"

	"github.com/nyelonong/boilerplate-go/core/format"
	"github.com/nyelonong/boilerplate-go/internal/entities"
)

var (
	querySelectUser                  = "SELECT id, name, email, phone, photo, status, user_type, last_action_time, timezone, created_by, created_time, updated_by, updated_time"
	queryGetUserByID                 = querySelectUser + " FROM users WHERE id = :id"
	queryGetUserByEmail              = querySelectUser + " FROM users WHERE email = :email AND user_type = :user_type"
	queryGetUserByPhone              = querySelectUser + " FROM users WHERE phone = :phone AND user_type = :user_type"
	queryGetUserDuplicateUpdateEmail = querySelectUser + " FROM users WHERE email = :email AND user_type = :user_type AND id != :id LIMIT 1"
	queryGetUserDuplicateUpdatePhone = querySelectUser + " FROM users WHERE phone = :phone AND user_type = :user_type AND id != :id LIMIT 1"
)

func (us *UserObject) GetUserByID(ctx context.Context, id int64) (entities.User, error) {
	result := entities.User{}
	param := map[string]interface{}{
		"id": id,
	}

	err := us.db.Get(ctx, &result, queryGetUserByID, param)
	return result, err
}

func (us *UserObject) GetUserByEmail(ctx context.Context, email string, userType int) (entities.User, error) {
	result := entities.User{}
	param := map[string]interface{}{
		"email":     email,
		"user_type": userType,
	}

	err := us.db.Get(ctx, &result, queryGetUserByEmail, param)
	return result, err
}

func (us *UserObject) GetUserByPhone(ctx context.Context, phone string, userType int) (entities.User, error) {
	result := entities.User{}
	phone = format.Phone(phone)
	param := map[string]interface{}{
		"phone":     phone,
		"user_type": userType,
	}

	err := us.db.Get(ctx, &result, queryGetUserByPhone, param)
	return result, err
}

func (us *UserObject) GetUserDuplicateUpdateEmail(ctx context.Context, id int64, email string, userType int) (entities.User, error) {
	result := entities.User{}
	param := map[string]interface{}{
		"id":        id,
		"email":     email,
		"user_type": userType,
	}

	err := us.db.Get(ctx, &result, queryGetUserDuplicateUpdateEmail, param)
	return result, err
}

func (us *UserObject) GetUserDuplicateUpdatePhone(ctx context.Context, id int64, phone string, userType int) (entities.User, error) {
	result := entities.User{}
	phone = format.Phone(phone)
	param := map[string]interface{}{
		"id":        id,
		"phone":     phone,
		"user_type": userType,
	}

	err := us.db.Get(ctx, &result, queryGetUserDuplicateUpdatePhone, param)
	return result, err
}
