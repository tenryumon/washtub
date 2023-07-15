package staff_dash

import (
	"context"

	"github.com/nyelonong/boilerplate-go/core/log"
	"github.com/nyelonong/boilerplate-go/internal/entities"
	"github.com/nyelonong/boilerplate-go/internal/models"
)

func (sd *StaffDashboard) UploadImage(ctx context.Context, param models.UploadFileReq) (resp models.UploadFileResp) {
	// Initialize Usecase Parameter
	resp.InitError()
	ucContext := map[string]interface{}{
		"usecase":  "UploadImage",
		"staff_id": param.ActionBy,
	}

	if !entities.IsValidImageMime(param.FileExt) {
		resp.ErrorGeneral("Tipe gambar tidak valid.")
		return
	}
	if len(param.FileContent) == 0 {
		resp.ErrorGeneral("Gambar tidak valid.")
		return
	}
	if len(param.FileContent) > entities.MaxImageSize {
		resp.ErrorGeneral("Gambar terlalu besar.")
		return
	}

	image := entities.UploadFile{
		FileType:        entities.FileTypeImage,
		FileContent:     param.FileContent,
		FileName:        param.FileName,
		CreatedBy:       param.ActionBy,
		UpdatedBy:       param.ActionBy,
		FileContentType: param.FileExt,
	}
	file, err := sd.upload.UploadFile(ctx, image)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do UploadFile because %s", err)
		return
	}

	resp.Token = file.Token
	resp.FileURL = sd.upload.GenerateUserContentURL(ctx, file.FileURL)
	resp.Success()
	return
}
