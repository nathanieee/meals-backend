package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/services/illnessservice"
	"project-skbackend/packages/utils/utrequest"
	"project-skbackend/packages/utils/utresponse"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	illnessroutes struct {
		cfg      *configs.Config
		sillness illnessservice.IIllnessService
	}
)

func newIllnessRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	sillness illnessservice.IIllnessService,
) {
	r := &illnessroutes{
		cfg:      cfg,
		sillness: sillness,
	}

	gillnesspub := rg.Group("illnesses")
	{
		gillnesspub.GET("", r.findIllnesses)
		gillnesspub.GET("raw", r.findIllnessesRaw)
		gillnesspub.GET(":illid", r.getIllnesses)
	}
}

func (r *illnessroutes) findIllnesses(ctx *gin.Context) {
	var (
		entity  = "illnesses"
		reqpage = utrequest.GeneratePaginationFromRequest(ctx)
	)

	illnesses, err := r.sillness.FindAll(reqpage)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utresponse.GeneralNotFound(
				entity,
				ctx,
				err,
			)
			return
		}

		utresponse.GeneralInternalServerError(
			entity,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessFetch(
		entity,
		ctx,
		illnesses,
	)
}

func (r *illnessroutes) findIllnessesRaw(ctx *gin.Context) {
	var (
		entity = "illnesses"
	)

	illnesses, err := r.sillness.Read()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utresponse.GeneralNotFound(
				entity,
				ctx,
				err,
			)
			return
		}

		utresponse.GeneralInternalServerError(
			entity,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessFetch(
		entity,
		ctx,
		illnesses,
	)
}

func (r *illnessroutes) getIllnesses(ctx *gin.Context) {
	var (
		entity = "illnesses"
	)

	illid := ctx.Param("illid")
	illiduuid, err := uuid.Parse(illid)
	if err != nil {
		utresponse.GeneralInputRequiredError(
			entity,
			ctx,
			err,
		)
		return
	}

	illness, err := r.sillness.GetByID(illiduuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utresponse.GeneralNotFound(
				entity,
				ctx,
				err,
			)
			return
		}

		utresponse.GeneralInternalServerError(
			entity,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessFetch(
		entity,
		ctx,
		illness,
	)
}
