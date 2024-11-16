package util

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/eyko139/immich-notifier/internal/models"
)

type Helper struct {
	TemplateCache map[string]*template.Template
}

func New(templateCache map[string]*template.Template) *Helper {
	return &Helper{
		TemplateCache: templateCache,
	}
}

func (h *Helper) Render(w http.ResponseWriter, template string, data any) {
	if ts, ok := h.TemplateCache[template]; !ok {
		panic("Could not fetch template from cache")
	} else {

		// writing template to a buffer first catches runtime errors
		buf := new(bytes.Buffer)

		err := ts.ExecuteTemplate(buf, "base", data)
		if err != nil {
			panic(err)
		}
		buf.WriteTo(w)
	}
}

func (h *Helper) ReturnHtml(w http.ResponseWriter, templateName string, data any) {
	cwd, err := os.Getwd() 
	if err != nil {
		panic(err)
	}
	staticPath := filepath.Join(cwd, "ui/html/singles")

	ts, err := template.ParseFiles(staticPath + "/" + templateName)
	if err != nil {
		panic("Error parsing partial")
	}
	ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		fmt.Println(err)
	}
}

func (h *Helper) ReturnPlainHtml(w http.ResponseWriter, templateName string, data any) {
	cwd, err := os.Getwd() 
	if err != nil {
		panic(err)
	}
	staticPath := filepath.Join(cwd, "ui/html/singles")

	ts, err := template.ParseFiles(staticPath + "/" + templateName)
	if err != nil {
		panic("Error parsing partial")
	}
	ts.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

func (h *Helper) NewTemplateData(albums []models.Album, email, name string) *models.TemplateData {
	return &models.TemplateData{
		Albums: albums,
        User:   models.UserContext{Email: email, Name: name},
	}
}
