package interfaces

import (
	"context"

	"github.com/nyelonong/boilerplate-go/internal/entities"
)

type Session interface {
	CreateSession(ctx context.Context, user entities.User, rememberMe bool) (entities.Session, error)
	DeleteSession(ctx context.Context, token string) error
	GetSession(ctx context.Context, token string) (int64, error)
	ExtendSession(ctx context.Context, sessionToken string, rememberMe bool) (entities.Session, error)
}
