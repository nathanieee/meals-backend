package authservice

import (
	"fmt"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/userrepo"
	"project-skbackend/internal/services/mailservice"
	"project-skbackend/packages/utils/utlogger"
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
		SendVerificationEmail(id uuid.UUID, token string) error
		SendResetPasswordEmail(id uuid.UUID, token string) error
		VerifyToken(req requests.VerifyToken) error
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
		utlogger.LogError(err)
		return nil, nil, err
	}

	err = verifyPassword(*user, req.Password)
	if err != nil {
		utlogger.LogError(err)
		return nil, nil, err
	}

	tokenHeader, err := s.generateAuthTokens(user, ctx)
	if err != nil {
		utlogger.LogError(err)
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
		utlogger.LogError(err)
		return err
	}

	err = s.SendResetPasswordEmail(user.ID, token)
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
		return err
	}

	if !time.Now().UTC().Before(user.ResetPasswordSentAt.Add(time.Minute * 5)) {
		return err
	}

	user.Password = req.Password
	user.ResetPasswordToken = ""

	_, err = s.userrepo.Update(*user)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) SendResetPasswordEmail(id uuid.UUID, token string) error {
	user, err := s.userrepo.FindByID(id)
	if err != nil {
		return err
	}

	if time.Now().UTC().Before(user.ResetPasswordSentAt.Add(time.Minute * 5)) {
		return err
	}

	user.ResetPasswordToken = token
	user.ResetPasswordSentAt = time.Now().UTC()

	_, err = s.userrepo.Update(*user)
	if err != nil {
		return err
	}

	// TODO - change this request to reset password and sent correct data.
	emreq := requests.SendEmail{
		Template: "email_verification.html",
		Subject:  "Reset Password",
		Email:    user.Email,
		Token:    token,
	}

	err = s.mailsvc.SendVerificationEmail(emreq)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) SendVerificationEmail(id uuid.UUID, token string) error {
	user, err := s.userrepo.FindByID(id)
	if err != nil {
		return err
	}

	if time.Now().UTC().Before(user.ConfirmationSentAt.Add(time.Minute * 5)) {
		return err
	}

	user.ConfirmationToken = token
	user.ConfirmationSentAt = time.Now().UTC()

	_, err = s.userrepo.Update(*user)
	if err != nil {
		return err
	}

	// TODO - change this request to send correct data.
	emreq := requests.SendEmail{
		Template: "email_verification.html",
		Subject:  "Reset Password",
		Email:    user.Email,
		Token:    token,
	}

	err = s.mailsvc.SendVerificationEmail(emreq)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) VerifyToken(req requests.VerifyToken) error {
	user, err := s.userrepo.FindByEmail(req.Email)
	if err != nil {
		return err
	}

	if !time.Now().UTC().Before(user.ConfirmationSentAt.Add(time.Minute * 5)) {
		return err
	}

	if !user.ConfirmedAt.Equal(time.Time{}) {
		return err
	}

	if req.Token != user.ConfirmationToken {
		return err
	}

	user.ConfirmedAt = time.Now().UTC()

	_, err = s.userrepo.Update(*user)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RefreshAuthToken(refreshToken string, ctx *gin.Context) (*responses.User, *uttoken.TokenHeader, error) {
	now := time.Now()
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

	if time.Now().Unix() >= parsedToken.Expires.Unix() {
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
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return err
		default:
			return err
		}
	}

	return err
}

func (s *AuthService) generateAuthTokens(user *models.User, ctx *gin.Context) (*uttoken.TokenHeader, error) {
	now := time.Now()
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

	tokenHeader := uttoken.TokenHeader{
		AccessToken:         *accessToken.Token,
		AccessTokenExpires:  *accessToken.Expires,
		RefreshToken:        *refreshToken.Token,
		RefreshTokenExpires: *refreshToken.Expires,
	}

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
		"refresh_token",
		*refreshToken.Token,
		int(refreshToken.Expires.Unix()),
		"/",
		s.cfg.API.Domain,
		false,
		true,
	)

	return &tokenHeader, err
}
