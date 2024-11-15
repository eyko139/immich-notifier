package models

import (
	"html/template"
	"os"
	"path/filepath"
)

type TemplateData struct {
	Albums []Album
	User   UserContext

}

func NewTemplateCache() (map[string]*template.Template, error) {
    cwd, err := os.Getwd() // Get the current working directory
    if err != nil {
        panic(err)
    }
    htmlPath := filepath.Join(cwd, "ui/html/")
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(htmlPath + "/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles(htmlPath + "/base.html")

		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(htmlPath + "/partials/*.html")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseFiles(page)

		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
