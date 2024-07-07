package fileservice

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"project-skbackend/configs"
	"project-skbackend/packages/consttypes"

	"github.com/gin-gonic/gin"
)

type (
	ImageService struct {
		cfg *configs.Config
		ibd string
		ipd string
		imd string
	}

	IImageService interface {
		Upload(fileheader *multipart.FileHeader, imgtype consttypes.ImageType, ctx *gin.Context) error
	}
)

func NewImageService(cfg *configs.Config) *ImageService {
	return &ImageService{
		cfg: cfg,
		ibd: cfg.FileImage.BaseDir,
		ipd: cfg.FileImage.BaseDir + cfg.FileImage.ProfileDir,
		imd: cfg.FileImage.BaseDir + cfg.FileImage.MealDir,
	}
}

func (s *ImageService) Upload(fileheader *multipart.FileHeader, imgtype consttypes.ImageType, ctx *gin.Context) error {
	var (
		uppath, filename string
	)

	// * switch the image type based on the image type
	switch imgtype {
	case consttypes.IT_PROFILE:
		uppath = s.ipd
	case consttypes.IT_MEAL:
		uppath = s.imd
	default:
		uppath = s.ibd
	}

	// * create directory if it does not exist
	if err := os.MkdirAll(uppath, os.ModePerm); err != nil {
		return consttypes.ErrFailedToCreateDirectory
	}

	// * set the file name and destination
	filename = fmt.Sprintf("%s_%s", imgtype, filename)
	destination := filepath.Join(uppath, filename)
	if err := ctx.SaveUploadedFile(fileheader, destination); err != nil {
		return consttypes.ErrFailedToUploadFile
	}

	return nil
}
