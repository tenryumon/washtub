package interfaces

import (
	"context"

	"github.com/nyelonong/boilerplate-go/internal/entities"
)

type Upload interface {
	UploadFile(ctx context.Context, param entities.UploadFile) (entities.UploadFile, error)
	GetFile(ctx context.Context, token string) (entities.UploadFile, error)
	UpdateFileStatus(ctx context.Context, token string, status int, actionBy int64) error
	RemoveFile(ctx context.Context, file_url string, actionBy int64) error
	GenerateUserContentURL(ctx context.Context, relativePath string) string
}
