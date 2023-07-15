package interfaces

import (
	"context"

	"github.com/nyelonong/boilerplate-go/internal/entities"
	"moul.io/progress"
)

type Export interface {
	ExportFile(ctx context.Context, param entities.ExportFile) error
	TrackExportProgress(ctx context.Context, id string, progress []byte) error
	GetExportProgress(ctx context.Context, id string) (*progress.Snapshot, error)
	GetExportedFileByID(ctx context.Context, id string) (entities.ExportFile, error)
	GetExportedFilesByTag(ctx context.Context, orgID int64, tag string) ([]entities.ExportFile, error)
}
