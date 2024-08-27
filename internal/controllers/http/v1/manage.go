package controllers

import (
	"errors"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/middlewares"
	"project-skbackend/internal/services/allergyservice"
	"project-skbackend/internal/services/donationservice"
	"project-skbackend/internal/services/fileservice"
	"project-skbackend/internal/services/illnessservice"
	"project-skbackend/internal/services/mealservice"
	"project-skbackend/internal/services/memberservice"
	"project-skbackend/internal/services/partnerservice"
	"project-skbackend/internal/services/patronservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utrequest"
	"project-skbackend/packages/utils/utresponse"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type (
	manageroutes struct {
		cfg       *configs.Config
		smeal     mealservice.IMealService
		smember   memberservice.IMemberService
		spartner  partnerservice.IPartnerService
		spatron   patronservice.IPatronService
		sillness  illnessservice.IIllnessService
		sfile     fileservice.IFileService
		sallergy  allergyservice.IAllergyService
		sdonation donationservice.IDonationService
	}
)

func newManageRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	smeal mealservice.IMealService,
	smember memberservice.IMemberService,
	spartner partnerservice.IPartnerService,
	spatron patronservice.IPatronService,
	sillness illnessservice.IIllnessService,
	sfile fileservice.IFileService,
	sallergy allergyservice.IAllergyService,
	sdonation donationservice.IDonationService,
) {
	r := &manageroutes{
		cfg:       cfg,
		smeal:     smeal,
		smember:   smember,
		spartner:  spartner,
		spatron:   spatron,
		sillness:  sillness,
		sfile:     sfile,
		sallergy:  sallergy,
		sdonation: sdonation,
	}

	gmanage := rg.Group("manages")
	gmanage.Use(middlewares.JWTAuthMiddleware(
		cfg,
		consttypes.UR_ADMIN,
	))
	{
		gmeals := gmanage.Group("meals")
		{
			gmeals.POST("", r.createMeal)
			gmeals.GET("", r.findMeals)
			gmeals.GET("raw", r.findMealsRaw)
			gmeals.PUT("/:mid", r.updateMeal)
			gmeals.DELETE("/:mid", r.deleteMeal)
		}

		gmember := gmanage.Group("members")
		{
			gmember.POST("", r.createMember)
			gmember.GET("", r.findMembers)
			gmember.GET("raw", r.findMembersRaw)
			gmember.PUT("/:mid", r.updateMember)
			gmember.DELETE("/:mid", r.deleteMember)
		}

		gpartner := gmanage.Group("partners")
		{
			gpartner.POST("", r.createPartner)
			gpartner.GET("", r.findPartners)
			gpartner.GET("raw", r.findPartnersRaw)
			gpartner.PUT("/:pid", r.updatePartner)
			gpartner.DELETE("/:pid", r.deletePartner)
		}

		gpatron := gmanage.Group("patrons")
		{
			gpatron.POST("", r.createPatron)
			gpatron.GET("", r.findPatrons)
			gpatron.GET("raw", r.findPatronsRaw)
			gpatron.PUT("/:pid", r.updatePatron)
			gpatron.DELETE("/:pid", r.deletePatron)
		}

		gillness := gmanage.Group("illnesses")
		{
			gillness.POST("", r.createIllness)
			gillness.GET("", r.findIllnesses)
			gillness.GET("raw", r.findIllnessesRaw)
			gillness.PUT("/:iid", r.updateIllness)
			gillness.DELETE("/:iid", r.deleteIllness)
		}

		gallergy := gmanage.Group("allergies")
		{
			gallergy.POST("", r.createAllergy)
			gallergy.GET("", r.findAllergies)
			gallergy.GET("raw", r.findAllergiesRaw)
			gallergy.PUT("/:aid", r.updateAllergy)
			gallergy.DELETE("/:aid", r.deleteAllergy)
		}

		gdonation := gmanage.Group("donations")
		{
			// ! no create donation because admin cannot interfene
			gdonation.GET("", r.findDonations)
			gdonation.GET("raw", r.findDonationsRaw)
			gdonation.PUT("/:did", r.updateDonation)
			gdonation.DELETE("/:did", r.deleteDonation)
		}
	}
}

// ! -------------------------------------------------------------------------- ! //
// !                        start of meals routing group                        ! //
// ! -------------------------------------------------------------------------- ! //
func (r *manageroutes) createMeal(ctx *gin.Context) {
	var (
		function = "create meal"
		entity   = "meal"
		req      requests.CreateMeal
	)

	err := ctx.ShouldBind(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			function,
			ctx,
			ve,
			err,
		)
		return
	}

	resmeal, err := r.smeal.Create(req)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessCreate(
		entity,
		ctx,
		resmeal,
	)
}

func (r *manageroutes) findMeals(ctx *gin.Context) {
	var (
		entity  = "meals"
		reqpage = utrequest.GeneratePaginationFromRequest(ctx)
	)

	meals, err := r.smeal.FindAll(reqpage)
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
		meals,
	)
}

func (r *manageroutes) findMealsRaw(ctx *gin.Context) {
	var (
		entity = "meals"
	)

	meals, err := r.smeal.Read()
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
		meals,
	)
}

func (r *manageroutes) updateMeal(ctx *gin.Context) {
	var (
		function = "update meal"
		entity   = "meal"
		req      requests.UpdateMeal
	)

	err := ctx.ShouldBind(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			function,
			ctx,
			ve,
			err,
		)
		return
	}

	uuid, err := uuid.Parse(ctx.Param("mid"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

	_, err = r.smeal.GetByID(uuid)
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

	resmeal, err := r.smeal.Update(uuid, req)
	if err != nil {
		utresponse.GeneralInternalServerError(
			entity,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessUpdate(
		entity,
		ctx,
		resmeal,
	)
}

func (r *manageroutes) deleteMeal(ctx *gin.Context) {
	var (
		function = "delete meal"
		entity   = "meal"
	)

	uuid, err := uuid.Parse(ctx.Param("mid"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

	_, err = r.smeal.GetByID(uuid)
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

	err = r.smeal.Delete(uuid)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessDelete(
		entity,
		ctx,
		nil,
	)
}

// ! -------------------------------------------------------------------------- ! //
// !                         end of meals routing group                         ! //
// ! -------------------------------------------------------------------------- ! //

// ! -------------------------------------------------------------------------- ! //
// !                       start of members routing group                       ! //
// ! -------------------------------------------------------------------------- ! //

func (r *manageroutes) createMember(ctx *gin.Context) {
	var (
		function = "create member"
		entity   = "member"
		req      requests.CreateMember
	)

	err := ctx.ShouldBind(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			function,
			ctx,
			ve,
			err,
		)
		return
	}

	resmemb, err := r.smember.Create(req)
	if err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			if pgerrcode.IsIntegrityConstraintViolation(pgerr.SQLState()) {
				utresponse.GeneralDuplicate(
					pgerr.TableName,
					ctx,
					pgerr,
				)
				return
			}
		} else {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
		}
		return
	}

	// * define the image request
	reqimg := req.User.CreateImage
	// * if the image request is not empty
	// * validate and upload the image
	if reqimg != nil {
		if err := reqimg.Validate(); err != nil {
			utresponse.GeneralInvalidRequest(
				function,
				ctx,
				nil,
				err,
			)
			return
		}

		multipart, err := reqimg.GetMultipartFile()
		if err != nil {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
			return
		}

		err = r.sfile.UploadProfilePicture(resmemb.User.ID, multipart)
		if err != nil {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
			return
		}
	}

	utresponse.GeneralSuccessCreate(
		entity,
		ctx,
		resmemb,
	)
}

func (r *manageroutes) findMembers(ctx *gin.Context) {
	var (
		entity  = "members"
		reqpage = utrequest.GeneratePaginationFromRequest(ctx)
	)

	members, err := r.smember.FindAll(reqpage)
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
		members,
	)
}

func (r *manageroutes) findMembersRaw(ctx *gin.Context) {
	var (
		entity = "members"
	)

	resmemb, err := r.smember.Read()
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
		resmemb,
	)
}

func (r *manageroutes) updateMember(ctx *gin.Context) {
	var (
		function = "update member"
		entity   = "member"
		req      requests.UpdateMember
	)

	err := ctx.ShouldBind(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			function,
			ctx,
			ve,
			err,
		)
		return
	}

	uuid, err := uuid.Parse(ctx.Param("mid"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

	_, err = r.smember.GetByID(uuid)
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

	resmemb, err := r.smember.Update(uuid, req)
	if err != nil {
		utresponse.GeneralFailedUpdate(
			entity,
			ctx,
			err,
		)
		return
	}

	// * define the image request
	reqimg := req.User.UpdateImage
	// * if the image request is not empty
	// * validate and upload the image
	if reqimg != nil {
		if err := reqimg.Validate(); err != nil {
			utresponse.GeneralInvalidRequest(
				function,
				ctx,
				nil,
				err,
			)
			return
		}

		multipart, err := reqimg.GetMultipartFile()
		if err != nil {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
			return
		}

		err = r.sfile.UploadProfilePicture(resmemb.User.ID, multipart)
		if err != nil {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
			return
		}
	}

	utresponse.GeneralSuccessUpdate(
		entity,
		ctx,
		resmemb,
	)
}

func (r *manageroutes) deleteMember(ctx *gin.Context) {
	var (
		function = "delete member"
		entity   = "member"
	)

	uuid, err := uuid.Parse(ctx.Param("mid"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

	_, err = r.smember.GetByID(uuid)
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

	err = r.smember.Delete(uuid)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessDelete(
		entity,
		ctx,
		nil,
	)
}

// ! -------------------------------------------------------------------------- ! //
// !                        end of members routing group                        ! //
// ! -------------------------------------------------------------------------- ! //

// ! -------------------------------------------------------------------------- ! //
// !                       start of partners routing group                      ! //
// ! -------------------------------------------------------------------------- ! //
func (r *manageroutes) createPartner(ctx *gin.Context) {
	var (
		function = "create partner"
		entity   = "partner"
		req      requests.CreatePartner
	)

	err := ctx.ShouldBind(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			function,
			ctx,
			ve,
			err,
		)
		return
	}

	respartner, err := r.spartner.Create(req)
	if err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			if pgerrcode.IsIntegrityConstraintViolation(pgerr.SQLState()) {
				utresponse.GeneralDuplicate(
					pgerr.TableName,
					ctx,
					pgerr,
				)
				return
			}
		} else {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
		}
		return
	}

	// * define the image request
	reqimg := req.User.CreateImage
	// * if the image request is not empty
	// * validate and upload the image
	if reqimg != nil {
		if err := reqimg.Validate(); err != nil {
			utresponse.GeneralInvalidRequest(
				function,
				ctx,
				nil,
				err,
			)
			return
		}

		multipart, err := reqimg.GetMultipartFile()
		if err != nil {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
			return
		}

		err = r.sfile.UploadProfilePicture(respartner.User.ID, multipart)
		if err != nil {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
			return
		}
	}

	utresponse.GeneralSuccessCreate(
		entity,
		ctx,
		respartner,
	)
}

func (r *manageroutes) findPartners(ctx *gin.Context) {
	var (
		entity  = "partners"
		reqpage = utrequest.GeneratePaginationFromRequest(ctx)
	)

	partners, err := r.spartner.FindAll(reqpage)
	if err != nil {
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
	}

	utresponse.GeneralSuccessFetch(
		entity,
		ctx,
		partners,
	)
}

func (r *manageroutes) findPartnersRaw(ctx *gin.Context) {
	var (
		entity = "partners"
	)

	partners, err := r.spartner.Read()
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
		partners,
	)
}

func (r *manageroutes) updatePartner(ctx *gin.Context) {
	var (
		function = "update partner"
		entity   = "partner"
		req      requests.UpdatePartner
	)

	err := ctx.ShouldBind(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			function,
			ctx,
			ve,
			err,
		)
		return
	}

	uuid, err := uuid.Parse(ctx.Param("pid"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

	_, err = r.spartner.GetByID(uuid)
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

	respartner, err := r.spartner.Update(uuid, req)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	// * define the image request
	reqimg := req.User.UpdateImage
	// * if the image request is not empty
	// * validate and upload the image
	if reqimg != nil {
		if err := reqimg.Validate(); err != nil {
			utresponse.GeneralInvalidRequest(
				function,
				ctx,
				nil,
				err,
			)
			return
		}

		multipart, err := reqimg.GetMultipartFile()
		if err != nil {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
			return
		}

		err = r.sfile.UploadProfilePicture(respartner.User.ID, multipart)
		if err != nil {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
			return
		}
	}

	utresponse.GeneralSuccessUpdate(
		entity,
		ctx,
		respartner,
	)
}

func (r *manageroutes) deletePartner(ctx *gin.Context) {
	var (
		function = "delete partner"
		entity   = "partner"
	)

	uuid, err := uuid.Parse(ctx.Param("pid"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

	_, err = r.spartner.GetByID(uuid)
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

	err = r.spartner.Delete(uuid)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessDelete(
		entity,
		ctx,
		nil,
	)
}

// ! -------------------------------------------------------------------------- ! //
// !                        end of partners routing group                       ! //
// ! -------------------------------------------------------------------------- ! //

// ! -------------------------------------------------------------------------- ! //
// !                       start of patrons routing group                       ! //
// ! -------------------------------------------------------------------------- ! //

func (r *manageroutes) createPatron(ctx *gin.Context) {
	var (
		function = "create patron"
		entity   = "patron"
		req      requests.CreatePatron
	)

	err := ctx.ShouldBind(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			function,
			ctx,
			ve,
			err,
		)
		return
	}

	respatron, err := r.spatron.Create(req)
	if err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			if pgerrcode.IsIntegrityConstraintViolation(pgerr.SQLState()) {
				utresponse.GeneralDuplicate(
					pgerr.TableName,
					ctx,
					pgerr,
				)
				return
			}
		} else {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
		}
		return
	}

	// * define the image request
	reqimg := req.User.CreateImage
	// * if the image request is not empty
	// * validate and upload the image
	if reqimg != nil {
		if err := reqimg.Validate(); err != nil {
			utresponse.GeneralInvalidRequest(
				function,
				ctx,
				nil,
				err,
			)
			return
		}

		multipart, err := reqimg.GetMultipartFile()
		if err != nil {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
			return
		}

		err = r.sfile.UploadProfilePicture(respatron.User.ID, multipart)
		if err != nil {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
			return
		}
	}

	utresponse.GeneralSuccessCreate(
		entity,
		ctx,
		respatron,
	)
}

func (r *manageroutes) findPatrons(ctx *gin.Context) {
	var (
		entity  = "patrons"
		reqpage = utrequest.GeneratePaginationFromRequest(ctx)
	)

	patrons, err := r.spatron.FindAll(reqpage)
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
		patrons,
	)
}

func (r *manageroutes) findPatronsRaw(ctx *gin.Context) {
	var (
		entity = "patrons"
	)

	patrons, err := r.spatron.Read()
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
		patrons,
	)
}

func (r *manageroutes) updatePatron(ctx *gin.Context) {
	var (
		function = "update patron"
		entity   = "patron"
		req      requests.UpdatePatron
	)

	err := ctx.ShouldBind(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			function,
			ctx,
			ve,
			err,
		)
		return
	}

	uuid, err := uuid.Parse(ctx.Param("pid"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

	_, err = r.spatron.GetByID(uuid)
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

	respatron, err := r.spatron.Update(uuid, req)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	// * define the image request
	reqimg := req.User.UpdateImage
	// * if the image request is not empty
	// * validate and upload the image
	if reqimg != nil {
		if err := reqimg.Validate(); err != nil {
			utresponse.GeneralInvalidRequest(
				function,
				ctx,
				nil,
				err,
			)
			return
		}

		multipart, err := reqimg.GetMultipartFile()
		if err != nil {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
			return
		}

		err = r.sfile.UploadProfilePicture(respatron.User.ID, multipart)
		if err != nil {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
			return
		}
	}

	utresponse.GeneralSuccessUpdate(
		entity,
		ctx,
		respatron,
	)
}

func (r *manageroutes) deletePatron(ctx *gin.Context) {
	var (
		function = "delete patron"
		entity   = "patron"
	)

	uuid, err := uuid.Parse(ctx.Param("pid"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

	_, err = r.spatron.GetByID(uuid)
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

	err = r.spatron.Delete(uuid)
	if err != nil {
		utresponse.GeneralInternalServerError(
			entity,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessDelete(
		entity,
		ctx,
		nil,
	)
}

// ! -------------------------------------------------------------------------- ! //
// !                        end of patrons routing group                        ! //
// ! -------------------------------------------------------------------------- ! //

// ! -------------------------------------------------------------------------- ! //
// !                       start of illness routing group                       ! //
// ! -------------------------------------------------------------------------- ! //

func (r *manageroutes) createIllness(ctx *gin.Context) {
	var (
		function = "create illness"
		entity   = "illness"
		req      requests.CreateIllness
	)

	err := ctx.ShouldBind(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			function,
			ctx,
			ve,
			err,
		)
		return
	}

	resillness, err := r.sillness.Create(req)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessCreate(
		entity,
		ctx,
		resillness,
	)
}

func (r *manageroutes) findIllnesses(ctx *gin.Context) {
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

func (r *manageroutes) findIllnessesRaw(ctx *gin.Context) {
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

func (r *manageroutes) updateIllness(ctx *gin.Context) {
	var (
		function = "update illness"
		entity   = "illness"
		req      requests.UpdateIllness
	)

	err := ctx.ShouldBind(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			function,
			ctx,
			ve,
			err,
		)
		return
	}

	uuid, err := uuid.Parse(ctx.Param("iid"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

	_, err = r.sillness.GetByID(uuid)
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

	resillness, err := r.sillness.Update(uuid, req)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessUpdate(
		entity,
		ctx,
		resillness,
	)
}

func (r *manageroutes) deleteIllness(ctx *gin.Context) {
	var (
		function = "delete illness"
		entity   = "illness"
	)

	uuid, err := uuid.Parse(ctx.Param("iid"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

	_, err = r.sillness.GetByID(uuid)
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

	err = r.sillness.Delete(uuid)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessDelete(
		entity,
		ctx,
		nil,
	)
}

// ! -------------------------------------------------------------------------- ! //
// !                        end of illness routing group                        ! //
// ! -------------------------------------------------------------------------- ! //

// ! -------------------------------------------------------------------------- ! //
// !                       start of allergy routing group                       ! //
// ! -------------------------------------------------------------------------- ! //

func (r *manageroutes) createAllergy(ctx *gin.Context) {
	var (
		function = "create allergy"
		entity   = "allergy"
		req      requests.CreateAllergy
	)

	err := ctx.ShouldBind(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			function,
			ctx,
			ve,
			err,
		)
		return
	}

	resallergy, err := r.sallergy.Create(req)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessCreate(
		entity,
		ctx,
		resallergy,
	)
}

func (r *manageroutes) findAllergies(ctx *gin.Context) {
	var (
		entity  = "allergies"
		reqpage = utrequest.GeneratePaginationFromRequest(ctx)
	)

	allergies, err := r.sallergy.FindAll(reqpage)
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

func (r *manageroutes) findAllergiesRaw(ctx *gin.Context) {
	var (
		entity = "allergies"
	)

	allergies, err := r.sallergy.Read()
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

func (r *manageroutes) updateAllergy(ctx *gin.Context) {
	var (
		function = "update allergy"
		entity   = "allergy"
		req      requests.UpdateAllergy
	)

	err := ctx.ShouldBind(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			function,
			ctx,
			ve,
			err,
		)
		return
	}

	uuid, err := uuid.Parse(ctx.Param("aid"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

	_, err = r.sallergy.GetByID(uuid)
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

	resallergy, err := r.sallergy.Update(req, uuid)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessUpdate(
		entity,
		ctx,
		resallergy,
	)
}

func (r *manageroutes) deleteAllergy(ctx *gin.Context) {
	var (
		function = "delete allergy"
		entity   = "allergy"
	)

	uuid, err := uuid.Parse(ctx.Param("aid"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

	_, err = r.sallergy.GetByID(uuid)
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

	err = r.sallergy.Delete(uuid)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessDelete(
		entity,
		ctx,
		nil,
	)
}

// ! -------------------------------------------------------------------------- ! //
// !                        end of allergy routing group                        ! //
// ! -------------------------------------------------------------------------- ! //

// ! -------------------------------------------------------------------------- ! //
// !                       start of donation routing group                      ! //
// ! -------------------------------------------------------------------------- ! //

func (r *manageroutes) findDonations(ctx *gin.Context) {
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

func (r *manageroutes) findDonationsRaw(ctx *gin.Context) {
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

func (r *manageroutes) updateDonation(ctx *gin.Context) {
	var (
		function = "update donation"
		entity   = "donation"
		req      requests.UpdateDonation
	)

	err := ctx.ShouldBind(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			function,
			ctx,
			ve,
			err,
		)
		return
	}

	uuid, err := uuid.Parse(ctx.Param("did"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

	_, err = r.sdonation.GetByID(uuid)
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

	resdonation, err := r.sdonation.Update(req, uuid)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessUpdate(
		entity,
		ctx,
		resdonation,
	)
}

func (r *manageroutes) deleteDonation(ctx *gin.Context) {
	var (
		function = "delete donation"
		entity   = "donation"
	)

	uuid, err := uuid.Parse(ctx.Param("did"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

	_, err = r.sdonation.GetByID(uuid)
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

	err = r.sdonation.Delete(uuid)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessDelete(
		entity,
		ctx,
		nil,
	)
}
