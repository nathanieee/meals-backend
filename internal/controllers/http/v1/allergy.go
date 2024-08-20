package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/services/allergyservice"
	"project-skbackend/packages/utils/utrequest"
	"project-skbackend/packages/utils/utresponse"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	allergyroutes struct {
		cfg   *configs.Config
		salle allergyservice.IAllergyService
	}
)

func newAllergyRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	salle allergyservice.IAllergyService,
) {
	r := &allergyroutes{
		cfg:   cfg,
		salle: salle,
	}

	gallergypub := rg.Group("allergies")
	{
		gallergypub.GET("", r.findAllergies)
		gallergypub.GET("raw", r.findAllergiesRaw)
		gallergypub.GET(":alid", r.getAllergy)
	}
}

func (r *allergyroutes) findAllergies(ctx *gin.Context) {
	var (
		entity  = "allergies"
		reqpage = utrequest.GeneratePaginationFromRequest(ctx)
	)

	allergies, err := r.salle.FindAll(reqpage)
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
		allergies,
	)
}

func (r *allergyroutes) findAllergiesRaw(ctx *gin.Context) {
	var (
		entity = "allergies"
	)

	allergies, err := r.salle.Read()
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
		allergies,
	)
}

func (r *allergyroutes) getAllergy(ctx *gin.Context) {
	var (
		entity = "allergy"
	)

	alid := ctx.Param("alid")
	aliduuid, err := uuid.Parse(alid)
	if err != nil {
		utresponse.GeneralInputRequiredError(
			entity,
			ctx,
			err,
		)
		return
	}

	allergy, err := r.salle.GetByID(aliduuid)
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
		allergy,
	)
}
