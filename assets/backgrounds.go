package assets

import (
	"embed"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"log"
)

//go:embed backgrounds/*
var backgroundsFS embed.FS

func ListBackgrounds() []string {
	filePatterns := []string{
		"*.jpeg",
		"*.jpg",
		"*.png",
	}

	files := []string{}
	for _, filePattern := range filePatterns {
		matches, err := fs.Glob(backgroundsFS, fmt.Sprintf("backgrounds/%s", filePattern))
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, matches...)
	}

	return files
}

func LoadBackground(filePath string) (img image.Image, ok bool) {
	imgFile, err := backgroundsFS.Open(filePath)
	if err != nil {
		return nil, false
	}

	defer imgFile.Close()

	img, _, err = image.Decode(imgFile)
	if err != nil {
		return nil, false
	}

	return img, true
}
