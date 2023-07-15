package entities

import (
	"time"
)

const (
	FileRemoved = -1
	FilePending = 0
	FileActive  = 1

	FileTypeImage    = 1
	FileTypeDocument = 2
	FileTypeVideo    = 3
	FileTypeSheet    = 4

	MaxImageSize = 2 * 1024 * 1024
)

var (
	validImageMime = []string{"image/jpeg", "image/png"}
	validFileMime  = []string{"application/pdf"}
	validVideoMime = []string{"video/x-flv", "video/mp4"}

	fileTypeStringList = map[int]string{
		FileTypeImage:    "images",
		FileTypeDocument: "docs",
		FileTypeVideo:    "videos",
		FileTypeSheet:    "sheets",
	}
)

type UploadFile struct {
	Token       string    `db:"token"`
	FileURL     string    `db:"file_url"`
	FileType    int       `db:"file_type"`
	Status      int       `db:"status"`
	CreatedBy   int64     `db:"created_by"`
	CreatedTime time.Time `db:"created_time"`
	UpdatedBy   int64     `db:"updated_by"`
	UpdatedTime time.Time `db:"updated_time"`

	BucketName      string `db:"-"`
	OrgCode         string `db:"-"`
	FileName        string `db:"-"`
	FileContent     []byte `db:"-"`
	FileContentType string `db:"-"`
}

func (f UploadFile) Exist() bool {
	return f.Token != ""
}

func (f UploadFile) IsActive() bool {
	return f.Status == FileActive
}

func (f UploadFile) IsPending() bool {
	return f.Status == FilePending
}

func (f UploadFile) IsRemoved() bool {
	return f.Status == FileRemoved
}

func (f UploadFile) IsImage() bool {
	return f.FileType == FileTypeImage
}

func (f UploadFile) IsDocument() bool {
	return f.FileType == FileTypeDocument
}

func IsValidImageMime(ext string) bool {
	for _, mime := range validImageMime {
		if ext == mime {
			return true
		}
	}
	return false
}

func IsValidFileMime(ext string) bool {
	for _, mime := range validFileMime {
		if ext == mime {
			return true
		}
	}
	return false
}

func IsValidVideoMime(ext string) bool {
	for _, mime := range validVideoMime {
		if ext == mime {
			return true
		}
	}
	return false
}

func (uf UploadFile) GetFileTypeString() string {
	if value, ok := fileTypeStringList[uf.FileType]; ok {
		return value
	}
	return "unknown"
}
