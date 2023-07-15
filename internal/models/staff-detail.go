package models

type ShowStaffDetailReq struct {
	OrgID    int64
	ActionBy int64

	PlacementID int64 `json:"placement_id"`
}

type ShowStaffDetailResp struct {
	Errors `json:"-"`

	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Phone       string   `json:"phone"`
	Photo       string   `json:"photo"`
	Status      int      `json:"status"`
	RoleID      int64    `json:"role_id"`
	Accesses    []string `json:"accesses"`
	AccessNames []string `json:"access_names"`
}
