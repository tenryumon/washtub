package export

import (
	"context"

	"github.com/nyelonong/boilerplate-go/internal/entities"
)

var (
	queryInsertFiles = `INSERT INTO export_files (id, org_id, file_name, file_url, file_type, tag, status, created_by, updated_by) 
						VALUES (:id, :org_id, :file_name, :file_url, :file_type, :tag, :status, :created_by, :updated_by);`

	queryGetExportFileByID   = `SELECT * FROM export_files where id = :id`
	queryGetExportFilesByTag = `SELECT * FROM export_files where org_id = :org_id AND tag = :tag ORDER BY created_time DESC LIMIT 5`
)

func (eo *ExportObject) ExportFile(ctx context.Context, param entities.ExportFile) error {
	_, err := eo.db.Exec(ctx, queryInsertFiles, param)
	if err != nil {
		return err
	}

	return nil
}

func (eo *ExportObject) GetExportedFileByID(ctx context.Context, id string) (entities.ExportFile, error) {
	result := entities.ExportFile{}
	param := map[string]interface{}{
		"id": id,
	}

	err := eo.db.Get(ctx, &result, queryGetExportFileByID, param)
	return result, err
}

func (eo *ExportObject) GetExportedFilesByTag(ctx context.Context, orgID int64, tag string) ([]entities.ExportFile, error) {
	result := []entities.ExportFile{}
	param := map[string]interface{}{
		"org_id": orgID,
		"tag":    tag,
	}

	err := eo.db.Select(ctx, &result, queryGetExportFilesByTag, param)
	return result, err
}
