package consttypes

import (
	"strings"
)

type (
	FileSuffixSize string
	FileType       string
)

const (
	FSS_B  FileSuffixSize = "B"
	FSS_KB FileSuffixSize = "KB"
	FSS_MB FileSuffixSize = "MB"
	FSS_GB FileSuffixSize = "GB"
	FSS_TB FileSuffixSize = "TB"

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

func (enum FileSuffixSize) String() string {
	return string(enum)
}

func (enum FileType) String() string {
	return string(enum)
}

func MapFileExtensionToFileType(extension string) (FileType, error) {
	switch strings.ToLower(extension) {
	// * text file types
	case "txt":
		return FT_PLAIN_TEXT, nil
	case "html", "htm":
		return FT_HTML, nil
	case "css":
		return FT_CSS, nil
	case "js":
		return FT_JAVASCRIPT, nil

	// * image file types
	case "jpeg", "jpg":
		return FT_JPEG, nil
	case "png":
		return FT_PNG, nil
	case "gif":
		return FT_GIF, nil
	case "bmp":
		return FT_BMP, nil

	// * audio file types
	case "mp3":
		return FT_MP3, nil
	case "wav":
		return FT_WAV, nil
	case "ogg":
		return FT_OGG, nil

	// * video file types
	case "mp4":
		return FT_MP4, nil
	case "webm":
		return FT_WEBM, nil
	case "oggv":
		return FT_OGGV, nil

	// * application file types
	case "json":
		return FT_JSON, nil
	case "xml":
		return FT_XML, nil
	case "pdf":
		return FT_PDF, nil
	case "zip":
		return FT_ZIP, nil

	// * other file types
	case "multipart":
		return FT_MULTIPART_FORM, nil

	default:
		return "", ErrUnsupportedFileExtension(extension)
	}
}
