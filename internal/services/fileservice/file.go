package fileservice

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/repositories/userrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type (
	FileService struct {
		cfg *configs.Config
		ctx context.Context

		minio minio.Client

		ruser userrepo.IUserRepository

		me    string
		mussl bool
		mpubk string
		mpvtk string
		mb    string
		ml    string
	}

	IFileService interface {
		UploadImage(fileheader *multipart.FileHeader, imgtype consttypes.ImageType, ctx *context.Context) (string, error)
	}
)

func NewFileService(
	cfg *configs.Config,
	ctx context.Context,
	minio minio.Client,
	ruser userrepo.IUserRepository,
) *FileService {
	return &FileService{
		cfg: cfg,
		ctx: ctx,

		minio: minio,

		ruser: ruser,

		me:    cfg.Minio.Endpoint,
		mussl: cfg.Minio.UseSSL,
		mpubk: cfg.Minio.PublicKey,
		mpvtk: cfg.Minio.PrivateKey,
		mb:    cfg.Minio.Bucket,
		ml:    cfg.Minio.Location,
	}
}

func (s *FileService) UploadProfilePicture(uid uuid.UUID, fileheader *multipart.FileHeader) error {
	// TODO: implement uploading to the S3 and return the URL

	return nil
}

func (s *FileService) Upload(req requests.FileUpload) (string, error) {
	fileheader := req.File

	// * open the file
	file, err := fileheader.Open()
	if err != nil {
		utlogger.Error(err)
	}
	defer file.Close()

	// * generate object name and content type
	objname := filepath.Base(fileheader.Filename)
	contenttype := fileheader.Header.Get("Content-Type")

	// * upload the file
	info, err := s.minio.PutObject(s.ctx, s.mb, objname, file, fileheader.Size, minio.PutObjectOptions{ContentType: contenttype})
	if err != nil {
		utlogger.Error(err)
	}

	utlogger.Info(fmt.Sprintf("Successfully uploaded %s of size %d\n", objname, info.Size))
	return "", nil
}
