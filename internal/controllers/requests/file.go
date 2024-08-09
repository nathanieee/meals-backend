package requests

import (
	"mime/multipart"
	"project-skbackend/packages/consttypes"
)

type (
	FileBase64 struct {
		FileBase64Str string              `json:"file_base64" form:"file_base64" binding:"required"`
		FileName      string              `json:"file_name" form:"file_name" binding:"required"`
		FileType      consttypes.FileType `json:"file_type" form:"file_type" binding:"required"`
	}

	FileUpload struct {
		File *multipart.FileHeader `json:"file" form:"file" binding:"required,file"`
	}
)

func NewFileUpload(file *multipart.FileHeader) *FileUpload {
	return &FileUpload{
		File: file,
	}
}
