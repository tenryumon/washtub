package role

import (
	"context"
	"fmt"

	"github.com/nyelonong/boilerplate-go/core/database"
	"github.com/nyelonong/boilerplate-go/internal/entities"
)

var queryGetAllActiveRole = `SELECT sr.id, sr.placement_id, sr.role_id, rl.code, sr.created_by, sr.created_time, sr.updated_by, sr.updated_time 
                            FROM staff_roles sr INNER JOIN roles rl ON sr.role_id = rl.id 
                            WHERE sr.placement_id = :placement_id AND sr.status = 1 AND rl.status = 1`

func (ro *RoleObject) GetStaffRole(ctx context.Context, placementID int64) (roles []string, err error) {
	result := []entities.UserRole{}
	param := map[string]interface{}{
		"placement_id": placementID,
	}

	err = ro.db.Select(ctx, &result, queryGetAllActiveRole, param)
	if err != nil {
		return roles, err
	}

	for _, v := range result {
		roles = append(roles, v.RoleCode)
	}
	return roles, nil
}

var queryUpsertStaffRole = `INSERT INTO staff_roles (placement_id, role_id, status, created_by, updated_by)
							VALUES (:placement_id, :role_id, 1, :action_by, :action_by)
							ON DUPLICATE KEY UPDATE
								status = 1, updated_by = :action_by`

func (ro *RoleObject) addStaffRole(ctx context.Context, tx *database.Tx, userRole entities.InsertUserRole) error {
	param := map[string]interface{}{
		"placement_id": userRole.PlacementID,
		"action_by":    userRole.ActionBy,
	}
	for _, row := range userRole.Roles {
		param["role_id"] = row

		_, err := tx.Exec(ctx, queryUpsertStaffRole, param)
		if err != nil {
			return fmt.Errorf("Failed to execute query for id %d because %s", row, err)
		}
	}

	return nil
}

func (ro *RoleObject) AddStaffRole(ctx context.Context, userRole entities.InsertUserRole) error {
	tx, err := ro.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to create db Tx because %s", err)
	}
	defer tx.Rollback()

	err = ro.addStaffRole(ctx, tx, userRole)
	if err != nil {
		return fmt.Errorf("Failed to addStaffRole because %s", err)
	}
	tx.Commit()

	return nil
}

var queryRemoveExistingStaffRole = `UPDATE staff_roles 
									SET status = -9, updated_by = :action_by
									WHERE placement_id = :placement_id`

func (ro *RoleObject) removeExistingRole(ctx context.Context, tx *database.Tx, userRole entities.InsertUserRole) error {
	param := map[string]interface{}{
		"placement_id": userRole.PlacementID,
		"action_by":    userRole.ActionBy,
	}
	_, err := tx.Exec(ctx, queryRemoveExistingStaffRole, param)
	if err != nil {
		return fmt.Errorf("Failed to execute query for id %d because %s", userRole.PlacementID, err)
	}

	return nil
}

func (ro *RoleObject) UpdateStaffRole(ctx context.Context, userRoles []entities.InsertUserRole) error {
	tx, err := ro.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to create db Tx because %s", err)
	}
	defer tx.Rollback()

	for _, userRole := range userRoles {
		err = ro.removeExistingRole(ctx, tx, userRole)
		if err != nil {
			return fmt.Errorf("Failed to removeExistingRole because %s", err)
		}

		err = ro.addStaffRole(ctx, tx, userRole)
		if err != nil {
			return fmt.Errorf("Failed to addStaffRole because %s", err)
		}
	}
	tx.Commit()

	return nil
}
