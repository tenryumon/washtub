package export

import (
	"context"
	"encoding/json"
	"fmt"

	"moul.io/progress"
)

const (
	keyRedisExportProgressTracker = "export:id_%s"
)

func (eo *ExportObject) TrackExportProgress(ctx context.Context, id string, progress []byte) error {
	return eo.cache.SetExpire(ctx, fmt.Sprintf(keyRedisExportProgressTracker, id), string(progress), 3600)
}

func (eo *ExportObject) GetExportProgress(ctx context.Context, id string) (*progress.Snapshot, error) {
	var snapshot progress.Snapshot

	val, err := eo.cache.Get(ctx, fmt.Sprintf(keyRedisExportProgressTracker, id))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(val), &snapshot); err != nil {
		return nil, err
	}

	return &snapshot, nil
}
