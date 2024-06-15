package authservice

import (
	"errors"
	"fmt"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/userrepo"
	"project-skbackend/internal/services/mailservice"
	"project-skbackend/internal/services/userservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utstring"
	"project-skbackend/packages/utils/uttoken"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type (
	AuthService struct {
		cfg   *configs.Config
		rdb   *redis.Client
		ruser userrepo.IUserRepository
		suser userservice.IUserService
		smail mailservice.IMailService

		vtl int
		wu  string
	}

	IAuthService interface {
		Signin(req requests.Signin, ctx *gin.Context) (*responses.User, *uttoken.TokenHeader, error)
		ForgotPassword(req requests.ForgotPassword) error
		ResetPassword(req requests.ResetPassword) error
		SendResetPasswordEmail(user models.User) error
		RefreshAuthToken(trefresh string, ctx *gin.Context) (*responses.User, *uttoken.TokenHeader, error)
		SendVerificationEmail(id uuid.UUID) error
		VerifyToken(req requests.VerifyToken, ctx *gin.Context) (*responses.User, *uttoken.TokenHeader, error)
	}
)

func NewAuthService(
	cfg *configs.Config,
	rdb *redis.Client,
	ruser userrepo.IUserRepository,
	smail mailservice.IMailService,
	suser userservice.IUserService,
) *AuthService {
	return &AuthService{
		cfg:   cfg,
		ruser: ruser,
		smail: smail,
		suser: suser,
		rdb:   rdb,

		vtl: cfg.VerifyTokenLength,
		wu:  cfg.Web.URL,
	}
}

func (s *AuthService) Signin(req requests.Signin, ctx *gin.Context) (*responses.User, *uttoken.TokenHeader, error) {
	user, err := s.ruser.GetByEmail(req.Email)
	if err != nil {
		return nil, nil, err
	}

	err = verifyPassword(*user, req.Password)
	if err != nil {
		return nil, nil, err
	}

	thead, err := s.generateAuthTokens(user, ctx)
	if err != nil {
		return nil, nil, err
	}

	userres, err := user.ToResponse()
	if err != nil {
		return nil, nil, err
	}

	return userres, thead, nil
}

func (s *AuthService) ForgotPassword(req requests.ForgotPassword) error {
	user, err := s.ruser.GetByEmail(req.Email)
	if err != nil {
		return err
	}

	err = s.SendResetPasswordEmail(*user)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ResetPassword(req requests.ResetPassword) error {
	user, err := s.ruser.GetByEmail(req.Email)
	if err != nil && err == gorm.ErrRecordNotFound {
		return err
	}

	if req.Token != user.ResetPasswordToken {
		return consttypes.ErrTokenMismatch
	}

	if !consttypes.DateNow.Before(user.ResetPasswordSentAt.Add(time.Minute * time.Duration(s.cfg.ResetPassword.Cooldown))) {
		return consttypes.ErrTooQuickSendEmail
	}

	user, err = req.ToUserModel(*user)
	if err != nil {
		return err
	}

	_, err = s.ruser.Update(*user)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) SendResetPasswordEmail(user models.User) error {
	token, err := utstring.GenerateRandomToken(s.vtl)
	if err != nil {
		return err
	}

	if consttypes.DateNow.Before(user.ResetPasswordSentAt.Add(time.Minute * time.Duration(s.cfg.ResetPassword.Cooldown))) {
		return consttypes.ErrTooQuickSendEmail
	}

	user.ResetPasswordToken = token
	user.ResetPasswordSentAt = consttypes.DateNow

	_, err = s.ruser.Update(user)
	if err != nil {
		return err
	}

	firstname, lastname, err := s.suser.GetUserName(user.ID)
	if err != nil {
		return err
	}

	name := utstring.AppendName(firstname, lastname)
	emreq := requests.SendEmailResetPassword{
		Name:    name,
		Email:   user.Email,
		LinkUrl: fmt.Sprintf("%s/reset-password/%v", s.wu, token),
	}

	err = s.smail.SendResetPasswordEmail(emreq)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RefreshAuthToken(trefresh string, ctx *gin.Context) (*responses.User, *uttoken.TokenHeader, error) {
	now := consttypes.DateNow
	tparsed, err := uttoken.ParseToken(trefresh, s.cfg.JWT.RefreshToken.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.ruser.GetByID(tparsed.User.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, fmt.Errorf("the user with id %s does not exist", tparsed.User.ID)
		}

		return nil, nil, err
	}

	if now.Unix() >= tparsed.Expires.Unix() {
		return nil, nil, err
	}

	userres, err := user.ToResponse()
	if err != nil {
		return nil, nil, err
	}

	taccess, err := uttoken.
		GenerateToken(
			userres,
			s.cfg.JWT.AccessToken.Life,
			s.cfg.JWT.TimeUnit,
			s.cfg.JWT.AccessToken.PrivateKey,
		)

	if err != nil {
		return nil, nil, err
	}

	// * setting the access token into redis
	err = s.rdb.Set(ctx, taccess.TokenUUID.String(), user.ID, time.Unix(taccess.Expires.Unix(), 0).Sub(now)).Err()
	if err != nil {
		return nil, nil, err
	}

	// * setting the cookie for access token
	ctx.SetCookie(
		"access_token",
		*taccess.Token,
		int(taccess.Expires.Unix()),
		"/",
		s.cfg.API.URL,
		false,
		true,
	)

	theader, err := s.generateAuthTokens(user, ctx)
	if err != nil {
		return nil, nil, err
	}

	return userres, theader, err
}

func verifyPassword(user models.User, password string) error {
	ok := utstring.CheckPasswordHash(password, user.Password)
	if !ok {
		return fmt.Errorf("incorrect email or password")
	}

	return nil
}

func (s *AuthService) generateAuthTokens(user *models.User, ctx *gin.Context) (*uttoken.TokenHeader, error) {
	userres, err := user.ToResponse()
	if err != nil {
		return nil, err
	}

	now := consttypes.DateNow
	trefresh, err := uttoken.
		GenerateToken(
			userres,
			s.cfg.JWT.RefreshToken.Life,
			s.cfg.JWT.TimeUnit,
			s.cfg.JWT.RefreshToken.PrivateKey,
		)

	if err != nil {
		return nil, err
	}

	taccess, err := uttoken.
		GenerateToken(
			userres,
			s.cfg.JWT.AccessToken.Life,
			s.cfg.JWT.TimeUnit,
			s.cfg.JWT.AccessToken.PrivateKey,
		)

	if err != nil {
		return nil, err
	}

	theader := uttoken.NewTokenHeader(*taccess, *trefresh)

	// * setting the access token into redis
	err = s.rdb.Set(ctx, taccess.TokenUUID.String(), user.ID, time.Unix(taccess.Expires.Unix(), 0).Sub(now)).Err()
	if err != nil {
		return nil, err
	}

	// * setting the cookie for access token
	ctx.SetCookie(
		"access_token",
		*taccess.Token,
		int(taccess.Expires.Unix()),
		"/",
		s.cfg.API.URL,
		false,
		true,
	)

	// * setting the refresh token into redis
	err = s.rdb.Set(ctx, trefresh.TokenUUID.String(), user.ID, time.Unix(trefresh.Expires.Unix(), 0).Sub(now)).Err()
	if err != nil {
		return nil, err
	}

	// * setting the cookie for refresh token
	ctx.SetCookie(
		"refresh-token",
		*trefresh.Token,
		int(trefresh.Expires.Unix()),
		"/",
		s.cfg.API.URL,
		false,
		true,
	)

	return theader, err
}

func (s *AuthService) SendVerificationEmail(id uuid.UUID) error {
	tverif, err := utstring.GenerateRandomToken(s.vtl)
	if err != nil {
		return err
	}

	user, err := s.ruser.GetByID(id)
	if err != nil {
		return err
	}

	if consttypes.DateNow.Before(user.ConfirmationSentAt.Add(time.Minute * 5)) {
		return consttypes.ErrTooQuickSendEmail
	}

	user.ConfirmationToken = tverif
	user.ConfirmationSentAt = consttypes.DateNow

	user, err = s.ruser.Update(*user)
	if err != nil {
		return err
	}

	userres, err := user.ToResponse()
	if err != nil {
		return err
	}

	firstname, lastname, err := s.suser.GetUserName(user.ID)
	if err != nil {
		return err
	}

	name := utstring.AppendName(firstname, lastname)
	emailData := requests.SendEmailVerification{
		Name:  name,
		Email: userres.Email,
		Token: tverif,
	}

	err = s.smail.SendVerifyEmail(emailData)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) VerifyToken(req requests.VerifyToken, ctx *gin.Context) (*responses.User, *uttoken.TokenHeader, error) {
	user, err := s.ruser.GetByEmail(req.Email)
	if err != nil {
		return nil, nil, err
	}

	if !consttypes.DateNow.Before(user.ConfirmationSentAt.Add(time.Minute * 5)) {
		return nil, nil, errors.New("this token is expired")
	}

	if !user.ConfirmedAt.Equal(time.Time{}) {
		return nil, nil, errors.New("this user is already confirmed")
	}

	if req.Token != user.ConfirmationToken {
		return nil, nil, errors.New("this token is not the same")
	}

	user.ConfirmationSentAt = consttypes.DateNow
	user.ConfirmationToken = ""

	user, err = s.ruser.Update(*user)
	if err != nil {
		return nil, nil, err
	}

	theader, err := s.generateAuthTokens(user, ctx)
	if err != nil {
		return nil, nil, err
	}

	userres, err := user.ToResponse()
	if err != nil {
		return nil, nil, err
	}

	return userres, theader, nil
}
