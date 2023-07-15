package authorization

import (
	"context"

	"github.com/nyelonong/boilerplate-go/core/log"
)

// GetRolesByUserID will get User roles
func (au *Authorization) GetRolesByUserID(ctx context.Context, userID int64) (result []string, err error) {
	// Initialize Usecase Parameter
	ucContext := map[string]interface{}{
		"usecase": "GetRolesByUserID",
		"user_id": userID,
	}

	// Start Usecase Action
	roles, err := au.role.GetAdminRoles(ctx, userID)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do GetAdminRoles because %s", err)
		return
	}

	result = []string{}
	for _, role := range roles {
		result = append(result, role.RoleCode)
	}

	return
}
