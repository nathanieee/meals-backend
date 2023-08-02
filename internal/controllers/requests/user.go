package requests

type (
	CreateUserRequest struct {
		Name     string `json:"name" binding:"required" example:"user name"`
		Email    string `json:"email" binding:"required" example:"email@email.com"`
		Password string `json:"password" binding:"required" example:"password123"`
	}
)
