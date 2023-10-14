package requests

type (
	CreateUserRequest struct {
		Email    string `json:"email" binding:"required,email" example:"email@email.com"`
		Password string `json:"-" binding:"required" example:"password"`
	}
)
