package templates

import "embed"

//go:embed *.html
var files embed.FS

//go:embed static
var staticFiles embed.FS

func GetFiles() embed.FS {
	return files
}

func GetStaticFiles() embed.FS {
	return staticFiles
}