package utfile

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net/textproto"
	"path/filepath"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utmath"
	"strings"
)

var (
	suffixes = [5]consttypes.FileSuffixSize{
		consttypes.FSS_B,
		consttypes.FSS_KB,
		consttypes.FSS_MB,
		consttypes.FSS_GB,
		consttypes.FSS_TB,
	}

	// * define the default value of limit for each extension type
	extensionLimits = map[string]FileSizeLimit{
		"image": {MaxSize: 0.6, MaxSizeSuffix: consttypes.FSS_MB},
		"video": {MaxSize: 2.0, MaxSizeSuffix: consttypes.FSS_GB},
		"audio": {MaxSize: 5.0, MaxSizeSuffix: consttypes.FSS_MB},
		// * dd more extensions and their corresponding limits here
	}

	// * set the default accepted file extension type
	validExtension = map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
)

type (
	ValidateFileOpts struct {
		AllowedExtensions []string
		ValidateFileSizeOpts
	}

	ValidateFileSizeOpts struct {
		MaxImageSize       float64
		MaxImageSizeSuffix consttypes.FileSuffixSize
	}

	FileSizeLimit struct {
		MaxSize       float64
		MaxSizeSuffix consttypes.FileSuffixSize
	}
)

func ValidateFile(file *multipart.FileHeader, opts *ValidateFileOpts) (*string, error) {
	// Set default value
	if opts == nil {
		opts = &ValidateFileOpts{
			AllowedExtensions:    nil,
			ValidateFileSizeOpts: ValidateFileSizeOpts{},
		}
	}

	// Validate file
	filename := file.Filename
	extension := filepath.Ext(filename)
	if err := ValidateExtension(extension, opts.AllowedExtensions); err != nil {
		return new(string), err
	}

	reader, err := ReadRequestFile(file)
	if err != nil {
		return new(string), err
	}

	err = GetReadableFileSize(float64(reader.Size()), extension, &opts.ValidateFileSizeOpts)
	if err != nil {
		return new(string), err
	}

	return &filename, nil
}

func GetFileExtension(fileHeader *multipart.FileHeader) string {
	// * extract the filename from the file header
	filename := fileHeader.Filename

	// * get the file extension using the filepath package
	extension := filepath.Ext(filename)

	// * return the extension
	return extension
}

func GetReadableFileSize(size float64, ext string, opts *ValidateFileSizeOpts) error {
	// * set default options if not provided
	opts = setDefaultValidateFileSizeOpts(opts)

	base := math.Log(size) / math.Log(1024)
	roundedSize := utmath.Round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	suffix := suffixes[int(math.Floor(base))]

	// * convert KB to MB for simplicity and consistency
	if suffix == consttypes.FSS_KB {
		roundedSize = math.Ceil((roundedSize/1000)*100) / 100
		suffix = consttypes.FSS_MB
	}

	// * validate image size of extension matches
	if limit, exists := extensionLimits[ext]; exists {
		if opts.MaxImageSizeSuffix == consttypes.FileSuffixSize(suffix) {
			if roundedSize > limit.MaxSize || suffix == consttypes.FSS_GB || suffix == consttypes.FSS_TB {
				return consttypes.ErrFileSizeTooBig(ext, limit.MaxSize, limit.MaxSizeSuffix.String())
			}
		}
	} else {
		return consttypes.ErrUnsupportedFileExtension(ext)
	}
	return nil
}

func setDefaultValidateFileSizeOpts(opts *ValidateFileSizeOpts) *ValidateFileSizeOpts {
	if opts == nil {
		return &ValidateFileSizeOpts{
			MaxImageSize:       0.6,
			MaxImageSizeSuffix: consttypes.FSS_MB,
		}
	}
	if opts.MaxImageSize == 0 {
		opts.MaxImageSize = 0.6
	}
	if opts.MaxImageSizeSuffix == "" {
		opts.MaxImageSizeSuffix = consttypes.FSS_MB
	}
	return opts
}

func ValidateExtension(ext string, allowedExtensions []string) error {
	// * if custom extensions are provided, override the default map
	if len(allowedExtensions) > 0 {
		validExtension = make(map[string]bool, len(allowedExtensions))
		for _, extension := range allowedExtensions {
			validExtension[extension] = true
		}
	}

	if validExtension[ext] {
		return nil
	}

	// * extract the keys from the validExtension map and convert them to a string slice
	keys := make([]string, 0, len(validExtension))
	for key := range validExtension {
		keys = append(keys, key)
	}

	return consttypes.ErrUnsupportedFileExtension(keys)
}

func ReadRequestFile(file *multipart.FileHeader) (*bytes.Reader, error) {
	ogFile, err := file.Open()
	if err != nil {
		return nil, consttypes.ErrFailedToOpenFile
	}

	fileBytes, err := io.ReadAll(ogFile)
	if err != nil {
		return nil, consttypes.ErrFailedToReadFile
	}

	fileReader := bytes.NewReader(fileBytes)

	return fileReader, nil
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

func MultipartFileHeaderToBase64(fileHeader *multipart.FileHeader) (*FileBase64, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// * encode file content to base64
	encodedFile := base64.StdEncoding.EncodeToString(fileBytes)

	// * extract file extension
	ext := strings.TrimPrefix(filepath.Ext(fileHeader.Filename), ".")

	//  * map file extension to FileType
	fileType, err := consttypes.MapFileExtensionToFileType(ext)
	if err != nil {
		return nil, fmt.Errorf("invalid file type: %w", err)
	}

	return &FileBase64{
		FileBase64Str: encodedFile,
		FileName:      fileHeader.Filename,
		FileType:      fileType,
	}, nil
}

// * used to replace a request and make it so its not
// * dependent on other packages, hence delcaring it here
type (
	FileBase64 struct {
		FileBase64Str string              `json:"file_base64" form:"file_base64" binding:"required"`
		FileName      string              `json:"file_name" form:"file_name" binding:"required"`
		FileType      consttypes.FileType `json:"file_type" form:"file_type" binding:"required"`
	}

	FileMultipart struct {
		File *multipart.FileHeader `json:"file" form:"file"`
	}
)

func NewFileUpload(file *multipart.FileHeader) *FileMultipart {
	return &FileMultipart{
		File: file,
	}
}
