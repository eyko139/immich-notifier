package models

import (
	"html/template"
	"path/filepath"
)

type TemplateData struct {
	Albums []Album
	ApiKey string
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("../../ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles("../../ui/html/base.html")

		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("../../ui/html/partials/*.html")
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
