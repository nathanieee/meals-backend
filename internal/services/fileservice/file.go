package fileservice

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"project-skbackend/configs"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"

	"github.com/gin-gonic/gin"
)

type (
	ImageService struct {
		cfg         *configs.Config
		imgbaseddir string
		imgprofdir  string
		imgmealdir  string
	}

	IImageService interface {
		Upload(fileheader *multipart.FileHeader, imgtype consttypes.ImageType, ctx *gin.Context) error
	}
)

func NewImageService(cfg *configs.Config) *ImageService {
	return &ImageService{
		cfg:         cfg,
		imgbaseddir: cfg.Image.BaseDir,
		imgprofdir:  cfg.Image.BaseDir + cfg.Image.ProfileDir,
		imgmealdir:  cfg.Image.BaseDir + cfg.Image.MealDir,
	}
}

func (s *ImageService) Upload(fileheader *multipart.FileHeader, imgtype consttypes.ImageType, ctx *gin.Context) error {
	var uppath, filename string

	// * switch the image type based on the image type
	switch imgtype {
	case consttypes.IT_PROFILE:
		uppath = s.imgprofdir
	case consttypes.IT_MEAL:
		uppath = s.imgmealdir
	default:
		uppath = s.imgbaseddir
	}

	// * create directory if it does not exist
	if err := os.MkdirAll(uppath, os.ModePerm); err != nil {
		utlogger.LogError(err)
		return err
	}

	// * set the file name and destination
	filename = fmt.Sprintf("%s_%s", imgtype, filename)
	destination := filepath.Join(uppath, filename)
	if err := ctx.SaveUploadedFile(fileheader, destination); err != nil {
		utlogger.LogError(err)
		return err
	}

	return nil
}
