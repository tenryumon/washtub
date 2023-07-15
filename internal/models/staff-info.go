package models

import (
	"github.com/nyelonong/boilerplate-go/internal/entities"
)

type GetPlacementResp struct {
	Errors

	Placements []StaffPlacement
	Info       entities.User
}

type StaffPlacement struct {
	RoleID   int64
	RoleName string
	Default  bool
	Accesses []string
}

type DropdownTimezoneListResp struct {
	Errors `json:"-"`

	List []DropdownTimezoneList `json:"list"`
}
type DropdownTimezoneList struct {
	Code    string `json:"code"`
	Country string `json:"country"`
	Utc     string `json:"utc"`
}
