package fileservice

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"project-skbackend/configs"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/donationproofrepo"
	"project-skbackend/internal/repositories/donationrepo"
	"project-skbackend/internal/repositories/imagerepo"
	"project-skbackend/internal/repositories/userimagerepo"
	"project-skbackend/internal/repositories/userrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utfile"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type (
	FileService struct {
		cfg *configs.Config
		ctx context.Context

		minio minio.Client

		ruser userrepo.IUserRepository
		rimg  imagerepo.IImageRepo
		ruimg userimagerepo.IUserImageRepo
		rdona donationrepo.IDonationRepository
		rdnpr donationproofrepo.IDonationProofRepository

		me    string
		mussl bool
		mpubk string
		mpvtk string
		mb    string
		ml    string
	}

	IFileService interface {
		Upload(req utfile.FileMultipart) (string, error)
		UploadProfilePicture(uid uuid.UUID, fileheader *multipart.FileHeader) error
		UploadDonationProof(did uuid.UUID, fileheader *multipart.FileHeader) error
	}
)

func NewFileService(
	cfg *configs.Config,
	ctx context.Context,
	minio minio.Client,
	ruser userrepo.IUserRepository,
	rimg imagerepo.IImageRepo,
	ruimg userimagerepo.IUserImageRepo,
	rdona donationrepo.IDonationRepository,
	rdnpr donationproofrepo.IDonationProofRepository,
) *FileService {
	return &FileService{
		cfg: cfg,
		ctx: ctx,

		minio: minio,

		ruser: ruser,
		rimg:  rimg,
		ruimg: ruimg,
		rdona: rdona,
		rdnpr: rdnpr,

		me:    cfg.Minio.Endpoint,
		mussl: cfg.Minio.UseSSL,
		mpubk: cfg.Minio.PublicKey,
		mpvtk: cfg.Minio.PrivateKey,
		mb:    cfg.Minio.Bucket,
		ml:    cfg.Minio.Location,
	}
}

func (s *FileService) UploadProfilePicture(uid uuid.UUID, fileheader *multipart.FileHeader) error {
	var (
		err error
	)

	user, err := s.ruser.GetByID(uid)
	if err != nil {
		utlogger.Error(err)
		return consttypes.ErrUserNotFound
	}

	filename, err := utfile.ValidateFile(fileheader, &utfile.ValidateFileOpts{
		AllowedExtensions: []string{".jpg", ".jpeg", ".png", ".gif"},
		ValidateFileSizeOpts: utfile.ValidateFileSizeOpts{
			MaxImageSize:       2,
			MaxImageSizeSuffix: consttypes.FSS_MB,
		},
	})
	if err != nil {
		utlogger.Error(err)
		return err
	}

	fileupload := utfile.NewFileUpload(fileheader)
	url, err := s.Upload(*fileupload)
	if err != nil {
		utlogger.Error(err)
		return consttypes.ErrFailedToUploadFile
	}

	image := models.NewProfileImage(
		*filename,
		url,
	)
	image, err = s.rimg.Create(*image)
	if err != nil {
		utlogger.Error(err)
		return consttypes.ErrFailedToCreateImage
	}

	err = s.CheckAndSaveUserImage(*user, *image)
	if err != nil {
		utlogger.Error(err)
		return consttypes.ErrGeneralFailed("check and save user image", err.Error())
	}

	return nil
}

func (s *FileService) UploadDonationProof(did uuid.UUID, fileheader *multipart.FileHeader) error {
	var (
		err error
	)

	// * get donation by its id
	donation, err := s.rdona.GetByID(did)
	if err != nil {
		utlogger.Error(err)
		return consttypes.ErrUserNotFound
	}

	// * validate the file based on the custom options
	filename, err := utfile.ValidateFile(fileheader, &utfile.ValidateFileOpts{
		AllowedExtensions: []string{".jpg", ".jpeg", ".png"},
		ValidateFileSizeOpts: utfile.ValidateFileSizeOpts{
			MaxImageSize:       1,
			MaxImageSizeSuffix: consttypes.FSS_MB,
		},
	})
	if err != nil {
		utlogger.Error(err)
		return err
	}

	// * construct a new fileupload request
	// * and then upload the file using the upload service
	fileupload := utfile.NewFileUpload(fileheader)
	url, err := s.Upload(*fileupload)
	if err != nil {
		utlogger.Error(err)
		return consttypes.ErrFailedToUploadFile
	}

	// * construct a new donation proof image model
	// * then create it in the database
	image := models.NewDonationProof(
		*filename,
		url,
	)
	image, err = s.rimg.Create(*image)
	if err != nil {
		utlogger.Error(err)
		return consttypes.ErrFailedToCreateImage
	}

	// * construct a new donation proof model
	// * then create it in the database
	donationproof := image.CreateDonationProof(
		*donation,
	)
	if _, err := s.rdnpr.Create(*donationproof); err != nil {
		utlogger.Error(err)
		return consttypes.ErrGeneralFailed("creating donation proof", err.Error())
	}

	return nil
}

func (s *FileService) CheckAndSaveUserImage(u models.User, image models.Image) error {
	ui, err := s.ruimg.GetByUserID(u.ID)
	if ui == nil || (err != nil && errors.Is(err, gorm.ErrRecordNotFound)) {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return consttypes.ErrGeneralFailed("get by user id", err.Error())
		}

		// * handle the case where there is no user image
		updatedui := image.CreateUserImage(
			u,
		)
		if _, err := s.ruimg.Create(*updatedui); err != nil {
			utlogger.Error(err)
			return consttypes.ErrGeneralFailed("creating user image", err.Error())
		}
	} else {
		// * handle the case where there is a user image
		updatedui := image.UpdateUserImage(*ui)
		if _, err := s.ruimg.Update(*updatedui); err != nil {
			return consttypes.ErrGeneralFailed("updating user image", err.Error())
		}
	}

	return nil
}

func (s *FileService) Upload(req utfile.FileMultipart) (string, error) {
	fileheader := req.File

	// * open the file
	file, err := fileheader.Open()
	if err != nil {
		utlogger.Error(err)
		return "", err
	}
	defer file.Close()

	// * generate a new random uuid v7 to replace the filename
	// * reason: so people could not guess the path or pattern of the file
	randuuid, err := uuid.NewV7()
	if err != nil {
		utlogger.Error(err)
		return "", err
	}

	// * extract the extension from the fileheader
	// * and construct the object name
	fileext := utfile.GetFileExtension(fileheader)
	contenttype := fileheader.Header.Get("Content-Type")
	objname := fmt.Sprintf("%s-%s%s", contenttype, randuuid.String(), fileext)

	// * upload the file
	info, err := s.minio.PutObject(s.ctx, s.mb, objname, file, fileheader.Size, minio.PutObjectOptions{ContentType: contenttype})
	if err != nil {
		utlogger.Error(err)
		return "", err
	}

	utlogger.Info(fmt.Sprintf("Successfully uploaded %s of size %d", objname, info.Size))

	miniourl := fmt.Sprintf("%s/%s/%s", s.minio.EndpointURL().String(), s.mb, objname)
	environment := s.cfg.API.Environment

	// * if its on local, then replace the endpoint url to localhost ip.
	if environment == "local" {
		miniourl = fmt.Sprintf("%s/%s/%s", "http://127.0.0.1:9000", s.mb, objname)
	}

	return miniourl, nil
}
