package assets

import (
	"embed"
	"io/fs"
	"log"
	"math/rand"
	"strings"
)

//go:embed titles/*
var titlesFS embed.FS

func RandomTitle() []string {
	files, err := fs.Glob(titlesFS, "titles/*.txt")
	if err != nil {
		log.Fatal(err)
	}

	file := files[rand.Intn(len(files))]

	content, err := titlesFS.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(content), "\n")

	var result []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}

	return result
}
