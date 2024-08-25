package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/services/donationservice"
	"project-skbackend/packages/utils/utrequest"
	"project-skbackend/packages/utils/utresponse"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	donationroutes struct {
		cfg       *configs.Config
		sdonation donationservice.IDonationService
	}
)

func newDonationRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	sdonation donationservice.IDonationService,
) {
	r := &donationroutes{
		cfg:       cfg,
		sdonation: sdonation,
	}

	gdonationpub := rg.Group("donations")
	{
		gdonationpub.GET("", r.findDonations)
		gdonationpub.GET("raw", r.findDonationsRaw)
		gdonationpub.GET(":donid", r.getDonation)
	}
}

func (r *donationroutes) findDonations(ctx *gin.Context) {
	var (
		entity  = "donations"
		reqpage = utrequest.GeneratePaginationFromRequest(ctx)
	)

	donations, err := r.sdonation.FindAll(reqpage)
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
		donations,
	)
}

func (r *donationroutes) findDonationsRaw(ctx *gin.Context) {
	var (
		entity = "donations"
	)

	donations, err := r.sdonation.Read()
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
		donations,
	)
}

func (r *donationroutes) getDonation(ctx *gin.Context) {
	var (
		entity = "donation"
	)

	donid := ctx.Param("donid")
	doniduuid, err := uuid.Parse(donid)
	if err != nil {
		utresponse.GeneralInputRequiredError(
			entity,
			ctx,
			err,
		)
		return
	}

	donation, err := r.sdonation.GetByID(doniduuid)
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
		donation,
	)
}
