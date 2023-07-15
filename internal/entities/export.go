package entities

import "time"

type ExportFile struct {
	ID          string    `db:"id"`
	OrgID       int64     `db:"org_id"`
	FileName    string    `db:"file_name"`
	FileURL     string    `db:"file_url"`
	FileType    string    `db:"file_type"`
	Tag         string    `db:"tag"`
	Status      string    `db:"status"`
	CreatedBy   string    `db:"created_by"`
	CreatedTime time.Time `db:"created_time"`
	UpdatedBy   string    `db:"updated_by"`
	UpdatedTime time.Time `db:"updated_time"`
}
