package requests

type (
	CreateLevelRequest struct {
		Name string `json:"name,omitempty" binding:"required, unique"`
	}
)
