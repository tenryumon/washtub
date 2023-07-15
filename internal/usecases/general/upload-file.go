package general_uc

import (
	"context"

	"github.com/nyelonong/boilerplate-go/core/log"
	"github.com/nyelonong/boilerplate-go/internal/entities"
	"github.com/nyelonong/boilerplate-go/internal/models"
)

func (gu *GeneralUsecase) UploadImage(ctx context.Context, param models.UploadFileReq) (resp models.UploadFileResp) {
	// Initialize Usecase Parameter
	resp.InitError()
	ucContext := map[string]interface{}{
		"usecase": "UploadImage",
		"email":   param.ActionBy,
	}

	if !entities.IsValidImageMime(param.FileExt) {
		resp.ErrorGeneral("Tipe gambar tidak valid.")
		return
	}
	if len(param.FileContent) == 0 {
		resp.ErrorGeneral("Gambar tidak valid.")
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
	file, err := gu.upload.UploadFile(ctx, image)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do UploadFile because %s", err)
		return
	}

	resp.Token = file.Token
	resp.FileURL = gu.upload.GenerateUserContentURL(ctx, file.FileURL)
	resp.Success()
	return
}

func (gu *GeneralUsecase) UploadFile(ctx context.Context, param models.UploadFileReq) (resp models.UploadFileResp) {
	// Initialize Usecase Parameter
	resp.InitError()
	ucContext := map[string]interface{}{
		"usecase": "UploadFile",
		"email":   param.ActionBy,
	}

	if !entities.IsValidFileMime(param.FileExt) {
		resp.ErrorGeneral("Tipe gambar tidak valid.")
		return
	}
	if len(param.FileContent) == 0 {
		resp.ErrorGeneral("Gambar tidak valid.")
		return
	}

	doc := entities.UploadFile{
		FileType:        entities.FileTypeDocument,
		FileContent:     param.FileContent,
		FileName:        param.FileName,
		CreatedBy:       param.ActionBy,
		UpdatedBy:       param.ActionBy,
		FileContentType: param.FileExt,
	}
	file, err := gu.upload.UploadFile(ctx, doc)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do UploadFile because %s", err)
		return
	}

	resp.Token = file.Token
	resp.FileURL = gu.upload.GenerateUserContentURL(ctx, file.FileURL)
	resp.Success()
	return
}

func (gu *GeneralUsecase) UploadVideo(ctx context.Context, param models.UploadFileReq) (resp models.UploadFileResp) {
	resp.InitError()
	ucContext := map[string]interface{}{
		"usecase": "UploadVideo",
		"email":   param.ActionBy,
	}

	if !entities.IsValidVideoMime(param.FileExt) {
		resp.ErrorGeneral("Tipe video tidak valid.")
		return
	}
	if len(param.FileContent) == 0 {
		resp.ErrorGeneral("Video tidak valid.")
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
	file, err := gu.upload.UploadFile(ctx, image)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do UploadVideo because %s", err)
		return
	}

	resp.Token = file.Token
	resp.FileURL = gu.upload.GenerateUserContentURL(ctx, file.FileURL)
	resp.Success()
	return
}
