package requests

type (
	CreateImage struct {
		FileBase64 `binding:"dive"`
	}

	UpdateImage struct {
		FileBase64 `binding:"dive"`
	}
)
