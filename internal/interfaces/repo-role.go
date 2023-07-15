package interfaces

import (
	"context"

	"github.com/nyelonong/boilerplate-go/internal/entities"
)

type Role interface {
	GetAllRoles(ctx context.Context, param map[string]interface{}) ([]entities.Role, error)
	GetRoleByID(ctx context.Context, id int64) (entities.Role, error)
	GetRoleByCode(ctx context.Context, code string) (entities.Role, error)
	GetRoleByCodes(ctx context.Context, codes []string) ([]entities.Role, error)

	GetStaffRole(ctx context.Context, placementID int64) (roles []string, err error)
	AddStaffRole(ctx context.Context, userRole entities.InsertUserRole) error
	UpdateStaffRole(ctx context.Context, userRoles []entities.InsertUserRole) error

	GetAdminRoles(ctx context.Context, userID int64) (roles []entities.AdminRole, err error)
}
