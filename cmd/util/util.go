package util

import (
	"bytes"
	"fmt"
	"github.com/eyko139/immich-notifier/internal/models"
	"html/template"
	"net/http"
)

type Helper struct {
	TemplateCache map[string]*template.Template
}

func New(templateCache map[string]*template.Template) *Helper {
	return &Helper{
		TemplateCache: templateCache,
	}
}

func (h *Helper) Render(w http.ResponseWriter, template string, data *interface{}) {
	if ts, ok := h.TemplateCache[template]; !ok {
		panic("Could fetch template from cache")
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
	ts, err := template.ParseFiles(fmt.Sprintf("../../ui/html/singles/%s", templateName))
	if err != nil {
		panic("Error parsing partial")
	}
	ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		fmt.Println(err)
	}
}

func (h *Helper) NewTemplateData(albums []models.Album, apiKey string) *models.TemplateData {
	return &models.TemplateData{
		Albums: albums,
		ApiKey: apiKey,
	}
}
