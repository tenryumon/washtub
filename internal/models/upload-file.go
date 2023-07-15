package models

type UploadFileReq struct {
	ActionBy int64
	OrgID    int64

	FileName    string
	FileContent []byte
	FileExt     string
}

type ParentUploadFileReq struct {
	ActionBy int64
	OrgID    int64

	FileName    string
	FileContent []byte
	FileExt     string
	OrgCode     string
}

type UploadFileResp struct {
	Errors

	Token   string
	FileURL string
}

type DownloadFileResp struct {
	Errors `json:"-"`

	List []FileList `json:"list"`
}

type FileList struct {
	FileName    string `json:"name"`
	FileURL     string `json:"url"`
	CreatedTime string `json:"created_time"`
}
