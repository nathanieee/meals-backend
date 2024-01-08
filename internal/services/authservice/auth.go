package authservice

import (
	"encoding/json"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/userrepo"
	"project-skbackend/internal/services/mailservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/uttoken"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	AuthService struct {
		cfg      *configs.Config
		userrepo userrepo.IUserRepository
		mailsvc  mailservice.IMailService
	}

	IAuthService interface {
		Login(req requests.LoginRequest) (*responses.UserResponse, *uttoken.TokenHeader, error)
		Register(req requests.RegisterRequest) (*responses.UserResponse, *uttoken.TokenHeader, error)
		ForgotPassword(req requests.ForgotPasswordRequest) error
		ResetPassword(req requests.ResetPasswordRequest) error
		SendVerificationEmail(id uuid.UUID, token int) error
		VerifyToken(req requests.VerifyTokenRequest) error
		SendResetPasswordEmail(id uuid.UUID, token int) error
		RefreshAuthToken(token string) (*responses.UserResponse, *uttoken.TokenHeader, error)
	}
)

func NewAuthService(
	cfg *configs.Config,
	userrepo userrepo.IUserRepository,
	mailsvc mailservice.IMailService,
) *AuthService {
	return &AuthService{
		cfg:      cfg,
		userrepo: userrepo,
		mailsvc:  mailsvc,
	}
}

func (a *AuthService) Login(req requests.LoginRequest) (*responses.UserResponse, *uttoken.TokenHeader, error) {
	user, err := a.userrepo.FindByEmail(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, err
		}

		return nil, nil, err
	}

	err = verifyPassword(*user, req.Password)
	if err != nil {
		return nil, nil, err
	}

	tokenHeader, err := a.generateAuthTokens(user)
	if err != nil {
		return nil, nil, err
	}

	return user.ToResponse(), tokenHeader, nil
}

func (a *AuthService) Register(req requests.RegisterRequest) (*responses.UserResponse, *uttoken.TokenHeader, error) {
	var user *models.User
	req.Email = strings.ToLower(req.Email)

	user, err := a.userrepo.FindByEmail(req.Email)
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil, err
	}

	if user != nil {
		return nil, nil, err
	}

	user = &models.User{
		Email:    req.Email,
		Password: req.Password,
		Role:     consttypes.UR_USER,
	}

	user, err = a.userrepo.Create(*user)
	if err != nil {
		return nil, nil, err
	}

	marshaledUser, _ := json.Marshal(user)
	err = json.Unmarshal(marshaledUser, &user)
	if err != nil {
		return nil, nil, err
	}

	token, err := a.generateAuthTokens(user)
	if err != nil {
		return nil, nil, err
	}

	return user.ToResponse(), token, nil
}

func (a *AuthService) ForgotPassword(req requests.ForgotPasswordRequest) error {
	user, err := a.userrepo.FindByEmail(req.Email)
	if err != nil && err == gorm.ErrRecordNotFound {
		return err
	}

	token := uttoken.GenerateRandomToken()

	err = a.SendResetPasswordEmail(user.ID, token)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) ResetPassword(req requests.ResetPasswordRequest) error {
	user, err := a.userrepo.FindByEmail(req.Email)
	if err != nil && err == gorm.ErrRecordNotFound {
		return err
	}

	if req.Token != user.ResetPasswordToken {
		return err
	}

	if !time.Now().UTC().Before(user.ResetPasswordSentAt.Add(time.Minute * 5)) {
		return err
	}

	userUpdate := models.User{
		Password:           req.Password,
		ResetPasswordToken: 0,
	}

	_, err = a.userrepo.Update(userUpdate, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) SendResetPasswordEmail(id uuid.UUID, token int) error {
	user, err := a.userrepo.FindByID(id)
	if err != nil {
		return err
	}

	if time.Now().UTC().Before(user.ResetPasswordSentAt.Add(time.Minute * 5)) {
		return err
	}

	userUpdate := models.User{
		ResetPasswordToken:  token,
		ResetPasswordSentAt: time.Now().UTC(),
	}

	_, err = a.userrepo.Update(userUpdate, user.ID)
	if err != nil {
		return err
	}

	emreq := requests.SendEmailRequest{ // TODO - change this request accordingly
		Template: "email_verification.html",
		Subject:  "Reset Password",
		Email:    user.Email,
		Token:    token,
	}

	err = a.mailsvc.SendVerificationEmail(emreq)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) SendVerificationEmail(id uuid.UUID, token int) error {
	user, err := a.userrepo.FindByID(id)
	if err != nil {
		return err
	}

	if time.Now().UTC().Before(user.ConfirmationSentAt.Add(time.Minute * 5)) {
		return err
	}

	userUpdate := models.User{
		ConfirmationToken:  token,
		ConfirmationSentAt: time.Now().UTC(),
	}

	_, err = a.userrepo.Update(userUpdate, user.ID)
	if err != nil {
		return err
	}

	emreq := requests.SendEmailRequest{ // TODO - change this request accordingly
		Template: "email_verification.html",
		Subject:  "Reset Password",
		Email:    user.Email,
		Token:    token,
	}

	err = a.mailsvc.SendVerificationEmail(emreq)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) VerifyToken(req requests.VerifyTokenRequest) error {
	user, err := a.userrepo.FindByEmail(req.Email)
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

	userUpdate := models.User{
		ConfirmedAt: time.Now().UTC(),
	}

	_, err = a.userrepo.Update(userUpdate, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) RefreshAuthToken(refreshToken string) (*responses.UserResponse, *uttoken.TokenHeader, error) {
	parsedToken, err := uttoken.ParseToken(refreshToken, a.cfg.App.Secret)
	if err != nil {
		return nil, nil, err
	}

	user, err := a.userrepo.FindByID(parsedToken.User.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, err
		}

		return nil, nil, err
	}

	if time.Now().Unix() >= parsedToken.Expire {
		return nil, nil, err
	}

	tokenHeader, err := a.generateAuthTokens(user)
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

func (a *AuthService) generateAuthTokens(user *models.User) (*uttoken.TokenHeader, error) {
	refreshToken, err := uttoken.GenerateToken(user.ToResponse(), a.cfg.App.RefreshTokenLifespan, a.cfg.App.TokenLifespanDuration, a.cfg.App.Secret)
	if err != nil {
		return nil, err
	}

	token, err := uttoken.GenerateToken(user.ToResponse(), a.cfg.App.TokenLifespan, a.cfg.App.TokenLifespanDuration, a.cfg.App.Secret)
	if err != nil {
		return nil, err
	}

	tokenHeader := uttoken.TokenHeader{
		AuthToken:           token.Token,
		AuthTokenExpires:    token.Expires,
		RefreshToken:        refreshToken.Token,
		RefreshTokenExpires: refreshToken.Expires,
	}

	return &tokenHeader, err
}
