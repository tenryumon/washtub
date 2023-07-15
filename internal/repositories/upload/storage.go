package upload

import "context"

func (uo *UploadObject) GenerateUserContentURL(ctx context.Context, relativePath string) string {
	return uo.storage.GenerateUserContentURL(ctx, relativePath)
}
