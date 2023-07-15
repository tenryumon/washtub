package interfaces

import (
	"context"

	"github.com/nyelonong/boilerplate-go/internal/entities"
)

type User interface {
	GetUserByID(ctx context.Context, id int64) (entities.User, error)
	GetUserByEmail(ctx context.Context, email string, userType int) (entities.User, error)
	GetUserByPhone(ctx context.Context, phone string, userType int) (entities.User, error)
	SearchUser(ctx context.Context, param map[string]interface{}) ([]entities.User, error)
	PaginationUser(ctx context.Context, param map[string]interface{}) ([]entities.User, entities.Pagination, error)
	GetUserDuplicateUpdateEmail(ctx context.Context, id int64, email string, userType int) (entities.User, error)
	GetUserDuplicateUpdatePhone(ctx context.Context, id int64, phone string, userType int) (entities.User, error)

	GetLoginInfo(ctx context.Context, userID int64) (entities.Login, error)
	CheckPassword(ctx context.Context, id int64, password string) (bool, bool, error)
	GenerateLoginOTP(ctx context.Context, id int64) (string, error)
	CheckLoginOTP(ctx context.Context, id int64, otp string) error
	RemoveLoginData(ctx context.Context, id int64) error
	UpdateLastAction(ctx context.Context, id int64) error

	CreateNewUser(ctx context.Context, user entities.User) (entities.User, entities.PasswordRequest, error)
	UpdateUserInfo(ctx context.Context, user entities.User) error
	IsValidPasswordToken(ctx context.Context, token string, email string) (entities.ValidPasswordRequest, error)
	SetPassword(ctx context.Context, token string, email string, password string) (bool, bool, entities.ValidPasswordRequest, error)
	UpdatePassword(ctx context.Context, login entities.Login, newPassword string) error

	CreateNewPasswordRequest(ctx context.Context, param entities.PasswordRequest) (entities.PasswordRequest, error)
	GetPasswordRequestByToken(ctx context.Context, token string) (entities.PasswordRequest, error)
}
