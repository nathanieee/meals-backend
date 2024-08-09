package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/services/fileservice"
	"project-skbackend/packages/utils/utresponse"

	"github.com/gin-gonic/gin"
)

type (
	fileroutes struct {
		cfg   *configs.Config
		sfile fileservice.IFileService
	}
)

func newFileRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	sfile fileservice.IFileService,
) {
	r := &fileroutes{
		cfg:   cfg,
		sfile: sfile,
	}

	gfile := rg.Group("files")
	{
		gfile.POST("upload", r.uploadFile)
	}
}

func (r *fileroutes) uploadFile(ctx *gin.Context) {
	var (
		function = "upload file"
		req      requests.FileMultipart
		err      error
	)

	if err = ctx.ShouldBind(&req); err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			function,
			ctx,
			ve,
			err,
		)
		return
	}

	url, err := r.sfile.Upload(req)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccess(
		function,
		ctx,
		url,
	)
}
