package services

import (
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/packages/utils"
)

type (
	IUserService interface {
		CreateUser(req requests.CreateUserRequest) (*responses.UserResponse, error)
		GetUser(id uint) (*responses.UserResponse, error)
		GetUsers(paginationReq models.Pagination) (*models.Pagination, error)
		DeleteUser(id uint) error
	}

	IAuthService interface {
		Login(req requests.LoginRequest) (*responses.UserResponse, *utils.TokenHeader, error)
		Register(req requests.RegisterRequest) (*responses.UserResponse, *utils.TokenHeader, error)
		ForgotPassword(req requests.ForgotPasswordRequest) error
		ResetPassword(req requests.ResetPasswordRequest) error
		SendVerificationEmail(id uint, token int) error
		VerifyToken(req requests.VerifyTokenRequest) error
		SendResetPasswordEmail(id uint, token int) error
		RefreshAuthToken(token string) (*responses.UserResponse, *utils.TokenHeader, error)
	}

	ILevelService interface {
		CreateLevel(req requests.CreateLevelRequest) (*responses.LevelResponse, error)
		GetLevel(lid uint) (*responses.LevelResponse, error)
		GetLevels(p models.Pagination) (*models.Pagination, error)
		DeleteLevel(lid uint) error
	}

	IRoleService interface {
		CreateRole(req requests.CreateRoleRequest) (*responses.RoleResponse, error)
		GetRole(rid uint) (*responses.RoleResponse, error)
		GetRoles(p models.Pagination) (*models.Pagination, error)
		DeleteRole(rid uint) error
	}

	IMailService interface {
		SendVerificationEmail(req requests.SendEmailRequest) error
	}
)
