package staff_dash

import (
	"github.com/nyelonong/boilerplate-go/internal/interfaces"
)

type Configuration struct {
	User         interfaces.User
	Role         interfaces.Role
	Upload       interfaces.Upload
	Notification interfaces.Notification
	Export       interfaces.Export
}

type StaffDashboard struct {
	user         interfaces.User
	role         interfaces.Role
	upload       interfaces.Upload
	notification interfaces.Notification
	export       interfaces.Export
}

func New(config Configuration) interfaces.StaffDashboardUC {
	return &StaffDashboard{
		user:         config.User,
		role:         config.Role,
		upload:       config.Upload,
		notification: config.Notification,
		export:       config.Export,
	}
}
