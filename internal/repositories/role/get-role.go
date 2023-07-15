package role

import (
	"context"
	"strings"
	"time"

	"github.com/nyelonong/boilerplate-go/internal/entities"
)

var querySelectRole = "SELECT id, name, code, status, role_type, status, created_by, created_time, updated_by, updated_time"
var queryGetRoleByCode = querySelectRole + " FROM roles WHERE code = :code"
var queryGetRoleByCodes = querySelectRole + " FROM roles WHERE code IN(:codes)"
var queryGetRoleByID = querySelectRole + " FROM roles WHERE id = :id"
var queryGetAllRoles = querySelectRole + " FROM roles"
var queryRoleWhereStatus = "status = :status"
var queryRoleWhereRoleType = "role_type like :role_type"
var queryRoleWhereExcludeOwner = "code != '" + entities.RoleOwner + "'"
var queryRoleWhereExcludeAdmin = "code != '" + entities.RoleOrgAdmin + "'"

func (ro *RoleObject) GetAllRoles(ctx context.Context, param map[string]interface{}) ([]entities.Role, error) {
	query := queryGetAllRoles

	where := []string{}
	if _, ok := param["status"]; ok {
		where = append(where, queryRoleWhereStatus)
	}
	if _, ok := param["role_type"]; ok {
		where = append(where, queryRoleWhereRoleType)
	}
	if value, ok := param["exclude_owner_role"]; ok && value == true {
		where = append(where, queryRoleWhereExcludeOwner)
	}
	if value, ok := param["exclude_admin_role"]; ok && value == true {
		where = append(where, queryRoleWhereExcludeAdmin)
	}
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	result := []entities.Role{}
	err := ro.db.Select(ctx, &result, query, param)
	return result, err
}

func (ro *RoleObject) GetRoleByCode(ctx context.Context, code string) (entities.Role, error) {
	roleID, exist := ro.byCode[code]
	if exist {
		data := ro.getRoleFromMemory(ctx, roleID)
		if data.Exist() {
			return data, nil
		}
	}

	data, err := ro.getRoleByCode(ctx, code)
	if err != nil {
		return data, err
	}

	if data.Exist() {
		ro.setRoleInMemory(ctx, data)
	}

	return data, nil
}

func (ro *RoleObject) GetRoleByCodes(ctx context.Context, codes []string) ([]entities.Role, error) {
	result := []entities.Role{}
	param := map[string]interface{}{
		"codes": codes,
	}

	err := ro.db.Select(ctx, &result, queryGetRoleByCodes, param)
	return result, err
}

func (ro *RoleObject) getRoleByCode(ctx context.Context, code string) (entities.Role, error) {
	result := entities.Role{}
	param := map[string]interface{}{
		"code": code,
	}

	err := ro.db.Get(ctx, &result, queryGetRoleByCode, param)
	return result, err
}

func (ro *RoleObject) GetRoleByID(ctx context.Context, id int64) (entities.Role, error) {
	data := ro.getRoleFromMemory(ctx, id)
	if data.Exist() {
		return data, nil
	}

	data, err := ro.getRoleByID(ctx, id)
	if err != nil {
		return data, err
	}

	if data.Exist() {
		ro.setRoleInMemory(ctx, data)
	}

	return data, nil
}

func (ro *RoleObject) getRoleByID(ctx context.Context, id int64) (entities.Role, error) {
	result := entities.Role{}
	param := map[string]interface{}{
		"id": id,
	}

	err := ro.db.Get(ctx, &result, queryGetRoleByID, param)
	return result, err
}

func (ro *RoleObject) getRoleFromMemory(ctx context.Context, id int64) entities.Role {
	value, exist := ro.byID[id]
	if !exist || value.IsExpired() {
		return entities.Role{}
	}

	return value.data
}

func (ro *RoleObject) setRoleInMemory(ctx context.Context, data entities.Role) {
	ro.mtx.Lock()
	ro.byCode[data.Code] = data.ID
	ro.byID[data.ID] = RoleData{
		expiryTime: time.Now().Add(ro.expDuration),
		data:       data,
	}
	ro.mtx.Unlock()
}
