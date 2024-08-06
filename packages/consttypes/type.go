package consttypes

type (
	ImageType          string
	PatronType         string
	OrganizationType   string
	ResponseStatusType string
	FileType           string
)

const (
	IT_PROFILE       ImageType = "Profile"
	IT_MEAL          ImageType = "Meal"
	IT_MEAL_CATEGORY ImageType = "Meal Category"

	PT_ORGANIZATION PatronType = "Organization"
	PT_PERSONAL     PatronType = "Personal"

	OT_NURSINGHOME OrganizationType = "Nursing Home"

	RST_SUCCESS ResponseStatusType = "success"
	RST_FAIL    ResponseStatusType = "fail"
	RST_ERROR   ResponseStatusType = "error"

	// * text file types
	FT_PLAIN_TEXT FileType = "text/plain"
	FT_HTML       FileType = "text/html"
	FT_CSS        FileType = "text/css"
	FT_JAVASCRIPT FileType = "text/javascript"

	// * image file types
	FT_JPEG FileType = "image/jpeg"
	FT_PNG  FileType = "image/png"
	FT_GIF  FileType = "image/gif"
	FT_BMP  FileType = "image/bmp"

	// * audio file types
	FT_MP3 FileType = "audio/mpeg"
	FT_WAV FileType = "audio/wav"
	FT_OGG FileType = "audio/ogg"

	// * video file types
	FT_MP4  FileType = "video/mp4"
	FT_WEBM FileType = "video/webm"
	FT_OGGV FileType = "video/ogg"

	// * application file types
	FT_JSON FileType = "application/json"
	FT_XML  FileType = "application/xml"
	FT_PDF  FileType = "application/pdf"
	FT_ZIP  FileType = "application/zip"

	// * other file types
	FT_MULTIPART_FORM FileType = "multipart/form-data"
)

func (enum ImageType) String() string {
	return string(enum)
}

func (enum PatronType) String() string {
	return string(enum)
}

func (enum OrganizationType) String() string {
	return string(enum)
}

func (enum ResponseStatusType) String() string {
	return string(enum)
}

func (enum FileType) String() string {
	return string(enum)
}
