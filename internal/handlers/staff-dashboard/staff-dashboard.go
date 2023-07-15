package staff_dashboard

import (
	"github.com/nyelonong/boilerplate-go/internal/entities"
	"github.com/nyelonong/boilerplate-go/internal/interfaces"
)

type Configuration struct {
	DashboardUC     interfaces.StaffDashboardUC
	AuthorizationUC interfaces.AuthorizationUC
}

type StaffDashboard struct {
	dashboardUC     interfaces.StaffDashboardUC
	authorizationUC interfaces.AuthorizationUC

	cookieName       string
	autoExtendCookie string
	userType         int
}

func New(config Configuration) *StaffDashboard {
	return &StaffDashboard{
		dashboardUC:     config.DashboardUC,
		authorizationUC: config.AuthorizationUC,

		cookieName:       "APPSID",
		autoExtendCookie: "extExpiry",
		userType:         entities.UserTypeStaff,
	}
}
