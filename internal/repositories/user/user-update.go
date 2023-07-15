package user

import (
	"context"

	"github.com/nyelonong/boilerplate-go/core/format"
	"github.com/nyelonong/boilerplate-go/internal/entities"
)

var (
	queryUpdateUserInfo = `UPDATE users SET name = :name, email = :email, phone = :phone, updated_by = :updated_by WHERE id = :id`
)

func (us *UserObject) UpdateUserInfo(ctx context.Context, user entities.User) error {
	user.Phone = format.Phone(user.Phone)
	_, err := us.db.Exec(ctx, queryUpdateUserInfo, user)
	return err
}
