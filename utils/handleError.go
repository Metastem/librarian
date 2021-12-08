package utils

import (
	"html/template"
	"net/http"

	"codeberg.org/librarian/librarian/templates"
)

func HandleError(w http.ResponseWriter, err error) {
	errorTemplate, _ := template.ParseFS(templates.GetFiles(), "error.html")
	errorTemplate.Execute(w, map[string]interface{}{
		"err": err,
	})
}