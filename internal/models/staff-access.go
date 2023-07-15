package models

type DropdownStaffAccessListResp struct {
	Errors `json:"-"`

	List []DropdownStaffAccessList `json:"list"`
}
type DropdownStaffAccessList struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
