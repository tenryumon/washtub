package authorization

import (
	"github.com/nyelonong/boilerplate-go/internal/interfaces"
)

type Configuration struct {
	Session      interfaces.Session
	User         interfaces.User
	Role         interfaces.Role
	Notification interfaces.Notification
}

type Authorization struct {
	session      interfaces.Session
	user         interfaces.User
	role         interfaces.Role
	notification interfaces.Notification
}

func New(config Configuration) interfaces.AuthorizationUC {
	return &Authorization{
		session:      config.Session,
		user:         config.User,
		role:         config.Role,
		notification: config.Notification,
	}
}
