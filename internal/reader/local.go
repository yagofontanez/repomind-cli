package reader

import (
	"os"
	"path/filepath"
	"strings"
)

var ignoredDirs = map[string]bool{
	"node_modules": true,
	".git":         true,
	"dist":         true,
	"build":        true,
	".next":        true,
	"vendor":       true,
	"__pycache__":  true,
}

var ignoredExts = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".svg":  true,
	".ico":  true,
	".lock": true,
	".sum":  true,
}

type File struct {
	Path    string
	Content string
}

func ReadLocal(root string) ([]File, error) {
	var files []File

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			if ignoredDirs[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ignoredExts[ext] {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		files = append(files, File{
			Path:    path,
			Content: string(content),
		})

		return nil
	})

	return files, err
}
