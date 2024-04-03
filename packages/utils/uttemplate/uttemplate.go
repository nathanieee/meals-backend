package uttemplate

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"project-skbackend/packages/utils/utlogger"
)

func ParseTemplateDir(dir string, file string) (*template.Template, error) {
	var (
		paths []string
	)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (info.Name() == "base.html" || info.Name() == "styles.html" || info.Name() == file) {
			paths = append(paths, path)
		}
		return nil
	})

	utlogger.Info(fmt.Sprintf("Parsing template files: %s", paths))

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}
