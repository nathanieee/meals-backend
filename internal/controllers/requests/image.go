package requests

type (
	CreateImage struct {
		FileMultipart `binding:"dive"`
		FileBase64    `binding:"dive"`
	}

	UpdateImage struct {
		FileMultipart `binding:"dive"`
		FileBase64    `binding:"dive"`
	}
)
