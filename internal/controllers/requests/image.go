package requests

import (
	"mime/multipart"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utfile"
)

type (
	CreateImage struct {
		*utfile.FileMultipart `binding:"omitempty,dive"`
		*utfile.FileBase64    `binding:"omitempty,dive"`
	}

	UpdateImage struct {
		*utfile.FileMultipart `binding:"omitempty,dive"`
		*utfile.FileBase64    `binding:"omitempty,dive"`
	}
)

func (r *CreateImage) Validate() error {
	// * validate if there no file provided
	if r.FileMultipart == nil && r.FileBase64 == nil {
		return consttypes.ErrNoFiles
	}

	// * validate if there is more than one (1) file
	if r.FileMultipart != nil && r.FileBase64 != nil {
		return consttypes.ErrTooManyFiles
	}

	return nil
}

func (r *CreateImage) GetMultipartFile() (*multipart.FileHeader, error) {
	if r.FileBase64 != nil {
		multipart, err := utfile.Base64ToMultipartFileHeader(r.FileBase64.FileBase64Str, r.FileBase64.FileName, r.FileBase64.FileType)
		if err != nil {
			return nil, err
		}
		return multipart, nil
	}

	if r.FileMultipart != nil {
		return r.FileMultipart.File, nil
	}

	return nil, consttypes.ErrNoFiles
}

func (r *CreateImage) GetBase64File() (*utfile.FileBase64, error) {
	if r.FileBase64 != nil {
		return r.FileBase64, nil
	}

	if r.FileMultipart != nil {
		filebase64, err := utfile.MultipartFileHeaderToBase64(r.FileMultipart.File)
		if err != nil {
			return nil, err
		}
		return filebase64, nil
	}

	return nil, consttypes.ErrNoFiles
}

func (r *UpdateImage) Validate() error {
	// * validate if there no file provided
	if r.FileMultipart == nil && r.FileBase64 == nil {
		return consttypes.ErrNoFiles
	}

	// * validate if there is more than one (1) file
	if r.FileMultipart != nil && r.FileBase64 != nil {
		return consttypes.ErrTooManyFiles
	}

	return nil
}

func (r *UpdateImage) GetMultipartFile() (*multipart.FileHeader, error) {
	if r.FileBase64 != nil {
		multipart, err := utfile.Base64ToMultipartFileHeader(r.FileBase64.FileBase64Str, r.FileBase64.FileName, r.FileBase64.FileType)
		if err != nil {
			return nil, err
		}
		return multipart, nil
	}

	if r.FileMultipart != nil {
		return r.FileMultipart.File, nil
	}

	return nil, consttypes.ErrNoFiles
}

func (r *UpdateImage) GetBase64File() (*utfile.FileBase64, error) {
	if r.FileBase64 != nil {
		return r.FileBase64, nil
	}

	if r.FileMultipart != nil {
		filebase64, err := utfile.MultipartFileHeaderToBase64(r.FileMultipart.File)
		if err != nil {
			return nil, err
		}
		return filebase64, nil
	}

	return nil, consttypes.ErrNoFiles
}
