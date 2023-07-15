package admin_dashboard

import (
	"github.com/nyelonong/boilerplate-go/internal/interfaces"
)

type Configuration struct {
	DashboardUC     interfaces.AdminDashboardUC
	AuthorizationUC interfaces.AuthorizationUC
}

type AdminDashboard struct {
	dashboardUC     interfaces.AdminDashboardUC
	authorizationUC interfaces.AuthorizationUC

	cookieName       string
	autoExtendCookie string
}

func New(config Configuration) *AdminDashboard {
	return &AdminDashboard{
		dashboardUC:     config.DashboardUC,
		authorizationUC: config.AuthorizationUC,

		cookieName:       "APPSID",
		autoExtendCookie: "extExpiry",
	}
}
