package ekycmodel

import (
	"mime/multipart"
)

type UploadRequest struct {
	Image         *multipart.FileHeader `form:"image"`
	Title         string                `form:"title"`
	Description   string                `form:"description"`
	ClientSession string                `form:"client_session"`
}

type KYCImageUploadRes struct {
	Message string            `json:"message"`
	Object  ImageUploadObject `json:"object"`
}

type ImageUploadObject struct {
	FileName     string `json:"fileName"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Hash         string `json:"hash"`
	FileType     string `json:"fileType"`
	UploadedDate string `json:"uploadedDate"`
	StorageType  string `json:"storageType"`
	TokenId      string `json:"tokenId"`
}
