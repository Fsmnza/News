package main

import (
	"alexedwards.net/snippetbox/pkg/models"
	"html/template"
	"path/filepath"
)

type templateData struct {
	News                *models.News
	CurrentYear         int
	Flash               string
	Category            string
	IsAuthenticated     bool
	NewsArray           []*models.News
	UserRole            string
	UserArray           []*models.User
	CommentList         []*models.Comments
	AuthenticatedUserID int
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
