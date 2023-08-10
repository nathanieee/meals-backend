package auth

import (
	"encoding/json"
	"fmt"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	cfg      *configs.Config
	userRepo repositories.IUserRepo
}

func NewAuthService(userRepo repositories.IUserRepo, cfg *configs.Config) *AuthService {
	return &AuthService{userRepo: userRepo, cfg: cfg}
}

func (a *AuthService) Login(req requests.LoginRequest) (*responses.UserResponse, *utils.TokenHeader, error) {
	user, err := a.userRepo.FindByEmail(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, utils.ErrUserNotFound
		}

		return nil, nil, err
	}

	err = verifyPassword(user, req.Password)
	if err != nil {
		return nil, nil, err
	}

	tokenHeader, err := a.generateAuthTokens(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokenHeader, nil
}

func (a *AuthService) Register(req requests.RegisterRequest) (*responses.UserResponse, *utils.TokenHeader, error) {
	var user *responses.UserResponse
	req.Email = strings.ToLower(req.Email)

	user, err := a.userRepo.FindByEmail(req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, nil, utils.ErrUserNotFound
	}

	if user != nil {
		return nil, nil, utils.ErrUserAlreadyExist
	}

	hashedPassword, err := utils.EncryptPassword(req.Password)
	if err != nil {
		return nil, nil, err
	}

	userCreate := &models.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: hashedPassword,
		RoleID:   uint(consttypes.USER),
	}

	userModel, err := a.userRepo.Store(userCreate)
	if err != nil {
		return nil, nil, err
	}

	marshaledUser, _ := json.Marshal(userModel)
	err = json.Unmarshal(marshaledUser, &user)
	if err != nil {
		return nil, nil, err
	}

	token, err := a.generateAuthTokens(user)
	if err != nil {
		return nil, nil, err
	}

	return user, token, nil
}

func (a *AuthService) ForgotPassword(req requests.ForgotPasswordRequest) error {
	user, err := a.userRepo.FindByEmail(req.Email)
	if err != nil && err == gorm.ErrRecordNotFound {
		return utils.ErrUserNotFound
	}

	token := utils.GenerateRandomStringToken(16)

	err = a.SendResetPasswordEmail(user.ID, token)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) ResetPassword(req requests.ResetPasswordRequest) error {
	user, err := a.userRepo.FindByEmail(req.Email)
	if err != nil && err == gorm.ErrRecordNotFound {
		return utils.ErrUserNotFound
	}

	if req.Token != user.ResetPasswordToken {
		return utils.ErrTokenMismatch
	}

	if !time.Now().UTC().Before(user.ResetPasswordSentAt.Add(time.Minute * 5)) {
		return utils.ErrTokenExpired
	}

	hashedPassword, err := utils.EncryptPassword(req.Password)
	if err != nil {
		return err
	}

	userUpdate := models.User{
		Password:           hashedPassword,
		ResetPasswordToken: "",
	}

	_, err = a.userRepo.Update(userUpdate, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) SendResetPasswordEmail(id uint, token string) error {
	user, err := a.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	if time.Now().UTC().Before(user.ResetPasswordSentAt.Add(time.Minute * 5)) {
		return utils.ErrSendEmailResetRequest
	}

	userUpdate := models.User{
		ResetPasswordToken:  token,
		ResetPasswordSentAt: time.Now().UTC(),
	}

	_, err = a.userRepo.Update(userUpdate, user.ID)
	if err != nil {
		return err
	}

	// TODO DELETE THIS IN THE FUTURE BECAUSE THE TOKEN WILL BE SENT TO THE EMAIL
	fmt.Println(token)

	// TODO ADD AN EMAIL SENDING SERVICE HERE

	return nil
}

func (a *AuthService) SendVerificationEmail(id uint, token int) error {
	user, err := a.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	if time.Now().UTC().Before(user.ConfirmationSentAt.Add(time.Minute * 5)) {
		return utils.ErrSendEmailVerificationRequest
	}

	userUpdate := models.User{
		ConfirmationToken:  token,
		ConfirmationSentAt: time.Now().UTC(),
	}

	_, err = a.userRepo.Update(userUpdate, user.ID)
	if err != nil {
		return err
	}

	// TODO DELETE THIS IN THE FUTURE BECAUSE THE TOKEN WILL BE SENT TO THE EMAIL
	fmt.Println(token)

	// TODO ADD AN EMAIL SENDING SERVICE HERE

	return nil
}

func (a *AuthService) VerifyToken(req requests.VerifyTokenRequest) error {
	user, err := a.userRepo.FindByEmail(req.Email)
	if err != nil {
		return err
	}

	if !time.Now().UTC().Before(user.ConfirmationSentAt.Add(time.Minute * 5)) {
		return utils.ErrTokenExpired
	}

	if !user.ConfirmedAt.Equal(time.Time{}) {
		return utils.ErrUserAlreadyConfirmed
	}

	if req.Token != user.ConfirmationToken {
		return utils.ErrTokenIsNotTheSame
	}

	userUpdate := models.User{
		ConfirmedAt: time.Now().UTC(),
	}

	_, err = a.userRepo.Update(userUpdate, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) RefreshAuthToken(refreshToken string) (*responses.UserResponse, *utils.TokenHeader, error) {
	parsedToken, err := utils.ParseToken(refreshToken, a.cfg.App.Secret)
	if err != nil {
		return nil, nil, err
	}

	user, err := a.userRepo.FindByID(parsedToken.User.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, utils.ErrUserNotFound
		}

		return nil, nil, err
	}

	if time.Now().Unix() >= parsedToken.Expire {
		return nil, nil, utils.ErrTokenExpired
	}

	tokenHeader, err := a.generateAuthTokens(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokenHeader, err
}

func verifyPassword(u *responses.UserResponse, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return utils.ErrIncorrectPassword
		default:
			return err
		}
	}

	return err
}

func (a *AuthService) generateAuthTokens(user *responses.UserResponse) (*utils.TokenHeader, error) {
	refreshToken, err := utils.GenerateToken(user, a.cfg.App.RefreshTokenLifespan, a.cfg.App.TokenLifespanDuration, a.cfg.App.Secret)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user, a.cfg.App.TokenLifespan, a.cfg.App.TokenLifespanDuration, a.cfg.App.Secret)
	if err != nil {
		return nil, err
	}

	tokenHeader := utils.TokenHeader{
		AuthToken:           token.Token,
		AuthTokenExpires:    token.Expires,
		RefreshToken:        refreshToken.Token,
		RefreshTokenExpires: refreshToken.Expires,
	}

	return &tokenHeader, err
}
