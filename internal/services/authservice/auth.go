package authservice

import (
	"fmt"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/userrepo"
	"project-skbackend/internal/services/mailservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utresponse"
	"project-skbackend/packages/utils/utstring"
	"project-skbackend/packages/utils/uttoken"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	AuthService struct {
		cfg               *configs.Config
		userrepo          userrepo.IUserRepository
		mailsvc           mailservice.IMailService
		rdb               *redis.Client
		verifyTokenLength int
	}

	IAuthService interface {
		Login(req requests.Login, ctx *gin.Context) (*responses.User, *uttoken.TokenHeader, error)
		Register(req requests.Register, ctx *gin.Context) (*responses.User, *uttoken.TokenHeader, error)
		ForgotPassword(req requests.ForgotPassword) error
		ResetPassword(req requests.ResetPassword) error
		ResetPasswordEmail(id uuid.UUID, token string) error
		RefreshAuthToken(token string, ctx *gin.Context) (*responses.User, *uttoken.TokenHeader, error)
	}
)

func NewAuthService(
	cfg *configs.Config,
	userrepo userrepo.IUserRepository,
	mailsvc mailservice.IMailService,
	rdb *redis.Client,
) *AuthService {
	return &AuthService{
		cfg:               cfg,
		userrepo:          userrepo,
		mailsvc:           mailsvc,
		rdb:               rdb,
		verifyTokenLength: cfg.VerifyTokenLength,
	}
}

func (s *AuthService) Login(req requests.Login, ctx *gin.Context) (*responses.User, *uttoken.TokenHeader, error) {
	user, err := s.userrepo.FindByEmail(req.Email)
	if err != nil {
		return nil, nil, err
	}

	err = verifyPassword(*user, req.Password)
	if err != nil {
		return nil, nil, err
	}

	tokenHeader, err := s.generateAuthTokens(user, ctx)
	if err != nil {
		return nil, nil, err
	}

	return user.ToResponse(), tokenHeader, nil
}

func (s *AuthService) Register(req requests.Register, ctx *gin.Context) (*responses.User, *uttoken.TokenHeader, error) {
	req.Email = strings.ToLower(req.Email)

	user, err := s.userrepo.FindByEmail(req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, nil, err
	}

	if user != nil {
		return user.ToResponse(), nil, err
	}

	user = req.ToUserModel()

	user, err = s.userrepo.Create(*user)
	if err != nil {
		return nil, nil, err
	}

	token, err := s.generateAuthTokens(user, ctx)
	if err != nil {
		return nil, nil, err
	}

	return user.ToResponse(), token, nil
}

func (s *AuthService) ForgotPassword(req requests.ForgotPassword) error {
	user, err := s.userrepo.FindByEmail(req.Email)
	if err != nil && err == gorm.ErrRecordNotFound {
		return err
	}

	token, err := utstring.GenerateRandomToken(s.verifyTokenLength)
	if err != nil {
		return err
	}

	err = s.ResetPasswordEmail(user.ID, token)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ResetPassword(req requests.ResetPassword) error {
	user, err := s.userrepo.FindByEmail(req.Email)
	if err != nil && err == gorm.ErrRecordNotFound {
		return err
	}

	if req.Token != user.ResetPasswordToken {
		return utresponse.ErrTokenMismatch
	}

	if !consttypes.DateNow.Before(user.ResetPasswordSentAt.Add(time.Minute * time.Duration(s.cfg.ResetPassword.Cooldown))) {
		return utresponse.ErrSendEmailResetRequest
	}

	user.Password = req.Password
	user.ResetPasswordToken = ""

	_, err = s.userrepo.Update(*user)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ResetPasswordEmail(id uuid.UUID, token string) error {
	user, err := s.userrepo.FindByID(id)
	if err != nil {
		return err
	}

	if consttypes.DateNow.Before(user.ResetPasswordSentAt.Add(time.Minute * time.Duration(s.cfg.ResetPassword.Cooldown))) {
		return utresponse.ErrSendEmailResetRequest
	}

	user.ResetPasswordToken = token
	user.ResetPasswordSentAt = consttypes.DateNow

	_, err = s.userrepo.Update(*user)
	if err != nil {
		return err
	}

	emreq := requests.SendMail{
		To:       []string{user.Email},
		Subject:  "Reset Password",
		Template: "reset-password.html",
	}

	// TODO - change this request to reset password and sent correct data.
	emreqdata := requests.ResetPasswordEmail{
		Token: token,
	}

	err = s.mailsvc.SendResetPassword(emreq, emreqdata)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RefreshAuthToken(refreshToken string, ctx *gin.Context) (*responses.User, *uttoken.TokenHeader, error) {
	now := consttypes.DateNow
	parsedToken, err := uttoken.ParseToken(refreshToken, s.cfg.JWT.RefreshToken.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.userrepo.FindByID(parsedToken.User.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, fmt.Errorf("the user with id %s does not exist", parsedToken.User.ID)
		}

		return nil, nil, err
	}

	if now.Unix() >= parsedToken.Expires.Unix() {
		return nil, nil, err
	}

	accessToken, err := uttoken.
		GenerateToken(
			user.ToResponse(),
			s.cfg.JWT.AccessToken.Life,
			s.cfg.JWT.TimeUnit,
			s.cfg.JWT.AccessToken.PrivateKey,
		)

	if err != nil {
		return nil, nil, err
	}

	// * setting the access token into redis
	err = s.rdb.Set(ctx, accessToken.TokenUUID.String(), user.ID, time.Unix(accessToken.Expires.Unix(), 0).Sub(now)).Err()
	if err != nil {
		return nil, nil, err
	}

	// * setting the cookie for access token
	ctx.SetCookie(
		"access_token",
		*accessToken.Token,
		int(accessToken.Expires.Unix()),
		"/",
		s.cfg.API.Domain,
		false,
		true,
	)

	tokenHeader, err := s.generateAuthTokens(user, ctx)
	if err != nil {
		return nil, nil, err
	}

	return user.ToResponse(), tokenHeader, err
}

func verifyPassword(user models.User, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) generateAuthTokens(user *models.User, ctx *gin.Context) (*uttoken.TokenHeader, error) {
	now := consttypes.DateNow
	refreshToken, err := uttoken.
		GenerateToken(
			user.ToResponse(),
			s.cfg.JWT.RefreshToken.Life,
			s.cfg.JWT.TimeUnit,
			s.cfg.JWT.RefreshToken.PrivateKey,
		)

	if err != nil {
		return nil, err
	}

	accessToken, err := uttoken.
		GenerateToken(
			user.ToResponse(),
			s.cfg.JWT.AccessToken.Life,
			s.cfg.JWT.TimeUnit,
			s.cfg.JWT.AccessToken.PrivateKey,
		)

	if err != nil {
		return nil, err
	}

	tokenHeader := uttoken.NewTokenHeader(*accessToken, *refreshToken)

	// * setting the access token into redis
	err = s.rdb.Set(ctx, accessToken.TokenUUID.String(), user.ID, time.Unix(accessToken.Expires.Unix(), 0).Sub(now)).Err()
	if err != nil {
		return nil, err
	}

	// * setting the cookie for access token
	ctx.SetCookie(
		"access_token",
		*accessToken.Token,
		int(accessToken.Expires.Unix()),
		"/",
		s.cfg.API.Domain,
		false,
		true,
	)

	// * setting the refresh token into redis
	err = s.rdb.Set(ctx, refreshToken.TokenUUID.String(), user.ID, time.Unix(refreshToken.Expires.Unix(), 0).Sub(now)).Err()
	if err != nil {
		return nil, err
	}

	// * setting the cookie for refresh token
	ctx.SetCookie(
		"refresh-token",
		*refreshToken.Token,
		int(refreshToken.Expires.Unix()),
		"/",
		s.cfg.API.Domain,
		false,
		true,
	)

	return tokenHeader, err
}
