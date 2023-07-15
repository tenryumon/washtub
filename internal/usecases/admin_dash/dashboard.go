package admin_dash

import (
	"github.com/nyelonong/boilerplate-go/internal/interfaces"
)

type Configuration struct {
	User         interfaces.User
	Notification interfaces.Notification
}

type AdminDashboard struct {
	user         interfaces.User
	notification interfaces.Notification
}

func New(config Configuration) interfaces.AdminDashboardUC {
	return &AdminDashboard{
		user:         config.User,
		notification: config.Notification,
	}
}
