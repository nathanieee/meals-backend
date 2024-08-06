package utfile

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"project-skbackend/packages/consttypes"
)

func GetExtensionType(ext string) string {
	images := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	video := map[string]bool{".mp4": true}
	sound := map[string]bool{".mp3": true}

	filetype := "others"

	if images[ext] {
		filetype = "image"
	} else if video[ext] {
		filetype = "video"
	} else if sound[ext] {
		filetype = "sound"
	}

	return filetype
}

func Base64ToMultipartFileHeader(base64Str string, filename string, filetype consttypes.FileType) (*multipart.FileHeader, error) {
	// * decode the FileBase64 string
	fileBytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}

	// * create a buffer to store the multipart form
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// * create a form field with the file
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	partHeader.Set("Content-Type", filetype.String())

	part, err := writer.CreatePart(partHeader)
	if err != nil {
		return nil, err
	}

	// * write the file content to the form field
	if _, err := io.Copy(part, bytes.NewReader(fileBytes)); err != nil {
		return nil, err
	}

	// * close the writer to finalize the multipart form
	if err := writer.Close(); err != nil {
		return nil, err
	}

	// * create a multipart reader to read the form
	reader := multipart.NewReader(&b, writer.Boundary())
	form, err := reader.ReadForm(int64(len(fileBytes)))
	if err != nil {
		return nil, err
	}

	// * get the file header from the form
	fileHeaders := form.File["file"]
	if len(fileHeaders) == 0 {
		return nil, fmt.Errorf("no file headers found")
	}

	return fileHeaders[0], nil
}
