package upload

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/nyelonong/boilerplate-go/core/storage"
	"github.com/nyelonong/boilerplate-go/internal/entities"
)

var (
	queryInsertFiles = `INSERT INTO upload_files (token, file_url, file_type, status, created_by, updated_by) 
						VALUES (:token, :file_url, :file_type, 0, :created_by, :updated_by);`
)

func (uo *UploadObject) UploadFile(ctx context.Context, param entities.UploadFile) (entities.UploadFile, error) {
	token := uuid.New().String()

	// Upload File into somewhere and get URL
	// url := "https://www.unicef.org/indonesia/sites/unicef.org.indonesia/files/styles/hero_desktop/public/IDN-Education-Hero.jpg"
	if param.OrgCode == "" {
		param.OrgCode = "public"
	}

	object := storage.Object{
		Key:     fmt.Sprintf("%s/%s/%s/%s%s", param.GetFileTypeString(), param.OrgCode, time.Now().Format("2006/01/02"), token, filepath.Ext(param.FileName)),
		Content: param.FileContent,
		Type:    param.FileContentType,
	}

	if err := uo.storage.UploadFile(ctx, object); err != nil {
		return param, err
	}

	param.Token = token
	param.FileURL = object.Key
	_, err := uo.db.Exec(ctx, queryInsertFiles, param)
	if err != nil {
		return param, err
	}

	return param, nil
}

var (
	querySelectUploadFile     = "SELECT token, file_url, file_type, status, created_by, created_time, updated_by, updated_time"
	queryGetUploadFileByToken = querySelectUploadFile + " FROM upload_files WHERE token = :token"
)

func (uo *UploadObject) GetFile(ctx context.Context, token string) (entities.UploadFile, error) {
	result := entities.UploadFile{}
	param := map[string]interface{}{
		"token": token,
	}

	err := uo.db.Get(ctx, &result, queryGetUploadFileByToken, param)
	return result, err
}

var (
	queryUpdateFileStatus = "UPDATE upload_files SET status = :status, updated_by = :updated_by, updated_time = CURRENT_TIMESTAMP WHERE token = :token"
)

func (uo *UploadObject) UpdateFileStatus(ctx context.Context, token string, status int, actionBy int64) error {
	param := map[string]interface{}{
		"token":      token,
		"status":     status,
		"updated_by": actionBy,
	}

	_, err := uo.db.Exec(ctx, queryUpdateFileStatus, param)
	return err
}

var (
	queryRemoveFile = "UPDATE upload_files SET status = -1, updated_by = :updated_by, updated_time = CURRENT_TIMESTAMP WHERE file_url = :file_url"
)

func (uo *UploadObject) RemoveFile(ctx context.Context, file_url string, actionBy int64) error {
	param := map[string]interface{}{
		"file_url":   file_url,
		"updated_by": actionBy,
	}

	_, err := uo.db.Exec(ctx, queryRemoveFile, param)
	return err
}
