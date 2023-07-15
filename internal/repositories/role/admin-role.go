package role

import (
	"context"
	"strings"

	"github.com/nyelonong/boilerplate-go/internal/entities"
)

var (
	querySelectAdminRole      = `SELECT ar.id, ar.role_id, ar.created_by, ar.created_time, ar.updated_by, ar.updated_time`
	querySelectAdminRole_Role = `, rl.code`
	queryFromAdminRole        = " FROM admin_roles ar "
	queryJoinAdminRole_Role   = " INNER JOIN roles rl ON ar.role_id = rl.id "

	queryAdminRoleWhereUserID      = "ar.user_id = :user_id"
	queryAdminRoleWhereStatus      = "ar.status = :status"
	queryAdminRole_RoleWhereStatus = "rl.status = :role_status"
)

func (ro *RoleObject) GetAdminRoles(ctx context.Context, userID int64) (roles []entities.AdminRole, err error) {
	query := querySelectAdminRole + querySelectAdminRole_Role + queryFromAdminRole + queryJoinAdminRole_Role
	param := map[string]interface{}{
		"user_id":     userID,
		"status":      entities.AdminRoleStatusActive,
		"role_status": entities.RoleStatusActive,
	}

	where := []string{queryAdminRoleWhereUserID, queryAdminRoleWhereStatus, queryAdminRole_RoleWhereStatus}
	query += " WHERE " + strings.Join(where, " AND ")

	result := []entities.AdminRole{}
	err = ro.db.Select(ctx, &result, query, param)
	return result, err
}
