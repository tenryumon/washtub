package entities

const (
	// Important Roles Only
	RoleOrgAdmin   = "org-admin"
	RoleOwner      = "org-owner"
	RoleSuperAdmin = "super-admin"

	RoleStatusActive      = 1
	AdminRoleStatusActive = 1
	StaffRoleStatusActive = 1
)

var (
	forbiddenRoles = []string{RoleOrgAdmin, RoleOwner, RoleSuperAdmin}
)

type UserRole struct {
	ID          int64  `db:"id"`
	PlacementID int64  `db:"placement_id"`
	UserID      int64  `db:"user_id"`
	RoleID      int64  `db:"role_id"`
	RoleCode    string `db:"code"`
	Status      int    `db:"status"`
	CreatedBy   int64  `db:"created_by"`
	CreatedTime string `db:"created_time"`
	UpdatedBy   int64  `db:"updated_by"`
	UpdatedTime string `db:"updated_time"`
}

func (ur UserRole) Exist() bool {
	return ur.ID != 0
}

func (ur UserRole) IsActive() bool {
	return ur.Status == 1
}

const (
	RoleTypeStaff = 1
	RoleTypeAdmin = 9
)

type Role struct {
	ID          int64  `db:"id"`
	Code        string `db:"code"`
	Name        string `db:"name"`
	RoleType    int    `db:"role_type"`
	Status      int    `db:"status"`
	CreatedBy   int64  `db:"created_by"`
	CreatedTime string `db:"created_time"`
	UpdatedBy   int64  `db:"updated_by"`
	UpdatedTime string `db:"updated_time"`
}

func (rl Role) Exist() bool {
	return rl.ID != 0
}

func (rl Role) IsActive() bool {
	return rl.Status == 1
}

func (rl Role) IsStaffRole() bool {
	return rl.RoleType == RoleTypeStaff
}

func (rl Role) IsAdminRole() bool {
	return rl.RoleType == RoleTypeAdmin
}

func IsForbiddenRole(role string) bool {
	for _, v := range forbiddenRoles {
		if v == role {
			return true
		}
	}
	return false
}

type InsertUserRole struct {
	PlacementID int64
	Roles       []int64
	ActionBy    int64
}

type AdminRole struct {
	ID          int64  `db:"id"`
	UserID      int64  `db:"user_id"`
	RoleID      int64  `db:"role_id"`
	RoleCode    string `db:"code"`
	Status      int    `db:"status"`
	CreatedBy   int64  `db:"created_by"`
	CreatedTime string `db:"created_time"`
	UpdatedBy   int64  `db:"updated_by"`
	UpdatedTime string `db:"updated_time"`
}

func (ar AdminRole) Exist() bool {
	return ar.ID != 0
}

func (ar AdminRole) IsActive() bool {
	return ar.Status == 1
}
